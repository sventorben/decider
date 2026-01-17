package validate

import (
	"strings"
	"testing"
)

func TestValidateGitRef(t *testing.T) {
	tests := []struct {
		name    string
		ref     string
		wantErr bool
	}{
		{"valid branch", "main", false},
		{"valid branch with slash", "feature/my-feature", false},
		{"valid sha short", "abc1234", false},
		{"valid sha full", "abc1234567890def1234567890abc1234567890ab", false},
		{"valid tag", "v1.0.0", false},
		{"valid range", "main..feature", false},
		{"empty", "", true},
		{"starts with hyphen", "-flag", true},
		{"too long", strings.Repeat("a", 300), true},
		{"origin/main", "origin/main", false},
		{"HEAD~1", "HEAD~1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGitRef(tt.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGitRef(%q) error = %v, wantErr %v", tt.ref, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"valid title", "Use PostgreSQL for Persistence", false},
		{"empty", "", true},
		{"too long", strings.Repeat("a", MaxTitleLength+1), true},
		{"max length", strings.Repeat("a", MaxTitleLength), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTitle(%q) error = %v, wantErr %v", tt.title, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		wantErr bool
	}{
		{"valid lowercase", "database", false},
		{"valid with hyphen", "my-tag", false},
		{"valid with underscore", "my_tag", false},
		{"valid mixed case", "MyTag", false},
		{"empty", "", true},
		{"too long", strings.Repeat("a", MaxTagLength+1), true},
		{"invalid char space", "my tag", true},
		{"invalid char special", "my@tag", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTag(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTag(%q) error = %v, wantErr %v", tt.tag, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTags(t *testing.T) {
	tests := []struct {
		name    string
		tags    []string
		wantErr bool
	}{
		{"valid tags", []string{"tag1", "tag2"}, false},
		{"empty list", []string{}, false},
		{"too many tags", make([]string, MaxTagCount+1), true},
		{"one invalid tag", []string{"valid", "invalid@tag"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill placeholder tags for "too many" test
			if len(tt.tags) == MaxTagCount+1 {
				for i := range tt.tags {
					tt.tags[i] = "tag"
				}
			}
			err := ValidateTags(tt.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateScopePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid simple", "src/**", false},
		{"valid glob", "internal/**/*.go", false},
		{"empty", "", true},
		{"too long", strings.Repeat("a", MaxPathLength+1), true},
		{"too many doublestar", "**/**/**/**/**/**/**/**/**/**/**/*.go", true},
		{"valid doublestar count", "src/**/pkg/**/*.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateScopePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateScopePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid relative", "docs/adr", false},
		{"valid current dir", ".", false},
		{"empty", "", true},
		{"double dot cleaned", "docs/../docs/adr", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SanitizePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}
