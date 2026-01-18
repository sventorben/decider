package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestParseOutputFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    OutputFormat
		wantErr bool
	}{
		{"", FormatText, false},
		{"text", FormatText, false},
		{"TEXT", FormatText, false},
		{"toon", FormatTOON, false},
		{"TOON", FormatTOON, false},
		{"json", FormatJSON, false},
		{"JSON", FormatJSON, false},
		{"yaml", FormatYAML, false},
		{"YAML", FormatYAML, false},
		{"invalid", "", true},
		{"xml", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseOutputFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseOutputFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDefaultStructuredFormat(t *testing.T) {
	// TOON should be the default structured format
	if DefaultStructuredFormat != FormatTOON {
		t.Errorf("DefaultStructuredFormat = %v, want %v", DefaultStructuredFormat, FormatTOON)
	}
}

func TestOutputPrintStructured(t *testing.T) {
	data := map[string]interface{}{
		"name":  "test",
		"count": 42,
	}

	tests := []struct {
		format    OutputFormat
		checkTOON bool
		checkJSON bool
	}{
		{FormatTOON, true, false},
		{FormatJSON, false, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			var buf bytes.Buffer
			output := &Output{
				Format: tt.format,
				Writer: &buf,
			}

			err := output.PrintStructured(data)
			if err != nil {
				t.Fatalf("PrintStructured failed: %v", err)
			}

			result := buf.String()
			if result == "" {
				t.Error("PrintStructured produced empty output")
			}

			if tt.checkJSON {
				// Should be valid JSON
				var v interface{}
				if err := json.Unmarshal([]byte(result), &v); err != nil {
					t.Errorf("JSON output is not valid JSON: %v", err)
				}
			}

			if tt.checkTOON {
				// TOON should contain expected keys without JSON quotes
				if !containsStr(result, "name:") || !containsStr(result, "count:") {
					t.Errorf("TOON output missing expected keys: %s", result)
				}
			}
		})
	}
}

func TestIsStructuredFormat(t *testing.T) {
	tests := []struct {
		format OutputFormat
		want   bool
	}{
		{FormatText, false},
		{FormatTOON, true},
		{FormatJSON, true},
		{FormatYAML, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			output := &Output{Format: tt.format}
			if got := output.IsStructuredFormat(); got != tt.want {
				t.Errorf("IsStructuredFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
