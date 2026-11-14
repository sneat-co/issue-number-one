package payments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/sneat-co/paymentus/backend/payrail"
	"github.com/strongo/logus"
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

// testPayerUID is distinctive so a log line that leaks it is detectable.
const testPayerUID = "payer-uid-7Hq2"

func paidEvent() payrail.SettlementEvent {
	return payrail.SettlementEvent{ChargeRef: "ch_1", RefundRef: "pi_1", Consumer: Consumer, Kind: payrail.SettlementPaid, Amount: 100, Currency: "eur", Metadata: map[string]string{metadataUserID: testPayerUID, metadataPayerRef: testPayerUID, metadataPaymentMode: requiredPaymentMode}}
}

// logusEntry is one entry dispatched through logus, formatted exactly as the
// host's Cloud Logging handler formats it (sneat-go-host logginghost: severity
// mapped one-to-one, payload = fmt.Sprintf(MessageFormat, MessageArgs...)).
type logusEntry struct {
	severity logus.Severity
	message  string
}

type logusCapture struct {
	mu      sync.Mutex
	entries []logusEntry
}

func (c *logusCapture) Log(_ context.Context, entry logus.LogEntry) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = append(c.entries, logusEntry{severity: entry.Severity, message: fmt.Sprintf(entry.MessageFormat, entry.MessageArgs...)})
	return nil
}

var (
	capturedLogus     = &logusCapture{}
	capturedLogusOnce sync.Once
)

// captureLogus registers a handler on the process-wide logus dispatcher, the
// same one the host attaches its Cloud Logging handler to, and returns a reader
// for the entries written since the call.
func captureLogus(t *testing.T) func() []logusEntry {
	t.Helper()
	capturedLogusOnce.Do(func() { logus.AddLogEntryHandler(capturedLogus) })
	capturedLogus.mu.Lock()
	capturedLogus.entries = nil
	capturedLogus.mu.Unlock()
	return func() []logusEntry {
		capturedLogus.mu.Lock()
		defer capturedLogus.mu.Unlock()
		return append([]logusEntry(nil), capturedLogus.entries...)
	}
}

func onlyEntry(t *testing.T, entries []logusEntry) logusEntry {
	t.Helper()
	if len(entries) != 1 {
		t.Fatalf("want exactly one log entry, got %d: %v", len(entries), entries)
	}
	return entries[0]
}

// assertNoUserIdentifiers fails if an entry carries the event's user or payer
// identifier, or names the metadata keys that hold them.
func assertNoUserIdentifiers(t *testing.T, entries []logusEntry, e payrail.SettlementEvent) {
	t.Helper()
	for _, entry := range entries {
		for _, key := range []string{metadataUserID, metadataPayerRef} {
			if strings.Contains(entry.message, key) {
				t.Fatalf("log names user metadata key %q: %s", key, entry.message)
			}
			if v := e.Metadata[key]; v != "" && strings.Contains(entry.message, v) {
				t.Fatalf("log carries user identifier %q: %s", v, entry.message)
			}
		}
	}
}

// A paid event that fails a check can never become valid by being resent, so it
// is acknowledged (nil: no provider retry storm on the shared webhook), grants
// nothing, and is logged through logus at ERROR severity with the charge
// reference and reason, where the host ships it to Cloud Logging.
func TestSettlementInvalidPaidIsAcknowledgedAndLoggedAtErrorViaLogus(t *testing.T) {
	for name, test := range map[string]struct {
		mutate func(*payrail.SettlementEvent)
		reason string
	}{
		"consumer":            {func(e *payrail.SettlementEvent) { e.Consumer = "wallet-topup" }, reasonConsumerMismatch},
		"amount":              {func(e *payrail.SettlementEvent) { e.Amount = 99 }, reasonAmountMismatch},
		"currency":            {func(e *payrail.SettlementEvent) { e.Currency = "USD" }, reasonCurrencyMismatch},
		"user":                {func(e *payrail.SettlementEvent) { e.Metadata[metadataPayerRef] = "another-user-9Zx" }, reasonPayerMismatch},
		"empty uid":           {func(e *payrail.SettlementEvent) { delete(e.Metadata, metadataUserID) }, reasonPayerMissing},
		"empty payer":         {func(e *payrail.SettlementEvent) { delete(e.Metadata, metadataPayerRef) }, reasonPayerMissing},
		"test replay":         {func(e *payrail.SettlementEvent) { e.Metadata[metadataPaymentMode] = "test" }, reasonModeNotLive},
		"amount and currency": {func(e *payrail.SettlementEvent) { e.Amount = 5000; e.Currency = "GBP" }, reasonAmountMismatch},
		// Defence in depth: the real router rejects an empty charge reference
		// before calling the handler, so this is reachable only by a direct call.
		"charge": {func(e *payrail.SettlementEvent) { e.ChargeRef = "" }, reasonChargeMissing},
	} {
		t.Run(name, func(t *testing.T) {
			marker := &idempotentMarker{}
			s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
			logs := captureLogus(t)
			e := paidEvent()
			test.mutate(&e)
			if err := s.SettlementHandler()(context.Background(), e); err != nil {
				t.Fatalf("invalid paid event must be acknowledged, got err=%v", err)
			}
			if marker.calls != 0 || len(marker.grants) != 0 {
				t.Fatalf("invalid paid event unlocked eligibility: calls=%d grants=%v", marker.calls, marker.grants)
			}
			entries := logs()
			entry := onlyEntry(t, entries)
			if entry.severity != logus.SeverityError {
				t.Fatalf("severity=%d, want ERROR: %s", entry.severity, entry.message)
			}
			for _, want := range []string{
				msgPaidNotApplied,
				fmt.Sprintf("reason=%q", test.reason),
				fmt.Sprintf("chargeRef=%q", e.ChargeRef),
				fmt.Sprintf("kind=%q", payrail.SettlementPaid),
				fmt.Sprintf("amount=%d", e.Amount),
				fmt.Sprintf("currency=%q", e.Currency),
				ErrInvalidSettlement.Error(),
			} {
				if !strings.Contains(entry.message, want) {
					t.Fatalf("log lacks %q: %s", want, entry.message)
				}
			}
			assertNoUserIdentifiers(t, entries, e)
		})
	}
}

