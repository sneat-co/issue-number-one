# Feature: Git Storage

**Status:** Conceptual

## Summary

An alternative storage backend in which an org's data lives in a GitHub repository (public or private), persisted as strongly-structured, human-readable files following the [inGitDB](https://inGitDB.com) schema. Git history becomes the audit log; GitHub permissions govern access.

## Problem

Some organizations — especially developer-heavy teams, open-source projects, and compliance-sensitive groups — want their prioritization data to live in their own infrastructure. A GitHub repository is familiar, auditable, and already under their control. Git storage lets them keep their data where they already keep everything else, and use git history as the canonical change log.

## Behavior

### Backed by inGitDB

Git storage uses [inGitDB](https://inGitDB.com) for its on-disk schema. IssueNumber.one reads and writes files according to an inGitDB schema specific to its domain.

#### REQ: uses-ingitdb

Git storage MUST persist data via the inGitDB file schema.

### Dedicated inGitDB schema

IssueNumber.one's git-storage format requires its own schema definition, covering teams, issues, votes, permissions, and all other entities.

#### REQ: dedicated-schema

IssueNumber.one MUST have a dedicated inGitDB schema definition for its entities, separate from any other inGitDB user's schema.

### Repository choice

An org may connect an existing GitHub repository or create a new one at setup. The repository may be public or private; public topics require a public repository.

#### REQ: repo-public-or-private

Git storage MUST support both public and private GitHub repositories, except that public topics MUST use a public repository.

### API proxy for performance

Public topics — and any high-traffic git-backed scope — MUST be served through an API layer that caches and indexes the underlying git content for fast reads.

#### REQ: api-proxy-required

Public-topic git storage MUST be served via an API proxy layer rather than direct GitHub reads for every request.

### Write path

Writes flow through the IssueNumber.one API, which produces commits on behalf of the acting user. Commit metadata attributes each change to the correct actor (respecting anonymity rules).

#### REQ: writes-via-api

All writes to git storage MUST flow through the IssueNumber.one API so that validation, anonymity rules, and scheduling can be enforced.

#### REQ: commit-attribution-respects-anonymity

Commit author metadata MUST NOT reveal the identity of an anonymous issue author (see [issue/anonymity#req-anon-author-hidden](../../issue/anonymity/README.md)).

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [storage](../README.md) | One of the two supported backends |
| [organization/topic](../../organization/topic/README.md) | Public topics always use this backend |
| [issue/anonymity](../../issue/anonymity/README.md) | Commit metadata must preserve anonymity |

## Dependencies

- storage
- organization

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- The inGitDB schema for IssueNumber.one entities needs a dedicated specification — where should it live (in this repo, in a separate schema repo, in the inGitDB project)?
- How are anonymous commits implemented such that git history cannot be used to re-identify authors?
- Does the API proxy serve writes too (queued via background job), or only reads?
- How are conflicting concurrent writes resolved — optimistic merges, single-writer locks, or API-serialized writes?
- Acceptance criteria not yet defined for this feature.
