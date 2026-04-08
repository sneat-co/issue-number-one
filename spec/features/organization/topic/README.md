# Feature: Topic

**Status:** Conceptual

## Summary

A topic is a public discussion space that is not owned by any single organization. Anyone may create a public topic or a sub-topic inside an existing public topic. Issues in topics are always public, cannot be archived, and can only be withdrawn by their creator.

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

### Moderation constraints

Topic issues cannot be archived. The only way to close them is withdrawal by the creator, or — in extreme cases — banning by platform moderators.

#### REQ: topic-issues-not-archivable

Issues in topics MUST NOT be archivable by anyone. See also [issue/lifecycle#req-public-topic-no-archive](../../issue/lifecycle/README.md).

### Storage

Public topics are stored in a public GitHub repository, proxied by an API layer for performance. See [storage/git-storage](../../storage/git-storage/README.md).

#### REQ: topics-in-public-github

Public topics MUST be stored in a public GitHub repository via the git-storage backend.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [organization](../README.md) | Topics live outside the org hierarchy but use the same nesting model |
| [issue/visibility](../../issue/visibility/README.md) | Topics are always `public` |
| [storage/git-storage](../../storage/git-storage/README.md) | Topics are persisted as git-storage |

## Dependencies

- organization
- issue/visibility
- storage/git-storage

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Who moderates spam or abuse in public topics — platform-level moderators, or delegated?
- Can a topic be "claimed" by an org and transitioned into org-owned storage?
- Are there voting budgets in topics, and if so, whose budget applies?
- Should anonymous issues be allowed in public topics?
- Acceptance criteria not yet defined for this feature.
