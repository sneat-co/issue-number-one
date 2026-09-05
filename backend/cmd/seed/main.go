package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sneat-co/issue-number-one/backend/publicqa"
	"github.com/sneat-co/issue-number-one/backend/translations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
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

func main() {
	var file, project, confirm, space string
	var translateAll bool
	flag.StringVar(&file, "file", "../catalog/seed.json", "catalog seed JSON")
	flag.StringVar(&project, "project", os.Getenv("GCLOUD_PROJECT"), "Google Cloud project")
	flag.StringVar(&confirm, "confirm-production-project", "", "must exactly equal --project outside emulator")
	flag.StringVar(&space, "space", publicqa.DefaultPublicSpaceID, "public Space id")
	flag.BoolVar(&translateAll, "translate-all", false, "translate every seeded question into the supported language set")
	flag.Parse()
	if project == "" {
		fatal("--project or GCLOUD_PROJECT is required")
	}
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" && confirm != project {
		fatal("production seeding requires exact --confirm-production-project")
	}
	b, e := os.ReadFile(file)
	if e != nil {
		fatal(e.Error())
	}
	var c catalog
	if e = json.Unmarshal(b, &c); e != nil {
		fatal(e.Error())
	}
	if c.SchemaVersion != 1 {
		fatal("unsupported schemaVersion")
	}
	concepts := map[string]publicqa.Concept{}
	for _, v := range c.Concepts {
		concepts[v.ID] = v.Concept
	}
	timeout := 2 * time.Minute
	if translateAll {
		timeout = 15 * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	db, e := firestore.NewClient(ctx, project)
	if e != nil {
		fatal(e.Error())
	}
	defer db.Close()
	root := db.Collection("spaces").Doc(space).Collection("ext").Doc(publicqa.ExtensionID)
	batch := db.Batch()
	n := 0
	flush := func() {
		if n == 0 {
			return
		}
		if _, e := batch.Commit(ctx); e != nil {
			fatal(e.Error())
		}
		batch = db.Batch()
		n = 0
	}
	set := func(ref *firestore.DocumentRef, v map[string]any) {
		batch.Set(ref, v, firestore.MergeAll)
		n++
		if n == 400 {
			flush()
		}
	}
	for _, v := range c.Categories {
		if v.ID == "" {
			fatal("category id required")
		}
		set(root.Collection("categories").Doc(v.ID), map[string]any{"id": v.ID, "slug": v.Slug, "name": v.Name, "question": v.Question, "description": v.Description, "scopeType": v.ScopeType, "expectedChildScopeType": v.ExpectedChildScopeType, "intent": v.Intent, "seoTitle": v.SEOTitle, "seoDescription": v.SEODescription, "status": v.Publication, "publication": v.Publication, "indexable": v.Indexable, "conceptIds": v.DefaultConceptIDs, "prioritySlotId": v.ID, "parentCategoryId": v.ParentCategoryID})
	}
	for _, v := range c.Concepts {
		if v.ID == "" {
			fatal("concept id required")
		}
		set(root.Collection("concepts").Doc(v.ID), map[string]any{"id": v.ID, "slug": v.Slug, "title": v.Title, "description": v.Description, "aliases": v.Aliases, "status": publicqa.StatusPublished})
	}
	for _, q := range c.Questions {
		if q.ID == "" {
			fatal("question id required")
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
		// Scoped questions are resolved from the catalog and may legitimately
		// reuse slugs at different hierarchy levels (county/city Dublin). The
		// global slug registry belongs only to /questions/{slug} community URLs.
		if q.CategoryID == "" {
			set(root.Collection("questionSlugs").Doc(q.Slug), map[string]any{"questionId": q.ID, "seeded": true})
		}
		for _, cid := range q.ConceptIDs {
			v, ok := concepts[cid]
			if !ok {
				fatal("question references unknown concept " + cid)
			}
			iref := root.Collection("questions").Doc(q.ID).Collection("issues").Doc(cid)
			fields := map[string]any{"id": cid, "slug": v.Slug, "title": v.Title, "description": v.Description, "conceptId": cid, "targetType": "issue", "targetRef": cid, "attribution": "anonymous"}
			if _, e := iref.Get(ctx); status.Code(e) == codes.NotFound {
				fields["status"] = publicqa.StatusPublished
			} else if e != nil {
				fatal(e.Error())
			}
			set(iref, fields)
			for _, a := range append([]string{v.Title}, v.Aliases...) {
				norm, e := publicqa.NormalizeTitle(a)
				if e != nil {
					fatal(e.Error())
				}
				set(root.Collection("questions").Doc(q.ID).Collection("aliases").Doc(hash(norm)), map[string]any{"issueId": cid, "normalized": norm})
			}
		}
	}
	flush()
	fmt.Printf("seeded %d categories, %d concepts, %d questions; operational counts and moderation fields preserved\n", len(c.Categories), len(c.Concepts), len(c.Questions))
	if translateAll {
		translationRepository, err := translations.NewFirestoreRepository(db, space)
		if err != nil {
			fatal(err.Error())
		}
		translator, err := translations.NewGoogleTranslator(ctx)
		if err != nil {
			fatal(err.Error())
		}
		defer translator.Close()
		translationService, err := translations.NewService(translationRepository, translator, translations.Config{SupportedLanguages: translations.CanonicalSupportedLanguages()})
		if err != nil {
			fatal(err.Error())
		}
		questionService := publicqa.NewService(db, space, publicqa.WithTranslations(translationService, translationRepository))
		for _, question := range c.Questions {
			if _, err = questionService.TranslateAllQuestionLanguages(ctx, question.ID, translations.Actor{Trusted: true}); err != nil {
				fatal(fmt.Sprintf("translate question %s: %v", question.ID, err))
			}
		}
		fmt.Printf("translated %d questions into %d supported languages\n", len(c.Questions), len(translations.CanonicalSupportedLanguages()))
	}
}
func hash(s string) string { return publicqa.NormalizedKey(s) }
func fatal(s string)       { fmt.Fprintln(os.Stderr, s); os.Exit(1) }
