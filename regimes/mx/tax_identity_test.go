package mx_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityNormalization(t *testing.T) {
	r := mx.New()
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "GHI70123123Z",
			Expected: "GHI70123123Z",
		},
		{
			Code:     " GHI 701231 23Z ",
			Expected: "GHI70123123Z",
		},
		{
			Code:     "GHI-701231-23Z",
			Expected: "GHI70123123Z",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: l10n.MX, Code: ts.Code}
		err := r.CalculateObject(tID)
		assert.NoError(t, err)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityValidation(t *testing.T) {
	tests := []struct {
		Code     cbc.Code
		Expected error
	}{
		{
			Code:     "",
			Expected: nil,
		},
		{
			Code:     mx.TaxIdentityCodeForeign,
			Expected: nil,
		},
		{
			Code:     "MNOP8201019HJ",
			Expected: nil,
		},
		{
			Code:     "UVWX610715JKL",
			Expected: nil,
		},
		{
			Code:     "STU760612MN1",
			Expected: nil,
		},
		{
			Code:     "STU760612MN",
			Expected: tax.ErrIdentityCodeInvalid,
		},
		{
			Code:     "XXXX",
			Expected: tax.ErrIdentityCodeInvalid,
		},
	}
	r := mx.New()
	for _, ts := range tests {
		t.Run(string(ts.Code), func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.MX, Code: ts.Code}
			err := r.ValidateObject(tID)
			if ts.Expected == nil {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), ts.Expected.Error())
				}
			}
		})
	}
}

func TestTaxIdentityDetermineType(t *testing.T) {
	tests := []struct {
		Code cbc.Code
		Type cbc.Key
	}{
		{
			Code: "",
			Type: cbc.KeyEmpty,
		},
		{
			Code: mx.TaxIdentityCodeForeign,
			Type: mx.TaxIdentityTypePerson,
		},
		{
			Code: "MNOP8201019HJ",
			Type: mx.TaxIdentityTypePerson,
		},
		{
			Code: "ABC830720XYZ",
			Type: mx.TaxIdentityTypeCompany,
		},
		{
			Code: "XXXX",
			Type: cbc.KeyEmpty,
		},
	}
	for _, ts := range tests {
		t.Run(string(ts.Code), func(t *testing.T) {
			res := mx.DetermineTaxCodeType(ts.Code)
			assert.Equal(t, ts.Type, res)
		})
	}
}
