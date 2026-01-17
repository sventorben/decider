// Package glob provides glob pattern matching for scope paths.
package glob

import (
	pathpkg "path"
	"path/filepath"
	"strings"
)

// MaxGlobDepth limits the number of ** segments to prevent pathological patterns.
const MaxGlobDepth = 10

// Match checks if a path matches a glob pattern.
// Supports standard glob patterns including ** for recursive matching.
// Returns false for patterns with too many ** segments (DoS protection).
func Match(pattern, path string) bool {
	// Normalize path separators (handle both OS-native and explicit backslashes)
	pattern = filepath.ToSlash(pattern)
	path = filepath.ToSlash(path)
	// Also replace any remaining backslashes (for cross-platform compatibility)
	pattern = strings.ReplaceAll(pattern, "\\", "/")
	path = strings.ReplaceAll(path, "\\", "/")

	// Count ** segments to prevent pathological patterns
	if strings.Count(pattern, "**") > MaxGlobDepth {
		return false
	}

	// Handle ** patterns specially
	if strings.Contains(pattern, "**") {
		return matchDoublestar(pattern, path)
	}

	// Use path.Match for simple patterns (uses / as separator consistently)
	matched, err := pathpkg.Match(pattern, path)
	if err != nil {
		return false
	}
	return matched
}

// matchDoublestar handles patterns with ** (recursive directory matching).
func matchDoublestar(pattern, path string) bool {
	// Split pattern and path into parts
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	return matchParts(patternParts, pathParts)
}

func matchParts(pattern, path []string) bool {
	pi, pathi := 0, 0

	for pi < len(pattern) && pathi < len(path) {
		if pattern[pi] == "**" {
			// ** matches zero or more directories
			if pi == len(pattern)-1 {
				// ** at end matches everything
				return true
			}

			// Try matching ** with zero or more path segments
			for skip := 0; skip <= len(path)-pathi; skip++ {
				if matchParts(pattern[pi+1:], path[pathi+skip:]) {
					return true
				}
			}
			return false
		}

		// Match single part with potential wildcards
		matched, err := pathpkg.Match(pattern[pi], path[pathi])
		if err != nil || !matched {
			return false
		}

		pi++
		pathi++
	}

	// Handle trailing ** or exact match
	for pi < len(pattern) {
		if pattern[pi] != "**" {
			return false
		}
		pi++
	}

	return pathi == len(path)
}

// MatchAny checks if a path matches any of the given patterns.
func MatchAny(patterns []string, path string) bool {
	for _, pattern := range patterns {
		if Match(pattern, path) {
			return true
		}
	}
	return false
}

// FilterPaths returns only the paths that match at least one pattern.
func FilterPaths(patterns []string, paths []string) []string {
	var result []string
	for _, path := range paths {
		if MatchAny(patterns, path) {
			result = append(result, path)
		}
	}
	return result
}

// FindMatchingPatterns returns all patterns that match the given path.
func FindMatchingPatterns(patterns []string, path string) []string {
	var result []string
	for _, pattern := range patterns {
		if Match(pattern, path) {
			result = append(result, pattern)
		}
	}
	return result
}
