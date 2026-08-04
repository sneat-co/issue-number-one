---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Issue Visibility

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/visibility?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/visibility?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/visibility?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/visibility?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Every issue has a visibility level that controls who can see it. A person's full issue list remains their own view, while collective views show eligible candidates and the top N issues at each organizational scope. Only a child scope's current #1 issue automatically becomes a candidate in its parent scope; no manager approval is required.

## Problem

Without explicit visibility rules, teams either over-share every recorded issue or let important problems remain trapped below an unwilling management layer. Visibility in IssueNumber.one keeps a person's full issue list private to its appropriate scope, exposes only collectively eligible candidates, and automatically carries each scope's #1 issue upward.

## Behavior

### Visibility levels

| Level | Who can see it |
|-------|----------------|
| `team` (default) | Creator plus users entitled to see the issue as a candidate or ranked item in that team |
| `org` | Members entitled to see the issue as a candidate or ranked item in an enclosing organizational scope |
| `public` | Anyone, used for public topics and spaces |

#### REQ: default-visibility-team

A newly-raised issue in a team scope MUST default to `team` visibility.

#### REQ: visibility-levels

An issue's `visibility` MUST be one of `team`, `org`, or `public`.

### Automatic #1 bubble-up

Each scope contributes exactly one candidate to its parent: its current #1 issue. A manager does not publish, approve, defer, or block that promotion.

#### REQ: child-top-bubbles-automatically

The current #1 issue of every child scope MUST automatically become visible and vote-eligible in its immediate parent scope.

#### REQ: only-child-top-bubbles

No issue other than a child scope's current #1 MUST become a parent-scope candidate through ranking alone.

#### REQ: bubble-up-needs-no-manager-approval

A manager or scope owner MUST NOT have an approval or veto step in automatic #1 bubble-up.

### Top-N collective views

Within a voting scope, eligible members see the top N ranked issues. The display size is configurable, but N controls the view rather than how many issues move to the parent.

#### REQ: top-n-display-configurable

The number of ranked issues shown in a scope's top list MUST be configurable for the organization.

#### REQ: member-sees-own-and-ancestor-tops

A member MUST be able to see their own recorded issues, the top N issues in their department or immediate organizational scope, and the top N issues at the root company scope. When intermediate project or team scopes exist, the member MUST be able to see the top N for each enclosing scope.

### Drill-down and changing leaders

A leader may drill from the company view through department, project, team, person, and originating issue. When a scope's #1 changes, the former #1 stops being the candidate at the parent but remains open and visibly ageing at its source.

#### REQ: hierarchy-drill-down

An authorized viewer MUST be able to trace a bubbled issue through each child scope to its originating issue without gaining access to unrelated non-candidate issues.

#### REQ: replaced-top-remains-open

When a new issue becomes a scope's #1, the previous #1 MUST lose parent-scope candidate visibility but MUST remain `raised` with its original age unless separately closed under the lifecycle rules.

### Public topics

Public topics are always visible to anyone, logged in or not. Any nested sub-topics inherit public visibility.

#### REQ: public-topic-always-public

Issues in public topics MUST always be `public` visibility; they MUST NOT support `team` or `org` scoping.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../README.md) | The #1-issue rule depends on [issue#req-single-top-issue-per-team](../README.md) |
| [voting](../../voting/README.md) | Scope-specific support scores determine #1 and top-N ranking |
| [organization](../../organization/README.md) | Team/org/topic structure defines the scopes |
| [permissions](../../permissions/README.md) | Visibility is enforced via permissions |

## Dependencies

- issue
- voting
- organization
- permissions

## Acceptance Criteria

Not defined yet.

## Open Questions

- What is the default value of N, and may it vary by organizational scope?
- Should an issue's prior appearances at higher scopes remain visible as immutable ranking history after it stops being the child #1?
- Can a parent scope also contain directly raised issues, or only #1 issues supplied by its children?
- Acceptance criteria not yet defined for this feature.
