---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Vote Budget

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Every eligible voter has one configurable budget of one to seven support-vote units shared across all levels of an organization. They may concentrate several units on one eligible issue, distribute them among issues at different levels, and reallocate them at any time.

## Problem

Unlimited voting devalues every vote. A one-to-seven-unit budget forces users to make meaningful trade-offs: when all units are allocated, supporting something new requires moving support from something else. Concentrating several units also lets a voter express relative priority without receiving unlimited influence.

## Behavior

### One organization-wide one-to-seven budget

The configured budget is an integer from one through seven inclusive. Within one organization, every allocation at team, project, department, or company level counts against the same personal budget.

#### REQ: configured-budget-range

The vote budget MUST be configurable as an integer from 1 through 7 inclusive.

#### REQ: allocations-cannot-exceed-budget

A voter MUST NOT have more support units allocated across all organizational levels in total than their shared budget.

#### REQ: budget-shared-across-all-levels

A voter MUST receive one shared budget within an organization. Voting in a different team, project, department, or company-level ranking within that organization MUST NOT create an additional allocation.

#### REQ: cross-level-allocation-reduces-same-balance

Every unit allocated at any organizational level MUST reduce the same available balance, and every withdrawn or refunded unit MUST return to that balance.

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

Any member of an organizational scope, including members in descendant teams, may use their shared organization-wide budget to vote in that scope's ranking.

#### REQ: descendant-members-eligible

A member of a descendant team MUST be eligible to spend their shared organization-wide vote budget in an enclosing department or company ranking.

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
- Who configures the shared budget, and is the same value mandatory for every member of the organization?
- If a person belongs to multiple organizations, does each organization provide its own independent shared budget?
- Is a minimum membership size required before collective ranking and bubble-up apply?
- Acceptance criteria not yet defined for this feature.
