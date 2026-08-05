---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: GitHub Issues integration

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/github-issues?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/github-issues?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/github-issues?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/github-issues?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

Lets GitHub Issues participate in IssueNumber.one nomination, scarce voting, and bubble-up without replacing GitHub as the issue system of record.

## Problem

Software teams and public product communities already record actionable work in
GitHub Issues. Requiring them to recreate or migrate those issues into
IssueNumber.one would introduce duplicate records and competing status
workflows. GitHub's backlog, labels, reactions, and sort orders do not, however,
provide IssueNumber.one's scarce personal nomination, support budget, automatic
hierarchical bubble-up, or persistent organizational priority signal.

## Behavior

### A GitHub Issue can enter an IssueNumber.one scope

An issue created in GitHub can be linked into an eligible IssueNumber.one scope
and participate in the same nomination, voting, ranking, and bubble-up mechanism
as an issue first recorded in IssueNumber.one.

#### REQ: externally-created-github-issue-can-participate

An authorized person MUST be able to link an existing GitHub Issue into an
eligible IssueNumber.one scope without manually recreating the GitHub issue as
a separate descriptive source record.

#### REQ: linked-github-issue-uses-normal-priority-mechanism

A linked GitHub Issue MUST be eligible for the same personal-#1 nomination,
positive budgeted support, ranking, age display, and automatic bubble-up rules
as other issues in its IssueNumber.one trust domain.

### GitHub identity is preserved

GitHub may remain the system of record for the work item. IssueNumber.one adds
the organizational priority signal while retaining a durable reference to the
source repository and issue number.

#### REQ: preserve-github-issue-identity

A linked issue MUST retain its GitHub repository identity, issue number, and
canonical GitHub URL so an authorized viewer can open the source work item.

#### REQ: priority-overlay-does-not-require-migration

Using nomination, votes, ranking, or bubble-up for a GitHub Issue MUST NOT
require moving the source issue out of GitHub or replacing its GitHub identity.

### IssueNumber.one owns the bubble-up signal

GitHub reactions, comments, labels, assignees, and milestones can provide useful
context but are not automatically equivalent to IssueNumber.one nomination or
scarce support votes.

#### REQ: github-activity-does-not-create-priority-support

GitHub reactions, comments, labels, assignment, or other source activity MUST
NOT create or remove an IssueNumber.one creator star, budgeted vote, nomination,
or rank unless a future explicit mapping rule defines that behavior.

### The selected trust domain remains authoritative

Linking a GitHub Issue does not make private IssueNumber.one votes, discussion,
or organizational context public. Likewise, a private GitHub repository must
not leak source metadata or content through a public product surface.

#### REQ: github-link-respects-both-access-boundaries

The integration MUST require a viewer to have the applicable IssueNumber.one
access and any GitHub access needed for non-public source data. It MUST NOT
expose private-repository content, internal rankings, internal discussion, voter
identity, or author identity across a public trust boundary.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../issue/README.md) | Owns IssueNumber.one nomination, score, rank, age, and lifecycle state associated with the source work item |
| [voting](../voting/README.md) | Supplies scarce positive support rather than deriving votes from GitHub reactions |
| [organization/product](../organization/product/README.md) | Selects the internal or public product trust domain in which the linked issue participates |
| [permissions](../permissions/README.md) | Determines who may link and view an issue in IssueNumber.one |
| [storage/git-storage](../storage/git-storage/README.md) | Is a separate capability for storing IssueNumber.one data in a repository; it does not mean that GitHub Issues themselves are the data store |

## Dependencies

- issue
- voting
- organization/product
- permissions

## Acceptance Criteria

### AC: github-issue-bubbles-without-duplication

- **Given** an authorized person links an existing GitHub Issue to an eligible IssueNumber.one product scope
- **When** that issue is nominated, receives support, and becomes the scope's #1
- **Then** it participates in the normal parent bubble-up path
- **And** its repository, issue number, and canonical GitHub URL remain available
- **And** no second GitHub Issue is required

### AC: github-reaction-is-not-an-issue-number-one-vote

- **Given** a linked GitHub Issue has an established IssueNumber.one support score
- **When** someone adds or removes a GitHub reaction, comment, label, assignee, or milestone
- **Then** that activity alone does not change the IssueNumber.one nomination, vote total, or rank

### AC: private-github-content-does-not-leak-publicly

- **Given** a linked source issue is in a private GitHub repository
- **And** a viewer lacks access to that repository
- **When** the viewer opens a public IssueNumber.one product scope
- **Then** private source content and metadata are not exposed by the integration

## Open Questions

- Which GitHub fields are read into IssueNumber.one, and which system owns each
  field when values differ?
- Is synchronization one-way from GitHub, bidirectional, or selectable per
  product?
- Should closing or reopening a GitHub Issue change its IssueNumber.one
  lifecycle state, and should either product ever close the other record?
- May a person create a new GitHub Issue from IssueNumber.one as well as link an
  issue already created in GitHub?
- Who may link a GitHub Issue, and how is a duplicate link detected?
- What happens when a GitHub Issue is transferred, converted, deleted, or its
  repository visibility changes?
- Should selected public GitHub fields be cached for viewers without GitHub
  accounts, and how is staleness shown?
- What GitHub authentication and installation scope is required?

---
*This document follows the https://specscore.md/feature-specification*
