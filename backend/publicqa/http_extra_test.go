package publicqa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestWriteServiceError_AllBranches(t *testing.T) {
	for name, tc := range map[string]struct {
		err          error
		expectedCode int
		expectedMsg  string
		checkHeader  bool
	}{
		"not found": {
			err:          ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		"conflict": {
			err:          ErrConflict,
			expectedCode: http.StatusConflict,
		},
		"verification required": {
			err:          ErrVerificationRequired,
			expectedCode: http.StatusForbidden,
		},
		"rate limited": {
			err:          ErrRateLimited,
			expectedCode: http.StatusTooManyRequests,
			checkHeader:  true,
		},
		"validation": {
			err:          errors.Join(ErrValidation, errors.New("bad input")),
			expectedCode: http.StatusBadRequest,
		},
		"internal error": {
			err:          errors.New("unexpected explosion"),
			expectedCode: http.StatusInternalServerError,
		},
	} {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeServiceError(w, tc.err)
			if w.Code != tc.expectedCode {
				t.Fatalf("expected status %d, got %d", tc.expectedCode, w.Code)
			}
			if tc.checkHeader && w.Header().Get("Retry-After") != "86400" {
				t.Fatalf("expected Retry-After header 86400, got %q", w.Header().Get("Retry-After"))
			}
		})
	}
}

func TestHandler_OptionalCaller(t *testing.T) {
	// identity is nil
	hNil := NewHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if _, ok := hNil.optionalCaller(req); ok {
		t.Fatal("expected false when identity is nil")
	}

	// with identity and token
	h := NewHandler(nil, fakeIdentity{caller: Caller{UID: "user-123", PhoneVerified: true}, ok: true})
	req.Header.Set("Authorization", "Bearer valid-token")
	caller, ok := h.optionalCaller(req)
	if !ok || caller.UID != "user-123" {
		t.Fatalf("expected user-123, got %v (ok=%v)", caller, ok)
	}
}

