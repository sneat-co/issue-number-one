package publicqa

import (
	"github.com/sneat-co/issue-number-one/backend/translations"
	"time"
)

const (
	DefaultPublicSpaceID = "issuenumber-public"
	ExtensionID          = "issuenumber"
	StatusPublished      = "published"
	StatusPending        = "pending"
	AnswerKindCategory   = "category"
	AnswerKindPersonal   = "personal"
)

type Category struct {
	ID                     string   `json:"id" firestore:"id"`
	Slug                   string   `json:"slug" firestore:"slug"`
	Name                   string   `json:"name" firestore:"name"`
	Question               string   `json:"question" firestore:"question"`
	Description            string   `json:"description" firestore:"description"`
	ScopeType              string   `json:"scopeType" firestore:"scopeType"`
	ExpectedChildScopeType string   `json:"expectedChildScopeType,omitempty" firestore:"expectedChildScopeType,omitempty"`
	Intent                 string   `json:"intent,omitempty" firestore:"intent,omitempty"`
	SEOTitle               string   `json:"seoTitle" firestore:"seoTitle"`
	SEODescription         string   `json:"seoDescription" firestore:"seoDescription"`
	Status                 string   `json:"status" firestore:"status"`
	Indexable              bool     `json:"indexable" firestore:"indexable"`
	ConceptIDs             []string `json:"conceptIds,omitempty" firestore:"conceptIds,omitempty"`
	PrioritySlotID         string   `json:"prioritySlotId,omitempty" firestore:"prioritySlotId,omitempty"`
	ParentCategoryID       string   `json:"parentCategoryId,omitempty" firestore:"parentCategoryId,omitempty"`
	Publication            string   `json:"publication,omitempty" firestore:"publication,omitempty"`
}
type Concept struct {
	ID          string   `json:"id" firestore:"id"`
	Slug        string   `json:"slug" firestore:"slug"`
	Title       string   `json:"title" firestore:"title"`
	Description string   `json:"description,omitempty" firestore:"description,omitempty"`
	Aliases     []string `json:"aliases,omitempty" firestore:"aliases,omitempty"`
	Status      string   `json:"status" firestore:"status"`
}
type Question struct {
	ID                    string           `json:"id" firestore:"id"`
	Slug                  string           `json:"slug" firestore:"slug"`
	Title                 string           `json:"title" firestore:"title"`
	Description           string           `json:"description" firestore:"description"`
	CategoryID            string           `json:"categoryId" firestore:"categoryId"`
	PrioritySlotID        string           `json:"prioritySlotId,omitempty" firestore:"prioritySlotId,omitempty"`
	ScopeType             string           `json:"scopeType" firestore:"scopeType"`
	ScopeID               string           `json:"scopeId,omitempty" firestore:"scopeId,omitempty"`
	ScopeName             string           `json:"scopeName,omitempty" firestore:"scopeName,omitempty"`
	ParentQuestionID      string           `json:"parentQuestionId,omitempty" firestore:"parentQuestionId,omitempty"`
	ConceptIDs            []string         `json:"conceptIds,omitempty" firestore:"conceptIds,omitempty"`
	RelatedQuestionIDs    []string         `json:"relatedQuestionIds,omitempty" firestore:"relatedQuestionIds,omitempty"`
	Status                string           `json:"status" firestore:"status"`
	Indexable             bool             `json:"indexable" firestore:"indexable"`
	TotalRespondents      int64            `json:"-" firestore:"totalRespondents"`
	UpdatedAt             time.Time        `json:"-" firestore:"updatedAt"`
	AnswerTargetType      string           `json:"answerTargetType,omitempty" firestore:"answerTargetType,omitempty"`
	CreatorUID            string           `json:"-" firestore:"creatorUID,omitempty"`
	Scope                 Scope            `json:"scope" firestore:"scope"`
	Publication           string           `json:"publication,omitempty" firestore:"publication,omitempty"`
	MutationLock          string           `json:"-" firestore:"mutationLock,omitempty"`
	RespondentsByLanguage map[string]int64 `json:"-" firestore:"respondentsByLanguage,omitempty"`
	ChoiceSource          ChoiceSource     `json:"choiceSource" firestore:"choiceSource"`
	AllowSuggestions      bool             `json:"allowSuggestions" firestore:"allowSuggestions"`
	ContentRevision       int64            `json:"contentRevision" firestore:"contentRevision"`
	SourceLanguage        string           `json:"sourceLanguage" firestore:"sourceLanguage"`
	TranslationStatus     string           `json:"translationStatus,omitempty" firestore:"translationStatus,omitempty"`
	AvailableLanguages    []string         `json:"availableLanguages,omitempty" firestore:"availableLanguages,omitempty"`
}
type Scope struct {
	ID       string `json:"id" firestore:"id"`
	Type     string `json:"type" firestore:"type"`
	Name     string `json:"name" firestore:"name"`
	ParentID string `json:"parentId,omitempty" firestore:"parentId,omitempty"`
}
type ChoiceSource struct {
	Kind       string         `json:"kind" firestore:"kind"`
	EntityType string         `json:"entityType,omitempty" firestore:"entityType,omitempty"`
	Options    []ChoiceOption `json:"options,omitempty" firestore:"-"`
}
type ChoiceOption struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}
type Issue struct {
	ID                            string                       `json:"id" firestore:"id"`
	Slug                          string                       `json:"slug" firestore:"slug"`
	Title                         string                       `json:"title" firestore:"title"`
	Description                   string                       `json:"description,omitempty" firestore:"description,omitempty"`
	ConceptID                     string                       `json:"conceptId,omitempty" firestore:"conceptId,omitempty"`
	Status                        string                       `json:"status" firestore:"status"`
	Supporters                    int64                        `json:"supporters" firestore:"supporters"`
	PersonalTopSupporters         int64                        `json:"personalTopSupporters" firestore:"personalTopSupporters"`
	WeightedScore                 int64                        `json:"weightedScore" firestore:"weightedScore"`
	Attribution                   string                       `json:"attribution,omitempty" firestore:"attribution,omitempty"`
	AuthorDisplayName             string                       `json:"authorDisplayName,omitempty" firestore:"authorDisplayName,omitempty"`
	CreatorUID                    string                       `json:"-" firestore:"creatorUID,omitempty"`
	MergedIntoIssueID             string                       `json:"mergedIntoIssueId,omitempty" firestore:"mergedIntoIssueId,omitempty"`
	CreatedAt                     time.Time                    `json:"-" firestore:"createdAt,omitempty"`
	UpdatedAt                     time.Time                    `json:"-" firestore:"updatedAt,omitempty"`
	TargetType                    string                       `json:"targetType,omitempty" firestore:"targetType,omitempty"`
	TargetRef                     string                       `json:"targetRef,omitempty" firestore:"targetRef,omitempty"`
	LanguageSupporters            int64                        `json:"languageSupporters,omitempty" firestore:"-"`
	LanguagePersonalTopSupporters int64                        `json:"languagePersonalTopSupporters,omitempty" firestore:"-"`
	LanguageWeightedScore         int64                        `json:"languageWeightedScore,omitempty" firestore:"-"`
	LanguageStats                 map[string]LanguageAggregate `json:"-" firestore:"languageStats,omitempty"`
}
type LanguageAggregate struct {
	Supporters            int64 `firestore:"supporters"`
	PersonalTopSupporters int64 `firestore:"personalTopSupporters"`
	WeightedScore         int64 `firestore:"weightedScore"`
}
type Answer struct {
	CategoryID     string    `json:"categoryId" firestore:"categoryId"`
	PrioritySlotID string    `json:"prioritySlotId" firestore:"prioritySlotId"`
	QuestionID     string    `json:"questionId" firestore:"questionId"`
	IssueID        string    `json:"issueId" firestore:"issueId"`
	Revision       int64     `json:"revision" firestore:"revision"`
	CreatedAt      time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" firestore:"updatedAt"`
	LanguageCode   string    `json:"languageCode" firestore:"languageCode"`
}
type PersonalAnswer struct {
	CategoryID     string    `json:"categoryId" firestore:"categoryId"`
	PrioritySlotID string    `json:"prioritySlotId" firestore:"prioritySlotId"`
	QuestionID     string    `json:"questionId" firestore:"questionId"`
	IssueID        string    `json:"issueId" firestore:"issueId"`
	UpdatedAt      time.Time `json:"updatedAt" firestore:"updatedAt"`
	LanguageCode   string    `json:"languageCode" firestore:"languageCode"`
}
type CatalogResponse struct {
	Categories []Category `json:"categories"`
	Concepts   []Concept  `json:"concepts"`
	Questions  []Question `json:"questions"`
}
type QuestionResponse struct {
	Question            Question                  `json:"question"`
	Issues              []Issue                   `json:"issues"`
	TotalRespondents    int64                     `json:"totalRespondents"`
	UpdatedAt           time.Time                 `json:"updatedAt,omitempty"`
	LanguageCode        string                    `json:"languageCode,omitempty"`
	LanguageRespondents int64                     `json:"languageRespondents,omitempty"`
	ContentLanguage     string                    `json:"contentLanguage"`
	TranslationFallback bool                      `json:"translationFallback,omitempty"`
	Translation         *translations.Translation `json:"translation,omitempty"`
}
type AnswerRequest struct {
	AnswerKind   string `json:"answerKind,omitempty"`
	QuestionID   string `json:"questionId"`
	IssueID      string `json:"issueId,omitempty"`
	Title        string `json:"title,omitempty"`
	OperationID  string `json:"operationId"`
	Attribution  string `json:"attribution,omitempty"`
	LanguageCode string `json:"languageCode,omitempty"`
}
type Caller struct {
	UID           string
	PhoneVerified bool
	DisplayName   string
}
type AnswerResponse struct {
	Answer       Answer `json:"answer"`
	Issue        Issue  `json:"issue"`
	Changed      bool   `json:"changed"`
	PersonalTop  bool   `json:"personalTop"`
	CreatedIssue bool   `json:"createdIssue"`
	Replayed     bool   `json:"replayed"`
}
type CreateQuestionRequest struct {
	Title            string       `json:"title"`
	Description      string       `json:"description"`
	AnswerTargetType string       `json:"answerTargetType"`
	OperationID      string       `json:"operationId"`
	ChoiceSource     ChoiceSource `json:"choiceSource"`
	AllowSuggestions bool         `json:"allowSuggestions,omitempty"`
	SourceLanguage   string       `json:"sourceLanguage,omitempty"`
}
type CreateQuestionResponse struct {
	Question Question `json:"question"`
	Replayed bool     `json:"replayed"`
}
