# Architecture Decision Records

This directory contains the ADRs for the demo e-commerce application.

## What Are ADRs?

Architecture Decision Records (ADRs) document significant technical decisions with:
- **Context**: Why the decision was needed
- **Decision**: What was chosen and why (using the mandatory rationale pattern)
- **Alternatives**: What else was considered and why it was rejected
- **Consequences**: Positive and negative impacts
- **Constraints**: Rules that MUST be followed
- **Invariants**: Properties that must always hold true

## Current ADRs

| ID | Title | Status | Key Constraints |
|----|-------|--------|-----------------|
| ADR-0001 | Use PostgreSQL for Persistence | adopted | Repository pattern, prepared statements |
| ADR-0002 | REST API with Versioning | adopted | Version prefix, OpenAPI spec required |
| ADR-0003 | Redux for Frontend State | deprecated | (No longer enforced) |

## Directory Structure

```
docs/adr/
├── README.md           # This file
├── index.yaml          # Auto-generated index (do not edit manually)
├── templates/
│   └── adr.md          # Template for new ADRs
├── 0001-*.md           # Individual ADRs
├── 0002-*.md
└── 0003-*.md
```

## Mandatory Rationale Pattern

All ADRs MUST use explicit rationale sections:

### For Adopted Options
```markdown
### [Option Name]: Adopted

**Adopted because:**
- Clear reason why this was chosen
- Tied to decision drivers

**Adopted despite:**
- Known downside we accepted
- Trade-off compared to alternatives
```

### For Rejected Alternatives
```markdown
### [Alternative Name]: Rejected

**Rejected because:**
- Clear reason for rejection

**Rejected despite:**
- Legitimate strength of this option
```

## Working with ADRs

### Validate All ADRs

```bash
decider check adr --dir docs/adr --strict
```

### List ADRs

```bash
decider list --dir docs/adr
```

### Find ADRs for a File

```bash
decider list --dir docs/adr --path "src/db/users.go"
```

### View ADR Details

```bash
decider show ADR-0001 --dir docs/adr
```

### Update Index

```bash
decider index --dir docs/adr
```

## Adding New ADRs

1. Create the ADR:
   ```bash
   decider new "Your Decision" --dir docs/adr --tags tag1 --paths "affected/**"
   ```

2. Edit the generated file to fill in all sections using the rationale pattern

3. Validate:
   ```bash
   decider check adr --dir docs/adr --strict
   ```

4. Update the index:
   ```bash
   decider index --dir docs/adr
   ```
