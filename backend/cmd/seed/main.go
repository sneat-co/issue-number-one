package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/sneat-co/issue-number-one/backend/publicqa"
	"github.com/sneat-co/issue-number-one/backend/translations"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type catalog struct {
	SchemaVersion int            `json:"schemaVersion"`
	Categories    []seedCategory `json:"categories"`
	Concepts      []seedConcept  `json:"concepts"`
	Questions     []seedQuestion `json:"questions"`
}
type seedCategory struct {
	publicqa.Category
	Publication       string   `json:"publication"`
	DefaultConceptIDs []string `json:"defaultConceptIds"`
	ParentCategoryID  string   `json:"parentCategoryId"`
}
type seedConcept struct{ publicqa.Concept }
type seedScope struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	ParentID string `json:"parentId"`
}
type seedQuestion struct {
	publicqa.Question
	Publication string    `json:"publication"`
	Scope       seedScope `json:"scope"`
}

type seedTranslator interface {
	translations.Translator
	Close() error
}

var (
	exit                = os.Exit
	fatal               = func(s string) { fmt.Fprintln(os.Stderr, s); exit(1) }
	newFirestoreClient  = firestore.NewClient
	newTranslationRepo  = translations.NewFirestoreRepository
	newGoogleTranslator = func(ctx context.Context, opts ...option.ClientOption) (seedTranslator, error) {
		return translations.NewGoogleTranslator(ctx, opts...)
	}
	newTranslationService = translations.NewService
	newContext            = context.Background
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fatal(err.Error())
	}
}

