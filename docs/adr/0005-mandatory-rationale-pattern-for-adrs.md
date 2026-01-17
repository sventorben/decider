---
adr_id: ADR-0005
title: "Mandatory Rationale Pattern for ADR Decisions and Alternatives"
status: adopted
date: 2026-01-17
scope:
  paths:
    - "docs/adr/**"
    - ".claude/**"
    - "internal/**"
    - "cmd/**"
tags:
  - format
  - methodology
  - enforcement
constraints:
  - "All ADRs MUST include 'Adopted because:' and 'Adopted despite:' for the chosen option"
  - "All rejected alternatives MUST include 'Rejected because:' and 'Rejected despite:'"
  - "Pros/cons tables without explicit rationale are prohibited"
  - "The decider CLI MUST validate the presence of rationale patterns"
invariants:
  - "Every ADR decision has explicit trade-off documentation"
  - "Every rejected alternative acknowledges its strengths"
  - "Rationale pattern validation is part of decider check adr"
supersedes: []
superseded_by: []
related_adrs:
  - ADR-0002
---

# ADR-0005: Mandatory Rationale Pattern for ADR Decisions and Alternatives

## Context

ADRs document architectural decisions, but the quality of decision documentation varies significantly. Common problems include:

- **Hidden trade-offs**: Decisions list benefits but omit costs
- **Unfair comparisons**: Rejected alternatives are dismissed without acknowledging their merits
- **Hindsight bias**: Decisions appear obvious in retrospect, hiding the difficulty of the choice
- **Vague language**: "Better fit" or "more suitable" without concrete reasoning

Decision drivers:
- Future engineers need to understand WHY decisions were made, not just WHAT was decided
- AI agents need structured, predictable rationale sections to parse decision context
- Code review of ADRs requires clear criteria for completeness
- Reversing decisions requires understanding original trade-offs

## Decision

All ADRs in this repository MUST use an explicit rationale pattern with mandatory sections for both adopted and rejected options.

### Explicit Rationale Pattern: Adopted

**Adopted because:**
- Forces authors to articulate concrete reasons, not just preferences
- Creates symmetry between adoption and rejection reasoning
- Makes trade-offs explicit and reviewable
- Enables automated validation via `decider check adr`
- Provides AI agents with predictable structure for parsing decisions
- Reduces hindsight bias by requiring "despite" sections upfront

**Adopted despite:**
- Increases ADR writing effort compared to freeform prose
- May feel bureaucratic for simple decisions
- Requires discipline to maintain quality in "despite" sections
- Some trade-offs may be difficult to articulate at decision time

### Pattern Structure

For the **adopted option**:
```markdown
### [Option Name]: Adopted

**Adopted because:**
- Concrete reason tied to decision drivers
- Technical, operational, or strategic justification

**Adopted despite:**
- Known downside we consciously accepted
- Trade-off compared to alternatives
```

For each **rejected option**:
```markdown
### [Option Name]: Rejected

**Rejected because:**
- Concrete reason for rejection
- How it failed to meet decision drivers

**Rejected despite:**
- Legitimate strength of this option
- Benefit that made it attractive
```

### Enforcement

- `decider new` generates ADR templates with the pattern structure
- `decider check adr` validates presence of required headings
- Default mode: warns on missing pattern
- Strict mode (`--strict`): fails with exit code 2 on missing pattern

## Alternatives Considered

### Pros/Cons Tables Only: Rejected

**Rejected because:**
- Tables list attributes without explaining why they matter for THIS decision
- No explicit weighting or prioritization of concerns
- Easy to cherry-pick pros/cons to justify predetermined conclusions
- Does not force acknowledgment of trade-offs in the chosen option

**Rejected despite:**
- Familiar format that many engineers know
- Quick to write and easy to scan
- Works well for simple comparisons
- Can be combined with narrative for hybrid approach

### Freeform Narrative: Rejected

**Rejected because:**
- No predictable structure for AI agents to parse
- Quality varies dramatically between authors
- Easy to omit trade-offs or rejected alternative analysis
- Difficult to validate completeness programmatically

**Rejected despite:**
- Maximum flexibility for nuanced explanations
- Natural writing style for experienced technical writers
- Can capture context that structured formats miss
- No learning curve for new contributors

### Optional Rationale Pattern: Rejected

**Rejected because:**
- Optional patterns become unused patterns over time
- Inconsistent ADR quality across the repository
- Cannot rely on pattern presence for tooling or agents
- Creates two tiers of ADR quality

**Rejected despite:**
- Lower barrier to creating ADRs initially
- Allows gradual adoption and learning
- Respects author judgment on when depth is needed
- Avoids bureaucratic feel for simple decisions

## Consequences

**Positive:**
- Consistent, reviewable decision documentation
- AI agents can reliably parse ADR structure
- Trade-offs are explicit, reducing future confusion
- Automated validation catches incomplete ADRs

**Negative:**
- Higher effort to write ADRs (mitigated by templates and CLI)
- May discourage ADR creation for minor decisions (acceptable trade-off)
- Requires updating existing ADRs for compliance (one-time cost)

## Agent Guidance

When creating or reviewing ADRs:
- Always include both "because" AND "despite" sections
- Never use only pros/cons tables
- If "despite" is difficult to write, consider whether the analysis is complete
- Run `decider check adr --strict` before finalizing
