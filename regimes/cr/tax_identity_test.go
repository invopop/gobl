package cr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/regimes/cr"
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
		{name: "already normalized", input: "3101123456", expected: "3101123456"},
		{name: "cedula juridica with hyphens", input: "3-101-123456", expected: "3101123456"},
		{name: "cedula fisica with hyphens", input: "1-0998-0456", expected: "109980456"},
		{name: "with country prefix", input: "CR3101123456", expected: "3101123456"},
		{name: "with spaces", input: "3 101 123456", expected: "3101123456"},
		{name: "lowercase with prefix and hyphens", input: "cr 3-101-123456", expected: "3101123456"},
		{name: "dimex with spaces", input: "1558 1234 5678", expected: "155812345678"},
		{name: "empty", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: cr.CountryCode, Code: tt.input}
			norm.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		code    cbc.Code
		invalid bool
	}{
		// Cédula física: 9 digits, no leading zero.
		{name: "cedula fisica", code: "109980456"},
		{name: "cedula fisica high province", code: "901230456"},
		{name: "cedula fisica leading zero", code: "012345678", invalid: true},
		{name: "cedula fisica with letter", code: "1A9980456", invalid: true},
		// Cédula jurídica or NITE: 10 characters, leading zeros allowed.
		{name: "cedula juridica", code: "3101123456"},
		{name: "cedula juridica public entity", code: "2100042005"},
		{name: "cedula juridica leading zero", code: "0101123456"},
		{name: "nite natural person", code: "3120456789"},
		{name: "nite corporation", code: "3130456789"},
		{name: "cedula juridica with letter", code: "3101A23456"},
		{name: "ten characters with symbol", code: "3101-23456", invalid: true},
		// DIMEX: 11 or 12 digits, no leading zero.
		{name: "dimex 11 digits", code: "15581234567"},
		{name: "dimex 12 digits", code: "155812345678"},
		{name: "dimex leading zero", code: "015581234567", invalid: true},
		{name: "dimex with letter", code: "15581234567A", invalid: true},
		// Anything else.
		{name: "too short", code: "12345678", invalid: true},
		{name: "too long", code: "1234567890123", invalid: true},
		{name: "empty code", code: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: cr.CountryCode, Code: tt.code}
			err := rules.Validate(tID)
			if tt.invalid {
				assert.ErrorContains(t, err, "CR-TAX-IDENTITY-01")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
