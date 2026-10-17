package publicqa

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sneat-co/issue-number-one/backend/translations"
)

func TestService_OptionsAndHelpers(t *testing.T) {
	s := NewService(nil, "", WithPersonalWeight(25))
	if s.personalWeight != 25 {
		t.Fatalf("expected personalWeight 25, got %d", s.personalWeight)
	}

	// invalid weight (< 1) does not overwrite
	WithPersonalWeight(0)(s)
	if s.personalWeight != 25 {
		t.Fatalf("expected personalWeight unchanged at 25, got %d", s.personalWeight)
	}

	// WithClock
	fixedTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s.WithClock(func() time.Time { return fixedTime })
	if s.now() != fixedTime {
		t.Fatalf("expected fixedTime, got %v", s.now())
	}

	// hasIssue
	issues := []Issue{{ID: "iss-1"}, {ID: "iss-2"}}
	if !hasIssue(issues, "iss-1") {
		t.Fatal("expected true for iss-1")
	}
	if hasIssue(issues, "iss-3") {
		t.Fatal("expected false for iss-3")
	}

	// containsString
	list := []string{"en", "es", "fr"}
	if !containsString(list, "es") {
		t.Fatal("expected true for es")
	}
	if containsString(list, "de") {
		t.Fatal("expected false for de")
	}
}

func TestService_ValidationAndEligibility(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// HasPaid with empty uid
	paid, err := s.HasPaid(ctx, "")
	if err != nil || paid {
		t.Fatalf("expected (false, nil) for empty uid, got (%v, %v)", paid, err)
	}

	// requireEligible with unverified caller and not paid
	unverifiedCaller := Caller{UID: "user-unverified", PhoneVerified: false}
	if err := s.requireEligible(ctx, unverifiedCaller); !errors.Is(err, ErrVerificationRequired) {
		t.Fatalf("expected ErrVerificationRequired, got %v", err)
	}

	// MarkPaid validation
	if err := s.MarkPaid(ctx, "", "charge-1"); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for empty uid, got %v", err)
	}
	if err := s.MarkPaid(ctx, "user-1", ""); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for empty chargeID, got %v", err)
	}

	// QuestionBySlug and Question
	seedQuestion(t, s, "q-slug-test", "general", "issue-slug-1")
	// QuestionBySlug not found
	if _, err := s.QuestionBySlug(ctx, "nonexistent-slug", "u1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for nonexistent slug, got %v", err)
	}
	// Question method
	resp, err := s.Question(ctx, "q-slug-test", "u1")
	if err != nil {
		t.Fatalf("unexpected error calling Question: %v", err)
	}
	if resp.Question.ID != "q-slug-test" {
		t.Fatalf("expected q-slug-test, got %s", resp.Question.ID)
	}
}

func TestService_ChoiceSourceValidation(t *testing.T) {
	// Unknown kind
	if err := validateChoiceSource(ChoiceSource{Kind: "unknown-kind"}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for unknown kind, got %v", err)
	}

	// Custom: < 2 options
	if err := validateChoiceSource(ChoiceSource{Kind: "custom", Options: []ChoiceOption{{Title: "Option 1"}}}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for < 2 options, got %v", err)
	}

	// Custom: > 30 options
	thirtyOneOptions := make([]ChoiceOption, 31)
	for i := range thirtyOneOptions {
		thirtyOneOptions[i] = ChoiceOption{Title: fmt.Sprintf("Option %d", i+1)}
	}
	if err := validateChoiceSource(ChoiceSource{Kind: "custom", Options: thirtyOneOptions}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for > 30 options, got %v", err)
	}

	// Custom: description too long (> 500)
	if err := validateChoiceSource(ChoiceSource{Kind: "custom", Options: []ChoiceOption{
		{Title: "Option 1", Description: strings.Repeat("d", 501)},
		{Title: "Option 2"},
	}}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for description > 500 chars, got %v", err)
	}

	// Free: invalid entityType or options
	if err := validateChoiceSource(ChoiceSource{Kind: "free", EntityType: "city"}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for free kind with EntityType, got %v", err)
	}
	if err := validateChoiceSource(ChoiceSource{Kind: "free", Options: []ChoiceOption{{Title: "Option 1"}}}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for free kind with options, got %v", err)
	}
}

