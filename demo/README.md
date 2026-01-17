# DECIDER Demo: Learn the Methodology

This demo teaches you the DECIDER methodology end-to-end in 5-10 minutes.

You'll learn how to:
- Install decider (pinned, checksum-verified)
- Validate ADRs
- Find which ADRs apply to code changes
- Understand the mandatory rationale pattern
- Use ADR constraints to guide development

## Scenario

This demo simulates an **e-commerce application** with documented decisions about:
- **Database layer**: PostgreSQL with repository pattern (ADR-0001)
- **API design**: REST with versioning (ADR-0002)
- **Frontend state**: Redux - deprecated (ADR-0003)

## Quick Start (5 Commands)

```bash
# 1. Install decider (pinned version with checksum verification)
./scripts/install-decider.sh
export PATH="$PWD/tools/decider:$PATH"

# 2. Validate all ADRs
decider check adr --dir docs/adr --strict

# 3. List all ADRs
decider list --dir docs/adr

# 4. Find ADRs that apply to database code
decider list --dir docs/adr --path "src/db/users.go"

# 5. View constraints for ADR-0001
decider show ADR-0001 --dir docs/adr
```

**What success looks like:**
- Step 2 shows: `✓ All 3 ADR(s) are valid`
- Step 4 shows ADR-0001 applies to `src/db/users.go`
- Step 5 shows the constraints you must follow

---

## Full Walkthrough

### Step 1: Install decider

Install the pinned version with checksum verification:

```bash
# macOS/Linux
./scripts/install-decider.sh
export PATH="$PWD/tools/decider:$PATH"

# Windows PowerShell
.\scripts\install-decider.ps1
$env:PATH = "$PWD\tools\decider;$env:PATH"
```

**Expected output:**
```
Installing decider v0.1.0 for darwin/arm64...
Checksum verified.
Installed to: /path/to/demo/tools/decider/decider

decider version 0.1.0
```

### Step 2: Validate All ADRs

Run the ADR pre-flight check:

```bash
decider check adr --dir docs/adr --strict
```

**Expected output:**
```
Validating ADRs in docs/adr/...
✓ 0001-use-postgresql-for-persistence.md: valid
✓ 0002-rest-api-with-versioning.md: valid
✓ 0003-redux-for-frontend-state.md: valid
✓ All 3 ADR(s) are valid
```

This verifies:
- All ADRs have valid YAML frontmatter
- Required sections are present
- The mandatory rationale pattern is followed

### Step 3: List All ADRs

See what decisions are documented:

```bash
decider list --dir docs/adr
```

**Expected output:**
```
ADR-0001  adopted     Use PostgreSQL for Persistence       [database, storage, infrastructure]
ADR-0002  adopted     REST API with Versioning             [api, http, architecture]
ADR-0003  deprecated  Redux for Frontend State Management  [frontend, state, react]
```

### Step 4: Find ADRs That Apply to Your Changes

Before modifying `src/db/users.go`, find applicable ADRs:

```bash
decider list --dir docs/adr --path "src/db/users.go"
```

**Expected output:**
```
ADR-0001  adopted  Use PostgreSQL for Persistence  [database, storage, infrastructure]
```

This tells you: ADR-0001's constraints apply to database code.

### Step 5: View ADR Details and Constraints

Read the constraints you must follow:

```bash
decider show ADR-0001 --dir docs/adr
```

**Expected output:**
```
ADR-0001: Use PostgreSQL for Persistence

Status: adopted
Date:   2026-01-10

Tags: database, storage, infrastructure

Scope:
  - src/db/**
  - migrations/**
  - docker-compose.yaml

Constraints:
  - All database access must go through the repository pattern
  - Use prepared statements for all queries
  - Never expose raw SQL outside the db package

Invariants:
  - Database connections are pooled
  - All migrations are reversible
  - Foreign keys enforce referential integrity
```

---

## Change → ADR → Compliance Flow

Here's how to use DECIDER when making changes:

### Scenario: Add a "Get User by Email" Function

**1. Find applicable ADRs:**
```bash
decider list --dir docs/adr --path "src/db/users.go"
```
Result: ADR-0001 applies.

**2. Read the constraints:**
```bash
decider show ADR-0001 --dir docs/adr
```
Constraints:
- All database access must go through the repository pattern
- Use prepared statements for all queries
- Never expose raw SQL outside the db package

**3. Write compliant code:**

❌ **Non-compliant** (raw SQL in handler):
```go
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE email = '" + email + "'")
    // Violates: repository pattern, prepared statements
}
```

✅ **Compliant** (uses repository pattern):
```go
func (h *Handler) GetUserByEmail(ctx context.Context, email string) (*User, error) {
    return h.userRepo.FindByEmail(ctx, email)  // Repository handles SQL
}
```

**4. Verify compliance:**
```bash
decider check adr --dir docs/adr --strict
```

---

## Understanding the Rationale Pattern

Open `docs/adr/0001-use-postgresql-for-persistence.md` to see the mandatory pattern:

### Adopted Option
```markdown
### PostgreSQL with Repository Pattern: Adopted

**Adopted because:**
- ACID compliance ensures financial transactions are never corrupted
- Rich feature set supports complex queries
- Repository pattern improves testability

**Adopted despite:**
- Requires operational expertise for production
- More complex setup than embedded databases
```

### Rejected Alternative
```markdown
### MongoDB: Rejected

**Rejected because:**
- No multi-document ACID transactions
- Eventual consistency unsuitable for financial data

**Rejected despite:**
- Flexible schema would ease initial development
- Horizontal scaling is simpler
```

This pattern ensures every decision documents **why** choices were made.

---

## Additional Commands

### Filter by Status
```bash
decider list --dir docs/adr --status adopted
```

### Filter by Tag
```bash
decider list --dir docs/adr --tag database
```

### Get JSON Output
```bash
decider list --dir docs/adr --format json
```

### Check Index Consistency
```bash
decider index --dir docs/adr --check
```

### Regenerate Index
```bash
decider index --dir docs/adr
```

---

## Directory Structure

```
demo/
├── README.md                    # This walkthrough
├── AGENTS.md                    # Agent instructions (tool-agnostic)
├── CLAUDE.md                    # Claude-specific instructions
├── scripts/
│   ├── install-decider.sh       # macOS/Linux installer
│   └── install-decider.ps1      # Windows installer
├── tools/
│   ├── decider.version          # Pinned version: v0.1.0
│   └── decider/                 # Installed binary
├── docs/adr/
│   ├── README.md                # ADR system overview
│   ├── index.yaml               # Auto-generated index
│   ├── templates/adr.md         # ADR template
│   ├── 0001-use-postgresql-for-persistence.md
│   ├── 0002-rest-api-with-versioning.md
│   └── 0003-redux-for-frontend-state.md
└── src/                         # Demo source directories
    ├── db/
    ├── api/
    └── frontend/
```

---

## What You Learned

1. **Pinned installation**: Install decider with checksum verification
2. **Pre-flight checks**: Run `decider check adr --strict` before working
3. **Scope-based queries**: Find ADRs that apply to specific files
4. **Constraint-driven development**: Read and follow ADR constraints
5. **Rationale pattern**: Document decisions with "Adopted because/despite" format

## Next Steps

- Read the [main DECIDER documentation](../README.md)
- Try creating a new ADR: `decider new "Add Caching" --dir docs/adr --tags performance`
- Explore the [ADR Steward agent guide](../docs/guides/using-adr-steward.md)
