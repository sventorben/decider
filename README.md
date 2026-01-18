# DECIDER

**Turn architectural decisions into machine-readable constraints that both humans and AI coding agents can follow.**

DECIDER is a Git-native system for managing Architecture Decision Records (ADRs) with structured metadata. It combines a fast CLI, an ADR Steward agent, and CI workflows to create a living decision system that scales with your team.

```bash
# Install (replace v0.1.0 with the latest release version)
go install github.com/sventorben/decider/cmd/decider@v0.1.0

# Initialize and create your first ADR
decider init
decider new "Use PostgreSQL for persistence" --tags database --paths "src/db/**"
```

> **⚠️ Early Draft / Exploratory Project**
>
> This repository is an early-stage draft of an idea and methodology that is still evolving.
> The concepts, structure, and tooling are intentionally experimental and subject to change.
>
> The goal at this stage is to explore and validate the approach in public, gather feedback,
> and iterate openly — not to provide a finished or production-hardened solution.
> 
> If you are evaluating this project, please treat it as a foundation and conversation starter rather than a final recommendation.

## Why DECIDER?

Software teams make hundreds of architectural decisions. Most get lost in Slack threads, PR comments, or tribal knowledge. When a new team member joins—or an AI agent starts writing code—they have no way to know what constraints exist.

DECIDER solves this by treating decisions as code:

- **Decisions become constraints**: Each ADR defines rules that MUST be followed
- **Constraints have scope**: Glob patterns specify which code paths they apply to
- **Agents can query them**: `decider check diff --base main` returns applicable constraints
- **CI can enforce them**: Fail builds when changes violate documented decisions

### The Problem with AI Coding Agents

AI agents write code fast, but they don't know your architecture. Without explicit constraints:

- An agent might use raw SQL when you've decided on the repository pattern
- It might add a new database when you've standardized on PostgreSQL
- It might violate security invariants documented nowhere but someone's memory

DECIDER gives agents the context they need to make decisions that align with your architecture.

## How It Works

```
┌─────────────────────────────────────────────────────────────────┐
│                        Your Repository                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  docs/adr/                        src/                          │
│  ├── index.yaml    ◄─────────────► ├── db/                      │
│  ├── 0001-postgresql.md           │   └── users.go              │
│  ├── 0002-rest-api.md             ├── api/                      │
│  └── templates/                   │   └── handlers.go           │
│                                   └── ...                       │
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   decider   │  │ ADR Steward │  │  CI Checks  │              │
│  │    CLI      │  │   Agent     │  │  (GitHub)   │              │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │
│         │                │                │                     │
│         └────────────────┴────────────────┘                     │
│                          │                                      │
│                   Read/Write ADRs                               │
│                   Validate/Index                                │
│                   Match to Changes                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**Components:**

1. **ADRs with YAML frontmatter**: Machine-readable metadata (constraints, invariants, scope paths)
2. **Index file**: Auto-generated `index.yaml` for fast ADR discovery
3. **CLI (`decider`)**: Create, validate, list, and query ADRs
4. **ADR Steward agent**: Claude Code integration for ADR lifecycle management
5. **CI workflows**: Validate ADRs and check index consistency on every PR

## Quickstart

### 1. Install

```bash
# Pinned install with checksum verification (recommended)
./scripts/install-decider.sh
export PATH="$PWD/tools/decider:$PATH"

# Or with Go (replace v0.1.0 with pinned version)
go install github.com/sventorben/decider/cmd/decider@v0.1.0
```

### 2. Initialize

```bash
decider init
# Created: docs/adr/
# Created: docs/adr/templates/adr.md
# Created: docs/adr/index.yaml
```

### 3. Create an ADR

```bash
decider new "Use PostgreSQL for persistence" \
  --tags database,storage \
  --paths "src/db/**,migrations/**"
```

This creates `docs/adr/0001-use-postgresql-for-persistence.md`:

```yaml
---
adr_id: ADR-0001
title: "Use PostgreSQL for persistence"
status: proposed
date: 2026-01-16
scope:
  paths:
    - "src/db/**"
    - "migrations/**"
tags:
  - database
  - storage
constraints: []
invariants: []
---

# ADR-0001: Use PostgreSQL for persistence

## Context
<!-- Why is this decision needed? -->

## Decision
<!-- What did we decide and why? -->

## Alternatives Considered
<!-- What other options were evaluated? -->

