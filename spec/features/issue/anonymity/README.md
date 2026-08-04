---
format: https://specscore.md/feature-specification
status: Conceptual
---

# Feature: Issue Anonymity

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/anonymity?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/anonymity?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/anonymity?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue/anonymity?op=request-change) |
**Status:** Conceptual
**Source Ideas:** —

## Summary

Teams and orgs can configure whether issues may be raised anonymously and whether anonymous issues require a small fee to deter artificial boosting. Anonymous issues must still be controllable by their original author — without revealing who that is — and participate in the same personal-#1 nomination and vote-budget rules as authored issues.

## Problem

Some topics only surface when people feel safe raising them. Without anonymity, culturally sensitive, political, or risky issues stay hidden. But pure anonymity invites abuse — sockpuppet accounts and artificial boosting. Anonymity in IssueNumber.one is therefore a team-level opt-in with safeguards: caps, optional fees, and strict author-identity handling.

## Behavior

### Opt-in per team or org

Anonymity is off by default. A team or org MUST explicitly enable anonymous issues before any can be raised.

#### REQ: anonymity-off-by-default

A newly-created team or org MUST NOT allow anonymous issues until explicitly enabled.

### Authored and anonymous issues share one nomination

A person may keep multiple authored or anonymous issues, but may nominate at most one personal #1 across both types in the same scope.

#### REQ: anonymity-shares-personal-top-nomination

Authored and anonymous issues MUST share the same one-personal-#1 nomination limit for their creator in a scope.

### Optional fee for anonymous issues

To prevent artificial boosting, a team MAY require payment of a small fee before an anonymous issue is accepted. The fee mode MAY be:

- `per-card` — one flat fee per anonymous issue raised
- `per-vote` — one unit of fee per vote the issue receives (e.g., 1 vote = $1)

#### REQ: anon-fee-optional

An anonymous-issue fee MUST be optional per team/org. When enabled, the fee mode MUST be recorded.

### Author hiding with controllable identity

An anonymous issue MUST hide the author's identity from all other users but MUST still allow the original author to update, withdraw, or resolve their own anonymous issue.

#### REQ: anon-author-hidden

The author of an anonymous issue MUST NOT be visible to any other user, including team admins.

#### REQ: anon-author-can-modify

The system MUST allow the original anonymous author to update, nominate, withdraw, or resolve their issue without revealing their identity.

### Anonymous self-support

The original author may allocate support votes to their own anonymous nominated issue under the same budget as any other voter. The system enforces the budget and nomination rule without exposing the author relationship.

#### REQ: anonymous-author-may-self-support

An anonymous author MUST be able to allocate available support votes to their own nominated issue without revealing their identity to other users.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../README.md) | Anonymity is a property of an issue; the personal-#1 nomination rule still applies |
| [voting](../../voting/README.md) | Vote-budget enforcement must work without exposing an anonymous author |
| [permissions](../../permissions/README.md) | Enabling anonymity is an org/team admin action |

## Dependencies

- issue
- voting
- permissions

## Acceptance Criteria

Not defined yet.

## Open Questions

- How is author control represented so the system can enforce author-only update/withdraw/resolve and budget rules while being unable to reveal identity to admins or AI analysis?
- Should anonymous-issue fees be refundable when the issue is withdrawn or resolved?
- For the `per-vote` fee mode, when is the fee charged — per vote received, or capped at a maximum?
- Does AI analysis of issues need special handling for anonymous authors to avoid re-identification through writing style?
- Acceptance criteria not yet defined for this feature.
