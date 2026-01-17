package cli

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sventorben/decider/internal/adr"
	"github.com/sventorben/decider/internal/glob"
	"github.com/sventorben/decider/internal/validate"
)

// CheckADRConfig holds configuration for the check adr command.
type CheckADRConfig struct {
	Dir    string
	Strict bool
	Format OutputFormat
	Output *Output
}

// CheckADRResult holds the result of the check adr command.
type CheckADRResult struct {
	Valid    bool                 `json:"valid"`
	Count    int                  `json:"count"`
	Errors   []CheckADRError      `json:"errors,omitempty"`
	Warnings []CheckADRError      `json:"warnings,omitempty"`
	Results  []CheckADRFileResult `json:"results,omitempty"`
}

// CheckADRError represents an error or warning found during ADR validation.
type CheckADRError struct {
	File     string `json:"file"`
	Field    string `json:"field"`
	Message  string `json:"message"`
	Severity string `json:"severity"` // "error" or "warning"
	Code     string `json:"code"`     // Machine-readable error code
}

// CheckADRFileResult represents the validation result for a single file.
type CheckADRFileResult struct {
	File     string          `json:"file"`
	Valid    bool            `json:"valid"`
	Errors   []CheckADRError `json:"errors,omitempty"`
	Warnings []CheckADRError `json:"warnings,omitempty"`
}

// RunCheckADR validates all ADRs in the directory.
func RunCheckADR(cfg *CheckADRConfig) (*CheckADRResult, error) {
	adrs, err := adr.LoadAllADRs(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("loading ADRs: %w", err)
	}

	result := &CheckADRResult{
		Valid: true,
		Count: len(adrs),
	}

	for _, a := range adrs {
		vr := adr.Validate(a)
		fileResult := CheckADRFileResult{
			File:  a.Filename,
			Valid: vr.IsValid(),
		}

		// Process errors
		if !vr.IsValid() {
			result.Valid = false
			for _, ve := range vr.Errors {
				checkErr := CheckADRError{
					File:     ve.File,
					Field:    ve.Field,
					Message:  ve.Message,
					Severity: string(ve.Severity),
					Code:     ve.Code,
				}
				fileResult.Errors = append(fileResult.Errors, checkErr)
				result.Errors = append(result.Errors, checkErr)
			}
		}

		// Process warnings
		for _, vw := range vr.Warnings {
			checkWarn := CheckADRError{
				File:     vw.File,
				Field:    vw.Field,
				Message:  vw.Message,
				Severity: string(vw.Severity),
				Code:     vw.Code,
			}
			fileResult.Warnings = append(fileResult.Warnings, checkWarn)
			result.Warnings = append(result.Warnings, checkWarn)
		}

		// In strict mode, warnings invalidate the file
		if cfg.Strict && len(fileResult.Warnings) > 0 {
			fileResult.Valid = false
			result.Valid = false
		}

		result.Results = append(result.Results, fileResult)
	}

	// Output
	if cfg.Format == FormatJSON {
		_ = cfg.Output.PrintJSON(result)
	} else {
		hasErrors := len(result.Errors) > 0
		hasWarnings := len(result.Warnings) > 0

		if !hasErrors && !hasWarnings {
			cfg.Output.Success("All %d ADR(s) are valid", result.Count)
		} else {
			if hasErrors {
				cfg.Output.Error("Found %d validation error(s):", len(result.Errors))
				for _, e := range result.Errors {
					cfg.Output.Println("  [error] %s: %s: %s", e.File, e.Field, e.Message)
				}
			}
			if hasWarnings {
				if cfg.Strict {
					cfg.Output.Error("Found %d warning(s) (strict mode):", len(result.Warnings))
				} else {
					cfg.Output.Warn("Found %d warning(s):", len(result.Warnings))
				}
				for _, w := range result.Warnings {
					cfg.Output.Println("  [warning] %s: %s: %s", w.File, w.Field, w.Message)
				}
			}
			if !hasErrors && hasWarnings && !cfg.Strict {
				cfg.Output.Success("All %d ADR(s) are valid (with warnings)", result.Count)
			}
		}
	}

	return result, nil
}

// CheckDiffConfig holds configuration for the check diff command.
type CheckDiffConfig struct {
	Dir    string
	Base   string
	Format OutputFormat
	Output *Output
}

// CheckDiffResult holds the result of the check diff command.
type CheckDiffResult struct {
	ChangedFiles   []string           `json:"changed_files"`
	ApplicableADRs []ApplicableADR    `json:"applicable_adrs"`
	Summary        ConstraintsSummary `json:"summary"`
}

