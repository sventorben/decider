# AI Agents Need Constraints

How structured decisions help AI coding agents make better choices.

## The Agent Problem

AI coding agents are powerful. Give them a task and they'll generate code—often correct, sometimes elegant, occasionally brilliant.

But they have a fundamental limitation: **they don't know your architecture**.

When you ask an agent to "add a function to get user by email," it doesn't know:
- You've standardized on the repository pattern
- Raw SQL outside the db package violates your security policy
- You use prepared statements, never string concatenation
- The UserRepository interface already has a similar method

Without this context, the agent makes reasonable guesses. Those guesses might violate constraints your team spent months establishing.

## Why Context Matters

Consider this prompt:

> "Add a function to fetch user by email in the API handler"

An agent might generate:

```go
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
    query := "SELECT * FROM users WHERE email = '" + email + "'"
    row := db.QueryRow(query)
    // ...
}
```

This code works. It's also:
- A SQL injection vulnerability
- Violating the repository pattern
- Putting database logic in the wrong layer
- Using string concatenation instead of prepared statements

The agent isn't stupid. It just didn't know your rules.

## Constraints as Agent Context

DECIDER makes constraints queryable:

```bash
$ decider check diff --base main --format json
{
  "applicable_adrs": [
    {
      "adr_id": "ADR-0001",
      "title": "Use PostgreSQL for persistence",
      "constraints": [
        "All database access must go through the repository pattern",
        "Use prepared statements for all queries",
        "Never expose raw SQL outside the db package"
      ]
    }
  ]
}
```

An agent that queries this before coding can generate:

```go
func (h *Handler) GetUserByEmail(ctx context.Context, email string) (*User, error) {
    return h.userRepo.FindByEmail(ctx, email)
}
```

Same task. Different result. The difference is context.

## What Good Constraints Look Like

Constraints should be:

### Specific
**Bad:** "Write secure code"
**Good:** "Use prepared statements for all database queries"

### Actionable
**Bad:** "Follow best practices"
**Good:** "All public API endpoints must validate input using the schema package"

### Scoped
**Bad:** "Always use TypeScript" (applies everywhere)
**Good:** "Frontend components must use TypeScript" (scoped to `src/frontend/**`)

### Verifiable
**Bad:** "Code should be efficient"
**Good:** "Database queries must use indexes; explain any full table scans in comments"

## The Constraint Hierarchy

Not all rules are equal. DECIDER distinguishes:

### Constraints
Rules that MUST be followed. Violations are bugs.
- "All database access must go through repositories"
- "Never store passwords in plain text"
- "API endpoints must be versioned"

### Invariants
Properties that must ALWAYS be true. The system breaks if they don't hold.
- "Database connections are pooled"
- "User sessions expire after 24 hours"
- "All migrations are reversible"

This distinction helps agents prioritize. Constraints guide implementation. Invariants define correctness.

## How Agents Should Use ADRs

An effective agent workflow:

### 1. Query Before Coding
```bash
decider check diff --base main
```

Find ADRs that apply to the files being changed.

### 2. Read Constraints
Parse the constraints from applicable ADRs. These are hard requirements.

### 3. Follow the Rules
Generate code that complies with every constraint. If a constraint is unclear, ask for clarification.

### 4. Document New Decisions
If the task requires an architectural choice not covered by existing ADRs, suggest creating one.

## The Feedback Loop

DECIDER creates a feedback loop:

1. **Team documents decision** → ADR with constraints
2. **Agent queries ADR** → Understands rules
3. **Agent writes compliant code** → Follows constraints
4. **Review catches violations** → Team refines constraints
5. **Constraints improve** → Better agent behavior

Over time, ADRs become more precise and agents become more effective.

## What Agents Can't Do

DECIDER surfaces constraints. It doesn't enforce them semantically.

An agent can read "use prepared statements" but might still make mistakes. That's what code review is for.

The value is in making constraints explicit. Implicit rules can't be followed. Explicit rules can.

## Getting Started with Agent Integration

1. **Document one critical constraint**: The rule agents violate most often
2. **Add scope paths**: Be specific about where it applies
3. **Tell agents to query**: Add instructions to AGENTS.md
4. **Review and refine**: Improve constraints based on what agents get wrong

You don't need comprehensive coverage to get value. Start with the painful cases and expand.

## The Alternative

Without constraints:
- Agents generate plausible but wrong code
- Reviews become adversarial ("you violated our unwritten rules")
- Technical debt accumulates
- Team velocity suffers

With constraints:
- Agents generate compliant code
- Reviews focus on logic, not patterns
- Architecture stays coherent
- Team velocity increases

The choice is between implicit rules that agents can't follow and explicit constraints they can.
