---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Voting

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/voting?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Voting is the mechanism by which eligible members decide which issue is #1 at each scope. Every person has a configurable budget of one to seven support-vote units, may allocate several units to the same issue, may support their own nominated issue, and may move allocations at any time.

## Contents

| Directory | Description |
|-----------|-------------|
| [budget/](budget/README.md) | One-to-seven vote budgets, multi-unit allocation, and reallocation rules |
| [rating/](rating/README.md) | Scope-specific support scores, sort orders, and public/anonymous display |

### budget

Votes are a limited resource. Members allocate a small number of units among the candidates eligible in a scope and must move existing units when their priority changes. Raising or recording an issue does not consume this vote budget; nominating one personal #1 is the separate scarcity mechanism.

### rating

An issue's support score in a scope is the total number of support-vote units allocated to it in that scope. Vote attribution may be public or hidden. Issue lists sort by support score, creation date, or last activity date.

## Problem

Traditional voting systems give every user unlimited votes, which turns voting into a popularity contest rather than a prioritization tool. In IssueNumber.one, votes are limited on purpose: a member has only one to seven units and may concentrate or move them as priorities change. Combined with one personal #1 nomination, this forces a hierarchy toward shared focus without preventing people from recording other unresolved issues.

## Behavior

### Eligible candidates at each scope

At a person's immediate team scope, voting candidates are members' nominated personal #1 issues. At a parent scope, voting candidates are the current #1 issues of its direct child scopes.

#### REQ: leaf-candidates-are-personal-tops

Only members' nominated personal #1 issues MUST be eligible in the immediate team voting context.

#### REQ: parent-candidates-are-child-tops

Only each direct child scope's current #1 issue MUST be eligible in a parent voting context through automatic bubble-up.

### Budget, concentration, and self-support

Each eligible voter has a configurable total budget between one and seven units. They may distribute the units or place several, including all available units, on one eligible issue. They may vote for an issue they authored.

#### REQ: vote-budget-one-to-seven

The configured vote budget for an eligible voter MUST be an integer from 1 through 7 inclusive.

#### REQ: multiple-units-per-issue

A voter MUST be allowed to allocate more than one available vote unit to the same eligible issue.

#### REQ: self-support-allowed

A voter MUST be allowed to allocate available vote units to their own nominated issue.

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

### Candidate replacement refunds scope allocations

When an issue ceases to be eligible in a voting scope because a person changes their nomination or a child scope gets a different #1, allocations made in that scope return to their voters. The issue remains open at its source.

#### REQ: refund-when-candidate-leaves-scope

When an issue stops being eligible in a voting scope, every vote unit allocated to it in that scope MUST immediately return to its voter without changing the issue's lifecycle status.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../issue/README.md) | Votes drive issue scores and ranking |
| [issue/lifecycle](../issue/lifecycle/README.md) | Exiting `raised` triggers refunds |
| [issue/visibility](../issue/visibility/README.md) | Scope-specific support rankings determine each #1 and automatic bubble-up |
| [permissions](../permissions/README.md) | Who may vote in which scopes |

## Dependencies

- issue
- permissions

## Acceptance Criteria

Not defined yet.

## Open Questions

- What is the default budget size within the confirmed 1–7 range?
- Does each person receive a separate 1–7 budget at every hierarchy level, or one shared budget across team, department, and company rankings?
- May the configured budget vary by scope or membership type, or must every eligible voter receive the same number?
- Are negative/downvotes part of the MVP, or is voting support-only?
- Acceptance criteria not yet defined for this feature.
