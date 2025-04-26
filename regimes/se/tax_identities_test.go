package se_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/se"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    cbc.Code
		expected cbc.Code
	}{
		{name: "already normalized", input: "SE012345678901", expected: "SE012345678901"},
		{name: "lowercase input", input: "se012345678901", expected: "SE012345678901"},
		{name: "mixed case input", input: "Se012345678901", expected: "SE012345678901"},
		{name: "extra spaces", input: "  SE 0123456789 01  ", expected: "SE012345678901"},
		{name: "special characters", input: "SE-0123456789-01", expected: "SE012345678901"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SE", Code: tt.input}

			se.Normalize(tID)

			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid TAX ID 1", code: "SE012345678901"},
		{
			name: "too short",
			code: "SE123456701",
			err:  "invalid length",
		},
		{
			name: "too long",
			code: "SE1234567891012",
			err:  "invalid length",
		},
		{
			name: "invalid check digit",
			code: "SE123456789100",
			err:  "invalid check digit, expected 01",
		},
		{
			name: "doesn't start with SE",
			code: "ES123456789101",
			err:  se.ErrInvalidTaxIDCountryPrefix.Error(),
		},
		{
			name: "contains non-numeric characters",
			code: "SE123456789A01",
			err:  "invalid characters, expected numeric",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SE", Code: tt.code}

			err := se.Validate(tID)

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
