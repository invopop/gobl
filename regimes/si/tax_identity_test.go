package si_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "82646716",
			Expected: "82646716",
		},
		{
			Code:     "SI82646716",
			Expected: "82646716",
		},
		{
			Code:     "si 82646716",
			Expected: "82646716",
		},
		{
			Code:     " 12345679 ",
			Expected: "12345679",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "SI", Code: ts.Code}
		norm.Normalize(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{
			name: "empty",
			code: "",
		},
		{
			// Krka d.d., Novo mesto's publicly listed VAT number.
			name: "valid",
			code: "82646716",
		},
		{
			name: "valid 2",
			code: "12345679",
		},
		{
			// Weighted sum leaves a remainder of 1, so the check digit is 0.
			name: "valid with zero check digit",
			code: "10001000",
		},
		{
			name: "too short",
			code: "1234567",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "123456789",
			err:  "IDENTITY-01",
		},
		{
			name: "leading zero",
			code: "02345679",
			err:  "IDENTITY-01",
		},
		{
			name: "non numeric",
			code: "123456A9",
			err:  "IDENTITY-01",
		},
		{
			name: "non numeric check digit",
			code: "1234567A",
			err:  "IDENTITY-01",
		},
		{
			name: "checksum mismatch",
			code: "12345678",
			err:  "IDENTITY-02",
		},
		{
			// The weighted sum leaves a remainder of 0, which is never
			// assigned: no valid tax number produces it.
			name: "remainder zero",
			code: "10000100",
			err:  "IDENTITY-02",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SI", Code: tt.code}
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
