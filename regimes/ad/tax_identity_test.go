package ad_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	_ "github.com/invopop/gobl/regimes/ad"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityValidation(t *testing.T) {
	tests := []struct {
		name string
		code string
		err  string
	}{
		{name: "resident individual (F)", code: "F123456A"},
		{name: "non-resident individual (E)", code: "E123456B"},
		{name: "joint-stock company (A)", code: "A059888N"},
		{name: "limited liability company (L)", code: "L132950X"},
		{name: "other entity (U)", code: "U132950X"},
		{name: "too short", code: "L12345X", err: "IDENTITY-01"},
		{name: "too long", code: "L1234567XX", err: "IDENTITY-01"},
		{name: "starts with digit", code: "1123456A", err: "IDENTITY-01"},
		{name: "ends with digit", code: "L1234561", err: "IDENTITY-01"},
		{name: "with hyphens", code: "L-132950-X", err: "IDENTITY-01"},
		{name: "empty code", code: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{
				Country: "AD",
				Code:    cbc.Code(tt.code), // ← explicit cast string → cbc.Code
			}
			err := rules.Validate(id)
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

func TestTaxIdentityNormalisation(t *testing.T) {
	r := tax.RegimeDefFor("AD")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "hyphens stripped", input: "L-132950-X", expected: "L132950X"},
		{name: "lowercase", input: "l132950x", expected: "L132950X"},
		{name: "AD prefix stripped", input: "ADL132950X", expected: "L132950X"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{
				Country: "AD",
				Code:    cbc.Code(tt.input),
			}
			r.NormalizeObject(id)
			assert.Equal(t, tt.expected, id.Code.String())
		})
	}
}
