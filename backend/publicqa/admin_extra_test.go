package publicqa

import (
	"context"
	"errors"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestAdmin_ValidationAndBranches(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. SetQuestionStatus validation
	if err := s.SetQuestionStatus(ctx, "q1", "invalid-status"); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for invalid question status, got %v", err)
	}

	// SetQuestionStatus non-existent question
	if err := s.SetQuestionStatus(ctx, "nonexistent-q", StatusPublished); err == nil {
		t.Fatal("expected error setting status for nonexistent question")
	}

	// Seed question
	seedQuestion(t, s, "q-admin", "general", "issue-admin-1")

	// SetQuestionStatus "hidden" and "archived"
	if err := s.SetQuestionStatus(ctx, "q-admin", "hidden"); err != nil {
		t.Fatalf("unexpected error setting status to hidden: %v", err)
	}
	if err := s.SetQuestionStatus(ctx, "q-admin", "archived"); err != nil {
		t.Fatalf("unexpected error setting status to archived: %v", err)
	}

	// 2. SetIssueStatus validation
	if err := s.SetIssueStatus(ctx, "q-admin", "issue-admin-1", "invalid-issue-status"); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for invalid issue status, got %v", err)
	}

	// SetIssueStatus for non-existent issue
	if err := s.SetIssueStatus(ctx, "q-admin", "nonexistent-issue", "hidden"); err == nil {
		t.Fatal("expected error setting status for nonexistent issue")
	}

	// SetIssueStatus "hidden" and "rejected"
	if err := s.SetIssueStatus(ctx, "q-admin", "issue-admin-1", "hidden"); err != nil {
		t.Fatalf("unexpected error setting issue to hidden: %v", err)
	}
	if err := s.SetIssueStatus(ctx, "q-admin", "issue-admin-1", "rejected"); err != nil {
		t.Fatalf("unexpected error setting issue to rejected: %v", err)
	}

	// 3. PopulatePredefinedChoices validation
	if err := s.PopulatePredefinedChoices(ctx, "q-admin", "invalid-type", []PredefinedChoice{{ID: "c1", Title: "City 1"}}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for invalid entityType, got %v", err)
	}
	if err := s.PopulatePredefinedChoices(ctx, "q-admin", "city", nil); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for empty choices, got %v", err)
	}
	if err := s.PopulatePredefinedChoices(ctx, "q-admin", "city", []PredefinedChoice{{ID: "", Title: "Title"}}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for empty choice ID, got %v", err)
	}
	if err := s.PopulatePredefinedChoices(ctx, "q-admin", "city", []PredefinedChoice{{ID: "c1", Title: ""}}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for empty choice Title, got %v", err)
	}
	if err := s.PopulatePredefinedChoices(ctx, "q-admin", "city", []PredefinedChoice{{ID: "c1", Title: "ab"}}); err == nil {
		t.Fatal("expected error for invalid choice Title < 3 chars")
	}

	// 4. MergeIssue validation
	for name, args := range map[string][]string{
		"empty questionID":       {"", "src", "tgt", "op"},
		"empty sourceID":         {"q", "", "tgt", "op"},
		"empty targetID":         {"q", "src", "", "op"},
		"same source and target": {"q", "src", "src", "op"},
		"empty operationID":      {"q", "src", "tgt", ""},
	} {
		t.Run(name, func(t *testing.T) {
			if err := s.MergeIssue(ctx, args[0], args[1], args[2], args[3]); !errors.Is(err, ErrValidation) {
				t.Fatalf("expected ErrValidation, got %v", err)
			}
		})
	}

	// MergeIssue question not found
	if err := s.MergeIssue(ctx, "nonexistent-q", "src", "tgt", "op-merge-1"); err == nil {
		t.Fatal("expected error merging on nonexistent question")
	}

	// MergeIssue with existing mutationLock conflict
	qLockRef := s.root().Collection("questions").Doc("q-locked")
	if _, err := qLockRef.Set(ctx, map[string]any{
		"title": "Locked Q", "mutationLock": "other-lock", "contentRevision": int64(1),
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-locked", "src", "tgt", "op-merge-2"); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict on locked question, got %v", err)
	}

	// MergeIssue with conflicting mergeJob (different target)
	jobRef := s.root().Collection("mergeJobs").Doc("op-existing-job")
	if _, err := jobRef.Set(ctx, mergeJob{
		QuestionID:    "q-admin",
		SourceIssueID: "src-1",
		TargetIssueID: "different-target",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-admin", "src-1", "tgt-1", "op-existing-job"); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict on conflicting merge job, got %v", err)
	}

	// MergeIssue source issue not found
	if err := s.MergeIssue(ctx, "q-admin", "nonexistent-src", "issue-admin-1", "op-merge-src-missing"); err == nil {
		t.Fatal("expected error when source issue not found")
	}

	// MergeIssue target issue not found
	if err := s.MergeIssue(ctx, "q-admin", "issue-admin-1", "nonexistent-target", "op-merge-tgt-missing"); err == nil {
		t.Fatal("expected error when target issue not found")
	}

	// 5. Corrupt documents handling in SetQuestionStatus and recomputeQuestionIndexable
	corruptQRef := s.root().Collection("questions").Doc("q-corrupt-status")
	if _, err := corruptQRef.Set(ctx, map[string]any{"title": []int{1, 2}, "status": StatusPublished}); err != nil {
		t.Fatal(err)
	}
	if err := s.SetQuestionStatus(ctx, "q-corrupt-status", StatusPublished); err == nil {
		t.Fatal("expected error from DataTo on corrupt question in SetQuestionStatus")
	}
	if err := s.recomputeQuestionIndexable(ctx, corruptQRef); err == nil {
		t.Fatal("expected error from DataTo on corrupt question in recomputeQuestionIndexable")
	}

	// 6. Merge with target.LanguageStats == nil
	seedQuestion(t, s, "q-merge-lang", "general")
	qMLRef := s.root().Collection("questions").Doc("q-merge-lang")
	srcIssueRef := qMLRef.Collection("issues").Doc("src-lang")
	tgtIssueRef := qMLRef.Collection("issues").Doc("tgt-lang")
	if _, err := srcIssueRef.Set(ctx, Issue{
		ID:     "src-lang",
		Slug:   "src-lang",
		Title:  "Source Language Issue",
		Status: StatusPublished,
		LanguageStats: map[string]LanguageAggregate{
			"en": {Supporters: 5, WeightedScore: 50},
		},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := tgtIssueRef.Set(ctx, map[string]any{
		"id":     "tgt-lang",
		"slug":   "tgt-lang",
		"title":  "Target Issue",
		"status": StatusPublished,
		// No languageStats map -> will be nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-merge-lang", "src-lang", "tgt-lang", "op-merge-nil-lang"); err != nil {
		t.Fatalf("unexpected error merging with nil target languageStats: %v", err)
	}
}

func TestAdmin_SetQuestionStatus_CustomPending(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. Seed custom pending question with pending issues
	qRef := s.root().Collection("questions").Doc("q-custom-pending")
	if _, err := qRef.Set(ctx, Question{
		ID:           "q-custom-pending",
		Title:        "What is the best custom option?",
		Status:       StatusPending,
		ChoiceSource: ChoiceSource{Kind: "custom"},
	}); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"opt-1", "opt-2"} {
		if _, err := qRef.Collection("issues").Doc(id).Set(ctx, Issue{
			ID:     id,
			Title:  "Option " + id,
			Status: StatusPending,
		}); err != nil {
			t.Fatal(err)
		}
	}
	// Publish custom question -> batch updates pending issues to published and recomputes indexable
	if err := s.SetQuestionStatus(ctx, "q-custom-pending", StatusPublished); err != nil {
		t.Fatalf("unexpected error setting custom pending question to published: %v", err)
	}

	// 2. SetQuestionStatus error on update when question doesn't exist and status is hidden
	if err := s.SetQuestionStatus(ctx, "nonexistent-hidden-q", "hidden"); err == nil {
		t.Fatal("expected error setting hidden status on nonexistent question")
	}

	// 3. PopulatePredefinedChoices error on update when question doesn't exist
	if err := s.PopulatePredefinedChoices(ctx, "nonexistent-choices-q", "city", []PredefinedChoice{{ID: "c1", Title: "City One"}}); err == nil {
		t.Fatal("expected error on PopulatePredefinedChoices for nonexistent question")
	}

	// 4. recomputeQuestionIndexable on reserved doc ID
	if err := s.recomputeQuestionIndexable(ctx, s.root().Collection("questions").Doc("__reserved__")); err == nil {
		t.Fatal("expected error on recomputeQuestionIndexable for reserved doc id")
	}
}

func TestAdmin_MergeIssue_ExtraBranches(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. Corrupt question document in tx1
	seedQuestion(t, s, "q-merge-extra", "general", "issue-m1", "issue-m2")
	corruptQ := s.root().Collection("questions").Doc("q-merge-corrupt-data")
	if _, err := corruptQ.Set(ctx, map[string]any{"title": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-merge-corrupt-data", "issue-m1", "issue-m2", "op-corrupt-q"); err == nil {
		t.Fatal("expected error on corrupt question doc in MergeIssue")
	}

	// 2. Non-NotFound error reading jobref in tx1 (reserved doc id)
	if err := s.MergeIssue(ctx, "q-merge-extra", "issue-m1", "issue-m2", "__reserved__"); err == nil {
		t.Fatal("expected error on reserved operation ID in MergeIssue")
	}

	// 3. Corrupt job document in tx1
	corruptJob := s.root().Collection("mergeJobs").Doc("op-corrupt-job")
	if _, err := corruptJob.Set(ctx, map[string]any{"questionId": []int{1, 2}}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-merge-extra", "issue-m1", "issue-m2", "op-corrupt-job"); err == nil {
		t.Fatal("expected error on corrupt job doc in MergeIssue")
	}

	// 4. Completed job resumption (job.Status == "complete")
	completedJob := s.root().Collection("mergeJobs").Doc("op-completed-job")
	if _, err := completedJob.Set(ctx, mergeJob{
		QuestionID:    "q-merge-extra",
		SourceIssueID: "issue-m1",
		TargetIssueID: "issue-m2",
		Status:        "complete",
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-merge-extra", "issue-m1", "issue-m2", "op-completed-job"); err != nil {
		t.Fatalf("unexpected error resuming complete job: %v", err)
	}

	// 5. In tx2: job is already in "migrating" state
	// 5a. Missing source issue in tx2
	jobMissingSrc := s.root().Collection("mergeJobs").Doc("op-migrating-missing-src")
	if _, err := jobMissingSrc.Set(ctx, mergeJob{
		QuestionID:    "q-merge-extra",
		SourceIssueID: "issue-missing-src",
		TargetIssueID: "issue-m2",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-merge-extra", "issue-missing-src", "issue-m2", "op-migrating-missing-src"); err == nil {
		t.Fatal("expected error on missing source issue in MergeIssue tx2")
	}

	// 5b. Missing target issue in tx2
	jobMissingTgt := s.root().Collection("mergeJobs").Doc("op-migrating-missing-tgt")
	if _, err := jobMissingTgt.Set(ctx, mergeJob{
		QuestionID:    "q-merge-extra",
		SourceIssueID: "issue-m1",
		TargetIssueID: "issue-missing-tgt",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-merge-extra", "issue-m1", "issue-missing-tgt", "op-migrating-missing-tgt"); err == nil {
		t.Fatal("expected error on missing target issue in MergeIssue tx2")
	}

	// 5c. Missing job doc in tx2
	// Seed job for tx1, but delete it right before tx2
	jobDeleted := s.root().Collection("mergeJobs").Doc("op-migrating-deleted-job")
	if _, err := jobDeleted.Set(ctx, mergeJob{
		QuestionID:    "q-merge-extra",
		SourceIssueID: "issue-m1",
		TargetIssueID: "issue-m2",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	// If we delete jobdoc right now, tx1 sees NotFound.
	// But what if tx2 runs when job doc is corrupt?
	// We can test corrupt source and corrupt target in tx2:
	// 5d. Corrupt source issue doc in tx2
	jobCorruptSrc := s.root().Collection("mergeJobs").Doc("op-migrating-corrupt-src")
	if _, err := jobCorruptSrc.Set(ctx, mergeJob{
		QuestionID:    "q-merge-extra",
		SourceIssueID: "issue-m1-corrupt",
		TargetIssueID: "issue-m2",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	corruptSrcRef := s.root().Collection("questions").Doc("q-merge-extra").Collection("issues").Doc("issue-m1-corrupt")
	if _, err := corruptSrcRef.Set(ctx, map[string]any{"supporters": "not-a-number"}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-merge-extra", "issue-m1-corrupt", "issue-m2", "op-migrating-corrupt-src"); err == nil {
		t.Fatal("expected error on corrupt source issue in MergeIssue tx2")
	}

	// 5e. Corrupt target issue doc in tx2
	jobCorruptTgt := s.root().Collection("mergeJobs").Doc("op-migrating-corrupt-tgt")
	if _, err := jobCorruptTgt.Set(ctx, mergeJob{
		QuestionID:    "q-merge-extra",
		SourceIssueID: "issue-m1",
		TargetIssueID: "issue-m2-corrupt",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	corruptTgtRef := s.root().Collection("questions").Doc("q-merge-extra").Collection("issues").Doc("issue-m2-corrupt")
	if _, err := corruptTgtRef.Set(ctx, map[string]any{"supporters": "not-a-number"}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-merge-extra", "issue-m1", "issue-m2-corrupt", "op-migrating-corrupt-tgt"); err == nil {
		t.Fatal("expected error on corrupt target issue in MergeIssue tx2")
	}

	// 6. Test migration loops with personal answer and answer on source
	seedQuestion(t, s, "q-merge-loops", "general", "src-loop", "tgt-loop")
	cLoop := Caller{UID: "user-merge-loop", PhoneVerified: true}
	if _, err := s.Answer(ctx, cLoop, AnswerRequest{
		AnswerKind:  AnswerKindPersonal,
		QuestionID:  "q-merge-loops",
		IssueID:     "src-loop",
		OperationID: "op-loop-ans",
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.MergeIssue(ctx, "q-merge-loops", "src-loop", "tgt-loop", "op-merge-loops"); err != nil {
		t.Fatalf("unexpected error in MergeIssue with answers and personalAnswers: %v", err)
	}
}

type atomicInt struct {
	val int64
}

func (a *atomicInt) Add(n int64) int64 {
	a.val += n
	return a.val
}

func (a *atomicInt) Load() int64 {
	return a.val
}

type cancelAfterNContext struct {
	context.Context
	doneChan     chan struct{}
	doneCalls    atomicInt
	triggerAfter int64
}

func newCancelAfterNContext(ctx context.Context, n int64) *cancelAfterNContext {
	return &cancelAfterNContext{
		Context:      ctx,
		doneChan:     make(chan struct{}),
		triggerAfter: n,
	}
}

func (c *cancelAfterNContext) cancel() {
	select {
	case <-c.doneChan:
	default:
		close(c.doneChan)
	}
}

func (c *cancelAfterNContext) Done() <-chan struct{} {
	if c.doneCalls.Add(1) > c.triggerAfter {
		c.cancel()
		return c.doneChan
	}
	return c.Context.Done()
}

func (c *cancelAfterNContext) Err() error {
	select {
	case <-c.doneChan:
		return context.Canceled
	default:
		if c.doneCalls.Load() > c.triggerAfter {
			c.cancel()
			return context.Canceled
		}
		return c.Context.Err()
	}
}

func TestRecomputeQuestionIndexable_QueryError(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	seedQuestion(t, s, "q-cancel-test", "general", "issue-ct-1")
	qref := s.root().Collection("questions").Doc("q-cancel-test")
	// Trigger cancellation after Get completes (after 2 Done calls)
	ctx := newCancelAfterNContext(context.Background(), 2)
	err := s.recomputeQuestionIndexable(ctx, qref)
	if err == nil {
		t.Fatal("expected error on canceled query in recomputeQuestionIndexable")
	}
	t.Logf("got expected error: %v", err)
}

func TestAdmin_CanceledContextBranches(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	// 1. PopulatePredefinedChoices with canceled context -> line 125
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()
	err := s.PopulatePredefinedChoices(canceledCtx, "q-predef-cancel", "city", []PredefinedChoice{{ID: "c1", Title: "City One"}})
	if err == nil {
		t.Fatal("expected error on canceled context in PopulatePredefinedChoices")
	}

	// 2. SetQuestionStatus on custom pending question with canceled context during query (line 60)
	qRef := s.root().Collection("questions").Doc("q-custom-cancel-query")
	if _, err := qRef.Set(ctx, Question{
		ID:           "q-custom-cancel-query",
		Title:        "Custom question for cancel query?",
		Status:       StatusPending,
		ChoiceSource: ChoiceSource{Kind: "custom"},
	}); err != nil {
		t.Fatal(err)
	}
	ctxAfterGet := newCancelAfterNContext(ctx, 2)
	if err := s.SetQuestionStatus(ctxAfterGet, "q-custom-cancel-query", StatusPublished); err == nil {
		t.Fatal("expected error on canceled query in SetQuestionStatus")
	}

	// 3. SetQuestionStatus on custom pending question with canceled context during batch commit (line 72)
	qRef2 := s.root().Collection("questions").Doc("q-custom-cancel-commit")
	if _, err := qRef2.Set(ctx, Question{
		ID:           "q-custom-cancel-commit",
		Title:        "Custom question for cancel commit?",
		Status:       StatusPending,
		ChoiceSource: ChoiceSource{Kind: "custom"},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := qRef2.Collection("issues").Doc("opt-cancel-commit").Set(ctx, Issue{
		ID:     "opt-cancel-commit",
		Title:  "Option Cancel Commit",
		Status: StatusPending,
	}); err != nil {
		t.Fatal(err)
	}
	ctxAfterQuery := newCancelAfterNContext(ctx, 4)
	err = s.SetQuestionStatus(ctxAfterQuery, "q-custom-cancel-commit", StatusPublished)
	if err == nil {
		t.Fatal("expected error on canceled commit in SetQuestionStatus")
	}
	t.Logf("step 3 error: %v, done calls: %d", err, ctxAfterQuery.doneCalls.Load())
}

func TestAdmin_MergeIssue_CanceledLoops(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	seedQuestion(t, s, "q-merge-cancel", "general", "src-c", "tgt-c")
	qRef := s.root().Collection("questions").Doc("q-merge-cancel")
	cUser := Caller{UID: "u-merge-cancel", PhoneVerified: true}

	// 1. Answers GetAll error (line 193)
	jobRef := s.root().Collection("mergeJobs").Doc("op-mc-answers-get")
	if _, err := jobRef.Set(ctx, mergeJob{
		QuestionID:    "q-merge-cancel",
		SourceIssueID: "src-c",
		TargetIssueID: "tgt-c",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	ctx1 := newLineCancelCtx(ctx, "MergeIssue", 187, 192)
	if err := s.MergeIssue(ctx1, "q-merge-cancel", "src-c", "tgt-c", "op-mc-answers-get"); err == nil {
		t.Fatal("expected error on canceled answers query in MergeIssue")
	}

	// 2. Answers batch commit error (line 200)
	ansRef := qRef.Collection("answers").Doc(cUser.UID)
	if _, err := ansRef.Set(ctx, Answer{
		CategoryID:     "general",
		PrioritySlotID: "slot-1",
		QuestionID:     "q-merge-cancel",
		IssueID:        "src-c",
	}); err != nil {
		t.Fatal(err)
	}
	jobRef2 := s.root().Collection("mergeJobs").Doc("op-mc-answers-commit")
	if _, err := jobRef2.Set(ctx, mergeJob{
		QuestionID:    "q-merge-cancel",
		SourceIssueID: "src-c",
		TargetIssueID: "tgt-c",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	ctx2 := newLineCancelCtx(ctx, "MergeIssue", 196, 202)
	if err := s.MergeIssue(ctx2, "q-merge-cancel", "src-c", "tgt-c", "op-mc-answers-commit"); err == nil {
		t.Fatal("expected error on canceled answers commit in MergeIssue")
	}

	// 3. PersonalAnswers GetAll error (line 206)
	_, _ = ansRef.Delete(ctx)
	jobRef3 := s.root().Collection("mergeJobs").Doc("op-mc-personal-get")
	if _, err := jobRef3.Set(ctx, mergeJob{
		QuestionID:    "q-merge-cancel",
		SourceIssueID: "src-c",
		TargetIssueID: "tgt-c",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	ctx3 := newLineCancelCtx(ctx, "MergeIssue", 205, 207)
	if err := s.MergeIssue(ctx3, "q-merge-cancel", "src-c", "tgt-c", "op-mc-personal-get"); err == nil {
		t.Fatal("expected error on canceled personalAnswers query in MergeIssue")
	}

	// 4. PersonalAnswers batch commit error (line 218)
	pref := s.root().Collection("personalAnswers").Doc(cUser.UID)
	if _, err := pref.Set(ctx, PersonalAnswer{
		CategoryID:     "general",
		PrioritySlotID: "slot-1",
		QuestionID:     "q-merge-cancel",
		IssueID:        "src-c",
	}); err != nil {
		t.Fatal(err)
	}
	jobRef4 := s.root().Collection("mergeJobs").Doc("op-mc-personal-commit")
	if _, err := jobRef4.Set(ctx, mergeJob{
		QuestionID:    "q-merge-cancel",
		SourceIssueID: "src-c",
		TargetIssueID: "tgt-c",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	ctx4 := newLineCancelCtx(ctx, "MergeIssue", 214, 220)
	if err := s.MergeIssue(ctx4, "q-merge-cancel", "src-c", "tgt-c", "op-mc-personal-commit"); err == nil {
		t.Fatal("expected error on canceled personalAnswers commit in MergeIssue")
	}
}

type lineCancelCtx struct {
	context.Context
	matchFunc string
	minLine   int
	maxLine   int
	canceled  bool
	doneChan  chan struct{}
}

func newLineCancelCtx(ctx context.Context, matchFunc string, minLine, maxLine int) *lineCancelCtx {
	return &lineCancelCtx{
		Context:   ctx,
		matchFunc: matchFunc,
		minLine:   minLine,
		maxLine:   maxLine,
		doneChan:  make(chan struct{}),
	}
}

func (c *lineCancelCtx) check() {
	if c.canceled {
		return
	}
	var pcs [32]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.Function, c.matchFunc) && frame.Line >= c.minLine && frame.Line <= c.maxLine {
			c.canceled = true
			close(c.doneChan)
			return
		}
		if !more {
			break
		}
	}
}

func (c *lineCancelCtx) Done() <-chan struct{} {
	c.check()
	if c.canceled {
		return c.doneChan
	}
	return c.Context.Done()
}

func (c *lineCancelCtx) Err() error {
	c.check()
	if c.canceled {
		return context.Canceled
	}
	return c.Context.Err()
}

func (c *lineCancelCtx) Value(key any) any {
	c.check()
	return c.Context.Value(key)
}

func (c *lineCancelCtx) Deadline() (time.Time, bool) {
	c.check()
	return c.Context.Deadline()
}

type triggerOnCallerCtx struct {
	context.Context
	targetSubstr string
	action       func()
	fired        bool
}

func (c *triggerOnCallerCtx) check() {
	if c.fired || c.action == nil {
		return
	}
	var pcs [32]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.Function, c.targetSubstr) {
			c.fired = true
			c.action()
			return
		}
		if !more {
			break
		}
	}
}

func (c *triggerOnCallerCtx) Done() <-chan struct{} {
	c.check()
	return c.Context.Done()
}

func (c *triggerOnCallerCtx) Value(key any) any {
	c.check()
	return c.Context.Value(key)
}

func (c *triggerOnCallerCtx) Deadline() (time.Time, bool) {
	c.check()
	return c.Context.Deadline()
}

func TestAdmin_MergeIssue_JobMissingOrCorruptInTx2(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}
	s := emulatorService(t)
	ctx := context.Background()

	seedQuestion(t, s, "q-merge-tx2-err", "general", "src-tx2", "tgt-tx2")

	// 1. Missing job in tx2 (line 231-234)
	jobDelRef := s.root().Collection("mergeJobs").Doc("op-tx2-del-job")
	if _, err := jobDelRef.Set(ctx, mergeJob{
		QuestionID:    "q-merge-tx2-err",
		SourceIssueID: "src-tx2",
		TargetIssueID: "tgt-tx2",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	ctxDel := &triggerOnCallerCtx{
		Context:      ctx,
		targetSubstr: "MergeIssue.func2",
		action: func() {
			_, _ = jobDelRef.Delete(context.Background())
		},
	}
	if err := s.MergeIssue(ctxDel, "q-merge-tx2-err", "src-tx2", "tgt-tx2", "op-tx2-del-job"); err == nil {
		t.Fatal("expected error when job is missing in tx2")
	}

	// 2. Corrupt job in tx2 (line 243-245)
	jobCorruptRef := s.root().Collection("mergeJobs").Doc("op-tx2-corrupt-job")
	if _, err := jobCorruptRef.Set(ctx, mergeJob{
		QuestionID:    "q-merge-tx2-err",
		SourceIssueID: "src-tx2",
		TargetIssueID: "tgt-tx2",
		Status:        "migrating",
	}); err != nil {
		t.Fatal(err)
	}
	ctxCorrupt := &triggerOnCallerCtx{
		Context:      ctx,
		targetSubstr: "MergeIssue.func2",
		action: func() {
			_, _ = jobCorruptRef.Set(context.Background(), map[string]any{"status": 12345})
		},
	}
	if err := s.MergeIssue(ctxCorrupt, "q-merge-tx2-err", "src-tx2", "tgt-tx2", "op-tx2-corrupt-job"); err == nil {
		t.Fatal("expected error on corrupt job in MergeIssue tx2")
	}
}
