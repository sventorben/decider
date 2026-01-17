---
adr_id: ADR-0002
title: "ADR Format: Markdown + YAML Frontmatter + Required Sections"
status: adopted
date: 2026-01-16
scope:
  paths:
    - "docs/adr/**/*.md"
    - "internal/adr/**"
tags:
  - format
  - schema
  - documentation
constraints:
  - Frontmatter must be valid YAML between `---` delimiters
  - All required frontmatter keys must be present
  - Filename must match pattern `NNNN-kebab-title.md`
  - ADR ID in frontmatter must match filename number
invariants:
  - ADRs are always valid Markdown
  - Frontmatter is always machine-parseable YAML
  - Required sections are present in body
supersedes: []
superseded_by: []
related_adrs:
  - ADR-0003
---

# ADR-0002: ADR Format: Markdown + YAML Frontmatter + Required Sections

## Context

ADRs need to be:
1. **Human-readable**: Engineers read and write them in editors/GitHub
2. **Machine-readable**: AI agents and tools need structured metadata
3. **Versionable**: Git-friendly plain text format
4. **Discoverable**: Searchable by status, tags, scope paths

Decision drivers:
- AI agents need structured metadata to query constraints before writing code
- Human engineers must be able to read/write ADRs without special tooling
- Format must be widely supported by editors, linters, and CI tools
- Single file per decision reduces sync issues

Traditional ADR formats (Nygard-style) lack structured metadata for automation.

## Decision

We adopt a **Markdown + YAML frontmatter** format with required fields and sections.

### Markdown + YAML Frontmatter: Adopted

**Adopted because:**
- Single file per ADR eliminates synchronization issues between metadata and content
- YAML frontmatter is a widely-recognized convention (Jekyll, Hugo, Docusaurus)
- Markdown body renders correctly on GitHub, GitLab, and most documentation tools
- Structured frontmatter enables machine parsing for constraints, tags, and scope
- IDE support is excellent (syntax highlighting, linting, previews)
- Validation is straightforward with standard YAML parsers

**Adopted despite:**
- Frontmatter must be kept in sync with body content manually
- YAML syntax errors can break parsing of entire file
- Learning curve for scope.paths glob patterns
- Some fields (constraints, invariants) may drift from implementation reality

### Frontmatter Schema (Required Keys)

```yaml
---
adr_id: ADR-NNNN           # String, must match filename number
title: "Short descriptive title"
status: proposed|adopted|rejected|deprecated|superseded
date: YYYY-MM-DD           # ISO 8601 date
scope:
  paths:                   # List of glob patterns this ADR applies to
    - "src/api/**"
    - "*.config.js"
tags:                      # List of categorization tags
  - security
  - api
constraints:               # List of rules that MUST be followed
  - "Constraint statement"
invariants:                # List of properties that must always hold
  - "Invariant statement"
supersedes: []             # List of ADR IDs this supersedes
superseded_by: []          # List of ADR IDs that supersede this
related_adrs: []           # List of related ADR IDs
---
```

### Status Enum

| Status | Meaning |
|--------|---------|
| `proposed` | Under discussion, not yet decided |
| `adopted` | Approved and in effect |
| `rejected` | Considered and explicitly rejected |
| `deprecated` | Was adopted, now discouraged |
| `superseded` | Replaced by another ADR |

### Required Markdown Sections

1. **Context** - Why this decision is needed
2. **Decision** - What we decided and why
3. **Alternatives Considered** - Other options evaluated
4. **Consequences** - Positive and negative impacts

### Optional Section

- **Agent Guidance** - Specific instructions for AI agents

### Filename Convention

```
docs/adr/NNNN-kebab-case-title.md
```

- `NNNN`: Zero-padded 4-digit number (0001, 0002, ...)
- `kebab-case-title`: Lowercase, hyphens, derived from title
- Example: `0042-use-postgres-for-persistence.md`

## Alternatives Considered

### Pure Markdown (No Frontmatter): Rejected

**Rejected because:**
- No structured metadata for machine parsing
- Tags, status, and scope would need to be parsed from prose
- AI agents cannot reliably extract constraints without structured fields
- Filtering and searching requires full-text search instead of indexed lookups

**Rejected despite:**
- Simplest possible format with zero learning curve
- Maximum compatibility with any Markdown renderer
- No risk of YAML syntax errors breaking the file

### JSON Files: Rejected

**Rejected because:**
- JSON is difficult to read and write for humans
- No comment support makes documentation harder
- Merge conflicts are more likely with strict JSON syntax
- Engineers would need IDE plugins for comfortable editing

**Rejected despite:**
- Machine parsing is trivial and unambiguous
- Schema validation is well-supported (JSON Schema)
- No frontmatter delimiter confusion
- Widely supported in all programming languages

### TOML Frontmatter: Rejected

**Rejected because:**
- Less common than YAML frontmatter in the Markdown ecosystem
- Tooling support is weaker (fewer linters, fewer IDE extensions)
- Engineers are generally less familiar with TOML syntax
- Hugo and Jekyll ecosystem uses YAML, reducing knowledge transfer

**Rejected despite:**
- Cleaner syntax than YAML for simple key-value pairs
- Less ambiguity in parsing (no "Norway problem")
- Better error messages for syntax issues

### Separate YAML + Markdown Files: Rejected

**Rejected because:**
- Two files per ADR doubles the file count and management overhead
- Synchronization between metadata file and content file can drift
- Git operations (rename, delete) must handle two files atomically
- Cognitive overhead of jumping between files when reading/writing

**Rejected despite:**
- Clean separation of concerns between metadata and prose
- YAML file can be validated independently
- Markdown file can be rendered without frontmatter stripping

## Consequences

**Positive:**
- Single file per ADR is easy to manage in Git
- Standard format recognized by GitHub, GitLab, documentation tools
- IDE support for Markdown + YAML is mature
- Easy to validate programmatically with existing libraries

**Negative:**
- Frontmatter must be kept in sync with content (mitigated by validation)
- Learning curve for scope.paths glob patterns (mitigated by examples)
- YAML syntax errors can be subtle (mitigated by `decider check adr`)

## Agent Guidance

When creating or modifying ADRs:
- Always validate frontmatter YAML syntax
- Ensure `adr_id` matches the filename number
- Use lowercase kebab-case for filenames
- Keep constraints actionable and specific
- Scope paths should use glob patterns matching project structure
