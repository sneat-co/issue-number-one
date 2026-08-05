---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Product scope

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization/product?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization/product?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization/product?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/organization/product?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

Scopes IssueNumber.one to a product while keeping its internal and public issues in separate trust domains.

## Problem

A product can have two legitimate but incompatible priority conversations. Its
makers need a private place to surface delivery, quality, security, commercial,
and organizational problems, while customers and a wider community may need a
public place to surface product problems. Combining those audiences into one
ranking would expose internal information, distort support, and make people
uncertain about who can see what. Treating them as unrelated products would
lose the shared product identity and make its two priority signals difficult to
compare deliberately.

## Behavior

### Product is an issue scope

A product can be selected as the context in which issues are raised, nominated,
supported, ranked, and discussed. The product context has its own current #1
rather than relying only on a department, project, company, or generic public
topic board.

#### REQ: product-has-ranked-issue-scope

IssueNumber.one MUST support a product as an issue-ranking scope with a visible
current #1 and top-N candidate view.

### Internal and public are separate trust domains

The internal product scope is part of an organization and is visible only under
that organization's permissions. The public product scope is a public
participation surface. They may refer to the same product identity, but they do
not share a ranking, issue list, discussion, or participant attribution.

#### REQ: product-has-separated-internal-and-public-domains

A product that has both internal and public participation MUST maintain a
separate issue set, candidate set, ranking, current #1, discussion space, and
access boundary for each trust domain.

#### REQ: internal-product-data-remains-private

The existence, content, discussion, nomination, support, voter identity, and
author identity of an internal product issue MUST NOT become visible in the
public product domain merely because both domains refer to the same product.

#### REQ: no-implicit-cross-domain-transfer

An issue, vote, nomination, discussion entry, or lifecycle event in one product
trust domain MUST NOT be copied, moved, counted, or applied in the other domain
without a separately defined and explicitly authorized cross-domain action.

### Internal issues can be deliberately published as separate public issues

An internal issue can become the source of a public issue when someone with the
required authority deliberately publishes it. Publication creates a new public
record linked to the internal source; it does not change the internal issue's
visibility or place the internal record into the public ranking.

#### REQ: deliberate-publication-creates-linked-public-issue

IssueNumber.one MUST allow an internal product issue to be published as a
separate linked issue in that product's public trust domain only after the
required organization-administrator approval.

### Publication requires organization-administrator approval

Moving selected content across the internal/public boundary is an
organizational disclosure decision. Creating the internal issue does not by
itself give its creator authority to make that disclosure.

#### REQ: organization-admin-approves-publication

A public issue derived from an internal product issue MUST NOT be created until
a current administrator of the owning organization explicitly approves the
publication and the content selected to become public.

#### REQ: creator-cannot-publish-without-admin-approval

The creator of an internal issue MUST NOT be able to publish it into the public
product domain without organization-administrator approval, even when the
creator is permitted to view, edit, nominate, or close the internal issue.

#### REQ: published-issue-has-independent-public-state

The linked public issue MUST have its own identity, content, nomination,
support, rank, discussion, age, and lifecycle state. Subsequent activity on the
public issue MUST NOT automatically alter the internal source issue.

#### REQ: publication-does-not-expose-internal-record

Publication MUST NOT expose the internal issue's record, discussion, votes,
rank, age, author identity, voter identities, or other private metadata. Only
content explicitly selected for the new public issue may cross the boundary.

#### REQ: private-source-link-remains-private

The link from the public issue to its internal source MUST be visible to
authorized internal viewers, but MUST NOT reveal the existence or identity of
the private source to unauthorized public viewers.

### Shared identity does not imply shared access

A relationship between the public and internal product scopes can let an
authorized user understand that both concern the same product. That relationship
is context, not permission, consent, or authorization to combine their data.

#### REQ: product-identity-does-not-grant-access

Linking public and internal scopes to one product identity MUST NOT grant access
to either scope or make their issues eligible in each other's ranking.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [organization](../README.md) | Owns the organizational hierarchy containing an internal product scope |
| [organization/topic](../topic/README.md) | Provides the existing public-scope rules that a public product mode must reconcile with rather than silently replace |
| [issue](../../issue/README.md) | Supplies issue identity, nomination, lifecycle, anonymity, visibility, moderation, and discussion |
| [issue/visibility](../../issue/visibility/README.md) | Enforces the public/internal trust boundary |
| [voting](../../voting/README.md) | Supplies positive support and ranking inside each product trust domain without combining them |
| [permissions](../../permissions/README.md) | Determines who may enter, see, and act in an internal product scope |
| [github-issues](../../github-issues/README.md) | Lets an issue originating in GitHub participate in the selected product trust domain |

## Dependencies

- issue/visibility
- voting
- permissions

## Acceptance Criteria

### AC: product-domains-have-independent-number-ones

- **Given** one product has both an internal scope and a public scope
- **And** each scope contains eligible issues with support
- **When** IssueNumber.one calculates the current #1 for that product
- **Then** it produces one ranking and current #1 for the internal domain and a separate ranking and current #1 for the public domain
- **And** support from one domain does not affect rank in the other

### AC: public-viewer-cannot-discover-internal-product-issue

- **Given** an issue exists only in a product's internal domain
- **And** a public viewer is not authorized to enter that internal domain
- **When** the viewer opens the public product scope or follows the shared product identity
- **Then** the internal issue's existence, content, age, discussion, author, voters, support, and rank are not exposed

### AC: shared-product-identity-does-not-merge-rankings

- **Given** public and internal scopes are linked to the same product identity
- **When** either scope's ranking changes
- **Then** the other scope's issue set, votes, and ranking remain unchanged

### AC: deliberate-publication-creates-safe-public-copy

- **Given** an organization administrator has explicitly approved selected content from an internal product issue for publication
- **When** the approved publication is completed in the product's public domain
- **Then** IssueNumber.one creates a separate linked public issue with its own identity, ranking, votes, discussion, age, and lifecycle
- **And** the internal issue remains private and unchanged
- **And** public viewers cannot discover the internal record, discussion, votes, rank, age, author, voters, or private source link

### AC: creator-cannot-publish-without-admin-approval

- **Given** the creator of an internal issue has not received organization-administrator approval to publish it
- **When** the creator attempts to publish that issue into the public product domain
- **Then** no public issue is created and no internal content or metadata crosses the trust boundary

## Open Questions

- Is a product an organization hierarchy node, a specialized public topic, or a
  product identity linked to one scope of each kind?
- Does an internal product #1 automatically become eligible in a parent
  department or company ranking?
- Is a person's personal-#1 nomination and support-vote budget separate per
  product domain, shared with their organization or public topic, or governed by
  another rule?
- Who may request publication: only the issue creator, any internal viewer, or
  another defined role?
- If the issue creator is also an organization administrator, may that person
  approve their own publication request?
- May one administrator approve publication alone, or can an organization
  require more than one approval?
- Which issue fields are selected for publication, and does later editing ever
  synchronize between the two records?
- Does the public issue's visible age start at publication, preserve the
  internal issue's original age, or show both dates?
- What public authorship should a published issue show when its internal source
  was authored anonymously?
- Who may create a public product scope and link it to an internal product?
- How are public product issues moderated?

---
*This document follows the https://specscore.md/feature-specification*
