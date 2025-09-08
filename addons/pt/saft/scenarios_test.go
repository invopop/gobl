package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoice(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "FT", inv.Tax.Ext[saft.ExtKeyInvoiceType].String())
		assert.Equal(t, "NOR", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyTaxRate].String())

		assert.NoError(t, inv.Validate())
	})

	t.Run("prepaid", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Percent:     num.NewPercentage(1, 0),
					Description: "prepaid",
				},
			},
		}
		inv.Series = "FR SERIES-A"
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "FR", inv.Tax.Ext[saft.ExtKeyInvoiceType].String())
		assert.NoError(t, inv.Validate())
	})

	t.Run("reduced", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Taxes[0].Rate = tax.RateReduced
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "RED", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyTaxRate].String())
		assert.NoError(t, inv.Validate())
	})

	t.Run("intermediate", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Taxes[0].Rate = tax.RateIntermediate
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "INT", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyTaxRate].String())
		assert.NoError(t, inv.Validate())
	})

	t.Run("exempt", func(t *testing.T) {
		inv := validInvoice()
		tc := inv.Lines[0].Taxes[0]
		tc.Key = tax.KeyExempt
		tc.Rate = ""
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "ISE", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M07", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyExemption].String())

		// Allow override
		inv.Lines[0].Taxes[0].Ext[saft.ExtKeyExemption] = "M04"
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "M07", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyExemption].String())

		// Allow override
		inv.Lines[0].Taxes[0].Ext[saft.ExtKeyExemption] = "M01"
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "M01", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyExemption].String())
	})
}
