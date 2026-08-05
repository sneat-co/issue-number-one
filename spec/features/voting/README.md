---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Voting

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Voting is the positive-support mechanism by which eligible members decide which issue is #1 at each scope. The MVP has no downvotes or opposing votes. A personal #1 automatically receives one free creator star. In addition, every person has one support-vote unit by default, configurable from one to seven. In the MVP, every eligible member of an organization receives the same configured budget, shared across all levels of that organization. Every budgeted vote is assigned directly to an issue, not to a team, department, company, or other hierarchy level. Members may assign several budgeted units to one issue, support their own nominated issue, and reassign their votes at any time.

## Contents

| Directory | Description |
|-----------|-------------|
| [budget/](budget/README.md) | One-to-seven vote budgets, multi-unit allocation, and reallocation rules |
| [rating/](rating/README.md) | One issue-bound support score, sort orders, and public/anonymous display |

### budget

Budgeted votes are a limited resource. Members assign a small number of units directly to eligible issues and must reassign existing units when their priority changes. Raising or recording an issue does not consume this vote budget. Nominating one personal #1 automatically supplies one separate creator star without consuming the budget.

### rating

An issue's support display distinguishes the automatic creator star from budgeted support-vote units assigned to the issue. The issue has one score wherever it appears in the hierarchy. Individual vote attribution may be public or hidden, but the total star count remains visible. Issue lists sort by issue score, creation date, or last activity date.

## Problem

Traditional voting systems give every user unlimited votes, which turns voting into a popularity contest rather than a prioritization tool. In IssueNumber.one, a personal nomination contributes one free creator star, while additional support is limited on purpose: the default budget is one unit, configurable up to seven, and may be concentrated or moved as priorities change. This forces a hierarchy toward shared focus without preventing people from recording other unresolved issues.

## Behavior

### Eligible candidates at each scope

At a person's immediate team scope, voting candidates are members' nominated personal #1 issues. At a parent scope, voting candidates are the current #1 issues of its direct child scopes.

#### REQ: leaf-candidates-are-personal-tops

Only members' nominated personal #1 issues MUST be eligible in the immediate team voting context.

#### REQ: parent-candidates-are-child-tops

Only each direct child scope's current #1 issue MUST be eligible in a parent voting context through automatic bubble-up.

#### REQ: public-topic-candidates-are-topic-personal-tops

Within a public topic, only each participant's personal #1 nomination for that topic MUST be eligible for topic voting. A nomination in one organization or public topic MUST NOT automatically enter another public topic's candidate set.

### Budget, concentration, and self-support

Each eligible voter has one budgeted support unit by default, configurable between one and seven. They may distribute the units or place several, including all available units, on one eligible issue. They may allocate budgeted support to an issue they authored in addition to its free creator star.

#### REQ: vote-budget-one-to-seven

The configured vote budget for an eligible voter MUST be an integer from 1 through 7 inclusive.

#### REQ: default-vote-budget-one

A new organization MUST default each eligible voter to one budgeted support-vote unit.

#### REQ: mvp-same-budget-for-all-org-members

In the MVP, an organization MUST configure one vote-budget size that applies equally to every eligible member. The MVP MUST NOT assign different budget sizes by person, role, membership type, team, or hierarchy level.

#### REQ: one-budget-shared-across-org-levels

A voter MUST have one vote budget shared across every team, project, department, and company ranking within the same organization. Allocating a unit at any level MUST reduce the same remaining balance; entering another level MUST NOT create or reset a separate budget.

#### REQ: multiple-units-per-issue

A voter MUST be allowed to allocate more than one available vote unit to the same eligible issue.

#### REQ: self-support-allowed

A voter MUST be allowed to allocate available budgeted vote units to their own nominated issue in addition to its automatic creator star.

#### REQ: creator-star-is-outside-budget

The automatic creator star on a personal #1 MUST NOT reduce the creator's available support-vote budget.

### Positive support only

The MVP lets members express support by assigning positive votes. It does not let a voter oppose an issue or reduce its score with a negative vote.

#### REQ: mvp-voting-is-positive-only

The MVP MUST NOT provide downvotes, negative votes, opposition votes, or any other vote that subtracts from an issue's support score.

### Vote refunds on closure

When an issue exits the `raised` state, every supporter gets all units allocated to it back.

#### REQ: votes-refunded-on-closure

When an issue leaves the `raised` state, the system MUST refund every vote unit allocated to it to its voter.

### Withdrawing a vote to vote elsewhere

A voter may change priorities at any time by removing one or more units from an issue and allocating them elsewhere.

#### REQ: withdraw-vote-to-vote-again

Users MUST be able to withdraw any number of units they allocated. Withdrawn units MUST immediately return to the user's available budget.

### Voting in a department or other organizational scope

Every member of a department is eligible to vote in that department's ranking, including members whose direct membership is in a nested project or team. The same rule applies to any organizational scope and its descendant members.

#### REQ: all-scope-members-may-vote

Every member of an organizational scope, including members of its descendant scopes, MUST be allowed to vote on the candidates eligible in that scope.

### Votes remain assigned to issues

Hierarchy changes affect an issue's candidate eligibility and visibility, not its votes. When an issue becomes visible at a higher level, its votes are not copied, moved, reset, or re-cast. Newly eligible members may add their own votes to that same issue. If the issue later loses candidate eligibility, its existing votes remain assigned and continue consuming their voters' budgets until those voters reassign them or the issue leaves the `raised` state.

#### REQ: budgeted-vote-belongs-to-issue

Every budgeted support-vote unit MUST reference one issue and MUST NOT belong to an organizational hierarchy level or ranking view.

#### REQ: hierarchy-change-does-not-move-votes

An issue becoming or ceasing to be a candidate at any hierarchy level MUST NOT copy, move, reset, refund, or re-scope any vote assigned to it.

#### REQ: ineligible-issue-keeps-assigned-votes

If an issue ceases to be eligible in a ranking view, its existing budgeted votes MUST remain assigned and MUST continue consuming their voters' budgets until individually withdrawn or released when the issue leaves the `raised` state.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../issue/README.md) | Votes drive issue scores and ranking |
| [issue/lifecycle](../issue/lifecycle/README.md) | Exiting `raised` triggers refunds |
| [issue/visibility](../issue/visibility/README.md) | An issue's single score determines its order within each eligible candidate set and therefore automatic bubble-up |
| [permissions](../permissions/README.md) | Who may vote in which scopes |

## Dependencies

- issue
- permissions

## Acceptance Criteria

Not defined yet.

## Open Questions

- After the MVP, should organizations be allowed to vary vote budgets by role or membership type?
- Acceptance criteria not yet defined for this feature.
