package payments

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sneat-co/paymentus/backend/payrail"
)

type failingMarker struct {
	hasPaidErr  error
	markPaidErr error
}

func (m *failingMarker) HasPaid(context.Context, string) (bool, error) {
	return false, m.hasPaidErr
}

func (m *failingMarker) MarkPaid(context.Context, string, string) error {
	return m.markPaidErr
}

type customProvider struct {
	charge payrail.Charge
	err    error
}

func (p *customProvider) CreateCharge(context.Context, payrail.ChargeRequest) (payrail.Charge, error) {
	return p.charge, p.err
}

func (p *customProvider) RefundCharge(context.Context, string, int64) error {
	return nil
}

func TestNewService_NilArguments(t *testing.T) {
	marker := &idempotentMarker{}
	provider := &payrail.MockProvider{}
	identity := fixedIdentity{uid: "u1"}

	if _, err := NewService(nil, provider, marker); err == nil {
		t.Fatal("expected error with nil identity")
	}
	if _, err := NewService(identity, nil, marker); err == nil {
		t.Fatal("expected error with nil provider")
	}
	if _, err := NewService(identity, provider, nil); err == nil {
		t.Fatal("expected error with nil marker")
	}
}

func TestCreateCheckout_MethodNotAllowed(t *testing.T) {
	s := newTestService(t, fixedIdentity{uid: "u1"}, &payrail.MockProvider{}, &idempotentMarker{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	s.CreateCheckoutHandler()(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestCreateCheckout_HasPaidError(t *testing.T) {
	marker := &failingMarker{hasPaidErr: errors.New("db down")}
	s := newTestService(t, fixedIdentity{uid: "u1"}, &payrail.MockProvider{}, marker)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"actionId":"a1"}`))
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()
	s.CreateCheckoutHandler()(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestCreateCheckout_InvalidAttributionIDs(t *testing.T) {
	s := newTestService(t, fixedIdentity{uid: "u1"}, &payrail.MockProvider{}, &idempotentMarker{})

	for name, body := range map[string]string{
		"invalid category": `{"categoryId":"bad category with spaces","actionId":"a1"}`,
		"invalid question": `{"questionId":"bad/question","actionId":"a1"}`,
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			req.Header.Set("Authorization", "Bearer token")
			w := httptest.NewRecorder()
			s.CreateCheckoutHandler()(w, req)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d", w.Code)
			}
		})
	}
}

func TestCreateCheckout_ProviderFailures(t *testing.T) {
	marker := &idempotentMarker{}

	t.Run("provider returns error", func(t *testing.T) {
		provider := &customProvider{err: errors.New("provider err")}
		s := newTestService(t, fixedIdentity{uid: "u1"}, provider, marker)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"actionId":"a1"}`))
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()
		s.CreateCheckoutHandler()(w, req)
		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected 503, got %d", w.Code)
		}
	})

	t.Run("empty charge ref", func(t *testing.T) {
		provider := &customProvider{charge: payrail.Charge{ChargeRef: ""}}
		s := newTestService(t, fixedIdentity{uid: "u1"}, provider, marker)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"actionId":"a1"}`))
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()
		s.CreateCheckoutHandler()(w, req)
		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected 503, got %d", w.Code)
		}
	})

	t.Run("settled but markPaid fails", func(t *testing.T) {
		provider := &customProvider{charge: payrail.Charge{ChargeRef: "ch_sync", Settled: true}}
		fMarker := &failingMarker{markPaidErr: errors.New("mark err")}
		s := newTestService(t, fixedIdentity{uid: "u1"}, provider, fMarker)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"actionId":"a1"}`))
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()
		s.CreateCheckoutHandler()(w, req)
		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected 503, got %d", w.Code)
		}
	})
}

func TestBearerToken_NoPrefix(t *testing.T) {
	if got := bearerToken("Basic xyz"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
	if got := bearerToken(""); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestSettlementHandler_UserMismatchAndMarkError(t *testing.T) {
	for name, mutate := range map[string]func(*payrail.SettlementEvent){
		"empty uid":          func(e *payrail.SettlementEvent) { delete(e.Metadata, metadataUserID) },
		"empty payer":        func(e *payrail.SettlementEvent) { delete(e.Metadata, metadataPayerRef) },
		"uid payer mismatch": func(e *payrail.SettlementEvent) { e.Metadata[metadataPayerRef] = "diff-payer" },
	} {
		t.Run(name, func(t *testing.T) {
			s := newTestService(t, fixedIdentity{uid: "u1"}, &payrail.MockProvider{}, &idempotentMarker{})
			e := paidEvent()
			mutate(&e)
			if err := s.SettlementHandler()(context.Background(), e); !errors.Is(err, ErrInvalidSettlement) {
				t.Fatalf("expected ErrInvalidSettlement, got %v", err)
			}
		})
	}

	t.Run("mark paid error returned", func(t *testing.T) {
		expectedErr := errors.New("db fail")
		fMarker := &failingMarker{markPaidErr: expectedErr}
		s := newTestService(t, fixedIdentity{uid: "u1"}, &payrail.MockProvider{}, fMarker)
		e := paidEvent()
		if err := s.SettlementHandler()(context.Background(), e); !errors.Is(err, expectedErr) {
			t.Fatalf("expected %v, got %v", expectedErr, err)
		}
	})
}

func TestCreateCheckout_InvalidJSONBody(t *testing.T) {
	s := newTestService(t, fixedIdentity{uid: "u1"}, &payrail.MockProvider{}, &idempotentMarker{})
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid-json{"))
	w := httptest.NewRecorder()
	s.CreateCheckoutHandler()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
