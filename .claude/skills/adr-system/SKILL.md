# ADR System Skill

This skill enables Claude to work with Architecture Decision Records (ADRs) managed by DECIDER.

## Overview

This repository uses DECIDER to manage ADRs. ADRs document significant architectural decisions with structured metadata that enables automated tooling and AI agent consumption.

## Mandatory Rationale Pattern

**ALL ADRs MUST use the explicit rationale pattern. This is enforced by the CLI.**

### For Adopted Decisions

The chosen option MUST include:

```markdown
### [Option Name]: Adopted

**Adopted because:**
- Concrete reason tied to decision drivers
- Technical, operational, or strategic justification

**Adopted despite:**
- Known downside we consciously accepted
- Trade-off compared to alternatives
```

### For Rejected Alternatives

Each rejected option MUST include:

```markdown
### [Option Name]: Rejected

**Rejected because:**
- Concrete reason for rejection
- How it failed to meet decision drivers

**Rejected despite:**
- Legitimate strength of this option
- Benefit that made it attractive
```

### Prohibited Patterns

- **DO NOT** use only pros/cons tables
- **DO NOT** omit "despite" sections
- **DO NOT** use vague language like "better fit"
- **DO NOT** imply rejection without explicit rationale

## ADR Format

ADRs use Markdown with YAML frontmatter:

```markdown
---
adr_id: ADR-NNNN
title: "Decision Title"
status: proposed|adopted|rejected|deprecated|superseded
date: YYYY-MM-DD
scope:
  paths:
    - "path/to/affected/**"
tags:
  - relevant-tag
constraints:
  - "Rule that MUST be followed"
invariants:
  - "Property that must always hold"
supersedes: []
superseded_by: []
related_adrs: []
---

# ADR-NNNN: Decision Title

## Context

Decision drivers:
- Driver 1
- Driver 2

## Decision

### [Chosen Option]: Adopted

**Adopted because:**
- Reason 1
- Reason 2

**Adopted despite:**
- Trade-off 1
- Trade-off 2

## Alternatives Considered

### [Alternative A]: Rejected

**Rejected because:**
- Reason 1

**Rejected despite:**
- Strength 1

## Consequences

**Positive:**
- Benefit 1

**Negative:**
- Cost 1
```

## Key Concepts

### Scope Paths
Glob patterns that indicate which parts of the codebase an ADR applies to. Used by `decider check diff` to find relevant ADRs.

### Constraints
Rules that MUST be followed when working in the affected scope. These are prescriptive.

### Invariants
Properties that must ALWAYS be true. These are conditions to preserve.

### Status Lifecycle
- `proposed` → Under discussion
- `adopted` → Approved and in effect
- `rejected` → Considered but not adopted
- `deprecated` → Was adopted, now discouraged
- `superseded` → Replaced by another ADR

## Commands

```bash
# Initialize ADR structure
decider init

# Create new ADR (includes rationale pattern template)
decider new "Title" --tags foo,bar --paths "src/**"

# List ADRs
decider list
decider list --status adopted --tag security

# Show ADR details
decider show ADR-0001

# Validate ADRs (warns on missing rationale pattern)
decider check adr

# Validate ADRs strictly (fails on missing rationale pattern)
decider check adr --strict

# Check which ADRs apply to changes
decider check diff --base main

# Explain why ADRs apply
decider explain --base main

# Update index
decider index
```

## Working with ADRs

### Before Making Changes

1. Run `decider check diff --base main` to find applicable ADRs
2. Read the constraints and invariants
3. Ensure your changes comply

### Creating a New ADR

1. Run `decider new "Title" --tags ... --paths ...`
2. Fill in the Context section with decision drivers
3. Document the chosen option with:
   - "Adopted because:" (concrete reasons)
   - "Adopted despite:" (trade-offs accepted)
4. Document each rejected alternative with:
   - "Rejected because:" (concrete reasons)
   - "Rejected despite:" (acknowledged strengths)
5. Run `decider check adr --strict` to validate
6. Run `decider index` to update the index

### Reviewing an ADR

Check that:
- [ ] "Adopted because:" section exists with concrete reasons
- [ ] "Adopted despite:" section exists with trade-offs
- [ ] Each alternative has "Rejected because:" and "Rejected despite:"
- [ ] No prose-only sections without the pattern
- [ ] `decider check adr --strict` passes

## Index File

The `docs/adr/index.yaml` file is auto-generated and provides quick access to ADR metadata:

```yaml
generated_at: "2026-01-16T10:00:00Z"
adr_count: 5
adrs:
  - adr_id: ADR-0001
    title: "Decision Title"
    status: adopted
    date: "2026-01-16"
    tags: [foundation]
    scope_paths: ["src/**"]
    file: "0001-decision-title.md"
```

Never edit this file manually. Always use `decider index` to regenerate it.
