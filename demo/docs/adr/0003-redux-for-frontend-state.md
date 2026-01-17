---
adr_id: ADR-0003
title: "Redux for Frontend State Management"
status: deprecated
date: 2026-01-08
scope:
  paths:
    - "src/frontend/store/**"
    - "src/frontend/reducers/**"
    - "src/frontend/actions/**"
tags:
  - frontend
  - state
  - react
constraints:
  - "All global state must go through Redux store"
  - "Use Redux Toolkit for reduced boilerplate"
  - "Async actions use createAsyncThunk"
invariants:
  - "State is immutable (never mutate directly)"
  - "All state changes are traceable via actions"
supersedes: []
superseded_by: []
related_adrs: []
---

# ADR-0003: Redux for Frontend State Management

**Status: DEPRECATED** - See Deprecation Notice below.

## Context

Our React frontend needs state management for:
- User authentication state
- Shopping cart contents
- Product catalog cache
- UI state (modals, notifications)

Decision drivers:
- Predictable state updates across complex UI
- Debugging capabilities for state-related bugs
- Team familiarity with Redux patterns
- Need for time-travel debugging during development

## Decision

We will use **Redux with Redux Toolkit** for global state management.

### Redux with Redux Toolkit: Adopted (Now Deprecated)

**Adopted because:**
- Predictable state updates via unidirectional data flow
- Excellent debugging with Redux DevTools and time-travel
- Redux Toolkit reduces boilerplate significantly vs vanilla Redux
- Large ecosystem with extensive middleware options
- Well-documented patterns for common scenarios
- Team had prior Redux experience

**Adopted despite:**
- Significant boilerplate even with Redux Toolkit
- Learning curve for developers new to flux patterns
- May be overkill for simple state requirements
- Action/reducer indirection adds complexity
- Performance overhead for frequent updates

Structure:
- `/store` - Redux store configuration
- `/reducers` - Slice reducers (using createSlice)
- `/actions` - Async thunks for API calls

## Alternatives Considered

### React Context: Rejected

**Rejected because:**
- Re-render performance issues at scale (entire tree re-renders on context change)
- No built-in devtools for debugging state changes
- Difficult to implement time-travel debugging
- State logic scattered across components

**Rejected despite:**
- Built into React (no additional dependency)
- Simpler mental model for small applications
- Less boilerplate than Redux
- Sufficient for truly simple state needs

### MobX: Rejected

**Rejected because:**
- "Magic" reactive behavior makes debugging harder
- Less explicit about state mutations
- Smaller ecosystem than Redux
- Team had no MobX experience

**Rejected despite:**
- Significantly less boilerplate than Redux
- More intuitive reactive programming model
- Better performance for frequent updates
- Simpler learning curve for OOP developers

### Zustand: Rejected

**Rejected because:**
- Newer library with smaller ecosystem at time of decision
- Fewer middleware and integration options
- Less tooling and documentation available
- Team familiarity favored Redux

**Rejected despite:**
- Minimal boilerplate with hooks-based API
- No providers required (simpler setup)
- Good TypeScript support
- Growing adoption and community

## Consequences

**Positive:**
- Predictable state updates make bugs reproducible
- Excellent debugging with Redux DevTools
- Time-travel debugging accelerates development
- Large ecosystem provides solutions for common patterns

**Negative:**
- Significant boilerplate even with RTK (confirmed in practice)
- Learning curve for new developers (higher than expected)
- Overkill for simple state (confirmed - most state is server-derived)

## Deprecation Notice

This ADR is deprecated as of 2026-01-15 because:

1. **Complexity vs. value**: Redux adds significant complexity for our relatively simple client-side state needs. Most "state" is actually server data that React Query handles better.

2. **Server state vs. client state**: We've adopted React Query for server state, which eliminates most Redux use cases (product catalog, user data, order history).

3. **Remaining state is simple**: Authentication status and UI state (modals, notifications) don't justify Redux's complexity.

**Guidance for new code:**
- Evaluate Zustand or React Context first
- Use React Query for server state
- Redux code will be maintained but not expanded
- Eventual migration planned but not scheduled

**Existing Redux code:**
- Continue using Redux patterns where already implemented
- No immediate migration required
- New features should not add to Redux store without explicit approval
