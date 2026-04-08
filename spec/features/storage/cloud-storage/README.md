# Feature: Cloud Storage

**Status:** Conceptual

## Summary

The default IssueNumber.one storage backend — a managed cloud service implemented on Firestore. Orgs that choose cloud storage get a zero-configuration experience: no repositories to create, no tokens to configure, no hosting to maintain.

## Problem

Most teams adopting a new tool don't want to provision infrastructure before they can try it. A managed backend removes every setup obstacle. At the same time, the implementation must not leak into the product contract — we should be able to swap Firestore for another database in the future without rewriting feature specs.

## Behavior

### Default choice

New orgs default to cloud storage unless the user explicitly selects git-storage.

#### REQ: cloud-is-default

When an organization is created without specifying a backend, IssueNumber.one MUST use cloud storage.

### Implementation

The current cloud implementation uses Firestore. This is an implementation detail and MAY change.

#### REQ: firestore-current-implementation

The current implementation of cloud storage MUST use Firestore, but the product contract MUST NOT depend on Firestore-specific features.

### Data locality

Data in cloud storage lives in IssueNumber.one-controlled Firestore projects. Teams do not have direct database access; they interact via the IssueNumber.one API.

#### REQ: no-direct-database-access

Cloud-storage users MUST NOT have direct access to the underlying Firestore instance. All access MUST go through the IssueNumber.one API.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [storage](../README.md) | Cloud is the default of the two supported backends |
| [organization](../../organization/README.md) | Orgs choose this backend at creation time |

## Dependencies

- storage
- organization

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- What regions are supported, and can a team pin their data to a specific region?
- What is the data-retention policy for archived issues, teams, and orgs?
- Is there an export-to-git tool for users who want to migrate away from cloud storage?
- Acceptance criteria not yet defined for this feature.
