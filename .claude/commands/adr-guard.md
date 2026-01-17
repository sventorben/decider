# /adr-guard

Check that current changes comply with applicable ADRs.

## Usage

```
/adr-guard [--base <ref>]
```

**Arguments:**
- `--base <ref>`: Git reference to compare against (default: `origin/main`)

## What This Command Does

1. Runs `decider check diff --base <ref>` to find applicable ADRs
2. For each applicable ADR:
   - Reads constraints and invariants
   - Checks current changes for violations
3. If violations exist:
   - Proposes and applies fixes
   - Re-runs checks until compliant
4. If compliance requires changing the ADR:
   - Stops and proposes an ADR update
   - Does not force non-compliant code
5. Reports final status: compliant, blocked, or no ADRs apply

## Examples

```
/adr-guard
```
Check against origin/main (default).

```
/adr-guard --base main
```
Check against local main branch.

```
/adr-guard --base HEAD~5
```
Check against 5 commits ago.

## When to Use

- After implementing changes, before final output
- When CI reports ADR constraint violations
- Before creating a pull request
- After modifying code in ADR-scoped paths

## Exit Conditions

The command finishes when:
- **Compliant**: All ADR constraints satisfied
- **Blocked**: Cannot proceed without ADR change (proposes update)
- **No ADRs Apply**: No ADRs match the changed files