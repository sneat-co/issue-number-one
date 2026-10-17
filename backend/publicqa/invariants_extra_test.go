package publicqa

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
)

func TestService_CounterInvariantsAndLanguageDefaults(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()
	caller := Caller{UID: "user-invariants", PhoneVerified: true}

	seedQuestion(t, s, "q-inv-1", "cat-inv", "issue-inv-1", "issue-inv-2")
	seedQuestion(t, s, "q-inv-2", "cat-inv", "issue-inv-3")

	// 1. Initial answer with empty LanguageCode
	ansRef := s.root().Collection("questions").Doc("q-inv-1").Collection("answers").Doc(caller.UID)
	now := s.now().UTC()
	_, _ = ansRef.Set(ctx, Answer{
		CategoryID:     "cat-inv",
		PrioritySlotID: "slot-1",
		QuestionID:     "q-inv-1",
		IssueID:        "issue-inv-1",
		LanguageCode:   "", // empty -> triggers line 501 / 516
		Revision:       1,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	// Set question total respondents and issue stats
	qRef := s.root().Collection("questions").Doc("q-inv-1")
	_, _ = qRef.Update(ctx, []firestore.Update{
		{Path: "totalRespondents", Value: int64(1)},
		{Path: "respondentsByLanguage.en", Value: int64(1)},
	})
	iRef := qRef.Collection("issues").Doc("issue-inv-1")
	_, _ = iRef.Update(ctx, []firestore.Update{
		{Path: "supporters", Value: int64(1)},
		{Path: "weightedScore", Value: int64(1)},
		{Path: "languageStats.en.supporters", Value: int64(1)},
		{Path: "languageStats.en.weightedScore", Value: int64(1)},
	})

	// Re-answer unchanged with empty language in old answer -> exercises lines 500-502
	_, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-inv-1",
		IssueID:     "issue-inv-1",
		OperationID: "op-inv-unchanged",
	})
	if err != nil {
		t.Fatalf("unexpected error on unchanged answer with empty old lang: %v", err)
	}

	// Change answer to issue-inv-2 with empty old.LanguageCode -> exercises lines 514-517
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-inv-1",
		IssueID:     "issue-inv-2",
		OperationID: "op-inv-change-empty-lang",
	})
	if err != nil {
		t.Fatalf("unexpected error on changed answer with empty old lang: %v", err)
	}

	// 2. Trigger counter invariant (line 510) by zeroing supporters on old issue
	_, _ = qRef.Collection("issues").Doc("issue-inv-2").Update(ctx, []firestore.Update{
		{Path: "supporters", Value: int64(0)},
	})
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-inv-1",
		IssueID:     "issue-inv-1",
		OperationID: "op-inv-counter-fail",
	})
	if err == nil || err.Error() != "counter invariant" {
		t.Fatalf("expected 'counter invariant' error, got %v", err)
	}

	// Fix counter and trigger language counter invariant (line 520)
	_, _ = qRef.Collection("issues").Doc("issue-inv-2").Update(ctx, []firestore.Update{
		{Path: "supporters", Value: int64(1)},
		{Path: "languageStats.en.supporters", Value: int64(0)},
	})
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-inv-1",
		IssueID:     "issue-inv-1",
		OperationID: "op-inv-lang-counter-fail",
	})
	if err == nil || err.Error() != "language counter invariant" {
		t.Fatalf("expected 'language counter invariant' error, got %v", err)
	}

	// Fix language counter and trigger question language counter invariant (line 530)
	_, _ = qRef.Collection("issues").Doc("issue-inv-2").Update(ctx, []firestore.Update{
		{Path: "languageStats.en.supporters", Value: int64(1)},
	})
	_, _ = qRef.Update(ctx, []firestore.Update{
		{Path: "respondentsByLanguage.en", Value: int64(0)},
	})
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-inv-1",
		IssueID:     "issue-inv-1",
		OperationID: "op-inv-q-lang-counter-fail",
	})
	if err == nil || err.Error() != "question language counter invariant" {
		t.Fatalf("expected 'question language counter invariant' error, got %v", err)
	}
}

