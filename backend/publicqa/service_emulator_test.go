package publicqa

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
)

func emulatorService(t *testing.T) *Service {
	t.Helper()
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	ctx := context.Background()
	db, e := firestore.NewClient(ctx, "issuenumber-public-test")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { _ = db.Close() })
	return NewService(db, fmt.Sprintf("test-%d", time.Now().UnixNano()))
}
func seedQuestion(t *testing.T, s *Service, qid, cat string, issues ...string) {
	t.Helper()
	ctx := context.Background()
	q := Question{ID: qid, CategoryID: cat, Status: StatusPublished, AnswerTargetType: "issue", ChoiceSource: ChoiceSource{Kind: "free"}, AllowSuggestions: true}
	if _, e := s.root().Collection("questions").Doc(qid).Set(ctx, q); e != nil {
		t.Fatal(e)
	}
	for _, id := range issues {
		i := Issue{ID: id, Slug: id, Title: id, Status: StatusPublished}
		if _, e := s.root().Collection("questions").Doc(qid).Collection("issues").Doc(id).Set(ctx, i); e != nil {
			t.Fatal(e)
		}
	}
}
func seedQuestionSlot(t *testing.T, s *Service, qid, cat, slot string, issues ...string) {
	t.Helper()
	ctx := context.Background()
	q := Question{ID: qid, CategoryID: cat, PrioritySlotID: slot, Status: StatusPublished, AnswerTargetType: "issue", ChoiceSource: ChoiceSource{Kind: "free"}, AllowSuggestions: true}
	if _, e := s.root().Collection("questions").Doc(qid).Set(ctx, q); e != nil {
		t.Fatal(e)
	}
	for _, id := range issues {
		if _, e := s.root().Collection("questions").Doc(qid).Collection("issues").Doc(id).Set(ctx, Issue{ID: id, Slug: id, Title: id, Status: StatusPublished}); e != nil {
			t.Fatal(e)
		}
	}
}
func getIssue(t *testing.T, s *Service, q, i string) Issue {
	t.Helper()
	d, e := s.root().Collection("questions").Doc(q).Collection("issues").Doc(i).Get(context.Background())
	if e != nil {
		t.Fatal(e)
	}
	var out Issue
	if e = d.DataTo(&out); e != nil {
		t.Fatal(e)
	}
	return out
}

func TestAnswerQuestionUniquenessPersonalWeightAndClearing(t *testing.T) {
	s := emulatorService(t)
	seedQuestion(t, s, "ireland", "country", "housing", "cost")
	seedQuestion(t, s, "france", "country", "health")
	caller := Caller{UID: "u1", PhoneVerified: true}
	r1, e := s.Answer(context.Background(), caller, AnswerRequest{QuestionID: "ireland", IssueID: "housing", OperationID: "op1"})
	if e != nil {
		t.Fatal(e)
	}
	if r1.Answer.Revision != 1 || r1.PersonalTop {
		t.Fatalf("unexpected first response %+v", r1)
	}
	r2, e := s.Answer(context.Background(), caller, AnswerRequest{AnswerKind: AnswerKindPersonal, QuestionID: "ireland", IssueID: "housing", OperationID: "op2"})
	if e != nil {
		t.Fatal(e)
	}
	if !r2.PersonalTop {
		t.Fatal("personal top not set")
	}
	i := getIssue(t, s, "ireland", "housing")
	if i.Supporters != 1 || i.PersonalTopSupporters != 1 || i.WeightedScore != 10 {
		t.Fatalf("weighted issue %+v", i)
	}
	r3, e := s.Answer(context.Background(), caller, AnswerRequest{QuestionID: "france", IssueID: "health", OperationID: "op3"})
	if e != nil {
		t.Fatal(e)
	}
	if r3.Answer.Revision != 1 || r3.PersonalTop {
		t.Fatalf("unexpected independent answer %+v", r3)
	}
	old := getIssue(t, s, "ireland", "housing")
	next := getIssue(t, s, "france", "health")
	if old.Supporters != 1 || old.PersonalTopSupporters != 1 || old.WeightedScore != 10 || next.Supporters != 1 || next.WeightedScore != 1 {
		t.Fatalf("old=%+v next=%+v", old, next)
	}
	r4, e := s.Answer(context.Background(), caller, AnswerRequest{QuestionID: "ireland", IssueID: "cost", OperationID: "op4"})
	if e != nil {
		t.Fatal(e)
	}
	if r4.PersonalTop || r4.Answer.Revision != 2 {
		t.Fatalf("replacement=%+v", r4)
	}
	old = getIssue(t, s, "ireland", "housing")
	replacement := getIssue(t, s, "ireland", "cost")
	if old.Supporters != 0 || old.PersonalTopSupporters != 0 || old.WeightedScore != 0 || replacement.Supporters != 1 || replacement.WeightedScore != 1 {
		t.Fatalf("old=%+v replacement=%+v", old, replacement)
	}
}

