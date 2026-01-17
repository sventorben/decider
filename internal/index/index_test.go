package index

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sventorben/decider/internal/adr"
)

func TestGenerate(t *testing.T) {
	adrs := []*adr.ADR{
		{
			Filename: "0001-first.md",
			Frontmatter: adr.Frontmatter{
				ADRID:  "ADR-0001",
				Title:  "First Decision",
				Status: adr.StatusAdopted,
				Date:   "2026-01-15",
				Tags:   []string{"foundation"},
				Scope:  adr.Scope{Paths: []string{"src/**"}},
			},
		},
		{
			Filename: "0002-second.md",
			Frontmatter: adr.Frontmatter{
				ADRID:  "ADR-0002",
				Title:  "Second Decision",
				Status: adr.StatusProposed,
				Date:   "2026-01-16",
				Tags:   []string{"api"},
				Scope:  adr.Scope{Paths: []string{"api/**"}},
			},
		},
	}

	idx := Generate(adrs)

	if idx.ADRCount != 2 {
		t.Errorf("ADRCount = %d, want 2", idx.ADRCount)
	}

	if len(idx.ADRs) != 2 {
		t.Errorf("len(ADRs) = %d, want 2", len(idx.ADRs))
	}

	// Verify generated_at is recent
	ts, err := time.Parse(time.RFC3339, idx.GeneratedAt)
	if err != nil {
		t.Errorf("GeneratedAt is not valid RFC3339: %v", err)
	}
	if time.Since(ts) > time.Minute {
		t.Errorf("GeneratedAt is too old: %v", idx.GeneratedAt)
	}

	// Check first entry
	entry := idx.ADRs[0]
	if entry.ADRID != "ADR-0001" {
		t.Errorf("ADRs[0].ADRID = %q, want %q", entry.ADRID, "ADR-0001")
	}
	if entry.File != "0001-first.md" {
		t.Errorf("ADRs[0].File = %q, want %q", entry.File, "0001-first.md")
	}
	if len(entry.Tags) != 1 || entry.Tags[0] != "foundation" {
		t.Errorf("ADRs[0].Tags = %v, want [foundation]", entry.Tags)
	}
}

func TestWriteAndLoad(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, IndexFilename)

	// Create index
	idx := &Index{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		ADRCount:    1,
		ADRs: []Entry{
			{
				ADRID:      "ADR-0001",
				Title:      "Test",
				Status:     "adopted",
				Date:       "2026-01-16",
				Tags:       []string{"test"},
				ScopePaths: []string{"src/**"},
				File:       "0001-test.md",
			},
		},
	}

	// Write
	if err := idx.Write(indexPath); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify file exists and has header
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index file: %v", err)
	}
	if !strings.HasPrefix(string(content), IndexHeader) {
		t.Error("Index file should start with header comment")
	}

	// Load
	loaded, err := Load(indexPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.ADRCount != idx.ADRCount {
		t.Errorf("Loaded ADRCount = %d, want %d", loaded.ADRCount, idx.ADRCount)
	}
	if len(loaded.ADRs) != 1 {
		t.Errorf("Loaded ADRs length = %d, want 1", len(loaded.ADRs))
	}
	if loaded.ADRs[0].ADRID != "ADR-0001" {
		t.Errorf("Loaded ADRs[0].ADRID = %q, want %q", loaded.ADRs[0].ADRID, "ADR-0001")
	}
}

func TestExists(t *testing.T) {
	dir := t.TempDir()

	// Should not exist initially
	if Exists(dir) {
		t.Error("Exists() should return false for empty directory")
	}

	// Create index file
	indexPath := filepath.Join(dir, IndexFilename)
	if err := os.WriteFile(indexPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Should exist now
	if !Exists(dir) {
		t.Error("Exists() should return true after creating index file")
	}
}

func TestEntriesEqual(t *testing.T) {
	a := Entry{
		ADRID:      "ADR-0001",
		Title:      "Test",
		Status:     "adopted",
		Date:       "2026-01-16",
		Tags:       []string{"a", "b"},
		ScopePaths: []string{"src/**"},
		File:       "0001-test.md",
	}

	b := Entry{
		ADRID:      "ADR-0001",
		Title:      "Test",
		Status:     "adopted",
		Date:       "2026-01-16",
		Tags:       []string{"a", "b"},
		ScopePaths: []string{"src/**"},
		File:       "0001-test.md",
	}

	if !entriesEqual(a, b) {
		t.Error("entriesEqual() should return true for identical entries")
	}

	// Different title
	b.Title = "Different"
	if entriesEqual(a, b) {
		t.Error("entriesEqual() should return false for different titles")
	}
	b.Title = "Test"

	// Different tags
	b.Tags = []string{"a", "c"}
	if entriesEqual(a, b) {
		t.Error("entriesEqual() should return false for different tags")
	}
	b.Tags = []string{"a", "b"}

	// Different number of tags
	b.Tags = []string{"a"}
	if entriesEqual(a, b) {
		t.Error("entriesEqual() should return false for different tag counts")
	}
}
