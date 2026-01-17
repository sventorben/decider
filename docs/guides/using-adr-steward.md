# Using the ADR Steward

The ADR Steward is a Claude Code agent that manages ADR lifecycle operations. It handles creation, supersession, validation, and index maintenance while leaving architectural decisions to humans.

## What the Steward Does

| Task | Steward Handles | Human Decides |
|------|-----------------|---------------|
| **Create ADR** | File creation, formatting, index update | Decision content, constraints, scope |
| **Supersede ADR** | Linking old/new, status updates, index sync | Why the decision changed |
| **Validate** | Format checks, rationale pattern | Whether to accept warnings |
| **Index** | Regeneration, consistency checks | Nothing—fully automated |

The steward is a tool, not a decision-maker.

## Invoking the Steward

### Slash Commands

The most direct way to use the steward:

```
/adr-new Use Redis for session storage
```

```
/adr-supersede ADR-0003 Migrate from REST to GraphQL
```

These commands trigger the steward with specific tasks.

### Natural Language

You can also describe what you need:

> "Create an ADR for our decision to use PostgreSQL for the main database. It should apply to src/db/** and require all queries to use prepared statements."

> "We need to supersede ADR-0005 because we're switching from Memcached to Redis. The new decision should keep the same scope."

The steward interprets your intent and executes the appropriate commands.

### Direct CLI

For scripting or when you don't need agent assistance:

```bash
decider new "Your Decision Title" --tags tag1,tag2 --paths "src/**"
decider index
decider check adr --strict
```

## Creating an ADR

### Via Slash Command

```
/adr-new Adopt OpenTelemetry for distributed tracing
```

The steward will:
1. Run `decider new "Adopt OpenTelemetry for distributed tracing"`
2. Open the generated file
3. Help you fill in sections if prompted
4. Run `decider index` to update the index

### What You Provide

After the steward creates the file, you fill in:

- **Context**: Why this decision is needed
- **Decision**: What you chose and why (using the rationale pattern)
- **Constraints**: Specific rules to follow
- **Invariants**: Properties that must hold
- **Scope paths**: Which code the decision applies to

### Example Interaction

```
User: /adr-new Use repository pattern for data access

Steward: Created docs/adr/0008-use-repository-pattern-for-data-access.md

I've created the ADR with proposed status. You'll need to fill in:
- Context section explaining why this pattern is needed
- Decision section with "Adopted because:" and "Adopted despite:"
- Constraints (e.g., "All database queries must go through repositories")
- Scope paths (e.g., "src/db/**", "internal/storage/**")

Would you like help drafting any of these sections?
```

## Superseding an ADR

When a decision changes, don't edit the old ADR. Create a new one that supersedes it.

### Via Slash Command

```
/adr-supersede ADR-0003 Switch from MySQL to PostgreSQL
```

The steward will:
1. Create a new ADR with reference to ADR-0003
2. Update ADR-0003's status to `superseded`
3. Add cross-references in both ADRs' frontmatter
4. Regenerate the index

### What Gets Updated

**Old ADR (ADR-0003):**
```yaml
status: superseded
superseded_by:
  - ADR-0009
```

**New ADR (ADR-0009):**
```yaml
status: proposed
supersedes:
  - ADR-0003
```

## Validation

### Check ADR Format

```
User: Check if all our ADRs are valid

Steward: Running decider check adr --strict...

All 8 ADR(s) are valid.
```

Or if there are issues:

```
Steward: Found 2 warning(s) (strict mode):
  [warning] 0005-caching-strategy.md: rationale: missing 'Adopted despite:' section
  [warning] 0007-api-versioning.md: rationale: missing 'Rejected because:' for alternative

These ADRs are missing required rationale sections. Would you like me to help add them?
```

### Check Index Consistency

```
User: Is the ADR index up to date?

Steward: Running decider index --check...

Index is up to date.
```

## Troubleshooting

### "decider: command not found"

The CLI isn't installed or isn't in PATH.

**Fix:**
```bash
go install github.com/sventorben/decider/cmd/decider@v0.1.0
```

Ensure `$GOPATH/bin` is in your PATH.

### "ADR directory not found"

The `docs/adr/` directory doesn't exist.

**Fix:**
```bash
decider init
```

### "Index out of date"

The index doesn't match current ADR files.

**Fix:**
```bash
decider index
```

This regenerates `index.yaml` from the current ADR files.

### "Validation failed: missing rationale pattern"

An ADR doesn't follow the required rationale format.

**Fix:**
Add the missing sections to the ADR:

```markdown
### [Option]: Adopted

**Adopted because:**
- Concrete reason 1
- Concrete reason 2

**Adopted despite:**
- Known trade-off 1
- Known trade-off 2
```

For rejected alternatives:

```markdown
### [Alternative]: Rejected

**Rejected because:**
- Why it wasn't chosen

**Rejected despite:**
- Its legitimate strengths
```

### Steward Creates Wrong File Location

The steward might use a different directory than expected.

**Fix:**
Specify the directory explicitly:
```bash
decider new "Title" --dir docs/adr
```

Or set it in the slash command context:
```
/adr-new Title (in docs/adr/)
```

### Steward Doesn't Update Index

The `--no-index` flag might be set, or there was an error.

**Fix:**
Run index update manually:
```bash
decider index
```

Check for errors in the ADR files that might prevent indexing.

## Best Practices

1. **Let the steward handle mechanics.** Use it for file creation, formatting, and index updates. Don't fight the tooling.

2. **Review generated content.** The steward fills in templates, but you're responsible for the actual decision content.

3. **Run validation before committing.** A quick `decider check adr --strict` catches issues early.

4. **Keep the index in version control.** The `index.yaml` should be committed so other developers (and agents) can use it without regenerating.

5. **Use supersession, not editing.** When decisions change, create new ADRs that supersede old ones. This preserves history.

## Related Guides

- [Adopting DECIDER in a Team](adopting-decider-in-a-team.md) - Progressive rollout strategy
- [Writing ADRs](writing-adrs.md) - How to write effective constraints
- [Agent Integration](agent-integration.md) - Connecting other AI agents to DECIDER
