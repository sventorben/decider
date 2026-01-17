package adr

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		adr       *ADR
		wantValid bool
		wantErrs  int
	}{
		{
			name: "valid ADR",
			adr: &ADR{
				Filename: "0001-test-decision.md",
				Frontmatter: Frontmatter{
					ADRID:  "ADR-0001",
					Title:  "Test Decision",
					Status: StatusAdopted,
					Date:   "2026-01-16",
				},
				Body: `# ADR-0001: Test

## Context

Some context.

## Decision

We decided.

## Alternatives Considered

Other options.

## Consequences

What happens.`,
			},
			wantValid: true,
			wantErrs:  0,
		},
		{
			name: "missing required fields",
			adr: &ADR{
				Filename:    "0001-test.md",
				Frontmatter: Frontmatter{},
				Body:        "# Test\n## Context\n## Decision\n## Alternatives Considered\n## Consequences",
			},
			wantValid: false,
			wantErrs:  4, // adr_id, title, status, date
		},
		{
			name: "invalid status",
			adr: &ADR{
				Filename: "0001-test.md",
				Frontmatter: Frontmatter{
					ADRID:  "ADR-0001",
					Title:  "Test",
					Status: "invalid",
					Date:   "2026-01-16",
				},
				Body: "## Context\n## Decision\n## Alternatives Considered\n## Consequences",
			},
			wantValid: false,
			wantErrs:  1,
		},
		{
			name: "invalid date format",
			adr: &ADR{
				Filename: "0001-test.md",
				Frontmatter: Frontmatter{
					ADRID:  "ADR-0001",
					Title:  "Test",
					Status: StatusAdopted,
					Date:   "January 16, 2026",
				},
				Body: "## Context\n## Decision\n## Alternatives Considered\n## Consequences",
			},
			wantValid: false,
			wantErrs:  1,
		},
		{
			name: "filename/ID mismatch",
			adr: &ADR{
				Filename: "0002-test.md",
				Frontmatter: Frontmatter{
					ADRID:  "ADR-0001",
					Title:  "Test",
					Status: StatusAdopted,
					Date:   "2026-01-16",
				},
				Body: "## Context\n## Decision\n## Alternatives Considered\n## Consequences",
			},
			wantValid: false,
			wantErrs:  1,
		},
		{
			name: "missing sections",
			adr: &ADR{
				Filename: "0001-test.md",
				Frontmatter: Frontmatter{
					ADRID:  "ADR-0001",
					Title:  "Test",
					Status: StatusAdopted,
					Date:   "2026-01-16",
				},
				Body: "# Test\n\nJust some content without sections.",
			},
			wantValid: false,
			wantErrs:  4, // Missing all 4 required sections
		},
		{
			name: "invalid ADR ID format",
			adr: &ADR{
				Filename: "0001-test.md",
				Frontmatter: Frontmatter{
					ADRID:  "0001", // Missing ADR- prefix
					Title:  "Test",
					Status: StatusAdopted,
					Date:   "2026-01-16",
				},
				Body: "## Context\n## Decision\n## Alternatives Considered\n## Consequences",
			},
			wantValid: false,
			wantErrs:  1, // Invalid format
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.adr)
			if result.IsValid() != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v", result.IsValid(), tt.wantValid)
			}
			if len(result.Errors) != tt.wantErrs {
				t.Errorf("Validate() errors = %d, want %d", len(result.Errors), tt.wantErrs)
				for _, e := range result.Errors {
					t.Logf("  Error: %s: %s", e.Field, e.Message)
				}
			}
		})
	}
}

