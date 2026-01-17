package cli

import (
	"fmt"
	"path/filepath"

	"github.com/sventorben/decider/internal/index"
	"gopkg.in/yaml.v3"
)

// IndexConfig holds configuration for the index command.
type IndexConfig struct {
	Dir    string
	Check  bool
	Format OutputFormat
	Output *Output
}

// IndexResult holds the result of the index command.
type IndexResult struct {
	Generated bool   `json:"generated"`
	UpToDate  bool   `json:"up_to_date,omitempty"`
	File      string `json:"file"`
	ADRCount  int    `json:"adr_count"`
}

// RunIndex generates or checks the ADR index.
func RunIndex(cfg *IndexConfig) (*IndexResult, error) {
	indexPath := filepath.Join(cfg.Dir, index.IndexFilename)

	if cfg.Check {
		// Check mode: verify index is up-to-date
		upToDate, err := index.Check(cfg.Dir)
		if err != nil {
			return nil, fmt.Errorf("checking index: %w", err)
		}

		result := &IndexResult{
			Generated: false,
			UpToDate:  upToDate,
			File:      indexPath,
		}

		if cfg.Format == FormatJSON {
			_ = cfg.Output.PrintJSON(result)
		} else {
			if upToDate {
				cfg.Output.Success("Index is up-to-date: %s", indexPath)
			} else {
				cfg.Output.Error("Index is out of date. Run 'decider index' to update.")
				return result, fmt.Errorf("index out of date")
			}
		}

		return result, nil
	}

	// Generate mode: create/update index
	idx, err := index.GenerateFromDir(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("generating index: %w", err)
	}

	if err := idx.Write(indexPath); err != nil {
		return nil, fmt.Errorf("writing index: %w", err)
	}

	result := &IndexResult{
		Generated: true,
		File:      indexPath,
		ADRCount:  idx.ADRCount,
	}

	switch cfg.Format {
	case FormatJSON:
		_ = cfg.Output.PrintJSON(result)
	case FormatYAML:
		data, _ := yaml.Marshal(idx)
		_, _ = fmt.Fprint(cfg.Output.Writer, string(data))
	default:
		cfg.Output.Success("Generated %s with %d ADRs", indexPath, idx.ADRCount)
	}

	return result, nil
}
