# /adr-new Command

Create a new Architecture Decision Record.

## Usage

```
/adr-new <title>
```

## Description

This command creates a new ADR using the DECIDER CLI. It will:

1. Determine the next ADR number automatically
2. Generate a kebab-case filename
3. Create the ADR with standard template
4. Update the index

## Examples

```
/adr-new Use PostgreSQL for persistence
/adr-new Adopt GraphQL for API layer
/adr-new Implement feature flags system
```

## Workflow

When invoked, the agent should:

1. Run `decider new "<title>"` with appropriate flags:
   - `--tags` for relevant categories
   - `--paths` for affected code paths
   - `--status proposed` (default)

2. Edit the generated ADR file to fill in:
   - Context section
   - Decision section
   - Alternatives Considered
   - Consequences
   - Constraints and invariants

3. Run `decider check adr` to validate

4. Run `decider index` to update the index

## Flags to Consider

- `--tags`: Comma-separated list of tags (e.g., `api,security`)
- `--paths`: Comma-separated glob patterns for scope
- `--status`: Initial status (default: proposed)

## Post-Creation Checklist

- [ ] Context explains why this decision is needed
- [ ] Decision clearly states what we chose and why
- [ ] Alternatives lists other options considered
- [ ] Consequences covers both positive and negative impacts
- [ ] Constraints are specific and actionable
- [ ] Invariants describe properties to preserve
- [ ] Scope paths accurately reflect affected code
- [ ] Tags are relevant and follow conventions
