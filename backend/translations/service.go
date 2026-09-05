package translations

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Service struct {
	repository  Repository
	translator  Translator
	supported   map[string]struct{}
	ordered     []string
	maxAttempts int
	now         func() time.Time
	mu          sync.Mutex
	inflight    map[string]*translationCall
}

type translationCall struct {
	done  chan struct{}
	value Translation
	err   error
}

var canonicalSupportedLanguages = []string{"en", "ru", "uk", "pl", "de", "fr", "es", "it", "pt", "nl", "ga"}

// CanonicalSupportedLanguages returns the product-wide bounded translation
// language set. Returning a copy prevents callers from mutating shared config.
func CanonicalSupportedLanguages() []string {
	return append([]string(nil), canonicalSupportedLanguages...)
}

func NewService(repository Repository, translator Translator, config Config) (*Service, error) {
	if repository == nil || translator == nil {
		return nil, errors.New("translation repository and translator are required")
	}
	languages := config.SupportedLanguages
	if len(languages) == 0 {
		languages = CanonicalSupportedLanguages()
	}
	supported := make(map[string]struct{}, len(languages))
	ordered := make([]string, 0, len(languages))
	for _, raw := range languages {
		language := normalizeLanguage(raw)
		if language == "" {
			return nil, fmt.Errorf("%w: %q", ErrUnsupportedLanguage, raw)
		}
		if _, exists := supported[language]; exists {
			continue
		}
		supported[language] = struct{}{}
		ordered = append(ordered, language)
	}
	attempts := config.MaxAttempts
	if attempts == 0 {
		attempts = 3
	}
	if attempts < 1 || attempts > 3 {
		return nil, errors.New("translation max attempts must be between 1 and 3")
	}
	now := config.Now
	if now == nil {
		now = time.Now
	}
	return &Service{repository: repository, translator: translator, supported: supported, ordered: ordered, maxAttempts: attempts, now: now, inflight: map[string]*translationCall{}}, nil
}

func (s *Service) SupportedLanguages() []string { return append([]string(nil), s.ordered...) }

func (s *Service) TranslateQuestion(ctx context.Context, request TranslateQuestionRequest) (Translation, error) {
	language := normalizeLanguage(request.Language)
	if _, ok := s.supported[language]; !ok {
		return Translation{}, fmt.Errorf("%w: %q", ErrUnsupportedLanguage, request.Language)
	}
	question, err := s.repository.GetQuestion(ctx, request.QuestionID)
	if err != nil {
		return Translation{}, err
	}
	if err := authorize(question, request.Actor); err != nil {
		return Translation{}, err
	}
	sourceLanguage := normalizeLanguage(question.SourceLanguage)
	if sourceLanguage == "" {
		sourceLanguage = "en"
	}
	hash := SourceHash(question, sourceLanguage)
	if language == sourceLanguage {
		return Translation{Language: language, Title: question.Title, Description: question.Description, SourceLanguage: sourceLanguage, SourceRevision: question.ContentRevision, SourceHash: hash, MachineTranslated: false}, nil
	}
	if existing, err := s.repository.GetTranslation(ctx, question.ID, language); err == nil && isFresh(existing, question, hash) {
		return existing, nil
	}
	key := question.ID + "\x00" + language + "\x00" + fmt.Sprint(question.ContentRevision) + "\x00" + hash
	return s.singleflight(ctx, key, func() (Translation, error) { return s.translateAndStore(ctx, question, sourceLanguage, language, hash) })
}

func (s *Service) TranslateEnabledQuestion(ctx context.Context, questionID string, actor Actor) ([]Translation, error) {
	result := make([]Translation, 0, len(s.ordered))
	for _, language := range s.ordered {
		translated, err := s.TranslateQuestion(ctx, TranslateQuestionRequest{QuestionID: questionID, Language: language, Actor: actor})
		if err != nil {
			return nil, err
		}
		result = append(result, translated)
	}
	return result, nil
}

func (s *Service) translateAndStore(ctx context.Context, question QuestionSource, sourceLanguage, language, hash string) (Translation, error) {
	var values []string
	var err error
	for attempt := 1; attempt <= s.maxAttempts; attempt++ {
		values, err = s.translator.Translate(ctx, sourceLanguage, language, []string{question.Title, question.Description})
		if err == nil {
			break
		}
		if !isRetryable(err) || attempt == s.maxAttempts {
			return Translation{}, fmt.Errorf("translate question %q to %s: %w", question.ID, language, err)
		}
	}
	if len(values) != 2 || strings.TrimSpace(values[0]) == "" {
		return Translation{}, errors.New("translation provider returned an invalid result")
	}
	now := s.now().UTC()
	translated := Translation{Language: language, Title: strings.TrimSpace(values[0]), Description: strings.TrimSpace(values[1]), SourceLanguage: sourceLanguage, SourceRevision: question.ContentRevision, SourceHash: hash, MachineTranslated: true, CreatedAt: now, UpdatedAt: now}
	stored, err := s.repository.PutTranslationIfFresh(ctx, question, translated)
	if err != nil {
		return Translation{}, err
	}
	return stored, nil
}

func (s *Service) singleflight(ctx context.Context, key string, work func() (Translation, error)) (Translation, error) {
	s.mu.Lock()
	if call := s.inflight[key]; call != nil {
		s.mu.Unlock()
		select {
		case <-ctx.Done():
			return Translation{}, ctx.Err()
		case <-call.done:
			return call.value, call.err
		}
	}
	call := &translationCall{done: make(chan struct{})}
	s.inflight[key] = call
	s.mu.Unlock()
	call.value, call.err = work()
	s.mu.Lock()
	delete(s.inflight, key)
	close(call.done)
	s.mu.Unlock()
	return call.value, call.err
}

func SourceHash(question QuestionSource, sourceLanguage string) string {
	sum := sha256.Sum256([]byte(normalizeLanguage(sourceLanguage) + "\x00" + question.Title + "\x00" + question.Description))
	return hex.EncodeToString(sum[:])
}

func normalizeLanguage(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if len(value) != 2 {
		return ""
	}
	for _, char := range value {
		if char < 'a' || char > 'z' {
			return ""
		}
	}
	return value
}

func authorize(question QuestionSource, actor Actor) error {
	switch question.Publication {
	case "published":
		return nil
	case "pending", "draft":
		if actor.Trusted || (actor.UID != "" && actor.UID == question.CreatorUID) {
			return nil
		}
	}
	return ErrTranslationDenied
}

func isFresh(value Translation, question QuestionSource, hash string) bool {
	return value.Language != "" && value.SourceRevision == question.ContentRevision && value.SourceHash == hash
}

type retryable interface{ Retryable() bool }

func isRetryable(err error) bool {
	var value retryable
	return errors.As(err, &value) && value.Retryable()
}
