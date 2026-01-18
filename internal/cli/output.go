// Package cli provides CLI command implementations for the decider tool.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sventorben/decider/internal/toon"
)

// OutputFormat represents the output format for commands.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatTOON OutputFormat = "toon"
	FormatJSON OutputFormat = "json"
	FormatYAML OutputFormat = "yaml"
)

// DefaultStructuredFormat is the default format for machine-readable output.
const DefaultStructuredFormat = FormatTOON

// ParseOutputFormat parses a string into an OutputFormat.
func ParseOutputFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(s) {
	case "text", "":
		return FormatText, nil
	case "toon":
		return FormatTOON, nil
	case "json":
		return FormatJSON, nil
	case "yaml":
		return FormatYAML, nil
	default:
		return "", fmt.Errorf("invalid format %q: must be text, toon, json, or yaml", s)
	}
}

// Output handles writing output in different formats.
type Output struct {
	Format OutputFormat
	Writer io.Writer
}

// NewOutput creates a new Output with the given format.
func NewOutput(format OutputFormat) *Output {
	return &Output{
		Format: format,
		Writer: os.Stdout,
	}
}

// Print outputs a message (text format only).
func (o *Output) Print(format string, args ...interface{}) {
	if o.Format == FormatText {
		_, _ = fmt.Fprintf(o.Writer, format, args...)
	}
}

// Println outputs a line (text format only).
func (o *Output) Println(format string, args ...interface{}) {
	if o.Format == FormatText {
		_, _ = fmt.Fprintf(o.Writer, format+"\n", args...)
	}
}

// PrintStructured outputs data in the configured structured format (TOON, JSON, or YAML).
func (o *Output) PrintStructured(data interface{}) error {
	switch o.Format {
	case FormatTOON:
		return o.PrintTOON(data)
	case FormatJSON:
		return o.PrintJSON(data)
	default:
		return o.PrintJSON(data)
	}
}

// PrintTOON outputs data as TOON.
func (o *Output) PrintTOON(data interface{}) error {
	enc := toon.NewEncoder(o.Writer)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// PrintJSON outputs data as JSON.
func (o *Output) PrintJSON(data interface{}) error {
	encoder := json.NewEncoder(o.Writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// Error outputs an error message to stderr.
func (o *Output) Error(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
}

// Warn outputs a warning message to stderr.
func (o *Output) Warn(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
}

// Success outputs a success message.
func (o *Output) Success(format string, args ...interface{}) {
	o.Println("✓ "+format, args...)
}

// Info outputs an informational message.
func (o *Output) Info(format string, args ...interface{}) {
	o.Println(format, args...)
}

// IsStructuredFormat returns true if the format is a structured data format.
func (o *Output) IsStructuredFormat() bool {
	return o.Format == FormatTOON || o.Format == FormatJSON || o.Format == FormatYAML
}
