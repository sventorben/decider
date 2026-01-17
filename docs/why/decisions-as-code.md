# Decisions as Code

Why architectural decisions deserve the same rigor as source code.

## The Problem with Tribal Knowledge

Every software team accumulates decisions:

- "We use PostgreSQL because MongoDB didn't handle our transaction requirements"
- "All API endpoints must be versioned because we broke mobile clients once"
- "Authentication goes through the gateway, never directly to services"

These decisions live in:
- Slack threads that scroll away
- Meeting notes nobody reads
- The memories of engineers who might leave
- Code comments that drift from reality

When a new engineer joins, they learn through osmosis or painful mistakes. When an AI agent writes code, it has no access to this context at all.

## Code Has Structure. Why Don't Decisions?

Source code is:
- **Version controlled**: Every change is tracked
- **Validated**: Compilers and linters catch errors
- **Tested**: CI ensures it works
- **Scoped**: Files and modules have clear boundaries
- **Queryable**: IDEs help you navigate

Architectural decisions should have the same properties.

## What "Decisions as Code" Means

In DECIDER, a decision isn't just documentation. It's structured data:

```yaml
---
adr_id: ADR-0001
title: "Use PostgreSQL for persistence"
status: adopted
date: 2026-01-16
scope:
  paths:
    - "src/db/**"
    - "migrations/**"
constraints:
  - "All database access must go through the repository pattern"
  - "Use prepared statements for all queries"
invariants:
  - "Database connections are pooled"
  - "All migrations are reversible"
---
```

This structure enables:

### Version Control
ADRs live in your Git repository. Changes are tracked, reviewed in PRs, and tied to the commits that implement them.

### Validation
`decider check adr` validates format, required fields, and cross-references. CI fails if ADRs are malformed.

### Testing
`decider index --check` ensures the index matches reality. Stale metadata is caught automatically.

### Scope
Glob patterns in `scope.paths` define exactly which code an ADR governs. No ambiguity about where rules apply.

### Queryability
`decider check diff --base main` tells you which ADRs apply to your changes. Agents can query constraints programmatically.

## The Lifecycle

Decisions evolve. DECIDER models this explicitly:

```
proposed → adopted → deprecated/superseded
                ↓
            rejected
```

You don't edit old decisions to reflect new thinking. You supersede them:

```yaml
# Old ADR
status: superseded
superseded_by:
  - ADR-0015

# New ADR
status: adopted
supersedes:
  - ADR-0005
```

This preserves history. You can always trace why things changed.

## Why This Matters

### For Teams

- **Onboarding**: New engineers read ADRs to understand architecture
- **Consistency**: Everyone follows the same constraints
- **Evolution**: Decisions can be revisited and superseded cleanly
- **Accountability**: PRs that violate ADRs get flagged

### For AI Agents

- **Context**: Agents can query applicable constraints before writing code
- **Compliance**: Constraints are specific enough to follow
- **Learning**: Agents can reference ADRs to understand why code is structured a certain way

### For the Codebase

- **Alignment**: Code reflects documented decisions
- **Maintainability**: Future changes respect established patterns
- **Quality**: Constraints encode hard-won lessons

## The Alternative

Without structured decisions:

- Every PR is a negotiation about patterns
- Agents write code that violates undocumented rules
- New engineers make avoidable mistakes
- Technical debt accumulates as decisions drift

With decisions as code:

- Patterns are explicit and queryable
- Agents follow documented constraints
- Engineers understand the "why" behind the "what"
- Evolution is deliberate, not accidental

## Getting Started

You don't need to document everything at once. Start with:

1. **One decision that caused pain**: The rule everyone learns the hard way
2. **One external dependency**: Why you chose it and how to use it
3. **One security constraint**: The thing that must never be violated

Then grow from there. The index stays current. The system stays healthy. And your architecture becomes as explicit as your code.
