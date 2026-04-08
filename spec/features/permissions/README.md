# Feature: Permissions

**Status:** Conceptual

## Summary

Defines the access matrix for every action in IssueNumber.one: who can raise issues, who can see them, who can vote, and who can create or archive teams and topics. Permissions are deliberately simple — they derive from membership in a team, org, or topic, plus the small set of explicit rules defined here.

## Problem

A product with several scopes (team, org, topic, public) and several actor classes (creator, supporter, team member, org member, platform user) needs a single place that spells out who can do what. Scattering permission rules across individual feature specs risks contradictions and makes it hard to reason about safety. This feature centralizes the rules and cross-links each one to the features that enforce it.

## Behavior

### Scopes and actor classes

| Scope | Actors |
|-------|--------|
| Team / space | Team member, non-member org member, outsider |
| Org | Org member, outsider |
| Public topic | Any authenticated user, anonymous visitor |
| Issue | Creator, supporter, anyone with read access to the enclosing scope |

### The access matrix

| Action | Team-scoped | Org-scoped | Topic-scoped |
|--------|-------------|------------|--------------|
| See issues | Team members only | Any org member (top issues + bubble-ups) | Anyone |
| Raise an issue | Only team members | Only team members (via their team) | Any authenticated user |
| Vote on an issue | Subject to per-team budget and cross-team support config | Same | Subject to topic rules |
| Withdraw an issue | Creator only | Creator only | Creator only |
| Archive an issue | Any team member | Any org admin | No one |
| Archive a team | — | Any org member | — |
| Create a team | — | Any org member | — |
| Create a topic/sub-topic | — | — | Any authenticated user |
| Ban an issue | Moderators | Moderators | Platform moderators |

#### REQ: team-raise-requires-membership

Only members of a team MUST be allowed to raise an issue in that team.

#### REQ: team-read-requires-membership

By default, only members of a team MUST be allowed to see team-level (non-top) issues for that team.

#### REQ: org-wide-team-visibility

Any member of an organization MUST be allowed to see every team inside that organization and each team's top issue and org-level issues.

#### REQ: creator-only-withdraw

Only the creator of an issue MUST be allowed to withdraw it. See also [issue#req-creator-only-withdraw](../issue/README.md).

#### REQ: any-team-member-archives-team-issue

Any member of a team MUST be allowed to archive any issue in that team. See also [issue/moderation#req-any-team-member-archives-team-issue](../issue/moderation/README.md).

#### REQ: any-org-member-archives-team

Any member of an organization MUST be allowed to archive any team within that organization. See also [organization#req-any-member-archives-team](../organization/README.md).

#### REQ: any-org-member-creates-team

Any member of an organization MUST be allowed to create a team within it. See also [organization#req-any-member-creates-team](../organization/README.md).

#### REQ: anyone-creates-public-topic

Any authenticated user MUST be allowed to create a public top-level topic or a sub-topic within an existing public topic. See also [organization/topic#req-anyone-creates-topic](../organization/topic/README.md).

### Derivation, not roles

Permissions derive from scope membership, not from a named role system. There is no "admin" role in the base product. Explicit moderator roles are reserved for the `ban` action.

#### REQ: permissions-derive-from-membership

All permission checks MUST derive from scope membership plus the rules in this document. The base product MUST NOT require a named-role system.

## Interaction with Other Features

Every feature that performs an action references the relevant `REQ:` in this document.

## Dependencies

- organization

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Should teams be allowed to promote explicit admin roles beyond what the base product defines?
- How are platform-level moderators (for ban actions) assigned, and by whom?
- Should there be a "guest" or "observer" role that can read a team without being able to raise or vote?
- Should permissions support per-team overrides of the defaults, or are the defaults inviolable?
- Acceptance criteria not yet defined for this feature.
