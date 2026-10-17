package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/sneat-co/issue-number-one/backend/publicqa"
	"github.com/sneat-co/issue-number-one/backend/translations"
	"google.golang.org/api/option"
)

type mockTranslator struct{}

func (m *mockTranslator) Translate(ctx context.Context, sourceLanguage, targetLanguage string, texts []string) ([]string, error) {
	return texts, nil
}

func (m *mockTranslator) Close() error {
	return nil
}

func emulatorClient(t *testing.T) *firestore.Client {
	t.Helper()
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "in1-test")
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return client
}

func seedAdminQuestion(t *testing.T, db *firestore.Client, qid, space string) {
	t.Helper()
	ctx := context.Background()
	root := db.Collection("spaces").Doc(space).Collection("ext").Doc(publicqa.ExtensionID)
	qref := root.Collection("questions").Doc(qid)
	if _, err := qref.Set(ctx, map[string]any{
		"id":             qid,
		"title":          "Admin Test Question?",
		"status":         publicqa.StatusPending,
		"sourceLanguage": "en",
	}); err != nil {
		t.Fatal(err)
	}
	iref := qref.Collection("issues").Doc("issue-1")
	if _, err := iref.Set(ctx, map[string]any{
		"id":     "issue-1",
		"title":  "Issue 1",
		"status": publicqa.StatusPending,
	}); err != nil {
		t.Fatal(err)
	}
	iref2 := qref.Collection("issues").Doc("issue-2")
	if _, err := iref2.Set(ctx, map[string]any{
		"id":     "issue-2",
		"title":  "Issue 2",
		"status": publicqa.StatusPending,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestAdmin_RunValidation(t *testing.T) {
	// Flag parse error
	if err := run([]string{"--invalid-flag"}); err == nil {
		t.Fatal("expected error on invalid flag")
	}

	// Missing project or question
	if err := run([]string{}); err == nil {
		t.Fatal("expected error on empty args")
	}
	if err := run([]string{"--project", "p"}); err == nil {
		t.Fatal("expected error on missing question")
	}

	// Production confirmation check when FIRESTORE_EMULATOR_HOST is unset
	oldEnv := os.Getenv("FIRESTORE_EMULATOR_HOST")
	_ = os.Unsetenv("FIRESTORE_EMULATOR_HOST")
	defer func() {
		if oldEnv != "" {
			_ = os.Setenv("FIRESTORE_EMULATOR_HOST", oldEnv)
		}
	}()
	if err := run([]string{"--project", "prod-proj", "--question", "q1", "--confirm-production-project", "wrong"}); err == nil {
		t.Fatal("expected error on unconfirmed production project")
	}
	_ = os.Setenv("FIRESTORE_EMULATOR_HOST", oldEnv)

	// Invalid action
	if err := run([]string{"--project", "in1-test", "--question", "q1", "--action", "unknown-action"}); err == nil {
		t.Fatal("expected error on unknown action")
	}
}

func TestAdmin_Actions(t *testing.T) {
	db := emulatorClient(t)
	space := "in1-admin-test-space"
	seedAdminQuestion(t, db, "q-admin-1", space)

	// publish-question
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-1", "--action", "publish-question"}); err != nil {
		t.Fatalf("publish-question failed: %v", err)
	}

	// hide-question
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-1", "--action", "hide-question"}); err != nil {
		t.Fatalf("hide-question failed: %v", err)
	}

	// publish-issue
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-1", "--issue", "issue-1", "--action", "publish-issue"}); err != nil {
		t.Fatalf("publish-issue failed: %v", err)
	}

	// hide-issue
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-1", "--issue", "issue-1", "--action", "hide-issue"}); err != nil {
		t.Fatalf("hide-issue failed: %v", err)
	}

	// reject-issue
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-1", "--issue", "issue-1", "--action", "reject-issue"}); err != nil {
		t.Fatalf("reject-issue failed: %v", err)
	}

	// merge-issue
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-1", "--issue", "issue-1", "--target-issue", "issue-2", "--operation", "op-admin-merge", "--action", "merge-issue"}); err != nil {
		t.Fatalf("merge-issue failed: %v", err)
	}
}

