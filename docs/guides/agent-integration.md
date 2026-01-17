# Agent Integration Guide

This guide explains how to use the ADR Steward agent and integrate DECIDER with AI coding agents.

## Overview

The ADR Steward is a Claude Code agent that manages ADR lifecycle operations. It handles:

- Creating new ADRs with proper formatting
- Superseding old decisions with new ones
- Validating ADR compliance
- Maintaining the index file

## Using the ADR Steward

### Slash Commands

The steward provides two slash commands:

#### `/adr-new` - Create a New ADR

```
/adr-new Use Redis for session storage
```

The steward will:
1. Run `decider new "Use Redis for session storage"`
2. Prompt you for tags and scope paths
3. Fill in the Context, Decision, Alternatives, and Consequences sections
4. Add appropriate constraints and invariants
5. Run `decider check adr` to validate
6. Run `decider index` to update the index

#### `/adr-supersede` - Supersede an Existing ADR

```
/adr-supersede ADR-0003 Migrate from Redux to Zustand
```

The steward will:
1. Read the old ADR (ADR-0003) to understand context
2. Create a new ADR documenting the change
3. Update the old ADR's status to `superseded`
4. Link both ADRs with `supersedes` and `superseded_by`
5. Validate both ADRs
6. Update the index

### Direct Invocation

You can also invoke the steward with natural language:

> "Create an ADR for our decision to use GraphQL for the new mobile API"

> "We decided to switch from Memcached to Redis. Please create an ADR that supersedes ADR-0005"

> "Check if all our ADRs are valid and the index is current"

> "What ADRs apply to the files I'm changing?"

### Example Session

```
You: /adr-new Use PostgreSQL for analytics data warehouse

Steward: I'll create an ADR for using PostgreSQL as the analytics data warehouse.

What tags should I use? (e.g., database, analytics, infrastructure)

You: database, analytics, data-warehouse

Steward: What code paths does this decision affect?

You: src/analytics/**, internal/warehouse/**, migrations/analytics/**

Steward: I've created ADR-0005: Use PostgreSQL for analytics data warehouse.

[Shows the generated ADR content]

I've filled in the sections based on the decision. Please review:
- Context explains the need for an analytics-specific database
- Decision documents PostgreSQL choice with reasoning
- Alternatives include BigQuery, Snowflake, ClickHouse
- Consequences list both benefits and operational costs

Constraints I've added:
- "Analytics queries must not impact production database"
- "Use read replicas for heavy analytical workloads"

Do you want me to adjust anything before committing?
```

## How Agents Query ADRs

AI agents can query DECIDER to understand constraints before writing code.

### Before Making Changes

```bash
# Find ADRs that apply to current changes
decider check diff --base main

# Get detailed explanation
decider explain --base main

# Get machine-readable output
decider check diff --base main --format json
```

### Programmatic Access

```bash
# Get constraints as JSON
$ decider check diff --base main --format json | jq '.applicable_adrs[].constraints'
[
  "All database access must go through the repository pattern",
  "Use prepared statements for all queries"
]

# Get scope paths
$ decider show ADR-0001 --format json | jq '.scope.paths'
[
  "src/db/**",
  "migrations/**"
]
```

### Agent Workflow Pattern

An effective agent workflow:

1. **Before coding**: Query applicable ADRs
   ```bash
   decider check diff --base main --format json
   ```

2. **During coding**: Follow constraints in applicable ADRs

3. **After coding**: Verify changes don't violate invariants

4. **If architectural change needed**: Create/update ADR first
   ```
   /adr-new <decision>
   ```

## Configuring Agents to Use ADRs

### AGENTS.md

The repository's `AGENTS.md` file tells agents how to work with ADRs:

```markdown
## Before Making Changes

1. Run `decider check diff --base main` to find applicable ADRs
2. Read the constraints and invariants in each applicable ADR
3. Ensure your changes comply with documented constraints

## When to Create ADRs

Create an ADR when:
- Adding a new external dependency
- Changing data storage approach
- Modifying API contracts
- Introducing new patterns
```

### Claude Code Integration

For Claude Code specifically, the `.claude/` directory contains:

```
.claude/
├── agents/
│   └── adr-steward.md       # Steward agent definition
├── commands/
│   ├── adr-new.md           # /adr-new command
│   └── adr-supersede.md     # /adr-supersede command
└── skills/
    └── adr-system/
        └── SKILL.md         # ADR system skill
```

## Troubleshooting

### Agent didn't follow ADRs

**Symptoms:** Agent wrote code that violates documented constraints.

**Causes:**
1. Agent didn't query ADRs before coding
2. Scope paths don't match the files being changed
3. Constraints are too vague

**Solutions:**
- Ensure AGENTS.md instructs agents to query ADRs
- Review scope paths in relevant ADRs
- Make constraints more specific and actionable

### Index out of date

**Symptoms:** `decider index --check` fails.

**Cause:** ADRs were modified without regenerating the index.

**Solution:**
```bash
decider index
# Commit the updated index.yaml
```

### ADR scope too broad

**Symptoms:** ADR applies to unrelated files.

**Cause:** Scope paths like `**` or `src/**` are too general.

**Solution:** Use more specific paths:
```yaml
# Too broad
scope:
  paths:
    - "src/**"

# Better
scope:
  paths:
    - "src/db/**"
    - "src/repositories/**"
```

### Claude is stalling / narrating tasks

**Symptoms:** Agent describes what it will do instead of doing it.

**Cause:** Agent is in planning mode or unsure how to proceed.

**Solutions:**
- Use explicit commands: `/adr-new` instead of asking
- Provide complete information upfront (title, tags, paths)
- Check if there are validation errors blocking progress

### Supersession not working correctly

**Symptoms:** Old and new ADRs aren't properly linked.

**Checklist:**
1. Old ADR status is `superseded`
2. Old ADR has new ADR in `superseded_by`
3. New ADR has old ADR in `supersedes`
4. Both ADRs pass `decider check adr`
5. Index has been regenerated

## Best Practices

### For Teams

1. **Document early**: Create ADRs during planning, not after implementation
2. **Keep scope tight**: Specific paths prevent false positives
3. **Review constraints**: Ensure they're actionable, not vague
4. **Validate in CI**: Catch drift before it reaches main

### For Agents

1. **Query first**: Always check applicable ADRs before coding
2. **Follow constraints**: Treat constraints as hard requirements
3. **Ask when unsure**: If a constraint is unclear, ask for clarification
4. **Create ADRs proactively**: Suggest ADRs for significant decisions

### For Steward Usage

1. **Provide context**: Give the steward enough information to write good ADRs
2. **Review output**: Always review generated content before committing
3. **Use supersession**: Don't edit old ADRs; supersede them
4. **Keep index current**: Let the steward handle index updates
