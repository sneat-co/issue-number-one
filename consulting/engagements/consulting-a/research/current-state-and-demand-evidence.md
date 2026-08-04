---
title: "IssueNumber.one Current-State and Demand-Evidence Baseline"
artifact: research
status: Draft
engagement: "consulting-a"
author: "Codex (Managing Partner)"
model: "GPT-5 (Codex)"
created: "2026-08-04"
updated: "2026-08-04"
confidence: "Medium for repository and live-surface state; low for demand because no private aggregate, customer or usage evidence was supplied"
dependencies:
  - "../README.md"
  - "../checkpoints/iteration-001-bootstrap.md"
related_specs:
  - "../../../../spec/features/README.md"
  - "https://github.com/sneat-co/backstage/blob/5f29c3424de8ffda0dfe7f7ac09c510741043381/spec/features/ai-consulting-framework/README.md"
promotion_target: "None — evidence remains provisional until reviewed"
---
# IssueNumber.one Current-State and Demand-Evidence Baseline

## Owned Question

What product, demand-probe and public-discoverability evidence exists today,
what does it actually establish, and which missing evidence most affects the
next research decisions?

## Concise Answer

IssueNumber.one has a live Cloudflare-served landing page, a working waitlist
write path in source, explicit pricing/segment hypotheses, an owned product
repository and a broad Conceptual Feature tree. The authenticated application
remains a scaffold rather than an observable product experience.

The evidence establishes that a **public demand-probe instrument exists**. It
does not establish demand, qualified interest, conversion, active use, repeat
use, willingness to pay, revenue or a viable buyer. No raw waitlist data was
accessed, no non-identifying aggregate was supplied, and the public surface does
not expose conversion analytics. Targeted public searches also found no
relevant indexed third-party product profile or review, but search absence is
weak evidence and must not be turned into a claim of zero awareness.

The highest-leverage next input remains the approved non-identifying evidence
inventory: unique waitlist signups by month, aggregate source/plan counts,
landing traffic/conversion if available, prior problem conversations and any
active/repeated-use or budget-owner evidence.

## Method and Scope

Research performed on 2026-08-04:

1. Read the current repository Feature tree, landing, worker, application shell,
   routes, package metadata, screen documents and recent Git history at product
   commit `570ef67fd1a84274de6f7b35f68ef9a8e2b26793`.
2. Read the source-linked Backstage seed, Sneat Work Idea, ecosystem review,
   relationship-graph and Moodometer context at Backstage commit
   `5f29c3424de8ffda0dfe7f7ac09c510741043381`.