## Consequences
<!-- What are the positive and negative impacts? -->
```

### 4. Fill in the ADR

Edit the file to add your decision using the **mandatory rationale pattern**:

```markdown
## Decision

### PostgreSQL: Adopted

**Adopted because:**
- Mature ecosystem with proven reliability at scale
- Team has extensive PostgreSQL experience
- Excellent JSON support for semi-structured data

**Adopted despite:**
- Higher operational complexity than SQLite
- Requires dedicated infrastructure

## Alternatives Considered

### MongoDB: Rejected

**Rejected because:**
- Team lacks MongoDB expertise
- Query patterns are primarily relational

**Rejected despite:**
- Better horizontal scaling characteristics
- More flexible schema evolution
```

Also add constraints and invariants:

```yaml
constraints:
  - "All database access must go through the repository pattern"
  - "Use prepared statements for all queries"
  - "Never expose raw SQL outside the db package"
invariants:
  - "Database connections are pooled"
  - "All migrations are reversible"
```

The rationale pattern ensures every decision documents:
- **Why it was adopted** (concrete reasons tied to drivers)
- **What trade-offs were accepted** (despite sections)
- **Why alternatives were rejected** (with acknowledged strengths)

### 5. Query Applicable ADRs

When making changes, check which ADRs apply:

```bash
$ decider check diff --base main
ADR-0001: Use PostgreSQL for persistence
  Status: adopted
  Matched files:
    - src/db/users.go
  Constraints:
    - All database access must go through the repository pattern
    - Use prepared statements for all queries
```

## Agent Integration

### How Agents Use ADRs

AI coding agents can query DECIDER before making changes:

```bash
# Agent workflow: get constraints for current changes (TOON format, default)
$ decider check diff --base main --format toon
{applicable_adrs:[{adr_id:ADR-0001 constraints:["Use repository pattern" "Use prepared statements"]}]}

# Use JSON format for traditional tooling
$ decider check diff --base main --format json | jq '.applicable_adrs[].constraints'
[
  "All database access must go through the repository pattern",
  "Use prepared statements for all queries"
]

# Get detailed explanation of why ADRs apply
$ decider explain --base main
```

### Example: Agent Constrained by ADRs

User prompt:
> "Add a function to get user by email"

Without ADRs, an agent might write:

```go
// BAD: Raw SQL in handler
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE email = '" + email + "'")
    // SQL injection vulnerability, not using repository pattern
}
```

With DECIDER, the agent first checks constraints:

```bash
$ decider check diff --base main --path src/api/handlers.go
ADR-0001 applies. Constraints:
- All database access must go through the repository pattern
- Use prepared statements for all queries
```

The agent then writes compliant code:

```go
// GOOD: Uses repository pattern per ADR-0001
func (h *Handler) GetUserByEmail(ctx context.Context, email string) (*User, error) {
    return h.userRepo.FindByEmail(ctx, email)
}
```

## ADR Steward Agent

DECIDER includes an **ADR Steward** agent for Claude Code that manages the ADR lifecycle. The steward handles creation, supersession, validation, and index maintenance.

### Using the Steward

The ADR Steward is available through Claude Code slash commands:

```bash
# Create a new ADR
/adr-new Use Redis for session storage

