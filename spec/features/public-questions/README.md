---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Public questions and primary issues

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/public-questions?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/public-questions?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/public-questions?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/public-questions?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

IssueNumber.one helps a person state their most important issue in a real scope,
understand the priorities of other participants, and invite people whose views
matter to them. Issues persist as objects; this is not an arbitrary survey builder.
The public acquisition surface leads naturally to Sneat.app relationships and
Sneat.work audience-specific feedback, without claiming those products already
provide private question workflows.

## Source and decision boundary

The September 5, 2026 public-data implementation request explicitly authorizes
Firestore persistence and this public acquisition surface. It supersedes the older
consulting discovery-only authority for this effort, without claiming market validation.
The existing organisational creator-star and vote-budget model is retained.
Public Q&A is a separate participation mode. The founder clarified that public
participation has one issue per concrete scoped question and one personal top issue
across all slots. Country, county and city nest in navigation and URLs. Each concrete scope owns
an independent choice: Ireland, France, Dublin and Cork can all have separate answers. Organisational vote budgets and propagation are presented on
`/for-work/`; they do not govern public participation.

**Weight assumption awaiting confirmation:** a scoped-question choice contributes one
point; designation as the overall personal top replaces that with ten points total,
not eleven. Raw people counts always remain separate from weighted priority points.

## User journey and observable results

1. “What is the #1 issue?” A visitor arrives from search or the homepage. Without
   signing in or clicking anything, useful question copy, candidates, real counts,
   methodology and navigation appear in the HTTP HTML.
2. “What is the top issue in your country?” The category shows aggregated child
   responses and a labelled scope selector. Choosing Ireland opens Ireland; no
   answer is stored against the generic template and no IP sets a profile country.
3. “For me too.” The visitor chooses a candidate or types their issue before auth.
   The pending action survives login/signup and executes automatically once auth
   succeeds. On failure it remains retryable; no count is shown as saved prematurely.
4. “This is my #1.” One current choice exists for the verified person/concrete question.
   Switching it, moves the old/new issue counts atomically. Choices for different concrete scopes do not replace one another.
   Retrying an already completed action cannot cast another vote.
5. “What do people I care about think?” After answering, a voluntary question/issue
   share invite lets another visitor follow the same journey. The original visitor
   can continue to the existing Sneat.app contacts experience. Business categories
   offer Sneat.work instead. No contacts are uploaded or invitations sent automatically.
6. Divergent epilogues: close the page (answer persists); return and change the answer
   (current ranking changes); share (canonical URL opens readable content); follow
   the contextual product link (existing destination opens with attribution).

## Architecture

Cloudflare serves public HTML before auth JavaScript starts. The existing Astro
landing remains the content/design base; interactive answer continuation uses the
shared Angular/Firebase Sneat authentication. The extension owns backend domain and
persistence in `issue-number-one/backend`; `sneat-go` only mounts handlers and binds
platform identity/storage. All mutations use authenticated sneat-go HTTPS endpoints.

Space boundary: a platform-managed public Space contains the public catalog under
`/spaces/{publicSpaceID}/ext/issuenumber/`. A public read projection contains only
published content and aggregate counts. Future organisation/household records use
their existing Space and membership authority, never a new owner/user container.

## Question axes

The founder proposed comparing which country creates the most issues for the world.
Distinguish the evaluation scope (world/Ireland/team), subject axis (issues/countries/
organisations/products), and question intent (experienced priority/attributed concern).
A scoped issue question returns an Issue; a country-comparison question returns a
country reference with optional linked Issue explanation. Countries are not Issue
concepts. Do not encode a country comparison by creating fake issue concepts named
USA/Iran/Russia/Vatican, and never label perceived concern as measured causal harm.

A proposed public wording is “Which country's actions cause you the most concern
globally?” with country search, no preselected answer and an optional explanation.
Offer the complete eligible country set rather than a prejudicial hand-picked four.
Explain whether the question concerns government/state actions rather than inhabitants.
Use stable country IDs; display-name/translation changes do not create new identities.
The data model can add answerTargetType and targetRef to a distinct question family;
its aggregation must not be mixed with issue-priority rankings. Scope hierarchy and
answer target taxonomy are independent. The seed demonstrates this axis with the
neutral question “Which country's government actions cause you the most concern
globally?” backed by canonical country entities. Broader axis analysis remains a
separate Feature and is not permission to turn the current MVP into an arbitrary
poll builder.

## Firestore model

| Relative path beneath extension           | Meaning                                                                  |
| ----------------------------------------- | ------------------------------------------------------------------------ |
| `categories/{id}`                         | Generic template, slug, scope type, SEO and intent                       |
| `concepts/{id}`                           | Shared issue identity, English presentation and curated aliases          |
| `questions/{id}`                          | Concrete scope, category and generic parent references, publication      |
| `questions/{id}/issues/{issueId}`         | Contextual title/description, optional conceptId, moderation, counters   |
| `questions/{questionId}/answers/{uid}`    | Private current questionId/issueId, revision, createdAt and updatedAt    |
| `personalAnswers/{uid}`                   | One overall personal top referring to a current slot choice              |
| `verification/{uid}`                      | Private trusted payment eligibility; phone comes from identity authority |
| `questions/{id}/aliases/{normalisedHash}` | Trusted normalized text to issue identity mapping                        |
| `operations/{operationId}`                | Private answer idempotency receipt bound to actor and payload            |
| `questionOperations/{operationId}`        | Private question-creation idempotency receipt                            |

