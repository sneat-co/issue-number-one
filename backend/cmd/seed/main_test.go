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

type mockTranslator struct {
	err error
}

func (m *mockTranslator) Translate(ctx context.Context, sourceLanguage, targetLanguage string, texts []string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return texts, nil
}

func (m *mockTranslator) Close() error {
	return nil
}

func emulatorSeedClient(t *testing.T) *firestore.Client {
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

func seedJSONPath() string {
	dir, _ := os.Getwd()
	for d := dir; d != "/" && d != "."; d = filepath.Dir(d) {
		target := filepath.Join(d, "catalog", "seed.json")
		if _, err := os.Stat(target); err == nil {
			return target
		}
	}
	return "../../catalog/seed.json"
}

func TestSeed_Validation(t *testing.T) {
	// 1. Invalid flag
	if err := run([]string{"--invalid-flag"}); err == nil {
		t.Fatal("expected error on invalid flag")
	}

	// 2. Missing project
	oldProj := os.Getenv("GCLOUD_PROJECT")
	_ = os.Unsetenv("GCLOUD_PROJECT")
	defer func() {
		if oldProj != "" {
			_ = os.Setenv("GCLOUD_PROJECT", oldProj)
		}
	}()
	if err := run([]string{}); err == nil {
		t.Fatal("expected error on missing project")
	}

	// 3. Production confirmation check when FIRESTORE_EMULATOR_HOST is unset
	oldEnv := os.Getenv("FIRESTORE_EMULATOR_HOST")
	_ = os.Unsetenv("FIRESTORE_EMULATOR_HOST")
	defer func() {
		if oldEnv != "" {
			_ = os.Setenv("FIRESTORE_EMULATOR_HOST", oldEnv)
		}
	}()
	if err := run([]string{"--project", "prod-proj", "--confirm-production-project", "wrong"}); err == nil {
		t.Fatal("expected error on unconfirmed production project")
	}
	_ = os.Setenv("FIRESTORE_EMULATOR_HOST", oldEnv)

	// 4. File read error
	if err := run([]string{"--project", "in1-test", "--file", "nonexistent-seed.json"}); err == nil {
		t.Fatal("expected error on nonexistent file")
	}

	// 5. Invalid JSON
	tmpDir := t.TempDir()
	badJSON := filepath.Join(tmpDir, "bad.json")
	if err := os.WriteFile(badJSON, []byte("invalid json"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"--project", "in1-test", "--file", badJSON}); err == nil {
		t.Fatal("expected error on invalid json")
	}

	// 6. SchemaVersion != 1
	badVer := filepath.Join(tmpDir, "bad-version.json")
	if err := os.WriteFile(badVer, []byte(`{"schemaVersion": 2}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"--project", "in1-test", "--file", badVer}); err == nil {
		t.Fatal("expected error on unsupported schema version")
	}
}

func TestSeed_SuccessAndIdempotency(t *testing.T) {
	emulatorSeedClient(t)
	seedPath := seedJSONPath()

	// Run 1: seed catalog from catalog seed JSON
	space := "in1-seed-test-space"
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", seedPath}); err != nil {
		t.Fatalf("first seed run failed: %v", err)
	}

	// Run 2: idempotent rerun (exercises existing issues where status is already published)
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", seedPath}); err != nil {
		t.Fatalf("idempotent seed run failed: %v", err)
	}
}

func TestSeed_TranslateAll(t *testing.T) {
	emulatorSeedClient(t)
	space := "in1-seed-trans-space"
	seedPath := seedJSONPath()

	// Mock translations
	oldRepo := newTranslationRepo
	oldGT := newGoogleTranslator
	oldTS := newTranslationService
	defer func() {
		newTranslationRepo = oldRepo
		newGoogleTranslator = oldGT
		newTranslationService = oldTS
	}()

	// 1. Error from newTranslationRepo (empty space)
	if err := run([]string{"--project", "in1-test", "--space", "   ", "--file", seedPath, "--translate-all"}); err == nil {
		t.Fatal("expected error on empty space")
	}

	// 2. Error from newGoogleTranslator
	newGoogleTranslator = func(ctx context.Context, opts ...option.ClientOption) (seedTranslator, error) {
		return nil, errors.New("translator init error")
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", seedPath, "--translate-all"}); err == nil {
		t.Fatal("expected error on translator failure")
	}
	newGoogleTranslator = oldGT

	// 3. Error from newTranslationService
	newTranslationService = func(repo translations.Repository, translator translations.Translator, config translations.Config) (*translations.Service, error) {
		return nil, errors.New("service init error")
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", seedPath, "--translate-all"}); err == nil {
		t.Fatal("expected error on service failure")
	}
	newTranslationService = oldTS

	// 4. Error translating question
	newGoogleTranslator = func(ctx context.Context, opts ...option.ClientOption) (seedTranslator, error) {
		return &mockTranslator{err: errors.New("translate failed")}, nil
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", seedPath, "--translate-all"}); err == nil {
		t.Fatal("expected error on translate question failure")
	}

	// 5. Success translating questions
	newGoogleTranslator = func(ctx context.Context, opts ...option.ClientOption) (seedTranslator, error) {
		return &mockTranslator{}, nil
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", seedPath, "--translate-all"}); err != nil {
		t.Fatalf("unexpected error on translate-all: %v", err)
	}
}

func TestSeed_CatalogValidationBranches(t *testing.T) {
	emulatorSeedClient(t)
	space := "in1-seed-branch-space"
	tmpDir := t.TempDir()

	writeCatalog := func(c catalog) string {
		path := filepath.Join(tmpDir, "cat.json")
		b, _ := json.Marshal(c)
		_ = os.WriteFile(path, b, 0644)
		return path
	}

	// 1. Category with empty ID
	catFile1 := writeCatalog(catalog{
		SchemaVersion: 1,
		Categories:    []seedCategory{{Category: publicqa.Category{ID: ""}}},
	})
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", catFile1}); err == nil {
		t.Fatal("expected error on empty category id")
	}

	// 2. Concept with empty ID
	catFile2 := writeCatalog(catalog{
		SchemaVersion: 1,
		Categories:    []seedCategory{{Category: publicqa.Category{ID: "c1"}}},
		Concepts:      []seedConcept{{Concept: publicqa.Concept{ID: ""}}},
	})
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", catFile2}); err == nil {
		t.Fatal("expected error on empty concept id")
	}

	// 3. Question with empty ID
	catFile3 := writeCatalog(catalog{
		SchemaVersion: 1,
		Categories:    []seedCategory{{Category: publicqa.Category{ID: "c1"}}},
		Concepts:      []seedConcept{{Concept: publicqa.Concept{ID: "con1"}}},
		Questions:     []seedQuestion{{Question: publicqa.Question{ID: ""}}},
	})
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", catFile3}); err == nil {
		t.Fatal("expected error on empty question id")
	}

	// 4. Question referencing unknown concept
	catFile4 := writeCatalog(catalog{
		SchemaVersion: 1,
		Categories:    []seedCategory{{Category: publicqa.Category{ID: "c1"}}},
		Concepts:      []seedConcept{{Concept: publicqa.Concept{ID: "con1"}}},
		Questions: []seedQuestion{{
			Question: publicqa.Question{
				ID:         "q1",
				ConceptIDs: []string{"unknown-concept-xyz"},
			},
		}},
	})
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", catFile4}); err == nil {
		t.Fatal("expected error on unknown concept id")
	}

	// 5. Question with invalid alias title causing NormalizeTitle error
	catFile5 := writeCatalog(catalog{
		SchemaVersion: 1,
		Categories:    []seedCategory{{Category: publicqa.Category{ID: "c1"}}},
		Concepts: []seedConcept{{Concept: publicqa.Concept{
			ID:      "con1",
			Title:   "Valid Concept",
			Aliases: []string{"a"}, // less than 2 runes
		}}},
		Questions: []seedQuestion{{
			Question: publicqa.Question{
				ID:         "q1",
				ConceptIDs: []string{"con1"},
			},
		}},
	})
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", catFile5}); err == nil {
		t.Fatal("expected error on invalid alias title")
	}

	// 6. Question with free ChoiceSource, question with custom ChoiceSource
	catFile6 := writeCatalog(catalog{
		SchemaVersion: 1,
		Questions: []seedQuestion{
			{Question: publicqa.Question{ID: "q-free", Slug: "free-slug", ChoiceSource: publicqa.ChoiceSource{Kind: "free"}}},
			{Question: publicqa.Question{ID: "q-predef", Slug: "predef-slug", ChoiceSource: publicqa.ChoiceSource{Kind: "predefined", EntityType: "country"}}},
		},
	})
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", catFile6}); err != nil {
		t.Fatalf("unexpected error on question with choices: %v", err)
	}

	// 7. Flush with 400 batch size
	var categories []seedCategory
	for i := 0; i < 405; i++ {
		categories = append(categories, seedCategory{Category: publicqa.Category{ID: filepath.Join("cat", string(rune('a'+i%26)), string(rune('a'+i/26)))}})
	}
	catFile7 := writeCatalog(catalog{
		SchemaVersion: 1,
		Categories:    categories,
	})
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", catFile7}); err != nil {
		t.Fatalf("unexpected error on 405 categories: %v", err)
	}
}

func TestSeed_ClientAndMain(t *testing.T) {
	emulatorSeedClient(t)
	space := "in1-seed-main-space"
	seedPath := seedJSONPath()

	// 1. Client error
	oldClient := newFirestoreClient
	newFirestoreClient = func(ctx context.Context, projectID string, opts ...option.ClientOption) (*firestore.Client, error) {
		return nil, errors.New("client creation error")
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", seedPath}); err == nil {
		t.Fatal("expected error on client creation failure")
	}
	newFirestoreClient = oldClient

	// 2. Fatal function test
	oldExit := exit
	var exitCode int
	exit = func(code int) { exitCode = code }
	defer func() { exit = oldExit }()
	fatal("test fatal message")
	if exitCode != 1 {
		t.Fatalf("expected exitCode 1, got %d", exitCode)
	}

	// 3. main() successful
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"seed", "--project", "in1-test", "--space", space, "--file", seedPath}
	main()

	// 4. main() error
	var fatalMsg string
	oldFatal := fatal
	fatal = func(s string) { fatalMsg = s }
	defer func() { fatal = oldFatal }()
	os.Args = []string{"seed", "--project", "in1-test", "--file", "nonexistent.json"}
	main()
	if fatalMsg == "" {
		t.Fatal("expected fatal to be called on error")
	}
}

func TestSeed_EmptyCatalog(t *testing.T) {
	emulatorSeedClient(t)
	space := "in1-seed-empty-space"
	tmpDir := t.TempDir()
	emptyPath := filepath.Join(tmpDir, "empty.json")
	if err := os.WriteFile(emptyPath, []byte(`{"schemaVersion": 1}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"--project", "in1-test", "--space", space, "--file", emptyPath}); err != nil {
		t.Fatalf("unexpected error on empty catalog: %v", err)
	}
}

func TestSeed_CommitError(t *testing.T) {
	emulatorSeedClient(t)
	space := "in1-seed-commit-err-space"
	tmpDir := t.TempDir()
	catPath := filepath.Join(tmpDir, "cat.json")
	var categories []seedCategory
	for i := 0; i < 405; i++ {
		categories = append(categories, seedCategory{Category: publicqa.Category{ID: filepath.Join("cat", string(rune('a'+i%26)), string(rune('a'+i/26)))}})
	}
	b, _ := json.Marshal(catalog{
		SchemaVersion: 1,
		Categories:    categories,
	})
	_ = os.WriteFile(catPath, b, 0644)

	oldContext := newContext
	defer func() { newContext = oldContext }()
	newContext = func() context.Context {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		return ctx
	}

	if err := run([]string{"--project", "in1-test", "--space", space, "--file", catPath}); err == nil {
		t.Fatal("expected error on batch.Commit cancellation")
	}
}

func TestSeed_IrefGetError(t *testing.T) {
	emulatorSeedClient(t)
	space := "in1-seed-get-err-space"
	tmpDir := t.TempDir()
	catPath := filepath.Join(tmpDir, "cat.json")
	b, _ := json.Marshal(catalog{
		SchemaVersion: 1,
		Categories:    []seedCategory{{Category: publicqa.Category{ID: "c1"}}},
		Concepts:      []seedConcept{{Concept: publicqa.Concept{ID: "con1", Title: "Valid Concept"}}},
		Questions: []seedQuestion{{
			Question: publicqa.Question{
				ID:         "q1",
				ConceptIDs: []string{"con1"},
			},
		}},
	})
	_ = os.WriteFile(catPath, b, 0644)

	oldContext := newContext
	defer func() { newContext = oldContext }()
	newContext = func() context.Context {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		return ctx
	}

	if err := run([]string{"--project", "in1-test", "--space", space, "--file", catPath}); err == nil {
		t.Fatal("expected error on iref.Get cancellation")
	}
}
