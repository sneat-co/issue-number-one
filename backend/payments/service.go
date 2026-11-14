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
	"github.com/strongo/logus"
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

// Log messages and reasons for settlement events that are acknowledged without
// granting eligibility. They are stable so operators can search and alert on them.
const (
	msgPaidNotApplied          = "issuenumber payments: paid settlement NOT applied, operator review required"
	msgRefundAcknowledged      = "issuenumber payments: refund acknowledged, eligibility unchanged (refund policy undecided)"
	msgFailedAcknowledged      = "issuenumber payments: failed charge acknowledged, nothing to do"
	msgUnknownKindAcknowledged = "issuenumber payments: unknown settlement kind acknowledged, nothing applied"

	reasonChargeMissing    = "charge reference missing"
	reasonConsumerMismatch = "consumer mismatch"
	reasonAmountMismatch   = "amount mismatch"
	reasonCurrencyMismatch = "currency mismatch"
	reasonModeNotLive      = "payment mode is not live"
	reasonPayerMissing     = "user or payer metadata missing"
	reasonPayerMismatch    = "user and payer metadata differ"
)

// SettlementHandler is registered with the host's one signed settlement
// router. It still verifies the consumer, price and trusted identity metadata
// so a valid provider event cannot be repurposed from another product or user.
//
// The host answers the provider with a retry for any returned error. An event
// that can never become valid by being resent is therefore acknowledged (nil)
// and logged through logus, which the host ships to Cloud Logging with
// severity: a refund or failure, an unknown kind, and a paid event that fails a
// check (ERROR). Nothing is granted for any of them. Every eligibility-store
// error is returned and so retried. That covers transient storage failures and
// also the store's two permanent refusals (a charge ID it cannot key, a receipt
// already bound to another user); neither is reachable with provider charge IDs
// and server-written metadata, and retrying them keeps the provider's
// failed-delivery signal.
func (s *Service) SettlementHandler() payrail.ConfirmHandler {
	return func(ctx context.Context, event payrail.SettlementEvent) error {
		switch event.Kind {
		case payrail.SettlementPaid:
		case payrail.SettlementRefunded:
			// Whether a refund revokes paid eligibility is an open product decision;
			// until it is made, eligibility is left exactly as it was.
			logus.Warningf(ctx, "%s: %s", msgRefundAcknowledged, settlementFacts(event))
			return nil
		case payrail.SettlementFailed:
			logus.Infof(ctx, "%s: %s", msgFailedAcknowledged, settlementFacts(event))
			return nil
		default:
			logus.Warningf(ctx, "%s: %s", msgUnknownKindAcknowledged, settlementFacts(event))
			return nil
		}
		if reason := invalidPaidReason(event); reason != "" {
			logus.Errorf(ctx, "%s: reason=%q %s err=%q", msgPaidNotApplied, reason, settlementFacts(event),
				fmt.Errorf("%w: %s", ErrInvalidSettlement, reason).Error())
			return nil
		}
		return s.marker.MarkPaid(ctx, event.Metadata[metadataUserID], event.ChargeRef)
	}
}

// invalidPaidReason returns why a paid event must not grant eligibility, or ""
// when it may.
func invalidPaidReason(event payrail.SettlementEvent) string {
	uid := event.Metadata[metadataUserID]
	payer := event.Metadata[metadataPayerRef]
	switch {
	case event.ChargeRef == "":
		// Defence in depth: payrail's router rejects an empty charge reference
		// before calling any handler, so only a direct call reaches this.
		return reasonChargeMissing
	case event.Consumer != Consumer:
		return reasonConsumerMismatch
	case event.Amount != AmountMinorUnits:
		return reasonAmountMismatch
	case !strings.EqualFold(event.Currency, Currency):
		return reasonCurrencyMismatch
	case event.Metadata[metadataPaymentMode] != requiredPaymentMode:
		return reasonModeNotLive
	case uid == "" || payer == "":
		return reasonPayerMissing
	case uid != payer:
		return reasonPayerMismatch
	}
	return ""
}

// settlementFacts are the provider facts an operator needs to find the charge,
// quoted so event values cannot forge log text. User identifiers are
// deliberately left out of logs.
func settlementFacts(event payrail.SettlementEvent) string {
	return fmt.Sprintf("consumer=%q kind=%q chargeRef=%q refundRef=%q amount=%d currency=%q paymentMode=%q",
		event.Consumer, string(event.Kind), event.ChargeRef, event.RefundRef,
		event.Amount, event.Currency, event.Metadata[metadataPaymentMode])
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
