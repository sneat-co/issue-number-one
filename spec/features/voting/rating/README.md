---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Rating

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/rating?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Rating covers the visible output of voting: an issue's support score in each organizational scope, vote attribution, and the sort orders used in issue lists. The founder-defined MVP mechanism uses positive support allocations; whether negative votes exist later remains open.

## Problem

Raw voting data needs a presentation layer. Users want to understand an issue's support at the current scope, whether their allocation is public, and how the candidate list is ranked. Scores must remain scope-specific because the same issue may be the #1 of a team while competing separately at department and company levels.

## Behavior

### Scope-specific support score

An eligible issue has a support score for each voting scope in which it appears.

#### REQ: score-formula

An issue's support score in a scope MUST equal the sum of support-vote units currently allocated to it in that scope.

#### REQ: score-is-scope-specific

A team, department, and company score for the same issue MUST be stored and presented as separate scope-specific values.

#### REQ: multiple-units-from-one-voter-count

A voter's multiple allocated units on one issue MUST all count toward that issue's support score.

### Support-only ranking

The current ranking mechanism counts positive support units. Negative voting is not required by the founder-defined MVP.

#### REQ: support-score-nonnegative

A scope's support score MUST be a non-negative integer.

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
| `score` | Scope-specific support score, highest first |
| `created` | Creation timestamp, newest first |
| `activity` | Last-activity timestamp, most recent first |

#### REQ: supported-sort-orders

Issue list views MUST support sorting by `score`, `created`, and `activity`.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [voting](../README.md) | Rating is the visible layer over the voting mechanism |
| [issue/visibility](../../issue/visibility/README.md) | Scope-specific score determines each #1 and automatic bubble-up |
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
