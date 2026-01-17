# AGENTS.md — Instructions for AI Coding Agents

This document tells AI agents how to work with this repository.

## Overview

This repository uses **DECIDER** to manage Architecture Decision Records (ADRs). ADRs document significant architectural decisions with machine-readable constraints that you MUST follow.

## Before Making Changes

### 1. Query Applicable ADRs

Before writing or modifying code, check which ADRs apply:

```bash
decider check diff --base main
```

Or for a detailed explanation:

```bash
decider explain --base main
```

### 2. Read Constraints and Invariants

Each applicable ADR contains:
- **Constraints**: Rules you MUST follow
- **Invariants**: Properties that must always be true

Example output:
```
ADR-0001: Adopt Go for DECIDER CLI
  Constraints:
    - Use Go 1.25 or later for generics and improved stdlib
    - Keep external dependencies minimal
    - Prefer stdlib flag package over heavy CLI frameworks
  Invariants:
    - All code compiles with `go build ./...`
    - All tests pass with `go test ./...`
```

### 3. Follow the Rules

Treat constraints as hard requirements. If you're unsure how to comply, ask for clarification rather than guessing.

### 4. Before Finalizing Changes

Before completing any change, run the ADR compliance loop:

```bash
decider check diff --base <ref>
```

For each applicable ADR:
1. Read its constraints and invariants
2. Verify your changes satisfy them
3. If violations exist, fix them and re-check
4. If compliance requires changing the ADR, propose an update instead of forcing non-compliant code

Only finalize when all applicable constraints are satisfied.

## Project Setup

### Install decider CLI

If `decider` is not available, install the pinned version:

```bash
# macOS/Linux
./scripts/install-decider.sh

# Windows
.\scripts\install-decider.ps1
```

This installs to `./tools/decider/`. Add to PATH:

```bash
export PATH="$PWD/tools/decider:$PATH"
```

### Build from Source

```bash
go build -o decider ./cmd/decider
```

### Test

```bash
go test ./...
```

### Lint

```bash
golangci-lint run
```

### Validate ADRs

```bash
decider check adr
```

## Project Structure

```
decider/
├── cmd/decider/          # CLI entry point
├── internal/
│   ├── adr/              # ADR parsing and validation
│   ├── cli/              # Command implementations
│   ├── glob/             # Glob pattern matching
│   ├── index/            # Index generation
│   └── validate/         # Input validation
├── docs/
│   ├── adr/              # Project ADRs
│   ├── guides/           # How-to documentation
│   ├── reference/        # CLI reference
│   └── why/              # Philosophy
├── demo/                 # Demo files
└── .claude/              # Claude Code integration
```

## Coding Standards

### Go Code

- Format with `gofmt` or `go fmt`
- Check with `go vet ./...`
- Follow standard Go conventions
- Prefer returning errors over panicking
- Keep functions small and testable

### Commits

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(cli): add new command
fix(adr): handle edge case in parsing
docs(readme): update quickstart
test(index): add coverage for empty index
```

### Error Handling

- Wrap errors with context
- Use meaningful error messages
- Return errors, don't panic

## When to Create ADRs

Create an ADR when a decision:

- Affects system architecture
- Has long-term implications
- Involves significant tradeoffs
- Isn't obvious to future developers

Use the DECIDER CLI:

```bash
decider new "Your Decision Title" --tags tag1,tag2 --paths "affected/**"
```

Then fill in the Context, Decision, Alternatives, and Consequences sections.

## ADR Workflow

### Creating

1. Run `decider new "Title" --tags ... --paths ...`
2. Edit the generated file to fill in all sections
3. Add specific constraints and invariants
4. Run `decider check adr` to validate
5. Run `decider index` to update the index

### Superseding

When a decision changes, don't edit the old ADR. Supersede it:

1. Create a new ADR explaining the change
2. Update the old ADR: `status: superseded`, add to `superseded_by`
3. Update the new ADR: add old ADR to `supersedes`
4. Run `decider check adr` and `decider index`

## Guardrails

- **No breaking API changes** without an ADR and migration plan
- **No secrets** in the repository
- **Security-sensitive changes** require human review
- **All PRs** must pass CI (tests, lint, ADR validation)

## Useful Commands

```bash
# List all ADRs
decider list

# List adopted ADRs
decider list --status adopted

# Show ADR details
decider show ADR-0001

# Find ADRs for your changes
decider check diff --base main

# Validate all ADRs
decider check adr

# Update index after ADR changes
decider index
```

## Getting Help

- Read relevant ADRs in `docs/adr/`
- Check guides in `docs/guides/`
- See specification in `SPEC.md`
- If constraints are unclear, ask for clarification
