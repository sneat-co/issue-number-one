---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Vote Budget

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting/budget?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Every eligible voter has one budgeted support-vote unit by default, configurable from one to seven and shared across all levels of an organization. This budget is additional to the free automatic creator star on the voter's nominated personal #1. Budgeted units are assigned directly to issues and may be concentrated, distributed among issues, and reallocated at any time.

## Problem

Unlimited voting devalues every vote. A one-to-seven-unit budget forces users to make meaningful trade-offs: when all units are allocated, supporting something new requires moving support from something else. Concentrating several units also lets a voter express relative priority without receiving unlimited influence.

## Behavior

### One organization-wide one-to-seven budget

The budget defaults to one support-vote unit and is configurable from one through seven inclusive. Within one organization, every budgeted vote assigned through a team, project, department, or company view counts against the same personal budget. A view does not own the vote. The automatic creator star does not count against the budget.

#### REQ: configured-budget-range

The vote budget MUST be configurable as an integer from 1 through 7 inclusive.

#### REQ: default-budget-one

A new organization MUST give each eligible voter one budgeted support-vote unit by default.

#### REQ: allocations-cannot-exceed-budget

A voter MUST NOT have more budgeted support units assigned to issues in an organization than their shared budget. The automatic creator star MUST be excluded from this total.

#### REQ: budget-shared-across-all-levels

A voter MUST receive one shared budget within an organization. Voting from a different team, project, department, or company-level view within that organization MUST NOT create an additional budget or a view-specific vote.

#### REQ: cross-level-allocation-reduces-same-balance

Every unit assigned to an issue in the organization MUST reduce the same available balance. Withdrawing a unit or closing its issue MUST make that unit available again.

### Concentrated or distributed allocation

A voter may put several budgeted units on one candidate, including their own nominated issue, or distribute units among several candidates. Their own nominated issue already has one separate automatic creator star.

#### REQ: budget-may-be-concentrated

A voter MUST be allowed to allocate any number of their available units to one eligible issue, up to their remaining budget.

#### REQ: own-issue-uses-same-budget

Any budgeted units allocated to the voter's own nominated issue MUST count against the same budget as units allocated to other candidates. Its one automatic creator star MUST NOT count against that budget.

### Issue recording is separate from voting budget

Raising or keeping an issue does not consume support-vote units. Nominating one personal #1 within the organization automatically contributes one creator star, which is separate from the shared support-vote budget.

#### REQ: raising-does-not-consume-votes

Raising, recording, or keeping an issue open MUST NOT consume vote-budget units.

#### REQ: nomination-star-free

Nominating the person's one personal #1 within an organization MUST create one automatic creator star and MUST leave the creator's full configured support-vote budget available.

### Trading a vote for a new issue

If a voter has allocated their full budget and wants to support another issue, they must remove one or more existing allocations first.

#### REQ: withdraw-to-vote-again

A user at their vote cap MUST withdraw one or more existing vote units before allocating those units to another issue.

### Scope eligibility

Any member of an organizational scope, including members in descendant teams, may use their shared organization-wide budget to support an issue while it is eligible in that scope's ranking. The resulting vote remains assigned to the issue rather than the scope.

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

- Who configures the shared budget, and is the same value mandatory for every member of the organization?
- If a person belongs to multiple organizations, does each organization provide its own independent shared budget?
- Is a minimum membership size required before collective ranking and bubble-up apply?
- Acceptance criteria not yet defined for this feature.
