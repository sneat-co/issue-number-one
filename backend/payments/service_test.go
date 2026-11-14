package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/sneat-co/paymentus/backend/payrail"
)

type fixedIdentity struct{ uid string }

func (i fixedIdentity) UserID(context.Context, string) (string, bool) { return i.uid, i.uid != "" }

type idempotentProvider struct {
	created map[string]payrail.ChargeRequest
	charges map[string]payrail.Charge
}

func (p *idempotentProvider) CreateCharge(_ context.Context, request payrail.ChargeRequest) (payrail.Charge, error) {
	if p.created == nil {
		p.created = map[string]payrail.ChargeRequest{}
		p.charges = map[string]payrail.Charge{}
	}
	if charge, ok := p.charges[request.IdempotencyKey]; ok {
		previous := p.created[request.IdempotencyKey]
		if previous.Metadata[metadataCategoryID] != request.Metadata[metadataCategoryID] || previous.Metadata[metadataQuestionID] != request.Metadata[metadataQuestionID] || previous.Metadata[metadataActionID] != request.Metadata[metadataActionID] {
			return payrail.Charge{}, errors.New("idempotency metadata conflict")
		}
		return charge, nil
	}
	charge := payrail.Charge{ChargeRef: "ch_" + request.Metadata[metadataActionID], CheckoutURL: "https://checkout.example/" + request.Metadata[metadataActionID]}
	p.created[request.IdempotencyKey] = request
	p.charges[request.IdempotencyKey] = charge
	return charge, nil
}

func (*idempotentProvider) RefundCharge(context.Context, string, int64) error { return nil }

type idempotentMarker struct {
	mu     sync.Mutex
	grants map[string]string
	calls  int
	paid   bool
}

func (m *idempotentMarker) HasPaid(context.Context, string) (bool, error) { return m.paid, nil }

func (m *idempotentMarker) MarkPaid(_ context.Context, uid, chargeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls++
	if m.grants == nil {
		m.grants = map[string]string{}
	}
	if previous, ok := m.grants[chargeID]; ok && previous != uid {
		return errors.New("charge belongs to another user")
	}
	m.grants[chargeID] = uid
	return nil
}

