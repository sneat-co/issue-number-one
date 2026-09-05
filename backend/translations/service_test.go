package translations

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type memoryRepository struct {
	mu           sync.Mutex
	question     QuestionSource
	translations map[string]Translation
	puts         int
}

func (repository *memoryRepository) GetQuestion(context.Context, string) (QuestionSource, error) {
	return repository.question, nil
}
func (repository *memoryRepository) GetTranslation(_ context.Context, _ string, language string) (Translation, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	value, ok := repository.translations[language]
	if !ok {
		return Translation{}, errors.New("missing")
	}
	return value, nil
}
func (repository *memoryRepository) PutTranslationIfFresh(_ context.Context, question QuestionSource, value Translation) (Translation, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if question.ContentRevision != repository.question.ContentRevision {
		return Translation{}, ErrStaleSource
	}
	repository.translations[value.Language] = value
	repository.puts++
	return value, nil
}

type fakeTranslator struct {
	mu    sync.Mutex
	calls int
	wait  chan struct{}
}

func (translator *fakeTranslator) Translate(_ context.Context, source, target string, texts []string) ([]string, error) {
	translator.mu.Lock()
	translator.calls++
	translator.mu.Unlock()
	if translator.wait != nil {
		<-translator.wait
	}
	return []string{"RU: " + texts[0], "RU: " + texts[1]}, nil
}

func testService(t *testing.T, question QuestionSource, translator *fakeTranslator) (*Service, *memoryRepository) {
	t.Helper()
	repository := &memoryRepository{question: question, translations: map[string]Translation{}}
	service, err := NewService(repository, translator, Config{Now: func() time.Time { return time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC) }})
	if err != nil {
		t.Fatal(err)
	}
	return service, repository
}

func TestTranslateQuestionCachesByRevisionAndHash(t *testing.T) {
	translator := &fakeTranslator{}
	service, repository := testService(t, QuestionSource{ID: "ireland", Title: "What is the top issue in Ireland?", Description: "Choose one.", SourceLanguage: "en", ContentRevision: 2, Publication: "published"}, translator)
	first, err := service.TranslateQuestion(context.Background(), TranslateQuestionRequest{QuestionID: "ireland", Language: "ru"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.TranslateQuestion(context.Background(), TranslateQuestionRequest{QuestionID: "ireland", Language: "RU"})
	if err != nil {
		t.Fatal(err)
	}
	if first.Title != second.Title || translator.calls != 1 || repository.puts != 1 {
		t.Fatalf("cache failed: calls=%d puts=%d", translator.calls, repository.puts)
	}
	if !first.MachineTranslated || first.SourceRevision != 2 || first.SourceHash == "" {
		t.Fatalf("bad metadata: %+v", first)
	}
}

func TestSourceLanguageReturnsSourceWithoutProviderOrWrite(t *testing.T) {
	translator := &fakeTranslator{}
	service, repository := testService(t, QuestionSource{ID: "q", Title: "Вопрос", SourceLanguage: "ru", ContentRevision: 1, Publication: "published"}, translator)
	value, err := service.TranslateQuestion(context.Background(), TranslateQuestionRequest{QuestionID: "q", Language: "ru"})
	if err != nil {
		t.Fatal(err)
	}
	if value.MachineTranslated || translator.calls != 0 || repository.puts != 0 {
		t.Fatalf("source should not be generated: %+v", value)
	}
}

func TestPendingQuestionRequiresCreator(t *testing.T) {
	service, _ := testService(t, QuestionSource{ID: "q", Title: "Private", SourceLanguage: "en", ContentRevision: 1, Publication: "pending", CreatorUID: "owner"}, &fakeTranslator{})
	if _, err := service.TranslateQuestion(context.Background(), TranslateQuestionRequest{QuestionID: "q", Language: "ru"}); !errors.Is(err, ErrTranslationDenied) {
		t.Fatalf("got %v", err)
	}
	if _, err := service.TranslateQuestion(context.Background(), TranslateQuestionRequest{QuestionID: "q", Language: "ru", Actor: Actor{UID: "owner"}}); err != nil {
		t.Fatal(err)
	}
}

func TestUnsupportedLanguageNeverFallsBack(t *testing.T) {
	service, _ := testService(t, QuestionSource{ID: "q", Title: "Question", SourceLanguage: "en", ContentRevision: 1, Publication: "published"}, &fakeTranslator{})
	if _, err := service.TranslateQuestion(context.Background(), TranslateQuestionRequest{QuestionID: "q", Language: "fr"}); !errors.Is(err, ErrUnsupportedLanguage) {
		t.Fatalf("got %v", err)
	}
}

func TestConcurrentRequestsDeduplicateProviderCall(t *testing.T) {
	wait := make(chan struct{})
	translator := &fakeTranslator{wait: wait}
	service, _ := testService(t, QuestionSource{ID: "q", Title: "Question", SourceLanguage: "en", ContentRevision: 1, Publication: "published"}, translator)
	errs := make(chan error, 2)
	for range 2 {
		go func() {
			_, err := service.TranslateQuestion(context.Background(), TranslateQuestionRequest{QuestionID: "q", Language: "ru"})
			errs <- err
		}()
	}
	time.Sleep(10 * time.Millisecond)
	close(wait)
	for range 2 {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}
	if translator.calls != 1 {
		t.Fatalf("provider calls=%d", translator.calls)
	}
}
