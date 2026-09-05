package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sneat-co/issue-number-one/backend/publicqa"
	"github.com/sneat-co/issue-number-one/backend/translations"
	"os"
	"sort"
	"time"
)

type countryFile struct {
	Countries map[string]struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"countriesByID"`
}

func main() {
	var action, project, confirm, space, qid, iid, target, op, file string
	flag.StringVar(&action, "action", "", "publish-question|hide-question|publish-issue|hide-issue|reject-issue|populate-countries|merge-issue|translate-question")
	flag.StringVar(&project, "project", os.Getenv("GCLOUD_PROJECT"), "Google Cloud project")
	flag.StringVar(&confirm, "confirm-production-project", "", "must exactly equal project outside emulator")
	flag.StringVar(&space, "space", publicqa.DefaultPublicSpaceID, "public Space")
	flag.StringVar(&qid, "question", "", "question id")
	flag.StringVar(&iid, "issue", "", "issue id")
	flag.StringVar(&target, "target-issue", "", "merge target issue id")
	flag.StringVar(&op, "operation", "", "stable merge operation id")
	flag.StringVar(&file, "countries-file", "", "Sneat Libs countries.json")
	flag.Parse()
	if project == "" || qid == "" {
		fatal("--project and --question required")
	}
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" && confirm != project {
		fatal("production administration requires exact --confirm-production-project")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	db, e := firestore.NewClient(ctx, project)
	if e != nil {
		fatal(e.Error())
	}
	defer db.Close()
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
			fatal("--countries-file must point to Sneat Libs countries.json")
		}
		b, x := os.ReadFile(file)
		if x != nil {
			fatal(x.Error())
		}
		var src countryFile
		if x = json.Unmarshal(b, &src); x != nil {
			fatal(x.Error())
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
			fatal(x.Error())
		}
		translator, x := translations.NewGoogleTranslator(ctx)
		if x != nil {
			fatal(x.Error())
		}
		defer translator.Close()
		translationService, x := translations.NewService(repository, translator, translations.Config{SupportedLanguages: translations.CanonicalSupportedLanguages()})
		if x != nil {
			fatal(x.Error())
		}
		s = publicqa.NewService(db, space, publicqa.WithTranslations(translationService, repository))
		_, e = s.TranslateAllQuestionLanguages(ctx, qid, translations.Actor{Trusted: true})
	default:
		fatal("invalid --action")
	}
	if e != nil {
		fatal(e.Error())
	}
	fmt.Println("admin operation completed")
}
func fatal(v string) { fmt.Fprintln(os.Stderr, v); os.Exit(1) }
