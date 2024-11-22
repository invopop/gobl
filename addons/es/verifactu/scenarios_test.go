package verifactu

import (
	"testing"
)

func TestScenarios(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		tags     []string
		want     bool
	}{
		{
			name:     "standard invoice",
			scenario: "standard",
			tags:     []string{},
			want:     true,
		},
		{
			name:     "simplified invoice",
			scenario: "simplified",
			tags:     []string{"simplified"},
			want:     true,
		},
		{
			name:     "corrective invoice",
			scenario: "corrective",
			tags:     []string{"corrective"},
			want:     true,
		},
		{
			name:     "invalid scenario",
			scenario: "invalid",
			tags:     []string{},
			want:     false,
		},
		{
			name:     "simplified with wrong tags",
			scenario: "simplified",
			tags:     []string{"corrective"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateScenario(tt.scenario, tt.tags)
			if got != tt.want {
				t.Errorf("ValidateScenario(%v, %v) = %v, want %v",
					tt.scenario, tt.tags, got, tt.want)
			}
		})
	}
}
