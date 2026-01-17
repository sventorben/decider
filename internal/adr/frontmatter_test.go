package adr

import (
	"strings"
	"testing"
)

func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		wantFM        bool
		wantBodyStart string
		wantErr       bool
	}{
		{
			name: "valid frontmatter",
			content: `---
adr_id: ADR-0001
title: Test ADR
---

# Test

Body content`,
			wantFM:        true,
			wantBodyStart: "# Test",
			wantErr:       false,
		},
		{
			name:          "no frontmatter",
			content:       "# Just markdown\n\nNo frontmatter here.",
			wantFM:        false,
			wantBodyStart: "# Just markdown",
			wantErr:       false,
		},
		{
			name: "unclosed frontmatter",
			content: `---
adr_id: ADR-0001
title: Test`,
			wantFM:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body, err := ExtractFrontmatter(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			hasFM := fm != ""
			if hasFM != tt.wantFM {
				t.Errorf("ExtractFrontmatter() hasFM = %v, want %v", hasFM, tt.wantFM)
			}

			if tt.wantBodyStart != "" && !strings.HasPrefix(strings.TrimSpace(body), tt.wantBodyStart) {
				t.Errorf("ExtractFrontmatter() body starts with %q, want %q", strings.TrimSpace(body)[:20], tt.wantBodyStart)
			}
		})
	}
}

func TestParseFrontmatter(t *testing.T) {
	yaml := `adr_id: ADR-0001
title: Test Decision
status: adopted
date: "2026-01-16"
scope:
  paths:
    - "src/**"
    - "*.go"
tags:
  - foundation
  - testing
constraints:
  - Must be fast
invariants:
  - Always returns valid data
supersedes: []
superseded_by: []
related_adrs:
  - ADR-0002`

	fm, err := ParseFrontmatter(yaml)
	if err != nil {
		t.Fatalf("ParseFrontmatter() error = %v", err)
	}

	if fm.ADRID != "ADR-0001" {
		t.Errorf("ADRID = %q, want %q", fm.ADRID, "ADR-0001")
	}
	if fm.Title != "Test Decision" {
		t.Errorf("Title = %q, want %q", fm.Title, "Test Decision")
	}
	if fm.Status != StatusAdopted {
		t.Errorf("Status = %q, want %q", fm.Status, StatusAdopted)
	}
	if len(fm.Scope.Paths) != 2 {
		t.Errorf("Scope.Paths length = %d, want 2", len(fm.Scope.Paths))
	}
	if len(fm.Tags) != 2 {
		t.Errorf("Tags length = %d, want 2", len(fm.Tags))
	}
	if len(fm.Constraints) != 1 {
		t.Errorf("Constraints length = %d, want 1", len(fm.Constraints))
	}
}

func TestParseADR(t *testing.T) {
	content := `---
adr_id: ADR-0042
title: Use PostgreSQL
status: adopted
date: "2026-01-16"
scope:
  paths:
    - "db/**"
tags:
  - database
constraints: []
invariants: []
supersedes: []
superseded_by: []
related_adrs: []
---

# ADR-0042: Use PostgreSQL

## Context

We need a database.

## Decision

We use PostgreSQL.

## Alternatives Considered

MySQL, SQLite.

## Consequences

Positive: great features.`

	adr, err := ParseADR(content, "0042-use-postgresql.md", "/path/to/adr/0042-use-postgresql.md")
	if err != nil {
		t.Fatalf("ParseADR() error = %v", err)
	}

	if adr.Frontmatter.ADRID != "ADR-0042" {
		t.Errorf("ADRID = %q, want %q", adr.Frontmatter.ADRID, "ADR-0042")
	}
	if adr.Filename != "0042-use-postgresql.md" {
		t.Errorf("Filename = %q, want %q", adr.Filename, "0042-use-postgresql.md")
	}
	if !strings.Contains(adr.Body, "## Context") {
		t.Error("Body should contain '## Context'")
	}
}

func TestSerializeFrontmatter(t *testing.T) {
	fm := &Frontmatter{
		ADRID:  "ADR-0001",
		Title:  "Test",
		Status: StatusProposed,
		Date:   "2026-01-16",
		Scope:  Scope{Paths: []string{"src/**"}},
		Tags:   []string{"test"},
	}

	result, err := SerializeFrontmatter(fm)
	if err != nil {
		t.Fatalf("SerializeFrontmatter() error = %v", err)
	}

	if !strings.HasPrefix(result, "---\n") {
		t.Error("Result should start with '---\\n'")
	}
	if !strings.HasSuffix(result, "---\n") {
		t.Error("Result should end with '---\\n'")
	}
	if !strings.Contains(result, "adr_id: ADR-0001") {
		t.Error("Result should contain 'adr_id: ADR-0001'")
	}
}
