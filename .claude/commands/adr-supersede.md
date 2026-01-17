# /adr-supersede Command

Supersede an existing ADR with a new one.

## Usage

```
/adr-supersede <old-adr-id> <new-title>
```

## Description

This command creates a new ADR that supersedes an existing one, properly linking them together.

## Example
```
/adr-supersede ADR-0005 Migrate from REST to GraphQL
/adr-supersede ADR-0012 Use Redis instead of Memcached
```

## Workflow

When invoked, the agent should:

1. Read the old ADR to understand context:
   ```bash
   decider show <old-adr-id>
   ```

2. Create the new ADR:
   ```bash
   decider new "<new-title>" --tags <relevant-tags> --paths "<scope-paths>"
   ```

3. Update the new ADR's frontmatter:
   - Add old ADR to `supersedes` list
   - Copy relevant scope paths from old ADR
   - Reference old ADR in context

4. Update the old ADR's frontmatter:
   - Change `status` to `superseded`
   - Add new ADR ID to `superseded_by` list

5. Fill in the new ADR content:
   - Reference the old decision in Context
   - Explain what changed and why
   - Document new alternatives considered
   - Describe consequences of the change

6. Validate and update index:
   ```bash
   decider check adr
   decider index
   ```

## Example Update

Old ADR (ADR-0005):
```yaml
---
adr_id: ADR-0005
title: Use REST API
status: superseded          # Changed from adopted
superseded_by:
  - ADR-0015               # Added
---
```

New ADR (ADR-0015):
```yaml
---
adr_id: ADR-0015
title: Migrate from REST to GraphQL
status: adopted
supersedes:
  - ADR-0005               # Added
related_adrs:
  - ADR-0005
---
```

## Post-Supersession Checklist

- [ ] Old ADR status changed to `superseded`
- [ ] Old ADR `superseded_by` includes new ADR
- [ ] New ADR `supersedes` includes old ADR
- [ ] New ADR explains why the change was made
- [ ] Both ADRs pass validation
- [ ] Index has been regenerated
