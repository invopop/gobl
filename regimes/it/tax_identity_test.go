package it_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		code     cbc.Code
		expected cbc.Code
		typ      cbc.Key
	}{
		{
			code:     "12345678901",
			expected: "12345678901",
		},
		{
			code:     "123-456-789-01",
			expected: "12345678901",
		},
		{
			code:     "123456 789 01",
			expected: "12345678901",
		},
		{
			code:     "IT 12345678901",
			expected: "12345678901",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: l10n.IT, Code: ts.code}
		err := it.Calculate(tID)
		assert.NoError(t, err)
		assert.Equal(t, ts.expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		zone l10n.Code
		err  string
	}{
		{name: "good 1", code: "12345678903"},
		{name: "good 2", code: "13029381004"},
		{name: "good 3", code: "10182640150"},
		{
			name: "empty",
			code: "",
			err:  "",
		},
		{
			name: "too long",
			code: "123456789001",
			err:  "invalid length",
		},
		{
			name: "too short",
			code: "1234567890",
			err:  "invalid length",
		},
		{
			name: "not normalized",
			code: "12.449.965-439",
			err:  "contains invalid characters",
		},
		{
			name: "includes non-numeric characters",
			code: "A764352056Z",
			err:  "contains invalid characters",
		},
		{
			name: "invalid check digit",
			code: "12345678901",
			err:  "invalid check digit",
		},
		{
			name: "invalid check digit",
			code: "13029381009",
			err:  "invalid check digit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.IT, Code: tt.code}
			err := it.Validate(tID)
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

func TestTaxIdentityValidateGeneralCases(t *testing.T) {
	tests := []struct {
		name string
		tID  *tax.Identity
		err  string
	}{
		{
			name: "just country",
			tID:  &tax.Identity{Country: l10n.IT},
			err:  "",
		},
		{
			name: "no type, assume biz",
			tID:  &tax.Identity{Country: l10n.IT, Code: "12345678903"},
			err:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := it.Validate(tt.tID)
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
