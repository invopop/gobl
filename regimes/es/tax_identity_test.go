package es_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	r := es.New()
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "93471790-C",
			Expected: "93471790C",
		},
		{
			Code:     " 4359 6386 R ",
			Expected: "43596386R",
		},
		{
			Code:     "Z-8327649-K",
			Expected: "Z8327649K",
		},
		{
			Code:     "ES93471790C",
			Expected: "93471790C",
		},
		{
			Code:     " ES-93 471 790-C ",
			Expected: "93471790C",
		},
		{
			Code:     "3157928M", // short, should zero pad
			Expected: "03157928M",
		},
		{
			Code:     "15S", // King Felipe VI
			Expected: "00000015S",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "ES", Code: ts.Code}
		r.NormalizeObject(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityKey(t *testing.T) {
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Key
	}{
		{
			Code:     "93471790C",
			Expected: es.TaxIdentityNational,
		},
		{
			Code:     "X5102754C",
			Expected: es.TaxIdentityForeigner,
		},
		{
			Code:     "A58818501",
			Expected: es.TaxIdentityOrg,
		},
		{
			Code:     "K9514336H",
			Expected: es.TaxIdentityOther,
		},
		{
			Code:     "XXX",
			Expected: cbc.KeyEmpty,
		},
		{
			Code:     "",
			Expected: cbc.KeyEmpty,
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "ES", Code: ts.Code}
		assert.Equal(t, ts.Expected, es.TaxIdentityKey(tID))
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	const errCode01 = "GOBL-ES-TAX-IDENTITY-01"
	tests := []struct {
		Code     cbc.Code
		Expected string
	}{
		// *** EMPTY ***
		{
			Code:     "",
			Expected: errCode01,
		},
		// *** NATIONAL ***
		{
			Code: "93471790C",
		},
		{
			Code: "43596386R",
		},
		{
			Code: "00000010X",
		},
		{
			Code:     "93471790A",
			Expected: errCode01,
		},
		{
			Code:     "00000000A",
			Expected: errCode01,
		},
		{
			Code:     "0111111C",
			Expected: errCode01,
		},
		// *** FOREIGN ***
		{
			Code: "X5102754C",
		},
		{
			Code: "Z8327649K",
		},
		{
			Code: "Y4174455S",
		},
		{
			Code:     "X5102755C",
			Expected: errCode01,
		},
		{
			Code:     "X111111C",
			Expected: errCode01,
		},
		// **** Org ****
		{
			Code: "A58818501",
		},
		{
			Code: "B65410011",
		},
		{
			Code: "V7565938C",
		},
		{
			Code: "V75659383",
		},
		{
			Code: "F0605378I",
		},
		{
			Code: "Q2238877A",
		},
		{
			Code: "D40022956",
		},
		{
			Code:     "A5881850B",
			Expected: errCode01,
		},
		{
			Code:     "B65410010",
			Expected: errCode01,
		},
		{
			Code:     "V75659382",
			Expected: errCode01,
		},
		{
			Code:     "V7565938B",
			Expected: errCode01,
		},
		{
			Code:     "F06053787",
			Expected: errCode01,
		},
		{
			Code:     "Q22388770",
			Expected: errCode01,
		},
		{
			Code:     "D4002295J",
			Expected: errCode01,
		},
		{
			Code:     "00000001A",
			Expected: errCode01,
		},
		{
			Code:     "B0111111",
			Expected: errCode01,
		},
		// *** Other ***
		{
			Code: "K9514336H",
		},
		{
			Code:     "K95143363",
			Expected: errCode01,
		},
		{
			Code:     "X111111C",
			Expected: errCode01,
		},
	}
	for _, ts := range tests {
		t.Run(string(ts.Code), func(t *testing.T) {
			tID := &tax.Identity{Country: "ES", Code: ts.Code}
			err := rules.Validate(tID)
			if ts.Expected == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, ts.Expected)
			}
		})
	}
}
