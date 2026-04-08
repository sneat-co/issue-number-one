# Feature: Issue Visibility

**Status:** Conceptual

## Summary

Every issue has a visibility level that controls who can see it. Visibility can escalate as an issue gains support — a sufficiently upvoted team issue bubbles up to the org level — and only the team's current top-priority issue is visible externally to other teams in the same org.

## Problem

Without explicit visibility rules, teams either over-share (every draft issue cluttering the org view) or under-share (critical cross-team concerns stay buried). Visibility in IssueNumber.one is tiered so that drafts stay local while top-priority items surface at the right scope, and so that public topics are discoverable beyond any single org.

## Behavior

### Visibility levels

| Level | Who can see it |
|-------|----------------|
| `team` (default) | Only members of the issue's team or space |
| `org` | Any member of the enclosing organization |
| `public` | Anyone, used for public topics and spaces |

#### REQ: default-visibility-team

A newly-raised issue in a team scope MUST default to `team` visibility.

#### REQ: visibility-levels

An issue's `visibility` MUST be one of `team`, `org`, or `public`.

### Bubble-up by upvotes

Team issues automatically escalate to `org` visibility once they cross a configurable upvote threshold. Only the top N issues per team bubble up.

#### REQ: bubble-up-by-upvotes

A `team` issue MUST be promoted to `org` visibility when it accumulates enough upvotes to be among the top-N for its team.

#### REQ: top-n-configurable

The value of `N` (how many issues bubble up per team) MUST be configurable per org.

### External top-issue rule

Only the team's current #1 issue is visible to other teams in the same org. All other team issues remain scoped to team members, regardless of upvote count.

#### REQ: external-shows-top-only

To users outside the team but inside the org, only the team's current #1 issue (per [issue/README.md#req-single-top-issue-per-team](../README.md)) MUST be visible.

### Public topics

Public topics are always visible to anyone, logged in or not. Any nested sub-topics inherit public visibility.

#### REQ: public-topic-always-public

Issues in public topics MUST always be `public` visibility; they MUST NOT support `team` or `org` scoping.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../README.md) | The #1-issue rule depends on [issue#req-single-top-issue-per-team](../README.md) |
| [voting](../../voting/README.md) | Upvote counts drive bubble-up |
| [organization](../../organization/README.md) | Team/org/topic structure defines the scopes |
| [permissions](../../permissions/README.md) | Visibility is enforced via permissions |

## Dependencies

- issue
- voting
- organization
- permissions

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- What is the default value of `N` for bubble-up, and how is it configured (per-team or per-org only)?
- When an org-level issue drops out of the top-N, does it revert to `team` visibility or stay at `org`?
- Should visibility be downgradable by the author, or only monotonically escalatable?
- Acceptance criteria not yet defined for this feature.
