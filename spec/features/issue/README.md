# Feature: Issue

**Status:** Conceptual

## Summary

An issue is the atomic unit of IssueNumber.one — a raised priority item that a team, org, or public topic must identify and address. Unlike a traditional issue tracker, IssueNumber.one caps how many issues any one person can have active, forcing surfaced priorities to reflect what actually matters today.

## Contents

| Directory | Description |
|-----------|-------------|
| [lifecycle/](lifecycle/README.md) | Issue statuses and the transitions between them |
| [anonymity/](anonymity/README.md) | Anonymous versus authored issues and author-hiding rules |
| [visibility/](visibility/README.md) | Team, org, and public visibility levels and bubble-up behavior |
| [moderation/](moderation/README.md) | Archiving, banning, and who is allowed to moderate |

### lifecycle

Defines the statuses an issue can be in (`raised`, `withdrawn`, `resolved`, `banned/moderated`, `archived`) and who can transition an issue between them. Only the creator can withdraw their own issue; team members can archive team issues; public topic issues cannot be archived.

### anonymity

Teams may allow anonymous issues (and configure a small fee to prevent artificial boosting) while still letting the originator update or withdraw their anonymous issue. Some teams allow a mix of personal and anonymous issues per member.

### visibility

Every issue has a visibility level: `team/space`, `org`, or `public`. Top-ranked team issues bubble up to the org level when sufficiently upvoted, and only the team's current #1 issue is visible externally to other teams.

### moderation

Any team member can archive any team issue. Any org member can archive any team. Public topic issues cannot be archived — only withdrawn by their creator. Supporters and creators are notified when an issue is archived.

## Problem

Traditional issue trackers encourage backlog growth: everything accumulates and nothing is truly prioritized. Teams end up with hundreds of issues that are all nominally important, managers can't reliably answer "what is the one thing we should do next?", and critical work disappears under daily tasks. IssueNumber.one treats an issue as a scarce resource — by limiting how many each person can have active at once, the list that surfaces is the list that matters.

## Behavior

### Per-person issue budget

By default each team member has exactly one active issue per team. A team may relax this cap, but never above three. This is not an issue tracker.

#### REQ: one-active-issue-default

By default, every team member MUST be limited to exactly one active issue per team.

#### REQ: max-three-per-person-cap

A team MAY allow its members to raise more than one active issue, but MUST NOT allow more than three active issues per member per team.

### Issue fields

Every issue has the following fields:

| Field | Description |
|-------|-------------|
| `id` | Unique identifier |
| `title` | Short summary of the issue |
| `description` | Long-form description |
| `author` | Creator, or marker that the issue is anonymous (see [anonymity](anonymity/README.md)) |
| `status` | Current state (see [lifecycle](lifecycle/README.md)) |
| `visibility` | `team` / `org` / `public` (see [visibility](visibility/README.md)) |
| `assignee` | Person responsible for addressing the issue (optional until priority #1) |
| `deadline` | Target date for resolution (optional until priority #1) |
| `progress` | Progress indicator, typically shown as a bar on the team page |
| `scope` | The team, org, or topic this issue belongs to |
| `createdAt` | Creation timestamp |
| `updatedAt` | Last update timestamp |

#### REQ: issue-required-fields

Every issue MUST have `id`, `title`, `status`, `visibility`, and `scope`. All other fields are optional.

### Raising an issue

Anyone with the required permission in a scope can raise an issue in that scope. If the user is already at their active-issue cap, they MUST first withdraw an existing issue before raising a new one.

#### REQ: raise-requires-capacity

A user MUST NOT be able to raise a new issue if doing so would exceed their per-team active-issue cap. They MUST withdraw an existing issue first.

#### REQ: raise-requires-permission

A user MUST have permission in the target scope before raising an issue there (see [permissions](../permissions/README.md)).

### Withdrawing an issue

Only the issue's creator can withdraw their own issue. Supporters of a withdrawn issue are notified and their votes are refunded (see [voting](../voting/README.md)).

#### REQ: creator-only-withdraw

Only the issue's creator MAY withdraw an issue. No other team member, including team admins, may withdraw on the creator's behalf.

#### REQ: notify-supporters-on-withdraw

When an issue is withdrawn, all supporters MUST be notified.

### The team's #1 issue

Each team always has at most one currently-active #1 issue (the one with the highest score). This issue is the only one visible externally to other teams and is the team's current focus.

#### REQ: single-top-issue-per-team

At any given moment, a team MUST have at most one "#1 issue" — the issue with the highest score. Ties MUST be resolved deterministically (e.g., by creation time).

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [voting](../voting/README.md) | Votes determine an issue's score and therefore its rank and bubble-up eligibility |
| [organization](../organization/README.md) | Every issue is scoped to a team, org, or topic |
| [permissions](../permissions/README.md) | Permissions gate who can raise/see/archive each issue |
| [storage](../storage/README.md) | Issues are persisted in either cloud or git-backed storage |
| [ai-integration](../ai-integration/README.md) | AI executive summaries analyze the current set of issues |

## Dependencies

- organization
- voting
- permissions

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Should `assignee` and `deadline` be mandatory once an issue becomes the team's #1, or remain optional?
- Should `progress` be a free-form field, a 0–100 integer, or discrete milestones?
- How are ties in score broken when selecting the team's #1 issue?
- Acceptance criteria not yet defined for this feature.
