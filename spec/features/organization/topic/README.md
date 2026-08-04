---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Topic

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization/topic?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization/topic?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization/topic?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization/topic?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

A topic is a public discussion space that is not owned by any single organization. Anyone may create a public topic or a sub-topic inside an existing public topic. A participant may nominate one personal #1 per topic, which receives one free creator star. Issues are always public and cannot be directly archived while raised.

## Problem

Some priority discussions do not belong to a specific company or team — for example, open-source projects, community concerns, or cross-org initiatives. Without a neutral public space, these conversations either get forced into an org hierarchy they don't fit or happen off-platform entirely. Topics give them a home.

## Behavior

### Creation is open

Any authenticated user MAY create a new top-level public topic. Anyone MAY create a sub-topic inside an existing public topic.

#### REQ: anyone-creates-topic

Any authenticated user MUST be allowed to create a public top-level topic or a sub-topic within an existing public topic.

### Topics nest

Topics support sub-topics to arbitrary depth, mirroring team nesting but in the public namespace.

#### REQ: topic-nesting

Topics MUST support nested sub-topics to arbitrary depth.

### Always public

Topic issues are always `public` visibility. There is no team-scoped mode.

#### REQ: topic-issues-public

All issues in topics MUST have `public` visibility. See also [issue/visibility#req-public-topic-always-public](../../issue/visibility/README.md).

### One personal #1 per participant and topic

A participant may keep multiple issues in a public topic but may nominate at most one personal #1 in that topic. The topic nomination is separate from their nomination in any organization or other public topic.

#### REQ: one-personal-top-per-public-topic

A participant MUST have at most one personal #1 nomination within a public topic at a time.

#### REQ: topic-personal-top-gets-creator-star

A participant's personal #1 in a public topic MUST automatically receive one free creator star without consuming any budgeted support allocation.

### Moderation constraints

Raised topic issues cannot be directly archived. They may be withdrawn or resolved by the creator, resolved by eligible peers when the creator is unavailable under the lifecycle rules, or banned by platform moderators.

#### REQ: topic-issues-not-archivable

Raised issues in topics MUST NOT be directly archivable by anyone. See also [issue/lifecycle#req-public-topic-no-archive](../../issue/lifecycle/README.md).

### Storage

Public topics are stored in a public GitHub repository, proxied by an API layer for performance. See [storage/git-storage](../../storage/git-storage/README.md).

#### REQ: topics-in-public-github

Public topics MUST be stored in a public GitHub repository via the git-storage backend.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [organization](../README.md) | Topics live outside the org hierarchy but use the same nesting model |
| [issue](../../issue/README.md) | Each participant may nominate one personal #1 per topic and receives one free creator star |
| [issue/visibility](../../issue/visibility/README.md) | Topics are always `public` |
| [voting](../../voting/README.md) | Topic candidate eligibility uses each participant's personal #1 for that topic; budget rules remain open |
| [storage/git-storage](../../storage/git-storage/README.md) | Topics are persisted as git-storage |

## Dependencies

- organization
- issue
- issue/visibility
- voting
- storage/git-storage

## Acceptance Criteria

Not defined yet.

## Open Questions

- Who moderates spam or abuse in public topics — platform-level moderators, or delegated?
- Can a topic be "claimed" by an org and transitioned into org-owned storage?
- Are there voting budgets in topics, and if so, whose budget applies?
- Should anonymous issues be allowed in public topics?
- Acceptance criteria not yet defined for this feature.
