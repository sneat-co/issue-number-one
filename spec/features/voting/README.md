# Feature: Voting

**Status:** Conceptual

## Summary

Voting is the mechanism by which team members signal which issues matter most. Votes are deliberately scarce: every member has a small per-team budget that must be spent carefully. Voting on your own issue is never allowed, votes are refunded when an issue closes, and teams may optionally enable downvotes.

## Contents

| Directory | Description |
|-----------|-------------|
| [budget/](budget/README.md) | Per-team vote budgets, cross-team support votes, and withdraw-to-vote-again rules |
| [rating/](rating/README.md) | Upvote and downvote mechanics, sort orders, and public/anonymous display |

### budget

Votes are a limited per-team resource. Members begin with a small budget, can spend extra votes to support issues instead of raising new ones, and can withdraw their vote to free it for something else. This forces prioritization of what actually matters today.

### rating

An issue's score is the net of upvotes and downvotes. Teams may disable downvotes entirely. By default upvotes are public and downvotes are anonymous. Issue lists sort by total score, creation date, or last activity date.

## Problem

Traditional voting systems give every user unlimited votes, which turns voting into a popularity contest rather than a prioritization tool. In IssueNumber.one, votes are limited on purpose — if you only have a handful, you spend them on what truly matters. Combined with the issue budget (one active issue per person by default), this forces the whole team toward shared focus.

## Behavior

### Scarcity is the point

Every user has a small per-team vote budget. The budget is deliberately smaller than the number of live issues a team might have.

#### REQ: per-team-vote-budget

Every team member MUST have a per-team vote budget. Votes in one team MUST NOT affect budgets in any other team.

### No self-voting

Users cannot spend votes on their own issues, whether authored or anonymous.

#### REQ: no-self-voting

A user MUST NOT be able to vote on an issue they authored (or raised anonymously — see [issue/anonymity#req-anon-no-self-voting](../issue/anonymity/README.md)).

### Vote refunds on closure

When an issue exits the `raised` state for any reason — withdrawal, resolution, archiving, or banning — every supporter gets their vote back.

#### REQ: votes-refunded-on-closure

When an issue leaves the `raised` state, the system MUST refund every vote that was cast on it to its voter.

### Withdrawing a vote to vote elsewhere

A user who has reached their vote limit and wants to support a new issue MUST be able to withdraw one of their existing votes to free up budget.

#### REQ: withdraw-vote-to-vote-again

Users MUST be able to withdraw any vote they have cast. A withdrawn vote MUST immediately return to the user's budget.

### Supporting across teams

Users MAY be allowed to spend a vote on an issue in a team other than their own when the scope configuration allows it. This is the "support vote" surfaced on the landing page.

#### REQ: cross-team-support-optional

Cross-team support voting MUST be opt-in per team. When enabled, the vote still counts against the voter's own per-team budget for their home team.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../issue/README.md) | Votes drive issue scores and ranking |
| [issue/lifecycle](../issue/lifecycle/README.md) | Exiting `raised` triggers refunds |
| [issue/visibility](../issue/visibility/README.md) | Upvote totals drive bubble-up |
| [permissions](../permissions/README.md) | Who may vote in which scopes |

## Dependencies

- issue
- permissions

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- What is the default budget size (number of votes per member per team)?
- Should the budget scale with team size?
- Should a cross-team support vote consume a slot in the home team's budget, the visited team's budget, or a separate "support" pool?
- Acceptance criteria not yet defined for this feature.
