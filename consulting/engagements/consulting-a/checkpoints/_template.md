---
title: "{{ENGAGEMENT_TITLE}} Iteration {{ITERATION_ID}} Checkpoint"
artifact: iteration-checkpoint
status: Draft
engagement: "{{ENGAGEMENT_ID}}"
author: "{{AUTHOR}}"
model: "{{ACTUAL_MODEL_OR_HUMAN}}"
created: "{{YYYY-MM-DD}}"
updated: "{{YYYY-MM-DD}}"
confidence: "{{LOW_MEDIUM_HIGH_WITH_REASON}}"
dependencies: []
related_specs: []
promotion_target: "None — navigation state is not promoted"
iteration: "{{ITERATION_ID}}"
predecessor: "{{PREVIOUS_CHECKPOINT_OR_NONE}}"
---
# {{ENGAGEMENT_TITLE}} Iteration {{ITERATION_ID}} Checkpoint

This checkpoint is a compact navigation aid. Its linked evidence and decisions
remain authoritative for their claims; this summary does not replace them.

## What we learned

- {{LEARNING_WITH_EVIDENCE_LINK}}

## What changed

- {{CHANGED_SCOPE_PRIORITY_PLAN_OR_NONE_WITH_SOURCE}}

## Decisions

- {{DECISION_ID_LINK_AND_OUTCOME_OR_NONE}}

## Hypothesis transitions

| Hypothesis | Previous state | Current state | Evidence / next test |
|---|---|---|---|
| {{HYPOTHESIS_ID}} | {{STATE}} | Validated / Rejected / Changed / Unresolved | {{SOURCE_OR_NEXT_TEST}} |

## Assumption changes

- {{ASSUMPTION_INTRODUCED_RETIRED_OR_CHANGED_WITH_REASON}}

## Open questions

- {{QUESTION_DECISION_AFFECTED_AND_BLOCKING_STATE_OR_NONE}}

## Next iteration

- Objective: {{OBJECTIVE}}
- Highest-leverage unknown: {{UNKNOWN}}
- Smallest safe learning action: {{ACTION}}
- Owner and review trigger: {{OWNER_TRIGGER}}
