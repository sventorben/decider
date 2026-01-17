# Getting Started with DECIDER

This guide walks you through installing DECIDER and creating your first ADR.

## Prerequisites

- Git repository (DECIDER is designed for Git-managed projects)
- Go 1.25+ (if installing from source)

## Installation

### Option 1: Download Binary

```bash
# Linux (amd64) - replace v0.1.0 with the latest release version
curl -L https://github.com/sventorben/decider/releases/download/v0.1.0/decider_Linux_x86_64.tar.gz | tar xz
sudo mv decider /usr/local/bin/

# macOS (Intel) - replace v0.1.0 with the latest release version
curl -L https://github.com/sventorben/decider/releases/download/v0.1.0/decider_Darwin_x86_64.tar.gz | tar xz
sudo mv decider /usr/local/bin/

# macOS (Apple Silicon) - replace v0.1.0 with the latest release version
curl -L https://github.com/sventorben/decider/releases/download/v0.1.0/decider_Darwin_arm64.tar.gz | tar xz
sudo mv decider /usr/local/bin/
```

### Option 2: Go Install

```bash
go install github.com/sventorben/decider/cmd/decider@v0.1.0
```

### Verify Installation

```bash
decider version
```

## Initialize Your Repository

Navigate to your project root and run:

```bash
decider init
```

This creates:
```
docs/adr/
├── templates/
│   └── adr.md          # Template for new ADRs
└── index.yaml          # Auto-generated index (empty)
```

## Create Your First ADR

```bash
decider new "Use PostgreSQL for persistence" \
  --tags database,storage \
  --paths "src/db/**,migrations/**"
```

Output:
```
Created: docs/adr/0001-use-postgresql-for-persistence.md
Updated: docs/adr/index.yaml
```

## Understand the ADR Structure

Open the generated file:

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
supersedes: []
superseded_by: []
related_adrs: []
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

### Key Fields

| Field | Purpose |
|-------|---------|
| `adr_id` | Unique identifier (ADR-NNNN format) |
| `status` | Lifecycle state (proposed, adopted, rejected, deprecated, superseded) |
| `scope.paths` | Glob patterns for affected code |
| `constraints` | Rules that MUST be followed |
| `invariants` | Properties that must always hold |

## Fill In Your ADR

Edit the file to document your decision:

```yaml
constraints:
  - "All database access must go through the repository pattern"
  - "Use prepared statements for all queries"
invariants:
  - "Database connections are pooled"
  - "All migrations are reversible"
```

Then fill in the markdown sections:

```markdown
## Context

We need a reliable database for our e-commerce application that handles:
- User accounts and authentication
- Product catalog with inventory
- Order history and transactions

We require ACID compliance for financial transactions.

## Decision

We will use PostgreSQL 15+ as our primary database.

All database access will follow the repository pattern:
- SQL queries encapsulated in repository structs
- Business logic uses interfaces, not raw SQL
- Prepared statements prevent SQL injection

## Alternatives Considered

| Alternative | Pros | Cons |
|-------------|------|------|
| MySQL | Widely known | Less feature-rich |
| MongoDB | Flexible schema | No ACID for transactions |

## Consequences

**Positive:**
- ACID compliance for orders
- Rich feature set (JSON columns, CTEs)

**Negative:**
- Requires operational knowledge
```

## Validate Your ADR

```bash
decider check adr
```

Expected output:
```
Validating ADRs in docs/adr/...
✓ 0001-use-postgresql-for-persistence.md: valid
All ADRs valid (1 checked)
```

## Update the Index

The index updates automatically when you create ADRs, but you can regenerate it:

```bash
decider index
```

## Check ADRs Applicable to Changes

When you've made changes to your code, find relevant ADRs:

```bash
decider check diff --base main
```

This compares your current branch to `main` and lists ADRs whose `scope.paths` match changed files.

## Next Steps

- [Writing Effective ADRs](writing-adrs.md) - Learn to write high-quality ADRs
- [CI Integration](ci-integration.md) - Set up validation in GitHub Actions
- [Agent Integration](agent-integration.md) - Use the ADR Steward

## Common Commands Cheat Sheet

```bash
# Initialize
decider init

# Create ADR
decider new "Title" --tags tag1,tag2 --paths "path/**"

# List all ADRs
decider list

# List by status
decider list --status adopted

# Show ADR details
decider show ADR-0001

# Validate all ADRs
decider check adr

# Find ADRs for changes
decider check diff --base main

# Explain why ADRs apply
decider explain --base main

# Update index
decider index

# Check if index is current
decider index --check
```
