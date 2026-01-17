# Guides

Practical guidance for using DECIDER effectively.

## By Intent

| I want to... | Read this |
|--------------|-----------|
| Install and run my first command | [Getting Started](getting-started.md) |
| Write ADRs that actually get followed | [Writing ADRs](writing-adrs.md) |
| Roll out DECIDER to my team | [Adopting DECIDER in a Team](adopting-decider-in-a-team.md) |
| Add validation to CI pipelines | [CI Integration](ci-integration.md) |
| Use the ADR Steward agent | [Using ADR Steward](using-adr-steward.md) |
| Connect AI agents to ADRs | [Agent Integration](agent-integration.md) |

## Guide Index

### Getting Started
First-time setup: installation, initialization, creating your first ADR.

### Writing ADRs
How to write ADRs that are specific, scoped, and actionable. Covers constraints, invariants, and the rationale pattern.

### Adopting DECIDER in a Team
Progressive adoption path from documentation-only to full CI enforcement. Addresses common concerns and anti-patterns.

### CI Integration
GitHub Actions workflow for validating ADRs and checking index consistency on every pull request.

### Using ADR Steward
How to use the Claude Code agent for ADR lifecycle management: creation, supersession, validation.

### Agent Integration
Connecting AI coding agents to DECIDER so they can query and respect architectural constraints.

## Prerequisites

All guides assume:
- Go 1.25+ installed (for building from source)
- Git repository initialized
- Basic familiarity with command-line tools

## Related Documentation

- [Why DECIDER?](../why/) - Philosophy and rationale
- [Reference](../reference/) - CLI command reference
- [SPEC.md](../../SPEC.md) - Complete specification
