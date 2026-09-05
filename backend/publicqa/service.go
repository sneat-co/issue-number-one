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
	"github.com/sneat-co/issue-number-one/backend/translations"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isNotFound(err error) bool { return status.Code(err) == codes.NotFound }

var supportedLanguages = func() map[string]bool {
	result := make(map[string]bool)
	for _, language := range translations.CanonicalSupportedLanguages() {
		result[language] = true
	}
	return result
}()

func normalizeLanguage(v string) (string, error) {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return "en", nil
	}
	if i := strings.IndexByte(v, '-'); i > 0 {
		v = v[:i]
	}
	if !supportedLanguages[v] {
		return "", fmt.Errorf("%w: unsupported languageCode", ErrValidation)
	}
	return v, nil
}
func langAgg(i *Issue, lang string) *LanguageAggregate {
	if i.LanguageStats == nil {
		i.LanguageStats = map[string]LanguageAggregate{}
	}
	v := i.LanguageStats[lang]
	return &v
}
func saveLangAgg(i *Issue, lang string, v *LanguageAggregate) { i.LanguageStats[lang] = *v }

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

// WithTranslations enables revision-aware stored question translations. Reads
// never invoke the external translator; they use the repository only.
func WithTranslations(service *translations.Service, repository translations.Repository) Option {
	return func(s *Service) {
		s.translationService = service
		s.translationRepository = repository
	}
}

