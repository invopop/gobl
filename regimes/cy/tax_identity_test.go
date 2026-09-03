package cy_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "structurally valid historic format", code: "12345678A"},
		{name: "structurally valid new format", code: "60000000A"},
		{name: "empty code", code: ""},
		{name: "too short", code: "1234567A", err: "IDENTITY-01"},
		{name: "too long", code: "123456789A", err: "IDENTITY-01"},
		{name: "missing final letter", code: "123456789", err: "IDENTITY-01"},
		{name: "letter in numeric section", code: "1234A6789", err: "IDENTITY-01"},
		{name: "lowercase letter", code: "12345678a", err: "IDENTITY-01"},
		{name: "not normalized", code: "1234-5678-A", err: "IDENTITY-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CY", Code: tt.code}
			err := rules.Validate(tID)
			if tt.err == "" {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.err)
			}
		})
	}
}

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		want cbc.Code
	}{
		{name: "no change", code: "12345678A", want: "12345678A"},
		{name: "normalize prefix, separators, and case", code: "cy 1234-5678.a", want: "12345678A"},
		{name: "empty", code: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CY", Code: tt.code}
			norm.Normalize(tID)
			assert.Equal(t, tt.want, tID.Code)
		})
	}
}
