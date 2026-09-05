package publicqa

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isNotFound(err error) bool { return status.Code(err) == codes.NotFound }

var (
	ErrNotFound             = errors.New("not found")
	ErrConflict             = errors.New("operation conflicts with its original actor or payload")
	ErrValidation           = errors.New("validation failed")
	ErrRateLimited          = errors.New("candidate issue submission rate limit exceeded")
	ErrVerificationRequired = errors.New("phone or payment verification required")
)

type Option func(*Service)

func WithPersonalWeight(weight int64) Option {
	return func(s *Service) {
		if weight >= 1 {
			s.personalWeight = weight
		}
	}
}

type Service struct {
	db             *firestore.Client
	spaceID        string
	now            func() time.Time
	personalWeight int64
}

func NewService(db *firestore.Client, spaceID string, opts ...Option) *Service {
	if spaceID == "" {
		spaceID = DefaultPublicSpaceID
	}
	s := &Service{db: db, spaceID: spaceID, now: time.Now, personalWeight: 10}
	for _, o := range opts {
		o(s)
	}
	return s
}
func (s *Service) WithClock(now func() time.Time) *Service { s.now = now; return s }
func (s *Service) root() *firestore.DocumentRef {
	return s.db.Collection("spaces").Doc(s.spaceID).Collection("ext").Doc(ExtensionID)
}

type paidVerification struct {
	Eligible   bool      `firestore:"eligible"`
	ChargeID   string    `firestore:"chargeId"`
	VerifiedAt time.Time `firestore:"verifiedAt"`
}

func (s *Service) HasPaid(ctx context.Context, uid string) (bool, error) {
	if uid == "" {
		return false, nil
	}
	d, e := s.root().Collection("verification").Doc(uid).Get(ctx)
	if isNotFound(e) {
		return false, nil
	}
	if e != nil {
		return false, e
	}
	var v paidVerification
	if e = d.DataTo(&v); e != nil {
		return false, e
	}
	return v.Eligible, nil
}
func (s *Service) requireEligible(ctx context.Context, caller Caller) error {
	if caller.PhoneVerified {
		return nil
	}
	paid, e := s.HasPaid(ctx, caller.UID)
	if e != nil {
		return e
	}
	if !paid {
		return ErrVerificationRequired
	}
	return nil
}

// MarkPaid is for a trusted settled-payment consumer; it is deliberately not
// mounted as an HTTP route. The caller must prove an exact EUR 1.00 settlement.
func (s *Service) MarkPaid(ctx context.Context, uid, chargeID string) error {
	if uid == "" || chargeID == "" || strings.Contains(chargeID, "/") {
		return fmt.Errorf("%w: uid and safe chargeID required", ErrValidation)
	}
	return s.db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		r := s.root().Collection("paymentReceipts").Doc(chargeID)
		if d, e := tx.Get(r); e == nil {
			var x struct {
				UID string `firestore:"uid"`
			}
			if e = d.DataTo(&x); e != nil {
				return e
			}
			if x.UID != uid {
				return ErrConflict
			}
			return nil
		} else if !isNotFound(e) {
			return e
		}
		now := s.now().UTC()
		if e := tx.Set(r, map[string]any{"uid": uid, "createdAt": now}); e != nil {
			return e
		}
		return tx.Set(s.root().Collection("verification").Doc(uid), paidVerification{Eligible: true, ChargeID: chargeID, VerifiedAt: now})
	})
}

