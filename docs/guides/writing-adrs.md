# Writing Effective ADRs

This guide covers how to write ADRs that are useful for both humans and AI agents.

## When to Write an ADR

Write an ADR when a decision:

- **Affects architecture**: Changes system structure, component boundaries, or data flow
- **Has long-term implications**: Will influence code for months or years
- **Involves tradeoffs**: Choosing one approach means giving up another
- **Isn't obvious**: Future team members might ask "why did we do it this way?"

### Examples of ADR-worthy decisions

- Database choice (PostgreSQL vs MongoDB)
- API style (REST vs GraphQL)
- Authentication mechanism (JWT vs sessions)
- State management pattern (Redux vs Context)
- Error handling strategy
- Logging and observability approach

### Not ADR-worthy

- Which logging library to use (unless it's a significant investment)
- Code formatting preferences (use a linter config)
- Individual bug fixes

## ADR Structure

### Frontmatter

The YAML frontmatter makes ADRs machine-readable:

```yaml
---
adr_id: ADR-0001
title: "Use PostgreSQL for persistence"
status: adopted
date: 2026-01-16
scope:
  paths:
    - "src/db/**"
    - "migrations/**"
tags:
  - database
  - storage
constraints:
  - "All database access must go through the repository pattern"
  - "Use prepared statements for all queries"
invariants:
  - "Database connections are pooled"
  - "All migrations are reversible"
supersedes: []
superseded_by: []
related_adrs: []
---
```

### Constraints vs Invariants

**Constraints** are prescriptive rules:
- "All API endpoints must be versioned"
- "Never store passwords in plain text"
- "Use prepared statements for all queries"

**Invariants** are properties that must always hold:
- "Database connections are pooled"
- "All migrations are reversible"
- "User sessions expire after 24 hours"

The difference: constraints tell you what to DO; invariants describe what must BE TRUE.

### Scope Paths

Scope paths are glob patterns that indicate which code the ADR applies to:

```yaml
scope:
  paths:
    - "src/db/**"           # All files under src/db/
    - "**/*.proto"          # All .proto files anywhere
    - "cmd/server/**"       # Server command code
    - "migrations/*.sql"    # SQL migrations
```

**Tips for scope paths:**
- Be specific enough to avoid false positives
- Use `**` for recursive matching
- Include config files if they affect the decision
- Cover test files if constraints apply to tests

### Status Lifecycle

```
proposed → adopted → (active use)
    │          │
    │          ├── deprecated → (discouraged)
    │          │
    │          └── superseded → (replaced by new ADR)
    │
    └── rejected → (not adopted)
```

| Status | Meaning |
|--------|---------|
| `proposed` | Under discussion, not yet decided |
| `adopted` | Approved and in effect |
| `rejected` | Considered and explicitly rejected |
| `deprecated` | Was adopted, now discouraged but not replaced |
| `superseded` | Replaced by another ADR |

## Writing Good Sections

### Context

Explain WHY this decision is needed. Include:
- The problem you're solving
- Current situation or pain points
- Requirements driving the decision
- Constraints from the environment

**Good:**
> Our e-commerce application needs a database for storing user accounts, product catalog, orders, and shopping cart state. We require ACID compliance for financial transactions and need to support complex queries for reporting.

**Bad:**
> We need a database.

### Decision

State WHAT you decided and WHY. Include:
- The choice you made
- Key reasons for choosing it
- How it will be implemented

**Good:**
> We will use PostgreSQL 15+ as our primary database. All database access will follow the repository pattern, with SQL queries encapsulated in repository structs and business logic using interfaces.

**Bad:**
> We're using PostgreSQL.

### Alternatives Considered

Document what else you evaluated. Use a table:

| Alternative | Pros | Cons |
|-------------|------|------|
| MySQL | Widely adopted, good performance | Weaker JSON support, fewer advanced features |
| MongoDB | Flexible schema, easy horizontal scaling | No ACID for multi-document transactions |
| SQLite | Zero setup, embedded | Not suitable for concurrent web traffic |

This helps future readers understand why alternatives weren't chosen.

### Consequences

List both positive and negative impacts:

**Positive:**
- ACID compliance for order transactions
- Rich feature set (JSON columns, full-text search, CTEs)
- Mature ecosystem with excellent tooling

**Negative:**
- Requires operational knowledge for production deployment
- More complex than NoSQL for simple document storage

Be honest about tradeoffs. Every decision has costs.

### Agent Guidance (Optional)

Add specific instructions for AI agents:

```markdown
## Agent Guidance

When working in `src/db/`:
- Never write raw SQL in handlers or services; use repository methods
- Always use parameterized queries (`$1`, `$2`) not string concatenation
- Add indexes for columns used in WHERE clauses
- Write reversible migrations (both UP and DOWN)
```

## Common Mistakes

### 1. Constraints too vague

**Bad:** "Write good code"
**Good:** "All public functions must have doc comments"

### 2. Scope too broad

**Bad:** `paths: ["**"]` (matches everything)
**Good:** `paths: ["src/api/**", "internal/handlers/**"]`

### 3. Missing context

Don't assume readers know your situation. Explain the problem before the solution.

### 4. No alternatives

Even if the choice seemed obvious, document what else was considered. Future readers need to know the decision was thoughtful.

### 5. Forgetting consequences

Every decision has downsides. Documenting them shows mature thinking and helps future decisions.

## Template

```markdown
---
adr_id: ADR-NNNN
title: "Decision Title"
status: proposed
date: YYYY-MM-DD
scope:
  paths:
    - "affected/path/**"
tags:
  - relevant-tag
constraints:
  - "Rule that must be followed"
invariants:
  - "Property that must hold"
supersedes: []
superseded_by: []
related_adrs: []
---

# ADR-NNNN: Decision Title

## Context

[What is the problem? What requirements or constraints exist?]

## Decision

[What did we decide? Why this choice?]

## Alternatives Considered

| Alternative | Pros | Cons |
|-------------|------|------|
| Option A | ... | ... |
| Option B | ... | ... |

## Consequences

**Positive:**
- Benefit 1
- Benefit 2

**Negative:**
- Cost 1
- Cost 2

## Agent Guidance

[Optional: Specific instructions for AI agents]
```

## Checklist

Before finalizing an ADR:

- [ ] Title clearly describes the decision
- [ ] Context explains why the decision is needed
- [ ] Decision states what was chosen and why
- [ ] Alternatives are documented with pros/cons
- [ ] Consequences include both positive and negative impacts
- [ ] Constraints are specific and actionable
- [ ] Scope paths accurately match affected code
- [ ] Status reflects current state (usually `proposed` initially)
- [ ] Related ADRs are linked
- [ ] Runs `decider check adr` without errors