func TestPersonalEnsuresCategoryChoiceAndIdempotency(t *testing.T) {
	s := emulatorService(t)
	seedQuestion(t, s, "q1", "country", "a")
	seedQuestion(t, s, "q2", "country", "b")
	c := Caller{UID: "u", PhoneVerified: true}
	if _, e := s.Answer(context.Background(), c, AnswerRequest{QuestionID: "q1", IssueID: "a", OperationID: "a1"}); e != nil {
		t.Fatal(e)
	}
	req := AnswerRequest{AnswerKind: AnswerKindPersonal, QuestionID: "q2", IssueID: "b", OperationID: "a2"}
	r, e := s.Answer(context.Background(), c, req)
	if e != nil {
		t.Fatal(e)
	}
	replay, e := s.Answer(context.Background(), c, req)
	if e != nil || !replay.Replayed || replay.Answer.IssueID != "b" {
		t.Fatalf("replay=%+v err=%v", replay, e)
	}
	if r.Changed || !r.PersonalTop {
		t.Fatalf("response=%+v", r)
	}
	i := getIssue(t, s, "q2", "b")
	if i.Supporters != 1 || i.PersonalTopSupporters != 1 || i.WeightedScore != 10 {
		t.Fatalf("issue=%+v", i)
	}
	if first := getIssue(t, s, "q1", "a"); first.Supporters != 1 || first.WeightedScore != 1 {
		t.Fatalf("first question answer was displaced: %+v", first)
	}
}

func TestNestedGeographyUsesIndependentPrioritySlots(t *testing.T) {
	s := emulatorService(t)
	seedQuestionSlot(t, s, "ireland", "geography", "country", "housing")
	seedQuestionSlot(t, s, "dublin", "geography", "city", "transport")
	c := Caller{UID: "u", PhoneVerified: true}
	if _, e := s.Answer(context.Background(), c, AnswerRequest{QuestionID: "ireland", IssueID: "housing", OperationID: "g1"}); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Answer(context.Background(), c, AnswerRequest{QuestionID: "dublin", IssueID: "transport", OperationID: "g2"}); e != nil {
		t.Fatal(e)
	}
	if getIssue(t, s, "ireland", "housing").Supporters != 1 || getIssue(t, s, "dublin", "transport").Supporters != 1 {
		t.Fatal("nested navigation incorrectly shared an answer slot")
	}
}

func TestPaidEligibilityAndReceiptBinding(t *testing.T) {
	s := emulatorService(t)
	seedQuestion(t, s, "q", "country", "a")
	c := Caller{UID: "paid"}
	if _, e := s.Answer(context.Background(), c, AnswerRequest{QuestionID: "q", IssueID: "a", OperationID: "before"}); e != ErrVerificationRequired {
		t.Fatalf("got %v", e)
	}
	if e := s.MarkPaid(context.Background(), "paid", "charge-1"); e != nil {
		t.Fatal(e)
	}
	if e := s.MarkPaid(context.Background(), "paid", "charge-1"); e != nil {
		t.Fatal(e)
	}
	if e := s.MarkPaid(context.Background(), "other", "charge-1"); e != ErrConflict {
		t.Fatalf("got %v", e)
	}
	if _, e := s.Answer(context.Background(), c, AnswerRequest{QuestionID: "q", IssueID: "a", OperationID: "after"}); e != nil {
		t.Fatal(e)
	}
}

