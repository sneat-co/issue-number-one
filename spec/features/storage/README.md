# Feature: Storage

**Status:** Conceptual

## Summary

Teams and orgs choose where their IssueNumber.one data is stored: the IssueNumber.one cloud (Firestore-backed by default) or a GitHub repository via [inGitDB](https://inGitDB.com). Public topics always use a public GitHub repository, proxied by an API layer for performance.

## Contents

| Directory | Description |
|-----------|-------------|
| [cloud-storage/](cloud-storage/README.md) | The default IssueNumber.one cloud backend, implemented on Firestore |
| [git-storage/](git-storage/README.md) | The GitHub-backed alternative using inGitDB's structured, human-readable file layout |

### cloud-storage

The turnkey option. Teams create an org in IssueNumber.one and their data is persisted in a managed Firestore instance. Low-friction, zero-configuration, zero-operations.

### git-storage

The user-owned option. Teams connect a GitHub repository and their data is persisted as strongly-structured, human-readable files using the inGitDB schema. Git history becomes the audit log; repository permissions govern access.

## Problem

Different teams trust different infrastructure. Some want a managed cloud service with no setup; others want their data to live in their own GitHub org, auditable, and versioned alongside their code. Forcing a single backend alienates one group or the other. IssueNumber.one supports both through a storage-choice feature that keeps the product contract identical regardless of backend.

## Behavior

### Per-org storage choice

The storage backend is chosen at org creation. It MAY be migratable later, but this is not a base-product guarantee.

#### REQ: storage-chosen-at-creation

Every organization MUST have its storage backend chosen at creation time.

#### REQ: supported-backends

The supported backends MUST be `cloud` (default) and `git`.

### Public topics are always git

Public topics always live in a public GitHub repository. They are never stored in the cloud backend.

#### REQ: public-topics-always-git

Public topics MUST always use the git-storage backend with a public GitHub repository.

### Product contract is backend-independent

Every feature MUST behave identically regardless of storage backend. Backend-specific behavior (e.g., commit author attribution) MUST NOT leak into feature specs.

#### REQ: backend-independent-contract

All user-visible behavior MUST be identical across supported storage backends.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [organization](../organization/README.md) | Orgs choose their backend |
| [organization/topic](../organization/topic/README.md) | Topics always use git-storage |
| [issue/anonymity](../issue/anonymity/README.md) | Backend must support hidden-author requirements |

## Dependencies

- organization

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Can an org migrate from cloud to git (or vice versa) after creation?
- Is there a hybrid mode — e.g., metadata in cloud, long-form content in git?
- Needs a comparison doc of cloud vs git pros and cons (to live in `docs/`).
- Acceptance criteria not yet defined for this feature.
