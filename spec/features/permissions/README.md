---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Permissions

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/permissions?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/permissions?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/permissions?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/permissions?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Defines the access matrix for every action in IssueNumber.one: who can raise issues, see personal and ranked views, vote at each organizational scope, confirm closure, moderate content, and create or archive teams and topics. Permissions derive from membership and the issue's place in the #1 bubble-up path; a manager receives no special power to block promotion or close an inconvenient issue.

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
| See issues | Creator sees all their own; eligible members see voting candidates and top N | Members see top N at each enclosing scope and the #1 from each direct child | Anyone |
| Raise an issue | Only team members | Only team members (via their team) | Any authenticated user |
| Vote on an issue | Any member eligible in that voting scope, subject to budget | Any member of the scope, including descendant-team members, subject to budget | Subject to topic rules |
| Withdraw an issue | Creator only | Creator only | Creator only |
| Resolve an issue | Creator; eligible peers by vote only if creator unavailable | Same | Same |
| Archive a raised issue | No one | No one | No one |
| Archive a team | — | Any org member | — |
| Create a team | — | Any org member | — |
| Create a topic/sub-topic | — | — | Any authenticated user |
| Ban an issue | Moderators | Moderators | Platform moderators |

#### REQ: team-raise-requires-membership

Only members of a team MUST be allowed to raise an issue in that team.

#### REQ: team-read-requires-membership

Only the creator MUST automatically see their complete list of non-nominated issues. Other members MUST see issues when those issues are eligible candidates or appear in a top-N view for a scope in which the viewer is a member.

#### REQ: org-wide-team-visibility

Any member of an organization MUST be allowed to see the top N at their enclosing scopes and the current #1 supplied by each direct child of a scope they may view. This MUST NOT grant access to unrelated non-candidate personal issues.

#### REQ: descendant-member-may-vote-in-scope

Any member of a department or other organizational scope, including a member whose direct membership is in a descendant team, MUST be allowed to vote on candidates in that scope.

#### REQ: manager-cannot-block-bubble-up

No manager or scope owner MUST be able to block a child scope's current #1 from automatically becoming a parent-scope candidate.

#### REQ: creator-only-withdraw

Only the creator of an issue MUST be allowed to withdraw it. See also [issue#req-creator-only-withdraw](../issue/README.md).

#### REQ: creator-controls-normal-resolution

Only the creator MUST normally be allowed to resolve their issue. Eligible peers MAY resolve it by vote only when the creator is unavailable, as defined by [issue/lifecycle#req-peer-closure-only-when-creator-unavailable](../issue/lifecycle/README.md).

#### REQ: no-direct-archive-of-raised-issue

No ordinary member, manager, or administrator MUST be able to archive another person's `raised` issue. Policy-violating content may be banned only through the visible moderation process.

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

## Open Questions

- Should teams be allowed to promote explicit admin roles beyond what the base product defines?
- How are platform-level moderators (for ban actions) assigned, and by whom?
- Should there be a "guest" or "observer" role that can read a team without being able to raise or vote?
- Should permissions support per-team overrides of the defaults, or are the defaults inviolable?
- Acceptance criteria not yet defined for this feature.