IDs are immutable; editorial slugs are routes, not database authority. A concept can
occur in many question memberships. New free-form candidates have no automatically
asserted global equivalence. Curated aliases join obvious synonyms; semantic clustering
and cross-question merging need explicit moderation. A merge retains old identity
with mergedIntoIssueId and old route redirect; it must reconcile answer counts.

Geo scope identity uses country codes and stable city identifiers, with geographic
parent references independent of aggregation edges. City votes are not also counted
in a country category through geographic ancestry. Generic hierarchy must be acyclic;
aggregation follows one explicit parent tree rather than every related link.

## Answering and integrity

(spaceId, questionId, uid) identifies one current answer, where questionId denotes
a concrete scoped question. This follows the founder's explicit clarification:
“Yes. I can have issue per each country and per each city”. A user may independently
choose an issue in Ireland, France, Dublin and Cork. The preceding interpretation of
one shared country slot was the agent's mistake, not a founder decision.

One private personal-top record per user references one current scoped-question answer.
Marking an issue as personal top ensures that question's answer and updates the
designation atomically. Changing an answer that was personal top clears that designation;
the system never silently promotes its replacement. Changing another question preserves
personal top. Navigation/category depth never defines mutual exclusion.

Trusted transactions validate publication, membership, identity and participation
eligibility; read before writing; move counts across both old/new questions and
issues. Public DTOs distinguish supporters, personalTopSupporters and weightedScore.
Clients never supply authoritative totals, actor IDs, payment status or phone status.
Same selection is a no-op; operation receipts bind actor, operation kind and payload.
Delayed replay must not undo a newer choice. Transaction tests cover competing requests
for the same question, independent different-question answers, and global personal-top races.

### Participation verification and authorship

Founder decision: either verified phone OR a one-time EUR 1 card payment qualifies
participation. Payment supports the project and makes mass participation more costly;
it does not buy extra votes or prove unique personhood. Phone verification also
provides friction, not a guarantee against multi-account manipulation. Authentication
and a draft precede verification, which precedes authoritative count changes.

Offer both methods with equal prominence and neither selected. Randomize their display
order once per attempt using browser cryptographic randomness, persist that assignment
through retries and auth/payment returns, and measure exposure, selection, completion,
failure and abandonment by known category/question IDs. Order is an
experiment dimension so first-position bias can be distinguished from preference.
No analytics phone numbers, card details, free-form issue text or private graph edges.
The card amount is fixed server-side at EUR 1. Only the existing payrail signed
settlement pipeline grants payment eligibility, never the return URL or client flag.
Reuse Firebase account-linked phone verification, preserving the same authenticated UID.

Issue attribution is separate: authored or anonymous to public readers. Retain creator
identity privately for moderation. Never publish a UID, phone or email as an author
fallback. Public attribution needs explicit choice and a trusted display name. Pending
submissions cannot expose other users' identities or unreviewed text publicly.

Normalize Unicode, trim and collapse whitespace, lower-case for equality; curated
aliases handle common variants such as living costs/cost of living. Reject controls,
URLs, empty and oversized titles. Candidate creation has an authenticated actor rate
limit and remains pending/unindexed until editorial publication. The submitter can
see their own pending selection without exposing other private candidate answers.

## Aggregation and ranking

Current question answers, not cumulative clicks. Rankings use separately labelled weighted points. Parent raw aggregation groups child
question memberships by canonical concept. Label it “Responses across participating
scopes”; each question includes a person at most once, but a category sum counts answers and
may include the same person across multiple scopes. Never label parent totals as unique
people or a representative population sample.
Percentages use the stated eligible-answer denominator. Zero counts have no winner.
Equal counts are ties; stable editorial ordering only breaks display ties.

A separate metric counts child scopes where each concept is joint first (positive
counts only), labelled “Joint or sole #1 in N participating scopes”. It is never
presented as respondent share. Equal-scope weighting and average rank remain future
analyses with explicit missing-scope rules. Current answers and timestamped operation
receipts provide reconciliation input for future daily buckets; the MVP does not claim
historical movement or retain an unbounded event log.

## Privacy and relationship bridge

Contactus is the authority: Space-bounded contacts, linked users and typed relationships.
The first bridge is user-initiated sharing and entry into the real contacts experience;
no copied social graph, fabricated contact counts, or public respondent identities.
Future graph filtering resolves authorized linked users server-side, intersects only
answers with explicit relationship visibility consent, and returns aggregate DTOs.
Minimum cohort: at least 10 consenting respondents, at least 5 in a disclosed cell;
suppress complementary cells and arbitrary filters. These are initial design defaults,
not a guarantee against inference; conduct a privacy review before enabling statistics.
Distance 0 is self, 1 direct authorized contacts, 2 separately consented traversal.
Space/team/company membership and contact edges never imply answer-sharing consent.

