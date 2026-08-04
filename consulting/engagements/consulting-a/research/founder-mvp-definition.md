---
title: "IssueNumber.one Founder MVP Definition"
artifact: research
status: Draft
engagement: "consulting-a"
author: "Codex (Managing Partner)"
model: "GPT-5 (Codex)"
created: "2026-08-04"
updated: "2026-08-04"
confidence: "High on the sponsor's stated product intent and case description; low on customer demand, usability, buyer fit, and outcome improvement because this is founder evidence only"
dependencies:
  - "../decisions/sponsor-discovery-authorization.md"
  - "../experiments/buyer-problem-discovery.md"
related_specs:
  - "../../../../spec/features/issue/README.md"
  - "../../../../spec/features/voting/README.md"
  - "../../../../spec/features/organization/README.md"
promotion_target: "Confirmed behavior is recorded in the owning Feature Specs; founder evidence is not promoted as external validation"
---
# IssueNumber.one Founder MVP Definition

## Owned Question

What problem and minimum product behavior does the sponsor intend
IssueNumber.one to embody, independently of whether the market evidence later
supports building or selling it?

## Concise Answer

IssueNumber.one is intended to prevent consequential unresolved problems from
remaining trapped in meetings and management layers. People may record multiple
issues but nominate only one personal #1 per organization or public topic, and
each nomination receives one free creator star.
Each person also receives one budgeted support vote by default, configurable to
1–7, shared across levels, and continuously movable. Votes are assigned to
issues rather than hierarchy levels, so the same issue keeps the same votes and
score wherever it appears. These signals select one #1 at each organizational scope; that issue automatically
enters the parent scope until a company-level view can show the most important
issue from each department and allow drill-down to its source. Losing rank does
not close an issue or reset its visible age.

## Method and Scope

The Managing Partner conducted a plain-language founder interview in this task
on 2026-08-04 and compared each answer with the Backstage seed, the complete
conceptual Feature tree, and current screen notes. This memo records founder
intent and one retrospective founder case. It is not a customer interview,
demand study, usability test, buyer test, or evidence that the mechanism would
have changed the described outcome.

## Evidence

| ID | Classification | Claim or observation | Source and date | Freshness / limitations |
|---|---|---|---|---|
| E1 | Founder-reported case | Integration tests were disabled to permit an urgent release, were expected to return the next week, and remained disabled for more than three months. | Sponsor interview, 2026-08-04 | Retrospective single case; no independent record inspected |
| E2 | Founder-reported case | New build servers did not support the containers and database setup assumed by the integration tests; developers had raised the problem in stand-ups, retrospectives, and meetings for roughly two years without sufficient escalation or resolution. | Sponsor interview, 2026-08-04 | Retrospective single case; employer and project intentionally not recorded |
| E3 | Founder requirement | A person may keep multiple issues but nominate only one personal #1 per organization or per public topic; only context-specific nominations enter collective voting, and each nomination automatically supplies one creator star without consuming support-vote budget. | Sponsor interview, 2026-08-04 | Product intent, not validation |
| E4 | Founder requirement | Each person receives one budgeted support vote by default, configurable from 1–7 and shared across team, project, department, and company levels; may put multiple budgeted votes on one issue; may spend them on their own nominated issue in addition to its creator star; and may reallocate them at any time. | Sponsor interview, 2026-08-04 | Membership in multiple organizations remains unresolved in the owning Vote Budget Feature |
| E5 | Founder requirement | Anyone in a department may vote; every scope's #1 moves upward automatically without a manager gate. A company leader sees each department's most important issue and can drill down through department, project, team, person, and issue. | Sponsor interview, 2026-08-04 | Exact candidate and budget behavior is specified where confirmed; remaining questions stay with the owning Features |
| E6 | Founder requirement | Members see their own issues, the top N department issues, and the top N company issues. | Sponsor interview, 2026-08-04 | Default N remains unresolved |
| E7 | Founder requirement | Issues may be authored or anonymous. Management activity does not close an issue; it stays open with a visible age, can lose rank without closing, is normally closed by its creator, and may be closed by peer vote when the creator is unavailable. | Sponsor interview, 2026-08-04 | Unavailable-author proof and peer threshold remain unresolved |
| E8 | Observed specification conflict | The prior Specs capped active issues, prohibited self-voting, limited one vote per issue, inconsistently described top-N versus #1 bubble-up, and let ordinary members archive active issues. | IssueNumber.one Feature tree inspected 2026-08-04 before correction | Shows source-of-truth drift, not market evidence |
| E9 | Founder requirement | Votes are assigned to issues, not hierarchy levels. When an issue becomes visible or ineligible at another level, its votes do not travel, split, reset, or receive a separate scope score. | Sponsor interview, 2026-08-04 | Product intent, not validation |

## Findings

### Observed or sourced

- The founder's concrete case is a priority-to-action and escalation failure,
  not merely a missing suggestion box (E1–E2).
- The intended scarcity mechanism is one nominated personal #1 per organization
  or public topic with a free creator star plus a separate scarce support budget
  that defaults to one, not a hard cap on recording open issues (E3–E4).
- Promotion is recursive and automatic: a child #1 becomes a parent candidate;
  top N is a scope view, not the number of issues promoted (E5–E6).
- Votes remain attached to an issue. Hierarchy levels change candidate
  eligibility and visibility, while the issue retains one support score (E9).
- Persistence is a product signal. An unresolved issue remaining open for, for
  example, 247 days is intended to be visible rather than administratively
  deferred or archived away (E7).

### Inferred

- The mechanism is better described as a hierarchy of collective priority
  elections than as generic survey, feedback, or backlog software (E3–E7).
- An issue-bound score makes the hierarchy a sequence of candidate filters
  rather than separate elections that reset or duplicate support at every
  level (E5, E9). Whether users understand or value this mechanism remains
  unvalidated.

## Contradictions and Disconfirming Evidence

- Before this interview was reconciled, several conceptual requirements
  contradicted the founder's intent (E8). The owning Specs have been corrected;
  this memo does not override them.
- No evidence yet shows that users will participate honestly, that voting will
  reduce filtering or retaliation, that leaders will act, or that a 1–7 budget
  produces better priorities than existing alternatives.
- The founder case cannot substitute for the accepted independent buyer/problem
  discovery cohort and is excluded from that experiment's numerator.

## Assumptions and Unknowns

The remaining product decisions are recorded only in the relevant Feature
`## Open Questions` sections. The founder has confirmed one free creator star
plus a separate support budget that defaults to one, is configurable from one
to seven, and is shared across all levels within an organization. Votes are
assigned directly to issues, so there is no separate lower-level or parent-level
vote total. The next central rating unknown is whether the product is
support-only or will include negative votes.

## Confidence

Confidence is high that the corrected Feature Specs now reflect the sponsor's
stated intent on the points above. Confidence remains low that this is the
right MVP to build or a product customers will adopt or buy; no independent
participant, usage, payment, or outcome evidence was added.

## Risks, Dependencies and Next Research

- Continue the founder interview only for unresolved MVP rules already routed
  to their owning Feature Specs.
- Keep the accepted external discovery experiment paused or unexecuted rather
  than describing the founder interview as its replacement.
- Test whether users understand nomination, issue-bound voting, automatic
  bubble-up, anonymity, ageing, and creator-controlled closure before any build
  authority is requested.

## Proposed Experiments

- After the founder definition is complete, construct a no-code scenario using
  the integration-test case and ask development and product managers to walk
  through what they would see, vote on, and do. This would test comprehension
  and workflow hypotheses, not demand or willingness to pay.