func TestAdmin_PopulateCountries(t *testing.T) {
	db := emulatorClient(t)
	space := "in1-admin-countries-space"
	seedAdminQuestion(t, db, "q-admin-countries", space)

	// Missing --countries-file
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-countries", "--action", "populate-countries"}); err == nil {
		t.Fatal("expected error on missing countries-file")
	}

	// File not found
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-countries", "--action", "populate-countries", "--countries-file", "nonexistent.json"}); err == nil {
		t.Fatal("expected error on nonexistent countries-file")
	}

	// Invalid json
	badFile := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(badFile, []byte("invalid json"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-countries", "--action", "populate-countries", "--countries-file", badFile}); err == nil {
		t.Fatal("expected error on invalid json countries-file")
	}

	// Valid countries file
	goodFile := filepath.Join(t.TempDir(), "countries.json")
	data := map[string]any{
		"countriesByID": map[string]any{
			"ie": map[string]string{"id": "ie", "title": "Ireland"},
			"gb": map[string]string{"id": "gb", "title": "United Kingdom"},
		},
	}
	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(goodFile, raw, 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-countries", "--action", "populate-countries", "--countries-file", goodFile}); err != nil {
		t.Fatalf("populate-countries failed: %v", err)
	}
}

func TestAdmin_TranslateQuestion(t *testing.T) {
	db := emulatorClient(t)
	space := "in1-admin-translate-space"
	seedAdminQuestion(t, db, "q-admin-trans", space)

	// 1. repository creation error (empty space) -> line 101
	if err := run([]string{"--project", "in1-test", "--space", "   ", "--question", "q-admin-trans", "--action", "translate-question"}); err == nil {
		t.Fatal("expected error on empty space")
	}

	// 2. translator creation error -> line 105
	oldGT := newGoogleTranslator
	newGoogleTranslator = func(ctx context.Context, opts ...option.ClientOption) (adminTranslator, error) {
		return nil, errors.New("translator init error")
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-trans", "--action", "translate-question"}); err == nil {
		t.Fatal("expected error on translator init failure")
	}
	newGoogleTranslator = oldGT

	// 3. newTranslationService error -> line 110
	oldTS := newTranslationService
	newGoogleTranslator = func(ctx context.Context, opts ...option.ClientOption) (adminTranslator, error) {
		return &mockTranslator{}, nil
	}
	newTranslationService = func(repo translations.Repository, translator translations.Translator, config translations.Config) (*translations.Service, error) {
		return nil, errors.New("translation service init error")
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-trans", "--action", "translate-question"}); err == nil {
		t.Fatal("expected error on translation service init failure")
	}
	newTranslationService = oldTS

	// 4. translate-question lines 118-119
	newGoogleTranslator = func(ctx context.Context, opts ...option.ClientOption) (adminTranslator, error) {
		return &mockTranslator{}, nil
	}
	defer func() { newGoogleTranslator = oldGT }()
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q-admin-trans", "--action", "translate-question"}); err != nil {
		t.Fatalf("translate-question failed: %v", err)
	}
}

func TestAdmin_ClientAndActionErrors(t *testing.T) {
	space := "in1-admin-err-space"

	// 1. newFirestoreClient error -> line 58
	oldFS := newFirestoreClient
	newFirestoreClient = func(ctx context.Context, projectID string, opts ...option.ClientOption) (*firestore.Client, error) {
		return nil, errors.New("client creation error")
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "q1", "--action", "publish-question"}); err == nil {
		t.Fatal("expected error on firestore client creation failure")
	}
	newFirestoreClient = oldFS

	// 2. action returning error -> line 118
	if err := run([]string{"--project", "in1-test", "--space", space, "--question", "nonexistent-q-for-err", "--action", "publish-question"}); err == nil {
		t.Fatal("expected error on action failure")
	}

	// 3. fatal function exercising line 27
	oldExit := exit
	var exitCode int
	exit = func(code int) { exitCode = code }
	defer func() { exit = oldExit }()
	fatal("test fatal message")
	if exitCode != 1 {
		t.Fatalf("expected exitCode 1, got %d", exitCode)
	}
}

func TestAdmin_Main(t *testing.T) {
	db := emulatorClient(t)
	space := "in1-admin-main-space"
	seedAdminQuestion(t, db, "q-admin-main", space)

	// Test successful main()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"admin", "--project", "in1-test", "--space", space, "--question", "q-admin-main", "--action", "publish-question"}
	main()

	// Test main() on failure calls fatal
	var fatalMsg string
	oldFatal := fatal
	fatal = func(v string) { fatalMsg = v }
	defer func() { fatal = oldFatal }()

	os.Args = []string{"admin", "--project", "in1-test", "--space", space, "--question", "q-admin-main", "--action", "invalid-action"}
	main()
	if fatalMsg == "" {
		t.Fatal("expected fatal to be called on error")
	}
}

func TestAdmin_NewGoogleTranslator(t *testing.T) {
	_, _ = newGoogleTranslator(context.Background())
}
