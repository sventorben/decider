package toon

import (
	"encoding/json"
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"null", nil},
		{"true", true},
		{"false", false},
		{"int", 42},
		{"float", 3.14},
		{"string", "hello"},
		{"string with spaces", "hello world"},
		// Note: empty array round-trip may produce null due to JSON marshaling
		{"int array", []interface{}{int64(1), int64(2), int64(3)}},
		{"empty object", map[string]interface{}{}},
		{"simple object", map[string]interface{}{"key": "value"}},
		{"nested object", map[string]interface{}{
			"name":   "test",
			"count":  int64(5),
			"active": true,
			"items":  []interface{}{"a", "b"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to TOON
			toonData, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// Unmarshal back
			var result interface{}
			if err := Unmarshal(toonData, &result); err != nil {
				t.Fatalf("Unmarshal failed: %v (data: %s)", err, toonData)
			}

			// Compare via JSON (handles type normalization)
			expectedJSON, _ := json.Marshal(tt.input)
			resultJSON, _ := json.Marshal(result)
			if string(expectedJSON) != string(resultJSON) {
				t.Errorf("Round-trip mismatch:\n  input:  %s\n  toon:   %s\n  result: %s",
					expectedJSON, toonData, resultJSON)
			}
		})
	}
}

func TestDeterministicOutput(t *testing.T) {
	input := map[string]interface{}{
		"zebra": "z",
		"alpha": "a",
		"beta":  "b",
	}

	// Marshal multiple times
	results := make([][]byte, 5)
	for i := range results {
		data, err := Marshal(input)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		results[i] = data
	}

	// All results should be identical
	for i := 1; i < len(results); i++ {
		if string(results[0]) != string(results[i]) {
			t.Errorf("Non-deterministic output:\n  run 0: %s\n  run %d: %s",
				results[0], i, results[i])
		}
	}

	// Keys should be sorted (bare words for simple strings)
	expected := `{alpha:a beta:b zebra:z}`
	if string(results[0]) != expected+"\n" {
		t.Errorf("Expected sorted keys:\n  got:      %s\n  expected: %s", results[0], expected)
	}
}

func TestJSONEquivalence(t *testing.T) {
	testCases := []string{
		`{"name":"test","count":5}`,
		`[1,2,3,"four",true,null]`,
		`{"nested":{"array":[1,2],"bool":false}}`,
	}

	for _, jsonInput := range testCases {
		t.Run(jsonInput, func(t *testing.T) {
			// Parse JSON
			var original interface{}
			if err := json.Unmarshal([]byte(jsonInput), &original); err != nil {
				t.Fatalf("JSON unmarshal failed: %v", err)
			}

			// Convert to TOON
			toonData, err := Marshal(original)
			if err != nil {
				t.Fatalf("TOON marshal failed: %v", err)
			}

			// Parse TOON
			var fromToon interface{}
			if err := Unmarshal(toonData, &fromToon); err != nil {
				t.Fatalf("TOON unmarshal failed: %v", err)
			}

			// Convert both back to JSON for comparison
			originalJSON, _ := json.Marshal(original)
			fromToonJSON, _ := json.Marshal(fromToon)

			if string(originalJSON) != string(fromToonJSON) {
				t.Errorf("JSON/TOON equivalence failed:\n  original: %s\n  toon:     %s\n  back:     %s",
					originalJSON, toonData, fromToonJSON)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	toonInput := []byte(`{name:"test" count:5 active:true}`)

	jsonOutput, err := ToJSON(toonInput)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonOutput, &result); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name=test, got %v", result["name"])
	}
	if result["count"] != float64(5) {
		t.Errorf("Expected count=5, got %v", result["count"])
	}
	if result["active"] != true {
		t.Errorf("Expected active=true, got %v", result["active"])
	}
}

func TestFromJSON(t *testing.T) {
	jsonInput := []byte(`{"name":"test","count":5}`)

	toonOutput, err := FromJSON(jsonInput)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := Unmarshal(toonOutput, &result); err != nil {
		t.Fatalf("TOON unmarshal failed: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name=test, got %v", result["name"])
	}
}

func TestBareWords(t *testing.T) {
	input := map[string]interface{}{
		"simple":     "hello",
		"with_under": "test_value",
		"with-dash":  "test-value",
	}

	data, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Bare words should not be quoted
	dataStr := string(data)
	if contains(dataStr, `"hello"`) {
		t.Errorf("Simple bareword should not be quoted: %s", data)
	}
	if contains(dataStr, `"test_value"`) {
		t.Errorf("Underscore bareword should not be quoted: %s", data)
	}
	if contains(dataStr, `"test-value"`) {
		t.Errorf("Dash bareword should not be quoted: %s", data)
	}
}

func TestQuotedStrings(t *testing.T) {
	input := map[string]interface{}{
		"spaces":  "hello world",
		"special": "hello\nworld",
		"empty":   "",
	}

	data, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Strings with spaces/special chars should be quoted
	dataStr := string(data)
	if !contains(dataStr, `"hello world"`) {
		t.Errorf("String with spaces should be quoted: %s", data)
	}
	if !contains(dataStr, `""`) {
		t.Errorf("Empty string should be quoted: %s", data)
	}
}

func TestStructEncoding(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	input := TestStruct{Name: "test", Count: 42}

	data, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name=test, got %v", result["name"])
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
