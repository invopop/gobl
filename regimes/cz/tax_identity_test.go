package cz_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/cz"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{Code: "00177041", Expected: "00177041"},
		{Code: "CZ00177041", Expected: "00177041"},
		{Code: "001 770 41", Expected: "00177041"},
		{Code: "001-770-41", Expected: "00177041"},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "CZ", Code: ts.Code}
		cz.Normalize(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid legal entity codes (8 digits)
		{name: "valid legal entity - Škoda Auto", code: "00177041"},
		{name: "valid legal entity - ČEZ", code: "45274649"},
		{name: "valid legal entity", code: "25596641"},
		{name: "valid legal entity - T-Mobile CZ", code: "64949681"},
		{name: "valid legal entity - Komerční banka", code: "45317054"},
		// Valid individual code (10 digits, divisible by 11)
		{name: "valid individual", code: "7103192745"},
		// Valid legal entity with check digit 0 (expected == 10 branch)
		{name: "valid legal entity check digit 0", code: "00200000"},
		// Valid 9-digit codes
		{name: "valid 9-digit special", code: "612345679"},
		{name: "valid 9-digit individual", code: "710319274"},
		// Invalid format
		{
			name: "too short",
			code: "0017704",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "00177041234",
			err:  "IDENTITY-01",
		},
		{
			name: "contains letters",
			code: "0017704A",
			err:  "IDENTITY-01",
		},
		// Bad checksum - legal entity
		{
			name: "bad checksum legal entity",
			code: "00177042",
			err:  "IDENTITY-01",
		},
		{
			name: "bad checksum legal entity 2",
			code: "45274640",
			err:  "IDENTITY-01",
		},
		// Bad checksum - individual (10 digits, not divisible by 11)
		{
			name: "bad checksum individual",
			code: "7103192746",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CZ", Code: tt.code}
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
