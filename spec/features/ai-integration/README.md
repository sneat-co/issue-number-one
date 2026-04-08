# Feature: AI Integration

**Status:** Conceptual

## Summary

An optional feature a team or org can enable to have an AI model analyze the current set of issues and produce an executive summary — what's trending, where focus is converging, and what risks are emerging. Offered both as a free-prompt mode (bring-your-own-model) and as a hosted SaaS with optimized setup, privacy guarantees, and zero maintenance.

## Problem

Once a team has a healthy flow of issues, scanning them to understand "what's actually going on this week" becomes work. AI is a natural fit for this summary — but it also raises legitimate concerns around privacy (especially for anonymous issues), model choice, and setup burden. IssueNumber.one offers both a DIY and a hosted path so teams can pick the trade-off that matches their context.

## Behavior

### Opt-in per team or org

AI integration is off by default and requires explicit opt-in at the team or org level.

#### REQ: ai-off-by-default

AI integration MUST be off by default. A team or org MUST explicitly enable it before any analysis runs.

### Two modes

| Mode | Description |
|------|-------------|
| `free-prompt` | Team supplies a prompt and model configuration; IssueNumber.one sends issue data to the configured model |
| `saas` | IssueNumber.one-hosted offering with a vetted prompt, managed model access, and privacy controls |

#### REQ: two-modes

AI integration MUST support both a `free-prompt` mode and a hosted `saas` mode.

### Executive summary

The primary output is an executive summary: what the top priorities are, how they've shifted, which issues are gaining or losing support, and where the team's focus is converging.

#### REQ: executive-summary-output

The default analysis output MUST be an executive summary of the current issue set, structured for quick reading by a team lead or manager.

### Anonymity preservation

Analysis MUST respect issue anonymity. The AI MUST NOT be given information that would allow re-identification of an anonymous author, and summaries MUST NOT refer to anonymous authors by any identifier that exposes their identity.

#### REQ: ai-preserves-anonymity

When analyzing issues, the AI integration MUST NOT receive author identity for anonymous issues, and MUST NOT produce output that re-identifies anonymous authors.

### SaaS value proposition

The hosted SaaS mode exists because the DIY path requires teams to make several choices: which model, what prompt, how to handle privacy, where to store prompt history, how to scale. The SaaS mode offers:

- One-click setup with a curated prompt
- No infrastructure or API-key management
- Privacy guarantees including anonymity preservation
- Bundled compliance and billing

#### REQ: saas-value-prop-documented

The SaaS mode MUST document its value proposition versus DIY, covering at minimum: setup cost, maintenance, privacy, and compliance.

## Interaction with Other Features

| Feature | Interaction |
|---------|-------------|
| [issue](../issue/README.md) | Analyses operate on the current issue set |
| [issue/anonymity](../issue/anonymity/README.md) | Anonymity rules constrain the data the AI sees |
| [permissions](../permissions/README.md) | Enabling AI is an org-admin action |

## Dependencies

- issue
- permissions

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- What is the default curated prompt for the SaaS mode?
- How is the SaaS value proposition surfaced and priced versus DIY (free-prompt) mode?
- Are AI outputs stored (as summaries users can revisit) or ephemeral?
- Does the SaaS mode support multiple models, or is the choice hidden entirely?
- How does `free-prompt` mode handle rate limits and cost caps when a team connects a pay-per-token model?
- Acceptance criteria not yet defined for this feature.