func TestFreeformDuplicateUsesAlias(t *testing.T) {
	s := emulatorService(t)
	seedQuestion(t, s, "q", "country")
	c := Caller{UID: "u", PhoneVerified: true}
	r1, e := s.Answer(context.Background(), c, AnswerRequest{QuestionID: "q", Title: "Cost of living", OperationID: "f1"})
	if e != nil {
		t.Fatal(e)
	}
	r2, e := s.Answer(context.Background(), c, AnswerRequest{QuestionID: "q", Title: "  COST  of living ", OperationID: "f2"})
	if e != nil {
		t.Fatal(e)
	}
	if !r1.CreatedIssue || r2.CreatedIssue || r1.Issue.ID != r2.Issue.ID {
		t.Fatalf("first=%+v second=%+v", r1, r2)
	}
}

func TestCreateQuestionPendingSlugIdempotencyAndCreatorPreview(t *testing.T) {
	s := emulatorService(t)
	c := Caller{UID: "author", PhoneVerified: true}
	req := CreateQuestionRequest{Title: "What country is the number one issue for the free world?", Description: "Choose the country whose situation should receive attention first.", AnswerTargetType: "country", OperationID: "create-q-1"}
	created, e := s.CreateQuestion(context.Background(), c, req)
	if e != nil {
		t.Fatal(e)
	}
	if created.Question.Status != StatusPending || created.Question.Indexable || created.Question.CreatorUID != "author" {
		t.Fatalf("created=%+v", created)
	}
	replay, e := s.CreateQuestion(context.Background(), c, req)
	if e != nil || !replay.Replayed || replay.Question.ID != created.Question.ID {
		t.Fatalf("replay=%+v err=%v", replay, e)
	}
	if _, e = s.QuestionBySlug(context.Background(), created.Question.Slug, ""); e != ErrNotFound {
		t.Fatalf("anonymous preview got %v", e)
	}
	if _, e = s.QuestionBySlug(context.Background(), created.Question.Slug, "author"); e != nil {
		t.Fatalf("creator preview: %v", e)
	}
}
func TestCustomQuestionPublicationAndPredefinedPopulation(t *testing.T) {
	s := emulatorService(t)
	c := Caller{UID: "author", PhoneVerified: true}
	custom, e := s.CreateQuestion(context.Background(), c, CreateQuestionRequest{Title: "What should our neighbourhood fix first?", Description: "Choose the single local improvement that should receive attention first.", ChoiceSource: ChoiceSource{Kind: "custom", Options: []ChoiceOption{{Title: "Street lighting"}, {Title: "Playground repairs"}}}, OperationID: "custom-q"})
	if e != nil {
		t.Fatal(e)
	}
	preview, e := s.Question(context.Background(), custom.Question.ID, "author")
	if e != nil || len(preview.Issues) != 2 {
		t.Fatalf("preview issues=%d err=%v", len(preview.Issues), e)
	}
	if e = s.SetQuestionStatus(context.Background(), custom.Question.ID, StatusPublished); e != nil {
		t.Fatal(e)
	}
	published, e := s.Question(context.Background(), custom.Question.ID, "")
	if e != nil || len(published.Issues) != 2 {
		t.Fatalf("published issues=%d err=%v", len(published.Issues), e)
	}
	pre, e := s.CreateQuestion(context.Background(), c, CreateQuestionRequest{Title: "What country needs our attention first?", Description: "Choose one country whose current situation deserves priority attention.", ChoiceSource: ChoiceSource{Kind: "predefined", EntityType: "country"}, OperationID: "country-q"})
	if e != nil {
		t.Fatal(e)
	}
	if e = s.PopulatePredefinedChoices(context.Background(), pre.Question.ID, "country", []PredefinedChoice{{ID: "IE", Title: "Ireland"}, {ID: "FR", Title: "France"}}); e != nil {
		t.Fatal(e)
	}
	if e = s.SetQuestionStatus(context.Background(), pre.Question.ID, StatusPublished); e != nil {
		t.Fatal(e)
	}
	q, e := s.Question(context.Background(), pre.Question.ID, "")
	if e != nil || len(q.Issues) != 2 || q.Issues[0].TargetType != "country" {
		t.Fatalf("country choices=%+v err=%v", q.Issues, e)
	}
}