3. Fetched [https://issuenumber.one](https://issuenumber.one) with a bounded,
   read-only HTTP request. It returned HTTP/2 200 through Cloudflare at
   2026-08-04 11:18 UTC and served the current Astro landing content.
4. Searched the public web for `IssueNumber.one product`,
   `IssuerNumber.One platform` and `site:issuenumber.one`, then targeted
   LinkedIn, Product Hunt, G2 and Capterra for `issuenumber.one`.

Exclusions:

- No waitlist submission was made.
- No KV namespace, raw email, private analytics, customer record, billing
  system or private conversation was accessed.
- No person was contacted and no public claim was published.
- This report does not evaluate competitor categories in depth; that belongs to
  `customer-market-evidence.md`.

## Evidence

| ID | Classification | Claim or observation | Source and date | Freshness / limitations |
|---|---|---|---|---|
| E1 | Observed | The public root returned HTTP/2 200 through Cloudflare and served an Astro page | [Live landing](https://issuenumber.one), fetched 2026-08-04 11:18 UTC | One successful point-in-time availability probe; not uptime evidence |
| E2 | Observed | The served page identifies the product as a focused feedback tool, centres scarce issues/votes and anonymity, presents five launch-price tiers, links it to `sneat.team`, offers a waitlist and says public boards are “coming later” | [Live landing](https://issuenumber.one), accessed 2026-08-04 | Product-owned claims and intent, not customer or commercial evidence |
| E3 | Observed | The landing source sends `{email, source, plan}` to `POST /waitlist`; the Worker validates email and stores a record keyed by normalized email when KV is bound | [Landing form](../../../../landings/src/pages/index.astro) and [worker](../../../../landings/worker.js), product commit `570ef67` | Establishes the intended write mechanism, not production record count or deliverability |
| E4 | Observed | The served HTML and layout source contain a TODO to add Google Analytics after a GA4 property exists | [Layout source](../../../../landings/src/layouts/BaseLayout.astro) and [live landing](https://issuenumber.one), accessed 2026-08-04 | Does not rule out Cloudflare analytics, logs or another private measurement system |
| E5 | Observed | The authenticated app renders the Nx welcome component and declares no routes | [App shell](../../../../src/app/app.html) and [routes](../../../../src/app/app.routes.ts), product commit `570ef67` | Code-state evidence; does not prove no separate unpublished prototype exists |
| E6 | Observed | The product Feature index and its issue, voting, organisation, permissions, storage and AI children are marked Conceptual, with acceptance criteria generally not defined | [Feature index](../../../../spec/features/README.md) and linked Features, product commit `570ef67` | Specification maturity, not implementation or market maturity |
| E7 | Sourced | The Backstage seed records founder conviction, buyer, pricing and public-board viral-loop hypotheses and says the landing became live | [Backstage seed](https://github.com/sneat-co/backstage/blob/5f29c3424de8ffda0dfe7f7ac09c510741043381/spec/ideas/seeds/issue-number-one.md), captured 2026-07-09 | Founder-supplied hypotheses; no independent demand evidence |
| E8 | Sourced with contradiction | The July ecosystem catalogue reported no public usage or revenue evidence across the ecosystem and then-stale hosting status for IssueNumber.one | [Ecosystem catalogue](https://github.com/sneat-co/backstage/blob/5f29c3424de8ffda0dfe7f7ac09c510741043381/spec/research/ecosystem-review-2026-07/02-product-catalogue.md), dated 2026-07-09 | Hosting statement is superseded by E1 and the later seed; usage/revenue absence has not been re-audited privately |
| E9 | Observed failed search | General and targeted searches returned irrelevant similarly named products or no relevant IssueNumber.one profile on LinkedIn, Product Hunt, G2 or Capterra | Search queries listed in Method, run 2026-08-04 | Search indexing is incomplete and noisy; this supports only “not found in these searches,” not “does not exist” |
| E10 | Unknown | Unique waitlist records, traffic, conversion, qualified conversations, active users, repeat behaviour, purchase intent and revenue were not supplied or inspected | [Accepted checkpoint](../checkpoints/iteration-001-bootstrap.md) | Absence from the engagement evidence base is not proof that the data is zero |

## Findings

### Observed or sourced

- **A demand-probe door exists and is currently reachable** (E1–E3).
- **The probe tests the older scarce-priority and `sneat.team` suite story, not
  the sponsor's new relationship-aware issue-intelligence framing** (E2, E7).
- **The probe collects an email and optional plan selection, but the engagement
  has no result data** (E3, E10).
- **The public page includes explicit prices but no checkout**; its own copy says
  checkout is not open and waitlist selection locks in launch pricing (E2).
- **The application repository is not evidence of an operable product journey**
  because the current app shell exposes no product routes (E5).
- **The public-search footprint found in this bounded review is essentially the
  owned site/repository, not independently visible adoption or review evidence**
  (E9, with its stated search limitation).

### Inferred

- **The current asset is better described as a live positioning/pricing probe
  backed by conceptual product specifications than as an early product with
  evidenced traction.** This follows from E1–E6 and E10; it is not a claim that
  no private traction exists.
- **Plan selection can provide directional message/price interest but cannot by
  itself establish willingness to pay.** The visitor incurs no payment and the
  page offers no checkout (E2–E3).
- **Any current waitlist aggregate would primarily evaluate the older landing
  proposition.** It should not be attributed automatically to relationship-aware
  intelligence, public research or benchmarking, which the served page does not
  offer (E2).
- **The `sneat.team` positioning mismatch is measurement-relevant:** future
  traffic or signups could respond to a suite bundle that the newer Sneat Work
  source no longer treats as a standalone product boundary. This inference
  combines E2 with the [Sneat Work Idea](https://github.com/sneat-co/backstage/blob/5f29c3424de8ffda0dfe7f7ac09c510741043381/spec/ideas/sneat-work.md).

## Contradictions and Disconfirming Evidence

- The July ecosystem catalogue called the product unhosted; the current live
  probe disproves that hosting claim. The catalogue must still be treated as a
  point-in-time source for its separate usage-evidence observation, not either
  accepted or discarded wholesale.
- The repository contains detailed Feature documents and a live-looking board
  mock, which can create an impression of product maturity; the application
  shell and absent routes disconfirm an inference that those documents and
  visuals correspond to a current implemented journey.
- The landing publishes price points, but the absence of checkout and missing
  conversion/purchase data disconfirm treating those prices as validated.
- Public search did not find independent profiles or reviews, but search
  incompleteness disconfirms treating “not found” as proof of no awareness.

## Assumptions and Unknowns

- **Assumption:** The served landing is the only currently public
  IssueNumber.one product surface. If wrong, discoverability and maturity are
  understated. Trigger: another URL, store listing or integration is supplied.
- **Assumption:** No analytics visible in served markup means only that no
  client-side GA tag is visible. It says nothing about Cloudflare or private
  server-side measurement.
- **Unknown:** Whether the production KV binding contains zero, one or many
  unique records.
- **Unknown:** Whether any signup maps to a qualified target buyer, problem
  conversation, organisation or selected paid plan.
- **Unknown:** Whether any person has used a product prototype repeatedly or
  offered budget.
- **Unknown:** Whether landing traffic was intentionally distributed; a low
  signup count is uninterpretable without exposure/source data.
- **Unknown:** Whether the older suite/price positioning should be treated as an
  active probe or a stale instrument; deciding that belongs after the broader
  research, not in this baseline.

## Confidence

**Medium** for the current repository and public-surface description because it
is directly observed and source-linked. **Low** for any demand conclusion
because the decisive private aggregate, exposure denominator, customer context
and behavioural/payment evidence were not supplied.

The evidence most likely to change this answer is a privacy-preserving dataset
showing traffic by source, unique signups, selected plans and qualified follow-up
outcomes, plus evidence of repeated use or an actual budget-holder commitment.

## Risks, Dependencies and Next Research

- Do not collapse `waitlist exists` into `demand exists`.
- Do not collapse `explicit prices exist` into `willingness to pay exists`.
- Do not attribute interest in the old scarce-priority/suite story to the new
  relationship-aware or benchmark hypotheses.
- Preserve traffic source and exposure denominator if an aggregate is supplied;
  raw signup count alone cannot distinguish positioning from distribution.
- The customer/market specialist should use this as the maturity baseline and
  avoid vendor comparisons that imply IssueNumber.one already operates at the
  same evidence level.
- The growth specialist should treat current tiers as hypotheses and the absent
  public measurement tag as a measurement gap, not evidence of zero traffic.

## Alternatives for the Next Evidence Step

| Option | Information gained | Limitation / risk | Current status |
|---|---|---|---|
| Status quo: continue public research while saying “traction unknown” | Preserves safety and prevents invented demand | Cannot prioritize segments from existing behaviour | Authorized and in progress |
| Sponsor supplies aggregate waitlist/traffic/plan counts with no identities | Establishes exposure and directional interest in the old proposition | Still not willingness to pay; may be too small or distribution-biased | Highest-priority missing input; not supplied |
| Review de-identified prior problem conversations | Reveals language, actors and repeated pain | Retrospective selection bias; governance required for notes | Not authorized/supplied |
| Conduct new problem interviews | Tests urgent job, status quo and buyer context | Outreach authorization and sampling plan required | Not authorized in this iteration |
| Build or change the product/landing | Produces activity, not necessarily decision-changing evidence | Premature implementation; confounds message and mechanism | Explicitly out of scope |

## Recommended Research Priority

1. Accept a non-identifying aggregate evidence inventory if the sponsor supplies
   it under the approved boundary.
2. Otherwise keep traction explicitly unknown and use the parallel public
   research to identify which problem and buyer deserve later primary research.
3. Do not redesign the landing or instrumentation until the research clarifies
   which proposition and decision the instrument should test.

This is a recommendation about **research order**, not a business, product or
go-to-market recommendation.

## Proposed Experiments

None. A data inventory and source review are smaller than an experiment and can
change the next research decision without implementation or external action.

## Dissent or Conflicts Requiring Managing Partner Resolution

- Founder conviction and published prices support taking the concept seriously,
  but they remain sponsor evidence and pricing intent. They must not be promoted
  into market demand or willingness-to-pay facts.
- The live landing proves that the product has a public door; it does not resolve
  whether that door tests the correct product framing. Later synthesis must keep
  “instrument exists” separate from “instrument measures the new hypothesis.”
