package publicqa

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

type PredefinedChoice struct{ ID, Title, Description string }

func (s *Service) SetQuestionStatus(ctx context.Context, questionID, statusValue string) error {
	if statusValue != StatusPublished && statusValue != "hidden" && statusValue != "archived" {
		return fmt.Errorf("%w: invalid question status", ErrValidation)
	}
	qref := s.root().Collection("questions").Doc(questionID)
	if statusValue == StatusPublished {
		it := qref.Collection("issues").Where("status", "==", StatusPending).Documents(ctx)
		docs, e := it.GetAll()
		it.Stop()
		if e != nil {
			return e
		}
		for start := 0; start < len(docs); start += 400 {
			end := start + 400
			if end > len(docs) {
				end = len(docs)
			}
			b := s.db.Batch()
			for _, d := range docs[start:end] {
				b.Update(d.Ref, []firestore.Update{{Path: "status", Value: StatusPublished}, {Path: "updatedAt", Value: s.now().UTC()}})
			}
			if _, e = b.Commit(ctx); e != nil {
				return e
			}
		}
	}
	_, e := qref.Update(ctx, []firestore.Update{{Path: "status", Value: statusValue}, {Path: "publication", Value: statusValue}, {Path: "indexable", Value: statusValue == StatusPublished}, {Path: "updatedAt", Value: s.now().UTC()}})
	return e
}
func (s *Service) SetIssueStatus(ctx context.Context, questionID, issueID, statusValue string) error {
	if statusValue != StatusPublished && statusValue != "hidden" && statusValue != "rejected" {
		return fmt.Errorf("%w: invalid issue status", ErrValidation)
	}
	_, e := s.root().Collection("questions").Doc(questionID).Collection("issues").Doc(issueID).Update(ctx, []firestore.Update{{Path: "status", Value: statusValue}, {Path: "updatedAt", Value: s.now().UTC()}})
	return e
}

// PopulatePredefinedChoices is an idempotent trusted administration operation.
// It updates editorial fields but never counters or answer records.
func (s *Service) PopulatePredefinedChoices(ctx context.Context, questionID, entityType string, choices []PredefinedChoice) error {
	if entityType != "country" && entityType != "city" && entityType != "currency" {
		return fmt.Errorf("%w: invalid entity type", ErrValidation)
	}
	if len(choices) == 0 {
		return fmt.Errorf("%w: choices required", ErrValidation)
	}
	qref := s.root().Collection("questions").Doc(questionID)
	for start := 0; start < len(choices); start += 175 {
		end := start + 175
		if end > len(choices) {
			end = len(choices)
		}
		batch := s.db.Batch()
		for _, c := range choices[start:end] {
			if c.ID == "" || c.Title == "" {
				return fmt.Errorf("%w: choice id/title required", ErrValidation)
			}
			iref := qref.Collection("issues").Doc(c.ID)
			fields := map[string]any{"id": c.ID, "slug": Slugify(c.Title), "title": c.Title, "description": c.Description, "targetType": entityType, "targetRef": c.ID, "status": StatusPublished, "updatedAt": s.now().UTC()}
			batch.Set(iref, fields, firestore.MergeAll)
			norm, e := NormalizeTitle(c.Title)
			if e != nil {
				return e
			}
			batch.Set(qref.Collection("aliases").Doc(normalizedHash(norm)), map[string]any{"issueId": c.ID, "normalized": norm}, firestore.MergeAll)
		}
		if _, e := batch.Commit(ctx); e != nil {
			return e
		}
	}
	_, e := qref.Update(ctx, []firestore.Update{{Path: "answerTargetType", Value: entityType}, {Path: "choiceSource", Value: ChoiceSource{Kind: "predefined", EntityType: entityType}}, {Path: "allowSuggestions", Value: false}, {Path: "updatedAt", Value: s.now().UTC()}})
	return e
}

type mergeJob struct {
	QuestionID, SourceIssueID, TargetIssueID, Status string
	CreatedAt, UpdatedAt                             time.Time
}

