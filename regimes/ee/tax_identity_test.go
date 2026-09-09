package ee_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name        string
		inputCode   cbc.Code
		expectedErr string
	}{
		// Valid cases, verified against estonian VAT number validator and EU VIES
		{
			name:      "valid 1",
			inputCode: "100366327",
		},
		{
			name:      "valid 2",
			inputCode: "102090374",
		},
		{
			name:      "empty code",
			inputCode: "",
		},

		// Format errors
		{
			name:        "too short",
			inputCode:   "1234567",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "too long",
			inputCode:   "1234567890",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "not normalized",
			inputCode:   "EE100366327",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "contains letters",
			inputCode:   "1234567A",
			expectedErr: "IDENTITY-01",
		},
		// Checksum errors
		{
			name:        "bad checksum",
			inputCode:   "100366326",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "remainder 1 - no valid check digit",
			inputCode:   "80000000",
			expectedErr: "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "EE", Code: tt.inputCode}
			err := rules.Validate(tID)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
			}
		})
	}
}
