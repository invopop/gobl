package se_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/se"
	"github.com/invopop/gobl/rules"
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
		{name: "already normalized", input: "SE012345678901", expected: "012345678901"},
		{name: "lowercase input", input: "se012345678901", expected: "012345678901"},
		{name: "mixed case input", input: "Se012345678901", expected: "012345678901"},
		{name: "extra spaces", input: "  SE 0123456789 01  ", expected: "012345678901"},
		{name: "special characters", input: "SE-0123456789-01", expected: "012345678901"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SE", Code: tt.input}

			se.Normalize(tID)

			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestTaxIdentityRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid Tax ID", code: "202100548901"},
		{name: "empty Tax ID", code: ""},
		{
			name: "too short",
			code: "20210054890",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "2021005489001",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid check digit",
			code: "202100548900",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid checksum",
			code: "202100548801",
			err:  "IDENTITY-01",
		},
		{
			name: "contains non-numeric characters",
			code: "202100548A01",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SE", Code: tt.code}

			err := rules.Validate(tID)

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
