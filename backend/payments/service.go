// Package payments owns IssueNumber.one's payment-verification policy while
// delegating money movement and signed webhook handling to paymentus/payrail.
package payments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/sneat-co/paymentus/backend/payrail"
)

const (
	Consumer            = "issuenumber-payment-verification"
	AmountMinorUnits    = int64(100)
	Currency            = "EUR"
	MerchantDescription = "Support IssueNumber.one and verify participation"

	metadataUserID      = "issueNumberUserID"
	metadataPayerRef    = "sneatPayerRef"
	metadataCategoryID  = "categoryID"
	metadataQuestionID  = "questionID"
	metadataActionID    = "actionID"
	metadataPaymentMode = "sneatPaymentMode"
	requiredPaymentMode = "live"
)

var (
	ErrUnauthorized      = errors.New("issuenumber payments: authentication required")
	ErrUnavailable       = errors.New("issuenumber payments: payment verification unavailable")
	ErrInvalidRequest    = errors.New("issuenumber payments: invalid request")
	ErrInvalidSettlement = errors.New("issuenumber payments: invalid settlement")
	// Attribution values are identifiers only. This deliberately excludes spaces,
	// URLs and free-form issue text from payment metadata and analytics joins.
	safeAttributionID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
)

// Identity verifies a Firebase bearer token and returns its trusted user ID.
type Identity interface {
	UserID(ctx context.Context, bearerToken string) (string, bool)
}

// EligibilityStore reads and durably records paid verification.
// Repeating the same provider charge ID must not grant eligibility twice.
type EligibilityStore interface {
	HasPaid(ctx context.Context, userID string) (bool, error)
	MarkPaid(ctx context.Context, userID, chargeID string) error
}

// Service creates a fixed-price checkout and consumes verified settlements.
type Service struct {
	identity Identity
	provider payrail.PaymentProvider
	marker   EligibilityStore
}

func NewService(identity Identity, provider payrail.PaymentProvider, marker EligibilityStore) (*Service, error) {
	if identity == nil || provider == nil || marker == nil {
		return nil, fmt.Errorf("%w: identity, provider and marker are required", ErrUnavailable)
	}
	return &Service{identity: identity, provider: provider, marker: marker}, nil
}

type CreateCheckoutRequest struct {
	CategoryID string `json:"categoryId,omitempty"`
	QuestionID string `json:"questionId,omitempty"`
	ActionID   string `json:"actionId,omitempty"`
}

type CreateCheckoutResponse struct {
	ChargeID    string `json:"chargeId,omitempty"`
	CheckoutURL string `json:"checkoutUrl,omitempty"`
	Settled     bool   `json:"settled"`
}

// CreateCheckoutHandler returns the authenticated POST endpoint. Amount,
// currency, consumer and merchant copy are server-owned and cannot be supplied
// by the browser.
func (s *Service) CreateCheckoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		uid, ok := s.identity.UserID(r.Context(), bearerToken(r.Header.Get("Authorization")))
		if !ok || uid == "" {
			http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}
		paid, err := s.marker.HasPaid(r.Context(), uid)
		if err != nil {
			http.Error(w, ErrUnavailable.Error(), http.StatusServiceUnavailable)
			return
		}
		if paid {
			writeCheckoutResponse(w, http.StatusOK, CreateCheckoutResponse{Settled: true})
			return
		}
		var input CreateCheckoutRequest
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&input); err != nil {
			http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
			return
		}
		if !validOptionalID(input.CategoryID) || !validOptionalID(input.QuestionID) || !safeAttributionID.MatchString(input.ActionID) {
			http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
			return
		}
		metadata := map[string]string{metadataUserID: uid, metadataPayerRef: uid}
		putIfPresent(metadata, metadataCategoryID, input.CategoryID)
		putIfPresent(metadata, metadataQuestionID, input.QuestionID)
		putIfPresent(metadata, metadataActionID, input.ActionID)
		charge, err := s.provider.CreateCharge(r.Context(), payrail.ChargeRequest{
			Consumer: Consumer, PayerRef: uid, Amount: AmountMinorUnits,
			Currency: Currency, Description: MerchantDescription,
			Metadata: metadata, IdempotencyKey: "issuenumber-verification:" + uid + ":" + input.ActionID,
		})
		if err != nil {
			http.Error(w, ErrUnavailable.Error(), http.StatusServiceUnavailable)
			return
		}
		if charge.ChargeRef == "" {
			http.Error(w, ErrUnavailable.Error(), http.StatusServiceUnavailable)
			return
		}
		if charge.Settled {
			if err := s.marker.MarkPaid(r.Context(), uid, charge.ChargeRef); err != nil {
				http.Error(w, ErrUnavailable.Error(), http.StatusServiceUnavailable)
				return
			}
		}
		writeCheckoutResponse(w, http.StatusCreated, CreateCheckoutResponse{ChargeID: charge.ChargeRef, CheckoutURL: charge.CheckoutURL, Settled: charge.Settled})
	}
}

// SettlementHandler is registered with the host's one signed settlement
// router. It still verifies the consumer, price and trusted identity metadata
// so a valid provider event cannot be repurposed from another product or user.
func (s *Service) SettlementHandler() payrail.ConfirmHandler {
	return func(ctx context.Context, event payrail.SettlementEvent) error {
		uid := event.Metadata[metadataUserID]
		payer := event.Metadata[metadataPayerRef]
		if event.ChargeRef == "" || event.Consumer != Consumer || event.Kind != payrail.SettlementPaid ||
			event.Amount != AmountMinorUnits || !strings.EqualFold(event.Currency, Currency) ||
			event.Metadata[metadataPaymentMode] != requiredPaymentMode ||
			uid == "" || payer == "" || uid != payer {
			return fmt.Errorf("%w: consumer=%q kind=%q amount=%d currency=%q", ErrInvalidSettlement, event.Consumer, event.Kind, event.Amount, event.Currency)
		}
		return s.marker.MarkPaid(ctx, uid, event.ChargeRef)
	}
}

func writeCheckoutResponse(w http.ResponseWriter, status int, response CreateCheckoutResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func validOptionalID(value string) bool { return value == "" || safeAttributionID.MatchString(value) }
func putIfPresent(values map[string]string, key, value string) {
	if value != "" {
		values[key] = value
	}
}
