# DECIDER Documentation

Welcome to the DECIDER documentation. This is your guide to turning architectural decisions into machine-readable constraints.

## Quick Links

| I want to... | Go to |
|--------------|-------|
| Get started quickly | [Getting Started Guide](guides/getting-started.md) |
| Write effective ADRs | [Writing ADRs Guide](guides/writing-adrs.md) |
| Roll out DECIDER to my team | [Adopting DECIDER](guides/adopting-decider-in-a-team.md) |
| Set up CI validation | [CI Integration Guide](guides/ci-integration.md) |
| Use the ADR Steward agent | [Using ADR Steward](guides/using-adr-steward.md) |
| Connect AI agents to ADRs | [Agent Integration Guide](guides/agent-integration.md) |
| Understand the philosophy | [Why DECIDER?](why/) |
| Look up CLI commands | [Reference](reference/) |
| See the full specification | [SPEC.md](../SPEC.md) |

## Documentation Structure

```
docs/
├── README.md                    # This file
├── guides/                      # How-to guides
│   ├── README.md                # Guide index
│   ├── getting-started.md       # First steps with DECIDER
│   ├── writing-adrs.md          # Crafting effective ADRs
│   ├── adopting-decider-in-a-team.md  # Team rollout strategy
│   ├── ci-integration.md        # GitHub Actions setup
│   ├── using-adr-steward.md     # ADR Steward agent guide
│   └── agent-integration.md     # Connecting AI agents
├── why/                         # Philosophy and rationale
│   ├── README.md                # Overview
│   ├── decisions-as-code.md     # Why treat decisions as code
│   ├── decisions-as-constraints.md    # ADRs as enforceable rules
│   ├── shared-context-beats-bigger-models.md  # Context over capability
│   └── ai-agents-need-constraints.md  # Why agents need structure
├── adr/                         # ADR system documentation
│   ├── README.md                # ADR overview and conventions
│   ├── templates/               # ADR templates
│   └── *.md                     # Individual ADRs
├── reference/                   # Reference documentation
│   └── README.md                # CLI and schema reference
└── security/                    # Security documentation
    └── SECURITY_REVIEW_v0.1.md  # Security review artifacts
```

## Core Concepts

### ADRs (Architecture Decision Records)

An ADR documents a significant decision along with its context and consequences. In DECIDER, ADRs include:

- **Metadata**: Machine-readable YAML frontmatter
- **Constraints**: Rules that MUST be followed
- **Invariants**: Properties that must always hold
- **Scope**: Glob patterns defining where the decision applies

### Index

The `index.yaml` file is auto-generated and provides fast access to ADR metadata without parsing every file. Never edit it manually.

### Steward Agent

The ADR Steward is a Claude Code agent that manages ADR lifecycle: creation, supersession, validation, and indexing.

## The System at a Glance

```
Decisions → Constraints → Scope → Enforcement → Learning
    │            │          │          │            │
    │            │          │          │            └─ Team improves process
    │            │          │          └─ CI checks, human review
    │            │          └─ Glob patterns match code paths
    │            └─ Rules agents/humans must follow
    └─ ADRs document what and why
```

## Getting Help

- [GitHub Issues](https://github.com/sventorben/decider/issues) - Bug reports and feature requests
- [CONTRIBUTING.md](../CONTRIBUTING.md) - How to contribute
