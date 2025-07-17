package es_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
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
	tests := []struct {
		Code     cbc.Code
		Expected string
	}{
		// *** EMPTY ***
		{
			Code:     "",
			Expected: "",
		},
		// *** NATIONAL ***
		{
			Code:     "93471790C",
			Expected: "",
		},
		{
			Code:     "43596386R",
			Expected: "",
		},
		{
			Code:     "00000010X",
			Expected: "",
		},
		{
			Code:     "93471790A",
			Expected: "invalid check digit",
		},
		{
			Code:     "00000000A",
			Expected: "invalid format",
		},
		{
			Code:     "0111111C",
			Expected: "invalid format",
		},
		// *** FOREIGN ***
		{
			Code:     "X5102754C",
			Expected: "",
		},
		{
			Code:     "Z8327649K",
			Expected: "",
		},
		{
			Code:     "Y4174455S",
			Expected: "",
		},
		{
			Code:     "X5102755C",
			Expected: "invalid check digit",
		},
		{
			Code:     "X111111C",
			Expected: "invalid format",
		},
		// **** Org ****
		{
			Code:     "A58818501",
			Expected: "",
		},
		{
			Code:     "B65410011",
			Expected: "",
		},
		{
			Code:     "V7565938C",
			Expected: "",
		},
		{
			Code:     "V75659383",
			Expected: "",
		},
		{
			Code:     "F0605378I",
			Expected: "",
		},
		{
			Code:     "Q2238877A",
			Expected: "",
		},
		{
			Code:     "D40022956",
			Expected: "",
		},
		{
			Code:     "A5881850B",
			Expected: "invalid check digit",
		},
		{
			Code:     "B65410010",
			Expected: "invalid check digit",
		},
		{
			Code:     "V75659382",
			Expected: "invalid check digit",
		},
		{
			Code:     "V7565938B",
			Expected: "invalid check digit",
		},
		{
			Code:     "F06053787",
			Expected: "invalid check digit",
		},
		{
			Code:     "Q22388770",
			Expected: "invalid check digit",
		},
		{
			Code:     "D4002295J",
			Expected: "invalid check digit",
		},
		{
			Code:     "00000001A",
			Expected: "invalid check digit",
		},
		{
			Code:     "B0111111",
			Expected: "invalid format",
		},
		// *** Other ***
		{
			Code:     "K9514336H",
			Expected: "",
		},
		{
			Code:     "K95143363",
			Expected: "invalid check digit",
		},
		{
			Code:     "X111111C",
			Expected: "invalid format",
		},
	}
	r := es.New()
	for _, ts := range tests {
		t.Run(string(ts.Code), func(t *testing.T) {
			tID := &tax.Identity{Country: "ES", Code: ts.Code}
			err := r.ValidateObject(tID)
			if ts.Expected == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "code: "+ts.Expected)
			}
		})
	}
}
