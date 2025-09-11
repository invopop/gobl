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

	t.Run("exempt with addon added later", func(t *testing.T) {
		// This tests covers a typical use-case whereby a document is
		// created without addons but with the extensions to be used later.
		inv := validInvoice()
		inv.Addons = tax.Addons{}
		tc := inv.Lines[0].Taxes[0]
		tc.Key = ""
		tc.Rate = tax.KeyExempt
		tc.Ext = tc.Ext.Set(saft.ExtKeyExemption, "M40")

		require.NoError(t, inv.Calculate())

		assert.Empty(t, tc.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "exempt", tc.Key.String())
		assert.Equal(t, "M40", tc.Ext[saft.ExtKeyExemption].String())

		// Add the addon and re-calculate
		inv.Addons = tax.WithAddons(saft.V1)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "ISE", tc.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "reverse-charge", tc.Key.String())
		assert.Equal(t, "M40", tc.Ext[saft.ExtKeyExemption].String())
	})

	t.Run("exempt", func(t *testing.T) {
		inv := validInvoice()
		tc := inv.Lines[0].Taxes[0]
		tc.Key = tax.KeyExempt
		tc.Rate = ""
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "ISE", tc.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M07", tc.Ext[saft.ExtKeyExemption].String())

		// Allow override as this is "exempt"
		tc.Ext[saft.ExtKeyExemption] = "M04"
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "export", tc.Key.String())
		assert.Equal(t, "M04", tc.Ext[saft.ExtKeyExemption].String())

		// Do not allow override from "export" back to "exempt", but
		// force the code back to default "M05"
		tc.Ext[saft.ExtKeyExemption] = "M01"
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "export", tc.Key.String())
		assert.Equal(t, "M05", tc.Ext[saft.ExtKeyExemption].String())
	})
}
