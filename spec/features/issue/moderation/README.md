---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Issue Moderation

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/moderation?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/moderation?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/moderation?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/moderation?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Moderation removes abusive or policy-violating content without becoming a route for managers or ordinary members to suppress an inconvenient active issue. A ban requires an authorized moderator, a visible reason, and notification of the creator and supporters. Archival applies only to already-terminal records and is governed by lifecycle retention.

## Problem

A product intended to surface issues through collective support must handle abuse without allowing authority figures to erase persistent dissent. Moderation therefore addresses policy violations, not staleness, low rank, disagreement, or management inconvenience.

## Behavior

### Retention archive vs moderation ban

| Action | Who | Effect |
|--------|-----|--------|
| `archive` | Retention policy or authorized custodian | Moves an already-terminal record out of normal history views; never closes a raised issue |
| `ban` | Authorized moderators only | Removes policy-violating content; visible as a tombstone with reason |

#### REQ: archive-is-soft

Archiving MUST preserve an already-terminal issue and its closure reason for later reference. It MUST NOT be available as a transition from `raised`.

#### REQ: ban-is-tombstone

A banned issue's content MUST be hidden from normal views; only a tombstone (author redacted, reason shown) MUST remain.

### No suppression through moderation

An issue's age, low rank, disagreement, lack of management response, or inconvenience are not moderation grounds.

#### REQ: ordinary-members-cannot-moderate-close

An ordinary member or manager MUST NOT be able to archive, ban, withdraw, or resolve another creator's `raised` issue merely to remove it from active views.

#### REQ: ban-requires-reason

Every ban MUST record an authorized moderator, a policy reason, and the time of the action.

### Notifications

Every archive or ban action MUST notify the issue's creator and all current supporters.

#### REQ: notify-on-moderation

When an issue is archived or banned, the system MUST notify its creator and every current supporter.

### Vote refunds

A ban exits the `raised` state, which triggers refunds per [issue/lifecycle#req-refund-on-exit-raised](../lifecycle/README.md). Retention archival occurs only after votes have already been refunded by an earlier terminal transition.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue/lifecycle](../lifecycle/README.md) | Moderation may drive a transition to `banned`; archival only follows a terminal state |
| [voting](../../voting/README.md) | Moderation triggers vote refunds |
| [permissions](../../permissions/README.md) | Defines who is allowed to moderate |

## Dependencies

- issue/lifecycle
- voting
- permissions

## Acceptance Criteria

Not defined yet.

## Open Questions

- Should there be an appeal or review mechanism for banned issues?
- Who is authorized to ban — only org admins, or can teams elect moderators?
- Who is the authorized retention custodian for already-terminal issues, and can an archive action be reversed?
- Acceptance criteria not yet defined for this feature.
