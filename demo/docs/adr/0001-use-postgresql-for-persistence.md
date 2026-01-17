---
adr_id: ADR-0001
title: "Use PostgreSQL for Persistence"
status: adopted
date: 2026-01-10
scope:
  paths:
    - "src/db/**"
    - "migrations/**"
    - "docker-compose.yaml"
tags:
  - database
  - storage
  - infrastructure
constraints:
  - "All database access must go through the repository pattern"
  - "Use prepared statements for all queries"
  - "Never expose raw SQL outside the db package"
invariants:
  - "Database connections are pooled"
  - "All migrations are reversible"
  - "Foreign keys enforce referential integrity"
supersedes: []
superseded_by: []
related_adrs:
  - ADR-0002
---

# ADR-0001: Use PostgreSQL for Persistence

## Context

Our e-commerce application needs a reliable database for storing:
- User accounts and authentication data
- Product catalog with inventory
- Order history and transactions
- Shopping cart state

Decision drivers:
- Financial transactions require ACID compliance
- Complex reporting queries need rich SQL feature set
- Production system must handle concurrent traffic reliably
- Team has existing PostgreSQL operational expertise

## Decision

We will use **PostgreSQL 15+** as our primary database with the **repository pattern** for data access.

### PostgreSQL with Repository Pattern: Adopted

**Adopted because:**
- ACID compliance ensures financial transactions are never corrupted
- Rich feature set (JSON columns, CTEs, window functions) supports complex queries
- Mature ecosystem with excellent tooling (pgAdmin, pg_dump, replication)
- Repository pattern isolates SQL from business logic, improving testability
- Prepared statements with repository pattern eliminate SQL injection by design
- Strong community support and extensive documentation

**Adopted despite:**
- Requires operational expertise for production deployment (backups, replication)
- More complex setup than embedded databases for development
- Schema migrations require careful planning and testing
- Repository pattern adds boilerplate compared to raw SQL access

## Alternatives Considered

### MySQL: Rejected

**Rejected because:**
- JSON support is less mature than PostgreSQL's JSONB
- Fewer advanced SQL features (CTEs, window functions historically weaker)
- Default storage engine behavior can be surprising (MyISAM vs InnoDB)
- Less flexible constraint handling

**Rejected despite:**
- Widely adopted with large talent pool
- Good performance for read-heavy workloads
- Mature replication and clustering options
- Lower memory footprint for simple use cases

### MongoDB: Rejected

**Rejected because:**
- No multi-document ACID transactions (critical for order processing)
- Eventual consistency is unsuitable for financial data
- Schema flexibility leads to data inconsistency over time
- Query language less powerful than SQL for complex reporting

**Rejected despite:**
- Flexible schema would ease initial development
- Horizontal scaling is simpler than PostgreSQL
- Native JSON storage matches some data formats
- Good performance for document-oriented access patterns

### SQLite: Rejected

**Rejected because:**
- Single-writer limitation prevents concurrent web traffic handling
- No network access means no separate database server
- Limited scalability for production workloads
- Lacks enterprise features (replication, point-in-time recovery)

**Rejected despite:**
- Zero configuration required for development
- Excellent for local testing and prototyping
- No external dependencies or infrastructure
- Fast for single-user scenarios

## Consequences

**Positive:**
- ACID compliance for order transactions eliminates data corruption risk
- Rich feature set enables complex analytics queries
- Mature ecosystem with excellent tooling and community support
- Repository pattern makes data layer highly testable

**Negative:**
- Requires operational knowledge for production (mitigated by team experience)
- More complex than NoSQL for simple document storage (acceptable trade-off)
- Repository boilerplate increases code volume (offset by testability gains)

## Agent Guidance

When working in `src/db/`:
- Never write raw SQL in handlers or services; use repository methods
- Always use parameterized queries (`$1`, `$2`) not string concatenation
- Add indexes for columns used in WHERE clauses
- Write reversible migrations (both UP and DOWN)