func newTestService(t *testing.T, identity Identity, provider payrail.PaymentProvider, marker EligibilityStore) *Service {
	t.Helper()
	s, err := NewService(identity, provider, marker)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestCreateCheckoutIsAuthenticatedFixedAndAttributed(t *testing.T) {
	provider := &payrail.MockProvider{CheckoutURL: "https://checkout.example/session"}
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "user-1"}, provider, marker)
	req := httptest.NewRequest(http.MethodPost, "/v0/issuenumber/verification/payment", strings.NewReader(`{"categoryId":"country","questionId":"ireland","actionId":"attempt-7"}`))
	req.Header.Set("Authorization", "Bearer valid")
	w := httptest.NewRecorder()
	s.CreateCheckoutHandler()(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if len(provider.Charges) != 1 {
		t.Fatalf("charges=%d", len(provider.Charges))
	}
	got := provider.Charges[0]
	if got.Amount != 100 || got.Currency != "EUR" || got.Consumer != Consumer || got.PayerRef != "user-1" {
		t.Fatalf("unexpected charge: %+v", got)
	}
	if got.IdempotencyKey != "issuenumber-verification:user-1:attempt-7" {
		t.Fatalf("idempotency=%q", got.IdempotencyKey)
	}
	if got.Metadata[metadataUserID] != "user-1" || got.Metadata[metadataQuestionID] != "ireland" {
		t.Fatalf("metadata=%v", got.Metadata)
	}
	var out CreateCheckoutResponse
	if err := json.NewDecoder(w.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out.CheckoutURL == "" || out.ChargeID == "" {
		t.Fatalf("response=%+v", out)
	}
}

func TestCreateCheckoutRejectsAnonymousAndUnsafeAttribution(t *testing.T) {
	provider := &payrail.MockProvider{}
	marker := &idempotentMarker{}
	for name, test := range map[string]struct {
		identity Identity
		body     string
		want     int
	}{
		"anonymous": {fixedIdentity{}, `{"actionId":"attempt-1"}`, http.StatusUnauthorized},
		"no action": {fixedIdentity{uid: "u"}, `{}`, http.StatusBadRequest},
		"free form": {fixedIdentity{uid: "u"}, `{"questionId":"housing is expensive","actionId":"attempt-1"}`, http.StatusBadRequest},
	} {
		t.Run(name, func(t *testing.T) {
			s := newTestService(t, test.identity, provider, marker)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			req.Header.Set("Authorization", "Bearer token")
			w := httptest.NewRecorder()
			s.CreateCheckoutHandler()(w, req)
			if w.Code != test.want {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
		})
	}
	if len(provider.Charges) != 0 {
		t.Fatal("rejected requests must not reach provider")
	}
}

func paidEvent() payrail.SettlementEvent {
	return payrail.SettlementEvent{ChargeRef: "ch_1", RefundRef: "pi_1", Consumer: Consumer, Kind: payrail.SettlementPaid, Amount: 100, Currency: "eur", Metadata: map[string]string{metadataUserID: "u1", metadataPayerRef: "u1", metadataPaymentMode: requiredPaymentMode}}
}

// logRecord is one captured structured log line.
type logRecord map[string]any

// captureLogs points the service at an in-memory JSON logger and returns a
// reader for the records written so far.
func captureLogs(t *testing.T, s *Service) func() []logRecord {
	t.Helper()
	var buf bytes.Buffer
	s.log = slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return func() []logRecord {
		t.Helper()
		var records []logRecord
		for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
			if line == "" {
				continue
			}
			var r logRecord
			if err := json.Unmarshal([]byte(line), &r); err != nil {
				t.Fatalf("log line %q: %v", line, err)
			}
			records = append(records, r)
		}
		return records
	}
}

func onlyRecord(t *testing.T, records []logRecord) logRecord {
	t.Helper()
	if len(records) != 1 {
		t.Fatalf("want exactly one log record, got %d: %v", len(records), records)
	}
	return records[0]
}

// A paid event that fails a check can never become valid by being resent, so it
// is acknowledged (nil: no provider retry storm on the shared webhook), grants
// nothing, and is logged at ERROR with the reason for an operator.
func TestSettlementInvalidPaidIsAcknowledgedRecordedAndNotGranted(t *testing.T) {
	for name, test := range map[string]struct {
		mutate func(*payrail.SettlementEvent)
		reason string
	}{
		"consumer":            {func(e *payrail.SettlementEvent) { e.Consumer = "wallet-topup" }, reasonConsumerMismatch},
		"amount":              {func(e *payrail.SettlementEvent) { e.Amount = 99 }, reasonAmountMismatch},
		"currency":            {func(e *payrail.SettlementEvent) { e.Currency = "USD" }, reasonCurrencyMismatch},
		"user":                {func(e *payrail.SettlementEvent) { e.Metadata[metadataPayerRef] = "another-user" }, reasonPayerMismatch},
		"empty uid":           {func(e *payrail.SettlementEvent) { delete(e.Metadata, metadataUserID) }, reasonPayerMissing},
		"empty payer":         {func(e *payrail.SettlementEvent) { delete(e.Metadata, metadataPayerRef) }, reasonPayerMissing},
		"charge":              {func(e *payrail.SettlementEvent) { e.ChargeRef = "" }, reasonChargeMissing},
		"test replay":         {func(e *payrail.SettlementEvent) { e.Metadata[metadataPaymentMode] = "test" }, reasonModeNotLive},
		"amount and currency": {func(e *payrail.SettlementEvent) { e.Amount = 5000; e.Currency = "GBP" }, reasonAmountMismatch},
	} {
		t.Run(name, func(t *testing.T) {
			marker := &idempotentMarker{}
			s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
			logs := captureLogs(t, s)
			e := paidEvent()
			test.mutate(&e)
			if err := s.SettlementHandler()(context.Background(), e); err != nil {
				t.Fatalf("invalid paid event must be acknowledged, got err=%v", err)
			}
			if marker.calls != 0 || len(marker.grants) != 0 {
				t.Fatalf("invalid paid event unlocked eligibility: calls=%d grants=%v", marker.calls, marker.grants)
			}
			r := onlyRecord(t, logs())
			if r["level"] != "ERROR" || r["msg"] != msgPaidNotApplied || r["reason"] != test.reason {
				t.Fatalf("record=%v", r)
			}
			if r["kind"] != string(payrail.SettlementPaid) || r["chargeRef"] != e.ChargeRef || r["amount"] != float64(e.Amount) || r["currency"] != e.Currency {
				t.Fatalf("record lacks event facts: %v", r)
			}
			if !strings.Contains(r["err"].(string), ErrInvalidSettlement.Error()) {
				t.Fatalf("record err=%v", r["err"])
			}
		})
	}
}

// A refund is a terminal provider fact the consumer cannot change by retrying.
// The product has not decided whether a refund revokes paid eligibility, so the
// handler acknowledges it, leaves eligibility exactly as it was, and records it.
func TestSettlementRefundIsAcknowledgedRecordedAndChangesNothing(t *testing.T) {
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
	logs := captureLogs(t, s)
	handler := s.SettlementHandler()
	if err := handler(context.Background(), paidEvent()); err != nil {
		t.Fatal(err)
	}
	refund := paidEvent()
	refund.Kind = payrail.SettlementRefunded
	refund.RefundRef = "re_1"
	for i := 0; i < 2; i++ { // the provider re-lists every refund on re-delivery
		if err := handler(context.Background(), refund); err != nil {
			t.Fatalf("refund must be acknowledged, got err=%v", err)
		}
	}
	if marker.calls != 1 || len(marker.grants) != 1 || marker.grants["ch_1"] != "u1" {
		t.Fatalf("refund changed eligibility: calls=%d grants=%v", marker.calls, marker.grants)
	}
	records := logs()
	if len(records) != 2 {
		t.Fatalf("records=%v", records)
	}
	for _, r := range records {
		if r["level"] != "WARN" || r["msg"] != msgRefundAcknowledged || r["kind"] != string(payrail.SettlementRefunded) || r["chargeRef"] != "ch_1" || r["refundRef"] != "re_1" {
			t.Fatalf("record=%v", r)
		}
	}
}

func TestSettlementRefundWithoutPriorPaymentGrantsNothing(t *testing.T) {
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
	logs := captureLogs(t, s)
	refund := paidEvent()
	refund.Kind = payrail.SettlementRefunded
	if err := s.SettlementHandler()(context.Background(), refund); err != nil {
		t.Fatalf("err=%v", err)
	}
	if marker.calls != 0 {
		t.Fatalf("refund touched eligibility: calls=%d", marker.calls)
	}
	onlyRecord(t, logs())
}

func TestSettlementFailedAndUnknownKindsAreAcknowledged(t *testing.T) {
	for name, test := range map[string]struct {
		kind  payrail.SettlementKind
		level string
		msg   string
	}{
		"failed":  {payrail.SettlementFailed, "INFO", msgFailedAcknowledged},
		"unknown": {payrail.SettlementKind("disputed"), "WARN", msgUnknownKindAcknowledged},
		"empty":   {payrail.SettlementKind(""), "WARN", msgUnknownKindAcknowledged},
	} {
		t.Run(name, func(t *testing.T) {
			marker := &idempotentMarker{}
			s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
			logs := captureLogs(t, s)
			e := paidEvent()
			e.Kind = test.kind
			if err := s.SettlementHandler()(context.Background(), e); err != nil {
				t.Fatalf("%s must be acknowledged, got err=%v", name, err)
			}
			if marker.calls != 0 {
				t.Fatalf("%s touched eligibility", name)
			}
			r := onlyRecord(t, logs())
			if r["level"] != test.level || r["msg"] != test.msg || r["kind"] != string(test.kind) {
				t.Fatalf("record=%v", r)
			}
		})
	}
}

// Only a transient failure keeps the error, so the host answers 5xx and the
// provider resends; the next delivery is deduplicated by charge ID.
func TestSettlementTransientStoreErrorIsRetried(t *testing.T) {
	expected := errors.New("firestore unavailable")
	s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, &failingMarker{markPaidErr: expected})
	logs := captureLogs(t, s)
	if err := s.SettlementHandler()(context.Background(), paidEvent()); !errors.Is(err, expected) {
		t.Fatalf("transient store error must be returned for retry, got %v", err)
	}
	if records := logs(); len(records) != 0 {
		t.Fatalf("a retried event must not be recorded as acknowledged: %v", records)
	}
}

