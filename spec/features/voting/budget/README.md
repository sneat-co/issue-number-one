---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Vote Budget

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Every eligible voter has a configurable budget of one to seven support-vote units. They may concentrate several units on one eligible issue or distribute them, and may reallocate them at any time. The still-open hierarchy question is whether this budget is separate at each scope or shared across scopes.

## Problem

Unlimited voting devalues every vote. A one-to-seven-unit budget forces users to make meaningful trade-offs: when all units are allocated, supporting something new requires moving support from something else. Concentrating several units also lets a voter express relative priority without receiving unlimited influence.

## Behavior

### Configurable one-to-seven budget

The configured budget is an integer from one through seven inclusive. Every allocation in a voting context counts against the applicable budget.

#### REQ: configured-budget-range

The vote budget MUST be configurable as an integer from 1 through 7 inclusive.

#### REQ: allocations-cannot-exceed-budget

A voter MUST NOT have more allocated support units than the budget applicable to that voting context.

### Concentrated or distributed allocation

A voter may put several units on one candidate, including their own nominated issue, or distribute units among several candidates.

#### REQ: budget-may-be-concentrated

A voter MUST be allowed to allocate any number of their available units to one eligible issue, up to their remaining budget.

#### REQ: own-issue-uses-same-budget

Units allocated to the voter's own nominated issue MUST count against the same budget as units allocated to other candidates.

### Issue recording is separate from voting budget

Raising or keeping an issue does not consume support-vote units. The separate scarcity rule is that a person may nominate only one personal #1 per scope.

#### REQ: raising-does-not-consume-votes

Raising, recording, or keeping an issue open MUST NOT consume vote-budget units.

### Trading a vote for a new issue

If a voter has allocated their full budget and wants to support another issue, they must remove one or more existing allocations first.

#### REQ: withdraw-to-vote-again

A user at their vote cap MUST withdraw one or more existing vote units before allocating those units to another issue.

### Scope eligibility

Any member of an organizational scope, including members in descendant teams, may use the applicable budget to vote in that scope's ranking.

#### REQ: descendant-members-eligible

A member of a descendant team MUST be eligible for the vote budget used in an enclosing department or company ranking.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [voting](../README.md) | Budget is the mechanism voting uses |
| [issue](../../issue/README.md) | One personal #1 nomination is separate from vote-budget scarcity |
| [organization](../../organization/README.md) | Defines team membership |

## Dependencies

- voting
- issue
- organization

## Acceptance Criteria

Not defined yet.

## Open Questions

- What is the default budget within the confirmed 1–7 range?
- Does each hierarchy level provide a separate budget, or does one budget span all scopes in which a person can vote?
- Who configures the budget, and may different scopes choose different values?
- Is a minimum membership size required before collective ranking and bubble-up apply?
- Acceptance criteria not yet defined for this feature.
