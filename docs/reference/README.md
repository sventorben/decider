# Reference Documentation

This section contains reference documentation for DECIDER's CLI commands and file formats.

## CLI Commands

### decider init

Initialize ADR directory structure.

```bash
decider init [--dir PATH]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--dir` | ADR directory path | `docs/adr` |

Creates:
- Directory structure
- Template file (`templates/adr.md`)
- Empty index (`index.yaml`)

---

### decider new

Create a new ADR.

```bash
decider new [OPTIONS] "TITLE"
```

| Flag | Description | Default |
|------|-------------|---------|
| `--dir` | ADR directory | `docs/adr` |
| `--tags` | Comma-separated tags | none |
| `--paths` | Comma-separated scope paths | none |
| `--status` | Initial status | `proposed` |
| `--no-index` | Skip index update | false |
| `--format` | Output format (`text`, `json`) | `text` |

Example:
```bash
decider new "Use PostgreSQL for persistence" \
  --tags database,storage \
  --paths "src/db/**,migrations/**"
```

---

### decider list

List ADRs with optional filters.

```bash
decider list [OPTIONS]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--dir` | ADR directory | `docs/adr` |
| `--status` | Filter by status | none |
| `--tag` | Filter by tag (repeatable) | none |
| `--path` | Filter by scope path match | none |
| `--format` | Output format (`text`, `json`) | `text` |

Examples:
```bash
decider list
decider list --status adopted
decider list --tag security --tag api
decider list --path "src/api/handler.go"
```

---

### decider show

Display details of a specific ADR.

```bash
decider show [OPTIONS] IDENTIFIER
```

| Argument | Description |
|----------|-------------|
| `IDENTIFIER` | ADR-NNNN, NNNN, or filename |

| Flag | Description | Default |
|------|-------------|---------|
| `--dir` | ADR directory | `docs/adr` |
| `--format` | Output format (`text`, `json`) | `text` |

Examples:
```bash
decider show ADR-0001
decider show 0001
decider show 0001-use-postgresql.md
```

---

### decider check adr

Validate all ADRs for format compliance.

```bash
decider check adr [OPTIONS]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--dir` | ADR directory | `docs/adr` |
| `--format` | Output format (`text`, `json`) | `text` |

Validates:
- Required frontmatter fields
- Valid status values
- Date format (YYYY-MM-DD)
- ADR ID format (ADR-NNNN)
- Filename/ID consistency
- Required markdown sections

Exit codes:
- `0`: All valid
- `1`: Parse/usage error
- `2`: Validation failures

---

### decider check diff

Find ADRs applicable to changed files.

```bash
decider check diff --base REF [OPTIONS]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--base` | Git ref for comparison (required) | none |
| `--dir` | ADR directory | `docs/adr` |
| `--format` | Output format (`text`, `json`) | `text` |

Examples:
```bash
decider check diff --base main
decider check diff --base HEAD~5
decider check diff --base origin/develop --format json
```

---

### decider explain

Explain why ADRs apply to changes.

```bash
decider explain --base REF [OPTIONS]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--base` | Git ref for comparison (required) | none |
| `--dir` | ADR directory | `docs/adr` |
| `--format` | Output format (`text`, `json`) | `text` |

Like `check diff` but with narrative explanation showing which patterns matched which files.

---

### decider index

Generate or verify the ADR index.

```bash
decider index [OPTIONS]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--dir` | ADR directory | `docs/adr` |
| `--check` | Verify without modifying | false |
| `--format` | Output format (`text`, `json`, `yaml`) | `text` |

Exit codes:
- `0`: Success (or up-to-date with `--check`)
- `1`: Error
- `2`: Index out of date (with `--check`)

---

### decider version

Show version information.

```bash
decider version
```

Output:
```
decider version 0.1.0
  commit: abc123
  built:  2026-01-16T10:00:00Z
```

---

## File Formats

### ADR Frontmatter Schema

```yaml
---
adr_id: ADR-NNNN        # Required. Format: ADR-NNNN
title: "Title"           # Required. Human-readable title
status: adopted         # Required. See status values below
date: YYYY-MM-DD         # Required. ISO 8601 date

scope:                   # Optional
  paths:                 # Glob patterns for affected code
    - "path/**"

tags:                    # Optional. Categorization
  - tag1

constraints:             # Optional. Rules that MUST be followed
  - "Constraint"

invariants:              # Optional. Properties that must hold
  - "Invariant"

supersedes: []           # Optional. ADR IDs this supersedes
superseded_by: []        # Optional. ADR IDs that supersede this
related_adrs: []         # Optional. Related ADR IDs
---
```

### Status Values

| Status | Description |
|--------|-------------|
| `proposed` | Under discussion |
| `adopted` | Approved and in effect |
| `rejected` | Considered but not adopted |
| `deprecated` | Was adopted, now discouraged |
| `superseded` | Replaced by another ADR |

### Index Schema

```yaml
# AUTO-GENERATED by decider index - DO NOT EDIT
generated_at: "2026-01-16T10:00:00Z"   # RFC3339 UTC
adr_count: 4                            # Number of ADRs

adrs:
  - adr_id: ADR-0001                    # ADR identifier
    title: "Decision Title"             # ADR title
    status: adopted                    # Current status
    date: "2026-01-16"                  # Decision date
    tags:                               # Tags (may be empty)
      - foundation
    scope_paths:                        # Scope paths (may be empty)
      - "src/**"
    file: "0001-decision-title.md"      # Filename
```

### Glob Patterns

| Pattern | Matches |
|---------|---------|
| `*.go` | Go files in current directory |
| `src/*` | Direct children of src/ |
| `src/**` | All files under src/ recursively |
| `src/**/*.go` | All Go files under src/ |
| `**/*.proto` | All .proto files anywhere |

---

## Complete Specification

See [SPEC.md](../../SPEC.md) for the complete, stable specification including:

- Filename conventions
- All validation rules
- Exit code semantics
- Stability guarantees
