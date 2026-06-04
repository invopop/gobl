package no_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
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
		{name: "already normalized", input: "923456783", expected: "923456783"},
		{name: "with prefix and suffix", input: "NO 923 456 783 MVA", expected: "923456783"},
		{name: "lowercase prefix", input: "no923456783mva", expected: "923456783"},
		{name: "spaces only", input: "923 456 783", expected: "923456783"},
		{name: "with dashes", input: "923-456-783", expected: "923456783"},
		{name: "prefix only", input: "NO923456783", expected: "923456783"},
		{name: "suffix only", input: "923456783MVA", expected: "923456783"},
		{name: "empty code", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "NO", Code: tt.input}
			no.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}

	t.Run("nil identity", func(t *testing.T) {
		assert.NotPanics(t, func() {
			no.Normalize((*tax.Identity)(nil))
		})
	})
}

func TestValidateTaxIdentity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		code  cbc.Code
		valid bool
	}{
		{name: "valid code", code: "923456783", valid: true},
		{name: "valid code starting with 8", code: "889640782", valid: true},
		{name: "valid code not starting with 8 or 9", code: "123456785", valid: true},
		{name: "empty code", code: "", valid: true},
		{name: "too short", code: "92345678"},
		{name: "too long", code: "9234567830"},
		{name: "non-numeric", code: "92345678A"},
		{name: "bad check digit", code: "923456780"},
		// 850000000: sum = 8*3 + 5*2 = 34, 34 % 11 = 1, check = 10 → invalid
		{name: "check digit would be 10", code: "850000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "NO", Code: tt.code}
			err := rules.Validate(tID)
			if tt.valid {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "invalid organisasjonsnummer")
			}
		})
	}
}
