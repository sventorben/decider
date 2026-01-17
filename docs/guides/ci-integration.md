# CI Integration Guide

This guide shows how to integrate DECIDER validation into your CI pipeline.

## Why CI Integration?

CI checks ensure:
- All ADRs are valid (correct format, required fields)
- The index stays in sync with ADR files
- PRs that modify ADRs don't break the system

## GitHub Actions

### Basic Validation

Add this workflow to `.github/workflows/adr.yml`:

```yaml
name: ADR Validation

on:
  push:
    branches: [main]
    paths:
      - 'docs/adr/**'
  pull_request:
    branches: [main]
    paths:
      - 'docs/adr/**'

jobs:
  validate:
    name: Validate ADRs
    runs-on: ubuntu-24.04
    steps:
      - name: Checkout
        uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

      - name: Setup Go
        uses: actions/setup-go@d35c59abb061a4a6fb18e82ac0862c26744d6ab5 # v5.5.0
        with:
          go-version-file: go.mod

      - name: Install decider
        run: go install github.com/sventorben/decider/cmd/decider@v0.1.0

      - name: Validate ADRs
        run: decider check adr

      - name: Check index is up-to-date
        run: decider index --check
```

### Full Workflow with ADR Applicability

For PRs that touch code, show which ADRs apply:

```yaml
name: CI

on:
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-24.04
    steps:
      - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2
        with:
          fetch-depth: 0  # Need history for diff

      - uses: actions/setup-go@d35c59abb061a4a6fb18e82ac0862c26744d6ab5 # v5.5.0
        with:
          go-version-file: go.mod
          cache: true

      - name: Install decider
        run: go install github.com/sventorben/decider/cmd/decider@v0.1.0

      - name: Run tests
        run: go test ./...

      - name: Show applicable ADRs
        run: |
          echo "## Applicable ADRs" >> $GITHUB_STEP_SUMMARY
          echo "" >> $GITHUB_STEP_SUMMARY
          decider explain --base origin/main >> $GITHUB_STEP_SUMMARY || echo "No ADRs apply to these changes" >> $GITHUB_STEP_SUMMARY

  validate-adrs:
    runs-on: ubuntu-24.04
    steps:
      - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

      - uses: actions/setup-go@d35c59abb061a4a6fb18e82ac0862c26744d6ab5 # v5.5.0
        with:
          go-version-file: go.mod

      - name: Install decider
        run: go install github.com/sventorben/decider/cmd/decider@v0.1.0

      - name: Validate ADRs
        run: decider check adr

      - name: Check index
        run: decider index --check
```

### Using Pre-built Binaries

If you don't want to install via Go (replace v0.1.0 with the latest release version):

```yaml
- name: Install decider
  run: |
    curl -L https://github.com/sventorben/decider/releases/download/v0.1.0/decider_Linux_x86_64.tar.gz | tar xz
    sudo mv decider /usr/local/bin/
```

## Exit Codes

DECIDER uses specific exit codes for CI:

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (parse error, missing file, etc.) |
| 2 | Validation failure (invalid ADR, stale index) |

Use these in CI to distinguish between errors and validation failures:

```yaml
- name: Validate ADRs
  run: decider check adr
  continue-on-error: false

- name: Check index
  run: |
    if ! decider index --check; then
      echo "::error::ADR index is out of date. Run 'decider index' and commit."
      exit 1
    fi
```

## PR Comments

Add a step to comment on PRs with applicable ADRs:

```yaml
- name: Comment applicable ADRs
  if: github.event_name == 'pull_request'
  uses: actions/github-script@60a0d83039c74a4aee543508d2ffcb1c3799cdea # v7.0.1
  with:
    script: |
      const { execSync } = require('child_process');
      let output;
      try {
        output = execSync('decider explain --base origin/main').toString();
      } catch (e) {
        output = 'No ADRs apply to these changes.';
      }

      github.rest.issues.createComment({
        issue_number: context.issue.number,
        owner: context.repo.owner,
        repo: context.repo.repo,
        body: `## Applicable ADRs\n\n${output}`
      });
```

## Required Status Checks

After setting up the workflow:

1. Go to repository Settings > Branches
2. Add a branch protection rule for `main`
3. Enable "Require status checks to pass"
4. Select `validate-adrs` as required

## Other CI Systems

### GitLab CI

```yaml
validate-adrs:
  image: golang:1.25-bookworm
  stage: test
  script:
    - go install github.com/sventorben/decider/cmd/decider@v0.1.0
    - decider check adr
    - decider index --check
  rules:
    - changes:
        - docs/adr/**
```

### CircleCI

```yaml
version: 2.1

jobs:
  validate-adrs:
    docker:
      - image: cimg/go:1.25.0
    steps:
      - checkout
      - run:
          name: Install decider
          command: go install github.com/sventorben/decider/cmd/decider@v0.1.0
      - run:
          name: Validate ADRs
          command: decider check adr
      - run:
          name: Check index
          command: decider index --check

workflows:
  main:
    jobs:
      - validate-adrs
```

### Jenkins

```groovy
pipeline {
    agent any
    stages {
        stage('Validate ADRs') {
            steps {
                sh 'go install github.com/sventorben/decider/cmd/decider@v0.1.0'
                sh 'decider check adr'
                sh 'decider index --check'
            }
        }
    }
}
```

## Common Issues

### Index Out of Date

```
Error: index is out of date
```

**Solution:** Run `decider index` locally and commit the updated `index.yaml`.

### Missing Required Sections

```
Error: 0005-my-adr.md: missing required section: Alternatives Considered
```

**Solution:** Add the missing section to the ADR.

### Invalid Frontmatter

```
Error: 0005-my-adr.md: invalid status: "approved"
```

**Solution:** Use valid status values: `proposed`, `adopted`, `rejected`, `deprecated`, `superseded`.

## Best Practices

1. **Run validation early**: Add ADR checks to the same workflow as tests
2. **Fail fast**: Don't let invalid ADRs merge to main
3. **Show context**: Use `decider explain` to surface applicable ADRs in PRs
4. **Cache Go modules**: Speeds up `go install` significantly
5. **Pin decider version**: Consider pinning to a specific release for reproducibility

```yaml
- name: Install decider
  run: go install github.com/sventorben/decider/cmd/decider@v0.1.0
```
