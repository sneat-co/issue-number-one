---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Permissions

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/permissions?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/permissions?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/permissions?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/permissions?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Defines the access matrix for every action in IssueNumber.one: who can raise issues, see personal and ranked views, vote at each organizational scope, confirm closure, moderate content, administer organization-wide settings, and create or archive teams and topics. Most permissions derive from membership and the issue's place in the #1 bubble-up path. A constrained organization-administrator role controls settings such as the shared vote budget but receives no special power to block promotion, close an inconvenient issue, or cast extra votes.

## Problem

A product with several scopes (team, org, topic, public) and several actor classes (creator, supporter, team member, org member, platform user) needs a single place that spells out who can do what. Scattering permission rules across individual feature specs risks contradictions and makes it hard to reason about safety. This feature centralizes the rules and cross-links each one to the features that enforce it.

## Behavior

### Scopes and actor classes

| Scope | Actors |
|-------|--------|
| Team / space | Team member, non-member org member, outsider |
| Org | Org member, organization administrator, outsider |
| Public topic | Any authenticated user, anonymous visitor |
| Issue | Creator, supporter, anyone with read access to the enclosing scope |

### The access matrix

| Action | Team-scoped | Org-scoped | Topic-scoped |
|--------|-------------|------------|--------------|
| See issues | Creator sees all their own; eligible members see voting candidates and top N | Members see top N at each enclosing scope and the #1 from each direct child | Anyone |
| Raise an issue | Only team members | Only team members (via their team) | Any authenticated user |
| Vote on an issue | Any member eligible in that voting scope, subject to budget | Any member of the scope, including descendant-team members, subject to budget | Subject to topic rules |
| Change the shared vote budget | — | Organization administrators only | — |
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

### Membership permissions and a constrained administrator role

Ordinary participation permissions derive from scope membership. The MVP also has an organization-administrator role for explicitly named organization-wide settings, including the shared vote budget. Administrator status does not change issue-ranking power or lifecycle authority. Moderator roles remain separate and are reserved for the `ban` action.

#### REQ: permissions-derive-from-membership

All ordinary participation permissions MUST derive from scope membership plus the rules in this document. A named organization-administrator role MAY grant only the organization-wide setting permissions explicitly assigned to it.

#### REQ: organization-administrator-role

The MVP MUST support designating one or more organization members as organization administrators.

#### REQ: organization-creation-inherits-spaceus-admin-default

IssueNumber.one MUST create organizations through the shared Spaceus organization-creation behavior, which automatically designates the creator as the first administrator. IssueNumber.one MUST NOT define a conflicting product-local initial-administrator rule.

#### REQ: administrator-appointment-inherits-spaceus

IssueNumber.one MUST inherit the shared Spaceus behavior that allows any existing organization administrator to designate another existing member of that organization as an administrator. A non-administrator MUST NOT be allowed to grant organization-administrator status, and IssueNumber.one MUST NOT define a conflicting product-local appointment rule.

#### REQ: administrator-removal-inherits-spaceus

IssueNumber.one MUST inherit the shared Spaceus behavior that allows any existing organization administrator to remove administrator status from another administrator. A non-administrator MUST NOT be allowed to remove organization-administrator status, and IssueNumber.one MUST NOT define a conflicting product-local removal rule.

#### REQ: admin-only-vote-budget-configuration

Only an organization administrator MUST be allowed to change that organization's shared vote-budget size.

#### REQ: administrator-does-not-control-priority

Administrator status MUST NOT grant extra votes, a manager promotion gate, access to unrelated non-candidate personal issues, or authority to withdraw, resolve, or archive another person's raised issue.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| Every local feature performing an action | References the relevant `REQ:` in this document |
| [Spaceus platform Feature](https://github.com/sneat-co/backstage/blob/main/spec/features/spaceus/README.md) | Owns the universal organization-administrator bootstrap, appointment, and removal rules inherited by IssueNumber.one |

## Dependencies

- organization

## Acceptance Criteria

Not defined yet.

## Open Questions

- Should the product later support team-level administrators in addition to organization administrators?
- How are platform-level moderators (for ban actions) assigned, and by whom?
- Should there be a "guest" or "observer" role that can read a team without being able to raise or vote?
- Should permissions support per-team overrides of the defaults, or are the defaults inviolable?
- Acceptance criteria not yet defined for this feature.
