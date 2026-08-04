---
title: "IssueNumber.one Product Modes, Trust and Relationship-Context Research"
artifact: research
status: Draft
engagement: "consulting-a"
author: "Chief Product Officer specialist"
model: "gpt-5.6-sol"
created: "2026-08-04"
updated: "2026-08-04"
confidence: "Medium-low — strong standards and product-practice evidence define the trust and research-method boundaries, but no IssueNumber.one user research or comparative product outcome exists"
dependencies:
  - "../README.md"
  - "../checkpoints/iteration-001-bootstrap.md"
  - "../../../../spec/features/README.md"
  - "sneat-co/backstage@5f29c3424de8ffda0dfe7f7ac09c510741043381"
related_specs:
  - "../../../../spec/features/issue/README.md"
  - "../../../../spec/features/issue/anonymity/README.md"
  - "../../../../spec/features/issue/visibility/README.md"
  - "../../../../spec/features/organization/README.md"
  - "../../../../spec/features/organization/topic/README.md"
  - "../../../../spec/features/permissions/README.md"
  - "../../../../spec/features/voting/README.md"
  - "../../../../spec/features/ai-integration/README.md"
promotion_target: "None while Draft; possible future owners are the relevant IssueNumber.one Features and reviewed Sneat platform relationship-graph decisions"
---
# IssueNumber.one Product Modes, Trust and Relationship-Context Research

## Owned Question

How do the proposed public-research and internal-organisation-intelligence
modes differ in actors, jobs, incentives, journeys, trust, anonymity, consent,
data rights, actionability and harm; what relationship context could improve
value compared with a non-graph baseline; and which signal classes should
remain isolated?

