package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/sneat-co/issue-number-one/backend/publicqa"
	"github.com/sneat-co/issue-number-one/backend/translations"
	"google.golang.org/api/option"
)

type countryFile struct {
	Countries map[string]struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"countriesByID"`
}

type adminTranslator interface {
	translations.Translator
	Close() error
}

var (
	exit                = os.Exit
	fatal               = func(v string) { fmt.Fprintln(os.Stderr, v); exit(1) }
	newFirestoreClient  = firestore.NewClient
	newGoogleTranslator = func(ctx context.Context, opts ...option.ClientOption) (adminTranslator, error) {
		return translations.NewGoogleTranslator(ctx, opts...)
	}
	newTranslationService = translations.NewService
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fatal(err.Error())
	}
}

func run(args []string) error {
	var action, project, confirm, space, qid, iid, target, op, file string
	fs := flag.NewFlagSet("admin", flag.ContinueOnError)
	fs.StringVar(&action, "action", "", "publish-question|hide-question|publish-issue|hide-issue|reject-issue|populate-countries|merge-issue|translate-question")
	fs.StringVar(&project, "project", os.Getenv("GCLOUD_PROJECT"), "Google Cloud project")
	fs.StringVar(&confirm, "confirm-production-project", "", "must exactly equal project outside emulator")
	fs.StringVar(&space, "space", publicqa.DefaultPublicSpaceID, "public Space")
	fs.StringVar(&qid, "question", "", "question id")
	fs.StringVar(&iid, "issue", "", "issue id")
	fs.StringVar(&target, "target-issue", "", "merge target issue id")
	fs.StringVar(&op, "operation", "", "stable merge operation id")
	fs.StringVar(&file, "countries-file", "", "Sneat Libs countries.json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if project == "" || qid == "" {
		return errors.New("--project and --question required")
	}
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" && confirm != project {
		return errors.New("production administration requires exact --confirm-production-project")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	db, e := newFirestoreClient(ctx, project)
	if e != nil {
		return e
	}
	defer func() { _ = db.Close() }()
	s := publicqa.NewService(db, space)
	switch action {
	case "publish-question":
		e = s.SetQuestionStatus(ctx, qid, publicqa.StatusPublished)
	case "hide-question":
		e = s.SetQuestionStatus(ctx, qid, "hidden")
	case "publish-issue":
		e = s.SetIssueStatus(ctx, qid, iid, publicqa.StatusPublished)
	case "hide-issue":
		e = s.SetIssueStatus(ctx, qid, iid, "hidden")
	case "reject-issue":
		e = s.SetIssueStatus(ctx, qid, iid, "rejected")
	case "merge-issue":
		e = s.MergeIssue(ctx, qid, iid, target, op)
	case "populate-countries":
		if file == "" {
			return errors.New("--countries-file must point to Sneat Libs countries.json")
		}
		b, x := os.ReadFile(file)
		if x != nil {
			return x
		}
		var src countryFile
		if x = json.Unmarshal(b, &src); x != nil {
			return x
		}
		ids := make([]string, 0, len(src.Countries))
		for id := range src.Countries {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		choices := make([]publicqa.PredefinedChoice, 0, len(ids))
		for _, id := range ids {
			c := src.Countries[id]
			choices = append(choices, publicqa.PredefinedChoice{ID: c.ID, Title: c.Title})
		}
		e = s.PopulatePredefinedChoices(ctx, qid, "country", choices)
	case "translate-question":
		repository, x := translations.NewFirestoreRepository(db, space)
		if x != nil {
			return x
		}
		translator, x := newGoogleTranslator(ctx)
		if x != nil {
			return x
		}
		defer func() { _ = translator.Close() }()
		translationService, x := newTranslationService(repository, translator, translations.Config{SupportedLanguages: translations.CanonicalSupportedLanguages()})
		if x != nil {
			return x
		}
		s = publicqa.NewService(db, space, publicqa.WithTranslations(translationService, repository))
		_, e = s.TranslateAllQuestionLanguages(ctx, qid, translations.Actor{Trusted: true})
	default:
		return errors.New("invalid --action")
	}
	if e != nil {
		return e
	}
	fmt.Println("admin operation completed")
	return nil
}
