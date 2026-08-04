---
title: "{{EXPERIMENT_TITLE}}"
artifact: experiment
status: Draft
engagement: "{{ENGAGEMENT_ID}}"
author: "{{AUTHOR}}"
model: "{{ACTUAL_MODEL_OR_HUMAN}}"
created: "{{YYYY-MM-DD}}"
updated: "{{YYYY-MM-DD}}"
confidence: "{{PRE_EXPERIMENT_CONFIDENCE}}"
dependencies: []
related_specs: []
promotion_target: "{{TARGET_OR_NONE}}"
---
# {{EXPERIMENT_TITLE}}

## Hypothesis

{{FALSIFIABLE_HYPOTHESIS}}

## Assumptions

- {{ASSUMPTION}}

## Smallest Possible Implementation

{{TEST_AND_WHY_SMALLER_WOULD_NOT_DECIDE}}

## Owner, Cost and Duration

- Owner: {{OWNER}}
- Participants: {{PARTICIPANTS}}
- Cost / budget source: {{COST}}
- Duration: {{DURATION}}
- Stop conditions: {{STOP_CONDITIONS}}

## Metrics and Guardrails

- Success: {{METRIC_THRESHOLD}}
- Failure: {{METRIC_THRESHOLD}}
- Inconclusive: {{RANGE_OR_CONDITION}}
- Guardrails / ethics: {{GUARDRAILS}}

## Expected Learning

{{WHAT_CONFIDENCE_CAN_CHANGE}}

## Precommitted Decisions

- On success: {{DECISION}}
- On failure: {{DECISION}}
- If inconclusive: {{DECISION}}

## Results

- Actual execution: {{PENDING}}
- Deviations: {{PENDING}}
- Measurements: {{PENDING}}
- Confidence after: {{PENDING}}
- Decision taken: {{PENDING}}