func readPublished[T any](ctx context.Context, c *firestore.CollectionRef, dst *[]T) error {
	it := c.Where("status", "==", StatusPublished).Documents(ctx)
	defer it.Stop()
	for {
		d, e := it.Next()
		if e == iterator.Done {
			return nil
		}
		if e != nil {
			return e
		}
		var v T
		if e = d.DataTo(&v); e != nil {
			return e
		}
		*dst = append(*dst, v)
	}
}
func (s *Service) Catalog(ctx context.Context) (CatalogResponse, error) {
	var o CatalogResponse
	if e := readPublished(ctx, s.root().Collection("categories"), &o.Categories); e != nil {
		return o, e
	}
	if e := readPublished(ctx, s.root().Collection("concepts"), &o.Concepts); e != nil {
		return o, e
	}
	if e := readPublished(ctx, s.root().Collection("questions"), &o.Questions); e != nil {
		return o, e
	}
	sort.Slice(o.Categories, func(i, j int) bool { return o.Categories[i].Slug < o.Categories[j].Slug })
	sort.Slice(o.Concepts, func(i, j int) bool { return o.Concepts[i].Slug < o.Concepts[j].Slug })
	sort.Slice(o.Questions, func(i, j int) bool { return o.Questions[i].Slug < o.Questions[j].Slug })
	return o, nil
}
func (s *Service) Question(ctx context.Context, qid, uid string) (QuestionResponse, error) {
	var o QuestionResponse
	qref := s.root().Collection("questions").Doc(qid)
	d, e := qref.Get(ctx)
	if e != nil {
		if isNotFound(e) {
			return o, ErrNotFound
		}
		return o, e
	}
	if e = d.DataTo(&o.Question); e != nil {
		return o, e
	}
	if o.Question.Status != StatusPublished && !(o.Question.Status == StatusPending && o.Question.CreatorUID == uid) {
		return o, ErrNotFound
	}
	o.TotalRespondents = o.Question.TotalRespondents
	o.UpdatedAt = o.Question.UpdatedAt
	it := qref.Collection("issues").Where("status", "==", StatusPublished).Documents(ctx)
	defer it.Stop()
	for {
		d, e := it.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			return o, e
		}
		var i Issue
		if e = d.DataTo(&i); e != nil {
			return o, e
		}
		o.Issues = append(o.Issues, i)
	}
	if uid != "" {
		a, e := s.OwnAnswer(ctx, uid, qid)
		if e == nil && !hasIssue(o.Issues, a.Issue.ID) && a.Issue.Status == StatusPending {
			o.Issues = append(o.Issues, a.Issue)
		}
	}
	sort.SliceStable(o.Issues, func(i, j int) bool {
		if o.Issues[i].WeightedScore == o.Issues[j].WeightedScore {
			return o.Issues[i].Title < o.Issues[j].Title
		}
		return o.Issues[i].WeightedScore > o.Issues[j].WeightedScore
	})
	return o, nil
}
func (s *Service) QuestionBySlug(ctx context.Context, slug, uid string) (QuestionResponse, error) {
	if slug == "" {
		return QuestionResponse{}, ErrNotFound
	}
	d, e := s.root().Collection("questionSlugs").Doc(slug).Get(ctx)
	if e != nil {
		return QuestionResponse{}, ErrNotFound
	}
	var v struct {
		QuestionID string `firestore:"questionId"`
	}
	if e = d.DataTo(&v); e != nil {
		return QuestionResponse{}, e
	}
	return s.Question(ctx, v.QuestionID, uid)
}
func hasIssue(v []Issue, id string) bool {
	for _, i := range v {
		if i.ID == id {
			return true
		}
	}
	return false
}
func (s *Service) OwnAnswer(ctx context.Context, uid, qid string) (AnswerResponse, error) {
	var o AnswerResponse
	qd, e := s.root().Collection("questions").Doc(qid).Get(ctx)
	if e != nil {
		return o, ErrNotFound
	}
	var q Question
	if e = qd.DataTo(&q); e != nil {
		return o, e
	}
	ad, e := s.root().Collection("questions").Doc(qid).Collection("answers").Doc(uid).Get(ctx)
	if e != nil {
		if isNotFound(e) {
			return o, ErrNotFound
		}
		return o, e
	}
	if e = ad.DataTo(&o.Answer); e != nil {
		return o, e
	}
	id, e := s.root().Collection("questions").Doc(o.Answer.QuestionID).Collection("issues").Doc(o.Answer.IssueID).Get(ctx)
	if e != nil {
		return o, e
	}
	if e = id.DataTo(&o.Issue); e != nil {
		return o, e
	}
	pd, e := s.root().Collection("personalAnswers").Doc(uid).Get(ctx)
	if e == nil {
		var p PersonalAnswer
		if pd.DataTo(&p) == nil {
			o.PersonalTop = p.PrioritySlotID == o.Answer.PrioritySlotID && p.QuestionID == o.Answer.QuestionID && p.IssueID == o.Answer.IssueID
		}
	}
	return o, nil
}

