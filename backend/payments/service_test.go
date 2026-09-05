package payments

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/sneat-co/paymentus/backend/payrail"
)

type fixedIdentity struct{ uid string }

func (i fixedIdentity) UserID(context.Context, string) (string, bool) { return i.uid, i.uid != "" }

type idempotentMarker struct {
	mu     sync.Mutex
	grants map[string]string
	calls  int
}

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

func newTestService(t *testing.T, identity Identity, provider payrail.PaymentProvider, marker EligibilityMarker) *Service {
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
	if got.IdempotencyKey != "issuenumber-verification:user-1" {
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
		"anonymous": {fixedIdentity{}, `{}`, http.StatusUnauthorized},
		"free form": {fixedIdentity{uid: "u"}, `{"questionId":"housing is expensive"}`, http.StatusBadRequest},
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
	return payrail.SettlementEvent{ChargeRef: "ch_1", RefundRef: "pi_1", Consumer: Consumer, Kind: payrail.SettlementPaid, Amount: 100, Currency: "eur", Metadata: map[string]string{metadataUserID: "u1", metadataPayerRef: "u1"}}
}

func TestSettlementRejectsForgedConsumerAmountAndUser(t *testing.T) {
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
	for name, mutate := range map[string]func(*payrail.SettlementEvent){
		"consumer": func(e *payrail.SettlementEvent) { e.Consumer = "wallet-topup" },
		"amount":   func(e *payrail.SettlementEvent) { e.Amount = 99 },
		"currency": func(e *payrail.SettlementEvent) { e.Currency = "USD" },
		"outcome":  func(e *payrail.SettlementEvent) { e.Kind = payrail.SettlementRefunded },
		"user":     func(e *payrail.SettlementEvent) { e.Metadata[metadataPayerRef] = "another-user" },
	} {
		t.Run(name, func(t *testing.T) {
			e := paidEvent()
			mutate(&e)
			if err := s.SettlementHandler()(context.Background(), e); !errors.Is(err, ErrInvalidSettlement) {
				t.Fatalf("err=%v", err)
			}
		})
	}
	if len(marker.grants) != 0 {
		t.Fatalf("forged settlement granted: %v", marker.grants)
	}
}

func TestSettlementDuplicateIsIdempotentByProviderCharge(t *testing.T) {
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "u"}, &payrail.MockProvider{}, marker)
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
}

func TestSynchronousSettlementMarksPaid(t *testing.T) {
	provider := &payrail.MockProvider{SettleSync: true}
	marker := &idempotentMarker{}
	s := newTestService(t, fixedIdentity{uid: "u1"}, provider, marker)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()
	s.CreateCheckoutHandler()(w, req)
	if w.Code != http.StatusCreated || len(marker.grants) != 1 {
		t.Fatalf("status=%d grants=%v", w.Code, marker.grants)
	}
}