// ApplicableADR represents an ADR that applies to changed files.
type ApplicableADR struct {
	ADRID        string   `json:"adr_id"`
	Title        string   `json:"title"`
	MatchedPaths []string `json:"matched_paths"`
	MatchedFiles []string `json:"matched_files"`
	Constraints  []string `json:"constraints,omitempty"`
	Invariants   []string `json:"invariants,omitempty"`
}

// ConstraintsSummary provides a summary of all constraints that apply.
type ConstraintsSummary struct {
	TotalADRs        int      `json:"total_adrs"`
	TotalConstraints int      `json:"total_constraints"`
	TotalInvariants  int      `json:"total_invariants"`
	AllConstraints   []string `json:"all_constraints,omitempty"`
	AllInvariants    []string `json:"all_invariants,omitempty"`
}

// RunCheckDiff finds ADRs applicable to changed files.
func RunCheckDiff(cfg *CheckDiffConfig) (*CheckDiffResult, error) {
	// Get changed files from git
	changedFiles, err := getGitDiff(cfg.Base)
	if err != nil {
		return nil, fmt.Errorf("getting git diff: %w", err)
	}

	// Load all ADRs
	adrs, err := adr.LoadAllADRs(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("loading ADRs: %w", err)
	}

	result := &CheckDiffResult{
		ChangedFiles: changedFiles,
	}

	// Find applicable ADRs
	for _, a := range adrs {
		if len(a.Frontmatter.Scope.Paths) == 0 {
			continue
		}

		var matchedFiles []string
		var matchedPaths []string

		for _, changedFile := range changedFiles {
			for _, scopePath := range a.Frontmatter.Scope.Paths {
				if glob.Match(scopePath, changedFile) {
					matchedFiles = append(matchedFiles, changedFile)
					if !contains(matchedPaths, scopePath) {
						matchedPaths = append(matchedPaths, scopePath)
					}
					break
				}
			}
		}

		if len(matchedFiles) > 0 {
			result.ApplicableADRs = append(result.ApplicableADRs, ApplicableADR{
				ADRID:        a.Frontmatter.ADRID,
				Title:        a.Frontmatter.Title,
				MatchedPaths: matchedPaths,
				MatchedFiles: matchedFiles,
				Constraints:  a.Frontmatter.Constraints,
				Invariants:   a.Frontmatter.Invariants,
			})
		}
	}

	// Build summary
	result.Summary.TotalADRs = len(result.ApplicableADRs)
	for _, aa := range result.ApplicableADRs {
		result.Summary.TotalConstraints += len(aa.Constraints)
		result.Summary.TotalInvariants += len(aa.Invariants)
		result.Summary.AllConstraints = append(result.Summary.AllConstraints, aa.Constraints...)
		result.Summary.AllInvariants = append(result.Summary.AllInvariants, aa.Invariants...)
	}

	// Output
	if cfg.Format == FormatJSON {
		_ = cfg.Output.PrintJSON(result)
	} else {
		cfg.Output.Println("Changed files: %d", len(changedFiles))
		cfg.Output.Println("Applicable ADRs: %d", len(result.ApplicableADRs))
		cfg.Output.Println("")

		if len(result.ApplicableADRs) == 0 {
			cfg.Output.Println("No ADRs apply to the changed files.")
		} else {
			cfg.Output.Println("## Constraint Summary")
			cfg.Output.Println("")

			for _, aa := range result.ApplicableADRs {
				cfg.Output.Println("### %s: %s", aa.ADRID, aa.Title)
				cfg.Output.Println("Matches: %s", strings.Join(aa.MatchedPaths, ", "))

				if len(aa.Constraints) > 0 {
					cfg.Output.Println("Constraints:")
					for _, c := range aa.Constraints {
						cfg.Output.Println("  - %s", c)
					}
				}

				if len(aa.Invariants) > 0 {
					cfg.Output.Println("Invariants:")
					for _, i := range aa.Invariants {
						cfg.Output.Println("  - %s", i)
					}
				}
				cfg.Output.Println("")
			}
		}
	}

	return result, nil
}

func getGitDiff(base string) ([]string, error) {
	// Validate git ref to prevent injection
	if err := validate.ValidateGitRef(base); err != nil {
		return nil, fmt.Errorf("invalid git ref: %w", err)
	}

	cmd := exec.Command("git", "diff", "--name-only", base)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff failed: %w", err)
	}

	var files []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, filepath.ToSlash(line))
		}
	}

	return files, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
