package saft_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("PT"),
		Addons: tax.WithAddons(saft.V1),
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Code:    "123456789",
				Country: "PT",
			},
			Name: "Test Supplier",
		},
		Customer: &org.Party{
			Name: "Test Customer",
		},
		Code:      "INV/1",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2023, 1, 1),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
	}
}

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
		inv.Payment = &bill.Payment{
			Advances: []*pay.Advance{
				{
					Percent:     num.NewPercentage(1, 0),
					Description: "prepaid",
				},
			},
		}
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
		inv.Lines[0].Taxes[0].Rate = tax.RateExempt
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "ISE", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyTaxRate].String())
		assert.ErrorContains(t, inv.Validate(), "lines: (0: (taxes: (0: (ext: (pt-saft-exemption: required.).).).).)")

		inv.Lines[0].Taxes[0].Ext[saft.ExtKeyExemption] = "M04"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})
}
