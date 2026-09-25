package pe_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/regimes/pe"
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
		{name: "already normalized", input: "20131312955", expected: "20131312955"},
		{name: "with hyphens", input: "20-13131295-5", expected: "20131312955"},
		{name: "with spaces", input: "20 131 312 955", expected: "20131312955"},
		{name: "with dots", input: "20.131.312.955", expected: "20131312955"},
		{name: "with country prefix", input: "PE20131312955", expected: "20131312955"},
		{name: "lowercase country prefix", input: "pe 20131312955", expected: "20131312955"},
		{name: "empty code", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: pe.CountryCode, Code: tt.input}
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
		// Real RUCs whose check digits exercise the mod-11 convention.
		{name: "valid legal entity (SUNAT)", code: "20131312955", valid: true},
		{name: "valid legal entity (BCP)", code: "20100047218", valid: true},
		{name: "valid legal entity (Telefónica)", code: "20100017491", valid: true},
		// Structurally valid RUCs for the remaining taxpayer type prefixes.
		{name: "valid natural person (10)", code: "10401234565", valid: true},
		{name: "valid prefix 15", code: "15000000016", valid: true},
		{name: "valid prefix 17", code: "17401234560", valid: true},
		// Check digit normalization edge cases. When 11 - (sum mod 11)
		// yields 10 the check digit is 0, and when it yields 11 it is 1
		// (equivalent to (11 - sum mod 11) mod 10). SUNAT's public PVS
		// documentation confirms modulus 11, while these tests pin the exact
		// edge-case behavior used by the implementation.
		{name: "check digit from 10 is 0", code: "20000000010", valid: true},
		{name: "check digit from 11 is 1", code: "20000000061", valid: true},
		{name: "check digit from 11 is not 0", code: "20000000060"},
		// Presence is enforced at invoice level, not here.
		{name: "empty code", code: "", valid: true},
		// The check digit of 20547825443 should be 7, not 3.
		{name: "bad check digit", code: "20547825443"},
		{name: "undocumented type prefix (16)", code: "16401234563"},
		// 30401234566 passes the mod-11 check, but 30 is not a RUC type
		// prefix assigned by SUNAT.
		{name: "invalid type prefix", code: "30401234566"},
		{name: "too short", code: "2013131295"},
		{name: "too long", code: "201313129555"},
		{name: "non-numeric", code: "2013131295A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: pe.CountryCode, Code: tt.code}
			err := rules.Validate(tID)
			if tt.valid {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "invalid Peruvian RUC")
			}
		})
	}

	t.Run("foreign identity ignored", func(t *testing.T) {
		tID := &tax.Identity{Country: "CO", Code: "20547825443"}
		assert.NoError(t, rules.Validate(tID))
	})
}
