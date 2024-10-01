package de_test

import (
	"fmt"
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
		{name: "valid 10 digits", code: "12/345/67890"},
		{name: "valid 11 digits", code: "123/456/78901"},

		// Invalid formats
		{name: "too short", code: "12/345/678", err: "invalid length"},
		{name: "too long", code: "1234/567/89012", err: "invalid length"},
		{name: "non-numeric", code: "12/3AB/67890", err: "should only contain digits"},
		{name: "invalid separator", code: "12-345-67890", err: "invalid format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: de.IdentityKeyTaxNumber, Code: tt.code}
			fmt.Println(id.Code.String())
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
