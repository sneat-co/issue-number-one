# Feature: Issue Moderation

**Status:** Conceptual

## Summary

Moderation covers archiving and banning issues. Any team member may archive any team issue; any org member may archive any team; public topic issues cannot be archived. Moderation always notifies the creator and supporters.

## Problem

A product that caps active issues per person must also have a safety valve for clearing out stale, duplicate, off-topic, or abusive issues. Without moderation, junk issues consume budget slots and crowd out real priorities. But moderation must be visible and accountable so it isn't abused to silence dissent — every action notifies the affected parties.

## Behavior

### Archive vs ban

| Action | Who | Effect |
|--------|-----|--------|
| `archive` | Any team member (team issues); any org member (whole teams) | Closes the issue without marking it resolved; stays visible in archive |
| `ban` | Authorized moderators only | Moderated out; visible only as a tombstone with reason |

#### REQ: archive-is-soft

Archiving MUST close the issue in its current state and preserve its content for later reference.

#### REQ: ban-is-tombstone

A banned issue's content MUST be hidden from normal views; only a tombstone (author redacted, reason shown) MUST remain.

### Who can archive

| Scope | Who can archive an issue |
|-------|--------------------------|
| Team issue | Any member of that team |
| Org-level issue | Any org admin |
| Public topic issue | No one — only the creator may withdraw |

#### REQ: any-team-member-archives-team-issue

Any member of a team MUST be allowed to archive any issue within that team.

#### REQ: any-org-member-archives-team

Any member of an organization MUST be allowed to archive an entire team within that organization.

#### REQ: public-topic-no-archive

Public topic issues MUST NOT be archivable by anyone. See also [issue/lifecycle#req-public-topic-no-archive](../lifecycle/README.md).

### Notifications

Every archive or ban action MUST notify the issue's creator and all supporters.

#### REQ: notify-on-moderation

When an issue is archived or banned, the system MUST notify its creator and every current supporter.

### Vote refunds

Moderation exits the `raised` state, which triggers refunds per [issue/lifecycle#req-refund-on-exit-raised](../lifecycle/README.md).

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue/lifecycle](../lifecycle/README.md) | Moderation drives transitions to `archived` or `banned` |
| [voting](../../voting/README.md) | Moderation triggers vote refunds |
| [permissions](../../permissions/README.md) | Defines who is allowed to moderate |

## Dependencies

- issue/lifecycle
- voting
- permissions

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Can archived issues be un-archived, or is archival terminal?
- Should there be an appeal or review mechanism for banned issues?
- Who is authorized to ban — only org admins, or can teams elect moderators?
- Should archiving an entire team cascade-archive all its issues, or leave them readable as a frozen snapshot?
- Acceptance criteria not yet defined for this feature.
