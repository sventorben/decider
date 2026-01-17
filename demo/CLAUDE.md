# CLAUDE.md — Claude Code Instructions

Instructions for Claude Code when working in this demo project.

**Source of truth**: See [AGENTS.md](AGENTS.md) for the complete agent workflow. This file provides Claude-specific commands.

## Install decider (if missing)

If `decider` is not available, install the pinned version:

```bash
# macOS/Linux
./scripts/install-decider.sh
export PATH="$PWD/tools/decider:$PATH"

# Windows PowerShell
.\scripts\install-decider.ps1
$env:PATH = "$PWD\tools\decider;$env:PATH"
```

## Daily Pre-flight (Required)

Before starting work, run:

```bash
decider check adr --dir docs/adr
```

When touching `docs/adr/`, templates, or index files, use strict mode:

```bash
decider check adr --dir docs/adr --strict
```

**Do not proceed if ADR checks fail.** Fix issues first.

## Before Making Changes

1. Find applicable ADRs for files you'll modify:
   ```bash
   decider list --dir docs/adr --path "src/db/users.go"
   ```

2. Read the constraints and invariants in applicable ADRs

3. Ensure your changes comply with all constraints

## After Making Changes

Re-run the ADR check:

```bash
decider check adr --dir docs/adr --strict
```

## Commands Reference

```bash
# Validate ADRs
decider check adr --dir docs/adr --strict

# List all ADRs
decider list --dir docs/adr

# Show ADR details
decider show ADR-0001 --dir docs/adr

# Find ADRs by path
decider list --dir docs/adr --path "src/api/**"

# Find ADRs by tag
decider list --dir docs/adr --tag database
```

## What NOT To Do

- Do not proceed if `decider check adr` fails
- Do not ignore ADR constraints
- Do not edit `index.yaml` manually (use `decider index --dir docs/adr`)
- Do not skip the pre-flight check
