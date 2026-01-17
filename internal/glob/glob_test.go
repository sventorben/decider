package glob

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		// Simple patterns
		{"*.go", "main.go", true},
		{"*.go", "main.js", false},
		{"*.go", "src/main.go", false},

		// Directory patterns
		{"src/*", "src/file.go", true},
		{"src/*", "src/sub/file.go", false},

		// Recursive patterns
		{"src/**", "src/file.go", true},
		{"src/**", "src/sub/file.go", true},
		{"src/**", "src/sub/deep/file.go", true},
		{"src/**", "other/file.go", false},

		// Pattern with extension
		{"src/**/*.go", "src/main.go", true},
		{"src/**/*.go", "src/sub/main.go", true},
		{"src/**/*.go", "src/sub/main.js", false},
		{"**/*.go", "main.go", true},
		{"**/*.go", "src/main.go", true},
		{"**/*.go", "src/sub/main.go", true},

		// Complex patterns
		{"cmd/decider/**", "cmd/decider/main.go", true},
		{"cmd/decider/**", "cmd/other/main.go", false},
		{"internal/**", "internal/adr/adr.go", true},

		// Exact match
		{"go.mod", "go.mod", true},
		{"go.mod", "other/go.mod", false},

		// Windows-style paths (should work after normalization)
		{"src/**", "src\\file.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.path, func(t *testing.T) {
			got := Match(tt.pattern, tt.path)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.path, got, tt.want)
			}
		})
	}
}

func TestMatchAny(t *testing.T) {
	patterns := []string{"src/**/*.go", "internal/**", "*.md"}

	tests := []struct {
		path string
		want bool
	}{
		{"src/main.go", true},
		{"internal/pkg/file.go", true},
		{"README.md", true},
		{"docs/guide.md", false},
		{"other/file.js", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := MatchAny(patterns, tt.path)
			if got != tt.want {
				t.Errorf("MatchAny(%v, %q) = %v, want %v", patterns, tt.path, got, tt.want)
			}
		})
	}
}

func TestFilterPaths(t *testing.T) {
	patterns := []string{"src/**/*.go", "*.md"}
	paths := []string{
		"src/main.go",
		"src/pkg/util.go",
		"README.md",
		"docs/guide.md",
		"main.js",
	}

	result := FilterPaths(patterns, paths)

	expected := []string{"src/main.go", "src/pkg/util.go", "README.md"}
	if len(result) != len(expected) {
		t.Errorf("FilterPaths() returned %d paths, want %d", len(result), len(expected))
	}

	for i, p := range expected {
		if result[i] != p {
			t.Errorf("FilterPaths()[%d] = %q, want %q", i, result[i], p)
		}
	}
}

func TestFindMatchingPatterns(t *testing.T) {
	patterns := []string{"src/**", "**/*.go", "internal/**"}

	tests := []struct {
		path      string
		wantCount int
	}{
		{"src/main.go", 2},     // matches src/** and **/*.go
		{"internal/adr.go", 2}, // matches internal/** and **/*.go
		{"other/file.md", 0},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := FindMatchingPatterns(patterns, tt.path)
			if len(result) != tt.wantCount {
				t.Errorf("FindMatchingPatterns(%v, %q) = %v (len %d), want len %d",
					patterns, tt.path, result, len(result), tt.wantCount)
			}
		})
	}
}