// MergeIssue is restart-safe: it first locks the question, then migrates every
// private reference in bounded batches, and only then transfers counters and
// marks the retained source record merged. Retrying the same operation resumes.
func (s *Service) MergeIssue(ctx context.Context, questionID, sourceID, targetID, operationID string) error {
	if questionID == "" || sourceID == "" || targetID == "" || sourceID == targetID || operationID == "" {
		return fmt.Errorf("%w: invalid merge", ErrValidation)
	}
	qref := s.root().Collection("questions").Doc(questionID)
	jobref := s.root().Collection("mergeJobs").Doc(operationID)
	lock := "merge:" + operationID
	now := s.now().UTC()
	if e := s.db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		qd, e := tx.Get(qref)
		if e != nil {
			return e
		}
		var q Question
		if e = qd.DataTo(&q); e != nil {
			return e
		}
		if q.MutationLock != "" && q.MutationLock != lock {
			return ErrConflict
		}
		if d, e := tx.Get(jobref); e == nil {
			var j mergeJob
			if e = d.DataTo(&j); e != nil {
				return e
			}
			if j.QuestionID != questionID || j.SourceIssueID != sourceID || j.TargetIssueID != targetID {
				return ErrConflict
			}
			return nil
		} else if !isNotFound(e) {
			return e
		}
		if _, e = tx.Get(qref.Collection("issues").Doc(sourceID)); e != nil {
			return e
		}
		if _, e = tx.Get(qref.Collection("issues").Doc(targetID)); e != nil {
			return e
		}
		if e = tx.Update(qref, []firestore.Update{{Path: "mutationLock", Value: lock}}); e != nil {
			return e
		}
		return tx.Set(jobref, mergeJob{QuestionID: questionID, SourceIssueID: sourceID, TargetIssueID: targetID, Status: "migrating", CreatedAt: now, UpdatedAt: now})
	}); e != nil {
		return e
	}
	for {
		it := qref.Collection("answers").Where("issueId", "==", sourceID).Limit(400).Documents(ctx)
		docs, e := it.GetAll()
		it.Stop()
		if e != nil {
			return e
		}
		if len(docs) == 0 {
			break
		}
		b := s.db.Batch()
		for _, d := range docs {
			b.Update(d.Ref, []firestore.Update{{Path: "issueId", Value: targetID}, {Path: "updatedAt", Value: now}})
		}
		if _, e = b.Commit(ctx); e != nil {
			return e
		}
	}
	for {
		it := s.root().Collection("personalAnswers").Where("questionId", "==", questionID).Where("issueId", "==", sourceID).Limit(400).Documents(ctx)
		docs, e := it.GetAll()
		it.Stop()
		if e != nil {
			return e
		}
		if len(docs) == 0 {
			break
		}
		b := s.db.Batch()
		for _, d := range docs {
			b.Update(d.Ref, []firestore.Update{{Path: "issueId", Value: targetID}, {Path: "updatedAt", Value: now}})
		}
		if _, e = b.Commit(ctx); e != nil {
			return e
		}
	}
	return s.db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		sd, e := tx.Get(qref.Collection("issues").Doc(sourceID))
		if e != nil {
			return e
		}
		td, e := tx.Get(qref.Collection("issues").Doc(targetID))
		if e != nil {
			return e
		}
		jd, e := tx.Get(jobref)
		if e != nil {
			return e
		}
		var source, target Issue
		var job mergeJob
		if e = sd.DataTo(&source); e != nil {
			return e
		}
		if e = td.DataTo(&target); e != nil {
			return e
		}
		if e = jd.DataTo(&job); e != nil {
			return e
		}
		if job.Status == "complete" {
			return nil
		}
		target.Supporters += source.Supporters
		target.PersonalTopSupporters += source.PersonalTopSupporters
		target.WeightedScore += source.WeightedScore
		if target.LanguageStats == nil {
			target.LanguageStats = map[string]LanguageAggregate{}
		}
		for lang, v := range source.LanguageStats {
			x := target.LanguageStats[lang]
			x.Supporters += v.Supporters
			x.PersonalTopSupporters += v.PersonalTopSupporters
			x.WeightedScore += v.WeightedScore
			target.LanguageStats[lang] = x
		}
		source.Supporters = 0
		source.PersonalTopSupporters = 0
		source.WeightedScore = 0
		source.LanguageStats = map[string]LanguageAggregate{}
		source.Status = "merged"
		source.MergedIntoIssueID = targetID
		source.UpdatedAt = now
		target.UpdatedAt = now
		job.Status = "complete"
		job.UpdatedAt = now
		if e = tx.Set(qref.Collection("issues").Doc(sourceID), source); e != nil {
			return e
		}
		if e = tx.Set(qref.Collection("issues").Doc(targetID), target); e != nil {
			return e
		}
		if e = tx.Set(jobref, job); e != nil {
			return e
		}
		return tx.Update(qref, []firestore.Update{{Path: "mutationLock", Value: firestore.Delete}})
	})
}
