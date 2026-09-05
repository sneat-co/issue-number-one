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

Each verification attempt has a browser-generated opaque `actionId`, retained
across retries. The server combines it with the trusted Firebase user ID for the
provider idempotency key. A cancelled or expired attempt can start again with a
new ID; replaying one attempt returns the same provider Checkout Session. Only
bounded `categoryId`, `questionId`, and `actionId` identifiers enter metadata;
free-form issue text must never be sent.
