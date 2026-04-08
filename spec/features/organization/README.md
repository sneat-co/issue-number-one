# Feature: Organization

**Status:** Conceptual

## Summary

IssueNumber.one organizes everything around teams. An organization is simply a root-level team. Teams may contain sub-teams, have peer teams, and have peer colleagues; team members may specify their own peers. Public topics live outside the org hierarchy in their own namespace.

## Contents

| Directory | Description |
|-----------|-------------|
| [topic/](topic/README.md) | Public topics and their sub-topics — a parallel hierarchy for cross-org discussion |

### topic

Topics are public discussion spaces not owned by any one org. Anyone can create a public topic or a sub-topic within one. Topic issues are always public and cannot be archived.

## Problem

A prioritization tool needs a clear scope model. Without one, users can't tell whether an issue belongs to their immediate team, a larger initiative, or the wider organization. IssueNumber.one treats "org" and "team" uniformly — an org is just a team with no parent — and lets teams nest, cross-link as peers, and bubble priorities up the tree.

## Behavior

### Org is a root-level team

There is no separate "organization" entity. An organization is simply a team that has no parent. This lets any org also become a sub-team of another org in the future.

#### REQ: org-is-root-team

An organization MUST be modeled as a team with no parent. When an org joins another org as a sub-team, it becomes an ordinary team in that hierarchy.

### Teams nest

A team MAY have any number of sub-teams. A sub-team is itself a team and has its own issues, budget, and membership.

#### REQ: team-nesting

Teams MUST support nested sub-teams to arbitrary depth.

### Peer teams

A team MAY declare peer teams — other teams at a similar level of the hierarchy whose priorities are especially relevant. Peer relationships are informational; they do not grant cross-team voting by themselves.

#### REQ: peer-teams

Teams MUST be able to declare a list of peer teams. Peer declarations MUST NOT by themselves grant voting or read permissions.

### Peer colleagues

A team MAY declare peer colleagues — individuals outside the team who the team considers close collaborators. A team member MAY also declare their own peers.

#### REQ: peer-colleagues

Teams and individual members MUST be able to declare peer colleagues. Peer colleague relationships MUST NOT by themselves grant voting or read permissions.

### Minimum team size

Per [voting/budget#req-minimum-team-size](../voting/budget/README.md), a team must have at least three members for voting mechanics to apply.

### Creation and destruction

| Action | Who |
|--------|-----|
| Create a team inside an org | Any org member |
| Create a public topic or sub-topic | Anyone |
| Archive (dissolve) a team | Any org member |

#### REQ: any-member-creates-team

Any member of an organization MUST be allowed to create a new team within it.

#### REQ: any-member-archives-team

Any member of an organization MUST be allowed to archive any team within that organization (see also [issue/moderation#req-any-org-member-archives-team](../issue/moderation/README.md)).

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../issue/README.md) | Every issue is scoped to a team, org, or topic |
| [voting](../voting/README.md) | Budgets are per-team |
| [permissions](../permissions/README.md) | Membership determines visibility and action rights |
| [storage](../storage/README.md) | Orgs choose where their data lives |

## Dependencies

- permissions

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Can a user belong to multiple orgs, and if so, are budgets independent per org?
- When a team is archived, what happens to its sub-teams and issues?
- Can a sub-team outlive its parent being archived (i.e., be re-parented)?
- Is there any rate-limit on team creation to prevent spam?
- Acceptance criteria not yet defined for this feature.
