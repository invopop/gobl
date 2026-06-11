package no_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/regimes/no"
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
		{name: "already normalized", input: "923456783MVA", expected: "923456783MVA"},
		{name: "bare org number gains the suffix", input: "923456783", expected: "923456783MVA"},
		{name: "with prefix and suffix", input: "NO 923 456 783 MVA", expected: "923456783MVA"},
		{name: "lowercase prefix and suffix", input: "no923456783mva", expected: "923456783MVA"},
		{name: "spaces only", input: "923 456 783", expected: "923456783MVA"},
		{name: "with dashes", input: "923-456-783", expected: "923456783MVA"},
		{name: "prefix only", input: "NO923456783", expected: "923456783MVA"},
		{name: "empty code", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: no.CountryCode, Code: tt.input}
			norm.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		code  cbc.Code
		valid bool
	}{
		{name: "valid code", code: "923456783MVA", valid: true},
		{name: "valid code starting with 8", code: "889640782MVA", valid: true},
		{name: "valid code not starting with 8 or 9", code: "123456785MVA", valid: true},
		{name: "empty code", code: "", valid: true},
		{name: "missing MVA suffix", code: "923456783"},
		{name: "too short", code: "92345678MVA"},
		{name: "too long", code: "9234567830MVA"},
		{name: "non-numeric", code: "92345678AMVA"},
		{name: "bad check digit", code: "923456780MVA"},
		// 850000000: sum = 8*3 + 5*2 = 34, 34 % 11 = 1, check = 10 → invalid
		{name: "check digit would be 10", code: "850000000MVA"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "NO", Code: tt.code}
			err := rules.Validate(tID)
			if tt.valid {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "invalid Norwegian VAT number")
			}
		})
	}
}
