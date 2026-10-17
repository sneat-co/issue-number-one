package translations

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type customError struct {
	msg       string
	retryable bool
}

func (e *customError) Error() string   { return e.msg }
func (e *customError) Retryable() bool { return e.retryable }

type variableTranslator struct {
	fn func(ctx context.Context, source, target string, texts []string) ([]string, error)
}

func (v *variableTranslator) Translate(ctx context.Context, source, target string, texts []string) ([]string, error) {
	return v.fn(ctx, source, target, texts)
}

func TestNewService_OptionsAndValidation(t *testing.T) {
	repo := &memoryRepository{question: QuestionSource{ID: "q1"}}
	trans := &fakeTranslator{}

	// nil repository
	if _, err := NewService(nil, trans, Config{}); err == nil {
		t.Fatal("expected error with nil repository")
	}
	// nil translator
	if _, err := NewService(repo, nil, Config{}); err == nil {
		t.Fatal("expected error with nil translator")
	}
	// invalid language
	if _, err := NewService(repo, trans, Config{SupportedLanguages: []string{"invalid-lang"}}); err == nil {
		t.Fatal("expected error with invalid language")
	}
	// duplicate language and default now/attempts
	svc, err := NewService(repo, trans, Config{
		SupportedLanguages: []string{"en", "EN", "ru"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(svc.SupportedLanguages()) != 2 {
		t.Fatalf("expected 2 unique languages, got %v", svc.SupportedLanguages())
	}

	// invalid MaxAttempts
	if _, err := NewService(repo, trans, Config{MaxAttempts: -1}); err == nil {
		t.Fatal("expected error with MaxAttempts = -1")
	}
	if _, err := NewService(repo, trans, Config{MaxAttempts: 4}); err == nil {
		t.Fatal("expected error with MaxAttempts = 4")
	}
}

func TestTranslateQuestion_ErrorsAndDefaults(t *testing.T) {
	ctx := context.Background()

	// Mock repo returning error
	errSvc := &Service{
		repository: &failingRepo{getQuestionErr: errors.New("db error")},
		translator: &fakeTranslator{},
		supported:  map[string]struct{}{"en": {}, "ru": {}},
		ordered:    []string{"en", "ru"},
	}
	if _, err := errSvc.TranslateQuestion(ctx, TranslateQuestionRequest{QuestionID: "q1", Language: "ru"}); err == nil {
		t.Fatal("expected error from GetQuestion")
	}

	// Default empty sourceLanguage to "en"
	repo := &memoryRepository{
		question:     QuestionSource{ID: "q1", Title: "Title", SourceLanguage: "", ContentRevision: 1, Publication: "published"},
		translations: map[string]Translation{},
	}
	svc, _ := NewService(repo, &fakeTranslator{}, Config{
		SupportedLanguages: []string{"en", "ru"},
		Now:                time.Now,
	})
	res, err := svc.TranslateQuestion(ctx, TranslateQuestionRequest{QuestionID: "q1", Language: "en"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SourceLanguage != "en" {
		t.Fatalf("expected default sourceLanguage 'en', got %s", res.SourceLanguage)
	}
}

type failingRepo struct {
	getQuestionErr error
	putErr         error
}

func (f *failingRepo) GetQuestion(context.Context, string) (QuestionSource, error) {
	return QuestionSource{}, f.getQuestionErr
}

func (f *failingRepo) GetTranslation(context.Context, string, string) (Translation, error) {
	return Translation{}, errors.New("not found")
}

func (f *failingRepo) PutTranslationIfFresh(context.Context, QuestionSource, Translation) (Translation, error) {
	return Translation{}, f.putErr
}

func TestTranslateEnabledQuestion(t *testing.T) {
	ctx := context.Background()
	repo := &memoryRepository{
		question:     QuestionSource{ID: "q1", Title: "Hello", Description: "Desc", SourceLanguage: "en", ContentRevision: 1, Publication: "published"},
		translations: map[string]Translation{},
	}
	svc, _ := NewService(repo, &fakeTranslator{}, Config{
		SupportedLanguages: []string{"en", "ru"},
		Now:                time.Now,
	})

	// Success
	res, err := svc.TranslateEnabledQuestion(ctx, "q1", Actor{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 translations, got %d", len(res))
	}

	// Error when question not found
	errSvc := &Service{
		repository: &failingRepo{getQuestionErr: errors.New("not found")},
		translator: &fakeTranslator{},
		supported:  map[string]struct{}{"en": {}},
		ordered:    []string{"en"},
	}
	if _, err := errSvc.TranslateEnabledQuestion(ctx, "q1", Actor{}); err == nil {
		t.Fatal("expected error when TranslateQuestion fails")
	}
}

func TestTranslateAndStore_RetriesAndFailures(t *testing.T) {
	ctx := context.Background()
	repo := &memoryRepository{
		question:     QuestionSource{ID: "q1", Title: "Hello", Description: "Desc", SourceLanguage: "en", ContentRevision: 1, Publication: "published"},
		translations: map[string]Translation{},
	}

	// Non-retryable error fails immediately
	var attempts int32
	vt := &variableTranslator{
		fn: func(ctx context.Context, source, target string, texts []string) ([]string, error) {
			atomic.AddInt32(&attempts, 1)
			return nil, errors.New("permanent failure")
		},
	}
	svc, _ := NewService(repo, vt, Config{
		SupportedLanguages: []string{"en", "ru"},
		MaxAttempts:        3,
		Now:                time.Now,
	})
	if _, err := svc.TranslateQuestion(ctx, TranslateQuestionRequest{QuestionID: "q1", Language: "ru"}); err == nil {
		t.Fatal("expected error on permanent failure")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt for non-retryable error, got %d", attempts)
	}

	// Retryable error succeeds on retry
	attempts = 0
	repoRetry := &memoryRepository{
		question:     QuestionSource{ID: "q1", Title: "Hello", Description: "Desc", SourceLanguage: "en", ContentRevision: 1, Publication: "published"},
		translations: map[string]Translation{},
	}
	vtSuccessOnRetry := &variableTranslator{
		fn: func(ctx context.Context, source, target string, texts []string) ([]string, error) {
			cur := atomic.AddInt32(&attempts, 1)
			if cur == 1 {
				return nil, &customError{msg: "temporary error", retryable: true}
			}
			return []string{"Привет", "Описание"}, nil
		},
	}
	svcRetry, _ := NewService(repoRetry, vtSuccessOnRetry, Config{
		SupportedLanguages: []string{"en", "ru"},
		MaxAttempts:        3,
		Now:                time.Now,
	})
	tr, err := svcRetry.TranslateQuestion(ctx, TranslateQuestionRequest{QuestionID: "q1", Language: "ru"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Title != "Привет" || attempts != 2 {
		t.Fatalf("expected title 'Привет' after 2 attempts, got %v (%d attempts)", tr, attempts)
	}

	// Retryable error reaches max attempts
	attempts = 0
	repoExhaust := &memoryRepository{
		question:     QuestionSource{ID: "q1", Title: "Hello", Description: "Desc", SourceLanguage: "en", ContentRevision: 1, Publication: "published"},
		translations: map[string]Translation{},
	}
	vtExhaustRetries := &variableTranslator{
		fn: func(ctx context.Context, source, target string, texts []string) ([]string, error) {
			atomic.AddInt32(&attempts, 1)
			return nil, &customError{msg: "temporary error", retryable: true}
		},
	}
	svcExhaust, _ := NewService(repoExhaust, vtExhaustRetries, Config{
		SupportedLanguages: []string{"en", "ru"},
		MaxAttempts:        2,
		Now:                time.Now,
	})
	if _, err := svcExhaust.TranslateQuestion(ctx, TranslateQuestionRequest{QuestionID: "q1", Language: "ru"}); err == nil {
		t.Fatal("expected error when retries exhausted")
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}

	// Invalid result from translator (not enough elements or empty title)
	for _, badRes := range [][]string{
		{"only-one"},
		{"", "desc"},
		{"  ", "desc"},
	} {
		repoBad := &memoryRepository{
			question:     QuestionSource{ID: "q1", Title: "Hello", Description: "Desc", SourceLanguage: "en", ContentRevision: 1, Publication: "published"},
			translations: map[string]Translation{},
		}
		vtBad := &variableTranslator{
			fn: func(ctx context.Context, source, target string, texts []string) ([]string, error) {
				return badRes, nil
			},
		}
		svcBad, _ := NewService(repoBad, vtBad, Config{
			SupportedLanguages: []string{"en", "ru"},
			Now:                time.Now,
		})
		if _, err := svcBad.TranslateQuestion(ctx, TranslateQuestionRequest{QuestionID: "q1", Language: "ru"}); err == nil {
			t.Fatalf("expected error for bad result %v", badRes)
		}
	}

	// PutTranslationIfFresh error
	svcPutErr := &Service{
		repository:  &failingRepo{putErr: errors.New("put failed")},
		translator:  &fakeTranslator{},
		supported:   map[string]struct{}{"en": {}, "ru": {}},
		ordered:     []string{"en", "ru"},
		maxAttempts: 1,
		now:         time.Now,
		inflight:    map[string]*translationCall{},
	}
	if _, err := svcPutErr.translateAndStore(ctx, repo.question, "en", "ru", "hash"); err == nil {
		t.Fatal("expected error from PutTranslationIfFresh")
	}
}

func TestSingleflight_ContextCancellation(t *testing.T) {
	blockChan := make(chan struct{})
	vt := &variableTranslator{
		fn: func(ctx context.Context, source, target string, texts []string) ([]string, error) {
			<-blockChan
			return []string{"RU: title", "RU: desc"}, nil
		},
	}
	repo := &memoryRepository{
		question:     QuestionSource{ID: "q1", Title: "Hello", Description: "Desc", SourceLanguage: "en", ContentRevision: 1, Publication: "published"},
		translations: map[string]Translation{},
	}
	svc, _ := NewService(repo, vt, Config{
		SupportedLanguages: []string{"en", "ru"},
		Now:                time.Now,
	})

	// Start first call in flight
	go func() {
		_, _ = svc.TranslateQuestion(context.Background(), TranslateQuestionRequest{QuestionID: "q1", Language: "ru"})
	}()

	time.Sleep(20 * time.Millisecond)

	// Second call with cancelled context
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.TranslateQuestion(cancelCtx, TranslateQuestionRequest{QuestionID: "q1", Language: "ru"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	close(blockChan)
}

func TestNormalizeLanguage(t *testing.T) {
	if got := normalizeLanguage("e"); got != "" {
		t.Fatalf("expected empty for short language, got %q", got)
	}
	if got := normalizeLanguage("eng"); got != "" {
		t.Fatalf("expected empty for long language, got %q", got)
	}
	if got := normalizeLanguage("12"); got != "" {
		t.Fatalf("expected empty for numeric language, got %q", got)
	}
	if got := normalizeLanguage("e1"); got != "" {
		t.Fatalf("expected empty for alphanumeric language, got %q", got)
	}
	if got := normalizeLanguage(" EN "); got != "en" {
		t.Fatalf("expected 'en', got %q", got)
	}
}

func TestAuthorize_EdgeCases(t *testing.T) {
	draft := QuestionSource{Publication: "draft", CreatorUID: "creator1"}

	// Trusted actor
	if err := authorize(draft, Actor{Trusted: true}); err != nil {
		t.Fatalf("expected nil for trusted actor, got %v", err)
	}
	// Untrusted other user
	if err := authorize(draft, Actor{UID: "other"}); !errors.Is(err, ErrTranslationDenied) {
		t.Fatalf("expected ErrTranslationDenied, got %v", err)
	}
	// Unknown publication status
	archived := QuestionSource{Publication: "archived"}
	if err := authorize(archived, Actor{Trusted: true}); !errors.Is(err, ErrTranslationDenied) {
		t.Fatalf("expected ErrTranslationDenied, got %v", err)
	}
}
