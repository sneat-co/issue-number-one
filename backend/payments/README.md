# Payment verification operations

IssueNumber.one uses the shared paymentus one-off card rail. The product package
owns the fixed EUR 1 participation-verification purchase and settlement checks;
sneat-go owns provider credentials, Checkout Session creation, webhook signature
verification, and the single settlement router.

The checkout description is `Support IssueNumber.one and verify participation`.
It describes both purposes without claiming that payment proves one human has
only one account. Phone verification remains an equal alternative.

Before enabling real payments, the operator must complete the existing
sneat-go/payments Stripe test-to-live setup, configure the signed webhook modes,
and confirm the merchant's accounting, tax, receipt, chargeback, and refund
policy with the responsible owner. This package does not invent those policies
or create a second payment authority. Automated tests use paymentus test doubles
and never create a provider charge.

The public eligibility path is live-mode only. sneat-go stamps the selected mode
into provider-owned metadata and this consumer rejects a signed test-mode event,
so a Stripe test Checkout cannot make an answer count in production. Emulator
tests inject paymentus' mock provider directly and remain incapable of moving
money.

Each verification attempt has a browser-generated opaque `actionId`, retained
across retries. The server combines it with the trusted Firebase user ID for the
provider idempotency key. A cancelled or expired attempt can start again with a
new ID; replaying one attempt returns the same provider Checkout Session. Only
bounded `categoryId`, `questionId`, and `actionId` identifiers enter metadata;
free-form issue text must never be sent.

## Settlement events

The host answers the provider with a retry for any error the settlement handler
returns, and the payments webhook is shared with other products. So the handler
returns an error only for an eligibility-store failure. That covers transient
storage errors, and also two permanent store refusals (a charge ID the store
cannot key, a receipt bound to another user) that provider charge IDs and
server-written metadata cannot produce. Every other event is acknowledged,
grants nothing, and is logged through `logus`, which the host ships to Cloud
Logging with severity:

- a valid `paid` event grants eligibility once per provider charge ID;
- a `paid` event that fails a check (consumer, EUR 1.00 amount, currency, live
  mode, user and payer metadata) logs `paid settlement NOT applied, operator
  review required` at ERROR with the reason and charge reference. Fix the
  cause, then resend the event from the provider dashboard (the handler is
  idempotent by charge ID) or refund the payer. A log-based alert on this
  message is an operations setting for the owner;
- a `refunded` event logs at WARNING and leaves eligibility unchanged. Whether a
  refund should revoke paid eligibility is an open product decision;
- a `failed` event logs at INFO; an unknown kind logs at WARNING.

Logs carry the provider charge and refund references, never user identifiers.
