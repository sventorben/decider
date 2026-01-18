---
adr_id: ADR-0006
title: "TOON as Default Structured Format"
status: adopted
date: 2026-01-18
scope:
  paths:
    - "cmd/**"
    - "internal/**"
    - "docs/**"
    - "demo/**"
    - ".claude/**"
    - "SPEC.md"
tags:
  - format
  - cli
  - output
constraints:
  - TOON must be the default for machine-readable structured output
  - JSON must remain fully supported via explicit format selection
  - YAML support for index files must remain unchanged
  - All structured outputs must be deterministic
invariants:
  - Format selection is explicit and predictable
  - Backward compatibility with JSON-based tooling is preserved
  - Round-trip encoding/decoding produces equivalent data
---

# TOON as Default Structured Format

## Context

DECIDER produces structured output in several contexts:
- CLI command outputs (list, show, check, explain) with `--format` flag
- Index files (`index.yaml`)
- Machine-readable error information

Currently, JSON is the default machine-readable format. TOON (Token-Oriented Object Notation) is a compact encoding of the JSON data model that offers better token efficiency for LLM consumption while maintaining full compatibility with the JSON data model.

As AI agents become primary consumers of DECIDER output, optimizing for token efficiency without sacrificing interoperability is valuable.

## Decision

Adopt TOON as the default structured format for machine-readable output, while maintaining full JSON and YAML support.

### TOON with Deterministic Output: Adopted

**Adopted because:**
- TOON encodes the same data model as JSON (objects, arrays, strings, numbers, booleans, null)
- More token-efficient for LLM consumption, reducing API costs and context usage
- Deterministic output enables reliable diffing and caching
- Human-readable while being more compact than JSON
- Full round-trip compatibility with JSON ensures no data loss

**Adopted despite:**
- Less familiar to developers than JSON
- Requires implementation of encoder/decoder (no Go stdlib support)
- Tooling ecosystem for TOON is less mature than JSON
- May require format conversion for integration with JSON-only tools

## Alternatives Considered

### JSON-Only (Status Quo): Rejected

**Rejected because:**
- Higher token cost for LLM consumers
- No benefit for the primary use case (AI agent consumption)
- Missed opportunity to optimize for emerging workflows

**Rejected despite:**
- Universal tooling support
- Developer familiarity
- No implementation effort required

### YAML-Only for All Outputs: Rejected

**Rejected because:**
- YAML parsing is more complex and error-prone
- Indentation-sensitive format increases risk of corruption
- Larger output size than both JSON and TOON
- Not well-suited for programmatic generation

**Rejected despite:**
- Human readability for configuration
- Already used for index files
- Good support in Go ecosystem

### TOON Optional but Not Default: Rejected

**Rejected because:**
- Would not realize token efficiency benefits by default
- AI agents would need explicit configuration to benefit
- Inconsistent with goal of optimizing for LLM consumption

**Rejected despite:**
- Lower risk approach
- No breaking change for existing users
- Easier migration path

## Consequences

### Positive
- Reduced token usage for AI agent consumers
- Deterministic outputs improve caching and reproducibility
- Full backward compatibility via `--format=json`
- Clear format selection model: `--format=toon|json|yaml`

### Negative
- Users expecting JSON by default must update scripts
- Additional code to maintain (TOON encoder/decoder)
- Documentation must explain format selection

### Neutral
- Index files continue to use YAML by default
- Text format remains default for human-oriented output

## Agent Guidance

When processing DECIDER output:
1. Default format is TOON for structured data
2. Use `--format=json` if JSON is required for downstream tools
3. TOON and JSON represent identical data structures
4. Parse TOON using standard TOON parsing rules (whitespace-separated tokens)
