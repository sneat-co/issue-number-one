---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Issue Lifecycle

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/lifecycle?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/lifecycle?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/lifecycle?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/lifecycle?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

The set of statuses an issue can occupy and the rules governing transitions between them, while keeping rank, nomination, and management activity separate from closure.

## Problem

Without a defined lifecycle, it is unclear whether a stale issue is still open, whether a resolved issue can be reopened, or who is allowed to close an issue. A formal lifecycle gives the product, its users, and its tooling a shared vocabulary for the state of every issue.

## Behavior

### Statuses

| Status | Description |
|--------|-------------|
| `raised` | The default state after creation — actively open; nomination and vote eligibility are separate properties |
| `withdrawn` | The creator confirms that the item should no longer be treated as a current issue |
| `resolved` | The creator, or eligible peers when the creator is unavailable, confirms that the issue has been addressed |
| `banned` | Moderated out by authorized members; visible only as a tombstone |
| `archived` | A retained terminal record moved out of normal history views after it was withdrawn, resolved, or banned |

#### REQ: allowed-statuses

An issue's `status` MUST be one of: `raised`, `withdrawn`, `resolved`, `banned`, `archived`.

#### REQ: initial-status-raised

A newly-created issue MUST have status `raised`.

### Withdrawal sub-reasons

When an issue is withdrawn, a sub-reason MUST be captured. The supported sub-reasons are:

| Sub-reason | Meaning |
|------------|---------|
| `not-actual` | The issue is no longer relevant |

#### REQ: withdraw-requires-reason

A withdrawal MUST include one of the supported sub-reasons.

### Who can transition

| Transition | Who | Notes |
|------------|-----|-------|
| `raised → withdrawn` | Creator only | Supporters are notified |
| `raised → resolved` | Creator | Management activity alone cannot make this transition |
| `raised → resolved` | Eligible peers by closure vote | Allowed only when the creator is unavailable; closure method is recorded |
| `raised → banned` | Authorized moderators (see [moderation](../moderation/README.md)) | Supporters and creator notified |
| `withdrawn/resolved/banned → archived` | Retention policy or authorized custodian | Does not alter the recorded closure reason |

#### REQ: creator-controls-normal-resolution

Only the creator MUST normally be allowed to transition their issue from `raised` to `resolved`.

#### REQ: peer-closure-only-when-creator-unavailable

Eligible peers MAY transition an issue from `raised` to `resolved` by vote only when the creator is recorded as unavailable. The system MUST record that the closure used peer voting rather than creator confirmation.

#### REQ: no-direct-archive-from-raised

A `raised` issue MUST NOT transition directly to `archived`.

#### REQ: public-topic-no-archive

Raised issues in public topics MUST NOT be directly archivable. They can be withdrawn by their creator, resolved under the same creator-control rule, or banned through moderation.

#### REQ: terminal-status-irreversible

Once an issue enters `withdrawn`, `resolved`, `banned`, or `archived`, it MUST NOT return to `raised`.

### Rank and activity do not close an issue

An issue's nomination, votes, rank, assignment, comments, and management response are activity around the issue, not lifecycle transitions.

#### REQ: rank-change-is-not-lifecycle-transition

An issue that is de-nominated, loses votes, or drops out of a top-N list MUST remain `raised` and retain its original creation time.

#### REQ: management-response-is-not-lifecycle-transition

Acknowledgement, assignment, explanation, progress, or a claimed fix MUST NOT close an issue without creator confirmation or the unavailable-creator peer-closure process.

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

## Open Questions

- Should resolved issues be archivable afterward, or do they remain in `resolved` indefinitely?
- What proves that a creator is unavailable, which peers may participate, and what vote threshold closes the issue?
- Who or what retention policy may archive an already-terminal issue, and after how long?
- Acceptance criteria not yet defined for this feature.