type Service struct {
	db                    *firestore.Client
	spaceID               string
	now                   func() time.Time
	personalWeight        int64
	translationService    *translations.Service
	translationRepository translations.Repository
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
	return s.QuestionLocalized(ctx, qid, uid, "en")
}
func (s *Service) QuestionLocalized(ctx context.Context, qid, uid, language string) (QuestionResponse, error) {
	var o QuestionResponse
	language, e := normalizeLanguage(language)
	if e != nil {
		return o, e
	}
	o.LanguageCode = language
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
	o.LanguageRespondents = o.Question.RespondentsByLanguage[language]
	o.UpdatedAt = o.Question.UpdatedAt
	o.ContentLanguage = o.Question.SourceLanguage
	if o.ContentLanguage == "" {
		o.ContentLanguage = "en"
	}
	if language != o.ContentLanguage {
		o.TranslationFallback = true
		if s.translationRepository != nil {
			if translated, x := s.translationRepository.GetTranslation(ctx, qid, language); x == nil && translationIsFresh(o.Question, translated) && strings.TrimSpace(translated.Title) != "" {
				o.Translation = &translated
				o.ContentLanguage = language
				o.TranslationFallback = false
			}
		}
	}
	var it *firestore.DocumentIterator
	if o.Question.Status == StatusPending && o.Question.CreatorUID == uid {
		it = qref.Collection("issues").Documents(ctx)
	} else {
		it = qref.Collection("issues").Where("status", "==", StatusPublished).Documents(ctx)
	}
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
		la := i.LanguageStats[language]
		o.Issues[len(o.Issues)-1].LanguageSupporters = la.Supporters
		o.Issues[len(o.Issues)-1].LanguagePersonalTopSupporters = la.PersonalTopSupporters
		o.Issues[len(o.Issues)-1].LanguageWeightedScore = la.WeightedScore
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
	return s.QuestionBySlugLocalized(ctx, slug, uid, "en")
}
func (s *Service) QuestionBySlugLocalized(ctx context.Context, slug, uid, language string) (QuestionResponse, error) {
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
	return s.QuestionLocalized(ctx, v.QuestionID, uid, language)
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
	language, e := normalizeLanguage(req.LanguageCode)
	if e != nil {
		return result, e
	}
	req.LanguageCode = language
	norm := ""
	if req.Title != "" {
		norm, e = NormalizeTitle(req.Title)
		if e != nil {
			return result, fmt.Errorf("%w: %v", ErrValidation, e)
		}
	}
	raw := req.AnswerKind + "\x00" + req.QuestionID + "\x00" + req.IssueID + "\x00" + norm + "\x00" + language + "\x00" + req.Attribution
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
		targetLanguage := language
		if !changed && oldExists {
			targetLanguage = old.LanguageCode
			if targetLanguage == "" {
				targetLanguage = "en"
			}
		}
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
				oldLang := old.LanguageCode
				if oldLang == "" {
					oldLang = "en"
				}
				ola := langAgg(oi, oldLang)
				if ola.Supporters < 1 || ola.WeightedScore < 1 {
					return fmt.Errorf("language counter invariant")
				}
				ola.Supporters--
				ola.WeightedScore--
				saveLangAgg(oi, oldLang, ola)
				oq.TotalRespondents--
				if oq.RespondentsByLanguage == nil {
					oq.RespondentsByLanguage = map[string]int64{}
				}
				if oq.RespondentsByLanguage[oldLang] < 1 {
					return fmt.Errorf("question language counter invariant")
				}
				oq.RespondentsByLanguage[oldLang]--
				oq.UpdatedAt = now
			}
			ni := issues[key(req.QuestionID, st.issue.ID)]
			nq := questions[req.QuestionID]
			ni.Supporters++
			ni.WeightedScore++
			nla := langAgg(ni, targetLanguage)
			nla.Supporters++
			nla.WeightedScore++
			saveLangAgg(ni, targetLanguage, nla)
			ni.UpdatedAt = now
			nq.TotalRespondents++
			if nq.RespondentsByLanguage == nil {
				nq.RespondentsByLanguage = map[string]int64{}
			}
			nq.RespondentsByLanguage[targetLanguage]++
			nq.UpdatedAt = now
		}
		if req.AnswerKind == AnswerKindCategory && changed && personalExists && personal.QuestionID == req.QuestionID {
			pi := issues[key(personal.QuestionID, personal.IssueID)]
			if pi.PersonalTopSupporters < 1 || pi.WeightedScore < s.personalWeight-1 {
				return fmt.Errorf("personal counter invariant")
			}
			pi.PersonalTopSupporters--
			pi.WeightedScore -= s.personalWeight - 1
			pl := personal.LanguageCode
			if pl == "" {
				pl = "en"
			}
			pla := langAgg(pi, pl)
			if pla.PersonalTopSupporters < 1 || pla.WeightedScore < s.personalWeight-1 {
				return fmt.Errorf("personal language counter invariant")
			}
			pla.PersonalTopSupporters--
			pla.WeightedScore -= s.personalWeight - 1
			saveLangAgg(pi, pl, pla)
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
				pl := personal.LanguageCode
				if pl == "" {
					pl = "en"
				}
				pla := langAgg(pi, pl)
				if pla.PersonalTopSupporters < 1 || pla.WeightedScore < s.personalWeight-1 {
					return fmt.Errorf("personal language counter invariant")
				}
				pla.PersonalTopSupporters--
				pla.WeightedScore -= s.personalWeight - 1
				saveLangAgg(pi, pl, pla)
			}
			ni := issues[key(req.QuestionID, st.issue.ID)]
			ni.PersonalTopSupporters++
			ni.WeightedScore += s.personalWeight - 1
			nla := langAgg(ni, targetLanguage)
			nla.PersonalTopSupporters++
			nla.WeightedScore += s.personalWeight - 1
			saveLangAgg(ni, targetLanguage, nla)
			personal = PersonalAnswer{CategoryID: st.question.CategoryID, PrioritySlotID: st.question.PrioritySlotID, QuestionID: req.QuestionID, IssueID: st.issue.ID, LanguageCode: targetLanguage, UpdatedAt: now}
			personalExists = true
			personalTop = true
		}
		ans := Answer{CategoryID: st.question.CategoryID, PrioritySlotID: st.question.PrioritySlotID, QuestionID: req.QuestionID, IssueID: st.issue.ID, LanguageCode: targetLanguage, Revision: 1, CreatedAt: now, UpdatedAt: now}
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
	if req.ChoiceSource.Kind == "" {
		if req.AnswerTargetType == "country" {
			req.ChoiceSource = ChoiceSource{Kind: "predefined", EntityType: "country"}
		} else if req.AnswerTargetType == "issue" {
			req.ChoiceSource = ChoiceSource{Kind: "free"}
		}
	}
	if caller.UID == "" || req.OperationID == "" || len(req.OperationID) > 100 || strings.Contains(req.OperationID, "/") {
		return out, fmt.Errorf("%w: invalid question request", ErrValidation)
	}
	sourceLanguage, e := normalizeLanguage(req.SourceLanguage)
	if e != nil || (s.translationService != nil && !containsString(s.translationService.SupportedLanguages(), sourceLanguage)) {
		return out, fmt.Errorf("%w: sourceLanguage is not enabled for translation", ErrValidation)
	}
	req.SourceLanguage = sourceLanguage
	if e := validateChoiceSource(req.ChoiceSource); e != nil {
		return out, e
	}
	if req.ChoiceSource.Kind == "free" {
		req.AllowSuggestions = true
	}
	req.AnswerTargetType = req.ChoiceSource.EntityType
	if req.ChoiceSource.Kind == "free" || req.ChoiceSource.Kind == "custom" {
		req.AnswerTargetType = "issue"
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
	raw := title + "\x00" + description + "\x00" + req.SourceLanguage + "\x00" + req.ChoiceSource.Kind + "\x00" + req.ChoiceSource.EntityType + fmt.Sprint(req.ChoiceSource.Options, req.AllowSuggestions)
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
		q := Question{ID: qref.ID, Slug: slug, Title: title, Description: description, Status: StatusPending, Publication: StatusPending, Indexable: false, AnswerTargetType: req.AnswerTargetType, ChoiceSource: req.ChoiceSource, AllowSuggestions: req.AllowSuggestions, ContentRevision: 1, SourceLanguage: sourceLanguage, TranslationStatus: "pending", AvailableLanguages: []string{sourceLanguage}, CreatorUID: caller.UID, UpdatedAt: now}
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
		for _, option := range req.ChoiceSource.Options {
			iref := qref.Collection("issues").NewDoc()
			issue := Issue{ID: iref.ID, Slug: Slugify(option.Title), Title: strings.Join(strings.Fields(option.Title), " "), Description: strings.TrimSpace(option.Description), Status: StatusPending, TargetType: "custom", CreatorUID: caller.UID, CreatedAt: now, UpdatedAt: now}
			if e := tx.Set(iref, issue); e != nil {
				return e
			}
			norm, _ := NormalizeTitle(option.Title)
			if e := tx.Set(qref.Collection("aliases").Doc(normalizedHash(norm)), alias{IssueID: issue.ID, Normalized: norm}); e != nil {
				return e
			}
		}
		return tx.Set(opref, questionOperation{UID: caller.UID, PayloadHash: ph, Response: out, CreatedAt: now})
	})
	if err != nil || s.translationService == nil {
		return out, err
	}
	// Translation is intentionally outside the canonical creation transaction.
	// A provider failure cannot erase a valid question and is durably retryable.
	translated, translationErr := s.TranslateAllQuestionLanguages(ctx, out.Question.ID, translations.Actor{UID: caller.UID})
	if translationErr != nil {
		out.Question.TranslationStatus = "failed"
		return out, nil
	}
	out.Question.TranslationStatus = "ready"
	out.Question.AvailableLanguages = translated
	return out, nil
}

