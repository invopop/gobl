package cz_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/cz"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid legal entity codes (8 digits)
		{name: "valid legal entity - Skoda Auto", code: "00177041"},
		{name: "valid legal entity - CEZ", code: "45274649"},
		{name: "valid legal entity", code: "25596641"},
		{name: "valid legal entity - T-Mobile CZ", code: "64949681"},
		{name: "valid legal entity - Komercni banka", code: "45317054"},
		// Valid individual code (10 digits, divisible by 11)
		{name: "valid individual", code: "7103192745"},
		// Valid legal entity with check digit 0 (expected == 10 branch)
		{name: "valid legal entity check digit 0", code: "00200000"},
		// Valid 9-digit codes
		{name: "valid 9-digit special", code: "612345679"},
		{name: "valid 9-digit individual", code: "710319274"},
		// Empty code
		{
			name: "empty code",
			code: "",
			err:  "",
		},
		// Invalid format
		{
			name: "too short",
			code: "0017704",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "00177041234",
			err:  "invalid format",
		},
		{
			name: "contains letters",
			code: "0017704A",
			err:  "invalid format",
		},
		// Bad checksum - legal entity
		{
			name: "bad checksum legal entity",
			code: "00177042",
			err:  "checksum mismatch",
		},
		{
			name: "bad checksum legal entity 2",
			code: "45274640",
			err:  "checksum mismatch",
		},
		// Bad checksum - individual (10 digits, not divisible by 11)
		{
			name: "bad checksum individual",
			code: "7103192746",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CZ", Code: tt.code}
			err := cz.Validate(tID)
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
			name:     "with CZ prefix",
			code:     "CZ00177041",
			expected: "00177041",
		},
		{
			name:     "with spaces",
			code:     "001 770 41",
			expected: "00177041",
		},
		{
			name:     "with dashes",
			code:     "001-770-41",
			expected: "00177041",
		},
		{
			name:     "already normalized",
			code:     "00177041",
			expected: "00177041",
		},
		{
			name:     "empty",
			code:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CZ", Code: tt.code}
			cz.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}
