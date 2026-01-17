package adr

import "testing"

func TestParseStatus(t *testing.T) {
	tests := []struct {
		input   string
		want    Status
		wantErr bool
	}{
		{"proposed", StatusProposed, false},
		{"adopted", StatusAdopted, false},
		{"rejected", StatusRejected, false},
		{"deprecated", StatusDeprecated, false},
		{"superseded", StatusSuperseded, false},
		{"ADOPTED", StatusAdopted, false},
		{"  proposed  ", StatusProposed, false},
		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractNumber(t *testing.T) {
	tests := []struct {
		input   string
		want    int
		wantErr bool
	}{
		{"ADR-0001", 1, false},
		{"ADR-0042", 42, false},
		{"ADR-1234", 1234, false},
		{"adr-0001", 1, false},
		{"0001", 1, false},
		{"42", 42, false},
		{"invalid", 0, true},
		{"ADR-", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ExtractNumber(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractNumber(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractNumber(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidStatuses(t *testing.T) {
	statuses := ValidStatuses()
	if len(statuses) != 5 {
		t.Errorf("ValidStatuses() returned %d statuses, want 5", len(statuses))
	}

	expected := []Status{StatusProposed, StatusAdopted, StatusRejected, StatusDeprecated, StatusSuperseded}
	for i, s := range expected {
		if statuses[i] != s {
			t.Errorf("ValidStatuses()[%d] = %v, want %v", i, statuses[i], s)
		}
	}
}