// TranslateAllQuestionLanguages is the trusted retry/backfill boundary. The
// translation package makes writes idempotent by content revision and hash.
func (s *Service) TranslateAllQuestionLanguages(ctx context.Context, questionID string, actor translations.Actor) ([]string, error) {
	if s.translationService == nil {
		return nil, errors.New("question translation is not configured")
	}
	values, err := s.translationService.TranslateEnabledQuestion(ctx, questionID, actor)
	statusValue := "ready"
	if err != nil {
		statusValue = "failed"
	}
	languages := make([]string, 0, len(values))
	for _, value := range values {
		if value.Language != "" && strings.TrimSpace(value.Title) != "" {
			languages = append(languages, value.Language)
		}
	}
	update := map[string]any{"translationStatus": statusValue, "translationUpdatedAt": s.now().UTC()}
	if err == nil {
		update["availableLanguages"] = languages
	}
	if _, updateErr := s.root().Collection("questions").Doc(questionID).Set(ctx, update, firestore.MergeAll); updateErr != nil {
		return languages, updateErr
	}
	return languages, err
}

func translationIsFresh(question Question, value translations.Translation) bool {
	sourceLanguage := question.SourceLanguage
	if sourceLanguage == "" {
		sourceLanguage = "en"
	}
	source := translations.QuestionSource{ID: question.ID, Title: question.Title, Description: question.Description, SourceLanguage: sourceLanguage, ContentRevision: question.ContentRevision, Publication: question.Publication, CreatorUID: question.CreatorUID}
	return value.SourceRevision == question.ContentRevision && value.SourceHash == translations.SourceHash(source, sourceLanguage)
}

func containsString(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
func validateChoiceSource(source ChoiceSource) error {
	switch source.Kind {
	case "predefined":
		if source.EntityType != "country" && source.EntityType != "city" && source.EntityType != "currency" || len(source.Options) > 0 {
			return fmt.Errorf("%w: invalid predefined choice source", ErrValidation)
		}
	case "custom":
		if source.EntityType != "" || len(source.Options) < 2 || len(source.Options) > 30 {
			return fmt.Errorf("%w: custom choices require 2-30 options", ErrValidation)
		}
		seen := map[string]bool{}
		for _, o := range source.Options {
			n, e := NormalizeTitle(o.Title)
			if e != nil {
				return fmt.Errorf("%w: invalid custom option", ErrValidation)
			}
			if seen[n] {
				return fmt.Errorf("%w: duplicate custom option", ErrValidation)
			}
			seen[n] = true
			if len([]rune(o.Description)) > 500 {
				return fmt.Errorf("%w: custom option description too long", ErrValidation)
			}
		}
	case "free":
		if source.EntityType != "" || len(source.Options) > 0 {
			return fmt.Errorf("%w: invalid free choice source", ErrValidation)
		}
	default:
		return fmt.Errorf("%w: invalid choice source kind", ErrValidation)
	}
	return nil
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
	if st.question.MutationLock != "" {
		return st, ErrConflict
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
			if !st.question.AllowSuggestions && st.question.ChoiceSource.Kind != "free" {
				return st, fmt.Errorf("%w: this question does not accept suggested choices", ErrValidation)
			}
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