func TestService_MarkPaidAndHasPaid_Branches(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// ChargeID with slash
	if err := s.MarkPaid(ctx, "u1", "charge/slash"); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for chargeID with slash, got %v", err)
	}

	// First MarkPaid success
	if err := s.MarkPaid(ctx, "u1", "charge-replay-1"); err != nil {
		t.Fatalf("unexpected error on MarkPaid: %v", err)
	}

	// Replay with same UID -> idempotent nil error
	if err := s.MarkPaid(ctx, "u1", "charge-replay-1"); err != nil {
		t.Fatalf("unexpected error on idempotent MarkPaid: %v", err)
	}

	// Replay with different UID -> ErrConflict
	if err := s.MarkPaid(ctx, "u2", "charge-replay-1"); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict on charge replay with different user, got %v", err)
	}

	// HasPaid with corrupt verification document
	corruptVerRef := s.root().Collection("verification").Doc("u-corrupt-ver")
	if _, err := corruptVerRef.Set(ctx, map[string]any{"eligible": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.HasPaid(ctx, "u-corrupt-ver"); err == nil {
		t.Fatal("expected error on corrupt verification document")
	}

	// requireEligible with corrupt HasPaid error
	if err := s.requireEligible(ctx, Caller{UID: "u-corrupt-ver", PhoneVerified: false}); err == nil {
		t.Fatal("expected error from requireEligible when HasPaid fails")
	}
}

func TestService_CatalogWithSorting(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// Seed multiple categories, concepts, and questions with out-of-order slugs
	catRef1 := s.root().Collection("categories").Doc("cat-b")
	_, _ = catRef1.Set(ctx, Category{ID: "cat-b", Slug: "b-cat", Name: "Cat B", Status: StatusPublished})
	catRef2 := s.root().Collection("categories").Doc("cat-a")
	_, _ = catRef2.Set(ctx, Category{ID: "cat-a", Slug: "a-cat", Name: "Cat A", Status: StatusPublished})

	conRef1 := s.root().Collection("concepts").Doc("con-b")
	_, _ = conRef1.Set(ctx, Concept{ID: "con-b", Slug: "b-con", Title: "Con B", Status: StatusPublished})
	conRef2 := s.root().Collection("concepts").Doc("con-a")
	_, _ = conRef2.Set(ctx, Concept{ID: "con-a", Slug: "a-con", Title: "Con A", Status: StatusPublished})

	qRef1 := s.root().Collection("questions").Doc("q-b")
	_, _ = qRef1.Set(ctx, Question{ID: "q-b", Slug: "b-q", Title: "Q B?", Status: StatusPublished})
	qRef2 := s.root().Collection("questions").Doc("q-a")
	_, _ = qRef2.Set(ctx, Question{ID: "q-a", Slug: "a-q", Title: "Q A?", Status: StatusPublished})

	catalog, err := s.Catalog(ctx)
	if err != nil {
		t.Fatalf("unexpected error fetching catalog: %v", err)
	}
	if len(catalog.Categories) != 2 || len(catalog.Concepts) != 2 || len(catalog.Questions) != 2 {
		t.Fatalf("unexpected catalog count: %+v", catalog)
	}
	if catalog.Categories[0].Slug != "a-cat" || catalog.Concepts[0].Slug != "a-con" || catalog.Questions[0].Slug != "a-q" {
		t.Fatalf("expected elements sorted by slug, got categories=%v, concepts=%v, questions=%v",
			catalog.Categories, catalog.Concepts, catalog.Questions)
	}
}

func TestService_Answer_PersonalTransitionsAndUnchanged(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()
	caller := Caller{UID: "user-personal-transitions", PhoneVerified: true}

	seedQuestion(t, s, "q-trans-1", "cat-trans", "issue-trans-1", "issue-trans-2")
	seedQuestion(t, s, "q-trans-2", "cat-trans", "issue-trans-3")

	// 1. Submit personal answer on question 1
	res1, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-trans-1",
		IssueID:     "issue-trans-1",
		OperationID: "op-trans-1",
		AnswerKind:  AnswerKindPersonal,
	})
	if err != nil {
		t.Fatalf("unexpected error on first personal answer: %v", err)
	}
	if !res1.PersonalTop {
		t.Fatal("expected PersonalTop == true")
	}

	// 2. Submit personal answer on question 2 (switches personal top from q1 to q2 -> lines 573-591)
	res2, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-trans-2",
		IssueID:     "issue-trans-3",
		OperationID: "op-trans-2",
		AnswerKind:  AnswerKindPersonal,
	})
	if err != nil {
		t.Fatalf("unexpected error on second personal answer: %v", err)
	}
	if !res2.PersonalTop {
		t.Fatal("expected PersonalTop == true on q2")
	}

	// 3. Re-answer question 2 with same issue (unchanged: changed == false -> lines 498-503)
	res3, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-trans-2",
		IssueID:     "issue-trans-3",
		OperationID: "op-trans-3",
		AnswerKind:  AnswerKindCategory,
	})
	if err != nil {
		t.Fatalf("unexpected error on unchanged answer: %v", err)
	}
	if res3.Changed {
		t.Fatal("expected Changed == false")
	}

	// 4. Change answer on question 2 to personal answer again, then change category answer on question 2 to different issue
	// Seed issue-trans-4 on question 2
	i4 := Issue{ID: "issue-trans-4", Slug: "issue-trans-4", Title: "Issue 4", Status: StatusPublished}
	_, _ = s.root().Collection("questions").Doc("q-trans-2").Collection("issues").Doc("issue-trans-4").Set(ctx, i4)

	res4, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-trans-2",
		IssueID:     "issue-trans-4",
		OperationID: "op-trans-4",
		AnswerKind:  AnswerKindCategory,
	})
	if err != nil {
		t.Fatalf("unexpected error changing answer on q2: %v", err)
	}
	if !res4.Changed {
		t.Fatal("expected Changed == true")
	}
}