func TestHasSection(t *testing.T) {
	tests := []struct {
		body    string
		section string
		want    bool
	}{
		{"## Context\nSome text", "Context", true},
		{"# Context\nSome text", "Context", true},
		{"##Context\nSome text", "Context", true},
		{"## CONTEXT\nSome text", "Context", true},
		{"## context\nSome text", "Context", true},
		{"Some text without section", "Context", false},
		{"## Alternatives Considered\nOptions", "Alternatives Considered", true},
	}

	for _, tt := range tests {
		t.Run(tt.section, func(t *testing.T) {
			got := hasSection(tt.body, tt.section)
			if got != tt.want {
				t.Errorf("hasSection(%q, %q) = %v, want %v", tt.body[:min(20, len(tt.body))], tt.section, got, tt.want)
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestValidateRationalePattern(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantWarnings int
	}{
		{
			name: "compliant ADR with full rationale pattern",
			body: `# ADR-0001: Test

## Context

Decision drivers:
- Driver 1

## Decision

### Option A: Adopted

**Adopted because:**
- Reason 1

**Adopted despite:**
- Trade-off 1

## Alternatives Considered

### Option B: Rejected

**Rejected because:**
- Reason 1

**Rejected despite:**
- Strength 1

## Consequences

Positive and negative.`,
			wantWarnings: 0,
		},
		{
			name: "missing adopted because",
			body: `# ADR-0001: Test

## Context

Context here.

## Decision

### Option A: Adopted

**Adopted despite:**
- Trade-off 1

## Alternatives Considered

None.

## Consequences

Done.`,
			wantWarnings: 1, // missing "Adopted because:"
		},
		{
			name: "missing adopted despite",
			body: `# ADR-0001: Test

## Context

Context here.

## Decision

### Option A: Adopted

**Adopted because:**
- Reason 1

## Alternatives Considered

None.

## Consequences

Done.`,
			wantWarnings: 1, // missing "Adopted despite:"
		},
		{
			name: "missing both adopted rationale sections",
			body: `# ADR-0001: Test

## Context

Context here.

## Decision

### Option A: Adopted

Just some prose.

## Alternatives Considered

None.

## Consequences

Done.`,
			wantWarnings: 2, // missing both
		},
		{
			name: "missing rejected rationale",
			body: `# ADR-0001: Test

## Context

Context here.

## Decision

### Option A: Adopted

**Adopted because:**
- Reason 1

**Adopted despite:**
- Trade-off 1

## Alternatives Considered

### Option B: Rejected

Just some prose about rejection.

## Consequences

Done.`,
			wantWarnings: 2, // missing "Rejected because:" and "Rejected despite:"
		},
		{
			name: "no explicit adopted heading but has rationale",
			body: `# ADR-0001: Test

## Context

Context here.

## Decision

We chose to do X.

**Adopted because:**
- Reason 1

**Adopted despite:**
- Trade-off 1

## Alternatives Considered

None considered.

## Consequences

Done.`,
			wantWarnings: 0, // Has rationale, no explicit heading but that's ok
		},
		{
			name: "old-style pros/cons table",
			body: `# ADR-0001: Test

## Context

Context here.

## Decision

We chose option A.

## Alternatives Considered

| Alternative | Pros | Cons |
|-------------|------|------|
| Option A | Good | Bad |

## Consequences

Done.`,
			wantWarnings: 2, // missing adopted because/despite (no rejected heading = no rejected warning)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adr := &ADR{
				Filename: "0001-test.md",
				Frontmatter: Frontmatter{
					ADRID:  "ADR-0001",
					Title:  "Test",
					Status: StatusAdopted,
					Date:   "2026-01-16",
				},
				Body: tt.body,
			}

			result := Validate(adr)
			if len(result.Warnings) != tt.wantWarnings {
				t.Errorf("ValidateRationalePattern() warnings = %d, want %d", len(result.Warnings), tt.wantWarnings)
				for _, w := range result.Warnings {
					t.Logf("  Warning: %s: %s (code: %s)", w.Field, w.Message, w.Code)
				}
			}
		})
	}
}

func TestValidationSeverity(t *testing.T) {
	// Test that errors have error severity and warnings have warning severity
	adr := &ADR{
		Filename:    "0001-test.md",
		Frontmatter: Frontmatter{}, // Missing required fields
		Body: `## Context
## Decision
## Alternatives Considered
## Consequences`,
	}

	result := Validate(adr)

	// Should have errors for missing required fields
	if len(result.Errors) == 0 {
		t.Error("Expected errors for missing required fields")
	}

	for _, e := range result.Errors {
		if e.Severity != SeverityError {
			t.Errorf("Error severity = %q, want %q", e.Severity, SeverityError)
		}
	}

	// Should have warnings for missing rationale pattern
	if len(result.Warnings) == 0 {
		t.Error("Expected warnings for missing rationale pattern")
	}

	for _, w := range result.Warnings {
		if w.Severity != SeverityWarning {
			t.Errorf("Warning severity = %q, want %q", w.Severity, SeverityWarning)
		}
	}
}

func TestIsValidStrict(t *testing.T) {
	tests := []struct {
		name            string
		adr             *ADR
		wantValid       bool
		wantValidStrict bool
	}{
		{
			name: "valid with no warnings",
			adr: &ADR{
				Filename: "0001-test.md",
				Frontmatter: Frontmatter{
					ADRID:  "ADR-0001",
					Title:  "Test",
					Status: StatusAdopted,
					Date:   "2026-01-16",
				},
				Body: `## Context

Some context.

## Decision

### Option: Adopted

**Adopted because:**
- Reason

**Adopted despite:**
- Trade-off

## Alternatives Considered

None.

## Consequences

Some consequences.`,
			},
			wantValid:       true,
			wantValidStrict: true,
		},
		{
			name: "valid with warnings",
			adr: &ADR{
				Filename: "0001-test.md",
				Frontmatter: Frontmatter{
					ADRID:  "ADR-0001",
					Title:  "Test",
					Status: StatusAdopted,
					Date:   "2026-01-16",
				},
				Body: `## Context
## Decision
Just some decision.
## Alternatives Considered
## Consequences`,
			},
			wantValid:       true,
			wantValidStrict: false, // Has warnings
		},
		{
			name: "invalid with errors",
			adr: &ADR{
				Filename:    "0001-test.md",
				Frontmatter: Frontmatter{},
				Body:        "# Test",
			},
			wantValid:       false,
			wantValidStrict: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.adr)
			if result.IsValid() != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v", result.IsValid(), tt.wantValid)
			}
			if result.IsValidStrict() != tt.wantValidStrict {
				t.Errorf("IsValidStrict() = %v, want %v", result.IsValidStrict(), tt.wantValidStrict)
			}
		})
	}
}