func TestCountryTargetRejectsFreeformIssue(t *testing.T) {
	s := emulatorService(t)
	seedQuestion(t, s, "q", "geo")
	qref := s.root().Collection("questions").Doc("q")
	if _, e := qref.Update(context.Background(), []firestore.Update{{Path: "answerTargetType", Value: "country"}, {Path: "choiceSource", Value: ChoiceSource{Kind: "predefined", EntityType: "country"}}, {Path: "allowSuggestions", Value: false}}); e != nil {
		t.Fatal(e)
	}
	_, e := s.Answer(context.Background(), Caller{UID: "u", PhoneVerified: true}, AnswerRequest{QuestionID: "q", Title: "Ireland", OperationID: "country-freeform"})
	if !errors.Is(e, ErrValidation) {
		t.Fatalf("got %v", e)
	}
}

func TestLanguageAttributionDoesNotCreateExtraVote(t *testing.T) {
	s := emulatorService(t)
	seedQuestion(t, s, "q", "topic", "a", "b")
	c := Caller{UID: "u", PhoneVerified: true}
	if _, e := s.Answer(context.Background(), c, AnswerRequest{QuestionID: "q", IssueID: "a", LanguageCode: "ru-RU", OperationID: "l1"}); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Answer(context.Background(), c, AnswerRequest{QuestionID: "q", IssueID: "a", LanguageCode: "en", OperationID: "l2"}); e != nil {
		t.Fatal(e)
	}
	a := getIssue(t, s, "q", "a")
	if a.Supporters != 1 || a.LanguageStats["ru"].Supporters != 1 || a.LanguageStats["en"].Supporters != 0 {
		t.Fatalf("display language changed vote: %+v", a)
	}
	if _, e := s.Answer(context.Background(), c, AnswerRequest{AnswerKind: AnswerKindPersonal, QuestionID: "q", IssueID: "a", LanguageCode: "en", OperationID: "l3"}); e != nil {
		t.Fatal(e)
	}
	a = getIssue(t, s, "q", "a")
	if a.LanguageStats["ru"].WeightedScore != 10 {
		t.Fatalf("personal weight not attributed to answer language: %+v", a.LanguageStats)
	}
	if _, e := s.Answer(context.Background(), c, AnswerRequest{QuestionID: "q", IssueID: "b", LanguageCode: "en", OperationID: "l4"}); e != nil {
		t.Fatal(e)
	}
	a = getIssue(t, s, "q", "a")
	b := getIssue(t, s, "q", "b")
	if a.WeightedScore != 0 || a.LanguageStats["ru"].WeightedScore != 0 || b.LanguageStats["en"].Supporters != 1 {
		t.Fatalf("a=%+v b=%+v", a, b)
	}
}

func TestMergeIssueMigratesReferencesCountersAndIsResumable(t *testing.T) {
	s := emulatorService(t)
	seedQuestion(t, s, "q", "topic", "source", "target")
	c1 := Caller{UID: "u1", PhoneVerified: true}
	c2 := Caller{UID: "u2", PhoneVerified: true}
	if _, e := s.Answer(context.Background(), c1, AnswerRequest{AnswerKind: AnswerKindPersonal, QuestionID: "q", IssueID: "source", LanguageCode: "ru", OperationID: "m1"}); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Answer(context.Background(), c2, AnswerRequest{QuestionID: "q", IssueID: "target", LanguageCode: "en", OperationID: "m2"}); e != nil {
		t.Fatal(e)
	}
	if e := s.MergeIssue(context.Background(), "q", "source", "target", "merge-1"); e != nil {
		t.Fatal(e)
	}
	if e := s.MergeIssue(context.Background(), "q", "source", "target", "merge-1"); e != nil {
		t.Fatal(e)
	}
	source := getIssue(t, s, "q", "source")
	target := getIssue(t, s, "q", "target")
	if source.Status != "merged" || source.MergedIntoIssueID != "target" || source.Supporters != 0 || target.Supporters != 2 || target.PersonalTopSupporters != 1 || target.WeightedScore != 11 || target.LanguageStats["ru"].WeightedScore != 10 {
		t.Fatalf("source=%+v target=%+v", source, target)
	}
	a, e := s.OwnAnswer(context.Background(), "u1", "q")
	if e != nil || a.Answer.IssueID != "target" || !a.PersonalTop {
		t.Fatalf("answer=%+v err=%v", a, e)
	}
}
