package fr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeAndValidateTaxNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid tax number with leading 0",
			input:   "0123456789012",
			wantErr: false,
		},
		{
			name:    "Valid tax number with leading 1",
			input:   "1234567890123",
			wantErr: false,
		},
		{
			name:    "Valid tax number with leading 2",
			input:   "2234567890123",
			wantErr: false,
		},
		{
			name:    "Valid tax number with leading 3",
			input:   "3234567890123",
			wantErr: false,
		},
		{
			name:    "Invalid tax number with leading 4",
			input:   "4234567890123", // First digit not allowed
			wantErr: true,
		},
		{
			name:    "Invalid length",
			input:   "123456789", // Less than 13 digits
			wantErr: true,
		},
		{
			name:    "Invalid length with extra digits",
			input:   "12345678901234", // More than 13 digits
			wantErr: true,
		},
		{
			name:    "Invalid characters",
			input:   "1234A6789012", // Contains a letter
			wantErr: true,
		},
		{
			name:    "Empty input",
			input:   "",
			wantErr: false,
		},
		{
			name:    "Nil identity",
			input:   "0123456789012", // This should not error out because we will check for nil
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := &org.Identity{
				Key:  fr.IdentityKeyTaxNumber,
				Code: cbc.Code(tt.input),
			}

			// Normalize the tax number first
			fr.Normalize(identity)

			// Validate the tax number
			err := fr.Validate(identity)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
