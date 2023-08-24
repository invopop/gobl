package pt_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxRateMigration(t *testing.T) {
	inv := validInvoice()
	inv.Lines[0].Taxes[0].Rate = "exempt+outlay"

	err := inv.Calculate()
	require.NoError(t, err)

	t0 := inv.Lines[0].Taxes[0]
	assert.Equal(t, common.TaxRateExempt, t0.Rate)
	assert.Equal(t, cbc.Code("M01"), t0.Ext["pt-exemption-code"])

	inv = validInvoice()
	inv.Lines[0].Taxes[0].Rate = "exempt+invalid"

	err = inv.Calculate()
	require.Error(t, err)
	assert.Contains(t, "invalid tax rate", err.Error())
}
