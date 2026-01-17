---
adr_id: ADR-0002
title: "REST API with Versioning"
status: adopted
date: 2026-01-12
scope:
  paths:
    - "src/api/**"
    - "src/handlers/**"
    - "openapi/**"
tags:
  - api
  - http
  - architecture
constraints:
  - "All endpoints must be versioned with /api/v{N}/ prefix"
  - "Breaking changes require a new API version"
  - "Deprecate old versions for at least 6 months before removal"
invariants:
  - "All responses use JSON format"
  - "HTTP status codes follow REST conventions"
  - "All endpoints are documented in OpenAPI spec"
supersedes: []
superseded_by: []
related_adrs:
  - ADR-0001
---

# ADR-0002: REST API with Versioning

## Context

Our e-commerce platform needs a public API for:
- Mobile applications (iOS, Android)
- Third-party integrations
- Internal frontend SPA

Decision drivers:
- API must evolve without breaking existing clients
- Multiple client types with different update cycles
- Team familiarity with REST patterns
- Need for clear, stable contracts

## Decision

We will implement a **RESTful API with URL-based versioning**.

### REST with URL Versioning: Adopted

**Adopted because:**
- URL versioning is immediately visible and discoverable in documentation
- Clear contract between client and server at each version level
- Easy to test (each version is a distinct URL path)
- Clients can migrate at their own pace across versions
- REST is well-understood by team and third-party developers
- Extensive tooling ecosystem (Postman, OpenAPI, client generators)

**Adopted despite:**
- Must maintain multiple versions during transition periods
- Potential code duplication between version handlers
- URL pollution with version numbers
- Does not solve field-level evolution (still need deprecation strategy)

Version format: `/api/v1/`, `/api/v2/`, etc.

Versioning rules:
1. Non-breaking changes (new fields, new endpoints) go in current version
2. Breaking changes (removed fields, changed types) require new version
3. Old versions supported for minimum 6 months with deprecation warnings

## Alternatives Considered

### GraphQL: Rejected

**Rejected because:**
- Adds significant complexity for our current requirements
- Caching is more difficult than REST (POST-based queries)
- Learning curve for team members unfamiliar with GraphQL
- Overfetching/underfetching problems not significant for our use case
- Security surface area is larger (arbitrary queries)

**Rejected despite:**
- Flexible queries reduce API surface area
- Single endpoint simplifies routing
- Strong typing with schema introspection
- Efficient for mobile clients with bandwidth constraints
- Active ecosystem and tooling

### Header-Based Versioning: Rejected

**Rejected because:**
- Less discoverable than URL versioning
- Harder to test manually (requires header manipulation)
- Documentation tools handle it less gracefully
- Easy to forget header in client code

**Rejected despite:**
- Cleaner URLs without version numbers
- Same endpoint handles multiple versions (less routing complexity)
- Follows REST purist principles about URL structure
- Better for gradual migration (default version concept)

### No Versioning: Rejected

**Rejected because:**
- Any breaking change would break all existing clients immediately
- No graceful migration path for third parties
- Forces synchronized releases across all clients
- Creates support burden when changes cause failures

**Rejected despite:**
- Simplest implementation with no version logic
- Single codebase to maintain
- No deprecation tracking needed
- Works for internal-only APIs with controlled clients

## Consequences

**Positive:**
- Clear contract between client and server
- Easy to test and document (each version standalone)
- Clients can migrate at their own pace
- Third parties can depend on stable versioned endpoints

**Negative:**
- Multiple versions to maintain during transition (mitigated by 6-month deprecation policy)
- Possible code duplication between versions (mitigated by shared internal logic)
- Version sprawl if not managed carefully (mitigated by deprecation enforcement)

## Agent Guidance

When modifying API endpoints:
- Check if the change is breaking (removing fields, changing types)
- If breaking, create new version and deprecate old endpoint
- Update OpenAPI spec in `openapi/` directory
- Add deprecation headers to old endpoints: `Deprecation: true`