## User-created questions

Founder requirement: users can create questions. Canonical path namespace is
`/questions/{slug}` (correcting the illustrative spelling `/quesrions/`). Founder
also requested a visible question mark at the end: generated question links use
`/questions/{slug}?`. The final `?` is an empty query delimiter, not part of the
stored slug; both forms resolve to the same question. The founder explicitly requires canonical SEO URLs to include that trailing `?` too.
Canonical tags, Open Graph URLs, sitemap entries and generated links must agree.
Requests with or without the empty delimiter resolve to the same underlying question.

Draft title/description and answer source (predefined country/city/currency, custom choices,
or free answers) are editable before
authentication. Sign-in and either phone or EUR 1 verification preserve the draft.
A trusted, rate-limited create operation reserves a safe immutable route slug and
separate random document ID, creates a pending/unindexed question, and routes its
creator to the preview. Only moderation can publish or index it. Predefined answers reference shared entity IDs. Custom choices belong to the question
and require 2–30 distinct options. Free answers use normalization/duplicate handling.
Allowing additional suggestions is a separate option. City lists must disclose their
coverage rather than claim a nonexistent complete global city catalog. No survey
appearance builder or multiple-choice voting mode is introduced.

## SEO and discovery

`/issues` is discovery. Geographic category navigation nests country → county → city.
Concrete canonical paths follow scope parents, for example
`/issues/country/ireland/county/limerick/city/limerick`. Issue slugs append to their
question path. Other categories use `/issues/{category}/{scope}`. URLs and visual
nesting do not determine per-question uniqueness or aggregation edges.
The public host supplies title, one h1, description, canonical, Open Graph/Twitter,
breadcrumbs and crawlable content. Use WebPage/CollectionPage, BreadcrumbList and
ItemList; do not claim accepted answers or QAPage rich-result eligibility for rankings.
Related pages follow category/scope/concept relationships, not random links.
Curated questions with explanatory text and at least five useful candidates can be
indexed at zero participation, explicitly labelled candidate options. User-created
thin pages remain noindex. Editorial approval plus substantive context is required
for issue indexing; counts alone are not a quality or truth signal. Draft/hidden/private
content is inaccessible publicly. Canonicals omit tracking and filter parameters; the
founder-required empty `?` and explicit `?lang={code}` are the only question URL
variants. Redirects preserve approved aliases, and the sitemap includes only canonical
indexable pages. Later partition by
category with a sitemap index rather than generating filter combinations.

## Performance and scale

A question reads one bounded ranking page plus its summary, never its answers.
Public responses are cacheable briefly; own answers and relationship responses are
private/no-store. No auth initialization blocks first HTML. Synchronous counters are
appropriate initially but hot questions will contend: measure retries, then use
sharded deltas and versioned ranking projections with reconciliation. Never put an
unbounded answer/issue/relationship array on one document. Page large lists and cap
request sizes. Rebuild projections from authoritative answers as an administrative job.

## Analytics

Use existing analytics transport with allowlisted metadata: category/question/issue
IDs, intent, action kind and outcome, never free-form text or relationship edges.
Events: question_view, issue_view, for_me_too_clicked, answer_submitted, answer_changed,
freeform_started, freeform_submitted, auth_required_for_answer,
auth_completed_for_answer, question_shared, issue_shared,
relationship_results_viewed (only when actual results exist), sneat_app_cta_clicked,
sneat_work_cta_clicked. Compare public acquisition → answer → auth → voluntary graph
bridge → downstream consumer/business visit. Do not count CTA clicks as downstream use.

## Acceptance Criteria

- Domain, seed-idempotence and normalization tests verify identity/reference integrity.
- Firestore emulator integration verifies atomic changes, concurrent idempotency and
  private-answer/public-projection separation; rules deny all client mutations.
- Browser E2E traverses both anonymous existing and typed-answer auth journeys against
  the real local auth/persistence mechanism, not mocked successful API responses.
- HTTP SEO tests check representative category/scope/issue pages, canonical metadata,
  noindex policy, sitemap, headings, escaped user content and unavailable-data states.
- Final review covers product distinctiveness, non-intrusive conversion, network value,
  privacy, abuse, future dimensions, accessibility and server-rendered performance.

## Open Questions

- Confirm ten points total for a personal top versus a different weighting; this
  remains an implementation assumption, not a founder ruling.
- Future cohort consent UX and private question creation belong to their owning
  Features; the MVP links to existing product destinations without promising delivery.

## References

- [Firestore transactions](https://firebase.google.com/docs/firestore/manage-data/transactions)
- [Google QAPage guidance](https://developers.google.com/search/docs/appearance/structured-data/qapage)
- [Cloudflare Workers practices](https://developers.cloudflare.com/workers/best-practices/workers-best-practices/)
