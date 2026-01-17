# Adopting DECIDER in a Team

DECIDER works best when adopted incrementally. This guide walks through a progressive rollout that builds confidence without disrupting existing workflows.

## The Three Phases

| Phase | What You Do | Enforcement | Effort |
|-------|-------------|-------------|--------|
| 1. Documentation | Write ADRs, no tooling | None | Low |
| 2. Steward | Use CLI and agent | Advisory | Medium |
| 3. CI Checks | Automated validation | Required | Low |

Try to spend 1-2 month in each phase before progressing.

## Phase 1: Documentation Only

Start by writing ADRs without any tooling enforcement.

### Goals
- Establish the habit of documenting decisions
- Build initial ADR inventory
- Learn what makes good constraints

### Steps

1. **Initialize the ADR directory**
   ```bash
   decider init
   ```

2. **Document existing decisions**

   Start with 3-5 decisions your team already follows but hasn't written down:
   - Primary database choice
   - API authentication method
   - Core framework or language version
   - Error handling conventions

3. **Create ADRs for new decisions**

   When the team makes a decision in a meeting or PR discussion, capture it:
   ```bash
   decider new "Use OpenTelemetry for observability" --tags observability
   ```

4. **Keep it lightweight**

   At this phase, focus on constraints and scope. Don't worry about perfect index consistency or CI integration yet.

### Success Criteria
- Team has 5+ ADRs covering major decisions
- New architectural discussions reference ADRs
- No enforcement—just documentation

## Phase 2: Steward and CLI

Add tooling to maintain consistency and surface relevant decisions.

### Goals
- Automated index maintenance
- Developers query applicable ADRs before coding
- ADR Steward handles lifecycle tasks

### Steps

1. **Use the CLI regularly**

   Before making changes, check which ADRs apply:
   ```bash
   decider check diff --base main
   ```

2. **Adopt the ADR Steward**

   Use Claude Code's `/adr-new` and `/adr-supersede` commands instead of manual creation. The steward ensures:
   - Consistent formatting
   - Automatic index updates
   - Proper rationale pattern

3. **Add scope paths to existing ADRs**

   Review ADRs and add specific scope paths:
   ```yaml
   scope:
     paths:
       - "src/api/**"
       - "internal/handlers/**"
   ```

4. **Validate locally**
   ```bash
   decider check adr --strict
   decider index --check
   ```

### Success Criteria
- `decider check diff` is part of pre-coding routine
- ADR Steward handles most ADR operations
- Index stays in sync without manual intervention

## Phase 3: CI Integration

Add automated checks to enforce ADR health.

### Goals
- PRs can't merge with invalid ADRs
- Index drift is caught automatically
- Applicable constraints are visible in PR reviews

### Steps

1. **Add CI workflow**

   Create `.github/workflows/adr.yml`:
   ```yaml
   name: ADR Validation
   on: [pull_request]
   jobs:
     validate:
       runs-on: ubuntu-24.04
       steps:
         - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2
         - uses: actions/setup-go@d35c59abb061a4a6fb18e82ac0862c26744d6ab5 # v5.5.0
           with:
             go-version-file: go.mod
         - run: go install github.com/sventorben/decider/cmd/decider@v0.1.0
         - run: decider check adr --strict
         - run: decider index --check
   ```

2. **Optional: Comment applicable ADRs on PRs**

   Add a step that runs `decider check diff --base ${{ github.base_ref }}` and posts results as a PR comment.

3. **Monitor and adjust**

   If validation fails frequently, the issue is usually:
   - Scope paths too broad (everything matches)
   - Constraints too vague (hard to follow)
   - Missing ADRs for common decisions

### Success Criteria
- CI blocks PRs with ADR validation failures
- Developers fix ADR issues before requesting review
- False positives are rare

## Common Concerns

### "This feels like bureaucracy"

DECIDER is opt-in at every level. If writing an ADR feels bureaucratic, ask:
- Is this decision actually significant?
- Could the constraint be more specific?
- Is the scope too broad?

Not every decision needs an ADR. Focus on decisions that:
- Affect multiple parts of the codebase
- Involve trade-offs worth documenting
- Would confuse a new team member

### "We already have a wiki"

Wikis are fine for prose documentation. DECIDER handles a specific niche: decisions that need to be surfaced to developers (and agents) at the right moment.

Keep your wiki. Use DECIDER for constraints that should appear when someone touches affected code.

### "What about legacy decisions?"

Don't backfill everything. Instead:
1. Document decisions as they become relevant
2. When a question arises about "why do we do X?", capture the answer as an ADR
3. Let the ADR inventory grow organically

Backfilling creates stale documentation. Capturing decisions in context keeps them fresh.

### "Our team is too small"

Teams of 2-3 benefit from externalized decisions too. The value increases when:
- You onboard new team members
- You use AI coding assistants
- You revisit code after months away

Even if you're the only developer, future-you will appreciate past-you's documentation.

## Anti-Patterns to Avoid

**Scope: `**/*`**
If an ADR applies everywhere, it probably applies nowhere useful. Narrow the scope.

**Constraint: "Follow best practices"**
Vague constraints can't be followed or verified. Be specific about what "best practices" means.

**50 ADRs in week one**
Backfilling creates busy work without value. Document decisions as you make them.

**Ignoring CI failures**
If the team routinely skips ADR validation, either the checks are wrong or the process has broken down. Fix the root cause.

## Next Steps

- [Using ADR Steward](using-adr-steward.md) - Detailed guide to the Claude Code agent
- [CI Integration](ci-integration.md) - Complete CI setup instructions
- [Writing ADRs](writing-adrs.md) - How to write effective constraints
