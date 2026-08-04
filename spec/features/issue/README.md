---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Issue

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

An issue is the atomic unit of IssueNumber.one — a raised priority item that a team, company, or public topic must identify and address. A person may keep multiple open issues, but may nominate at most one personal #1 issue in a scope at a time. Only that nominated issue enters collective voting, so scarcity applies to attention rather than to recording problems.

## Contents

| Directory | Description |
|-----------|-------------|
| [lifecycle/](lifecycle/README.md) | Issue statuses and the transitions between them |
| [anonymity/](anonymity/README.md) | Anonymous versus authored issues and author-hiding rules |
| [visibility/](visibility/README.md) | Team, org, and public visibility levels and bubble-up behavior |
| [moderation/](moderation/README.md) | Archiving, banning, and who is allowed to moderate |

### lifecycle

Defines the statuses an issue can be in (`raised`, `withdrawn`, `resolved`, `banned/moderated`, `archived`) and who can transition an issue between them. An issue stays open and visibly ageing when it loses votes or rank. Only the creator normally withdraws or resolves it; eligible peers may close it by vote only when the creator is unavailable.

### anonymity

Teams may allow authored and anonymous issues while still letting the originator control an anonymous issue. Authored and anonymous issues share the same one-personal-#1 nomination rule.

### visibility

Every issue has a visibility level: `team/space`, `org`, or `public`. Each scope shows its top N issues to eligible members, while only its current #1 issue automatically becomes a candidate in the parent scope. No manager approval is required.

### moderation

Moderation handles abusive or policy-violating content without providing a way to suppress an inconvenient issue. A manager or ordinary member cannot archive an active issue merely to close it. A ban requires an authorized moderator, a visible reason, and notifications to the creator and supporters.

## Problem

Traditional issue trackers encourage backlog growth: everything accumulates and nothing is truly prioritized. Teams end up with hundreds of issues that are all nominally important, managers cannot reliably answer "what is the one thing we should do next?", and critical work disappears under daily tasks. IssueNumber.one lets people record the issues they know about while forcing each person to identify only one current #1 for collective prioritization.

## Behavior

### Personal issue list and #1 nomination

A person may keep multiple open issues in a scope. They choose at most one as their personal #1. Only that nominated issue is eligible for collective voting in the immediate team or scope.

#### REQ: multiple-open-issues-per-person

A member MUST be able to keep multiple `raised` issues in the same scope.

#### REQ: one-personal-top-nomination

A member MUST have at most one nominated personal #1 issue per scope at a time.

#### REQ: only-personal-top-enters-voting

Only a member's nominated personal #1 issue MUST be eligible for collective voting in that member's immediate scope.

#### REQ: changing-nomination-does-not-close

When a member changes their personal #1 nomination, the previously nominated issue MUST remain `raised` unless the creator separately withdraws or resolves it.

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
| `nominatedAt` | When the creator most recently made this issue their personal #1, or empty when it is not nominated |

#### REQ: issue-required-fields

Every issue MUST have `id`, `title`, `status`, `visibility`, and `scope`. All other fields are optional.

### Raising an issue

Anyone with the required permission in a scope can raise an issue in that scope. Raising it does not automatically make it the creator's personal #1.

#### REQ: raising-does-not-replace-nomination

Raising a new issue MUST NOT replace the creator's existing personal #1 nomination without an explicit nomination action by the creator.

#### REQ: raise-requires-permission

A user MUST have permission in the target scope before raising an issue there (see [permissions](../permissions/README.md)).

### Withdrawing an issue

Only the issue's creator can withdraw their own issue. Supporters of a withdrawn issue are notified and their votes are refunded (see [voting](../voting/README.md)).

#### REQ: creator-only-withdraw

Only the issue's creator MAY withdraw an issue. No other team member, including team admins, may withdraw on the creator's behalf.

#### REQ: notify-supporters-on-withdraw

When an issue is withdrawn, all supporters MUST be notified.

### Resolving an issue

Management acknowledgement, assignment, explanation, or claimed action does not close an issue. The creator confirms when the problem has been addressed. If the creator is unavailable, eligible peers may close it through the governed closure vote defined in [lifecycle](lifecycle/README.md).

#### REQ: creator-confirms-resolution

Only the creator MUST normally be allowed to mark a `raised` issue as `resolved`.

#### REQ: management-action-does-not-close

Acknowledging, assigning, commenting on, or recording work against an issue MUST NOT change its lifecycle status or reset its open age.

### The team's #1 issue

Each team or organizational scope always has at most one current #1 issue: the eligible candidate with the highest support score in that scope. This #1 automatically becomes eligible in the parent scope while remaining open in its source scope.

#### REQ: single-top-issue-per-team

At any given moment, a scope MUST have at most one "#1 issue" — the eligible issue with the highest scope-specific support score. Ties MUST be resolved deterministically.

### Visible age and persistence

An unresolved issue stays recorded and its age remains visible even when it is no longer nominated, receives no votes, or falls out of a top-N list.

#### REQ: open-age-visible

Every `raised` issue MUST expose how long it has remained open, derived from its original creation timestamp.

#### REQ: loss-of-rank-does-not-close

Removing votes from an issue, changing a nomination, or losing rank MUST NOT withdraw, resolve, archive, or reset the age of that issue.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [voting](../voting/README.md) | Votes determine an issue's score and therefore its rank and bubble-up eligibility |
| [organization](../organization/README.md) | Every issue is scoped to a team, org, or topic |
| [permissions](../permissions/README.md) | Permissions gate who can raise, see, vote, close, and moderate each issue |
| [storage](../storage/README.md) | Issues are persisted in either cloud or git-backed storage |
| [ai-integration](../ai-integration/README.md) | AI executive summaries analyze the current set of issues |

## Dependencies

- organization
- voting
- permissions

## Acceptance Criteria

Not defined yet.

## Open Questions

- Should `assignee` and `deadline` be mandatory once an issue becomes the team's #1, or remain optional?
- Should `progress` be a free-form field, a 0–100 integer, or discrete milestones?
- How are ties in score broken when selecting the team's #1 issue?
- Should a person have a configurable maximum number of recorded open issues, or no product-enforced maximum?
- Acceptance criteria not yet defined for this feature.
