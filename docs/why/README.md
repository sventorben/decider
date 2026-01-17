# Why DECIDER?

This section explains the philosophy behind DECIDER and why treating architectural decisions as structured, machine-readable artifacts matters.

## Essays

- [Decisions as Code](decisions-as-code.md) - Why architectural decisions deserve the same rigor as source code
- [Decisions as Constraints](decisions-as-constraints.md) - How ADRs become enforceable rules, not just documentation
- [Shared Context Beats Bigger Models](shared-context-beats-bigger-models.md) - Why externalized context matters more than model capability
- [AI Agents Need Constraints](ai-agents-need-constraints.md) - How structure helps AI coding agents make better decisions

## The Core Insight

Software architecture isn't about the code you write today. It's about the constraints that shape every future change.

Most teams document decisions in wikis, Confluence pages, or meeting notes. These documents:

- Get stale as the codebase evolves
- Aren't connected to the code they govern
- Can't be queried programmatically
- Don't help AI agents understand your architecture

DECIDER changes this by treating decisions as code:

```
Decisions → Constraints → Scope → Enforcement → Learning
```

1. **Decisions** are documented in ADRs with structured metadata
2. **Constraints** define rules that must be followed
3. **Scope** specifies which code paths the constraints apply to
4. **Enforcement** happens through CI and human review
5. **Learning** improves the process over time

## What Makes DECIDER Different

| Traditional ADRs | DECIDER ADRs |
|-----------------|--------------|
| Free-form markdown | Structured YAML frontmatter |
| No tooling | CLI for creation, validation, querying |
| Manual maintenance | Auto-generated index |
| Human-only consumption | Machine-readable for agents |
| No scope definition | Glob patterns specify affected code |
| Static documents | Living system with CI integration |

## The Methodology

DECIDER isn't just a CLI. It's a lightweight methodology:

1. **Document decisions early** - Create ADRs during planning
2. **Make constraints explicit** - Write rules agents and humans can follow
3. **Scope decisions tightly** - Use glob patterns to limit applicability
4. **Validate continuously** - CI checks keep ADRs healthy
5. **Evolve through supersession** - Don't edit old decisions; create new ones

This creates a living record of your architecture that grows with your codebase.
