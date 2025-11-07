package fi_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/fi"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	r := fi.New()
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "031251-1011",
			Expected: "0312511011",
		},
		{
			Code:     "1111 85-22 77",
			Expected: "1111852277",
		},
		{
			Code:     "0 8 0 4 6 2 - 6 0 2 T",
			Expected: "080462602T",
		},
		{
			Code:     "150600A905P",
			Expected: "150600905P",
		},
		{
			Code:     "150600+905P",
			Expected: "150600905P",
		},
		{
			Code:     "150600C905P",
			Expected: "150600905P",
		},
		{
			Code:     "FI0 8 0 4 6 2 - 6 0 2 T",
			Expected: "080462602T",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "FI", Code: ts.Code}
		r.NormalizeObject(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     cbc.Code
		Expected error
	}{
		// *** EMPTY ***
		{
			Code:     "",
			Expected: nil,
		},
		// *** NATIONAL ***
		{
			Code:     "0312511011", // 031251-1011
			Expected: nil,
		},
		{
			Code:     "1111852277", // 111185-2277
			Expected: nil,
		},
		{
			Code:     "080462602T", // 080462-602T
			Expected: nil,
		},
		{
			Code:     "0804626022",
			Expected: fi.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "0000000000",
			Expected: fi.ErrTaxCodeInvalidDate,
		},
		{
			Code:     "0101010000",
			Expected: fi.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "081362602T",
			Expected: fi.ErrTaxCodeInvalidDate,
		},
		{
			Code:     "93471790A",
			Expected: fi.ErrTaxCodeUnknownType,
		},
		{
			Code:     "00000000A",
			Expected: fi.ErrTaxCodeUnknownType,
		},
		{
			Code:     "0111111C",
			Expected: fi.ErrTaxCodeUnknownType,
		},
		// *** FOREIGN ***
		{
			Code:     "123456789012",
			Expected: nil,
		},
		{
			Code:     "987654321098",
			Expected: nil,
		},
		{
			Code:     "112233445566",
			Expected: nil,
		},
		// // **** Org ****
		{
			Code:     "50774741", // 5077474-1
			Expected: nil,
		},
		{
			Code:     "12345671", // 1234567-1
			Expected: nil,
		},
		{
			Code:     "07654322",
			Expected: nil,
		},
		{
			Code:     "55555556",
			Expected: nil,
		},
		{
			Code:     "00000000",
			Expected: nil,
		},
		{
			Code:     "99999992",
			Expected: nil,
		},
		{
			Code:     "33145556",
			Expected: nil,
		},
		{
			Code:     "33145557",
			Expected: fi.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "33145558",
			Expected: fi.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "33145559",
			Expected: fi.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "33145551",
			Expected: fi.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "33145552",
			Expected: fi.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "33145553",
			Expected: fi.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "33145554",
			Expected: fi.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "33145555",
			Expected: fi.ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "X111111C",
			Expected: fi.ErrTaxCodeUnknownType,
		},
		{
			Code:     "Y111111C",
			Expected: fi.ErrTaxCodeUnknownType,
		},
		{
			Code:     "0",
			Expected: fi.ErrTaxCodeUnknownType,
		},
		{
			Code:     "000000000",
			Expected: fi.ErrTaxCodeUnknownType,
		},
		{
			Code:     "123456789",
			Expected: fi.ErrTaxCodeUnknownType,
		},
	}
	r := fi.New()
	for _, ts := range tests {
		t.Run(string(ts.Code), func(t *testing.T) {
			tID := &tax.Identity{Country: "FI", Code: ts.Code}
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