type operation struct {
	UID, PayloadHash string
	Response         AnswerResponse
	CreatedAt        time.Time
}
type alias struct {
	IssueID    string `firestore:"issueId"`
	Normalized string `firestore:"normalized"`
}
type rateLimit struct {
	WindowStart time.Time `firestore:"windowStart"`
	Count       int64     `firestore:"count"`
}
type txState struct {
	question           Question
	qref               *firestore.DocumentRef
	issue              Issue
	iref               *firestore.DocumentRef
	created            bool
	aliasRef, limitRef *firestore.DocumentRef
	limit              rateLimit
}

func (s *Service) Answer(ctx context.Context, caller Caller, req AnswerRequest) (AnswerResponse, error) {
	var result AnswerResponse
	uid := caller.UID
	if req.AnswerKind == "" {
		req.AnswerKind = AnswerKindCategory
	}
	if req.Attribution == "" {
		req.Attribution = "anonymous"
	}
	if uid == "" || req.QuestionID == "" || req.OperationID == "" || len(req.OperationID) > 100 || strings.Contains(req.OperationID, "/") || (req.AnswerKind != AnswerKindCategory && req.AnswerKind != AnswerKindPersonal) || (req.Attribution != "anonymous" && req.Attribution != "authored") || (req.Attribution == "authored" && strings.TrimSpace(caller.DisplayName) == "") {
		return result, fmt.Errorf("%w: invalid answer request", ErrValidation)
	}
	if e := s.requireEligible(ctx, caller); e != nil {
		return result, e
	}
	if (req.IssueID == "") == (strings.TrimSpace(req.Title) == "") {
		return result, fmt.Errorf("%w: supply exactly one of issueId or title", ErrValidation)
	}
	norm := ""
	var e error
	if req.Title != "" {
		norm, e = NormalizeTitle(req.Title)
		if e != nil {
			return result, fmt.Errorf("%w: %v", ErrValidation, e)
		}
	}
	raw := req.AnswerKind + "\x00" + req.QuestionID + "\x00" + req.IssueID + "\x00" + norm
	sum := sha256.Sum256([]byte(raw))
	ph := hex.EncodeToString(sum[:])
	now := s.now().UTC()
	e = s.db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		opref := s.root().Collection("operations").Doc(req.OperationID)
		if d, x := tx.Get(opref); x == nil {
			var op operation
			if x = d.DataTo(&op); x != nil {
				return x
			}
			if op.UID != uid || op.PayloadHash != ph {
				return ErrConflict
			}
			result = op.Response
			result.Replayed = true
			return nil
		} else if !isNotFound(x) {
			return x
		}
		st, x := s.loadTarget(ctx, tx, caller, req, norm, now)
		if x != nil {
			return x
		}
		if norm != "" && st.question.AnswerTargetType != "" && st.question.AnswerTargetType != "issue" {
			return fmt.Errorf("%w: free-form answers are available only for issue-target questions", ErrValidation)
		}
		aref := st.qref.Collection("answers").Doc(uid)
		var old Answer
		oldExists := false
		if d, x := tx.Get(aref); x == nil {
			oldExists = true
			if x = d.DataTo(&old); x != nil {
				return x
			}
		} else if !isNotFound(x) {
			return x
		}
		pref := s.root().Collection("personalAnswers").Doc(uid)
		var personal PersonalAnswer
		personalExists := false
		if d, x := tx.Get(pref); x == nil {
			personalExists = true
			if x = d.DataTo(&personal); x != nil {
				return x
			}
		} else if !isNotFound(x) {
			return x
		}
		personalDocWasPresent := personalExists
		refs := map[string]*firestore.DocumentRef{}
		issues := map[string]*Issue{}
		questions := map[string]*Question{}
		key := func(q, i string) string { return q + "/" + i }
		refs[key(req.QuestionID, st.issue.ID)] = st.iref
		issues[key(req.QuestionID, st.issue.ID)] = &st.issue
		questions[req.QuestionID] = &st.question
		load := func(qid, iid string) error {
			k := key(qid, iid)
			if _, ok := issues[k]; ok {
				return nil
			}
			qr := s.root().Collection("questions").Doc(qid)
			qd, x := tx.Get(qr)
			if x != nil {
				return x
			}
			var q Question
			if x = qd.DataTo(&q); x != nil {
				return x
			}
			ir := qr.Collection("issues").Doc(iid)
			id, x := tx.Get(ir)
			if x != nil {
				return x
			}
			var i Issue
			if x = id.DataTo(&i); x != nil {
				return x
			}
			questions[qid] = &q
			refs[k] = ir
			issues[k] = &i
			return nil
		}
		if oldExists {
			if x := load(old.QuestionID, old.IssueID); x != nil {
				return x
			}
		}
		if personalExists {
			if x := load(personal.QuestionID, personal.IssueID); x != nil {
				return x
			}
		}
		changed := !oldExists || old.QuestionID != req.QuestionID || old.IssueID != st.issue.ID
		personalTop := personalExists && personal.QuestionID == req.QuestionID && personal.IssueID == st.issue.ID
		if changed {
			if oldExists {
				oi := issues[key(old.QuestionID, old.IssueID)]
				oq := questions[old.QuestionID]
				if oi.Supporters < 1 || oi.WeightedScore < 1 || oq.TotalRespondents < 1 {
					return fmt.Errorf("counter invariant")
				}
				oi.Supporters--
				oi.WeightedScore--
				oq.TotalRespondents--
				oq.UpdatedAt = now
			}
			ni := issues[key(req.QuestionID, st.issue.ID)]
			nq := questions[req.QuestionID]
			ni.Supporters++
			ni.WeightedScore++
			ni.UpdatedAt = now
			nq.TotalRespondents++
			nq.UpdatedAt = now
		}
		if req.AnswerKind == AnswerKindCategory && changed && personalExists && personal.QuestionID == req.QuestionID {
			pi := issues[key(personal.QuestionID, personal.IssueID)]
			if pi.PersonalTopSupporters < 1 || pi.WeightedScore < s.personalWeight-1 {
				return fmt.Errorf("personal counter invariant")
			}
			pi.PersonalTopSupporters--
			pi.WeightedScore -= s.personalWeight - 1
			personalExists = false
			personalTop = false
		}
		if req.AnswerKind == AnswerKindPersonal && !personalTop {
			if personalExists {
				pi := issues[key(personal.QuestionID, personal.IssueID)]
				if pi.PersonalTopSupporters < 1 || pi.WeightedScore < s.personalWeight-1 {
					return fmt.Errorf("personal counter invariant")
				}
				pi.PersonalTopSupporters--
				pi.WeightedScore -= s.personalWeight - 1
			}
			ni := issues[key(req.QuestionID, st.issue.ID)]
			ni.PersonalTopSupporters++
			ni.WeightedScore += s.personalWeight - 1
			personal = PersonalAnswer{CategoryID: st.question.CategoryID, PrioritySlotID: st.question.PrioritySlotID, QuestionID: req.QuestionID, IssueID: st.issue.ID, UpdatedAt: now}
			personalExists = true
			personalTop = true
		}
		ans := Answer{CategoryID: st.question.CategoryID, PrioritySlotID: st.question.PrioritySlotID, QuestionID: req.QuestionID, IssueID: st.issue.ID, Revision: 1, CreatedAt: now, UpdatedAt: now}
		if oldExists {
			ans.Revision = old.Revision
			if changed {
				ans.Revision++
			}
			ans.CreatedAt = old.CreatedAt
		}
		result = AnswerResponse{Answer: ans, Issue: *issues[key(req.QuestionID, st.issue.ID)], Changed: oldExists && changed, PersonalTop: personalTop, CreatedIssue: st.created}
		if changed {
			if x := tx.Set(aref, ans); x != nil {
				return x
			}
		}
		for k, i := range issues {
			if x := tx.Set(refs[k], i); x != nil {
				return x
			}
		}
		for qid, q := range questions {
			if x := tx.Set(s.root().Collection("questions").Doc(qid), q); x != nil {
				return x
			}
		}
		if personalExists {
			if x := tx.Set(pref, personal); x != nil {
				return x
			}
		} else if personalDocWasPresent {
			if x := tx.Delete(pref); x != nil {
				return x
			}
		}
		if st.created {
			st.limit.Count++
			if x := tx.Set(st.aliasRef, alias{IssueID: st.issue.ID, Normalized: norm}); x != nil {
				return x
			}
			if x := tx.Set(st.limitRef, st.limit); x != nil {
				return x
			}
		}
		return tx.Set(opref, operation{UID: uid, PayloadHash: ph, Response: result, CreatedAt: now})
	})
	return result, e
}

