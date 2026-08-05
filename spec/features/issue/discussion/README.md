---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Issue discussion

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/discussion?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/discussion?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/discussion?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/discussion?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

Lets people discuss an issue in its existing trust domain without turning discussion activity into votes, closure, or a reset of how long the issue has remained open.

## Problem

An issue's title, description, and support score can show that a problem matters without explaining disagreements, evidence, attempted fixes, or why the creator still considers it unresolved. Those details otherwise fragment across meetings and chat, recreating the same loss of context that IssueNumber.one is intended to prevent. Discussion needs to remain attached to the issue while staying separate from the scarce-vote priority signal and creator-controlled closure.

## Behavior

### Discussion belongs to one issue

Every discussion entry is activity attached to exactly one issue. It does not become an independent issue, receive votes, or enter the hierarchy.

#### REQ: discussion-attached-to-one-issue

Every top-level contribution and response MUST belong to exactly one issue and MUST NOT be ranked or promoted independently of that issue.

### Discussion does not alter priority or lifecycle

Writing, editing, receiving, or replying to a discussion entry is not a vote, nomination, management acknowledgement, or lifecycle transition.

#### REQ: discussion-does-not-change-priority

Discussion activity MUST NOT add or remove creator stars or budgeted votes, change an issue's support score, alter its nomination, or independently change its rank.

#### REQ: discussion-does-not-close-or-reset-age

Discussion activity MUST NOT withdraw, resolve, ban, or archive an issue and MUST NOT reset or hide its original open age.

### Discussion stays inside the issue trust domain

A discussion can reveal more sensitive context than the issue title. It therefore follows the issue's visibility boundary rather than becoming a separate broadly visible feed.

#### REQ: discussion-follows-issue-visibility

A person who cannot view an issue MUST NOT be able to view its discussion. Making an issue visible in a parent scope MUST NOT expose discussion content beyond the viewers authorized for that issue.

### Candidate conversation structure

The founder has proposed one top-level contribution per person on each issue, with responses shown as a single flat chronological discussion beneath that contribution. This could keep each participant's main position distinct while avoiding deeply nested threads. It is a candidate design, not yet a confirmed requirement.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../README.md) | Owns the issue record, creator, nomination, score, rank, and original open age that discussion must not alter |
| [issue/lifecycle](../lifecycle/README.md) | Defines closure and confirms that comments are activity rather than lifecycle transitions |
| [issue/anonymity](../anonymity/README.md) | Must prevent discussion behavior from revealing the creator of an anonymous issue |
| [issue/visibility](../visibility/README.md) | Supplies the trust boundary within which discussion may be read |
| [issue/moderation](../moderation/README.md) | Must distinguish moderating an individual contribution from suppressing the issue itself |
| [permissions](../../permissions/README.md) | Will define who may read and contribute to discussion |
| [voting](../../voting/README.md) | Remains the only source of budgeted support; discussion volume is not a vote signal |

## Dependencies

- issue
- permissions
- anonymity
- moderation

## Acceptance Criteria

### AC: discussion-preserves-issue-state

- **Given** a raised issue has a known nomination, support score, rank, and original creation time
- **When** discussion activity is recorded on that issue
- **Then** its lifecycle status, nomination, support score, rank inputs, and original open-age basis remain unchanged

### AC: discussion-does-not-cross-visibility-boundary

- **Given** a person is not authorized to view an issue
- **When** that person attempts to view the issue's discussion
- **Then** no discussion content or discussion-author metadata is exposed

## Open Questions

- Should each person be limited to one top-level contribution per issue?
- If a person has one top-level contribution, may they edit or replace it while preserving an audit history?
- Should all responses beneath a top-level contribution appear in one flat chronological list with no nested reply levels?
- Who may write top-level contributions and responses: every viewer, members of the issue's originating scope, or a narrower permission?
- May a response be anonymous, and how can the author of an anonymous issue participate without revealing that they created it?
- May authors delete discussion content, or only edit it with visible history?
- How do issue moderation rules apply to individual discussion entries without letting moderation suppress the issue itself?
- After an issue reaches a terminal status, does its discussion become read-only or remain open?
- Which discussion events notify the issue creator, supporters, assignee, or prior participants?

---
*This document follows the https://specscore.md/feature-specification*
