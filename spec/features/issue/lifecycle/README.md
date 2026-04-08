# Feature: Issue Lifecycle

**Status:** Conceptual

## Summary

The set of statuses an issue can occupy and the rules governing transitions between them.

## Problem

Without a defined lifecycle, it is unclear whether a stale issue is still open, whether a resolved issue can be reopened, or who is allowed to close an issue. A formal lifecycle gives the product, its users, and its tooling a shared vocabulary for the state of every issue.

## Behavior

### Statuses

| Status | Description |
|--------|-------------|
| `raised` | The default state after an issue is created — actively open and eligible for votes |
| `withdrawn` | The creator has pulled the issue; sub-reasons include "lost priority" and "not actual" |
| `resolved` | The issue has been addressed |
| `banned` | Moderated out by authorized members; visible only as a tombstone |
| `archived` | Closed by a team member without resolution; common for stale or superseded issues |

#### REQ: allowed-statuses

An issue's `status` MUST be one of: `raised`, `withdrawn`, `resolved`, `banned`, `archived`.

#### REQ: initial-status-raised

A newly-created issue MUST have status `raised`.

### Withdrawal sub-reasons

When an issue is withdrawn, a sub-reason MUST be captured. The supported sub-reasons are:

| Sub-reason | Meaning |
|------------|---------|
| `lost-priority` | The creator no longer considers this their top priority |
| `not-actual` | The issue is no longer relevant |

#### REQ: withdraw-requires-reason

A withdrawal MUST include one of the supported sub-reasons.

### Who can transition

| Transition | Who | Notes |
|------------|-----|-------|
| `raised → withdrawn` | Creator only | Supporters are notified |
| `raised → resolved` | Any team member (team issues); creator (public topics) | Supporters are notified |
| `raised → archived` | Any team member (team issues) | Not allowed for public topic issues |
| `raised → banned` | Authorized moderators (see [moderation](../moderation/README.md)) | Supporters and creator notified |

#### REQ: public-topic-no-archive

Issues in public topics MUST NOT be archivable. They can only be withdrawn by their creator.

#### REQ: terminal-status-irreversible

Once an issue enters `withdrawn`, `resolved`, `banned`, or `archived`, it MUST NOT return to `raised`.

### Vote refunds

When an issue exits the `raised` state for any reason, every vote cast on that issue MUST be refunded to its voter (see [voting](../../voting/README.md)).

#### REQ: refund-on-exit-raised

Exiting the `raised` state MUST trigger refunds of all votes cast on the issue.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue/moderation](../moderation/README.md) | Defines the actors allowed to move an issue to `banned` or `archived` |
| [voting](../../voting/README.md) | Transitions out of `raised` trigger vote refunds |
| [permissions](../../permissions/README.md) | Gates who is allowed to perform each transition |

## Dependencies

- issue
- voting
- permissions

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Should resolved issues be archivable afterward, or do they remain in `resolved` indefinitely?
- Should there be a `dormant` or `snoozed` intermediate state, or is withdrawal sufficient?
- Who may resolve a public topic issue — only the creator, or any member of the public topic?
- Acceptance criteria not yet defined for this feature.
