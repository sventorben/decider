package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sventorben/decider/internal/adr"
	"github.com/sventorben/decider/internal/index"
	"github.com/sventorben/decider/internal/validate"
)

// NewConfig holds configuration for the new command.
type NewConfig struct {
	Title   string
	Dir     string
	Tags    []string
	Paths   []string
	Owners  []string
	Status  string
	NoIndex bool
	Format  OutputFormat
	Output  *Output
}

// NewResult holds the result of creating a new ADR.
type NewResult struct {
	ADRID    string `json:"adr_id"`
	Title    string `json:"title"`
	File     string `json:"file"`
	FilePath string `json:"file_path"`
	Number   int    `json:"number"`
}

// RunNew creates a new ADR.
func RunNew(cfg *NewConfig) (*NewResult, error) {
	// Validate inputs
	if err := validate.ValidateTitle(cfg.Title); err != nil {
		return nil, fmt.Errorf("invalid title: %w", err)
	}
	if err := validate.ValidateTags(cfg.Tags); err != nil {
		return nil, fmt.Errorf("invalid tags: %w", err)
	}
	if err := validate.ValidateScopePaths(cfg.Paths); err != nil {
		return nil, fmt.Errorf("invalid paths: %w", err)
	}

	// Parse and validate status
	status, err := adr.ParseStatus(cfg.Status)
	if err != nil {
		return nil, err
	}

	// Find next number
	number, err := adr.FindNextNumber(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("finding next ADR number: %w", err)
	}

	// Generate filename
	filename := adr.GenerateFilename(number, cfg.Title)
	filePath := filepath.Join(cfg.Dir, filename)
	adrID := fmt.Sprintf("ADR-%04d", number)

	// Create frontmatter
	fm := &adr.Frontmatter{
		ADRID:        adrID,
		Title:        cfg.Title,
		Status:       status,
		Date:         time.Now().Format("2006-01-02"),
		Scope:        adr.Scope{Paths: cfg.Paths},
		Tags:         cfg.Tags,
		Constraints:  []string{},
		Invariants:   []string{},
		Supersedes:   []string{},
		SupersededBy: []string{},
		RelatedADRs:  []string{},
	}

	// Generate content
	content, err := generateADRContent(fm, cfg.Title)
	if err != nil {
		return nil, fmt.Errorf("generating ADR content: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		return nil, fmt.Errorf("creating directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("writing ADR file: %w", err)
	}

	result := &NewResult{
		ADRID:    adrID,
		Title:    cfg.Title,
		File:     filename,
		FilePath: filePath,
		Number:   number,
	}

	// Update index unless disabled
	if !cfg.NoIndex {
		if err := index.WriteToDir(cfg.Dir); err != nil {
			// Non-fatal: warn but don't fail
			cfg.Output.Error("warning: could not update index: %v", err)
		}
	}

	// Output result
	if cfg.Format == FormatJSON {
		cfg.Output.PrintJSON(result)
	} else {
		cfg.Output.Success("Created %s", filePath)
		cfg.Output.Println("  ADR ID: %s", adrID)
		cfg.Output.Println("  Title:  %s", cfg.Title)
		cfg.Output.Println("  Status: %s", status)
	}

	return result, nil
}

func generateADRContent(fm *adr.Frontmatter, title string) (string, error) {
	fmStr, err := adr.SerializeFrontmatter(fm)
	if err != nil {
		return "", err
	}

	body := fmt.Sprintf(`
# %s: %s

## Context

_Describe the context and background that led to this decision. What problem are we solving? What forces are at play?_

Decision drivers:
- _Key driver 1 that influenced the decision_
- _Key driver 2_
- _Key driver 3_

## Decision

_State the decision clearly and concisely._

### [Chosen Option]: Adopted

**Adopted because:**
- _Clear, concrete reason why this option was chosen_
- _Tie reasons to decision drivers above_
- _Technical, operational, or strategic justification_

**Adopted despite:**
- _Known downside or trade-off we consciously accepted_
- _Cost or weakness compared to alternatives_
- _Risk we are taking on_

## Alternatives Considered

### [Alternative A]: Rejected

**Rejected because:**
- _Clear, concrete reason why this option was not chosen_
- _Technical, organizational, or strategic reason_
- _How it failed to meet decision drivers_

**Rejected despite:**
- _Legitimate strength of this option_
- _Benefit that made it attractive_
- _Reason it was seriously considered_

## Consequences

**Positive:**
- _First positive consequence_
- _Second positive consequence_

**Negative:**
- _First negative consequence (and mitigation if any)_
- _Second negative consequence_
`, fm.ADRID, title)

	return fmStr + strings.TrimPrefix(body, "\n"), nil
}
