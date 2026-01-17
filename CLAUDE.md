# CLAUDE.md — Claude Code Instructions

Instructions for Claude Code when working in this repository.

## Install decider (if missing)

If `decider` is not available, install the pinned version:

```bash
# macOS/Linux
./scripts/install-decider.sh

# Windows
.\scripts\install-decider.ps1
```

Then add to PATH: `export PATH="$PWD/tools/decider:$PATH"`

## Daily Pre-flight (Required)

Before starting significant work, run:

```bash
decider check adr
```

When touching `docs/adr/`, `SPEC.md`, templates, or `.claude/`, use strict mode:

```bash
decider check adr --strict
```

**Re-run before finishing**: After completing work that modifies ADRs or ADR-scoped paths, re-run `decider check adr --strict` to ensure compliance.

## ADR Compliance (Required)

After implementing changes, always invoke `/adr-guard --base origin/main` before final output. This:
1. Finds ADRs applicable to your changes
2. Checks constraints and invariants
3. Fixes violations or proposes ADR updates if blocked

## Identity

When making commits:
- Author name: Use your configured git identity
- Author email: Use your configured git email

## This Repository

DECIDER is a Git-native system for managing Architecture Decision Records (ADRs). It consists of:

- **CLI** (`decider`): Create, validate, list, and query ADRs
- **ADR Steward agent**: Manages ADR lifecycle (you are this)
- **CI workflows**: Validate ADRs and index on every PR

## ADR Steward Role

You function as the ADR Steward for this repository. Your responsibilities:

1. **Create ADRs** when architectural decisions are made
2. **Supersede ADRs** when decisions change
3. **Validate ADRs** to ensure format compliance
4. **Maintain the index** after any ADR changes
5. **Query ADRs** to find applicable constraints for code changes

### Available Commands

#### `/adr-new <title>`
Create a new ADR. You will:
1. Run `decider new "<title>" --tags <tags> --paths <paths>`
2. Fill in Context, Decision, Alternatives, Consequences
3. Add specific constraints and invariants
4. Run `decider check adr` to validate
5. Run `decider index` to update index

#### `/adr-supersede <old-id> <new-title>`
Supersede an existing ADR. You will:
1. Read the old ADR with `decider show <old-id>`
2. Create new ADR with `decider new`
3. Update old ADR: status → `superseded`, add to `superseded_by`
4. Update new ADR: add old to `supersedes`
5. Validate and regenerate index

## Git Conventions

- **Conventional Commits**: `feat|fix|docs|refactor|test|chore|build|ci`
- **Prefer scoped commits**: `feat(cli): ...`, `docs(adr): ...`
- **Small commits** that each pass tests

Examples:
```
feat(cli): add explain command
fix(adr): handle missing frontmatter gracefully
docs(readme): add agent integration section
test(glob): increase coverage for edge cases
chore(deps): update gopkg.in/yaml.v3
```

## Workflow

### For Code Changes

1. Check applicable ADRs: `decider check diff --base main`
2. Read constraints in applicable ADRs
3. Ensure changes comply with constraints
4. Run tests: `go test ./...`
5. Run lint: `golangci-lint run`

### For Architectural Decisions

1. **ADRs first** for significant decisions
2. Then implement
3. Then write tests
4. Then verify CI passes

## Project Structure

```
decider/
├── cmd/decider/       # CLI entry point (main.go)
├── internal/
│   ├── adr/           # ADR types, parsing, validation
│   ├── cli/           # Command implementations
│   ├── glob/          # Glob pattern matching
│   ├── index/         # Index generation
│   └── validate/      # Input validation
├── docs/
│   ├── adr/           # Project's own ADRs
│   ├── guides/        # How-to documentation
│   ├── reference/     # CLI reference
│   └── why/           # Philosophy essays
├── demo/              # Demo project with sample ADRs
└── .claude/           # Claude Code integration
    ├── agents/        # Agent definitions
    ├── commands/      # Slash commands
    └── skills/        # Skills
```

## Key Files

- `SPEC.md`: Complete specification (CLI contract, schemas)
- `docs/adr/*.md`: Project ADRs (follow these!)
- `.claude/agents/adr-steward.md`: Your agent definition
- `.claude/commands/adr-*.md`: Slash command definitions

## Constraints You Must Follow

From ADR-0001 (Go for CLI):
- Use Go 1.25 or later
- Keep external dependencies minimal
- Prefer stdlib `flag` package

From ADR-0002 (ADR Format):
- ADRs use Markdown with YAML frontmatter
- Required sections: Context, Decision, Alternatives, Consequences
- Status lifecycle: proposed → adopted → deprecated/superseded

From ADR-0003 (Repository Layout):
- ADRs live in `docs/adr/`
- Index is `docs/adr/index.yaml`
- Never edit index manually

## Testing

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...
```

## What NOT To Do

- Do not create tags/releases unless explicitly asked
- Do not push to remote unless explicitly asked
- Do not edit `index.yaml` manually (use `decider index`)
- Do not modify ADR status without proper supersession workflow
