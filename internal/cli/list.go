package cli

import (
	"fmt"
	"path/filepath"

	"github.com/sventorben/decider/internal/adr"
	"github.com/sventorben/decider/internal/glob"
	"github.com/sventorben/decider/internal/index"
)

// ListConfig holds configuration for the list command.
type ListConfig struct {
	Dir    string
	Status string
	Tags   []string
	Path   string
	Format OutputFormat
	Output *Output
}

// ListEntry represents an ADR in the list output.
type ListEntry struct {
	ADRID  string `json:"adr_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Date   string `json:"date"`
	File   string `json:"file"`
}

// ListResult holds the result of the list command.
type ListResult struct {
	Count int         `json:"count"`
	ADRs  []ListEntry `json:"adrs"`
}

// RunList lists ADRs with optional filters.
func RunList(cfg *ListConfig) (*ListResult, error) {
	var entries []ListEntry

	// Try to use index if it exists
	indexPath := filepath.Join(cfg.Dir, index.IndexFilename)
	idx, err := index.Load(indexPath)
	if err == nil {
		// Use index
		for _, e := range idx.ADRs {
			if !matchesFilters(e, cfg) {
				continue
			}
			entries = append(entries, ListEntry{
				ADRID:  e.ADRID,
				Title:  e.Title,
				Status: e.Status,
				Date:   e.Date,
				File:   e.File,
			})
		}
	} else {
		// Fallback: scan ADR files
		adrs, err := adr.LoadAllADRs(cfg.Dir)
		if err != nil {
			return nil, fmt.Errorf("loading ADRs: %w", err)
		}

		for _, a := range adrs {
			entry := index.Entry{
				ADRID:      a.Frontmatter.ADRID,
				Title:      a.Frontmatter.Title,
				Status:     string(a.Frontmatter.Status),
				Date:       a.Frontmatter.Date,
				Tags:       a.Frontmatter.Tags,
				ScopePaths: a.Frontmatter.Scope.Paths,
				File:       a.Filename,
			}
			if !matchesFilters(entry, cfg) {
				continue
			}
			entries = append(entries, ListEntry{
				ADRID:  a.Frontmatter.ADRID,
				Title:  a.Frontmatter.Title,
				Status: string(a.Frontmatter.Status),
				Date:   a.Frontmatter.Date,
				File:   a.Filename,
			})
		}
	}

	result := &ListResult{
		Count: len(entries),
		ADRs:  entries,
	}

	// Output
	if cfg.Format == FormatJSON {
		_ = cfg.Output.PrintJSON(result)
	} else {
		if len(entries) == 0 {
			cfg.Output.Println("No ADRs found.")
		} else {
			cfg.Output.Println("%-10s %-12s %-12s %s", "ADR ID", "Status", "Date", "Title")
			cfg.Output.Println("%-10s %-12s %-12s %s", "------", "------", "----", "-----")
			for _, e := range entries {
				cfg.Output.Println("%-10s %-12s %-12s %s", e.ADRID, e.Status, e.Date, e.Title)
			}
			cfg.Output.Println("\nTotal: %d ADR(s)", len(entries))
		}
	}

	return result, nil
}

func matchesFilters(entry index.Entry, cfg *ListConfig) bool {
	// Filter by status
	if cfg.Status != "" && entry.Status != cfg.Status {
		return false
	}

	// Filter by tags (any match)
	if len(cfg.Tags) > 0 {
		found := false
		for _, filterTag := range cfg.Tags {
			for _, tag := range entry.Tags {
				if tag == filterTag {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by path (matches against scope.paths)
	if cfg.Path != "" {
		found := false
		for _, scopePath := range entry.ScopePaths {
			if glob.Match(scopePath, cfg.Path) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
