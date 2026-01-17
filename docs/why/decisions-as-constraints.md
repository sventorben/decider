# Decisions as Constraints

Traditional ADRs are documentation. You write them, file them, and hope someone reads them later. DECIDER treats ADRs differently: as constraints that actively shape development.

## The Documentation Trap

Most teams that adopt ADRs follow this pattern:

1. Make an important architectural decision
2. Write an ADR explaining the decision
3. Store it in a wiki or docs folder
4. Never look at it again

The ADR becomes a historical artifact—useful for archaeology, useless for guidance.

## From Records to Rules

DECIDER shifts the model. An ADR isn't just a record of what was decided. It's a specification of what must be followed.

This shift happens through three mechanisms:

### 1. Constraints

Constraints are explicit rules embedded in ADR frontmatter:

```yaml
constraints:
  - "All API endpoints must validate input before processing"
  - "Error responses must follow RFC 7807 Problem Details format"
  - "Authentication tokens must be validated on every request"
```

These aren't suggestions. They're requirements that apply to code within the ADR's scope.

### 2. Invariants

Invariants describe properties that must always hold:

```yaml
invariants:
  - "User sessions expire after 24 hours of inactivity"
  - "All personally identifiable information is encrypted at rest"
  - "Rate limits are enforced per API key"
```

Constraints tell you what to do. Invariants tell you what must remain true.

### 3. Scope Paths

Scope paths define where constraints apply:

```yaml
scope:
  paths:
    - "src/api/**"
    - "internal/handlers/**"
```

A constraint without scope is just a guideline. A constraint with scope is a rule for specific code.

## Before and After

### Before: Traditional ADR

```markdown
# ADR-007: Use Repository Pattern for Data Access

## Status
Adopted

## Context
We need a consistent approach to data access...

## Decision
We will use the repository pattern...

## Consequences
This adds a layer of abstraction...
```

This ADR is informative but passive. Nothing connects it to code. Nothing enforces it.

### After: DECIDER ADR

```yaml
---
adr_id: ADR-0007
title: "Use Repository Pattern for Data Access"
status: adopted
date: 2026-01-15
scope:
  paths:
    - "src/db/**"
    - "internal/storage/**"
constraints:
  - "All database queries must go through repository interfaces"
  - "Repositories must not expose database-specific types in their interfaces"
  - "Direct SQL queries outside repositories are prohibited"
invariants:
  - "Repository interfaces are defined in the domain layer"
  - "Repository implementations are in the infrastructure layer"
---

# ADR-0007: Use Repository Pattern for Data Access

## Context
...
```

Now when someone modifies `src/db/users.go`:

```bash
$ decider check diff --base main
ADR-0007: Use Repository Pattern for Data Access
  Matched: src/db/**
  Constraints:
    - All database queries must go through repository interfaces
    - Repositories must not expose database-specific types
    - Direct SQL queries outside repositories are prohibited
```

The constraint surfaces at the moment it matters.

## Enforcement Levels

DECIDER doesn't automatically block violations. Instead, it supports progressive enforcement:

| Level | Mechanism | Effect |
|-------|-----------|--------|
| **Advisory** | `decider check diff` in development | Developer sees applicable constraints |
| **Visible** | `decider check diff` in PR comments | Reviewers see constraints |
| **Required** | CI check with exit code | PR blocked until reviewed |
| **Automated** | Future: semantic analysis | Violations detected automatically |

Most teams start at advisory and move up as they build confidence.

## Writing Effective Constraints

Good constraints are:

**Specific.** "Code should be well-structured" is useless. "All public functions must have a single return statement" is verifiable.

**Actionable.** The reader should know exactly what to do. "Consider performance" fails. "Database queries must use indexed columns" succeeds.

**Scoped.** Apply constraints to the narrowest reasonable scope. Global constraints become noise.

**Testable.** If you can't tell whether code violates a constraint, the constraint is too vague.

## The Rationale Pattern

DECIDER requires explicit rationale for decisions:

```markdown
### Repository Pattern: Adopted

**Adopted because:**
- Isolates business logic from data access details
- Enables testing with mock repositories
- Allows switching storage backends without domain changes

**Adopted despite:**
- Adds indirection that increases initial complexity
- Requires discipline to avoid leaking implementation details
```

This forces authors to acknowledge trade-offs, making constraints more credible and easier to revisit later.

## When Constraints Conflict

Sometimes ADRs have overlapping scopes with different constraints. When this happens:

1. More specific scope wins over broader scope
2. Later decisions can supersede earlier ones
3. Explicit `supersedes` links clarify the relationship

If conflicts remain ambiguous, that's a signal to refine scope paths or consolidate ADRs.

## Further Reading

- [Shared Context Beats Bigger Models](shared-context-beats-bigger-models.md) - Why externalized decisions matter
- [Adopting DECIDER in a Team](../guides/adopting-decider-in-a-team.md) - Practical rollout guidance
