package es_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/es"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     string
		Expected string
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
	}
	for _, ts := range tests {
		tID := &org.TaxIdentity{Country: l10n.ES, Code: ts.Code}
		err := es.NormalizeTaxIdentity(tID)
		assert.NoError(t, err)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     string
		Expected error
	}{
		// *** NATIONAL ***
		{
			Code:     "93471790C",
			Expected: nil,
		},
		{
			Code:     "43596386R",
			Expected: nil,
		},
		{
			Code:     "00000010X",
			Expected: nil,
		},
		{
			Code:     "93471790A",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "00000000A",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "0111111C",
			Expected: es.ErrTaxCodeUnknownType,
		},
		// *** FOREIGN ***
		{
			Code:     "X5102754C",
			Expected: nil,
		},
		{
			Code:     "Z8327649K",
			Expected: nil,
		},
		{
			Code:     "Y4174455S",
			Expected: nil,
		},
		{
			Code:     "X5102755C",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "X111111C",
			Expected: es.ErrTaxCodeUnknownType,
		},
		// **** Org ****
		{
			Code:     "A58818501",
			Expected: nil,
		},
		{
			Code:     "B65410011",
			Expected: nil,
		},
		{
			Code:     "V7565938C",
			Expected: nil,
		},
		{
			Code:     "V75659383",
			Expected: nil,
		},
		{
			Code:     "F0605378I",
			Expected: nil,
		},
		{
			Code:     "Q2238877A",
			Expected: nil,
		},
		{
			Code:     "D40022956",
			Expected: nil,
		},
		{
			Code:     "A5881850B",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "B65410010",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "V75659382",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "V7565938B",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "F06053787",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "Q22388770",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "D4002295J",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "00000000A",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "B0111111",
			Expected: es.ErrTaxCodeUnknownType,
		},
		// *** Other ***
		{
			Code:     "K9514336H",
			Expected: nil,
		},
		{
			Code:     "K95143363",
			Expected: es.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "X111111C",
			Expected: es.ErrTaxCodeUnknownType,
		},
	}
	for _, ts := range tests {
		t.Run(ts.Code, func(t *testing.T) {
			tID := &org.TaxIdentity{Country: l10n.ES, Code: ts.Code}
			err := es.ValidateTaxIdentity(tID)
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