func TestServiceLoggerDefaultsToProcessLogger(t *testing.T) {
	if got := (&Service{}).logger(); got != slog.Default() {
		t.Fatalf("logger=%v", got)
	}
}

func TestSettlementDuplicateIsIdempotentByProviderCharge(t *testing.T) {
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
	logs := captureLogs(t, s)
	e := paidEvent()
	if err := s.SettlementHandler()(context.Background(), e); err != nil {
		t.Fatal(err)
	}
	if err := s.SettlementHandler()(context.Background(), e); err != nil {
		t.Fatal(err)
	}
	if len(marker.grants) != 1 || marker.grants["ch_1"] != "u1" {
		t.Fatalf("grants=%v", marker.grants)
	}
	if records := logs(); len(records) != 0 {
		t.Fatalf("a valid paid event is not an operator event: %v", records)
	}
}

func TestSynchronousSettlementMarksPaid(t *testing.T) {
	provider := &payrail.MockProvider{SettleSync: true}
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "u1"}, provider, marker)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"actionId":"attempt-sync"}`))
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()
	s.CreateCheckoutHandler()(w, req)
	if w.Code != http.StatusCreated || len(marker.grants) != 1 {
		t.Fatalf("status=%d grants=%v", w.Code, marker.grants)
	}
}

func TestAlreadyPaidReturnsSettledWithoutCreatingCharge(t *testing.T) {
	provider := &payrail.MockProvider{}
	marker := &idempotentMarker{paid: true}
	s := newTestService(t, fixedIdentity{uid: "u1"}, provider, marker)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"actionId":"attempt-1"}`))
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()
	s.CreateCheckoutHandler()(w, req)
	if w.Code != http.StatusOK || len(provider.Charges) != 0 {
		t.Fatalf("status=%d charges=%d body=%s", w.Code, len(provider.Charges), w.Body.String())
	}
	var response CreateCheckoutResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil || !response.Settled || response.ChargeID != "" {
		t.Fatalf("response=%+v err=%v", response, err)
	}
}

func TestCheckoutRetryUsesAttemptIdempotencyAndNewAttemptCanCreate(t *testing.T) {
	provider := &idempotentProvider{}
	s := newTestService(t, fixedIdentity{uid: "u1"}, provider, &idempotentMarker{})
	create := func(body string) CreateCheckoutResponse {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()
		s.CreateCheckoutHandler()(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
		}
		var response CreateCheckoutResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}
		return response
	}
	body := `{"categoryId":"country","questionId":"ireland","actionId":"attempt-1"}`
	first, retry := create(body), create(body)
	if first != retry || len(provider.created) != 1 {
		t.Fatalf("first=%+v retry=%+v provider creations=%d", first, retry, len(provider.created))
	}
	secondAttempt := create(`{"categoryId":"country","questionId":"ireland","actionId":"attempt-2"}`)
	if secondAttempt.ChargeID == first.ChargeID || len(provider.created) != 2 {
		t.Fatalf("first=%+v second=%+v provider creations=%d", first, secondAttempt, len(provider.created))
	}
}
