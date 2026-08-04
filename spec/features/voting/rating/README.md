---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Rating

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Rating covers the visible output of voting: the free creator star on a personal #1, budgeted support allocated in each organizational scope, vote attribution, ranking, and issue-list sort orders. The founder-defined MVP uses positive support; whether negative votes exist later remains open.

## Problem

Raw voting data needs a presentation layer. Users need to distinguish the creator's automatic nomination signal from scarce peer or self-support, understand support at each scope, know whether attribution is public, and see how candidates are ranked. The same issue may be the #1 of a team while competing at department and company levels.

## Behavior

### Creator star and budgeted support

Every current personal #1 has one automatic creator star. Budgeted support allocations are recorded separately with the organizational scope in which each unit was allocated.

#### REQ: one-creator-star-for-personal-top

A current personal #1 MUST have exactly one automatic creator star. A non-nominated issue MUST NOT have a creator star.

#### REQ: creator-star-free

The creator star MUST NOT reduce the creator's configured support-vote budget.

#### REQ: score-formula

An issue's local support score in a voting scope MUST equal its creator star plus the budgeted support-vote units allocated to it in that scope.

#### REQ: allocation-scope-recorded

Every budgeted support unit MUST retain the team, project, department, or company voting scope in which it was allocated so lower-level and parent-level support can be distinguished.

#### REQ: multiple-units-from-one-voter-count

A voter's multiple budgeted units on one issue MUST all count toward that issue's local support score.

### Support-only ranking

The current ranking mechanism counts positive support units. Negative voting is not required by the founder-defined MVP.

#### REQ: support-score-nonnegative

A scope's local support score MUST be a non-negative integer.

### Public vs anonymous vote display

By default, support allocations are publicly attributed. Teams may hide attribution without changing the total score.

#### REQ: support-public-by-default

Support allocations MUST default to public attribution so eligible viewers can see who allocated support.

#### REQ: vote-visibility-configurable

Teams and orgs MUST be able to hide individual support attribution while retaining aggregate scores.

### Sort orders

Issue lists support three sort orders:

| Order | Description |
|-------|-------------|
| `score` | Applicable ranking score, highest first |
| `created` | Creation timestamp, newest first |
| `activity` | Last-activity timestamp, most recent first |

#### REQ: supported-sort-orders

Issue list views MUST support sorting by `score`, `created`, and `activity`.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [voting](../README.md) | Rating is the visible layer over the voting mechanism |
| [issue/visibility](../../issue/visibility/README.md) | The applicable ranking score determines each #1 and automatic bubble-up |
| [permissions](../../permissions/README.md) | Controls who may see vote attribution |

## Dependencies

- voting
- issue
- permissions

## Acceptance Criteria

Not defined yet.

## Open Questions

- Should aggregate support totals remain visible when individual attribution is hidden?
- Are negative/downvotes intentionally excluded from the product, or only from the MVP?
- When an issue becomes a parent-scope candidate, does its ranking score carry budgeted support from lower scopes or use only the creator star plus support allocated at the parent scope?
- How are equal support scores ordered when determining a single #1?
- Acceptance criteria not yet defined for this feature.
