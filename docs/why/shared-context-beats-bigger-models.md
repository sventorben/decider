# Shared Context Beats Bigger Models

When teams adopt AI coding assistants, the natural instinct is to reach for the most capable model available. But model capability is only half the equation. The other half—often overlooked—is context.

## The Context Problem

A powerful model with no context about your system will:

- Suggest patterns your team has explicitly rejected
- Introduce dependencies you've standardized against
- Violate architectural boundaries nobody told it about
- Repeat mistakes the team learned from months ago

No amount of model intelligence compensates for missing information. A smaller model with the right context will outperform a larger model working blind.

## What "Context" Actually Means

Context isn't just "more text in the prompt." Useful context is:

1. **Relevant** - Applies to the current task
2. **Structured** - Can be parsed and reasoned about
3. **Authoritative** - Represents team decisions, not suggestions
4. **Scoped** - Doesn't overwhelm with irrelevant information

This is what DECIDER provides through ADRs with structured metadata.

## How DECIDER Externalizes Context

Traditional documentation fails because it's disconnected from code and unstructured for machines. DECIDER addresses both:

### Structured Metadata

```yaml
constraints:
  - "All database access must go through the repository pattern"
  - "Use prepared statements for all queries"
invariants:
  - "Database connections are pooled"
scope:
  paths:
    - "src/db/**"
```

An agent can parse this directly. No interpretation required.

### Scoped Relevance

When an agent works on `src/db/users.go`, it can query:

```bash
decider check diff --base main
```

This returns only the ADRs whose scope patterns match the changed files. The agent gets precisely the constraints that matter, not a dump of every decision ever made.

### The Index as a Context Layer

The `index.yaml` file provides fast access to ADR metadata:

```yaml
adrs:
  - adr_id: ADR-0001
    title: "Use PostgreSQL for persistence"
    constraints: ["All database access must go through repository pattern"]
    scope_paths: ["src/db/**"]
```

Agents can scan this index to understand which decisions exist before diving into details.

## Teams Benefit Too

This isn't just about AI. Human developers face the same context problem:

- New team members don't know historical decisions
- Senior engineers forget decisions made years ago
- Cross-team collaboration requires shared understanding

When decisions are externalized and structured, everyone—human or AI—works from the same source of truth.

## The Compound Effect

Context improves over time:

1. Team makes a decision → documents it as an ADR
2. Agent encounters similar situation → finds the ADR
3. Agent follows the constraint → produces consistent code
4. Team reviews and merges → validates the system works
5. Next decision builds on previous ones → context grows

Each documented decision makes future work easier. This compounds. A team with 50 well-scoped ADRs has a significant advantage over a team with none, regardless of which model they use.

## Practical Implications

**Start small.** Document decisions as you make them. Don't backfill everything at once.

**Scope tightly.** An ADR that applies to `**/*` is nearly useless. Specific scope paths make constraints actionable.

**Keep constraints concrete.** "Code should be clean" helps nobody. "All public functions must have doc comments" can be followed.

**Trust the system.** When you document a constraint, expect it to be applied. If agents or humans consistently violate it, either the constraint is wrong or the enforcement is missing.

## Further Reading

- [Decisions as Constraints](decisions-as-constraints.md) - How ADRs become enforceable rules
- [AI Agents Need Constraints](ai-agents-need-constraints.md) - Why structure matters for AI coding
