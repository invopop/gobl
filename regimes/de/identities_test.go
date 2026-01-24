package de_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/de"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestIdentityNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "11 digits with separators", input: "123 / 456.78901", expected: "123/456/78901"},
		{name: "11 digits without separators", input: "12345678901", expected: "123/456/78901"},
		{name: "10 digits with separators", input: "12 / 345.67890", expected: "12/345/67890"},
		{name: "10 digits without separators", input: "1234567890", expected: "12/345/67890"},
		{name: "mixed characters", input: "12a3b4c5d6e7f8g9h0i1", expected: "123/456/78901"},
		{name: "less than 10 digits", input: "123456789", expected: "123456789"},
		{name: "more than 11 digits", input: "1234567890123", expected: "1234567890123"},
		// NRW format (3/4/4) - should be preserved when explicitly using slashes
		{name: "NRW format explicit", input: "123/4567/8910", expected: "123/4567/8910"},
		// Non-slash separators should normalize to standard format
		{name: "dots with 11 digits", input: "123.4567.8910", expected: "123/456/78910"},
		{name: "dashes with 11 digits", input: "123-4567-8910", expected: "123/456/78910"},
		// Standard format should be preserved when already in 3/3/5
		{name: "standard format explicit", input: "123/456/78901", expected: "123/456/78901"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tax.RegimeDefFor("DE")
			id := &org.Identity{
				Key:  de.IdentityKeyTaxNumber,
				Code: cbc.Code(tt.input),
			}
			r.NormalizeObject(id)
			assert.Equal(t, tt.expected, id.Code.String())
		})
	}
}

func TestTaxNumberValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid 10 digits (2/3/5)", code: "12/345/67890"},
		{name: "valid 11 digits (3/3/5)", code: "123/456/78901"},
		{name: "valid 11 digits NRW (3/4/4)", code: "123/4567/8910"},

		// Invalid formats
		{name: "too short", code: "12/345/678", err: "code: must be in a valid format."},
		{name: "too long", code: "1234/567/89012", err: "code: must be in a valid format."},
		{name: "non-numeric", code: "12/3AB/67890", err: "code: must be in a valid format."},
		{name: "invalid separator", code: "12-345-67890", err: "code: must be in a valid format."},
		{name: "wrong NRW middle length", code: "123/456/8910", err: "code: must be in a valid format."},
		{name: "wrong NRW last length", code: "123/4567/89101", err: "code: must be in a valid format."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: de.IdentityKeyTaxNumber, Code: tt.code}
			err := de.Validate(id)
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
