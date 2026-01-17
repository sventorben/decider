package adr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Use PostgreSQL for Persistence", "use-postgresql-for-persistence"},
		{"ADR Format", "adr-format"},
		{"Hello World", "hello-world"},
		{"multiple   spaces", "multiple-spaces"},
		{"UPPERCASE", "uppercase"},
		{"with_underscores", "with-underscores"},
		{"special!@#chars", "specialchars"},
		{"  leading and trailing  ", "leading-and-trailing"},
		{"kebab-case-already", "kebab-case-already"},
		{"Numbers 123 Here", "numbers-123-here"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToKebabCase(tt.input)
			if got != tt.want {
				t.Errorf("ToKebabCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		number int
		title  string
		want   string
	}{
		{1, "Use PostgreSQL", "0001-use-postgresql.md"},
		{42, "ADR Format Specification", "0042-adr-format-specification.md"},
		{9999, "Test", "9999-test.md"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got := GenerateFilename(tt.number, tt.title)
			if got != tt.want {
				t.Errorf("GenerateFilename(%d, %q) = %q, want %q", tt.number, tt.title, got, tt.want)
			}
		})
	}
}

func TestValidateFilename(t *testing.T) {
	tests := []struct {
		filename string
		adrID    string
		wantErr  bool
	}{
		{"0001-test.md", "ADR-0001", false},
		{"0042-use-postgresql.md", "ADR-0042", false},
		{"0001-test.md", "ADR-0002", true},  // Number mismatch
		{"test.md", "ADR-0001", true},       // Wrong pattern
		{"01-test.md", "ADR-0001", true},    // Wrong number format
		{"0001-test.txt", "ADR-0001", true}, // Wrong extension
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			err := ValidateFilename(tt.filename, tt.adrID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilename(%q, %q) error = %v, wantErr %v", tt.filename, tt.adrID, err, tt.wantErr)
			}
		})
	}
}

func TestFindNextNumber(t *testing.T) {
	// Create temp directory
	dir := t.TempDir()

	// Test empty directory
	num, err := FindNextNumber(dir)
	if err != nil {
		t.Fatalf("FindNextNumber() error = %v", err)
	}
	if num != 1 {
		t.Errorf("FindNextNumber() empty dir = %d, want 1", num)
	}

	// Create some ADR files
	files := []string{
		"0001-first.md",
		"0002-second.md",
		"0005-fifth.md", // Gap
	}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	num, err = FindNextNumber(dir)
	if err != nil {
		t.Fatalf("FindNextNumber() error = %v", err)
	}
	if num != 6 {
		t.Errorf("FindNextNumber() = %d, want 6", num)
	}

	// Test non-existent directory
	num, err = FindNextNumber(filepath.Join(dir, "nonexistent"))
	if err != nil {
		t.Fatalf("FindNextNumber() error = %v", err)
	}
	if num != 1 {
		t.Errorf("FindNextNumber() nonexistent dir = %d, want 1", num)
	}
}

func TestListADRFiles(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	files := []string{
		"0003-third.md",
		"0001-first.md",
		"0002-second.md",
		"readme.md", // Should be ignored
		"templates", // Directory, should be ignored
	}
	for _, f := range files {
		if f == "templates" {
			if err := os.Mkdir(filepath.Join(dir, f), 0755); err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}
		} else {
			if err := os.WriteFile(filepath.Join(dir, f), []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}
	}

	result, err := ListADRFiles(dir)
	if err != nil {
		t.Fatalf("ListADRFiles() error = %v", err)
	}

	if len(result) != 3 {
		t.Errorf("ListADRFiles() returned %d files, want 3", len(result))
	}

	// Check sorting
	expected := []string{"0001-first.md", "0002-second.md", "0003-third.md"}
	for i, f := range expected {
		if result[i] != f {
			t.Errorf("ListADRFiles()[%d] = %q, want %q", i, result[i], f)
		}
	}
}
