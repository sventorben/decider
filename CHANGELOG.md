# Changelog

All notable changes to DECIDER will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-01-17

Initial release of DECIDER, a Git-native system for managing Architecture Decision Records with machine-readable constraints.

### Added

#### CLI Commands
- `decider init` - Initialize ADR directory structure with templates
- `decider new` - Create new ADRs with YAML frontmatter
- `decider index` - Generate/update machine-readable index
- `decider list` - List ADRs with status, tag, and path filtering
- `decider show` - Display ADR details
- `decider check adr` - Validate ADR format compliance
- `decider check diff` - Find ADRs applicable to git changes
- `decider explain` - Explain why ADRs apply to changes
- `decider version` - Show version information

#### ADR Format
- Markdown with YAML frontmatter for human + machine readability
- Structured fields: adr_id, title, status, date, scope, tags, constraints, invariants
- Status lifecycle: proposed, adopted, rejected, deprecated, superseded
- Scope paths with glob pattern matching
- Supersession tracking via supersedes/superseded_by fields

#### Agent Integration
- ADR Steward agent for Claude Code (`.claude/agents/adr-steward.md`)
- Slash commands: `/adr-new`, `/adr-supersede`
- Machine-readable JSON output for automation
- Constraint and invariant querying for AI agents

#### Security
- Input validation for all CLI arguments
- Git ref validation to prevent injection attacks
- File size limits (10MB) to prevent memory exhaustion
- Glob pattern depth limits (max 10 `**` segments)
- Path sanitization utilities

#### Documentation
- Comprehensive README with quickstart and examples
- Full specification in SPEC.md
- Guides: getting-started, writing-adrs, ci-integration, agent-integration
- Philosophy essays in docs/why/
- Security review documentation

#### Infrastructure
- Cross-platform support (Linux, macOS, Windows)
- GitHub Actions CI workflow
- GoReleaser configuration for releases
- Dependabot configuration for dependency updates

[Unreleased]: https://github.com/sventorben/decider/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/sventorben/decider/releases/tag/v0.1.0
