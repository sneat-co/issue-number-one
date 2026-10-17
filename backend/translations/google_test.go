package translations

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/api/option"
)

func TestGoogleTranslator_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"translations": []map[string]string{
					{"translatedText": "Hello"},
					{"translatedText": "World"},
				},
			},
		})
	}))
	defer ts.Close()

	ctx := context.Background()
	gt, err := NewGoogleTranslator(ctx, option.WithoutAuthentication(), option.WithEndpoint(ts.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		if err := gt.Close(); err != nil {
			t.Fatalf("Close failed: %v", err)
		}
	}()

	res, err := gt.Translate(ctx, "es", "en", []string{"Hola", "Mundo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 || res[0] != "Hello" || res[1] != "World" {
		t.Fatalf("unexpected result: %v", res)
	}
}

func TestGoogleTranslator_Errors(t *testing.T) {
	ctx := context.Background()

	// NewGoogleTranslator error
	if _, err := NewGoogleTranslator(ctx, option.WithCredentialsJSON([]byte("invalid json"))); err == nil { //nolint:staticcheck
		t.Fatal("expected error with invalid credentials JSON")
	}

	// Close with nil receiver and nil Client
	var nilGT *GoogleTranslator
	if err := nilGT.Close(); err != nil {
		t.Fatalf("unexpected error on nil Close: %v", err)
	}
	emptyGT := &GoogleTranslator{}
	if err := emptyGT.Close(); err != nil {
		t.Fatalf("unexpected error on empty Close: %v", err)
	}

	// Translate on nil receiver and nil Client
	if _, err := nilGT.Translate(ctx, "en", "es", []string{"test"}); err == nil {
		t.Fatal("expected error on nil receiver Translate")
	}
	if _, err := emptyGT.Translate(ctx, "en", "es", []string{"test"}); err == nil {
		t.Fatal("expected error on empty Client Translate")
	}

	// Translate with invalid source and target language
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	gt, err := NewGoogleTranslator(ctx, option.WithoutAuthentication(), option.WithEndpoint(ts.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = gt.Close() }()

	if _, err := gt.Translate(ctx, "invalid-source-language-tag-!!!", "es", []string{"test"}); err == nil {
		t.Fatal("expected error for invalid source language")
	}
	if _, err := gt.Translate(ctx, "en", "invalid-target-language-tag-!!!", []string{"test"}); err == nil {
		t.Fatal("expected error for invalid target language")
	}

	// API error from server
	if _, err := gt.Translate(ctx, "en", "es", []string{"test"}); err == nil {
		t.Fatal("expected error from 500 server response")
	}
}
