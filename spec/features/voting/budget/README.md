# Feature: Vote Budget

**Status:** Conceptual

## Summary

Every team member has a small per-team vote budget that they spend on issues they want to prioritize. The budget forces prioritization — if you want to back a new issue, you must withdraw one of your existing votes first.

## Problem

Unlimited voting devalues every vote. A budget forces users to make meaningful trade-offs: when you support something new, you must let go of something old. That trade-off is the whole point of IssueNumber.one — it surfaces what people genuinely think matters most *today*.

## Behavior

### Budget is per team

Budgets are not global. Each team a user belongs to gives them a separate allocation.

#### REQ: budget-is-per-team

A user's vote budget MUST be scoped to a single team and MUST NOT be shared across teams.

### Raising an issue consumes a slot too

Within the scope of the [issue](../../issue/README.md) active-issue cap, a user's own raised issue implicitly counts against the same prioritization budget. In other words: raising an issue OR voting on someone else's issue are both ways to spend your priority attention.

#### REQ: raise-counts-against-priority

A user's own active issue MUST count against their priority budget, so that raising more issues leaves fewer votes to spend.

### Trading a vote for a new issue

If a user is at their active-issue cap and wants to raise a new issue, they MUST first withdraw an existing issue (per [issue#req-raise-requires-capacity](../../issue/README.md)). Similarly, if a user is at their vote cap and wants to vote on a new issue, they MUST withdraw an existing vote.

#### REQ: withdraw-to-vote-again

A user at their vote cap MUST withdraw one of their existing votes before casting a new vote.

#### REQ: withdraw-vote-to-raise

If a user is at their combined priority cap, they MAY reclaim a slot by withdrawing a vote they have cast elsewhere.

### Minimum team size

A team must be large enough for budgeting to make sense. Teams of one or two create prioritization theater.

#### REQ: minimum-team-size

A team MUST have at least three members for voting to be meaningful. Teams below this size MAY behave as personal scratch space but MUST NOT participate in bubble-up or cross-team visibility.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [voting](../README.md) | Budget is the mechanism voting uses |
| [issue](../../issue/README.md) | Active-issue caps and budget caps interact |
| [organization](../../organization/README.md) | Defines team membership |

## Dependencies

- voting
- issue
- organization

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- What are the default budget values (total votes per member)?
- Does an org admin set the budget, or is it fixed per product version?
- Is the priority budget literally the sum of (issues raised + votes cast), or are they counted separately?
- Acceptance criteria not yet defined for this feature.