type questionOperation struct {
	UID, PayloadHash string
	Response         CreateQuestionResponse
	CreatedAt        time.Time
}

func (s *Service) CreateQuestion(ctx context.Context, caller Caller, req CreateQuestionRequest) (CreateQuestionResponse, error) {
	var out CreateQuestionResponse
	if caller.UID == "" || req.OperationID == "" || len(req.OperationID) > 100 || strings.Contains(req.OperationID, "/") || (req.AnswerTargetType != "issue" && req.AnswerTargetType != "country") {
		return out, fmt.Errorf("%w: invalid question request", ErrValidation)
	}
	if e := s.requireEligible(ctx, caller); e != nil {
		return out, e
	}
	title := strings.Join(strings.Fields(strings.TrimSpace(req.Title)), " ")
	description := strings.TrimSpace(req.Description)
	if len([]rune(title)) < 10 || len([]rune(title)) > 180 || !strings.HasSuffix(title, "?") || len([]rune(description)) < 20 || len([]rune(description)) > 1000 || urlPattern.MatchString(title) {
		return out, fmt.Errorf("%w: question must be 10-180 characters ending in ?, with a 20-1000 character description", ErrValidation)
	}
	slug := Slugify(strings.TrimSuffix(title, "?"))
	raw := title + "\x00" + description + "\x00" + req.AnswerTargetType
	sum := sha256.Sum256([]byte(raw))
	ph := hex.EncodeToString(sum[:])
	now := s.now().UTC()
	err := s.db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		opref := s.root().Collection("questionOperations").Doc(req.OperationID)
		if d, e := tx.Get(opref); e == nil {
			var op questionOperation
			if e = d.DataTo(&op); e != nil {
				return e
			}
			if op.UID != caller.UID || op.PayloadHash != ph {
				return ErrConflict
			}
			out = op.Response
			out.Replayed = true
			return nil
		} else if !isNotFound(e) {
			return e
		}
		slugref := s.root().Collection("questionSlugs").Doc(slug)
		if _, e := tx.Get(slugref); e == nil {
			return fmt.Errorf("%w: question slug already exists", ErrConflict)
		} else if !isNotFound(e) {
			return e
		}
		limitref := s.root().Collection("actors").Doc(caller.UID).Collection("limits").Doc("questions")
		var limit rateLimit
		if d, e := tx.Get(limitref); e == nil {
			if e = d.DataTo(&limit); e != nil {
				return e
			}
		} else if !isNotFound(e) {
			return e
		}
		if limit.WindowStart.IsZero() || now.Sub(limit.WindowStart) >= 24*time.Hour {
			limit = rateLimit{WindowStart: now}
		}
		if limit.Count >= 3 {
			return ErrRateLimited
		}
		qref := s.root().Collection("questions").NewDoc()
		q := Question{ID: qref.ID, Slug: slug, Title: title, Description: description, Status: StatusPending, Indexable: false, AnswerTargetType: req.AnswerTargetType, CreatorUID: caller.UID, UpdatedAt: now}
		out = CreateQuestionResponse{Question: q}
		limit.Count++
		if e := tx.Set(qref, q); e != nil {
			return e
		}
		if e := tx.Set(slugref, map[string]any{"questionId": q.ID, "createdAt": now}); e != nil {
			return e
		}
		if e := tx.Set(limitref, limit); e != nil {
			return e
		}
		return tx.Set(opref, questionOperation{UID: caller.UID, PayloadHash: ph, Response: out, CreatedAt: now})
	})
	return out, err
}
func (s *Service) loadTarget(ctx context.Context, tx *firestore.Transaction, caller Caller, req AnswerRequest, norm string, now time.Time) (txState, error) {
	uid := caller.UID
	var st txState
	st.qref = s.root().Collection("questions").Doc(req.QuestionID)
	d, e := tx.Get(st.qref)
	if e != nil {
		return st, ErrNotFound
	}
	if e = d.DataTo(&st.question); e != nil {
		return st, e
	}
	if st.question.Status != StatusPublished {
		return st, ErrNotFound
	}
	st.iref = st.qref.Collection("issues").Doc(req.IssueID)
	if norm != "" {
		st.aliasRef = st.qref.Collection("aliases").Doc(normalizedHash(norm))
		if d, e = tx.Get(st.aliasRef); e == nil {
			var a alias
			if e = d.DataTo(&a); e != nil {
				return st, e
			}
			st.iref = st.qref.Collection("issues").Doc(a.IssueID)
		} else if isNotFound(e) {
			st.limitRef = s.root().Collection("actors").Doc(uid).Collection("limits").Doc("candidate-issues")
			if d, e = tx.Get(st.limitRef); e == nil {
				if e = d.DataTo(&st.limit); e != nil {
					return st, e
				}
			} else if !isNotFound(e) {
				return st, e
			}
			if st.limit.WindowStart.IsZero() || now.Sub(st.limit.WindowStart) >= 24*time.Hour {
				st.limit = rateLimit{WindowStart: now}
			}
			if st.limit.Count >= 5 {
				return st, ErrRateLimited
			}
			st.iref = st.qref.Collection("issues").NewDoc()
			author := ""
			if req.Attribution == "authored" {
				author = strings.TrimSpace(caller.DisplayName)
			}
			st.issue = Issue{ID: st.iref.ID, Slug: Slugify(req.Title), Title: strings.Join(strings.Fields(req.Title), " "), Status: StatusPending, Attribution: req.Attribution, AuthorDisplayName: author, CreatorUID: uid, CreatedAt: now, UpdatedAt: now}
			st.created = true
		} else {
			return st, e
		}
	}
	if !st.created {
		d, e = tx.Get(st.iref)
		if e != nil {
			return st, ErrNotFound
		}
		if e = d.DataTo(&st.issue); e != nil {
			return st, e
		}
		if st.issue.Status != StatusPublished && !(st.issue.Status == StatusPending && st.issue.CreatorUID == uid) {
			return st, ErrNotFound
		}
	}
	return st, nil
}