# Supersede an existing ADR
/adr-supersede ADR-0003 Migrate from Redux to Zustand
```

You can also invoke the steward directly:

> "Create an ADR for our decision to use GraphQL instead of REST for the new API"

> "We need to supersede ADR-0005 because we're switching from Memcached to Redis"

> "Check if all our ADRs are valid and the index is up to date"

### What the Steward Does

| Action | Steward Behavior |
|--------|-----------------|
| **Create ADR** | Runs `decider new`, fills sections, sets scope paths, updates index |
| **Supersede ADR** | Creates new ADR, links both with `supersedes`/`superseded_by`, updates statuses |
| **Validate** | Runs `decider check adr`, reports issues, suggests fixes |
| **Update Index** | Runs `decider index` after any ADR changes |
| **Query** | Uses `decider check diff` to find applicable ADRs for changes |

### What Humans Decide

The steward assists but doesn't make architectural decisions:

- **Humans decide**: What to document, whether to accept/reject, constraint wording
- **Steward executes**: File creation, format compliance, linking, index updates

## How Teams Use This

### Sprint Workflow

1. **Planning**: Create ADRs for significant decisions (`decider new` or `/adr-new`)
2. **Development**: Agents and devs check constraints before coding (`decider check diff`)
3. **Review**: PRs include ADR validation in CI
4. **Evolution**: Supersede outdated decisions (`/adr-supersede`)

### Adoption Path

| Phase | What You Do | Benefit |
|-------|-------------|---------|
| **1. Docs only** | Write ADRs, no tooling | Decisions documented |
| **2. Steward agent** | Use Claude Code integration | Consistent format, auto-indexing |
| **3. CI checks** | Add validation to pipelines | Catch drift, enforce structure |
| **4. Policy enforcement** | (Future) Block PRs violating constraints | Automated governance |

## CLI Reference

| Command | Description |
|---------|-------------|
| `decider init` | Initialize ADR directory structure |
| `decider new "<title>"` | Create a new ADR |
| `decider list` | List ADRs with optional filters |
| `decider show <id>` | Display ADR details |
| `decider check adr` | Validate all ADRs |
| `decider check adr --strict` | Validate ADRs (fail on missing rationale pattern) |
| `decider check diff --base <ref>` | Find ADRs applicable to changes |
| `decider explain --base <ref>` | Explain why ADRs apply |
| `decider index` | Regenerate the ADR index |
| `decider version` | Show version info |

Common flags:
- `--dir PATH`: ADR directory (default: `docs/adr`)
- `--format text|toon|json`: Output format (TOON is default for structured output)
- `--status STATUS`: Filter by status
- `--tag TAG`: Filter by tag
- `--strict`: Treat warnings as errors (for `check adr`)

### Output Formats

DECIDER supports multiple output formats:
- **text**: Human-readable (default for CLI)
- **toon**: Token-Oriented Object Notation, compact and LLM-friendly (default for structured output)
- **json**: Standard JSON for traditional tooling integration

Use `--format=json` when integrating with tools that require JSON.

See [SPEC.md](SPEC.md) for complete CLI documentation.

## Installation

### Pinned Install (Recommended)

Use the installer scripts for a pinned, checksum-verified install:

```bash
# macOS/Linux
./scripts/install-decider.sh

# Windows PowerShell
.\scripts\install-decider.ps1
```

This downloads the version pinned in `tools/decider.version`, verifies the checksum, and installs to `./tools/decider/`.

### From GitHub Releases (Manual)

```bash
# Linux (amd64) - replace 0.1.0 with pinned version
curl -L https://github.com/sventorben/decider/releases/download/v0.1.0/decider_0.1.0_linux_amd64.tar.gz | tar xz

# macOS (arm64) - replace 0.1.0 with pinned version
curl -L https://github.com/sventorben/decider/releases/download/v0.1.0/decider_0.1.0_darwin_arm64.tar.gz | tar xz

# Windows - download zip from releases page
```

### From Source

```bash
go install github.com/sventorben/decider/cmd/decider@v0.1.0
```

### Verify Installation

```bash
decider version
# decider version 0.1.0
#   commit: abc123
#   built:  2026-01-16T10:00:00Z
```

## Documentation

| Document | Purpose |
|----------|---------|
| [docs/](docs/README.md) | Documentation portal |
| [Adopting DECIDER](docs/guides/adopting-decider-in-a-team.md) | Team rollout guide |
| [Using ADR Steward](docs/guides/using-adr-steward.md) | Claude Code agent guide |
| [Why: Decisions as Constraints](docs/why/decisions-as-constraints.md) | ADRs as enforceable rules |
| [Why: Shared Context](docs/why/shared-context-beats-bigger-models.md) | Context over model capability |
| [SPEC.md](SPEC.md) | Complete specification |
| [AGENTS.md](AGENTS.md) | Agent integration guide |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contribution guidelines |
| [demo/](demo/README.md) | Interactive demo |

## Non-Goals

DECIDER intentionally does not:

- **Enforce constraints semantically**: The CLI surfaces constraints; enforcement is your CI/human review
- **Manage approvals**: No built-in workflow for ADR acceptance; use your existing PR process
- **Support remote ADR repos**: ADRs live in your project repository
- **Provide policy engine**: No runtime enforcement; this is documentation with structure

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines. Run `make check` before opening a PR.

## AI Assistance

Parts of this project were developed with AI-assisted tooling. Contributors and users should review code and documentation critically, as with any project. All contributions—whether human or AI-assisted—must meet the same project standards: tests must pass, ADRs must be followed, and human review is expected.
