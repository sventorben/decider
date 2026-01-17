// Package validate provides input validation utilities for the DECIDER CLI.
package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Limits defines maximum lengths for various inputs.
const (
	MaxTitleLength      = 500
	MaxTagLength        = 100
	MaxPathLength       = 1000
	MaxConstraintLength = 2000
	MaxTagCount         = 50
	MaxPathCount        = 100
	MaxFileSizeBytes    = 10 * 1024 * 1024 // 10 MB
	MaxGlobDepth        = 10               // Maximum ** segments in a glob pattern
)

// gitRefPattern matches valid git references (branches, tags, SHAs).
// Allows: alphanumeric, hyphens, underscores, slashes, dots, tildes (~), carets (^), and SHA formats.
// Tildes and carets are used for ancestor/parent references (e.g., HEAD~1, HEAD^2).
var gitRefPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._/~^-]*$|^[a-fA-F0-9]{7,40}$`)

// ValidateGitRef checks if a string is a valid git reference.
func ValidateGitRef(ref string) error {
	if ref == "" {
		return fmt.Errorf("git ref cannot be empty")
	}
	if len(ref) > 256 {
		return fmt.Errorf("git ref too long (max 256 characters)")
	}
	// Reject refs that start with a hyphen (could be interpreted as flags)
	if strings.HasPrefix(ref, "-") {
		return fmt.Errorf("git ref cannot start with hyphen")
	}
	// Reject refs with suspicious patterns
	if strings.Contains(ref, "..") && strings.Contains(ref, "/") {
		// Allow ".." for git range syntax but not path traversal
		parts := strings.Split(ref, "..")
		for _, part := range parts {
			if strings.Contains(part, "../") || strings.HasPrefix(part, "/") {
				return fmt.Errorf("git ref contains invalid path pattern")
			}
		}
	}
	if !gitRefPattern.MatchString(strings.ReplaceAll(ref, "..", "")) {
		return fmt.Errorf("git ref contains invalid characters")
	}
	return nil
}

// ValidateTitle checks if an ADR title is valid.
func ValidateTitle(title string) error {
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if len(title) > MaxTitleLength {
		return fmt.Errorf("title too long (max %d characters)", MaxTitleLength)
	}
	return nil
}

// ValidateTag checks if a tag is valid.
func ValidateTag(tag string) error {
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if len(tag) > MaxTagLength {
		return fmt.Errorf("tag too long (max %d characters)", MaxTagLength)
	}
	// Tags should be simple identifiers
	for _, r := range tag {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_') {
			return fmt.Errorf("tag contains invalid character: %c", r)
		}
	}
	return nil
}

// ValidateTags checks a list of tags.
func ValidateTags(tags []string) error {
	if len(tags) > MaxTagCount {
		return fmt.Errorf("too many tags (max %d)", MaxTagCount)
	}
	for _, tag := range tags {
		if err := ValidateTag(tag); err != nil {
			return err
		}
	}
	return nil
}

// ValidateScopePath checks if a scope path pattern is valid.
func ValidateScopePath(path string) error {
	if path == "" {
		return fmt.Errorf("scope path cannot be empty")
	}
	if len(path) > MaxPathLength {
		return fmt.Errorf("scope path too long (max %d characters)", MaxPathLength)
	}
	// Count ** segments to prevent pathological patterns
	segments := strings.Count(path, "**")
	if segments > MaxGlobDepth {
		return fmt.Errorf("scope path has too many ** segments (max %d)", MaxGlobDepth)
	}
	return nil
}

// ValidateScopePaths checks a list of scope paths.
func ValidateScopePaths(paths []string) error {
	if len(paths) > MaxPathCount {
		return fmt.Errorf("too many scope paths (max %d)", MaxPathCount)
	}
	for _, path := range paths {
		if err := ValidateScopePath(path); err != nil {
			return err
		}
	}
	return nil
}

// SanitizePath cleans and validates a directory path.
// Returns the cleaned path or an error if the path is invalid.
func SanitizePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Clean the path
	cleaned := filepath.Clean(path)

	// Convert to absolute for validation if it looks relative
	var absPath string
	if filepath.IsAbs(cleaned) {
		absPath = cleaned
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("cannot get working directory: %w", err)
		}
		absPath = filepath.Join(wd, cleaned)
	}

	// Evaluate any symlinks to get real path
	realPath, err := filepath.EvalSymlinks(filepath.Dir(absPath))
	if err != nil && !os.IsNotExist(err) {
		// Parent must exist for new directories, but we allow new leaf directories
		parentDir := filepath.Dir(absPath)
		if _, statErr := os.Stat(parentDir); statErr != nil && !os.IsNotExist(statErr) {
			return "", fmt.Errorf("invalid path: %w", err)
		}
	}

	// If we could resolve the real path, use it
	if realPath != "" {
		absPath = filepath.Join(realPath, filepath.Base(absPath))
	}

	// Return the original (cleaned) path to preserve relative paths
	return cleaned, nil
}

// CheckFileSize verifies a file is within size limits.
func CheckFileSize(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.Size() > MaxFileSizeBytes {
		return fmt.Errorf("file too large: %d bytes (max %d)", info.Size(), MaxFileSizeBytes)
	}
	return nil
}
