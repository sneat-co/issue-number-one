# Storage sheet: revoke paid verification on a full refund

Date: 2026-11-15. **DRAFT for the owner's yes / no. Nothing is built.**

**The rule.** A full refund cancels the service; a partial refund is goodwill and the service stays (founder, 2026-11-14). Today IssueNumber.one cannot cancel it. A `refunded` event is logged, and paid verification stays ([service.go `SettlementHandler`](../../backend/payments/service.go)). The eligibility store offers only `HasPaid` and `MarkPaid`.

**What is stored today** (space-bounded; nothing at the top of the database):

```
spaces/{space}/ext/issuenumber/verification/{uid}          eligible, chargeId, verifiedAt
spaces/{space}/ext/issuenumber/paymentReceipts/{chargeId}  uid, createdAt
```

## Decision: option A (recommended) — reuse `eligible`, add two optional audit fields

On a full refund (EUR, amount ≥ 100 minor units) the handler sets the existing `eligible` to `false` on `verification/{uid}`. It also adds two new optional fields: `revokedAt` (timestamp) and `revokedRefundRef` (the provider's refund reference). Nothing else changes:

- the receipt stays, so a redelivered `paid` event for the same charge cannot grant again;
- a new payment, which is a new charge, grants again through `MarkPaid`, which rewrites the record and clears both fields;
- a partial refund is only logged.

Example, before → after:

```
verification/u_42  {eligible: true,  chargeId: "cs_a1", verifiedAt: 2026-11-10T09:00Z}
verification/u_42  {eligible: false, chargeId: "cs_a1", verifiedAt: 2026-11-10T09:00Z,
                    revokedAt: 2026-11-12T14:05Z, revokedRefundRef: "re_b2"}
```

**How the refund finds the record.** By user, read from the charge's own metadata (`issueNumberUserID`, set by the server and signed). It is not found by charge ID: a refund reports the provider charge (`ch_…`), while the record holds the checkout session (`cs_…`).

**Cost.** One read and one write per full refund, and refunds are rare. No new collection, no index, no client-readable change. Storage grows by about 60 bytes per revoked record.

## Alternatives

| Option | New fields | Lost | Verdict |
|---|---|---|---|
| **A. `eligible=false` + `revokedAt`, `revokedRefundRef`** | 2, optional | nothing | recommended |
| B. `eligible=false` only | none | when it was revoked and by which refund (the operator must search the logs) | acceptable if you want zero shape change |
| C. delete `verification/{uid}` | none | the audit of the original grant | no |

## Open points for the same yes / no

1. **Two payments by one user.** Two checkouts can both settle before the first marks the user paid. Recommendation: any full refund revokes. The user can pay again.
2. **Refunds already logged.** Refunds before this change only left a WARNING log line. Recommendation: no backfill.
3. **Phone verification is untouched.** It is a separate path to eligibility.

On a yes: a DRAFT PR adds `Revoke` to the eligibility store and handles the `refunded` event, with tests first and 100% coverage. Money mechanics stay in Paymentus; IssueNumber.one only reacts to the event.
