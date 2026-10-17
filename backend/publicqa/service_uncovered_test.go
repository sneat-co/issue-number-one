package publicqa

import (
	"context"
	"errors"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/sneat-co/issue-number-one/backend/translations"
)

func TestService_MarkPaid_Extra(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. Corrupt existing receipt
	r := s.root().Collection("paymentReceipts").Doc("corrupt-receipt")
	if _, err := r.Set(ctx, map[string]any{"uid": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if err := s.MarkPaid(ctx, "u1", "corrupt-receipt"); err == nil {
		t.Fatal("expected error on corrupt payment receipt")
	}

	// 2. Reserved chargeID (non-NotFound error)
	if err := s.MarkPaid(ctx, "u1", "__reserved__"); err == nil {
		t.Fatal("expected error on reserved chargeID")
	}
}

func TestService_Catalog_Extra(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. Corrupt category
	if _, err := s.root().Collection("categories").Doc("corrupt-cat").Set(ctx, map[string]any{
		"status": StatusPublished, "slug": []int{1, 2},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Catalog(ctx); err == nil {
		t.Fatal("expected error on corrupt category in Catalog")
	}
	_, _ = s.root().Collection("categories").Doc("corrupt-cat").Delete(ctx)

	// 2. Corrupt concept
	if _, err := s.root().Collection("concepts").Doc("corrupt-concept").Set(ctx, map[string]any{
		"status": StatusPublished, "slug": []int{1, 2},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Catalog(ctx); err == nil {
		t.Fatal("expected error on corrupt concept in Catalog")
	}
	_, _ = s.root().Collection("concepts").Doc("corrupt-concept").Delete(ctx)

	// 3. Corrupt question
	if _, err := s.root().Collection("questions").Doc("corrupt-q-cat").Set(ctx, map[string]any{
		"status": StatusPublished, "slug": []int{1, 2},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Catalog(ctx); err == nil {
		t.Fatal("expected error on corrupt question in Catalog")
	}
	_, _ = s.root().Collection("questions").Doc("corrupt-q-cat").Delete(ctx)
}

func TestService_Question_Extra(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. QuestionLocalized with nonexistent qid
	if _, err := s.QuestionLocalized(ctx, "nonexistent-qid", "", "en"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// 2. QuestionLocalized with reserved qid
	if _, err := s.QuestionLocalized(ctx, "__reserved__", "", "en"); err == nil {
		t.Fatal("expected error on reserved qid")
	}

	// 3. QuestionLocalized with corrupt issue
	seedQuestion(t, s, "q-corrupt-issue", "general")
	qRef := s.root().Collection("questions").Doc("q-corrupt-issue")
	if _, err := qRef.Collection("issues").Doc("issue-bad").Set(ctx, map[string]any{
		"status": StatusPublished, "supporters": "not-a-number",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.QuestionLocalized(ctx, "q-corrupt-issue", "", "en"); err == nil {
		t.Fatal("expected error on corrupt issue")
	}
	_, _ = qRef.Collection("issues").Doc("issue-bad").Delete(ctx)

	// 4. QuestionLocalized sort with different WeightedScores
	seedQuestion(t, s, "q-sort-scores", "general", "issue-s1", "issue-s2")
	qSortRef := s.root().Collection("questions").Doc("q-sort-scores")
	_, _ = qSortRef.Collection("issues").Doc("issue-s1").Update(ctx, []firestore.Update{
		{Path: "weightedScore", Value: int64(10)},
	})
	_, _ = qSortRef.Collection("issues").Doc("issue-s2").Update(ctx, []firestore.Update{
		{Path: "weightedScore", Value: int64(20)},
	})
	res, err := s.QuestionLocalized(ctx, "q-sort-scores", "", "en")
	if err != nil {
		t.Fatalf("unexpected error on QuestionLocalized: %v", err)
	}
	if len(res.Issues) >= 2 && res.Issues[0].WeightedScore < res.Issues[1].WeightedScore {
		t.Fatalf("expected descending sort, got %+v", res.Issues)
	}
}

func TestService_OwnAnswer_Extra(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	seedQuestion(t, s, "q-own-extra", "general", "issue-1")

	// 1. Reserved UID in OwnAnswer -> ad.Get returns InvalidArgument
	if _, err := s.OwnAnswer(ctx, "__reserved__", "q-own-extra"); err == nil {
		t.Fatal("expected error on reserved UID in OwnAnswer")
	}

	// 2. Answer document pointing to reserved issueID -> id.Get returns error
	ansRef := s.root().Collection("questions").Doc("q-own-extra").Collection("answers").Doc("user-bad-issue")
	if _, err := ansRef.Set(ctx, Answer{
		QuestionID: "q-own-extra",
		IssueID:    "__reserved__",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.OwnAnswer(ctx, "user-bad-issue", "q-own-extra"); err == nil {
		t.Fatal("expected error when answer points to reserved issueID")
	}
}

func TestService_Answer_TargetAndDocErrors(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()
	caller := Caller{UID: "user-target-errs", PhoneVerified: true}

	seedQuestion(t, s, "q-tgt-errs", "general", "issue-t1", "issue-t2")

	// 1. OperationID == "__reserved__" -> tx.Get(opref) returns non-NotFound
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		IssueID:     "issue-t1",
		OperationID: "__reserved__",
	}); err == nil {
		t.Fatal("expected error on reserved operationID in Answer")
	}

	// 2. Free-form answer on question where AnswerTargetType is not "issue"
	qCountryRef := s.root().Collection("questions").Doc("q-country-target")
	if _, err := qCountryRef.Set(ctx, Question{
		ID:               "q-country-target",
		Title:            "What is your country?",
		Status:           StatusPublished,
		AnswerTargetType: "country",
		AllowSuggestions: true,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-country-target",
		Title:       "Free form custom text",
		OperationID: "op-freeform-country",
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for free-form on country target, got %v", err)
	}

	// 3. Caller UID == "__reserved__" -> tx.Get(aref) returns non-NotFound
	badCaller := Caller{UID: "__reserved__", PhoneVerified: true}
	if _, err := s.Answer(ctx, badCaller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		IssueID:     "issue-t1",
		OperationID: "op-bad-uid",
	}); err == nil {
		t.Fatal("expected error on reserved UID in Answer")
	}

	// 4. In load(old.QuestionID, old.IssueID):
	// 4a. old.QuestionID is "__reserved__"
	ansBadQ := s.root().Collection("questions").Doc("q-tgt-errs").Collection("answers").Doc(caller.UID)
	if _, err := ansBadQ.Set(ctx, Answer{
		QuestionID: "__reserved__",
		IssueID:    "issue-t1",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		IssueID:     "issue-t2",
		OperationID: "op-load-bad-q",
	}); err == nil {
		t.Fatal("expected error when old answer questionID is reserved")
	}

	// 4b. old.QuestionID points to corrupt question
	corruptQ := s.root().Collection("questions").Doc("q-corrupt-load")
	if _, err := corruptQ.Set(ctx, map[string]any{"title": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := ansBadQ.Set(ctx, Answer{
		QuestionID: "q-corrupt-load",
		IssueID:    "issue-t1",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		IssueID:     "issue-t2",
		OperationID: "op-load-corrupt-q",
	}); err == nil {
		t.Fatal("expected error when old answer question is corrupt")
	}

	// 4c. old.IssueID is "__reserved__"
	if _, err := ansBadQ.Set(ctx, Answer{
		QuestionID: "q-tgt-errs",
		IssueID:    "__reserved__",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		IssueID:     "issue-t2",
		OperationID: "op-load-bad-issue",
	}); err == nil {
		t.Fatal("expected error when old answer issueID is reserved")
	}

	// 4d. old.IssueID points to corrupt issue
	corruptI := s.root().Collection("questions").Doc("q-tgt-errs").Collection("issues").Doc("issue-corrupt-load")
	if _, err := corruptI.Set(ctx, map[string]any{"title": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := ansBadQ.Set(ctx, Answer{
		QuestionID: "q-tgt-errs",
		IssueID:    "issue-corrupt-load",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		IssueID:     "issue-t2",
		OperationID: "op-load-corrupt-i",
	}); err == nil {
		t.Fatal("expected error when old answer issue is corrupt")
	}

	// 4e. personal.QuestionID is "__reserved__"
	prefBad := s.root().Collection("personalAnswers").Doc(caller.UID)
	if _, err := prefBad.Set(ctx, PersonalAnswer{
		QuestionID: "__reserved__",
		IssueID:    "issue-t1",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		IssueID:     "issue-t2",
		OperationID: "op-load-bad-personal",
	}); err == nil {
		t.Fatal("expected error when personal answer questionID is reserved")
	}
	_, _ = prefBad.Delete(ctx)

	// Clean up ansBadQ
	_, _ = ansBadQ.Delete(ctx)

	// 5. In loadTarget:
	// 5a. Corrupt question doc in loadTarget
	corruptTargetQ := s.root().Collection("questions").Doc("q-corrupt-target")
	if _, err := corruptTargetQ.Set(ctx, map[string]any{"title": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-corrupt-target",
		IssueID:     "issue-1",
		OperationID: "op-corrupt-target-q",
	}); err == nil {
		t.Fatal("expected error on corrupt target question in Answer")
	}

	// 5b. Corrupt alias doc in loadTarget
	qRef := s.root().Collection("questions").Doc("q-tgt-errs")
	norm, _ := NormalizeTitle("Existing Option")
	aliasRef := qRef.Collection("aliases").Doc(normalizedHash(norm))
	if _, err := aliasRef.Set(ctx, map[string]any{"issueId": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		Title:       "Existing Option",
		OperationID: "op-corrupt-alias",
	}); err == nil {
		t.Fatal("expected error on corrupt alias in Answer")
	}
	_, _ = aliasRef.Delete(ctx)

	// 5c. Corrupt candidate limit doc in loadTarget
	limitRef := s.root().Collection("actors").Doc(caller.UID).Collection("limits").Doc("candidate-issues")
	if _, err := limitRef.Set(ctx, map[string]any{"count": "not-int"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		Title:       "Brand New Suggestion",
		OperationID: "op-corrupt-limit",
	}); err == nil {
		t.Fatal("expected error on corrupt limit doc in Answer")
	}
	_, _ = limitRef.Delete(ctx)

	// 5c2. Reserved caller UID on new suggestion -> limitRef returns non-NotFound (line 887)
	resUser := Caller{UID: "__reserved__", PhoneVerified: true}
	if _, err := s.Answer(ctx, resUser, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		Title:       "Brand New Suggestion For Reserved UID",
		OperationID: "op-res-limit",
	}); err == nil {
		t.Fatal("expected error on reserved UID candidate limit in Answer")
	}

	// 5c3. Canceled context during alias lookup (line 903)
	ctxAfterTargetQ := newCancelAfterNContext(ctx, 2)
	if _, err := s.Answer(ctxAfterTargetQ, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		Title:       "Canceled Alias Lookup Suggestion",
		OperationID: "op-cancel-alias",
	}); err == nil {
		t.Fatal("expected error on canceled alias lookup in Answer")
	}

	// 5d. Corrupt issue doc in loadTarget
	corruptIssue := qRef.Collection("issues").Doc("issue-target-corrupt")
	if _, err := corruptIssue.Set(ctx, map[string]any{"supporters": "not-int"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		IssueID:     "issue-target-corrupt",
		OperationID: "op-corrupt-target-issue",
	}); err == nil {
		t.Fatal("expected error on corrupt target issue in Answer")
	}

	// 6. Canceled context on pref.Get in Answer (line 448)
	// loadTarget does: Get(qref) [2] + Get(iref) [2] = 4 calls. Then Get(aref) [2] = 6 calls.
	// Cancel after 6 so tx.Get(pref) fails with Canceled
	ctxAfterAref := newCancelAfterNContext(ctx, 6)
	if _, err := s.Answer(ctxAfterAref, caller, AnswerRequest{
		QuestionID:  "q-tgt-errs",
		IssueID:     "issue-t1",
		OperationID: "op-cancel-pref",
	}); err == nil {
		t.Fatal("expected error on canceled pref in Answer")
	}
}

func TestService_Answer_Invariants_Missing(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()
	caller := Caller{UID: "user-inv-missing", PhoneVerified: true}

	seedQuestion(t, s, "q-inv-miss", "general", "issue-im1", "issue-im2")
	qRef := s.root().Collection("questions").Doc("q-inv-miss")
	i1Ref := qRef.Collection("issues").Doc("issue-im1")
	now := s.now().UTC()

	// 1. oq.RespondentsByLanguage == nil (line 527)
	ansRef := qRef.Collection("answers").Doc(caller.UID)
	if _, err := ansRef.Set(ctx, Answer{
		CategoryID:     "general",
		PrioritySlotID: "slot-1",
		QuestionID:     "q-inv-miss",
		IssueID:        "issue-im1",
		LanguageCode:   "en",
		Revision:       1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := qRef.Set(ctx, Question{
		ID:                    "q-inv-miss",
		Title:                 "Invariant question test?",
		Status:                StatusPublished,
		TotalRespondents:      1,
		RespondentsByLanguage: nil, // nil triggers line 527 initialization, then line 529 invariant error
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := i1Ref.Set(ctx, Issue{
		ID:            "issue-im1",
		Title:         "Issue IM1",
		Status:        StatusPublished,
		Supporters:    1,
		WeightedScore: 1,
		LanguageStats: map[string]LanguageAggregate{
			"en": {Supporters: 1, WeightedScore: 1},
		},
	}); err != nil {
		t.Fatal(err)
	}
	// Change answer -> hits line 527, then line 529 ("question language counter invariant")
	_, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-inv-miss",
		IssueID:     "issue-im2",
		OperationID: "op-nil-lang-map",
	})
	if err == nil || err.Error() != "question language counter invariant" {
		t.Fatalf("expected question language counter invariant error, got %v", err)
	}

	// 2. Personal counter invariant in category answer (line 554)
	// Set valid respondents map
	_, _ = qRef.Update(ctx, []firestore.Update{
		{Path: "respondentsByLanguage", Value: map[string]int64{"en": 1}},
	})
	pref := s.root().Collection("personalAnswers").Doc(caller.UID)
	if _, err := pref.Set(ctx, PersonalAnswer{
		CategoryID:     "general",
		PrioritySlotID: "slot-1",
		QuestionID:     "q-inv-miss",
		IssueID:        "issue-im1",
		LanguageCode:   "en",
		UpdatedAt:      now,
	}); err != nil {
		t.Fatal(err)
	}
	// Issue has personalTopSupporters = 0 -> triggers line 554
	_, _ = i1Ref.Update(ctx, []firestore.Update{
		{Path: "personalTopSupporters", Value: int64(0)},
		{Path: "weightedScore", Value: int64(1)},
	})
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-inv-miss",
		IssueID:     "issue-im2",
		OperationID: "op-pi-zero",
		AnswerKind:  AnswerKindCategory,
	})
	if err == nil || err.Error() != "personal counter invariant" {
		t.Fatalf("expected personal counter invariant error, got %v", err)
	}

	// 3. Personal language counter invariant in category answer (line 564)
	// Fix personalTopSupporters = 1, but languageStats["en"].personalTopSupporters = 0
	_, _ = i1Ref.Update(ctx, []firestore.Update{
		{Path: "personalTopSupporters", Value: int64(1)},
		{Path: "weightedScore", Value: int64(10)},
		{Path: "languageStats.en.supporters", Value: int64(1)},
		{Path: "languageStats.en.weightedScore", Value: int64(10)},
		{Path: "languageStats.en.personalTopSupporters", Value: int64(0)},
	})
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-inv-miss",
		IssueID:     "issue-im2",
		OperationID: "op-pla-zero",
		AnswerKind:  AnswerKindCategory,
	})
	if err == nil || err.Error() != "personal language counter invariant" {
		t.Fatalf("expected personal language counter invariant error, got %v", err)
	}

	// 4. Personal counter invariant in AnswerKindPersonal switch (line 576)
	// Fix issue-im1 back to personalTopSupporters = 0 for personal switch test
	_, _ = i1Ref.Update(ctx, []firestore.Update{
		{Path: "personalTopSupporters", Value: int64(0)},
	})
	// User has personal top on issue-im1, now switches to issue-im2 with AnswerKindPersonal
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-inv-miss",
		IssueID:     "issue-im2",
		OperationID: "op-pi-switch-zero",
		AnswerKind:  AnswerKindPersonal,
	})
	if err == nil || err.Error() != "personal counter invariant" {
		t.Fatalf("expected personal counter invariant error on personal switch, got %v", err)
	}

	// 5. Personal language counter invariant in AnswerKindPersonal switch (line 586)
	_, _ = i1Ref.Update(ctx, []firestore.Update{
		{Path: "personalTopSupporters", Value: int64(1)},
		{Path: "weightedScore", Value: int64(10)},
		{Path: "languageStats.en.supporters", Value: int64(1)},
		{Path: "languageStats.en.weightedScore", Value: int64(10)},
		{Path: "languageStats.en.personalTopSupporters", Value: int64(0)},
	})
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-inv-miss",
		IssueID:     "issue-im2",
		OperationID: "op-pla-switch-zero",
		AnswerKind:  AnswerKindPersonal,
	})
	if err == nil || err.Error() != "personal language counter invariant" {
		t.Fatalf("expected personal language counter invariant error on personal switch, got %v", err)
	}
}

type fakeTranslator struct{}

func (fakeTranslator) Translate(ctx context.Context, sourceLanguage, targetLanguage string, texts []string) ([]string, error) {
	return texts, nil
}

func TestService_CreateQuestion_Extra(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. Translation service doesn't support sourceLanguage (lines 670-671)
	repo, err := translations.NewFirestoreRepository(s.db, "in1-test")
	if err != nil {
		t.Fatal(err)
	}
	transSvc, err := translations.NewService(repo, fakeTranslator{}, translations.Config{
		SupportedLanguages: []string{"en"},
	})
	if err != nil {
		t.Fatal(err)
	}
	sWithTrans := NewService(s.db, "in1-test", WithTranslations(transSvc, repo))
	caller := Caller{UID: "user-cq-extra", PhoneVerified: true}
	_, err = sWithTrans.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:          "What is the best language to learn?",
		Description:    "Valid description with more than 20 characters.",
		SourceLanguage: "ru", // not in transSvc.SupportedLanguages()
		OperationID:    "op-unsupported-lang",
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for unsupported language, got %v", err)
	}

	// 2. Caller not eligible (line 684)
	unverifiedCaller := Caller{UID: "unverified-user", PhoneVerified: false}
	_, err = s.CreateQuestion(ctx, unverifiedCaller, CreateQuestionRequest{
		Title:          "What is the question from unverified?",
		Description:    "Valid description with more than 20 characters.",
		SourceLanguage: "en",
		ChoiceSource:   ChoiceSource{Kind: "free"},
		OperationID:    "op-unverified",
	})
	if !errors.Is(err, ErrVerificationRequired) {
		t.Fatalf("expected ErrVerificationRequired, got %v", err)
	}

	// 3. Invalid titles / descriptions (lines 689-690)
	for name, req := range map[string]CreateQuestionRequest{
		"no question mark": {
			Title:          "What is the question without mark",
			Description:    "Valid description with more than 20 characters.",
			ChoiceSource:   ChoiceSource{Kind: "free"},
			SourceLanguage: "en",
			OperationID:    "op-no-mark",
		},
		"description too short": {
			Title:          "What is the question with mark?",
			Description:    "Too short.",
			ChoiceSource:   ChoiceSource{Kind: "free"},
			SourceLanguage: "en",
			OperationID:    "op-short-desc",
		},
		"title has url": {
			Title:          "What is http://example.com question?",
			Description:    "Valid description with more than 20 characters.",
			ChoiceSource:   ChoiceSource{Kind: "free"},
			SourceLanguage: "en",
			OperationID:    "op-url-title",
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := s.CreateQuestion(ctx, caller, req)
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("expected ErrValidation, got %v", err)
			}
		})
	}

	// 4. Corrupt questionOperation doc (line 701)
	opRef := s.root().Collection("questionOperations").Doc("op-corrupt-cq")
	if _, err := opRef.Set(ctx, map[string]any{"payloadHash": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	_, err = s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:          "What is the corrupt operation question?",
		Description:    "Valid description with more than 20 characters.",
		ChoiceSource:   ChoiceSource{Kind: "free"},
		SourceLanguage: "en",
		OperationID:    "op-corrupt-cq",
	})
	if err == nil {
		t.Fatal("expected error on corrupt questionOperation doc")
	}

	// 5. Reserved operationID (line 710)
	_, err = s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:          "What is the reserved operation question?",
		Description:    "Valid description with more than 20 characters.",
		ChoiceSource:   ChoiceSource{Kind: "free"},
		SourceLanguage: "en",
		OperationID:    "__reserved__",
	})
	if err == nil {
		t.Fatal("expected error on reserved operationID in CreateQuestion")
	}

	// 6. Corrupt limit doc (line 722)
	limitRef := s.root().Collection("actors").Doc(caller.UID).Collection("limits").Doc("questions")
	if _, err := limitRef.Set(ctx, map[string]any{"count": "not-int"}); err != nil {
		t.Fatal(err)
	}
	_, err = s.CreateQuestion(ctx, caller, CreateQuestionRequest{
		Title:          "What is the corrupt limit question?",
		Description:    "Valid description with more than 20 characters.",
		ChoiceSource:   ChoiceSource{Kind: "free"},
		SourceLanguage: "en",
		OperationID:    "op-cq-corrupt-limit",
	})
	if err == nil {
		t.Fatal("expected error on corrupt limit doc in CreateQuestion")
	}
	_, _ = limitRef.Delete(ctx)

	// 7. slugref non-NotFound error (line 716) via context cancel after opref.Get
	ctxAfterOp := newCancelAfterNContext(ctx, 2)
	_, err = s.CreateQuestion(ctxAfterOp, caller, CreateQuestionRequest{
		Title:          "What is the canceled slug question?",
		Description:    "Valid description with more than 20 characters.",
		ChoiceSource:   ChoiceSource{Kind: "free"},
		SourceLanguage: "en",
		OperationID:    "op-cq-cancel-slug",
	})
	if err == nil {
		t.Fatal("expected error on canceled slugref in CreateQuestion")
	}

	// 8. limitref non-NotFound error (line 725) via reserved UID
	reservedCaller := Caller{UID: "__reserved__", PhoneVerified: true}
	_, err = s.CreateQuestion(ctx, reservedCaller, CreateQuestionRequest{
		Title:          "What is the reserved caller question?",
		Description:    "Valid description with more than 20 characters.",
		ChoiceSource:   ChoiceSource{Kind: "free"},
		SourceLanguage: "en",
		OperationID:    "op-cq-res-uid",
	})
	if err == nil {
		t.Fatal("expected error on reserved UID in CreateQuestion")
	}

	// 9. validateChoiceSource with title causing NormalizeTitle error (line 832)
	err = validateChoiceSource(ChoiceSource{
		Kind: "custom",
		Options: []ChoiceOption{
			{Title: "ab"}, // less than 3 runes
			{Title: "Valid Option"},
		},
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for invalid custom option title, got %v", err)
	}
}

func TestService_TranslateAllQuestionLanguages_Reserved(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. Unconfigured translation service
	if _, err := s.TranslateAllQuestionLanguages(ctx, "q1", translations.Actor{UID: "u1"}); err == nil {
		t.Fatal("expected error when translation service is unconfigured")
	}

	// 2. Configured translation service with question that fails
	repo, err := translations.NewFirestoreRepository(s.db, "in1-test")
	if err != nil {
		t.Fatal(err)
	}
	transSvc, err := translations.NewService(repo, fakeTranslator{}, translations.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sWithTrans := NewService(s.db, "in1-test", WithTranslations(transSvc, repo))
	// Calling with reserved question causes doc.Set to fail at line 796
	langs, err := sWithTrans.TranslateAllQuestionLanguages(ctx, "__reserved__", translations.Actor{UID: "u1", Trusted: true})
	if err == nil {
		t.Fatal("expected error translating reserved question")
	}
	_ = langs
}

func TestService_QuestionLocalized_IteratorError(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	seedQuestion(t, s, "q-iter-err", "general", "issue-it-1")
	// Cancel after Get (2 Done calls) so it.Next() fails at line 258
	ctxAfterGet := newCancelAfterNContext(ctx, 2)
	_, err := s.QuestionLocalized(ctxAfterGet, "q-iter-err", "", "en")
	if err == nil {
		t.Fatal("expected error from iterator in QuestionLocalized")
	}
}

func TestService_RemainingUncoveredLines(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	seedQuestion(t, s, "q-rem-1", "general", "issue-rem-1")

	// 1. Answer: tx.Get(pref) non-NotFound error -> line 446
	cUser := Caller{UID: "u-rem-pref-err", PhoneVerified: true}
	ctxPref := newLineCancelCtx(ctx, "Answer", 437, 442)
	_, err := s.Answer(ctxPref, cUser, AnswerRequest{
		QuestionID:  "q-rem-1",
		IssueID:     "issue-rem-1",
		OperationID: "op-rem-pref-err",
	})
	if err == nil {
		t.Fatal("expected error on canceled pref in Answer")
	}

	// 2. Answer: load(personal.QuestionID, personal.IssueID) error -> line 491
	cUserPersonalBad := Caller{UID: "u-rem-personal-bad", PhoneVerified: true}
	pref := s.root().Collection("personalAnswers").Doc(cUserPersonalBad.UID)
	if _, err := pref.Set(ctx, PersonalAnswer{
		CategoryID:     "general",
		PrioritySlotID: "slot-1",
		QuestionID:     "nonexistent-question-for-load",
		IssueID:        "nonexistent-issue",
	}); err != nil {
		t.Fatal(err)
	}
	_, err = s.Answer(ctx, cUserPersonalBad, AnswerRequest{
		QuestionID:  "q-rem-1",
		IssueID:     "issue-rem-1",
		OperationID: "op-rem-bad-personal",
	})
	if err == nil {
		t.Fatal("expected error when personal answer references nonexistent question")
	}

	// 3. CreateQuestion: tx.Get(slugref) non-NotFound error -> line 700
	cCreator := Caller{UID: "u-rem-cq-err", PhoneVerified: true}
	ctxSlug := newLineCancelCtx(ctx, "CreateQuestion", 695, 699)
	_, err = s.CreateQuestion(ctxSlug, cCreator, CreateQuestionRequest{
		Title:            "Unique Question Title For Slug Cancel?",
		Description:      "This is a valid question description with more than 20 characters.",
		AnswerTargetType: "free",
		ChoiceSource:     ChoiceSource{Kind: "free"},
		OperationID:      "op-rem-cq-slug-err",
	})
	if err == nil {
		t.Fatal("expected error on canceled slug in CreateQuestion")
	}

	// 4. loadTarget: tx.Get(st.aliasRef) non-NotFound error -> line 877
	cTarget := Caller{UID: "u-rem-lt-err", PhoneVerified: true}
	ctxAlias := newLineCancelCtx(ctx, "loadTarget", 843, 847)
	_, err = s.Answer(ctxAlias, cTarget, AnswerRequest{
		QuestionID:  "q-rem-1",
		Title:       "A Suggested Issue Title For Alias Error",
		OperationID: "op-rem-alias-err",
	})
	if err == nil {
		t.Fatal("expected error on canceled alias in loadTarget")
	}
}
