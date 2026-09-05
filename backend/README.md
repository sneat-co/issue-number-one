# IssueNumber.one public Q&A backend

This module owns the public priority domain. Sneat-Go supplies a Firestore
client and a Firebase-token-backed `publicqa.UserIdentity`, constructs
`publicqa.NewService`, and mounts `publicqa.NewHandler(...).RegisterHttpRoutes`.
Clients cannot write Firestore.

## Public HTTP contract

- `GET /v0/issuenumber/catalog` returns published categories, concepts, and
  questions. It contains no answers or actor identifiers.
- `GET /v0/issuenumber/question?questionId=...` (or `?slug=...`) returns the published question,
  genuine unweighted respondent total, and bounded issue ranking fields:
  `supporters`, `personalTopSupporters`, and `weightedScore`.
- `GET /v0/issuenumber/answer?questionId=...` returns the caller's current
  answer for that concrete question and whether it is their personal top issue.
- `GET /v0/issuenumber/eligibility` returns private phone/payment eligibility.
- `POST /v0/issuenumber/answer` accepts `answerKind` (`category`, the default,
  or `personal`), `questionId`, exactly one of `issueId`/`title`, `operationId`,
  and optional `attribution` (`anonymous`, the default, or `authored`).

Country, county, and city are nested for URLs and discovery, but answer
uniqueness is per concrete question. A caller can answer Ireland, France,
Dublin, and Cork independently. Promoting an issue to `personal` atomically
makes it the question answer and the single personal top across all questions.
Replacing that exact question's answer clears a displaced personal designation;
it never promotes the new issue silently. Parent totals are labelled answers,
because one person can contribute to several child scopes.

The current score assumption is one point for a question answer and ten total
points for the personal top (a bonus of nine). `WithPersonalWeight` makes this
configurable. Counts remain genuine people counts and never include the weight.

## Firestore records

All paths are beneath `/spaces/{spaceID}/ext/issuenumber`:

- `categories`, `concepts`, `questions`, and per-question `issues`/`aliases`;
- `questions/{questionId}/answers/{uid}` for the private current choice;
- `personalAnswers/{uid}` for the private single personal top;
- `operations/{operationId}` for actor/payload-bound idempotency;
- `actors/{uid}/limits/candidate-issues` for candidate rate limiting;
- `verification/{uid}` and `paymentReceipts/{chargeID}` for private paid
  eligibility and settlement replay binding.

Phone verification is trusted identity metadata. Payment verification is
granted only by the non-HTTP `Service.MarkPaid` integration after the payment
consumer proves a settled EUR 1.00 charge. Either route grants eligibility;
using both grants no extra answer or score. This reduces cheap manipulation but
is not described as a unique-human guarantee.

New free-form issues are immediately usable by their creator but remain
`pending`, private through the API, and unindexed. Authored attribution is an
explicit opt-in and uses only a trusted display name supplied by the identity
adapter. Creator UID, email, answers, operations, and payment records never
appear in public DTOs.

The seed command imports `catalog/seed.json` schema version 1. It refuses a
real project unless `--confirm-production-project` exactly matches `--project`.
Reruns update curated copy and aliases while preserving counts, answers, issue
moderation status, authorship, and payment state.