The report evaluates H1, H2, H3, H4, H6, H8 and H10 from the approved
[engagement brief](../README.md#working-hypotheses) as hypotheses. It does not
choose a product strategy, product boundary or technical design.

## Concise Answer

The two proposed modes are not the same job with different visibility.

- **Public research** asks a participant body to inform a public, community or
  policy decision. Its credibility depends on a declared target population,
  recruitment and sampling method, accessible participation, question
  neutrality, transparent analysis and a public account of what happened to
  the input. Open participation can legitimately surface issues, but raw vote
  totals from a self-selected public board are not population estimates.
- **Internal organisation intelligence** asks people in an employment power
  relationship to surface concerns or patterns to actors who may control their
  work, evaluation and livelihood. Its value depends on psychological safety,
  a specific and credible purpose, separation from individual performance
  management, confidentiality that survives filtering and free text, and a
  visible follow-up process. Employee consent alone is usually not a sound
  trust basis because it may not be freely given.

Relationship context can plausibly add value where it answers a bounded
decision question: which affected cohort is missing; which organisational unit
owns an issue; where a problem crosses teams; whether materially different
groups agree; or which accountable actor must close the loop. That is a narrower
claim than “the graph makes the insight better.” The available evidence does
not show that IssueNumber.one relationship context improves decision quality,
adoption or willingness to pay over a non-graph workflow.

Several signal classes have strong reasons to remain inside their originating
trust domain unless future evidence establishes a specific, disclosed and
lawful use: anonymous-author control mappings, individual mood or wellbeing
pulses, free-text employee comments, moderation and retaliation reports,
respondent contact/recruitment records, public demographic microdata, inferred
opinion or health labels, and organisation-level benchmark contributions.
Aggregates do not automatically make these safe: small or overlapping groups,
rare relationship attributes and longitudinal releases can expose individuals
by singling out, linkage, inference or differencing.

The strongest current evidence therefore supports neither one shared product
nor two separate products. It supports treating the modes as distinct trust and
research propositions until user evidence, comparative outcomes and data-use
comprehension tests discriminate among separate products, tightly separated
modes, a single initial mode, the current scarce-priority product, or rejection
of both intelligence directions.

## Method and Scope

Research was conducted on 2026-08-04 from repository-visible evidence and
current public sources. The method was:

1. read the approved engagement brief and latest checkpoint;
2. inspect the complete current IssueNumber.one Feature tree;
3. inspect only decision-relevant Backstage sources linked from the brief:
   the IssueNumber.one seed, Sneat Work Idea, Contactius relationship-graph
   Feature, core-modules research, Moodometer sources and the explicitly
   opinionated graph-growth notes;
4. review primary standards, regulators, original research and official product
   documentation on survey research, public participation, worker monitoring,
   employee voice, confidentiality, anonymisation and organisational insight;
5. compare each graph-context claim against a non-graph baseline and identify
   its beneficiary, permission domain, possible consumer and harm mechanism.

The review did not inspect private waitlist records, employee data, customer
interviews, analytics or unpublished product behaviour. It did not design a
domain model, API, storage architecture or integration. Legal and regulatory
sources establish product-risk questions, not legal advice for a particular
deployment or jurisdiction.

This report covers product-facing graph value and harm. It does not establish
market defensibility, switching costs, cold-start economics, graph
replicability or portability, or measured network effects across participants
or organisations; those require separate commercial and strategic synthesis.

### Provenance and freshness

| Source class | Sources used | Access/freshness | Limitation |
|---|---|---|---|
| Product canon | Current IssueNumber.one Features under `spec/features/` | Checkout on 2026-08-04; all Features are Conceptual | No acceptance criteria or implemented product journey |
| Engagement context | [Approved brief](../README.md) and [iteration 001 checkpoint](../checkpoints/iteration-001-bootstrap.md) | Accepted 2026-08-04 | Sponsor framing and hypotheses, not market evidence |
| Sneat platform ecosystem | Backstage sources pinned by the brief to commit `5f29c342` | Mostly dated June–July 2026 | Several are Ideas or opinionated strategy, not validated user evidence |
| Data protection | EU GDPR text, EDPB and ICO guidance | Accessed 2026-08-04 | ICO monitoring/anonymisation pages state that UK guidance is under review after the Data (Use and Access) Act; applicability varies by jurisdiction |
| Survey/public participation | AAPOR 2026 Code, AAPOR best practices, OECD citizen-participation guidance, UK consultation principles | Accessed 2026-08-04 | These define quality practice; they do not validate this product or every type of community board |
| Employee voice | Peer-reviewed original/review research from 1999–2021 | Stable publications; accessed 2026-08-04 | Mostly observational or review evidence; not an IssueNumber.one trial |
| Product practice | Microsoft Viva Insights, Culture Amp, Qualtrics and Decidim official documentation | Current pages accessed 2026-08-04 | Vendor controls are examples, not proof that a particular threshold or product pattern is sufficient |

## Evidence

| ID | Classification | Claim or observation | Source and date | Freshness / limitations |
|---|---|---|---|---|
| E1 | Observed | Current IssueNumber.one canon defines scarce issues and votes, team/org/public scopes, opt-in anonymity, membership-derived permissions and optional AI summaries; it does not define research sampling, benchmarks, organisation-intelligence interpretation or relationship-aware analysis. | Current [Feature index](../../../../spec/features/README.md), [Issue](../../../../spec/features/issue/README.md), [Voting](../../../../spec/features/voting/README.md), [Anonymity](../../../../spec/features/issue/anonymity/README.md), [Visibility](../../../../spec/features/issue/visibility/README.md) and [AI Integration](../../../../spec/features/ai-integration/README.md), inspected 2026-08-04 | Conceptual; no acceptance criteria or implementation evidence |
| E2 | Observed | Public topics are always public, open to authenticated contributors and stored in a public GitHub repository; their sampling, moderation authority, anonymity and voting budget remain open. | Current [Topic](../../../../spec/features/organization/topic/README.md) and [Git Storage](../../../../spec/features/storage/git-storage/README.md), inspected 2026-08-04 | A public topic board is not yet a public-research method |
| E3 | Observed | The current anonymity Feature promises that even admins cannot see an anonymous author while the author can still edit/withdraw and self-voting is prevented; the Feature explicitly leaves the identity mechanism and writing-style re-identification unresolved. | Current [Issue Anonymity](../../../../spec/features/issue/anonymity/README.md), inspected 2026-08-04 | Promise is unimplemented and has no acceptance criteria |
| E4 | Observed | The shared Sneat platform relationship graph is an explicit, typed graph of people, organisations, Spaces and roles. The Contactius Feature prohibits inferred or parallel relationships; current ecosystem sources do not define IssueNumber.one signals in it. | Backstage Contactius relationship-graph Feature and core-modules research pinned in the [brief](../README.md#ecosystem-context-and-ownership-map), June–July 2026 | Establishes existing substrate and ownership, not permission to add issue-derived edges |
| E5 | Observed | Moodometer is framed as an anonymous team pulse and ritual aggregate. Its own Idea identifies a tension between knowing who submitted for an individual gate and showing an anonymous aggregate. | Backstage Moodometer seed and ritual-step Idea pinned in the [brief](../README.md#ecosystem-context-and-ownership-map), July 2026 | Conceptual; demand and combined-decision value unverified |
| E6 | Sourced | GDPR principles require specified purposes, data minimisation, accuracy, storage limitation, lawfulness, fairness, transparency and accountability; rights and safeguards depend on the processing context. | [Regulation (EU) 2016/679, especially Articles 5, 12–22, 25 and 35](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX%3A32016R0679), 2016 | Direct EU law; deployment-specific analysis is outside this report |
| E7 | Sourced | Pseudonymised data that can be attributed with additional information remains personal data. Pseudonymisation reduces risk but does not itself make a record anonymous. | [EDPB pseudonymisation announcement and Guidelines 01/2025](https://www.edpb.europa.eu/news/edpb-adopts-pseudonymisation-guidelines-and-paves-the-way-to-improve-cooperation-with_en), 17 January 2025 | The linked 2025 guidelines were published for consultation; the legal distinction is also stated in GDPR and current EDPB material |
| E8 | Sourced | Effective anonymisation must resist singling out, linkability and inference; simply removing names is insufficient. Combining records can create a mosaic/jigsaw re-identification effect. | [EDPB guidance for SMEs](https://www.edpb.europa.eu/sme/be-compliant/secure-personal-data_en) and [ICO effective-anonymisation guidance](https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/data-sharing/anonymisation/how-do-we-ensure-anonymisation-is-effective/), current pages accessed 2026-08-04 | Risk is contextual; no universal threshold establishes anonymity |
| E9 | Sourced | Employee consent is usually inappropriate for monitoring because employer–worker power imbalance can make it non-free; monitoring must have a defined purpose, necessity, proportionality, transparency and, where high risk, a DPIA. | [ICO data protection and monitoring workers](https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/employment/monitoring-workers/data-protection-and-monitoring-workers/), current page accessed 2026-08-04 | UK guidance under review; still direct regulator evidence of the power-imbalance problem |
| E10 | Sourced | Psychological safety is a shared belief that a team is safe for interpersonal risk-taking; in a 51-team field study it was associated with learning behaviour. | [Edmondson, “Psychological Safety and Learning Behavior in Work Teams”](https://doi.org/10.2307/2666999), 1999 | One multimethod field study; association does not validate anonymous software |
| E11 | Sourced | In a two-phase study of 3,149 employees and 223 managers, managerial openness was more consistently related to improvement-oriented employee voice, mediated by perceived psychological safety. | [Detert and Burris, “Leadership Behavior and Employee Voice: Is the Door Really Open?”](https://doi.org/10.5465/amj.2007.26279183), 2007 | One organisation/industry context; does not establish causation for a product feature |
| E12 | Sourced | Employee survey value depends on follow-up and action planning; a systematic review found the follow-up process often neglected and highlighted purpose, management endorsement, resources and trained change agents. | [Huebner and Zacher, “Following Up on Employee Surveys”](https://pmc.ncbi.nlm.nih.gov/articles/PMC8696015/), 2021 | Systematic review found the evidence base itself limited and heterogeneous |
| E13 | Sourced | AAPOR distinguishes probability and non-probability samples, warns that opt-in/social-network recruitment needs special care, and expects disclosure of population, sample construction, recruitment, mode, wording, weighting and uncertainty. It also warns that demographics plus responses can identify people in public datasets. | [AAPOR Best Practices for Survey Research](https://aapor.org/standards-and-ethics/best-practices/) and [2026 Code](https://aapor.org/standards-and-ethics/), current pages accessed 2026-08-04 | US professional standard, broadly relevant but not a substitute for jurisdiction-specific research governance |
| E14 | Sourced | Quality citizen participation requires clarity, accountability, transparency, inclusiveness, accessibility, privacy and feedback on how input was used or why it was not used. | [OECD Guidelines for Citizen Participation Processes](https://www.oecd.org/en/publications/oecd-guidelines-for-citizen-participation-processes_f765caf6-en.html), 2022 | Applies to participation processes; not every public board represents a government consultation |
| E15 | Sourced | UK consultation principles require a genuine purpose, adequate information, targeted and accessible engagement, scrutiny and a timely published response explaining how input informed policy. | [UK Cabinet Office Consultation Principles](https://www.gov.uk/government/publications/consultation-principles-guidance/consultation-principles-2018), 2018 | Government practice, useful as a public decision-loop benchmark rather than a universal legal rule |
| E16 | Sourced | Official employee-insight products protect results using minimum group sizes, access roles and additional controls. Microsoft documents minimum group size, differential privacy and masked distributions; it also warns that overlapping group comparisons can still identify a person. | [Microsoft Viva Insights technical privacy guide](https://learn.microsoft.com/en-us/viva/insights/advanced/privacy/privacy), updated 30 June 2025 | Vendor practice; configuration and downstream use still matter |
| E17 | Sourced | Employee-survey products distinguish confidential attributed responses from unattributed responses and suppress small groups/comments. Culture Amp recommends reporting-group minimums and indirect protections; Qualtrics defaults engagement and comment thresholds to five. | [Culture Amp participant confidentiality guidance](https://support.cultureamp.com/en/articles/8384294-general-survey-participant-faqs), [Culture Amp pre-launch checklist](https://support.cultureamp.com/en/articles/9267681-engagement-survey-pre-launch-checklist) and [Qualtrics confidentiality overview](https://www.qualtrics.com/support/confidentiality-overview-ex/), current pages accessed 2026-08-04 | Product defaults are not evidence that 3 or 5 is safe for every cohort, free-text set or repeated release |
| E18 | Sourced | Official statistics guidance warns that comparing overlapping tables or time windows can reveal small cells by differencing; multi-dimensional sparse tables are riskier. | [UK Office for National Statistics disclosure-control policy](https://www.ons.gov.uk/methodology/methodologytopicsandstatisticalconcepts/disclosurecontrol/policyonprotectingconfidentialityintablesofbirthanddeathstatistics), current page accessed 2026-08-04 | Domain-specific statistical guidance; the disclosure mechanism is broadly relevant |
| E19 | Sourced | Workplace relationships affect voice: research links employee workflow-network position and leader positions in informal networks with voice behaviour and perceived power/support. | [Venkataramani et al., “Social networks and employee voice”](https://doi.org/10.1016/j.obhdp.2015.12.001), 2016 | Shows relational context matters to voice behaviour; does not show a graph product improves decisions |
| E20 | Sourced | Contextual-integrity research evaluates privacy by whether information flow fits the actors, information, purpose and transmission norms of its social context, not by whether data was technically available. | [Nissenbaum, “Privacy as Contextual Integrity”](https://digitalcommons.law.uw.edu/wlr/vol79/iss1/10/), 2004 | Normative privacy framework, not a statutory test |
| E21 | Sourced | Public participation products demonstrate that issue collection, voting, opinion grouping and accountable follow-through can work without a pre-existing social graph. Polis groups participants by response similarity; Decidim links proposals to official results and progress. | [Computational Democracy Project opinion groups](https://compdemocracy.org/opinion-groups/) and [Decidim accountability](https://docs.decidim.org/en/v0.30/admin/components/accountability), current pages accessed 2026-08-04 | Official product documentation; neither proves population representativeness or outcome superiority |

## Findings

### Observed or sourced

1. **The current product has two scopes, not two validated modes.** Internal
   teams/organisations and public topics exist in the conceptual Feature tree,
   but they share only issue/vote vocabulary. Neither has a specified research
   or intelligence contract (E1–E3).
2. **The current scarce-priority mechanism is already a non-graph baseline.**
   A participant raises a limited number of issues and spends a limited vote
   budget. Any relationship-aware version can therefore be compared with a
   concrete alternative: the same issue/vote set without relationship joins,
   cohorts or inferred paths (E1).
3. **Public visibility is not public-research validity.** A self-selected public
   topic can show what its participants submitted and supported. Without a
   target population, recruitment account and bias analysis, it cannot claim
   that the result estimates what “Dublin,” “AI agents” or another population
   thinks (E2, E13).
4. **Internal identity is not made safe by hiding a name.** Author control,
   duplicate prevention, longitudinal participation, team membership,
   manager relationships and free text can all make a person linkable. Current
   absolute anonymity wording is therefore an unvalidated product promise,
   while pseudonymity and confidentiality are more accurate descriptions of
   many possible workflows (E3, E7–E8, E16–E18).
5. **Employment changes the consent and harm model.** A worker's refusal or
   critical response may have perceived career consequences. This makes
   managerial openness, psychological safety and an independently credible
   purpose part of the product value—not external compliance decoration
   (E9–E12).
6. **Collection without follow-through can destroy the next round's signal.**
   Public-participation guidance and employee-survey research both treat
   feedback/action after collection as part of process quality. A board whose
   terminal state is only “ranked” is incomplete for either mode (E12, E14–E15,
   E21).
7. **Relationship context can explain and route; it can also expose and bias.**
   Organisational network research supports the relevance of workflow and
   leader relationships to voice. The same context can reveal who likely spoke,
   amplify formal hierarchy, turn “insight” into worker evaluation or enable
   retaliation (E8–E11, E19–E20).
8. **Aggregation thresholds are necessary examples, not a guarantee.** Major
   employee-insight products use thresholds, masked distributions, differential
   privacy and indirect-reporting protections. Overlapping filters, rare
   attributes, comments and longitudinal releases still permit inference or
   differencing (E16–E18).
9. **A graph is not required for clustering opinions or closing action loops.**
   Polis derives opinion groups from response similarity, while Decidim links
   proposals to official results and progress. These are credible non-graph
   comparisons for H1/H2 rather than proof that either product fits the target
   job (E21).
10. **The existing Sneat platform graph is an explicit relationship record, not an
    inference sink.** Current ecosystem canon prohibits Contactius from
    inventing derived relationships. No current source authorises issue,
    sentiment, inferred health or benchmark labels as graph edges (E4–E5).

### Mode comparison

| Dimension | Public research / public participation | Internal organisation intelligence |
|---|---|---|
| Primary actors | Commissioner/convenor; researcher or moderator; contributor/respondent; affected but non-participating population; public/decision-maker; data steward | Worker/contributor; peer; manager/issue owner; leadership/buyer; people/HR or analyst; admin/data steward; union/works council or equivalent representative where relevant |
| Core contributor job | Express a view, propose an issue, allocate attention and understand whether the input influenced a public/community decision | Surface a concern or improvement safely, help prioritise it and see that an accountable actor responded without creating career harm |
| Core decision-user job | Understand a defined participant body or population well enough to make and justify a public/community decision | Detect actionable organisational problems across units, assign ownership and improve work without converting feedback into surveillance or individual evaluation |
| Legitimacy source | Declared population and purpose; inclusive recruitment; method transparency; neutral questions; response and limitation reporting | Specific employment purpose; necessity/proportionality; psychological safety; credible confidentiality; access separation; worker information/consultation; action and anti-retaliation practice |
| Contributor incentive | Influence a decision, advance an issue, civic/community efficacy, recognition or group advocacy | Improve work, protect colleagues, remove friction, escalate risk, be heard without retaliation |
| Distorting incentives | Mobilisation/brigading, multiple accounts, campaign capture, performative positions, strategic non-response, loud minorities | Fear, social desirability, loyalty pressure, manager gaming, self-censorship, score optimisation, using feedback against a worker |
| Sampling meaning | Must distinguish an open participation count from a sample-based estimate; relationship/demographic cohorts can diagnose reach but not repair unknown selection bias automatically | A census invitation can still have non-response bias; manager/team cohorts may support action but create small-cell and retaliation risk |
| Anonymity | Public authored speech can be deliberately attributable; confidential/unattributed research needs a separate promise for identity, metadata and publication | “Anonymous to manager” may still be confidential or pseudonymous to the service; identity required for author control or duplicate prevention remains personal data and needs a precisely stated boundary |
| Consent | Participation may be voluntary, but publication, linkage, future research and cross-product reuse are separate purposes; refusal must not be disguised as “no opinion” | Consent is usually a weak lawful/trust basis because of power imbalance; a defined alternative basis, genuine worker information and rights remain necessary (E9) |
| Data rights and expectations | Notice of purpose/audience; ability to skip; correction/withdrawal rules; publication and retention boundaries; access to methodology; portability/export expectations where applicable | Notice of purposes and recipients; access/correction/objection/erasure or restriction as applicable; contesting derived labels; retention; separation from performance use; representation/consultation rights where applicable |
| Actionability | Decision maker must publish acknowledgement, disposition or reason; research output can be informative without promising every proposal will be implemented | Accountable owner needs a safe level of detail, resources and authority to act; contributor needs status without being exposed; “insufficient data” must not mean “no concern” |
| Principal harms | False representativeness, manipulation, harassment, doxxing, group labelling, voter suppression, permanent publication, misleading benchmarks | Retaliation, discrimination, chilling voice, worker monitoring, manager misuse, re-identification, false health/performance inference, collective mistrust |
| Relationship context with plausible value | Geography/constituency/affectedness for recruitment and subgroup coverage; response-similarity groups for plural views; institution/decision ownership for routing | Team/unit/manager/workflow dependency for routing and cross-team patterns; role/tenure/location cohorts for carefully protected aggregate analysis; accountable ownership for follow-up |
| Relationship context with weak current justification | Personal social graph, friend/follow edges or cross-product activity used to weight opinions; inferred affiliation used without disclosure | Informal friendship/avoidance graph, household or personal Sneat platform graph, individual productivity/collaboration graph, manager-path labels attached to anonymous comments |
| Terminal good result | Participants and public can see method, limitations, result, decision response and reuse conditions | Workers and decision users can see an issue's safe aggregate/action status, what changed or why not, retention outcome and next review without exposing an individual |

## End-to-End User Journeys

These journeys describe observable product-research outcomes, including
null-action and divergent epilogues. They are not interface or architecture
designs.

### Journey A — public research or participation

1. **Start: a convenor frames a real open decision.** The convenor publishes
   the question, decision owner, target population or invited participant body,
   what input can change, field period, moderation, data use and publication
   conditions. **Observable good result:** a contributor and an outside reader
   can distinguish open participation from representative research before any
   response is collected (E13–E15).
2. **Contribute: a person raises an issue or allocates scarce support.** The
   person sees whether their response is public, confidential or unattributed;
   can skip sensitive questions; and is not silently enrolled into graph reuse.
   **Observable good result:** the contribution is accepted once, its audience
   and withdrawal rule are visible, and missing response remains different from
   disagreement (E13).
3. **Other actor sees and uses the signal.** A moderator/researcher sees
   participation quality, coverage gaps, issue/support patterns and limitations;
   the decision maker sees results with methodology rather than a bare
   leaderboard. **Observable good result:** relationship context is used only
   for the declared comparison or routing purpose, and the non-graph total
   remains available as a reference.
4. **Decision/action.** The decision owner accepts, rejects, combines, defers or
   records an issue as informative only. **Observable good result:** the public
   result links the input to a response and explains the method and important
   exclusions; it does not imply every participant represents a population or
   every popular item must be implemented (E14–E15, E21).
5. **Terminal result.** The field period closes with result, limitations,
   disposition, retention/publication state and next review. **Observable good
   result:** a participant can tell that their participation ended and what it
   achieved; an analyst can reproduce the interpretation without respondent
   identities.
6. **Null-action path.** If no contributor acts, the result is “no responses” or
   “insufficient participation,” never “no issue.” If the decision owner takes
   no action, the board keeps a dated awaiting-response/deferred state and does
   not fabricate adoption. **Observable good result:** silence is not converted
   into population opinion, and non-response by the authority is visible.
7. **Close epilogue.** The convenor closes the process. **Observable good
   result:** invitations and collection stop; retention/publication rules take
   effect; no provisional issue remains presented as an active commitment.
8. **Share/replay epilogue.** An authorised reader shares or reruns the work.
   **Observable good result:** the package includes prompt/question wording,
   population, recruitment, dates, mode, sample/participation counts, weighting
   or clustering, moderation, uncertainty, response and reuse conditions—but
   not respondent contact records or unnecessary microdata (E13–E15).

### Journey B — internal organisation intelligence

1. **Start: the organisation defines a bounded improvement purpose.** Workers
   are told the question, recipient, lawful/trust basis, whether responses are
   authored/confidential/unattributed, group thresholds, prohibited uses,
   retention, action owner and challenge route. Relevant worker representatives
   are involved where required. **Observable good result:** a worker can explain
   who can see what and can decline an optional contribution without believing
   refusal affects their employment (E9–E12).
2. **Contribute: a worker raises an issue or allocates scarce support.** The
   worker chooses the available attribution level and sees how free text,
   membership and repeated responses affect confidentiality. **Observable good
   result:** the system does not claim true anonymity merely because a manager
   cannot see a name; the contribution remains within the stated employment
   purpose (E7–E9, E16–E18).
3. **Other actor sees and uses the signal.** An authorised manager, issue owner
   or analyst receives only the level of issue, aggregate or cohort detail that
   their action requires. **Observable good result:** suppressed/insufficient
   groups remain suppressed across filters; the user cannot pivot through
   relationships to isolate a person; personal feedback is not turned into a
   performance score (E16–E18).
4. **Action and worker-visible follow-up.** An accountable actor acknowledges,
   investigates, assigns, resolves, defers or explains a constraint. The process
   separates genuine safeguarding/escalation from a popularity contest.
   **Observable good result:** the worker sees a safe status and the actor has
   resources and authority for follow-through; no need exists to reveal the
   author to show that action occurred (E12).
5. **Terminal result.** The issue is resolved, closed with reason, transferred to
   a named safe process, or marked not acted on with a review trigger.
   **Observable good result:** contributors and decision users get terminal
   state; individual-level evidence follows its retention/access rule; aggregate
   learning is not silently repurposed.
6. **Null-action path.** If nobody submits, the output is “no/insufficient
   response,” not organisational health. If a manager takes no action, the lack
   of response is visible to the authorised cohort and escalates only through a
   declared process; the product does not expose non-submitters or pressure
   them. **Observable good result:** neither silence nor manager inactivity is
   misrepresented as wellbeing or resolution.
7. **Close epilogue.** The organisation closes the cycle. **Observable good
   result:** reminders stop; raw and identity-linking material follows retention;
   unresolved safeguarding items are transferred explicitly rather than
   disappearing with a survey close.
8. **Share/replay epilogue.** An authorised successor reviews or reruns the
   cycle. **Observable good result:** they receive purpose, instrument, cohort
   definitions, thresholds, field dates, action history and limitations, not a
   portable employee-identity or free-text corpus. A changed purpose starts a
   new permission and rights assessment.

## Relationship Context Compared with a Non-Graph Baseline

The table separates plausible decision value from graph rhetoric. “Graph” here
means explicit relationship or cohort context beyond the issue/vote record; it
does not imply where data is stored.

| Decision job | Non-graph baseline | Relationship context that might add value | Observable beneficiary/value | Evidence state | Principal counter-risk |
|---|---|---|---|---|---|
| Find the public's top issue | Rank all self-selected responses | Declared geography/constituency/affected group used to show participation coverage and subgroup differences | Researcher and public avoid mistaking one mobilised group for everyone | Methodologically plausible from AAPOR/OECD (E13–E14); no IssueNumber.one outcome evidence | Demographics do not repair unknown self-selection; small cohorts expose people |
| Understand plural public positions | Overall support score | Response-similarity groups derived from issue votes, without a personal social graph | Decision user sees common ground and minority disagreement hidden by one average | Demonstrated product practice in Polis (E21), not validated for scarce-vote mechanics | Algorithm choices shape groups; participants seeing results can influence later responses |
| Route a public issue | Manual triage by topic | Explicit institution, jurisdiction or decision-owner relationship | Contributor reaches the actor capable of responding; convenor reduces manual routing | Plausible; public consultation guidance requires clear decision ownership (E14–E15) | Incorrect/inferred ownership misroutes or implies endorsement |
| Find an internal top issue | Rank all team issues equally | Explicit team/unit/workflow dependency used to show where the issue occurs or crosses boundaries | Issue owner can act at the right scope; worker avoids repeating the issue in every team | Plausible from organisational context and Sneat Work graph; no comparative result | Formal hierarchy can suppress cross-cutting or minority concerns |
| Interpret internal silence/voice | Count submissions and votes | Protected cohort response rates and organisational unit context | Analyst sees missing/underrepresented groups rather than treating silence as health | Supported as survey-quality need (E12–E13) | Revealing which small group did not respond pressures individuals |
| Identify a safe action owner | Free-text mention/manual assignment | Explicit responsibility or workflow relationship, separately verified | Contributor gets follow-through; manager sees only what is needed to act | Plausible; action loop supported (E12) | Automatically routing to the subject of a complaint creates retaliation risk |
| Detect cross-team friction | Separate team leaderboards | Aggregate issue co-occurrence across declared workflow dependencies | Leadership can see systemic rather than local causes | Hypothesis only; relationship relevance supported generally (E19) | Joins can expose rare teams/roles and turn narrative concerns into worker scores |
| Compare organisations | Overall customer averages | Methodologically comparable peer strata such as size, sector, region and field period | Organisation interprets its aggregate against relevant peers | Plausible only with comparable instrument/recruitment/definitions (E13) | Peer filtering and repeated releases can reveal organisations; opt-in customers bias benchmarks |
| Combine mood and issue patterns | View aggregate mood and issue ranks separately | Compare protected cohort/time aggregates under one declared decision question | Action owner may distinguish acute sentiment from persistent issue structure | Plausible complementarity only; no evidence of incremental decision quality (E5, E12) | Sensitive inference, false causality and re-identification by small/time-linked cohorts |
| Improve public distribution | Share a topic link | Existing organisation/community affiliation used only for disclosed recruitment coverage | Convenor reaches otherwise missing affected groups | OECD/AAPOR support inclusive targeted recruitment (E13–E14) | Social-graph recruitment over-represents connected users and enables pressure/brigading |

## Candidate Signal Assessment

“Plausible consumer” identifies where a future user benefit might exist. It is
not an integration recommendation, permission decision or assertion that the
consumer currently supports the signal.

| Candidate signal | Source relationship | Intended beneficiary | Permission / trust domain | Plausible consuming Sneat product | Non-graph baseline | Possible compounding mechanism | Abuse / re-identification risk and current boundary |
|---|---|---|---|---|---|---|---|
| Issue statement and scope | Contributor→declared public topic or work Space | Contributor and accountable issue owner | Public-process domain or one organisation Space; purpose-specific | IssueNumber.one; possibly Sneat Work Flow as a deliberately transferred action item | Standalone issue record with manually chosen scope | Repeated, consistently classified issues improve retrieval and action history | Free text can name or identify people and reveal special-category beliefs. Keep raw internal/public corpora separated; no automatic personal-graph write |
| Scarce support/vote and rank | Participant→issue within one process/team | Cohort and decision owner | The originating ballot/priority cycle | IssueNumber.one | Flat up/down/support total | Longitudinal ranks show changing priority under one stable instrument | Attribution, minority positions and strategic voting can harm; do not export individual votes or turn them into relationship edges |
| Issue lifecycle/action outcome | Accountable actor→issue | Contributor, affected cohort and future decision owner | Originating public process or organisation; disposition may be public at aggregate/item level | IssueNumber.one; Sneat Work Flow when action is explicitly handed off | Manual status note | Completed loops improve institutional memory and permit outcome-based learning | Owner assignment can expose complainant or accused; cross-product transfer needs explicit purpose and audience |
| Public recruitment/source channel | Convenor→invited population/channel | Researcher and public interpreting coverage | Research operations; separate from published responses | IssueNumber.one research reporting; Invitus only as a disclosed delivery channel | Campaign/source counts in research log | Better channel evidence can reduce repeated coverage gaps | Contact records and source affiliation can reveal politics/community ties. Keep recruitment identities isolated from response data and public graph |
| Public geography/constituency/affectedness | Respondent→declared place or affected group | Researcher, public and decision owner | Purpose-specific research cohort; public only as safe aggregate | IssueNumber.one; perhaps Localius for a separately consented public place reference | Self-declared category stored only with survey | Repeated comparable cohorts may improve reach diagnosis and local routing | Rare locations/roles plus text can identify people. Do not infer from personal address graph; suppress unsafe cells and overlapping slices |
| Demographic cohort attributes | Respondent→self-declared research categories | Researcher diagnosing inclusion/bias | Research-only; minimum necessary attributes | IssueNumber.one research analysis | No demographic segmentation | Consistent definitions can support comparability over time | High singling-out/inference risk, especially intersections. Keep microdata out of shared person graph; public release only with disclosure control |
| Opinion-similarity cluster | Response pattern→derived group | Participants and decision owner | One conversation/field period; explicitly labelled inference | IssueNumber.one | Overall mean/rank or manual thematic coding | More responses can stabilise within-process clusters and reveal cross-group common ground | Cluster can become a political/person label. Do not write to person/relationship graph or reuse outside stated process; expose method/uncertainty |
| Organisation/team membership | Explicit person→Space membership | Worker, analyst and action owner | Organisation's work Space; role/access governed independently | IssueNumber.one inside Sneat Work Team | Self-selected team field | Accurate group membership reduces manual routing and supports protected cohort analysis | Membership plus comment/time can identify author. Read as context only when necessary; do not attach response-derived labels back to member |
| Manager/direct-report relationship | Explicit work relationship | Action owner and governed analyst | Employment trust domain with especially restrictive purpose/access | Sneat Work Team; IssueNumber.one internal mode | Contributor selects an action owner manually | Current reporting lines can route responsibility and define protected groups | Manager is also a potential retaliator and evaluator. Never expose anonymous contribution or inferred score through graph traversal; stale lines misroute |
| Workflow/dependency relationship | Explicit team/service/process dependency | Cross-team issue owners | Organisation operations domain | Sneat Work Flow and IssueNumber.one | Free-text “blocked by” / manual tag | Repeated cross-boundary issues may reveal systemic bottlenecks | Can become productivity/performance surveillance. Use declared operational links, not inferred personal collaboration intensity |
| Informal friendship, avoidance or communication-centrality relation | Observed/inferred person↔person tie | Researcher might seek explanatory context | Highly sensitive employment/person domain | No plausible consumer established | No informal-network join | Research says such positions relate to voice (E19), but reuse could increase explanatory power | High chilling, inference and retaliation risk; evidence supports relevance, not product collection. Keep outside IssueNumber.one/shared graph absent a distinct authorised study |
| Anonymous-author control mapping | Hidden contributor identity→issue credential/control | Contributor and integrity function only | Narrow identity/control domain | IssueNumber.one only | One-time anonymous token with no graph join | No legitimate compounding value beyond integrity over that issue | Direct re-identification key. Must remain isolated from managers, analysts, AI and shared graph; pseudonymity is not anonymity (E7–E8) |
| Attributed/confidential/unattributed mode marker | Contribution→disclosure promise | Contributor and authorised viewer | Originating cycle; promise locked with result | IssueNumber.one | One global anonymity toggle | Accurate trust metadata can make repeated participation interpretable | Changing promise later destroys trust; marker does not itself protect data. Keep access/audience explicit and non-retroactive |
| Internal free-text comment | Worker→issue/feedback cycle | Worker and safe action owner | Restricted organisation process; comments may need stricter thresholds than scores | IssueNumber.one internal mode | Moderated text inbox | Thematic history might improve recurring-issue detection | Writing style, incidents and rare facts re-identify. Do not place raw comments in relationship graph, benchmarks or broad AI corpus |
| Individual mood pulse | Worker→ritual/session/time | Individual (if personal reflection) or protected team aggregate | Moodometer/ritual trust domain | Moodometer/Trackus; aggregate may inform authorised ritual owner | Separate anonymous pulse chart | Stable instrument can reveal within-team trend | Health/wellbeing inference, timestamp and required-submission mapping expose individuals. Keep individual pulse and identity gate isolated from IssueNumber.one and shared graph |
| Protected mood aggregate | Cohort/time window→aggregate pulse | Team and authorised action owner | Organisation cohort with threshold/differencing controls | Moodometer; possible read-only comparison in internal research | Separate chart viewed beside issues | Repeated safe aggregates may show trend alongside issue lifecycle | Overlapping teams/windows enable differencing; correlation is not causation. No cross-org or personal graph reuse by default |
| Cross-organisation benchmark contribution | Organisation→instrument/version/peer stratum | Contributing organisation and benchmark users | Benchmark consortium/research domain with purpose, method and contribution rules | IssueNumber.one benchmark report | Each organisation compares only to its own history or public external datasets | More comparable organisations may improve peer coverage and precision | Opt-in customer sample bias; organisation re-identification; commercial leakage. Raw organisation contributions remain isolated; publish protected aggregates/method only |
| AI-generated theme, risk or organisational-health inference | Issues/aggregates→derived label | Authorised decision user | Same or narrower domain than every input; explicitly derived and contestable | IssueNumber.one AI summary | Human coding/summary | Reviewed labels can make repeated issue sets easier to scan | False sensitive/performance labels can spread and look factual. Do not write inferred person/team traits to shared graph; retain provenance, expiry and human challenge |
| Moderation, abuse, safeguarding or retaliation report | Reporter→incident/subject/moderator | Reporter and authorised safety actor | Separate high-restriction case domain | IssueNumber.one moderation/safety process | Dedicated confidential reporting channel | Pattern detection may reveal systemic harm under independent governance | Highest retaliation/defamation/legal risk. Keep outside ordinary issue ranking, benchmarks, AI summaries and relationship graph; transfer only through declared safety process |

### Signal classes with strong present evidence for isolation

The evidence supports isolation as a research constraint, not a storage design:

1. **Identity-control secrets:** anonymous-author mappings, duplicate-prevention
   keys and respondent contact/recruitment records. Their beneficiary is
   integrity or follow-up, not analysis.
2. **Individual sensitive experience:** mood/wellbeing pulses, harassment,
   safeguarding, retaliation and health-adjacent inferences. Aggregated use does
   not authorise personal reuse.
3. **Free text and moderation evidence:** unique incidents and writing style can
   identify a person even when names are removed.
4. **Public research microdata:** demographic, location, source-channel and
   response joins can identify respondents or political/community affiliation.
5. **Organisation benchmark microdata:** raw organisation contributions and
   peer filters can reveal commercially or reputationally sensitive results.
6. **Derived labels:** opinion clusters, risk, performance, sentiment and
   organisational-health interpretations are contextual and contestable, not
   durable facts about a person, team or relationship.
7. **Cross-context personal graph:** household, family, sports, customer or
   unrelated Sneat platform product relationships have no evidenced purpose in either
   mode. A shared account is not consent to cross-context inference (E20).

Potentially shareable items are narrower: a deliberately public issue and its
public disposition; a separately authorised work action transferred to its
accountable workflow; or a methodologically safe aggregate accompanied by its
population, instrument, period, threshold and limitations. Even these remain
purpose- and audience-bound rather than “graph signals” by default.

### Inferred interpretations

- **Trust is part of measurement validity.** In internal mode, fear changes who
  speaks and what they say; in public mode, recruitment and perceived influence
  change who participates. Privacy and actionability controls therefore change
  the signal itself, not only compliance risk (E9–E15).
- **Relationship context has asymmetric value.** Explicit work ownership or
  affected-population context can improve routing and interpretation. Personal
  social ties and inferred attributes carry much greater exposure risk and have
  no demonstrated incremental value for this product.
- **Cross-mode reuse has a purpose problem before it has a technical problem.**
  The same issue text or cohort label changes meaning when the sender, recipient,
  purpose and transmission rule change. Technical commonality does not establish
  contextual integrity (E6, E20).
- **Benchmark value is conditional, not automatically cumulative.** More
  contributions improve a benchmark only when instrument, population,
  recruitment, field period, definitions and peer strata are sufficiently
  comparable. More opt-in customers can also compound selection bias and
  disclosure surface (E13, E16–E18).
- **Action history may be a more defensible learning signal than personal
  relationships.** An issue's acknowledged/deferred/resolved path can improve
  institutional learning without claiming a sensitive trait about its author.
  This remains a hypothesis; no retention or product outcome evidence exists.

## Hypothesis Assessment

| Hypothesis | Current assessment | Evidence for | Evidence against / disconfirming | Decision-changing next evidence |
|---|---|---|---|---|
| H1 — relationship context materially improves decision-useful insight or adoption over a non-graph baseline | **Unresolved; low confidence.** Plausible for bounded routing, coverage and protected cohort questions, not established as a general advantage | Organisational relationships affect voice (E19); public-research standards require population/recruitment context (E13–E15); explicit Sneat Work graph exists (E4) | Non-graph systems already cluster opinions and close action loops (E21); context increases exposure/bias; no IssueNumber.one comparison or user outcome | Blind comparison of the same synthetic/public issue set with and without one declared relationship context, scored by target decision users for changed decision, calibration and harm |
| H2 — the durable asset is a permissioned issue-relationship graph | **Unresolved, with stronger evidence against the asset claim than for it.** | Explicit issue/action histories and verified organisational relationships could reduce repeated classification/routing | GDPR purpose/data-minimisation constraints (E6), contextual-integrity limits (E20), portability of explicit org data, non-graph alternatives (E21), and high isolation burden weaken “graph as asset”; no switching-cost evidence | Identify one repeated user decision that cannot be performed comparably with explicit fields/manual routing, then measure outcome and switching cost without cross-domain reuse |
| H3 — public and internal modes can share stable primitives while maintaining separate trust/data-use domains | **Unresolved.** Surface vocabulary may be shareable; a shared trust proposition is not supported | Both involve a question, issue/contribution, response/support, disposition and method metadata | Actors, power, sampling, anonymity, lawful basis, publication, harms and terminal obligations diverge materially; shared identity/graph reuse risks context collapse (E6–E20) | Journey comprehension and misuse review of separate products vs visibly separated modes using identical synthetic tasks; require participants to accurately predict recipients and reuse |
| H4 — cross-organisation benchmarks become more valuable and defensible as participation grows | **Unresolved; network-effect language is premature.** | More comparable observations can improve peer coverage/precision; protected aggregation is established practice (E16–E18) | Participation volume cannot fix opt-in bias, incompatible instruments or peer definitions; more intersections/releases increase disclosure; no buyer/value evidence | Benchmark simulation with stable instrument and deliberately imbalanced organisations, measuring when added participants reduce uncertainty versus increase bias/disclosure |
| H6 — each capability can create a consented signal useful to another Sneat product | **Currently contradicted as a universal claim; low-medium confidence.** | Public dispositions and explicitly transferred work actions have plausible consumers | Anonymous mappings, mood pulses, free text, moderation reports, research microdata and derived labels have no evidenced safe cross-product beneficiary; employee consent is not a general cure (E7–E9) | Per-capability necessity test: named user decision, beneficiary, non-graph alternative, lawful/trust basis and comprehension. Any capability with no passing use remains isolated rather than being forced into the hypothesis |
| H8 — permitted combination of issue intelligence and Moodometer improves organisational-health decisions enough to justify trust cost | **Unresolved; very low confidence.** | Issue/action and aggregate mood offer different descriptive views; survey follow-up evidence supports connecting measurement to action (E12) | No causal model, decision-user evidence or comparative outcome; individual/temporal joins create health inference, small-cell and differencing risk (E5, E8, E16–E18) | Synthetic separate-versus-combined aggregate decision task with no personal data; test whether combination changes a correct action and whether users over-infer causes or individuals |
| H10 — public and internal modes belong under one product/brand | **Unresolved; low confidence, with evidence against assuming unity.** | Current canon already uses shared issue/vote language and both scopes; a common brand might reduce discovery cost | Buyers, contributors, legitimacy, consent, harm, distribution and terminal journeys diverge; shared brand can imply unsafe reuse or equal methodology (mode comparison above) | Message/comprehension test of one brand with separated modes, two products and single-mode/status-quo concepts; measure trust prediction and task fit, not preference alone |

## Contradictions and Disconfirming Evidence

- The current Anonymity Feature promises identity hidden even from admins, but
  requires author control and self-vote prevention while leaving the mechanism
  unresolved. External guidance says such linkable records are normally
  pseudonymous personal data, not proven anonymous (E3, E7–E8).
- The current Permissions Feature derives access from membership and explicitly
  has no base role system, while internal intelligence commonly needs distinct
  contributor, manager, analyst, administrator and safeguarding recipients.
  This is a product-trust gap, not an architecture conclusion.
- Current public topics prohibit archiving except by the author, while public
  research/participation practice requires clear moderation, closure,
  disposition and retention. An immutable public topic is not automatically a
  credible or safe research record (E2, E14–E15).
- The founder graph-growth notes treat graph enrichment as a flywheel, but also
  label themselves opinionated. Data minimisation, contextual integrity and the
  lack of a beneficiary for several signal classes directly disconfirm the
  universal “every capability strengthens the graph” interpretation (E6,
  E20).
- Vendor thresholds of three or five are inconsistent with treating any one
  threshold as an anonymity guarantee. Microsoft uses thresholding plus masked
  distributions and differential privacy and still documents overlap risk
  (E16–E18).
- Relationship research shows networks affect voice, but that finding cuts both
  ways: relationship context might explain the signal while collecting or
  exposing it changes perceived speaking risk. It cannot be counted only as
  upside (E10–E11, E19).
- Public product practice demonstrates valuable non-graph baselines. Opinion
  groups can be derived from response patterns and proposals can be linked to
  outcomes without a pre-existing person graph (E21).

## Assumptions, Working Hypotheses and Unknowns

### Assumptions

- “IssuerNumber.One” in the mandate refers to the current `IssueNumber.one`
  product. If not, evidence and canonical routing may target the wrong product.
- Public-research mode includes both open participation and sample-based
  research, but the report keeps their claims distinct. If the sponsor intends
  only public leaderboards, parts of the sampling standard remain an alternative
  baseline rather than an intended capability.
- Internal organisation intelligence means improvement-oriented issue/feedback
  use, not individual performance monitoring. If worker evaluation is intended,
  the harm, legal and evidence bar changes materially.
- Sneat platform relationship sources describe explicit current relationships. No
  inferred relationship, sentiment or trait is treated as a shared primitive.
- “Anonymous” is used only where a person cannot reasonably be identified;
  otherwise the report uses confidential, unattributed or pseudonymous.

### Unknowns

| Unknown | Decision affected | Leverage | Safe default for research |
|---|---|---|---|
| Which mode has an urgent repeated job and a trusted decision owner? | Product direction and sequencing | Highest | Treat both as candidates; do not generalise from current scope names |
| What exact decision does relationship context change, and for whom? | H1/H2 and graph value | Highest | Require a non-graph comparator for every claim |
| What do workers believe “anonymous,” “confidential” and “private” mean in this context? | Internal trust and feasible journey | High | Make no absolute anonymity claim |
| Which internal actors are safe recipients when the issue concerns the direct manager or HR? | Actionability and retaliation | High | Do not auto-route sensitive content through the subject relationship |
| Is public mode open participation, public consultation, public-opinion estimation, community prioritisation or several distinct jobs? | Method, claims, buyer and brand | High | Label output by actual method; never call self-selected totals representative |
| Which cohort attributes are actually necessary for action or bias diagnosis? | Data minimisation and re-identification | High | Collect none by default in a hypothetical design; test necessity field by field |
| What minimum group and comment protections survive overlapping filters and repeated releases for target organisations? | Internal confidentiality and benchmark feasibility | High | Treat small/overlapping cohorts as unsafe; vendor defaults are not guarantees |
| Can a cross-organisation benchmark define stable instrument, population, peers and field periods? | H4/H7 and defensibility | High | No benchmark claim from customer aggregates |
| Does combined issue+mood information change a correct decision or only confidence/storytelling? | H8 | High | Keep views separate; do not infer health or causation |
| What data-subject and organisation-contributor exit/withdrawal expectations apply to published or benchmarked results? | Rights, retention and replay | High | Avoid raw/public microdata and irreversible cross-product reuse |
| Can one brand communicate the two trust models without people predicting cross-use? | H3/H10 | Medium-high | Do not infer unity from shared issue vocabulary |

## Confidence

**Overall confidence is medium-low.** Confidence is medium on the mode
distinction and isolation risks because multiple independent primary sources,
regulators, standards, research and current product practices converge.
Confidence is low on product value, demand, willingness to pay, brand structure,
relationship lift and benchmark network effects because no IssueNumber.one user
or outcome evidence exists.

The most decision-sensitive unknown is whether a defined decision user makes a
better, safer or faster decision from relationship context than from the same
issues, support and action history without it. A credible comparative result
could raise confidence in H1 and narrow H2/H3/H10. Evidence that users
misunderstand recipients or reuse, or that context does not change decisions,
would lower confidence and strengthen the status-quo/rejection alternatives.

## Risks and Dependencies

- **Promise risk:** using “anonymous” when the service retains a linkable author
  record would create a trust failure before any breach occurs.
- **Retaliation and chilling risk:** relationships that improve routing can also
  identify likely authors or route a concern to its subject.
- **False-representativeness risk:** public rankings can be read as population
  opinion even when they only describe self-selected participants.
- **Function-creep risk:** a useful internal improvement signal can become a
  performance, productivity or health label when shared across products.
- **Benchmark harm:** weak peer definitions can produce confident but misleading
  comparisons; small or repeated cuts can expose organisations or individuals.
- **AI authority risk:** summaries can erase minority evidence, present an
  inference as fact or reproduce identifying detail from free text.
- **Action debt:** asking for difficult feedback without a resourced owner and
  visible response can reduce trust and future participation.
- **Graph accuracy risk:** stale manager/team relationships can route issues
  incorrectly or distort cohorts; current Sneat platform graph existence does not prove
  completeness or consent.
- **Jurisdiction dependency:** lawful basis, worker consultation, special-
  category handling, research safeguards and employment rights vary. Product
  evidence cannot replace deployment-specific review.
- **Evidence dependency:** target workers, public participants, buyers, data
  stewards and decision owners have not been interviewed; no current demand or
  usage evidence was supplied.

## Alternatives Kept Open

These are research alternatives, not recommendations.

| Alternative | What it preserves | What it risks/forgoes | Evidence needed to distinguish it |
|---|---|---|---|
| Separate public and internal products | Clearer buyer, promise, method, brand and trust boundary | Duplicate discovery/workflows; weak ecosystem coherence; two cold starts | Distinct buyer/job evidence plus proof that a shared brand causes material trust or acquisition harm |
| One product with visibly separate modes and non-transferable trust domains | Shared discovery and carefully limited common journey vocabulary | Users may still infer common identity/data reuse; operational misuse surface | Comprehension test showing participants correctly predict recipients, rights and reuse; adversarial misuse review |
| One initial mode, preserve the other only as an option | Lower evidence and trust burden; clearer first job | May delay a valuable adjacent mode or bias later commonality decisions | Comparative urgency, buyer and willingness-to-act evidence by mode |
| Current scarce-priority product without intelligence/graph claims | Uses the existing differentiated concept and cleanest non-graph baseline | May remain a narrow prioritisation tool with weak action/market evidence | Usage and decision-quality evidence for scarce issues/votes alone |
| IssueNumber.one as a Sneat Work Team capability rather than a standalone product | Clear internal context and existing work relationships | Public mode/independent buyers may not fit; brand and trust become ecosystem-dependent | Evidence that initial internal users already enter through Sneat Work and value module-level integration |
| Public participation capability only; reject population-research/benchmark claims | Avoids unsupported representativeness and cross-org comparison | Smaller research/analytics proposition | Buyer/user job that values open prioritisation and accountable response without population estimates |
| Reject relationship enrichment and keep explicit manual context | Minimises privacy, accuracy and cold-start cost | Loses automated cohort/routing hypotheses | Comparative test showing relationship context does not change a decision enough to justify harm/cost |
| Reject both modes/status quo | Avoids investing on concept-only evidence | Opportunity cost if a painful job exists | Problem, buyer and action evidence remains absent after bounded discovery |

## Recommended Research Priority

The next research priority is **not to choose a mode or shared product**. It is
to identify one decision in each mode for which relationship context could
plausibly change the outcome, then compare it with the current non-graph
scarce-issue/vote workflow under the correct trust conditions.

Priority order:

1. **Internal trust/action decision:** establish whether workers and safe action
   owners share a comprehensible confidentiality promise and whether the owner
   can act without identifying a contributor. This is highest risk because a
   false promise can harm workers and invalidate the signal.
2. **Public method/claim decision:** determine whether the intended job is open
   participation or population research, because population, recruitment and
   inference claims differ before relationship value is considered.
3. **Relationship lift:** compare one explicit context field at a time against
   the non-graph baseline; measure changed decision, accuracy, time and harm—not
   perceived cleverness.
4. **Cross-mode comprehension:** test whether one brand/two modes causes people
   to predict cross-use of identity, responses or graph data.
5. **Benchmark and issue+mood simulations:** investigate H4/H8 only after the
   instrument, recipient and safe aggregate question are defined.

This priority follows consequence and decision leverage. It does not select a
commercial or product direction.

## Smallest Learning Actions

No external outreach or sensitive-data collection is authorised. The smallest
safe actions are:

1. **Decision-pair specification (repository/public evidence only).** Write one
   synthetic public decision and one synthetic internal decision, each with the
   same issues/votes in a non-graph version and one relationship-context
   variant. Name exactly what a decision user could do differently. If no
   difference can be named, H1 does not yet merit a product test.
2. **Trust-promise adversarial read.** Apply the EDPB singling-out/linkability/
   inference tests and ONS differencing scenarios to the current conceptual
   anonymity, team visibility and repeated-filter claims. Record which promise
   fails without designing its implementation.
3. **Method classification.** Classify three plausible public cases—open topic,
   targeted consultation and population estimate—against AAPOR/OECD disclosure
   requirements. This tests whether “public research” is one mode or several
   incompatible jobs.
4. **Action-loop inventory.** For each mode, identify the actual other actor who
   can acknowledge, act, defer and close. If that role or authority is absent,
   record the case as collection-only rather than intelligence.
5. **Later, sponsor-authorised interviews.** Interview workers/contributors and
   decision users separately; do not put raw notes or identities in this
   repository. Ask them to predict recipients, re-identification routes, rights
   and terminal action from plain-language concepts rather than asking whether
   they “like” a graph.

## Proposed Experiments

No experiment was run, and none is authorised by this public-source iteration.
Three later experiments could change the named hypotheses if the sponsor first
authorises participants, confidentiality, retention and evidence routing:

1. **H1 relationship-lift test:** randomise decision users to the same synthetic
   issue set with (a) no relationship context or (b) one necessary explicit
   context such as accountable unit. Predefine success as a materially better
   correct decision or faster routing without increased identity guesses;
   failure as no change or greater harm; inconclusive if the decision itself is
   ambiguous.
2. **H3/H10 trust-comprehension test:** compare separate-product, separated-mode
   and current scarce-priority concepts. Success is accurate prediction of who
   sees identity/content, what is reused and what terminal action occurs—not
   brand preference. Any concept that creates confident false predictions fails
   its trust proposition.
3. **H4/H8 simulation:** use synthetic organisations and cohorts to test how
   opt-in imbalance, small groups, overlapping filters and added organisations
   affect benchmark accuracy/disclosure; separately test whether aggregate mood
   changes a correct issue action or induces unsupported causal/individual
   inference.

These are smallest decision-changing tests, not product builds.

## Dissent and Conflicts Requiring Managing Partner Resolution

- **Shared-substrate conflict:** Sneat Work sources value a connected,
  permission-aware graph, while privacy and research evidence says context-
  specific information flows must remain purpose-bound. A shared substrate can
  exist without shared derived signals; synthesis must not treat infrastructure
  reuse as evidence for cross-mode data reuse.
- **Anonymity conflict:** Current canon makes an absolute admin-proof anonymity
  promise while leaving identity control unresolved. This report treats the
  promise as unproven and distinguishes anonymity from confidentiality and
  pseudonymity. Canon remains unchanged pending product decision and evidence.
- **Public-mode conflict:** Founder sources frame shareable public leaderboards
  as a viral loop, while AAPOR/OECD evidence distinguishes participation reach
  from representative research and requires an accountable response. A viral
  board may still be useful, but it cannot silently inherit public-opinion
  claims.
- **Graph-value conflict:** Organisational-network evidence makes relationships
  relevant to voice, but the same relationships heighten retaliation and
  re-identification. The evidence does not justify counting relevance as net
  product value.
- **H6 dissent:** The universal claim that every capability should produce a
  reusable Sneat platform graph signal has concrete counterexamples. Future synthesis
  should either narrow the hypothesis to specific action/public signals or
  retain this report's dissent.
- **Threshold dissent:** Official products use minimum groups of three or five,
  but regulator/statistical evidence rejects a universal safe threshold. This
  report does not endorse a number.
- **Brand dissent:** Shared issue vocabulary supplies weak evidence for one
  brand; diverging users, power, method, rights and harms supply evidence
  against assuming it. No current source resolves H10.

## Research Completion Check

- Owned question answered without product-strategy or architecture selection.
- Both journeys cover start, other-actor use, terminal result, null action,
  close and share/replay branches.
- H1, H2, H3, H4, H6, H8 and H10 remain explicit hypotheses with evidence for,
  against and a discriminating next test.
- Each candidate signal states relationship source, beneficiary, trust domain,
  plausible consumer, non-graph baseline, possible compounding and harm.
- Current external claims link directly to regulators, standards, original
  research or official product documentation and state limitations.
- No domain model, API, storage architecture, integration design, canonical edit
  or product recommendation is included.