// A refund is a terminal provider fact the consumer cannot change by retrying.
// The product has not decided whether a refund revokes paid eligibility, so the
// handler acknowledges it, leaves eligibility exactly as it was, and records it.
func TestSettlementRefundIsAcknowledgedRecordedAndChangesNothing(t *testing.T) {
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
	handler := s.SettlementHandler()
	if err := handler(context.Background(), paidEvent()); err != nil {
		t.Fatal(err)
	}
	logs := captureLogus(t)
	refund := paidEvent()
	refund.Kind = payrail.SettlementRefunded
	refund.RefundRef = "re_1"
	for i := 0; i < 2; i++ { // the provider re-lists every refund on re-delivery
		if err := handler(context.Background(), refund); err != nil {
			t.Fatalf("refund must be acknowledged, got err=%v", err)
		}
	}
	if marker.calls != 1 || len(marker.grants) != 1 || marker.grants["ch_1"] != testPayerUID {
		t.Fatalf("refund changed eligibility: calls=%d grants=%v", marker.calls, marker.grants)
	}
	entries := logs()
	if len(entries) != 2 {
		t.Fatalf("entries=%v", entries)
	}
	for _, entry := range entries {
		if entry.severity != logus.SeverityWarning || !strings.Contains(entry.message, msgRefundAcknowledged) ||
			!strings.Contains(entry.message, `kind="refunded"`) || !strings.Contains(entry.message, `chargeRef="ch_1"`) || !strings.Contains(entry.message, `refundRef="re_1"`) {
			t.Fatalf("entry=%v", entry)
		}
	}
	assertNoUserIdentifiers(t, entries, refund)
}

func TestSettlementRefundWithoutPriorPaymentGrantsNothing(t *testing.T) {
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
	logs := captureLogus(t)
	refund := paidEvent()
	refund.Kind = payrail.SettlementRefunded
	if err := s.SettlementHandler()(context.Background(), refund); err != nil {
		t.Fatalf("err=%v", err)
	}
	if marker.calls != 0 {
		t.Fatalf("refund touched eligibility: calls=%d", marker.calls)
	}
	onlyEntry(t, logs())
}

func TestSettlementFailedAndUnknownKindsAreAcknowledged(t *testing.T) {
	for name, test := range map[string]struct {
		kind     payrail.SettlementKind
		severity logus.Severity
		msg      string
	}{
		"failed":  {payrail.SettlementFailed, logus.SeverityInfo, msgFailedAcknowledged},
		"unknown": {payrail.SettlementKind("disputed"), logus.SeverityWarning, msgUnknownKindAcknowledged},
		"empty":   {payrail.SettlementKind(""), logus.SeverityWarning, msgUnknownKindAcknowledged},
	} {
		t.Run(name, func(t *testing.T) {
			marker := &idempotentMarker{}
			s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
			logs := captureLogus(t)
			e := paidEvent()
			e.Kind = test.kind
			if err := s.SettlementHandler()(context.Background(), e); err != nil {
				t.Fatalf("%s must be acknowledged, got err=%v", name, err)
			}
			if marker.calls != 0 {
				t.Fatalf("%s touched eligibility", name)
			}
			entries := logs()
			entry := onlyEntry(t, entries)
			if entry.severity != test.severity || !strings.Contains(entry.message, test.msg) || !strings.Contains(entry.message, fmt.Sprintf("kind=%q", test.kind)) {
				t.Fatalf("entry=%v", entry)
			}
			assertNoUserIdentifiers(t, entries, e)
		})
	}
}

// Store errors keep the error, so the host answers 5xx and the provider
// resends; the next delivery is deduplicated by charge ID.
func TestSettlementTransientStoreErrorIsRetried(t *testing.T) {
	expected := errors.New("firestore unavailable")
	s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, &failingMarker{markPaidErr: expected})
	logs := captureLogus(t)
	if err := s.SettlementHandler()(context.Background(), paidEvent()); !errors.Is(err, expected) {
		t.Fatalf("transient store error must be returned for retry, got %v", err)
	}
	if entries := logs(); len(entries) != 0 {
		t.Fatalf("a retried event must not be recorded as acknowledged: %v", entries)
	}
}

func TestSettlementDuplicateIsIdempotentByProviderCharge(t *testing.T) {
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
	logs := captureLogus(t)
	e := paidEvent()
	if err := s.SettlementHandler()(context.Background(), e); err != nil {
		t.Fatal(err)
	}
	if err := s.SettlementHandler()(context.Background(), e); err != nil {
		t.Fatal(err)
	}
	if len(marker.grants) != 1 || marker.grants["ch_1"] != testPayerUID {
		t.Fatalf("grants=%v", marker.grants)
	}
	if entries := logs(); len(entries) != 0 {
		t.Fatalf("a valid paid event is not an operator event: %v", entries)
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
