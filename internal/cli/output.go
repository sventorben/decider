// Package cli provides CLI command implementations for the decider tool.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// OutputFormat represents the output format for commands.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatYAML OutputFormat = "yaml"
)

// ParseOutputFormat parses a string into an OutputFormat.
func ParseOutputFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(s) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "yaml":
		return FormatYAML, nil
	default:
		return "", fmt.Errorf("invalid format %q: must be text, json, or yaml", s)
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
		fmt.Fprintf(o.Writer, format, args...)
	}
}

// Println outputs a line (text format only).
func (o *Output) Println(format string, args ...interface{}) {
	if o.Format == FormatText {
		fmt.Fprintf(o.Writer, format+"\n", args...)
	}
}

// PrintJSON outputs data as JSON.
func (o *Output) PrintJSON(data interface{}) error {
	encoder := json.NewEncoder(o.Writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// Error outputs an error message to stderr.
func (o *Output) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
}

// Warn outputs a warning message to stderr.
func (o *Output) Warn(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
}

// Success outputs a success message.
func (o *Output) Success(format string, args ...interface{}) {
	o.Println("✓ "+format, args...)
}

// Info outputs an informational message.
func (o *Output) Info(format string, args ...interface{}) {
	o.Println(format, args...)
}
