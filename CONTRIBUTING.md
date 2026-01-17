# Contributing to DECIDER

Thank you for your interest in contributing to DECIDER!

## Before You Start

**Please initiate a discussion before writing code.** This helps ensure your contribution aligns with the project direction and avoids wasted effort.

1. **Create a GitHub Issue** describing the change you want to make
2. **Wait for feedback** from maintainers before starting work
3. **Reference the issue** in your pull request

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR-USERNAME/decider.git`
3. Configure commit signing (see below)
4. Create a branch: `git checkout -b feature/your-feature`
5. Make your changes
6. Run tests: `go test ./...`
7. Commit with DCO sign-off
8. Push and create a Pull Request

## Development Setup

### Prerequisites

- Go 1.25 or later
- Git (with GPG signing configured)

### Building

```bash
# Build the CLI
go build -o decider ./cmd/decider

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run linter (if installed)
golangci-lint run
```

### Project Structure

```
decider/
├── cmd/decider/       # CLI entry point
├── internal/
│   ├── adr/           # ADR parsing and validation
│   ├── cli/           # Command implementations
│   ├── glob/          # Glob pattern matching
│   ├── index/         # Index generation
│   └── validate/      # Input validation
├── docs/
│   ├── adr/           # Project ADRs
│   ├── guides/        # How-to documentation
│   ├── reference/     # CLI reference
│   ├── security/      # Security reviews
│   └── why/           # Philosophy
├── demo/              # Demo project
└── .claude/           # Claude Code integration
```

### Validation Commands

```bash
# Validate ADR format
decider check adr

# Check index is up-to-date
decider index --check
```

## Signing Requirements

### Developer Certificate of Origin (DCO)

All commits must be signed off to certify that you have the right to submit the code under the project's license. Add the following to your commit message:

```
Signed-off-by: Your Name <your.email@example.com>
```

You can do this automatically with the `-s` flag:

```bash
git commit -s -m "feat(cli): add new feature"
```

By signing off, you certify the following (from [developercertificate.org](https://developercertificate.org/)):

> I certify that I have the right to submit this contribution under the open source license indicated in the file.

### Commit Signature Verification

All commits to the main branch require GitHub-verified signatures. Configure commit signing:

1. [Generate a GPG key](https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key)
2. [Add the GPG key to your GitHub account](https://docs.github.com/en/authentication/managing-commit-signature-verification/adding-a-gpg-key-to-your-github-account)
3. [Configure Git to sign commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)

```bash
# Configure Git to sign all commits
git config --global commit.gpgsign true

# Verify signatures in log
git log --show-signature
```

## Commit Guidelines

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(cli): add new command for ADR search

Closes #42

Signed-off-by: Your Name <your.email@example.com>
```

### Commit Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `test`: Adding/updating tests
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `chore`: Maintenance tasks
- `build`: Build system changes
- `ci`: CI configuration changes

### Scopes

Common scopes: `cli`, `adr`, `index`, `glob`, `docs`, `deps`

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

- **Subject**: Imperative mood, no period, max 72 characters
- **Body**: Explain *what* and *why*, not *how*
- **Footer**: Reference issues, include sign-off

## Pull Request Requirements

Each pull request must:

1. **Address a single feature or fix** - Keep PRs focused
2. **Be rebased on main** - Use `git rebase main`, not merge commits
3. **Contain only relevant changes** - No unrelated formatting or refactoring
4. **Include tests** - Cover new functionality with unit tests
5. **Update documentation** - Reflect changes in README or SPEC if needed
6. **Reference the issue** - Link to the GitHub issue being addressed
7. **Have signed commits** - All commits must have DCO sign-off and GPG signature
8. **Pass CI checks** - All tests and linters must pass

### Rebasing Your Branch

Keep your branch up to date with main:

```bash
git fetch origin
git rebase origin/main
git push --force-with-lease
```

### Squashing Commits

If you have multiple commits, squash them into a single descriptive commit:

```bash
git rebase -i origin/main
# Mark commits as 'squash' or 'fixup'
git push --force-with-lease
```

## Reviewing AI-Assisted Contributions

Some contributions may be created with AI assistance. When reviewing such contributions, apply the same standards as any other code:

- All tests must pass
- ADRs and SPEC must be followed
- Code must meet style and quality guidelines
- Human review and approval is required before merging

AI-assisted tooling can accelerate development, but it does not replace careful review. Contributors are responsible for the correctness and quality of their submissions regardless of how they were produced.

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Run `go vet ./...` to catch issues
- Keep functions small and focused
- Add comments for exported functions
- Prefer returning errors over panicking

## ADRs

For significant changes, create an ADR first:

```bash
./decider new "Your proposed change" --status proposed
```

Discuss the ADR in an issue before implementing.

## Testing

- Write unit tests for new functions
- Use table-driven tests when appropriate
- Place test fixtures in `testdata/` directories
- Aim for good coverage of edge cases

## Questions?

Open an issue for questions or discussions.
