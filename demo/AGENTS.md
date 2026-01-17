# AGENTS.md — Instructions for AI Coding Agents

This document tells AI agents how to work with this demo project.

## Overview

This demo uses **DECIDER** to manage Architecture Decision Records (ADRs). ADRs document architectural decisions with machine-readable constraints that you MUST follow.

## Install decider (if missing)

If `decider` is not available, install the pinned version:

```bash
# macOS/Linux
./scripts/install-decider.sh

# Windows PowerShell
.\scripts\install-decider.ps1
```

This installs to `tools/decider/`. Add to PATH:

```bash
export PATH="$PWD/tools/decider:$PATH"
```

## Agent Pre-flight (Required)

Before starting significant work, run:

```bash
decider check adr --dir docs/adr
```

When touching ADR system files (`docs/adr/`, templates, index), use strict mode:

```bash
decider check adr --dir docs/adr --strict
```

**Stop on violations**: If ADR checks fail, fix the issues before proceeding.

**Re-run before finishing**: After completing work, re-run the check to ensure compliance.

## Before Making Changes

### 1. Query Applicable ADRs

Before modifying code, find which ADRs apply to the files you're changing:

```bash
decider list --dir docs/adr --path "src/db/users.go"
```

### 2. Read Constraints and Invariants

Each applicable ADR contains:
- **Constraints**: Rules you MUST follow
- **Invariants**: Properties that must always be true

Example:
```
ADR-0001: Use PostgreSQL for Persistence
  Constraints:
    - All database access must go through the repository pattern
    - Use prepared statements for all queries
    - Never expose raw SQL outside the db package
```

### 3. Follow the Rules

Treat constraints as hard requirements. If you're unsure how to comply, ask for clarification.

### 4. Verify Compliance

After making changes, verify compliance:

```bash
# Find ADRs that apply to your changed files
decider list --dir docs/adr --path "path/to/changed/file"

# Validate all ADRs are properly formatted
decider check adr --dir docs/adr --strict
```

## Demo Project Structure

```
demo/
├── AGENTS.md              # This file
├── CLAUDE.md              # Claude-specific instructions
├── README.md              # Demo walkthrough
├── scripts/
│   ├── install-decider.sh
│   └── install-decider.ps1
├── tools/
│   ├── decider.version    # Pinned version (v0.1.0)
│   └── decider/           # Installed binary
├── docs/adr/
│   ├── README.md          # ADR system overview
│   ├── index.yaml         # Auto-generated index
│   ├── templates/adr.md   # ADR template
│   └── *.md               # Individual ADRs
└── src/                   # Demo source directories
```

## ADR Workflow

### Reading ADRs

```bash
# List all ADRs
decider list --dir docs/adr

# Show details of a specific ADR
decider show ADR-0001 --dir docs/adr

# Find ADRs by tag
decider list --dir docs/adr --tag database

# Find ADRs by scope path
decider list --dir docs/adr --path "src/db/**"
```

### Validating ADRs

```bash
# Basic validation
decider check adr --dir docs/adr

# Strict validation (includes rationale pattern check)
decider check adr --dir docs/adr --strict
```

### Creating ADRs (when needed)

```bash
decider new "Your Decision Title" --dir docs/adr --tags tag1,tag2 --paths "affected/**"
```

Then edit the generated file to fill in Context, Decision (with rationale pattern), Alternatives, and Consequences.

## Mandatory Rationale Pattern

All ADRs in this project MUST use this pattern:

### For Adopted Options
```markdown
### [Option]: Adopted

**Adopted because:**
- Concrete reason tied to decision drivers

**Adopted despite:**
- Known trade-off or downside
```

### For Rejected Alternatives
```markdown
### [Alternative]: Rejected

**Rejected because:**
- Concrete reason for rejection

**Rejected despite:**
- Legitimate strength of this option
```

## Useful Commands

```bash
# Install decider
./scripts/install-decider.sh

# Validate all ADRs
decider check adr --dir docs/adr --strict

# List ADRs
decider list --dir docs/adr

# Show ADR details
decider show ADR-0001 --dir docs/adr

# Find ADRs for a file
decider list --dir docs/adr --path "src/db/users.go"

# Update index
decider index --dir docs/adr
```
