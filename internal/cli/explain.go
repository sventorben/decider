package cli

import (
	"fmt"
	"strings"

	"github.com/sventorben/decider/internal/adr"
	"github.com/sventorben/decider/internal/glob"
)

// ExplainConfig holds configuration for the explain command.
type ExplainConfig struct {
	Dir    string
	Base   string
	Format OutputFormat
	Output *Output
}

// ExplainResult holds the result of the explain command.
type ExplainResult struct {
	ChangedFiles []string       `json:"changed_files"`
	Explanations []ExplainEntry `json:"explanations"`
}

// ExplainEntry explains why an ADR applies to specific files.
type ExplainEntry struct {
	ADRID       string         `json:"adr_id"`
	Title       string         `json:"title"`
	Status      string         `json:"status"`
	Matches     []MatchExplain `json:"matches"`
	Constraints []string       `json:"constraints,omitempty"`
	Invariants  []string       `json:"invariants,omitempty"`
}

// MatchExplain explains a single path match.
type MatchExplain struct {
	File    string `json:"file"`
	Pattern string `json:"pattern"`
	Reason  string `json:"reason"`
}

// RunExplain provides narrative explanation of why ADRs apply.
func RunExplain(cfg *ExplainConfig) (*ExplainResult, error) {
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

	result := &ExplainResult{
		ChangedFiles: changedFiles,
	}

	// Find applicable ADRs with explanations
	for _, a := range adrs {
		if len(a.Frontmatter.Scope.Paths) == 0 {
			continue
		}

		var matches []MatchExplain

		for _, changedFile := range changedFiles {
			for _, scopePath := range a.Frontmatter.Scope.Paths {
				if glob.Match(scopePath, changedFile) {
					matches = append(matches, MatchExplain{
						File:    changedFile,
						Pattern: scopePath,
						Reason:  explainMatch(scopePath, changedFile),
					})
					break // Only match first pattern per file
				}
			}
		}

		if len(matches) > 0 {
			result.Explanations = append(result.Explanations, ExplainEntry{
				ADRID:       a.Frontmatter.ADRID,
				Title:       a.Frontmatter.Title,
				Status:      string(a.Frontmatter.Status),
				Matches:     matches,
				Constraints: a.Frontmatter.Constraints,
				Invariants:  a.Frontmatter.Invariants,
			})
		}
	}

	// Output
	if cfg.Format == FormatJSON {
		_ = cfg.Output.PrintJSON(result)
	} else {
		cfg.Output.Println("# ADR Applicability Analysis")
		cfg.Output.Println("")
		cfg.Output.Println("Analyzing %d changed file(s) against %d ADR(s)...", len(changedFiles), len(adrs))
		cfg.Output.Println("")

		if len(result.Explanations) == 0 {
			cfg.Output.Println("No ADRs apply to the changed files.")
			cfg.Output.Println("")
			cfg.Output.Println("This means the changes don't fall within the scope of any documented")
			cfg.Output.Println("architectural decisions. This is fine for routine changes, but consider")
			cfg.Output.Println("whether a new ADR might be needed for significant changes.")
		} else {
			cfg.Output.Println("Found %d applicable ADR(s):", len(result.Explanations))
			cfg.Output.Println("")

			for _, exp := range result.Explanations {
				cfg.Output.Println("## %s: %s", exp.ADRID, exp.Title)
				cfg.Output.Println("Status: %s", exp.Status)
				cfg.Output.Println("")
				cfg.Output.Println("### Why This ADR Applies")
				cfg.Output.Println("")

				for _, m := range exp.Matches {
					cfg.Output.Println("- **%s**", m.File)
					cfg.Output.Println("  Pattern: `%s`", m.Pattern)
					cfg.Output.Println("  %s", m.Reason)
				}

				if len(exp.Constraints) > 0 {
					cfg.Output.Println("")
					cfg.Output.Println("### Constraints to Follow")
					for _, c := range exp.Constraints {
						cfg.Output.Println("- %s", c)
					}
				}

				if len(exp.Invariants) > 0 {
					cfg.Output.Println("")
					cfg.Output.Println("### Invariants to Preserve")
					for _, i := range exp.Invariants {
						cfg.Output.Println("- %s", i)
					}
				}

				cfg.Output.Println("")
			}
		}
	}

	return result, nil
}

// explainMatch generates a human-readable explanation of why a pattern matches a file.
func explainMatch(pattern, file string) string {
	if strings.Contains(pattern, "**") {
		parts := strings.Split(pattern, "**")
		if len(parts) == 2 {
			prefix := strings.TrimSuffix(parts[0], "/")
			suffix := strings.TrimPrefix(parts[1], "/")

			if prefix != "" && suffix != "" {
				return fmt.Sprintf("File is under '%s/' and matches '%s'", prefix, suffix)
			} else if prefix != "" {
				return fmt.Sprintf("File is anywhere under '%s/'", prefix)
			} else if suffix != "" {
				return fmt.Sprintf("File matches pattern '%s' at any depth", suffix)
			}
		}
		return fmt.Sprintf("File matches recursive pattern '%s'", pattern)
	}

	if strings.Contains(pattern, "*") {
		return fmt.Sprintf("File matches wildcard pattern '%s'", pattern)
	}

	return fmt.Sprintf("File exactly matches '%s'", pattern)
}
