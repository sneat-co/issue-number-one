# Feature: Issue Anonymity

**Status:** Conceptual

## Summary

Teams and orgs can configure whether issues may be raised anonymously, how many anonymous issues each member may have, and whether anonymous issues require a small fee to deter artificial boosting. Anonymous issues must still be updatable and withdrawable by their original author — without revealing who that is.

## Problem

Some topics only surface when people feel safe raising them. Without anonymity, culturally sensitive, political, or risky issues stay hidden. But pure anonymity invites abuse — sockpuppet accounts and artificial boosting. Anonymity in IssueNumber.one is therefore a team-level opt-in with safeguards: caps, optional fees, and strict author-identity handling.

## Behavior

### Opt-in per team or org

Anonymity is off by default. A team or org MUST explicitly enable anonymous issues before any can be raised.

#### REQ: anonymity-off-by-default

A newly-created team or org MUST NOT allow anonymous issues until explicitly enabled.

### Separate caps for personal and anonymous

A team may allow a mix of personal and anonymous issues per member — for example, 1 personal + 1 anonymous. The per-person active-issue cap from [issue](../README.md) applies to the combined total.

#### REQ: separate-anonymous-cap

Teams MAY configure independent caps for personal (`X`) and anonymous (`Y`) issues per member, subject to the overall cap of three.

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

The system MUST allow the original anonymous author to update or withdraw their issue without revealing their identity.

### Self-voting still prohibited

Even though the author is hidden, anonymous issues MUST enforce the global no-self-vote rule from [voting](../../voting/README.md).

#### REQ: anon-no-self-voting

An anonymous author MUST NOT be able to vote on their own anonymous issue.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../README.md) | Anonymity is a property of an issue; active caps still apply |
| [voting](../../voting/README.md) | Self-vote prevention must still work despite hidden author |
| [permissions](../../permissions/README.md) | Enabling anonymity is an org/team admin action |

## Dependencies

- issue
- voting
- permissions

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- How is author identity stored such that the system can enforce author-only update/withdraw and no-self-voting while being unable to reveal identity to anyone (including admins or AI analysis)?
- Should anonymous-issue fees be refundable when the issue is withdrawn or resolved?
- For the `per-vote` fee mode, when is the fee charged — per vote received, or capped at a maximum?
- Does AI analysis of issues need special handling for anonymous authors to avoid re-identification through writing style?
- Acceptance criteria not yet defined for this feature.
