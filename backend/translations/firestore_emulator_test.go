package translations

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/status"
)

func setupFirestoreClient(t *testing.T) (*firestore.Client, context.Context) {
	t.Helper()
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "issuenumber-translation-test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return client, ctx
}

func TestFirestoreRepository_Validation(t *testing.T) {
	client, _ := setupFirestoreClient(t)

	if _, err := NewFirestoreRepository(nil, "space"); err == nil {
		t.Fatal("expected error with nil client")
	}
	if _, err := NewFirestoreRepository(client, "   "); err == nil {
		t.Fatal("expected error with empty spaceID")
	}
}

func TestFirestoreRepositoryStoresTranslationBelowQuestion(t *testing.T) {
	client, ctx := setupFirestoreClient(t)

	spaceID := "translation-test-" + time.Now().UTC().Format("20060102150405.000000000")
	questionID := "ireland"
	repository, err := NewFirestoreRepository(client, spaceID)
	if err != nil {
		t.Fatal(err)
	}
	questionRef := repository.questionRef(questionID)
	if _, err := questionRef.Set(ctx, map[string]any{
		"title": "What is the top issue in Ireland?", "description": "Choose one.",
		"sourceLanguage": "en", "contentRevision": int64(1), "publication": "published",
	}); err != nil {
		t.Fatal(err)
	}
	service, err := NewService(repository, &fakeTranslator{}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	value, err := service.TranslateQuestion(ctx, TranslateQuestionRequest{QuestionID: questionID, Language: "ru"})
	if err != nil {
		t.Fatal(err)
	}
	stored, err := questionRef.Collection("translations").Doc("ru").Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var document Translation
	if err := stored.DataTo(&document); err != nil {
		t.Fatal(err)
	}
	if document.Title != value.Title || document.SourceRevision != 1 || !document.MachineTranslated {
		t.Fatalf("unexpected translation: %+v", document)
	}
}

func TestFirestoreRepository_GetQuestion_Branches(t *testing.T) {
	client, ctx := setupFirestoreClient(t)
	spaceID := "q-branches-" + time.Now().UTC().Format("20060102150405.000000000")
	repo, err := NewFirestoreRepository(client, spaceID)
	if err != nil {
		t.Fatal(err)
	}

	// 1. NotFound
	if _, err := repo.GetQuestion(ctx, "nonexistent"); !errors.Is(err, ErrQuestionNotFound) {
		t.Fatalf("expected ErrQuestionNotFound, got %v", err)
	}

	// 2. Context cancelled / general error
	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := repo.GetQuestion(cancelledCtx, "any"); err == nil {
		t.Fatal("expected error on cancelled context")
	}

	// 3. Status fallback when publication is empty
	qRef := repo.questionRef("q-status")
	if _, err := qRef.Set(ctx, map[string]any{
		"title": "Status test", "status": "draft", "contentRevision": int64(2),
	}); err != nil {
		t.Fatal(err)
	}
	q, err := repo.GetQuestion(ctx, "q-status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Publication != "draft" {
		t.Fatalf("expected status 'draft', got %q", q.Publication)
	}

	// 4. Decode error (e.g. title is array instead of string)
	corruptRef := repo.questionRef("q-corrupt")
	if _, err := corruptRef.Set(ctx, map[string]any{
		"title": []int{1, 2, 3},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetQuestion(ctx, "q-corrupt"); err == nil {
		t.Fatal("expected decode error")
	}
}

func TestFirestoreRepository_GetTranslation_Branches(t *testing.T) {
	client, ctx := setupFirestoreClient(t)
	spaceID := "tr-branches-" + time.Now().UTC().Format("20060102150405.000000000")
	repo, err := NewFirestoreRepository(client, spaceID)
	if err != nil {
		t.Fatal(err)
	}

	// 1. NotFound error
	if _, err := repo.GetTranslation(ctx, "q1", "es"); err == nil {
		t.Fatal("expected error for nonexistent translation")
	}

	// 2. Decode error on corrupted translation
	trRef := repo.questionRef("q1").Collection("translations").Doc("es")
	if _, err := trRef.Set(ctx, map[string]any{
		"title": []int{1, 2, 3},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetTranslation(ctx, "q1", "es"); err == nil {
		t.Fatal("expected decode error on corrupted translation")
	}

	// 3. Success
	now := time.Now().UTC()
	validTr := Translation{
		Language:  "es",
		Title:     "Hola",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := trRef.Set(ctx, validTr); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetTranslation(ctx, "q1", "es")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != "Hola" {
		t.Fatalf("expected 'Hola', got %q", got.Title)
	}
}

func TestFirestoreRepository_PutTranslationIfFresh_Branches(t *testing.T) {
	client, ctx := setupFirestoreClient(t)
	spaceID := "put-branches-" + time.Now().UTC().Format("20060102150405.000000000")
	repo, err := NewFirestoreRepository(client, spaceID)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	qSource := QuestionSource{
		ID:              "q1",
		Title:           "Question Title",
		Description:     "Description",
		SourceLanguage:  "en",
		ContentRevision: 1,
		Publication:     "published",
	}
	hash := SourceHash(qSource, "en")
	tr := Translation{
		Language:       "ru",
		Title:          "Заголовок",
		Description:    "Описание",
		SourceLanguage: "en",
		SourceRevision: 1,
		SourceHash:     hash,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 1. Question does not exist -> transaction fails
	if _, err := repo.PutTranslationIfFresh(ctx, qSource, tr); err == nil {
		t.Fatal("expected error when question document does not exist")
	}

	// 2. Corrupt question document -> decode error in transaction
	qRef := repo.questionRef("q1")
	if _, err := qRef.Set(ctx, map[string]any{"title": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.PutTranslationIfFresh(ctx, qSource, tr); err == nil {
		t.Fatal("expected decode error for corrupt question")
	}

	// 3. Stale revision / hash
	if _, err := qRef.Set(ctx, map[string]any{
		"title":           "Question Title",
		"description":     "Description",
		"sourceLanguage":  "",       // tests default to "en" as well
		"contentRevision": int64(2), // mismatch with qSource.ContentRevision=1
		"publication":     "published",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.PutTranslationIfFresh(ctx, qSource, tr); !errors.Is(err, ErrStaleSource) {
		t.Fatalf("expected ErrStaleSource, got %v", err)
	}

	// Now set live question matching qSource
	if _, err := qRef.Set(ctx, map[string]any{
		"title":           qSource.Title,
		"description":     qSource.Description,
		"sourceLanguage":  "", // fallback to "en"
		"contentRevision": qSource.ContentRevision,
		"publication":     qSource.Publication,
	}); err != nil {
		t.Fatal(err)
	}

	// 4. Stale hash mismatch (same revision, changed title)
	staleHashTr := tr
	staleHashTr.SourceHash = "mismatched-hash"
	if _, err := repo.PutTranslationIfFresh(ctx, qSource, staleHashTr); !errors.Is(err, ErrStaleSource) {
		t.Fatalf("expected ErrStaleSource on hash mismatch, got %v", err)
	}

	// 5. Existing translation corrupted -> DataTo error in transaction
	tRef := qRef.Collection("translations").Doc("ru")
	if _, err := tRef.Set(ctx, map[string]any{"title": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.PutTranslationIfFresh(ctx, qSource, tr); err == nil {
		t.Fatal("expected error on corrupted existing translation")
	}

	// 6. Existing human corrected translation -> preserved without overwrite
	correctedTr := Translation{
		Language:       "ru",
		Title:          "Ручной перевод",
		CorrectedByUID: "translator-1",
		SourceRevision: 1,
		SourceHash:     hash,
		CreatedAt:      now.Add(-time.Hour),
		UpdatedAt:      now.Add(-time.Hour),
	}
	if _, err := tRef.Set(ctx, correctedTr); err != nil {
		t.Fatal(err)
	}
	// Try putting machine translation over corrected
	res, err := repo.PutTranslationIfFresh(ctx, qSource, tr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Title != "Ручной перевод" || res.CorrectedByUID != "translator-1" {
		t.Fatalf("expected existing corrected translation preserved, got %+v", res)
	}

	// 7. Existing uncorrected translation -> overwritten, but CreatedAt preserved
	uncorrectedTr := Translation{
		Language:       "ru",
		Title:          "Старый автоперевод",
		SourceRevision: 1,
		SourceHash:     hash,
		CreatedAt:      now.Add(-2 * time.Hour),
		UpdatedAt:      now.Add(-2 * time.Hour),
	}
	if _, err := tRef.Set(ctx, uncorrectedTr); err != nil {
		t.Fatal(err)
	}
	// Put new translation with zero CreatedAt
	zeroCreatedTr := tr
	zeroCreatedTr.CreatedAt = time.Time{}
	res, err = repo.PutTranslationIfFresh(ctx, qSource, zeroCreatedTr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Title != tr.Title {
		t.Fatalf("expected title updated to %q, got %q", tr.Title, res.Title)
	}
	if res.CreatedAt.IsZero() || !res.CreatedAt.Equal(uncorrectedTr.CreatedAt) {
		t.Fatalf("expected CreatedAt preserved from existing translation, got %v", res.CreatedAt)
	}

	// 8. Error from transaction.Get(tref) that is not codes.NotFound (e.g. invalid document ID > 1500 bytes)
	longIDTr := tr
	longIDTr.Language = strings.Repeat("a", 2000)
	if _, err := repo.PutTranslationIfFresh(ctx, qSource, longIDTr); err == nil {
		t.Fatal("expected error for overly long document ID")
	} else {
		t.Logf("longID err: %v, code: %v", err, status.Code(err))
	}

	// 9. Put new translation where document doesn't exist and CreatedAt is zero
	_ = repo.questionRef("q-new").Collection("translations").Doc("ru")
	newQRef := repo.questionRef("q-new")
	newQ := QuestionSource{ID: "q-new", Title: "New Q", ContentRevision: 1, Publication: "published"}
	newHash := SourceHash(newQ, "en")
	if _, err := newQRef.Set(ctx, map[string]any{
		"title": "New Q", "sourceLanguage": "en", "contentRevision": int64(1), "publication": "published",
	}); err != nil {
		t.Fatal(err)
	}
	zeroTr := Translation{
		Language:       "ru",
		Title:          "Новый",
		SourceLanguage: "en",
		SourceRevision: 1,
		SourceHash:     newHash,
	}
	res, err = repo.PutTranslationIfFresh(ctx, newQ, zeroTr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be populated automatically")
	}
}

func TestFirestoreRepository_PutTranslationIfFresh_NonNotFoundGetError(t *testing.T) {
	client, ctx := setupFirestoreClient(t)
	spaceID := "non-not-found-" + time.Now().UTC().Format("20060102150405.000000000")
	repo, err := NewFirestoreRepository(client, spaceID)
	if err != nil {
		t.Fatal(err)
	}
	qRef := repo.questionRef("q1")
	qSource := QuestionSource{ID: "q1", Title: "Title", ContentRevision: 1, Publication: "published"}
	if _, err := qRef.Set(ctx, map[string]any{
		"title": "Title", "sourceLanguage": "en", "contentRevision": int64(1), "publication": "published",
	}); err != nil {
		t.Fatal(err)
	}
	tr := Translation{
		Language:       "__reserved__",
		Title:          "Reserved",
		SourceLanguage: "en",
		SourceRevision: 1,
		SourceHash:     SourceHash(qSource, "en"),
	}
	if _, err := repo.PutTranslationIfFresh(ctx, qSource, tr); err == nil {
		t.Fatal("expected error on reserved document id")
	}
}