func TestService_CreateQuestion_ValidationAndReplay(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()
	caller := Caller{UID: "creator-1", PhoneVerified: true}

	// Empty caller UID
	if _, err := s.CreateQuestion(ctx, Caller{}, CreateQuestionRequest{
		Title:       "What is the most pressing problem?",
		Description: "A long enough description here.",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for empty caller UID, got %v", err)
	}

	// Empty OperationID
	if _, err := s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:       "What is the most pressing problem?",
		Description: "A long enough description here.",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for empty OperationID, got %v", err)
	}

	// Invalid Title: doesn't end with ?
	if _, err := s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:       "What is the most pressing problem",
		Description: "A long enough description here.",
		OperationID: "op-val-1",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for title missing ?, got %v", err)
	}

	// Invalid Title: too short (< 10)
	if _, err := s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:       "Why??",
		Description: "A long enough description here.",
		OperationID: "op-val-2",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for title < 10 chars, got %v", err)
	}

	// Invalid Title: too long (> 180)
	if _, err := s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:       strings.Repeat("a", 181) + "?",
		Description: "A long enough description here.",
		OperationID: "op-val-3",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for title > 180 chars, got %v", err)
	}

	// Invalid Description: too short (< 20)
	if _, err := s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:       "What is the top issue?",
		Description: "Too short.",
		OperationID: "op-val-4",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for description < 20 chars, got %v", err)
	}

	// Invalid Description: too long (> 1000)
	if _, err := s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:       "What is the top issue?",
		Description: strings.Repeat("d", 1001),
		OperationID: "op-val-5",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for description > 1000 chars, got %v", err)
	}

	// Valid CreateQuestion and Replay
	req := CreateQuestionRequest{
		Title:            "What is the main priority for schools?",
		Description:      "Select or propose the most important issue.",
		AnswerTargetType: "issue",
		OperationID:      "op-replay-1",
		ChoiceSource:     ChoiceSource{Kind: "free"},
		AllowSuggestions: true,
	}
	resp1, err := s.CreateQuestion(ctx, caller, req)
	if err != nil {
		t.Fatalf("unexpected error creating question: %v", err)
	}
	if resp1.Replayed {
		t.Fatal("expected Replayed == false on first creation")
	}

	// Replay same operationID
	resp2, err := s.CreateQuestion(ctx, caller, req)
	if err != nil {
		t.Fatalf("unexpected error on replay: %v", err)
	}
	if !resp2.Replayed || resp2.Question.ID != resp1.Question.ID {
		t.Fatalf("expected replayed with same question ID, got %+v", resp2)
	}
}

