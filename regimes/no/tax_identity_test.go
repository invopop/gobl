package no_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/no"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid 1", code: "923609016"},
		{name: "valid 2", code: "982463718"},
		{name: "valid 3", code: "889640782"},
		{name: "valid with check digit 0", code: "100000040"},
		{
			name: "empty code",
			code: "",
		},
		{
			name: "too short",
			code: "12345678",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "1234567890",
			err:  "invalid format",
		},
		{
			name: "contains letters",
			code: "12345678A",
			err:  "invalid format",
		},
		{
			name: "bad checksum",
			code: "923609017",
			err:  "checksum mismatch",
		},
		{
			name: "remainder 1 - no valid check digit",
			code: "100000130",
			err:  "invalid number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "NO", Code: tt.code}
			err := no.Validate(tID)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name     string
		code     cbc.Code
		expected cbc.Code
	}{
		{
			name:     "with NO prefix",
			code:     "NO923609016",
			expected: "923609016",
		},
		{
			name:     "with spaces",
			code:     "923 609 016",
			expected: "923609016",
		},
		{
			name:     "with dashes",
			code:     "923-609-016",
			expected: "923609016",
		},
		{
			name:     "with MVA suffix",
			code:     "923609016MVA",
			expected: "923609016",
		},
		{
			name:     "with NO prefix and MVA suffix",
			code:     "NO923609016MVA",
			expected: "923609016",
		},
		{
			name:     "lowercase input",
			code:     "no923609016mva",
			expected: "923609016",
		},
		{
			name:     "already normalized",
			code:     "923609016",
			expected: "923609016",
		},
		{
			name:     "empty",
			code:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "NO", Code: tt.code}
			no.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}
