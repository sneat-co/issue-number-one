package publicqa

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeIdentity struct {
	caller Caller
	ok     bool
}

func (f fakeIdentity) Caller(context.Context, string) (Caller, bool) { return f.caller, f.ok }
func TestRegisterRoutes(t *testing.T) {
	h := NewHandler(nil, fakeIdentity{})
	got := map[string]bool{}
	h.RegisterHttpRoutes(func(m, p string, _ http.HandlerFunc) { got[m+" "+p] = true })
	for _, r := range []string{"GET /v0/issuenumber/catalog", "GET /v0/issuenumber/question", "GET /v0/issuenumber/answer", "GET /v0/issuenumber/eligibility", "POST /v0/issuenumber/answer", "POST /v0/issuenumber/question"} {
		if !got[r] {
			t.Errorf("missing %s", r)
		}
	}
}
func TestRequiredCallerRejectsMissingIdentity(t *testing.T) {
	h := NewHandler(nil, fakeIdentity{})
	w := httptest.NewRecorder()
	_, ok := h.requiredCaller(w, httptest.NewRequest("GET", "/", nil))
	if ok || w.Code != http.StatusUnauthorized {
		t.Fatalf("ok=%v status=%d", ok, w.Code)
	}
}
func TestVerificationRequiredResponse(t *testing.T) {
	w := httptest.NewRecorder()
	writeServiceError(w, ErrVerificationRequired)
	if w.Code != http.StatusForbidden {
		t.Fatalf("status=%d", w.Code)
	}
}
