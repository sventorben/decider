---
adr_id: ADR-0004
title: Release Process with GoReleaser and GitHub Actions
status: adopted
date: 2026-01-16
scope:
  paths:
    - ".github/workflows/**"
    - ".goreleaser.yaml"
    - "cmd/decider/**"
tags:
  - ci
  - release
  - automation
constraints:
  - Releases triggered only by version tags (v*.*.*)
  - All releases must pass CI before publishing
  - Binaries built for linux/darwin/windows on amd64/arm64
invariants:
  - Version info embedded via ldflags at build time
  - Checksums generated for all release artifacts
  - GitHub Release created automatically on tag push
supersedes: []
superseded_by: []
related_adrs:
  - ADR-0001
---

# ADR-0004: Release Process with GoReleaser and GitHub Actions

## Context

DECIDER needs a reliable, automated release process that:
1. Builds cross-platform binaries
2. Embeds version information
3. Creates GitHub Releases with artifacts
4. Generates checksums for verification
5. Requires minimal manual intervention

Decision drivers:
- Releases should be fully automated from tag push to published artifacts
- Cross-platform builds must be reproducible and consistent
- Version information must be embedded at build time
- Process must work for a solo maintainer or small team

## Decision

We adopt **GoReleaser** with **GitHub Actions** for the release pipeline.

### GoReleaser + GitHub Actions: Adopted

**Adopted because:**
- GoReleaser is the industry standard for Go binary releases
- Native integration with GitHub Actions via official action
- Cross-compilation is handled automatically with correct settings
- Checksum files generated automatically (SHA256)
- Changelog generation from git history
- ldflags injection for version/commit/date embedding

**Adopted despite:**
- Adds dependency on external tool (GoReleaser)
- Configuration file (.goreleaser.yaml) adds complexity
- Requires understanding GoReleaser's templating system
- GitHub Actions lock-in (acceptable for GitHub-hosted project)

### Version Tagging

- Use semantic versioning: `vMAJOR.MINOR.PATCH`
- Tags must match pattern `v*.*.*` (e.g., `v0.1.0`, `v1.2.3`)
- Pre-release tags: `v0.1.0-rc.1`, `v1.0.0-beta.2`

### CI Pipeline (ci.yml)

Triggers: push to any branch, pull requests

Steps:
1. Checkout code
2. Setup Go 1.25+
3. Run `go test ./...`
4. Run `golangci-lint`
5. Build CLI: `go build ./cmd/decider`
6. Run `decider check adr` to validate ADRs

### Release Pipeline (release.yml)

Triggers: push tags matching `v*.*.*`

Steps:
1. Checkout code
2. Setup Go 1.25+
3. Run GoReleaser with `--clean` flag

### GoReleaser Configuration

```yaml
# .goreleaser.yaml
builds:
  - main: ./cmd/decider
    binary: decider
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}
```

### Version Command

`decider version` outputs:

```
decider version 0.1.0
  commit: abc1234
  built:  2026-01-16T10:30:00Z
```

### Release Checklist

1. Update CHANGELOG.md with release notes
2. Ensure all tests pass
3. Create and push tag: `git tag v0.1.0 && git push origin v0.1.0`
4. GitHub Actions builds and publishes release
5. Verify release artifacts on GitHub Releases page

## Alternatives Considered

### Manual Releases: Rejected

**Rejected because:**
- Error-prone (wrong flags, missing platforms, typos)
- Time-consuming for cross-platform builds
- No reproducibility guarantee between releases
- Requires maintainer to have all build tools locally

**Rejected despite:**
- Full control over every aspect of the release
- No dependency on external services or tools
- Can work without GitHub Actions availability
- Simpler to understand for first release

### Makefile + Shell Scripts: Rejected

**Rejected because:**
- More code to maintain than declarative GoReleaser config
- Shell script portability issues (bash vs sh, macOS vs Linux)
- Checksum generation must be implemented manually
- GitHub Release upload requires additional tooling (gh CLI)

**Rejected despite:**
- Flexible and customizable for unusual requirements
- No external tool dependencies beyond standard Unix tools
- Easier to debug with echo/print statements
- Works with any CI system, not just GitHub Actions

### GitHub Actions Only (No GoReleaser): Rejected

**Rejected because:**
- More verbose YAML configuration for same result
- Cross-compilation matrix must be defined manually
- Checksum file generation requires additional steps
- ldflags injection requires manual shell scripting
- Less Go-specific optimization (e.g., CGO handling)

**Rejected despite:**
- Fewer moving parts (just GitHub Actions)
- More explicit about what each step does
- No dependency on GoReleaser availability or updates
- Easier to customize individual build steps

## Consequences

**Positive:**
- Fully automated releases from tag to published artifacts
- Consistent cross-platform builds every time
- Industry-standard tooling familiar to Go developers
- Checksums and signatures support built-in

**Negative:**
- GoReleaser learning curve for configuration changes (minimal)
- Depends on GitHub Actions availability (acceptable)
- GoReleaser updates may require config changes (rare)

## Agent Guidance

When preparing releases:
- Never create tags directly; follow the release checklist
- Update CHANGELOG.md before tagging
- Use semantic versioning strictly
- Test the release build locally with `goreleaser build --snapshot`
