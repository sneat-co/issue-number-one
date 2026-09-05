package translations

import (
	"cloud.google.com/go/firestore"
	"context"
	"os"
	"testing"
	"time"
)

func TestFirestoreRepositoryStoresTranslationBelowQuestion(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "issuenumber-translation-test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close() })

	spaceID := "translation-test-" + time.Now().UTC().Format("20060102150405.000000000")
	questionID := "ireland"
	questionRef := client.Collection("spaces").Doc(spaceID).Collection("ext").Doc("issuenumber").Collection("questions").Doc(questionID)
	if _, err := questionRef.Set(ctx, map[string]any{
		"title": "What is the top issue in Ireland?", "description": "Choose one.",
		"sourceLanguage": "en", "contentRevision": int64(1), "publication": "published",
	}); err != nil {
		t.Fatal(err)
	}
	repository, err := NewFirestoreRepository(client, spaceID)
	if err != nil {
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
