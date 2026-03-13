package ch_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/ch"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	var tID *tax.Identity
	assert.NotPanics(t, func() {
		ch.Normalize(tID)
	}, "nil tax identity")

	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "CHE-901.458.652",
			Expected: "E901458652",
		},
		{
			Code:     "E-800.134.536",
			Expected: "E800134536",
		},
		{
			Code:     "CHE-284.156.502-MWST",
			Expected: "E284156502",
		},
		{
			Code:     "CHE-284.156.502-TVA",
			Expected: "E284156502",
		},
		{
			Code:     "CHE-284.156.502 IVA",
			Expected: "E284156502",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "CH", Code: ts.Code}
		ch.Normalize(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "E100416306"},
		{name: "good 2", code: "E284156502"},
		{name: "good 3", code: "E432825998"},
		{name: "good 4", code: "E115400550"},
		{name: "good 5", code: "E424414541"},
		{name: "good 6", code: "E123456788"},
		{
			name: "zeros",
			code: "E000000001",
			err:  "IDENTITY-01",
		},
		{
			name: "bad mid length",
			code: "E12345678910",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "E1234567890123",
			err:  "IDENTITY-01",
		},
		{
			name: "too short",
			code: "E123456",
			err:  "IDENTITY-01",
		},
		{
			name: "not normalized",
			code: "E-385.16.405",
			err:  "IDENTITY-01",
		},
		{
			name: "bad checksum",
			code: "E116276851",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid code",
			code: "E100426306",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CH", Code: tt.code}
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
