package cli

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sventorben/decider/internal/adr"
)

// ShowConfig holds configuration for the show command.
type ShowConfig struct {
	ID     string // ADR-NNNN, NNNN, or filename
	Dir    string
	Format OutputFormat
	Output *Output
}

// ShowResult holds the result of the show command.
type ShowResult struct {
	ADRID       string   `json:"adr_id"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	Date        string   `json:"date"`
	Tags        []string `json:"tags,omitempty"`
	ScopePaths  []string `json:"scope_paths,omitempty"`
	Constraints []string `json:"constraints,omitempty"`
	Invariants  []string `json:"invariants,omitempty"`
	Decision    string   `json:"decision,omitempty"`
	File        string   `json:"file"`
}

// RunShow displays details of a specific ADR.
func RunShow(cfg *ShowConfig) (*ShowResult, error) {
	// Resolve ADR file
	filePath, err := resolveADRPath(cfg.ID, cfg.Dir)
	if err != nil {
		return nil, err
	}

	// Load ADR
	a, err := adr.LoadADR(filePath)
	if err != nil {
		return nil, err
	}

	// Extract decision section
	decision := extractSection(a.Body, "Decision")

	result := &ShowResult{
		ADRID:       a.Frontmatter.ADRID,
		Title:       a.Frontmatter.Title,
		Status:      string(a.Frontmatter.Status),
		Date:        a.Frontmatter.Date,
		Tags:        a.Frontmatter.Tags,
		ScopePaths:  a.Frontmatter.Scope.Paths,
		Constraints: a.Frontmatter.Constraints,
		Invariants:  a.Frontmatter.Invariants,
		Decision:    decision,
		File:        a.Filename,
	}

	// Output
	if cfg.Format == FormatTOON || cfg.Format == FormatJSON {
		_ = cfg.Output.PrintStructured(result)
	} else {
		cfg.Output.Println("# %s: %s", result.ADRID, result.Title)
		cfg.Output.Println("")
		cfg.Output.Println("Status: %s", result.Status)
		cfg.Output.Println("Date:   %s", result.Date)
		cfg.Output.Println("File:   %s", result.File)

		if len(result.Tags) > 0 {
			cfg.Output.Println("Tags:   %s", strings.Join(result.Tags, ", "))
		}

		if len(result.ScopePaths) > 0 {
			cfg.Output.Println("")
			cfg.Output.Println("## Scope Paths")
			for _, p := range result.ScopePaths {
				cfg.Output.Println("  - %s", p)
			}
		}

		if len(result.Constraints) > 0 {
			cfg.Output.Println("")
			cfg.Output.Println("## Constraints")
			for _, c := range result.Constraints {
				cfg.Output.Println("  - %s", c)
			}
		}

		if len(result.Invariants) > 0 {
			cfg.Output.Println("")
			cfg.Output.Println("## Invariants")
			for _, i := range result.Invariants {
				cfg.Output.Println("  - %s", i)
			}
		}

		if decision != "" {
			cfg.Output.Println("")
			cfg.Output.Println("## Decision")
			cfg.Output.Println(decision)
		}
	}

	return result, nil
}

// resolveADRPath finds the file path for an ADR given an ID, number, or filename.
func resolveADRPath(id string, dir string) (string, error) {
	// If it's already a path or filename ending in .md
	if strings.HasSuffix(id, ".md") {
		if filepath.IsAbs(id) {
			return id, nil
		}
		return filepath.Join(dir, id), nil
	}

	// Extract number from ADR-NNNN or just NNNN
	var num int
	if strings.HasPrefix(strings.ToUpper(id), "ADR-") {
		_, _ = fmt.Sscanf(id[4:], "%d", &num)
	} else {
		_, _ = fmt.Sscanf(id, "%d", &num)
	}

	if num == 0 {
		return "", fmt.Errorf("cannot parse ADR identifier: %s", id)
	}

	// Find file matching the number
	files, err := adr.ListADRFiles(dir)
	if err != nil {
		return "", err
	}

	pattern := fmt.Sprintf("%04d-", num)
	for _, f := range files {
		if strings.HasPrefix(f, pattern) {
			return filepath.Join(dir, f), nil
		}
	}

	return "", fmt.Errorf("ADR not found: %s", id)
}

// extractSection extracts content from a markdown section.
func extractSection(body, sectionTitle string) string {
	// Pattern: ## Section Title followed by content until next ## or end
	pattern := fmt.Sprintf(`(?i)##\s*%s\s*\n([\s\S]*?)(?:\n##|\z)`, regexp.QuoteMeta(sectionTitle))
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 {
		return ""
	}

	return strings.TrimSpace(matches[1])
}