func run(args []string) error {
	var file, project, confirm, space string
	var translateAll bool
	fs := flag.NewFlagSet("seed", flag.ContinueOnError)
	fs.StringVar(&file, "file", "../catalog/seed.json", "catalog seed JSON")
	fs.StringVar(&project, "project", os.Getenv("GCLOUD_PROJECT"), "Google Cloud project")
	fs.StringVar(&confirm, "confirm-production-project", "", "must exactly equal --project outside emulator")
	fs.StringVar(&space, "space", publicqa.DefaultPublicSpaceID, "public Space id")
	fs.BoolVar(&translateAll, "translate-all", false, "translate every seeded question into the supported language set")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if project == "" {
		return errors.New("--project or GCLOUD_PROJECT is required")
	}
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" && confirm != project {
		return errors.New("production seeding requires exact --confirm-production-project")
	}
	b, e := os.ReadFile(file)
	if e != nil {
		return e
	}
	var c catalog
	if e = json.Unmarshal(b, &c); e != nil {
		return e
	}
	if c.SchemaVersion != 1 {
		return errors.New("unsupported schemaVersion")
	}
	concepts := map[string]publicqa.Concept{}
	for _, v := range c.Concepts {
		concepts[v.ID] = v.Concept
	}
	timeout := 2 * time.Minute
	if translateAll {
		timeout = 15 * time.Minute
	}
	ctx, cancel := context.WithTimeout(newContext(), timeout)
	defer cancel()
	db, e := newFirestoreClient(ctx, project)
	if e != nil {
		return e
	}
	defer func() { _ = db.Close() }()
	root := db.Collection("spaces").Doc(space).Collection("ext").Doc(publicqa.ExtensionID)
	batch := db.Batch() //nolint:staticcheck
	n := 0
	var flushErr error
	flush := func() error {
		if flushErr != nil {
			return flushErr
		}
		if n == 0 {
			return nil
		}
		if _, e := batch.Commit(ctx); e != nil {
			flushErr = e
			return e
		}
		batch = db.Batch() //nolint:staticcheck
		n = 0
		return nil
	}
	set := func(ref *firestore.DocumentRef, v map[string]any) {
		batch.Set(ref, v, firestore.MergeAll)
		n++
		if n == 400 {
			_ = flush()
		}
	}
	for _, v := range c.Categories {
		if v.ID == "" {
			return errors.New("category id required")
		}
		set(root.Collection("categories").Doc(v.ID), map[string]any{"id": v.ID, "slug": v.Slug, "name": v.Name, "question": v.Question, "description": v.Description, "scopeType": v.ScopeType, "expectedChildScopeType": v.ExpectedChildScopeType, "intent": v.Intent, "seoTitle": v.SEOTitle, "seoDescription": v.SEODescription, "status": v.Publication, "publication": v.Publication, "indexable": v.Indexable, "conceptIds": v.DefaultConceptIDs, "prioritySlotId": v.ID, "parentCategoryId": v.ParentCategoryID})
	}
	for _, v := range c.Concepts {
		if v.ID == "" {
			return errors.New("concept id required")
		}
		set(root.Collection("concepts").Doc(v.ID), map[string]any{"id": v.ID, "slug": v.Slug, "title": v.Title, "description": v.Description, "aliases": v.Aliases, "status": publicqa.StatusPublished})
	}
	for _, q := range c.Questions {
		if q.ID == "" {
			return errors.New("question id required")
		}
		choiceSource := q.ChoiceSource
		answerTargetType := q.AnswerTargetType
		allowSuggestions := q.AllowSuggestions
		if choiceSource.Kind == "" {
			choiceSource = publicqa.ChoiceSource{Kind: "free"}
		}
		if choiceSource.Kind == "free" {
			answerTargetType = "issue"
			allowSuggestions = true
		} else if answerTargetType == "" {
			answerTargetType = choiceSource.EntityType
		}
		set(root.Collection("questions").Doc(q.ID), map[string]any{"id": q.ID, "slug": q.Slug, "title": q.Title, "description": q.Description, "categoryId": q.CategoryID, "prioritySlotId": q.ID, "scopeType": q.Scope.Type, "scopeId": q.Scope.ID, "scopeName": q.Scope.Name, "scopeParentId": q.Scope.ParentID, "scope": map[string]any{"id": q.Scope.ID, "type": q.Scope.Type, "name": q.Scope.Name, "parentId": q.Scope.ParentID}, "parentQuestionId": q.ParentQuestionID, "conceptIds": q.ConceptIDs, "relatedQuestionIds": q.RelatedQuestionIDs, "status": q.Publication, "publication": q.Publication, "indexable": q.Indexable, "choiceSource": choiceSource, "answerTargetType": answerTargetType, "allowSuggestions": allowSuggestions, "contentRevision": int64(1), "sourceLanguage": "en", "translationStatus": "pending", "availableLanguages": []string{"en"}})
		if q.CategoryID == "" {
			set(root.Collection("questionSlugs").Doc(q.Slug), map[string]any{"questionId": q.ID, "seeded": true})
		}
		for _, cid := range q.ConceptIDs {
			v, ok := concepts[cid]
			if !ok {
				return errors.New("question references unknown concept " + cid)
			}
			iref := root.Collection("questions").Doc(q.ID).Collection("issues").Doc(cid)
			fields := map[string]any{"id": cid, "slug": v.Slug, "title": v.Title, "description": v.Description, "conceptId": cid, "targetType": "issue", "targetRef": cid, "attribution": "anonymous"}
			if _, e := iref.Get(ctx); status.Code(e) == codes.NotFound {
				fields["status"] = publicqa.StatusPublished
			} else if e != nil {
				return e
			}
			set(iref, fields)
			for _, a := range append([]string{v.Title}, v.Aliases...) {
				norm, e := publicqa.NormalizeTitle(a)
				if e != nil {
					return e
				}
				set(root.Collection("questions").Doc(q.ID).Collection("aliases").Doc(hash(norm)), map[string]any{"issueId": cid, "normalized": norm})
			}
		}
	}
	if err := flush(); err != nil {
		return err
	}
	fmt.Printf("seeded %d categories, %d concepts, %d questions; operational counts and moderation fields preserved\n", len(c.Categories), len(c.Concepts), len(c.Questions))
	if translateAll {
		translationRepository, err := newTranslationRepo(db, space)
		if err != nil {
			return err
		}
		translator, err := newGoogleTranslator(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = translator.Close() }()
		translationService, err := newTranslationService(translationRepository, translator, translations.Config{SupportedLanguages: translations.CanonicalSupportedLanguages()})
		if err != nil {
			return err
		}
		questionService := publicqa.NewService(db, space, publicqa.WithTranslations(translationService, translationRepository))
		for _, question := range c.Questions {
			if _, err = questionService.TranslateAllQuestionLanguages(ctx, question.ID, translations.Actor{Trusted: true}); err != nil {
				return fmt.Errorf("translate question %s: %w", question.ID, err)
			}
		}
		fmt.Printf("translated %d questions into %d supported languages\n", len(c.Questions), len(translations.CanonicalSupportedLanguages()))
	}
	return nil
}

func hash(s string) string { return publicqa.NormalizedKey(s) }