func TestHandler_HttpEndpoints(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	identity := fakeIdentity{caller: Caller{UID: "user-http", PhoneVerified: true}, ok: true}
	h := NewHandler(s, identity)

	// Seed question
	seedQuestion(t, s, "q-http", "general", "issue-1")

	// 1. getCatalog
	t.Run("getCatalog", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v0/issuenumber/catalog", nil)
		w := httptest.NewRecorder()
		h.getCatalog(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if w.Header().Get("Cache-Control") == "" {
			t.Fatal("expected Cache-Control header on public JSON")
		}
	})

	// 2. getEligibility
	t.Run("getEligibility unauthorized", func(t *testing.T) {
		hUnauth := NewHandler(s, fakeIdentity{})
		req := httptest.NewRequest(http.MethodGet, "/v0/issuenumber/eligibility", nil)
		w := httptest.NewRecorder()
		hUnauth.getEligibility(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("getEligibility success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v0/issuenumber/eligibility", nil)
		w := httptest.NewRecorder()
		h.getEligibility(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	// 3. getQuestion (by questionId and by slug)
	t.Run("getQuestion by questionId", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v0/issuenumber/question?questionId=q-http", nil)
		w := httptest.NewRecorder()
		h.getQuestion(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("getQuestion by slug not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v0/issuenumber/question?slug=nonexistent-slug", nil)
		w := httptest.NewRecorder()
		h.getQuestion(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	// 4. postQuestion
	t.Run("postQuestion unauthorized", func(t *testing.T) {
		hUnauth := NewHandler(s, fakeIdentity{})
		req := httptest.NewRequest(http.MethodPost, "/v0/issuenumber/question", nil)
		w := httptest.NewRecorder()
		hUnauth.postQuestion(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("postQuestion invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v0/issuenumber/question", bytes.NewReader([]byte("{invalid-json")))
		w := httptest.NewRecorder()
		h.postQuestion(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("postQuestion success and error", func(t *testing.T) {
		validBody, _ := json.Marshal(CreateQuestionRequest{
			Title:            "What is the top issue in Education?",
			Description:      "Please select or suggest the top issue.",
			AnswerTargetType: "issue",
			OperationID:      "op-create-q-1",
			ChoiceSource:     ChoiceSource{Kind: "free"},
			AllowSuggestions: true,
		})
		req := httptest.NewRequest(http.MethodPost, "/v0/issuenumber/question", bytes.NewReader(validBody))
		w := httptest.NewRecorder()
		h.postQuestion(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		// Invalid request (missing category)
		invalidBody, _ := json.Marshal(CreateQuestionRequest{
			Title: "Too short",
		})
		reqErr := httptest.NewRequest(http.MethodPost, "/v0/issuenumber/question", bytes.NewReader(invalidBody))
		wErr := httptest.NewRecorder()
		h.postQuestion(wErr, reqErr)
		if wErr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", wErr.Code)
		}
	})

	// 5. getAnswer
	t.Run("getAnswer unauthorized", func(t *testing.T) {
		hUnauth := NewHandler(s, fakeIdentity{})
		req := httptest.NewRequest(http.MethodGet, "/v0/issuenumber/answer?questionId=q-http", nil)
		w := httptest.NewRecorder()
		hUnauth.getAnswer(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("getAnswer not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v0/issuenumber/answer?questionId=nonexistent-q", nil)
		w := httptest.NewRecorder()
		h.getAnswer(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("getAnswer success", func(t *testing.T) {
		// First answer the question
		ansReq := httptest.NewRequest(http.MethodPost, "/v0/issuenumber/answer", bytes.NewReader([]byte(`{"questionId":"q-http","issueId":"issue-1","operationId":"op-get-ans-1"}`)))
		wAns := httptest.NewRecorder()
		h.postAnswer(wAns, ansReq)
		if wAns.Code != http.StatusOK {
			t.Fatalf("expected 200 on postAnswer, got %d: %s", wAns.Code, wAns.Body.String())
		}

		// Now get the answer
		req := httptest.NewRequest(http.MethodGet, "/v0/issuenumber/answer?questionId=q-http", nil)
		w := httptest.NewRecorder()
		h.getAnswer(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 on getAnswer, got %d: %s", w.Code, w.Body.String())
		}
	})

	// 6. postAnswer
	t.Run("postAnswer unauthorized", func(t *testing.T) {
		hUnauth := NewHandler(s, fakeIdentity{})
		req := httptest.NewRequest(http.MethodPost, "/v0/issuenumber/answer", nil)
		w := httptest.NewRecorder()
		hUnauth.postAnswer(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("postAnswer invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v0/issuenumber/answer", bytes.NewReader([]byte("{invalid-json")))
		w := httptest.NewRecorder()
		h.postAnswer(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("postAnswer success and error", func(t *testing.T) {
		answerBody, _ := json.Marshal(AnswerRequest{
			QuestionID:  "q-http",
			IssueID:     "issue-1",
			OperationID: "op-1",
		})
		req := httptest.NewRequest(http.MethodPost, "/v0/issuenumber/answer", bytes.NewReader(answerBody))
		w := httptest.NewRecorder()
		h.postAnswer(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		// Answer error (e.g. non-existent question)
		badAnswerBody, _ := json.Marshal(AnswerRequest{
			QuestionID:  "q-nonexistent",
			IssueID:     "issue-1",
			OperationID: "op-2",
		})
		reqErr := httptest.NewRequest(http.MethodPost, "/v0/issuenumber/answer", bytes.NewReader(badAnswerBody))
		wErr := httptest.NewRecorder()
		h.postAnswer(wErr, reqErr)
		if wErr.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d: %s", wErr.Code, wErr.Body.String())
		}
	})

	t.Run("handlers with cancelled context", func(t *testing.T) {
		cancelReq := httptest.NewRequest(http.MethodGet, "/v0/issuenumber/catalog", nil)
		cancelCtx, cancel := context.WithCancel(cancelReq.Context())
		cancel()
		cancelReq = cancelReq.WithContext(cancelCtx)

		// getCatalog with cancelled context
		wCat := httptest.NewRecorder()
		h.getCatalog(wCat, cancelReq)
		if wCat.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", wCat.Code)
		}

		// getEligibility with cancelled context
		cancelReqElig := httptest.NewRequest(http.MethodGet, "/v0/issuenumber/eligibility", nil).WithContext(cancelCtx)
		wElig := httptest.NewRecorder()
		h.getEligibility(wElig, cancelReqElig)
		if wElig.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", wElig.Code)
		}
	})

}
