package nl_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/nl"
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
			Code:     "000099995b57",
			Expected: "000099995B57",
		},
		{
			Code:     "NL000099995b57",
			Expected: "000099995B57",
		},
		{
			Code:     " 4359 6386 R ",
			Expected: "43596386R",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "NL", Code: ts.Code}
		nl.Normalize(tID)
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
			name: "too long",
			code: "a really really long string that's way too long",
			err:  "IDENTITY-01",
		},
		{
			name: "too short",
			code: "shorty",
			err:  "IDENTITY-01",
		},
		{
			name: "valid",
			code: "000099998B57",
		},
		{
			name: "valid 2",
			code: "808661863B01",
		},
		{
			name: "not normalized",
			code: "000099995b57",
			err:  "IDENTITY-01",
		},
		{
			name: "no B",
			code: "000099998X57",
			err:  "IDENTITY-01",
		},
		{
			name: "non numbers",
			code: "000099998B5a",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid checksum",
			code: "123456789B12",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "NL", Code: tt.code}
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
