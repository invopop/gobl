package pt_test

import (
	"testing"

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
