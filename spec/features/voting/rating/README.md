# Feature: Rating

**Status:** Conceptual

## Summary

Rating covers the visible output of voting: upvotes, downvotes, an issue's total score, and the sort orders used in issue lists. Upvotes are public by default; downvotes are anonymous by default; teams may disable downvotes entirely.

## Problem

Raw voting data needs a presentation layer. Users want to understand how an issue is trending, whether their vote is public, and how to sort a list of issues. Rating defines those conventions so every screen behaves consistently and every team can configure the amount of social transparency they want.

## Behavior

### Upvotes and downvotes

An issue has an upvote count, a downvote count, and a total score.

#### REQ: score-formula

An issue's `score` MUST equal `upvotes − downvotes`.

#### REQ: one-vote-per-user-per-issue

A user MUST NOT cast more than one vote on a given issue at a time. Changing their mind means withdrawing and re-casting.

### Users choose up or down

A user spends one unit of their vote budget on either an upvote or a downvote, not both.

#### REQ: one-direction-per-vote

A single vote MUST be either an upvote or a downvote, not both.

### Teams may disable downvotes

Downvotes are optional per team or org. When disabled, users may only cast upvotes.

#### REQ: downvotes-optional

Teams and orgs MUST be able to disable downvotes entirely. When disabled, attempts to cast a downvote MUST be rejected.

### Public vs anonymous vote display

By default, upvotes are publicly attributed and downvotes are anonymous. Teams may change these defaults.

#### REQ: upvotes-public-by-default

Upvotes MUST default to public attribution — other users can see who upvoted.

#### REQ: downvotes-anonymous-by-default

Downvotes MUST default to anonymous — the voter's identity is hidden from other users.

#### REQ: vote-visibility-configurable

Teams and orgs MUST be able to override the default public/anonymous setting for both upvotes and downvotes.

### Sort orders

Issue lists support three sort orders:

| Order | Description |
|-------|-------------|
| `score` | Total score (upvotes − downvotes), highest first |
| `created` | Creation timestamp, newest first |
| `activity` | Last-activity timestamp, most recent first |

#### REQ: supported-sort-orders

Issue list views MUST support sorting by `score`, `created`, and `activity`.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [voting](../README.md) | Rating is the visible layer over the voting mechanism |
| [issue/visibility](../../issue/visibility/README.md) | Score drives bubble-up |
| [permissions](../../permissions/README.md) | Controls who may see vote attribution |

## Dependencies

- voting
- issue
- permissions

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Should there be an additional sort mode like `controversial` (high up AND high down)?
- Should vote-count totals be visible even when individual attribution is anonymous?
- Should users be allowed to flip their vote direction directly, or must they withdraw and re-cast?
- Acceptance criteria not yet defined for this feature.