func TestService_PersonalInvariantsAndLanguageDefaults(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()
	caller := Caller{UID: "user-pers-inv", PhoneVerified: true}

	seedQuestion(t, s, "q-pinv-1", "cat-pinv", "issue-pinv-1", "issue-pinv-2")
	seedQuestion(t, s, "q-pinv-2", "cat-pinv", "issue-pinv-3")

	// Set initial personal answer on q-pinv-1 with empty LanguageCode
	now := s.now().UTC()
	pref := s.root().Collection("personalAnswers").Doc(caller.UID)
	_, _ = pref.Set(ctx, PersonalAnswer{
		CategoryID:     "cat-pinv",
		PrioritySlotID: "slot-1",
		QuestionID:     "q-pinv-1",
		IssueID:        "issue-pinv-1",
		LanguageCode:   "", // empty -> triggers line 560
		UpdatedAt:      now,
	})
	aref := s.root().Collection("questions").Doc("q-pinv-1").Collection("answers").Doc(caller.UID)
	_, _ = aref.Set(ctx, Answer{
		CategoryID:     "cat-pinv",
		PrioritySlotID: "slot-1",
		QuestionID:     "q-pinv-1",
		IssueID:        "issue-pinv-1",
		LanguageCode:   "en",
		Revision:       1,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	qRef := s.root().Collection("questions").Doc("q-pinv-1")
	_, _ = qRef.Update(ctx, []firestore.Update{
		{Path: "totalRespondents", Value: int64(1)},
		{Path: "respondentsByLanguage.en", Value: int64(1)},
	})
	iRef := qRef.Collection("issues").Doc("issue-pinv-1")
	_, _ = iRef.Update(ctx, []firestore.Update{
		{Path: "supporters", Value: int64(1)},
		{Path: "personalTopSupporters", Value: int64(1)},
		{Path: "weightedScore", Value: int64(10)},
		{Path: "languageStats.en.supporters", Value: int64(1)},
		{Path: "languageStats.en.personalTopSupporters", Value: int64(1)},
		{Path: "languageStats.en.weightedScore", Value: int64(10)},
	})

	// Change category answer on q-pinv-1 from issue-1 to issue-2 with empty personal.LanguageCode
	// Exercises lines 551-571 with pl == "" (line 560)
	_, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-pinv-1",
		IssueID:     "issue-pinv-2",
		OperationID: "op-pinv-change-empty-pl",
		AnswerKind:  AnswerKindCategory,
	})
	if err != nil {
		t.Fatalf("unexpected error on clearing personal top with empty pl: %v", err)
	}

	// Now set up personal answer on q-pinv-1 again with empty LanguageCode
	_, _ = pref.Set(ctx, PersonalAnswer{
		CategoryID:     "cat-pinv",
		PrioritySlotID: "slot-1",
		QuestionID:     "q-pinv-1",
		IssueID:        "issue-pinv-2",
		LanguageCode:   "", // empty -> triggers line 582
		UpdatedAt:      now,
	})
	i2Ref := qRef.Collection("issues").Doc("issue-pinv-2")
	_, _ = i2Ref.Update(ctx, []firestore.Update{
		{Path: "personalTopSupporters", Value: int64(1)},
		{Path: "weightedScore", Value: int64(10)},
		{Path: "languageStats.en.personalTopSupporters", Value: int64(1)},
		{Path: "languageStats.en.weightedScore", Value: int64(10)},
	})

	// Switch personal answer to q-pinv-2 with empty personal.LanguageCode -> exercises lines 573-591 with pl == "" (line 582)
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-pinv-2",
		IssueID:     "issue-pinv-3",
		OperationID: "op-pinv-switch-empty-pl",
		AnswerKind:  AnswerKindPersonal,
	})
	if err != nil {
		t.Fatalf("unexpected error on switching personal top with empty pl: %v", err)
	}
}

func TestService_CorruptDocumentsInAnswer(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()
	caller := Caller{UID: "user-corrupt-ans", PhoneVerified: true}

	seedQuestion(t, s, "q-corrupt-flow", "cat-corrupt", "issue-cf-1")

	// 1. Corrupt operation document
	opRef := s.root().Collection("operations").Doc("op-corrupt-doc")
	_, _ = opRef.Set(ctx, map[string]any{"payloadHash": []int{1, 2}})
	_, err := s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-corrupt-flow",
		IssueID:     "issue-cf-1",
		OperationID: "op-corrupt-doc",
	})
	if err == nil {
		t.Fatal("expected error on corrupt operation doc")
	}

	// 2. Corrupt old answer document
	ansRef := s.root().Collection("questions").Doc("q-corrupt-flow").Collection("answers").Doc(caller.UID)
	_, _ = ansRef.Set(ctx, map[string]any{"questionId": []int{1, 2}})
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-corrupt-flow",
		IssueID:     "issue-cf-1",
		OperationID: "op-corrupt-old-ans",
	})
	if err == nil {
		t.Fatal("expected error on corrupt old answer doc")
	}

	// 3. Corrupt personal answer document
	_, _ = ansRef.Delete(ctx)
	pref := s.root().Collection("personalAnswers").Doc(caller.UID)
	_, _ = pref.Set(ctx, map[string]any{"questionId": []int{1, 2}})
	_, err = s.Answer(ctx, caller, AnswerRequest{
		QuestionID:  "q-corrupt-flow",
		IssueID:     "issue-cf-1",
		OperationID: "op-corrupt-personal-doc",
	})
	if err == nil {
		t.Fatal("expected error on corrupt personal answer doc")
	}
}
