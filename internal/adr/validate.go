package adr

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationSeverity indicates whether an issue is an error or warning.
type ValidationSeverity string

const (
	SeverityError   ValidationSeverity = "error"
	SeverityWarning ValidationSeverity = "warning"
)

// ValidationError represents a validation failure for an ADR.
type ValidationError struct {
	File     string
	Field    string
	Message  string
	Severity ValidationSeverity
	Code     string // Machine-readable error code
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.File, e.Field, e.Message)
}

// ValidationResult contains all validation errors for an ADR.
type ValidationResult struct {
	File     string
	Errors   []ValidationError
	Warnings []ValidationError
}

// IsValid returns true if there are no validation errors.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// HasWarnings returns true if there are validation warnings.
func (r *ValidationResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// IsValidStrict returns true if there are no errors or warnings.
func (r *ValidationResult) IsValidStrict() bool {
	return len(r.Errors) == 0 && len(r.Warnings) == 0
}

// RequiredSections lists the sections that must be present in an ADR body.
var RequiredSections = []string{
	"Context",
	"Decision",
	"Alternatives Considered",
	"Consequences",
}

var adrIDRegex = regexp.MustCompile(`^ADR-\d{4}$`)
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// Validate checks an ADR for required fields and returns validation errors.
func Validate(adr *ADR) *ValidationResult {
	result := &ValidationResult{File: adr.Filename}

	// Validate frontmatter fields
	if adr.Frontmatter.ADRID == "" {
		result.addError("adr_id", "required field is missing")
	} else if !adrIDRegex.MatchString(adr.Frontmatter.ADRID) {
		result.addError("adr_id", "must match pattern ADR-NNNN")
	}

	if adr.Frontmatter.Title == "" {
		result.addError("title", "required field is missing")
	}

	if adr.Frontmatter.Status == "" {
		result.addError("status", "required field is missing")
	} else if _, err := ParseStatus(string(adr.Frontmatter.Status)); err != nil {
		result.addError("status", err.Error())
	}

	if adr.Frontmatter.Date == "" {
		result.addError("date", "required field is missing")
	} else if !dateRegex.MatchString(adr.Frontmatter.Date) {
		result.addError("date", "must be in YYYY-MM-DD format")
	}

	// Validate filename matches ADR ID
	if adr.Frontmatter.ADRID != "" && adr.Filename != "" {
		if err := ValidateFilename(adr.Filename, adr.Frontmatter.ADRID); err != nil {
			result.addError("filename", err.Error())
		}
	}

	// Validate required sections
	for _, section := range RequiredSections {
		if !hasSection(adr.Body, section) {
			result.addError("body", fmt.Sprintf("missing required section: %s", section))
		}
	}

	// Validate rationale pattern (warnings)
	ValidateRationalePattern(adr, result)

	return result
}

func (r *ValidationResult) addError(field, message string) {
	r.Errors = append(r.Errors, ValidationError{
		File:     r.File,
		Field:    field,
		Message:  message,
		Severity: SeverityError,
		Code:     "validation_error",
	})
}

func (r *ValidationResult) addErrorWithCode(field, message, code string) {
	r.Errors = append(r.Errors, ValidationError{
		File:     r.File,
		Field:    field,
		Message:  message,
		Severity: SeverityError,
		Code:     code,
	})
}

func (r *ValidationResult) addWarning(field, message, code string) {
	r.Warnings = append(r.Warnings, ValidationError{
		File:     r.File,
		Field:    field,
		Message:  message,
		Severity: SeverityWarning,
		Code:     code,
	})
}

// hasSection checks if the body contains a markdown heading with the given title.
func hasSection(body string, sectionTitle string) bool {
	// Look for ## Section Title or # Section Title
	patterns := []string{
		fmt.Sprintf("## %s", sectionTitle),
		fmt.Sprintf("# %s", sectionTitle),
		fmt.Sprintf("##%s", sectionTitle), // No space
		fmt.Sprintf("#%s", sectionTitle),  // No space
	}

	bodyLower := strings.ToLower(body)

	for _, pattern := range patterns {
		patternLower := strings.ToLower(pattern)
		if strings.Contains(bodyLower, patternLower) {
			return true
		}
	}

	return false
}

// ValidateAll validates multiple ADRs and returns all results.
func ValidateAll(adrs []*ADR) []*ValidationResult {
	var results []*ValidationResult
	for _, adr := range adrs {
		result := Validate(adr)
		if !result.IsValid() || result.HasWarnings() {
			results = append(results, result)
		}
	}
	return results
}

// Rationale pattern markers
var (
	adoptedBecausePattern  = regexp.MustCompile(`(?i)\*\*adopted because:\*\*`)
	adoptedDespitePattern  = regexp.MustCompile(`(?i)\*\*adopted despite:\*\*`)
	rejectedBecausePattern = regexp.MustCompile(`(?i)\*\*rejected because:\*\*`)
	rejectedDespitePattern = regexp.MustCompile(`(?i)\*\*rejected despite:\*\*`)
	adoptedHeadingPattern  = regexp.MustCompile(`(?i)###\s+.+:\s*adopted`)
	rejectedHeadingPattern = regexp.MustCompile(`(?i)###\s+.+:\s*rejected`)
)

// ValidateRationalePattern checks if an ADR follows the mandatory rationale pattern.
// Returns warnings for missing rationale sections.
func ValidateRationalePattern(adr *ADR, result *ValidationResult) {
	body := adr.Body

	// Check for adopted option rationale
	hasAdoptedHeading := adoptedHeadingPattern.MatchString(body)
	hasAdoptedBecause := adoptedBecausePattern.MatchString(body)
	hasAdoptedDespite := adoptedDespitePattern.MatchString(body)

	// If there's an adopted heading, check for rationale sections
	if hasAdoptedHeading {
		if !hasAdoptedBecause {
			result.addWarning("rationale", "missing 'Adopted because:' section for adopted option", "missing_adopted_because")
		}
		if !hasAdoptedDespite {
			result.addWarning("rationale", "missing 'Adopted despite:' section for adopted option", "missing_adopted_despite")
		}
	} else if hasSection(body, "Decision") {
		// Decision section exists but no explicit adopted option
		if !hasAdoptedBecause {
			result.addWarning("rationale", "missing 'Adopted because:' section in Decision", "missing_adopted_because")
		}
		if !hasAdoptedDespite {
			result.addWarning("rationale", "missing 'Adopted despite:' section in Decision", "missing_adopted_despite")
		}
	}

	// Check for rejected alternatives rationale
	hasRejectedHeading := rejectedHeadingPattern.MatchString(body)
	hasRejectedBecause := rejectedBecausePattern.MatchString(body)
	hasRejectedDespite := rejectedDespitePattern.MatchString(body)

	// If there are rejected alternatives, check for rationale sections
	if hasRejectedHeading {
		if !hasRejectedBecause {
			result.addWarning("rationale", "missing 'Rejected because:' section for rejected alternative", "missing_rejected_because")
		}
		if !hasRejectedDespite {
			result.addWarning("rationale", "missing 'Rejected despite:' section for rejected alternative", "missing_rejected_despite")
		}
	}
	// Note: We only warn about missing alternative rationale when there's an explicit
	// "### Something: Rejected" heading. This avoids false positives for ADRs that
	// simply say "No alternatives considered" or list alternatives without rejection.
}