func TestService_Answer_ValidationAndBranches(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()
	caller := Caller{UID: "user-ans-val", PhoneVerified: true}

	// 1. Validation: empty questionID, issueID, operationID
	if _, err := s.Answer(ctx, caller, AnswerRequest{IssueID: "i", OperationID: "op"}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for empty questionID, got %v", err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{QuestionID: "q", OperationID: "op"}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for empty issueID, got %v", err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{QuestionID: "q", IssueID: "i"}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for empty operationID, got %v", err)
	}

	// 2. Unverified caller
	if _, err := s.Answer(ctx, Caller{UID: "u-unver"}, AnswerRequest{QuestionID: "q", IssueID: "i", OperationID: "op"}); !errors.Is(err, ErrVerificationRequired) {
		t.Fatalf("expected ErrVerificationRequired, got %v", err)
	}

	// 3. Question not found
	if _, err := s.Answer(ctx, caller, AnswerRequest{QuestionID: "nonexistent", IssueID: "i", OperationID: "op"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for nonexistent question, got %v", err)
	}

	// 4. Question not published (e.g. pending) returns ErrNotFound
	qPendingRef := s.root().Collection("questions").Doc("q-pending")
	if _, err := qPendingRef.Set(ctx, Question{
		ID:     "q-pending",
		Title:  "Pending Question?",
		Status: StatusPending,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{QuestionID: "q-pending", IssueID: "i", OperationID: "op"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for unpublished question, got %v", err)
	}

	// 5. Question with MutationLock returns ErrConflict
	qLockRef := s.root().Collection("questions").Doc("q-ans-lock")
	if _, err := qLockRef.Set(ctx, Question{
		ID:           "q-ans-lock",
		Title:        "Locked Question?",
		Status:       StatusPublished,
		MutationLock: "some-lock",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{QuestionID: "q-ans-lock", IssueID: "i", OperationID: "op"}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict for locked question, got %v", err)
	}

	// 6. Issue not found in question
	seedQuestion(t, s, "q-ans-branch", "general", "issue-ok")
	if _, err := s.Answer(ctx, caller, AnswerRequest{QuestionID: "q-ans-branch", IssueID: "issue-nonexistent", OperationID: "op"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for nonexistent issue in question, got %v", err)
	}
}

func TestService_QuestionLocalized_Branches(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. Invalid language
	if _, err := s.QuestionLocalized(ctx, "q1", "u1", "invalid-lang"); err == nil {
		t.Fatal("expected error on invalid language")
	}

	// 2. Corrupt question document
	qCorruptRef := s.root().Collection("questions").Doc("q-corrupt-data")
	if _, err := qCorruptRef.Set(ctx, map[string]any{"title": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.QuestionLocalized(ctx, "q-corrupt-data", "u1", "en"); err == nil {
		t.Fatal("expected error for corrupt question document")
	}

	// 3. User with own pending answer that gets included in question response
	seedQuestion(t, s, "q-own-pending", "general")
	qRef := s.root().Collection("questions").Doc("q-own-pending")
	// Add pending issue created by user
	pendIssueRef := qRef.Collection("issues").Doc("issue-pending-1")
	if _, err := pendIssueRef.Set(ctx, Issue{
		ID:         "issue-pending-1",
		Title:      "Pending User Suggestion",
		Status:     StatusPending,
		CreatorUID: "user-suggest",
	}); err != nil {
		t.Fatal(err)
	}
	// Add user's answer pointing to that pending issue
	ansRef := qRef.Collection("answers").Doc("user-suggest")
	if _, err := ansRef.Set(ctx, Answer{
		QuestionID: "q-own-pending",
		IssueID:    "issue-pending-1",
	}); err != nil {
		t.Fatal(err)
	}
	// Call QuestionLocalized with uid="user-suggest"
	resp, err := s.QuestionLocalized(ctx, "q-own-pending", "user-suggest", "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	foundPending := false
	for _, iss := range resp.Issues {
		if iss.ID == "issue-pending-1" {
			foundPending = true
			break
		}
	}
	if !foundPending {
		t.Fatal("expected user-suggest's pending issue to be included in Issues")
	}

	// 4. QuestionBySlugLocalized branches
	if _, err := s.QuestionBySlugLocalized(ctx, "", "u1", "en"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for empty slug, got %v", err)
	}
	// Corrupt slug doc
	slugRef := s.root().Collection("questionSlugs").Doc("corrupt-slug")
	if _, err := slugRef.Set(ctx, map[string]any{"questionId": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.QuestionBySlugLocalized(ctx, "corrupt-slug", "u1", "en"); err == nil {
		t.Fatal("expected error for corrupt slug doc")
	}
}

func TestService_OwnAnswer_Branches(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. Corrupt question
	qCorruptRef := s.root().Collection("questions").Doc("q-corrupt-own")
	if _, err := qCorruptRef.Set(ctx, map[string]any{"title": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.OwnAnswer(ctx, "u1", "q-corrupt-own"); err == nil {
		t.Fatal("expected error on corrupt question")
	}

	// 2. Corrupt answer
	seedQuestion(t, s, "q-own-corrupt-ans", "general", "issue-ok")
	qRef := s.root().Collection("questions").Doc("q-own-corrupt-ans")
	ansRef := qRef.Collection("answers").Doc("u1")
	if _, err := ansRef.Set(ctx, map[string]any{"questionId": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.OwnAnswer(ctx, "u1", "q-own-corrupt-ans"); err == nil {
		t.Fatal("expected error on corrupt answer")
	}

	// 3. Corrupt issue
	ansRef2 := qRef.Collection("answers").Doc("u2")
	if _, err := ansRef2.Set(ctx, Answer{
		QuestionID: "q-own-corrupt-ans",
		IssueID:    "issue-corrupt",
	}); err != nil {
		t.Fatal(err)
	}
	corruptIssueRef := qRef.Collection("issues").Doc("issue-corrupt")
	if _, err := corruptIssueRef.Set(ctx, map[string]any{"title": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.OwnAnswer(ctx, "u2", "q-own-corrupt-ans"); err == nil {
		t.Fatal("expected error on corrupt issue")
	}
}

func TestService_Answer_ExtraBranches(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()
	caller := Caller{UID: "user-ans-extra", PhoneVerified: true, DisplayName: "Alex"}

	seedQuestion(t, s, "q-ans-extra", "general", "issue-1")

	// 1. Invalid AnswerKind
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-ans-extra",
		IssueID:     "issue-1",
		OperationID: "op-kind",
		AnswerKind:  "invalid-kind",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for invalid AnswerKind, got %v", err)
	}

	// 2. Invalid Attribution
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-ans-extra",
		IssueID:     "issue-1",
		OperationID: "op-attrib",
		Attribution: "invalid-attrib",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for invalid Attribution, got %v", err)
	}

	// 3. Authored attribution with empty DisplayName
	callerNoName := Caller{UID: "user-no-name", PhoneVerified: true, DisplayName: ""}
	if _, err := s.Answer(ctx, callerNoName, AnswerRequest{
		QuestionID:  "q-ans-extra",
		IssueID:     "issue-1",
		OperationID: "op-authored-empty",
		Attribution: "authored",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for authored with empty DisplayName, got %v", err)
	}

	// 4. OperationID with slash or > 100 chars
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-ans-extra",
		IssueID:     "issue-1",
		OperationID: "op/with/slash",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for operationID with slash, got %v", err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-ans-extra",
		IssueID:     "issue-1",
		OperationID: strings.Repeat("o", 101),
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for operationID > 100 chars, got %v", err)
	}

	// 5. Both issueId and title or neither
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-ans-extra",
		OperationID: "op-none",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation when neither issueId nor title is given, got %v", err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-ans-extra",
		IssueID:     "issue-1",
		Title:       "Also Title",
		OperationID: "op-both",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation when both issueId and title are given, got %v", err)
	}

	// 6. Invalid language in AnswerRequest
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:   "q-ans-extra",
		IssueID:      "issue-1",
		OperationID:  "op-lang",
		LanguageCode: "invalid-lang",
	}); err == nil {
		t.Fatal("expected error on invalid language in answer")
	}

	// 7. Invalid Title (< 3 chars)
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-ans-extra",
		Title:       "ab",
		OperationID: "op-bad-title",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for title < 3 chars, got %v", err)
	}

	// 8. Replay conflict (different UID or payload)
	validAns := AnswerRequest{
		QuestionID:  "q-ans-extra",
		IssueID:     "issue-1",
		OperationID: "op-conflict-test",
	}
	if _, err := s.Answer(ctx, caller, validAns); err != nil {
		t.Fatalf("unexpected error on first answer: %v", err)
	}
	diffCaller := Caller{UID: "different-caller", PhoneVerified: true}
	if _, err := s.Answer(ctx, diffCaller, validAns); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict on replay with different caller UID, got %v", err)
	}

	// 9. Free-form answer when question.AnswerTargetType != "issue"
	qNonIssueRef := s.root().Collection("questions").Doc("q-non-issue")
	if _, err := qNonIssueRef.Set(ctx, Question{
		ID:               "q-non-issue",
		Title:            "What is your favorite country?",
		Status:           StatusPublished,
		AnswerTargetType: "country",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-non-issue",
		Title:       "Valid Title Suggestion",
		OperationID: "op-freeform-non-issue",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for freeform on non-issue question, got %v", err)
	}

	// 10. Suggestions disabled on predefined question
	qNoSugRef := s.root().Collection("questions").Doc("q-no-sug")
	if _, err := qNoSugRef.Set(ctx, Question{
		ID:               "q-no-sug",
		Title:            "Pick a predefined option?",
		Status:           StatusPublished,
		AnswerTargetType: "issue",
		ChoiceSource:     ChoiceSource{Kind: "predefined", EntityType: "city"},
		AllowSuggestions: false,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-no-sug",
		Title:       "My Suggested City",
		OperationID: "op-disallowed-suggestion",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation when suggestions are disallowed, got %v", err)
	}

	// 11. Authored suggestion sets AuthorDisplayName
	authAns, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-ans-extra",
		Title:       "My Authored Suggestion",
		OperationID: "op-authored-suggestion",
		Attribution: "authored",
	})
	if err != nil {
		t.Fatalf("unexpected error on authored answer: %v", err)
	}
	if authAns.Issue.AuthorDisplayName != "Alex" || authAns.Issue.Attribution != "authored" {
		t.Fatalf("expected author Alex and attribution authored, got %+v", authAns.Issue)
	}

	// 12. Candidate issue rate limit (5 suggestions)
	for i := 1; i <= 4; i++ {
		_, err := s.Answer(ctx, caller, AnswerRequest{
			QuestionID:  "q-ans-extra",
			Title:       fmt.Sprintf("Extra Unique Suggestion Number %d", i),
			OperationID: fmt.Sprintf("op-rate-sug-%d", i),
		})
		if err != nil {
			t.Fatalf("unexpected error on suggestion %d: %v", i, err)
		}
	}
	// The 6th candidate suggestion exceeds rate limit (1 initial + 4 loops = 5 count, next is 6)
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-ans-extra",
		Title:       "Extra Unique Suggestion Exceeds Limit",
		OperationID: "op-rate-sug-exceed",
	})
	if !errors.Is(err, ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited on 6th suggestion, got %v", err)
	}

	// 13. Other user trying to answer with someone else's pending issue returns ErrNotFound
	otherCaller := Caller{UID: "other-user-999", PhoneVerified: true}
	if _, err := s.Answer(ctx, otherCaller, AnswerRequest{
		QuestionID:  "q-ans-extra",
		IssueID:     authAns.Issue.ID,
		OperationID: "op-ans-other-pending",
	}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for other user's pending issue, got %v", err)
	}
}

func TestService_TranslateAllQuestionLanguages_WithoutConfig(t *testing.T) {
	s := emulatorService(t)
	ctx := context.Background()
	if _, err := s.TranslateAllQuestionLanguages(ctx, "q1", translations.Actor{}); err == nil {
		t.Fatal("expected error when translation service is nil")
	}

	// translationIsFresh with empty question.SourceLanguage
	q := Question{ID: "q1", Title: "Title", ContentRevision: 1, Publication: "published"}
	tr := translations.Translation{
		SourceRevision: 1,
		SourceHash:     translations.SourceHash(translations.QuestionSource{ID: "q1", Title: "Title", ContentRevision: 1, Publication: "published"}, "en"),
	}
	if !translationIsFresh(q, tr) {
		t.Fatal("expected translation to be fresh with default en sourceLanguage")
	}
}

func TestService_CreateQuestion_MoreBranches(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()
	caller := Caller{UID: "user-cq-more", PhoneVerified: true}

	// 1. Inferred choice source when Kind is empty
	resCountry, err := s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:            "What is your target country?",
		Description:      "Select your favorite country now.",
		AnswerTargetType: "country",
		OperationID:      "op-cq-country-infer",
	})
	if err != nil {
		t.Fatalf("unexpected error on country infer: %v", err)
	}
	if resCountry.Question.ChoiceSource.Kind != "predefined" || resCountry.Question.ChoiceSource.EntityType != "country" {
		t.Fatalf("expected predefined country choiceSource, got %+v", resCountry.Question.ChoiceSource)
	}

	resIssue, err := s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:            "What is your top issue in general?",
		Description:      "Select or write your issue here.",
		AnswerTargetType: "issue",
		OperationID:      "op-cq-issue-infer",
	})
	if err != nil {
		t.Fatalf("unexpected error on issue infer: %v", err)
	}
	if resIssue.Question.ChoiceSource.Kind != "free" {
		t.Fatalf("expected free choiceSource, got %+v", resIssue.Question.ChoiceSource)
	}

	// 2. Slug already exists conflict
	_, err = s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:            "What is your top issue in general?",
		Description:      "Different description for duplicate slug.",
		AnswerTargetType: "issue",
		OperationID:      "op-cq-duplicate-slug",
	})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict for duplicate slug, got %v", err)
	}

	// 3. Rate limit: 3 questions per 24h
	res3, err := s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:            "What is the third question here?",
		Description:      "A valid description for question three.",
		AnswerTargetType: "issue",
		OperationID:      "op-cq-rate-3",
	})
	if err != nil {
		t.Fatalf("unexpected error on 3rd question: %v", err)
	}
	_ = res3

	// 4th question exceeds limit
	_, err = s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:            "What is the fourth question over limit?",
		Description:      "A valid description for question four.",
		AnswerTargetType: "issue",
		OperationID:      "op-cq-rate-4",
	})
	if !errors.Is(err, ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited on 4th question, got %v", err)
	}

	// 4. Operation conflict on replay with different UID
	diffCaller := Caller{UID: "other-user-cq", PhoneVerified: true}
	_, err = s.CreateQuestion(ctx, diffCaller, CreateQuestionRequest{
		Title:            "What is your top issue in general?",
		Description:      "Select or write your issue here.",
		AnswerTargetType: "issue",
		OperationID:      "op-cq-issue-infer",
	})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict on replay with different UID, got %v", err)
	}
}
