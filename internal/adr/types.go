// Package adr provides types and functions for working with Architecture Decision Records.
package adr

import (
	"fmt"
	"strings"
)

// Status represents the lifecycle status of an ADR.
type Status string

const (
	StatusProposed   Status = "proposed"
	StatusAdopted    Status = "adopted"
	StatusRejected   Status = "rejected"
	StatusDeprecated Status = "deprecated"
	StatusSuperseded Status = "superseded"
)

// ValidStatuses returns all valid ADR statuses.
func ValidStatuses() []Status {
	return []Status{
		StatusProposed,
		StatusAdopted,
		StatusRejected,
		StatusDeprecated,
		StatusSuperseded,
	}
}

// ParseStatus parses a string into a Status, returning an error if invalid.
func ParseStatus(s string) (Status, error) {
	status := Status(strings.ToLower(strings.TrimSpace(s)))
	for _, valid := range ValidStatuses() {
		if status == valid {
			return status, nil
		}
	}
	return "", fmt.Errorf("invalid status %q: must be one of %v", s, ValidStatuses())
}

// Scope represents the scope of an ADR.
type Scope struct {
	Paths []string `yaml:"paths"`
}

// Frontmatter represents the YAML frontmatter of an ADR.
type Frontmatter struct {
	ADRID        string   `yaml:"adr_id"`
	Title        string   `yaml:"title"`
	Status       Status   `yaml:"status"`
	Date         string   `yaml:"date"`
	Scope        Scope    `yaml:"scope"`
	Tags         []string `yaml:"tags"`
	Constraints  []string `yaml:"constraints"`
	Invariants   []string `yaml:"invariants"`
	Supersedes   []string `yaml:"supersedes"`
	SupersededBy []string `yaml:"superseded_by"`
	RelatedADRs  []string `yaml:"related_adrs"`
}

// ADR represents a complete Architecture Decision Record.
type ADR struct {
	Frontmatter Frontmatter
	Body        string
	Filename    string
	FilePath    string
}

// Number extracts the numeric portion from the ADR ID (e.g., "ADR-0042" -> 42).
func (a *ADR) Number() (int, error) {
	return ExtractNumber(a.Frontmatter.ADRID)
}

// ExtractNumber extracts the numeric portion from an ADR ID string.
func ExtractNumber(adrID string) (int, error) {
	adrID = strings.TrimSpace(adrID)
	if strings.HasPrefix(strings.ToUpper(adrID), "ADR-") {
		adrID = adrID[4:]
	}
	var num int
	_, err := fmt.Sscanf(adrID, "%d", &num)
	if err != nil {
		return 0, fmt.Errorf("cannot extract number from ADR ID %q: %w", adrID, err)
	}
	return num, nil
}
