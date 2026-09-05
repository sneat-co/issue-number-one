package publicqa

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type UserIdentity interface {
	Caller(context.Context, string) (Caller, bool)
}
type Handler struct {
	service  *Service
	identity UserIdentity
}

func NewHandler(service *Service, identity UserIdentity) *Handler {
	return &Handler{service: service, identity: identity}
}
func (h *Handler) RegisterHttpRoutes(handle func(string, string, http.HandlerFunc)) {
	handle(http.MethodGet, "/v0/issuenumber/catalog", h.getCatalog)
	handle(http.MethodGet, "/v0/issuenumber/question", h.getQuestion)
	handle(http.MethodGet, "/v0/issuenumber/answer", h.getAnswer)
	handle(http.MethodGet, "/v0/issuenumber/eligibility", h.getEligibility)
	handle(http.MethodPost, "/v0/issuenumber/answer", h.postAnswer)
	handle(http.MethodPost, "/v0/issuenumber/question", h.postQuestion)
}
func (h *Handler) postQuestion(w http.ResponseWriter, r *http.Request) {
	c, ok := h.requiredCaller(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	defer r.Body.Close()
	var req CreateQuestionRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&req); err != nil {
		writeError(w, 400, "bad_request", "invalid JSON body")
		return
	}
	v, err := h.service.CreateQuestion(r.Context(), c, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	privateJSON(w, v)
}
func (h *Handler) getEligibility(w http.ResponseWriter, r *http.Request) {
	c, ok := h.requiredCaller(w, r)
	if !ok {
		return
	}
	paid, err := h.service.HasPaid(r.Context(), c.UID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	privateJSON(w, map[string]bool{"phoneVerified": c.PhoneVerified, "paid": paid, "eligible": c.PhoneVerified || paid})
}
func bearer(r *http.Request) string {
	return strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
}
func (h *Handler) optionalCaller(r *http.Request) (Caller, bool) {
	if h.identity == nil {
		return Caller{}, false
	}
	return h.identity.Caller(r.Context(), bearer(r))
}
func (h *Handler) requiredCaller(w http.ResponseWriter, r *http.Request) (Caller, bool) {
	c, ok := h.optionalCaller(r)
	if !ok || c.UID == "" {
		writeError(w, 401, "unauthorized", "missing or invalid bearer token")
		return Caller{}, false
	}
	return c, true
}
func (h *Handler) getCatalog(w http.ResponseWriter, r *http.Request) {
	v, e := h.service.Catalog(r.Context())
	if e != nil {
		writeServiceError(w, e)
		return
	}
	publicJSON(w, v)
}
func (h *Handler) getQuestion(w http.ResponseWriter, r *http.Request) {
	c, _ := h.optionalCaller(r)
	qid, slug := r.URL.Query().Get("questionId"), r.URL.Query().Get("slug")
	language := r.URL.Query().Get("lang")
	var v QuestionResponse
	var e error
	if qid != "" {
		v, e = h.service.QuestionLocalized(r.Context(), qid, c.UID, language)
	} else {
		v, e = h.service.QuestionBySlugLocalized(r.Context(), slug, c.UID, language)
	}
	if e != nil {
		writeServiceError(w, e)
		return
	}
	publicJSON(w, v)
}
func (h *Handler) getAnswer(w http.ResponseWriter, r *http.Request) {
	c, ok := h.requiredCaller(w, r)
	if !ok {
		return
	}
	v, e := h.service.OwnAnswer(r.Context(), c.UID, r.URL.Query().Get("questionId"))
	if e != nil {
		writeServiceError(w, e)
		return
	}
	privateJSON(w, v)
}
func (h *Handler) postAnswer(w http.ResponseWriter, r *http.Request) {
	c, ok := h.requiredCaller(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	defer r.Body.Close()
	var req AnswerRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if e := d.Decode(&req); e != nil {
		writeError(w, 400, "bad_request", "invalid JSON body")
		return
	}
	v, e := h.service.Answer(r.Context(), c, req)
	if e != nil {
		writeServiceError(w, e)
		return
	}
	privateJSON(w, v)
}
func publicJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Cache-Control", "public, max-age=30, s-maxage=120, stale-while-revalidate=300")
	writeJSON(w, 200, v)
}
func privateJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Cache-Control", "private, no-store")
	writeJSON(w, 200, v)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"code": code, "message": message})
}
func writeServiceError(w http.ResponseWriter, e error) {
	switch {
	case errors.Is(e, ErrNotFound):
		writeError(w, 404, "not_found", "resource not found")
	case errors.Is(e, ErrConflict):
		writeError(w, 409, "idempotency_conflict", ErrConflict.Error())
	case errors.Is(e, ErrVerificationRequired):
		writeError(w, 403, "verification_required", "verify by phone or EUR 1.00 support payment before answering")
	case errors.Is(e, ErrRateLimited):
		w.Header().Set("Retry-After", "86400")
		writeError(w, 429, "rate_limited", ErrRateLimited.Error())
	case errors.Is(e, ErrValidation):
		writeError(w, 400, "bad_request", e.Error())
	default:
		writeError(w, 500, "internal", "internal error")
	}
}
