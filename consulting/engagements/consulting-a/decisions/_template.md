---
title: "{{DECISION_TITLE}}"
artifact: decision
status: Draft
engagement: "{{ENGAGEMENT_ID}}"
author: "{{AUTHOR}}"
model: "{{ACTUAL_MODEL_OR_HUMAN}}"
created: "{{YYYY-MM-DD}}"
updated: "{{YYYY-MM-DD}}"
confidence: "{{LOW_MEDIUM_HIGH_WITH_REASON}}"
dependencies: []
related_specs: []
promotion_target: "{{TARGET_OR_NONE}}"
---
# {{DECISION_TITLE}}

## Question and Owner

- Question: {{QUESTION}}
- Decision owner: {{OWNER}}
- Deadline: {{DATE_OR_TRIGGER}}

## Context

{{CONTEXT}}

## Options

| Option | Evidence | Trade-offs | Risks | Reversibility |
|---|---|---|---|---|
| Status quo | {{EVIDENCE}} | {{TRADE_OFFS}} | {{RISKS}} | {{REVERSIBILITY}} |
| {{OPTION}} | {{EVIDENCE}} | {{TRADE_OFFS}} | {{RISKS}} | {{REVERSIBILITY}} |

## Criteria and Assumptions

- Decision criteria: {{CRITERIA}}
- Assumptions: {{ASSUMPTIONS}}

## Dissent

- {{POSITION_EVIDENCE_AND_STATUS_OR_NONE}}

## Decision and Rationale

{{CHOSEN_OPTION_AND_WHY}}

## Linked Experiments

- {{PATH_OR_NONE}}

## Review Trigger and Successor

- Review trigger: {{TIME_EVIDENCE_DEPENDENCY_OR_EXPERIMENT}}
- Supersedes: {{DECISION_OR_NONE}}
- Superseded by: —
