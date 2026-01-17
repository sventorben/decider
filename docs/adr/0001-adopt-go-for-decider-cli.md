---
adr_id: ADR-0001
title: Adopt Go for DECIDER CLI
status: adopted
date: 2026-01-16
scope:
  paths:
    - "cmd/**"
    - "internal/**"
    - "pkg/**"
    - "go.mod"
    - "go.sum"
tags:
  - toolchain
  - language
  - foundation
constraints:
  - Use Go 1.25 or later for generics and improved stdlib
  - Keep external dependencies minimal
  - Prefer stdlib flag package over heavy CLI frameworks unless complexity warrants it
invariants:
  - All code compiles with `go build ./...`
  - All tests pass with `go test ./...`
  - Cross-platform compatibility (Linux, macOS, Windows)
supersedes: []
superseded_by: []
related_adrs: []
---

# ADR-0001: Adopt Go for DECIDER CLI

## Context

DECIDER needs a language and runtime for building a cross-platform CLI tool that:
- Produces single static binaries with no runtime dependencies
- Has excellent tooling for testing, linting, and formatting
- Can be easily distributed via GitHub releases
- Has good support for YAML parsing and file operations

Decision drivers:
- Distribution simplicity (users should not need to install runtimes)
- Cross-platform builds without complex toolchains
- Fast CI builds
- Maintainability by a small team or solo maintainer

Candidates considered: Go, Rust, Python, TypeScript/Node.

## Decision

We adopt **Go 1.25+** as the implementation language for DECIDER.

We will use the **stdlib `flag` package** for CLI parsing to minimize dependencies. If complexity grows significantly (subcommands with shared flags, auto-completion), we may revisit and adopt Cobra in a future ADR.

### Go: Adopted

**Adopted because:**
- Single binary distribution with static linking eliminates runtime dependencies
- Built-in cross-compilation via `GOOS`/`GOARCH` makes multi-platform releases trivial
- Fast compilation enables quick iteration and CI builds under 2 minutes
- Mature YAML library (`gopkg.in/yaml.v3`) is well-maintained and sufficient
- GoReleaser provides industry-standard release automation for Go projects
- Simple concurrency model if parallel ADR processing is needed later

**Adopted despite:**
- Go is more verbose than scripting languages; simple operations require more code
- Error handling is repetitive (`if err != nil` patterns)
- Generics (added in 1.18) are less mature than in other languages
- Team must know Go, though it's a common skill with low learning curve

## Alternatives Considered

### Rust: Rejected

**Rejected because:**
- Compilation times are significantly slower than Go, impacting CI feedback loops
- Steeper learning curve (ownership, lifetimes) increases contributor friction
- Build toolchain is heavier than Go's single binary distribution
- Performance benefits are unnecessary for a CLI that parses YAML and runs git commands

**Rejected despite:**
- Memory safety guarantees are stronger than Go's
- Excellent CLI frameworks (clap) with auto-generated help and completions
- Growing popularity and ecosystem
- Performance would be marginally better for large ADR sets

### Python: Rejected

**Rejected because:**
- Requires Python runtime to be installed on user machines
- Packaging and distribution (pip, venv, pyinstaller) adds friction and complexity
- Execution speed is noticeably slower for CLI responsiveness
- Version compatibility issues (Python 2 vs 3, minor version differences)

**Rejected despite:**
- Fastest development velocity for initial prototyping
- Excellent YAML libraries (PyYAML, ruamel.yaml)
- Large pool of potential contributors familiar with Python
- Rich ecosystem for text processing and file manipulation

### TypeScript/Node: Rejected

**Rejected because:**
- Requires Node.js runtime, adding ~100MB dependency for users
- Packaging with pkg or similar tools produces larger binaries than Go
- npm ecosystem has security concerns (dependency supply chain)
- Performance is slower than compiled languages for CLI cold starts

**Rejected despite:**
- Familiar to web developers who may use DECIDER
- Strong typing with TypeScript catches errors at compile time
- Rich ecosystem for YAML, Markdown, and JSON processing
- npm distribution would be convenient for JavaScript-heavy teams

## Consequences

**Positive:**
- Simple build and release pipeline with GoReleaser
- Easy for contributors familiar with Go to onboard
- Excellent cross-platform support without conditional compilation
- Single ~10MB binary for distribution

**Negative:**
- Contributors must know Go (mitigated: common skill, good documentation)
- More verbose than scripting languages (acceptable for maintainability)
- Limited metaprogramming compared to dynamic languages

## Agent Guidance

When modifying Go code:
- Run `go fmt` before committing
- Run `go vet ./...` to catch common issues
- Keep functions small and testable
- Prefer returning errors over panicking
