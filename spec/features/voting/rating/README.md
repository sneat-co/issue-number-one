---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Rating

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Rating covers the visible output of voting: the free creator star on a personal #1, budgeted support assigned directly to issues, vote attribution, ranking, and issue-list sort orders. Each issue has one support score wherever it appears in the hierarchy. The founder-defined MVP uses positive support; whether negative votes exist later remains open.

## Problem

Raw voting data needs a presentation layer. Users need to distinguish the creator's automatic nomination signal from scarce peer or self-support, know whether attribution is public, and see how candidates are ranked. The same issue may appear in team, department, and company views, but its votes and score do not split by view.

## Behavior

### Creator star and budgeted support

Every current personal #1 has one automatic creator star. Every budgeted support allocation is assigned directly to one issue.

#### REQ: one-creator-star-for-personal-top

A current personal #1 MUST have exactly one automatic creator star. A non-nominated issue MUST NOT have a creator star.

#### REQ: creator-star-free

The creator star MUST NOT reduce the creator's configured support-vote budget.

#### REQ: score-formula

An issue's support score MUST equal its creator star, if any, plus all budgeted support-vote units currently assigned to that issue.

#### REQ: one-global-issue-score

The same issue MUST expose the same support score in every team, project, department, company, or other ranking view in which it appears.

#### REQ: ranking-view-does-not-own-vote

The view or hierarchy level from which a voter assigns support MAY be retained as interaction metadata, but MUST NOT own, duplicate, or create a separate score for that vote.

#### REQ: multiple-units-from-one-voter-count

A voter's multiple budgeted units on one issue MUST all count toward that issue's support score.

### Support-only ranking

The current ranking mechanism counts positive support units. Negative voting is not required by the founder-defined MVP.

#### REQ: support-score-nonnegative

An issue's support score MUST be a non-negative integer.

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
| `score` | Issue support score, highest first |
| `created` | Creation timestamp, newest first |
| `activity` | Last-activity timestamp, most recent first |

#### REQ: supported-sort-orders

Issue list views MUST support sorting by `score`, `created`, and `activity`.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [voting](../README.md) | Rating is the visible layer over the voting mechanism |
| [issue/visibility](../../issue/visibility/README.md) | The issue score determines order within each eligible candidate set and therefore each #1 |
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
- How are equal support scores ordered when determining a single #1?
- Acceptance criteria not yet defined for this feature.
