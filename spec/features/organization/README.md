---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Organization

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

IssueNumber.one models a company as a root team with nested department, project, and team nodes. Each node has its own candidate set and top-N ranking, and automatically contributes only its current #1 issue to its parent. Votes remain assigned to the issue as its visibility changes. Leaders may drill from the root through every node to the originating person and issue. Public topics live outside this hierarchy.

## Contents

| Directory | Description |
|-----------|-------------|
| [topic/](topic/README.md) | Public topics and their sub-topics — a parallel hierarchy for cross-org discussion |

### topic

Topics are public discussion spaces not owned by any one org. Anyone can create a public topic or a sub-topic within one. Topic issues are always public and cannot be archived.

## Problem

A prioritization tool needs a clear scope model. Without one, users cannot tell whether an issue belongs to their immediate team, a project, a department, or the wider company, and leaders cannot see whether the issue reaching them represents collective support. IssueNumber.one treats every organizational level as a nested scope with its own ranking and automatic #1 bubble-up.

## Behavior

### Org is a root-level team

There is no separate "organization" entity. An organization is simply a team that has no parent. This lets any org also become a sub-team of another org in the future.

#### REQ: org-is-root-team

An organization MUST be modeled as a team with no parent. When an org joins another org as a sub-team, it becomes an ordinary team in that hierarchy.

#### REQ: organization-creation-uses-spaceus-admin-default

Creating a root organization in IssueNumber.one MUST use the shared Spaceus creation behavior, under which the authenticated creator becomes the first organization administrator. The universal administrator bootstrap is owned by Spaceus, not redefined by this product.

### Teams nest

A team MAY have any number of sub-teams. A sub-team is itself a team and has its own issues, candidate set, ranking, and membership.

#### REQ: team-nesting

Teams MUST support nested sub-teams to arbitrary depth.

### Organizational names are roles in one hierarchy

The root team represents the company or organization. Department, project, and team are names for nested team nodes rather than different domain entities.

#### REQ: hierarchy-role-names

The model MUST support labeling nested team nodes as company, department, project, team, or another organization-defined scope type without changing the ranking mechanism.

### #1 flows upward one level at a time

At the immediate team level, members' nominated personal #1 issues form the candidate set. At every parent level, each direct child contributes only its current #1. Members of the parent scope rank those candidates, and the parent's #1 automatically continues upward.

#### REQ: recursive-top-flow

Every non-root scope's current #1 issue MUST automatically become a candidate in its immediate parent scope, and no lower-ranked issue from that child MUST enter the parent ranking through the automatic path.

#### REQ: scope-top-n-local

Each organizational scope MUST maintain its own top-N ranking for eligible members while contributing only its #1 to its parent.

### Membership and voting

A department or company includes the members of all descendant teams for voting purposes. A manager does not control whether a collectively selected child #1 enters the parent ranking.

#### REQ: descendant-members-belong-to-voting-scope

Members of descendant teams MUST be treated as members of every enclosing organizational voting scope.

#### REQ: no-manager-promotion-gate

The organizational hierarchy MUST NOT add a manager approval or veto between a child scope selecting its #1 and that issue becoming a parent-scope candidate.

### Drill-down

Authorized viewers can inspect the top N at a company, select a department, then continue through project, team, person, and originating issue.

#### REQ: drill-down-preserves-origin

A bubbled issue MUST retain its complete source path so an authorized viewer can drill down to the originating scope, person attribution when not anonymous, and issue.

### Peer teams

A team MAY declare peer teams — other teams at a similar level of the hierarchy whose priorities are especially relevant. Peer relationships are informational; they do not grant cross-team voting by themselves.

#### REQ: peer-teams

Teams MUST be able to declare a list of peer teams. Peer declarations MUST NOT by themselves grant voting or read permissions.

### Peer colleagues

A team MAY declare peer colleagues — individuals outside the team who the team considers close collaborators. A team member MAY also declare their own peers.

#### REQ: peer-colleagues

Teams and individual members MUST be able to declare peer colleagues. Peer colleague relationships MUST NOT by themselves grant voting or read permissions.

### Creation and destruction

| Action | Who |
|--------|-----|
| Create a team inside an org | Any org member |
| Create a public topic or sub-topic | Anyone |
| Archive (dissolve) a team | Any org member |

#### REQ: any-member-creates-team

Any member of an organization MUST be allowed to create a new team within it.

#### REQ: any-member-archives-team

Any member of an organization MUST be allowed to archive any team within that organization (see also [permissions#req-any-org-member-archives-team](../permissions/README.md)).

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../issue/README.md) | Every issue is scoped to a team, org, or topic |
| [voting](../voting/README.md) | Each scope has an eligible candidate set and ranking, while every vote stays attached to its issue and each member's 1–7 budget is shared across the organization |
| [permissions](../permissions/README.md) | Membership determines visibility and action rights |
| [storage](../storage/README.md) | Orgs choose where their data lives |
| [Spaceus platform Feature](https://github.com/sneat-co/backstage/blob/main/spec/features/spaceus/README.md) | Owns organization creation and the universal creator-as-first-administrator default |

## Dependencies

- permissions

## Acceptance Criteria

Not defined yet.

## Open Questions

- Can a user belong to multiple orgs, and if so, are budgets independent per org?
- How does `space` relate to `team` and `topic`: a third scope, an alias, or a super-type?
- May parent scopes contain directly raised issues in addition to child #1 candidates?
- When a team is archived, what happens to its sub-teams and issues?
- Can a sub-team outlive its parent being archived (i.e., be re-parented)?
- Is there any rate-limit on team creation to prevent spam?
- Acceptance criteria not yet defined for this feature.
