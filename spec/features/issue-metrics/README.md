---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Issue performance metrics

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue-metrics?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue-metrics?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue-metrics?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/issue-number-one/spec/features/issue-metrics?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

Shows team- and department-level measures of issue closure speed and how often originating issues rise through the priority hierarchy.

## Problem

Current rankings show which issues matter now, but they do not show whether a team or department resolves issues promptly or how often its originating issues rise through the organization. Counts alone can also mislead when scopes create different numbers of issues. Leaders need comparable scope-level measures with the population and calculation made explicit.

## Behavior

### Metrics are available for teams and departments

Authorized viewers can inspect issue-performance metrics for an individual team or department over a stated reporting period.

#### REQ: team-and-department-metric-views

The product MUST provide issue-performance metrics separately for each team and department available to the viewer.

### Closure speed is measured per originating scope

The metrics show how quickly issues associated with a team or department are closed. The exact starting event and summary statistic remain open, so the product must not present an unlabeled number as closure speed.

#### REQ: closure-speed-by-originating-scope

For each team and department, the product MUST show a closure-speed measure for the applicable issue population and MUST display the reporting period and calculation definition used.

### Reach is shown as both a count and a percentage

For each team and department, the metrics show how many originating issues reach the defined top-priority milestone and what percentage of the applicable originating issue population that represents.

#### REQ: top-reach-count-and-percentage

For each team and department, the product MUST show both the absolute number and percentage of applicable originating issues that satisfy the defined “reaches the top” event.

#### REQ: top-reach-denominator-visible

The top-reach metric MUST identify its applicable issue population or denominator so viewers can interpret the percentage.

### Metric context remains visible

#### REQ: metric-context-visible

Every displayed metric MUST identify its organization scope, reporting period, applicable issue population, and calculation definition.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../issue/README.md) | Supplies issue origin, creation, nomination, and closure events |
| [issue/lifecycle](../issue/lifecycle/README.md) | Defines what counts as issue closure |
| [issue/visibility](../issue/visibility/README.md) | Supplies hierarchy promotion events and constrains which underlying issues a viewer may inspect |
| [organization](../organization/README.md) | Supplies team, department, and parent-scope relationships |

## Dependencies

- issue
- organization

## Acceptance Criteria

### Scenario: a viewer inspects closure speed

- **Given** a team has applicable issues in a selected reporting period and some are closed
- **When** an authorized viewer opens that team's issue-performance metrics
- **Then** the view shows the team's closure-speed measure together with its period, issue population, and calculation definition.

### Scenario: a viewer inspects top reach

- **Given** 20 applicable issues originated in a department and 5 satisfied the defined “reaches the top” event
- **When** an authorized viewer opens that department's issue-performance metrics
- **Then** the view shows both 5 issues and 25%, with the 20-issue denominator identified.

## Open Questions

- Does an issue “reach the top” whenever it becomes #1 at its next parent scope, only when it becomes the organization-wide #1, or when it enters a top-N view?
- Does closure speed start at issue creation, personal-#1 nomination, first vote, or first promotion?
- Should closure speed use a median, percentile distribution, average, or another summary?
- How should still-open issues and their current age affect the closure-speed view?
- Which reporting periods and issue cohorts are available?
- Who may view team and department metrics, and what privacy thresholds apply to small groups?
- Is this capability part of the MVP or a later release?

---
*This document follows the https://specscore.md/feature-specification*
