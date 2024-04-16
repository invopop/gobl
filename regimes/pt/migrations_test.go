package pt_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxRateMigration(t *testing.T) {
	// Valid old rate
	inv := validInvoice()
	inv.Lines[0].Taxes[0].Rate = "exempt+outlay"

	err := inv.Calculate()
	require.NoError(t, err)

	t0 := inv.Lines[0].Taxes[0]
	assert.Equal(t, tax.RateExempt, t0.Rate)
	assert.Equal(t, tax.ExtValue("M01"), t0.Ext["pt-exemption-code"])

	// Invalid old rate
	inv = validInvoice()
	inv.Lines[0].Taxes[0].Rate = "exempt+invalid"

	err = inv.Calculate()
	require.Error(t, err)
	assert.Contains(t, "invalid tax rate", err.Error())

	// Valid new rate
	inv = validInvoice()
	inv.Lines[0].Taxes[0].Rate = "exempt"
	inv.Lines[0].Taxes[0].Ext = tax.Extensions{"pt-exemption-code": "M02"}

	err = inv.Calculate()
	require.NoError(t, err)

	t0 = inv.Lines[0].Taxes[0]
	assert.Equal(t, tax.RateExempt, t0.Rate)
	assert.Equal(t, tax.ExtValue("M02"), t0.Ext["pt-exemption-code"])
}

func TestTaxZoneMigration(t *testing.T) {
	testCases := []struct {
		zone string
		tags []cbc.Key
	}{
		{
			zone: "20",
			tags: []cbc.Key{"azores"},
		},
		{
			zone: "30",
			tags: []cbc.Key{"madeira"},
		},
		{
			zone: "",
			tags: nil,
		},
	}

	for _, tc := range testCases {
		t.Run("Tax Zone Migration "+tc.zone, func(t *testing.T) {
			inv := validInvoice()
			inv.Supplier.TaxID.Zone = l10n.Code(tc.zone)

			err := inv.Calculate()
			require.NoError(t, err)

			assert.Empty(t, inv.Supplier.TaxID.Zone)

			var tags []cbc.Key
			if inv.Tax != nil {
				tags = inv.Tax.Tags
			}
			assert.Equal(t, tags, tc.tags)

			err = inv.Validate()
			require.NoError(t, err)
		})
	}
}
