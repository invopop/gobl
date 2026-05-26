package ad_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	adregime "github.com/invopop/gobl/regimes/ad"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestOrgIdentityNormalisation(t *testing.T) {
	r := tax.RegimeDefFor("AD") // ← get the regime, use its NormalizeObject

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "hyphens stripped", input: "L-132950-X", expected: "L132950X"},
		{name: "lowercased", input: "l132950x", expected: "L132950X"},
		{name: "AD prefix stripped", input: "ADL132950X", expected: "L132950X"},
		{name: "NRT label prefix stripped", input: "NRT L132950X", expected: "L132950X"},
		{name: "already clean", input: "L132950X", expected: "L132950X"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Type: adregime.IdentityTypeNRT,
				Code: cbc.Code(tt.input),
			}
			r.NormalizeObject(id) // ← correct call, not rules.Normalize
			assert.Equal(t, tt.expected, id.Code.String())
		})
	}
}

func TestOrgIdentityValidation(t *testing.T) {
	tests := []struct {
		name string
		code string
		err  string
	}{
		{name: "valid SL company", code: "L132950X"},
		{name: "valid SA company", code: "A059888N"},
		{name: "valid resident individual", code: "F123456A"},
		{name: "too short", code: "L12345X", err: "IDENTITY-01"},
		{name: "too long", code: "L12345678X", err: "IDENTITY-01"},
		{name: "starts with digit", code: "1123456A", err: "IDENTITY-01"},
		{name: "ends with digit", code: "L1234561", err: "IDENTITY-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Type: adregime.IdentityTypeNRT,
				Code: cbc.Code(tt.code), // ← explicit cast
			}
			// ← must pass RegimeContext so orgIdentityRules fires
			err := rules.Validate(id, tax.RegimeContext("AD"))
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