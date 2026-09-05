// Package translations stores bounded, revision-aware machine translations of
// IssueNumber.one question copy. Question identity and answers remain language independent.
package translations

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUnsupportedLanguage = errors.New("unsupported translation language")
	ErrQuestionNotFound    = errors.New("question not found")
	ErrTranslationDenied   = errors.New("translation generation is not allowed")
	ErrStaleSource         = errors.New("question source changed during translation")
)

type QuestionSource struct {
	ID              string
	Title           string
	Description     string
	SourceLanguage  string
	ContentRevision int64
	Publication     string
	CreatorUID      string
}

type Translation struct {
	Language          string    `json:"language" firestore:"language"`
	Title             string    `json:"title" firestore:"title"`
	Description       string    `json:"description,omitempty" firestore:"description,omitempty"`
	SourceLanguage    string    `json:"sourceLanguage" firestore:"sourceLanguage"`
	SourceRevision    int64     `json:"sourceRevision" firestore:"sourceRevision"`
	SourceHash        string    `json:"sourceHash" firestore:"sourceHash"`
	MachineTranslated bool      `json:"machineTranslated" firestore:"machineTranslated"`
	CorrectedByUID    string    `json:"-" firestore:"correctedByUID,omitempty"`
	CreatedAt         time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt" firestore:"updatedAt"`
}

type Actor struct {
	UID     string
	Trusted bool
}

type Repository interface {
	GetQuestion(ctx context.Context, questionID string) (QuestionSource, error)
	GetTranslation(ctx context.Context, questionID, language string) (Translation, error)
	PutTranslationIfFresh(ctx context.Context, question QuestionSource, translation Translation) (Translation, error)
}

// Translator receives plain strings as data. Implementations must not interpolate
// question text into provider instructions or prompts.
type Translator interface {
	Translate(ctx context.Context, sourceLanguage, targetLanguage string, texts []string) ([]string, error)
}

type TranslateQuestionRequest struct {
	QuestionID string
	Language   string
	Actor      Actor
}

type Config struct {
	SupportedLanguages []string
	MaxAttempts        int
	Now                func() time.Time
}
