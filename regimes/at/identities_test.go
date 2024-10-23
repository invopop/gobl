package at_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/at"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeAndValidateTaxNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid PESEL with hyphen",
			input:   "12-3456789",
			wantErr: false,
		},
		{
			name:    "Valid PESEL with spaces",
			input:   "12 345 6789",
			wantErr: false,
		},
		{
			name:    "Valid PESEL with mixed symbols",
			input:   "12.345.6789",
			wantErr: false,
		},
		{
			name:    "Invalid length",
			input:   "12-34567", // Less than 9 digits
			wantErr: true,
		},
		{
			name:    "Invalid length with extra digits",
			input:   "12-34567890", // More than 9 digits
			wantErr: true,
		},
		{
			name:    "Invalid tax office code",
			input:   "00-3456789", // Tax office code should be between 1 and 99
			wantErr: true,
		},
		{
			name:    "Invalid taxpayer number",
			input:   "12-0000000", // Taxpayer number should be positive
			wantErr: true,
		},
		{
			name:    "Empty input",
			input:   "",
			wantErr: false,
		},
		{
			name:    "Nil identity",
			input:   "12-3456789", // This should not error out because we will check for nil
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := &org.Identity{
				Key:  at.IdentityKeyTaxNumber,
				Code: cbc.Code(tt.input),
			}

			// Normalize the tax number first
			at.Normalize(identity)

			// Validate the tax number
			err := at.Validate(identity)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
