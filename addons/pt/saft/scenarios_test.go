package saft_test

import (
	"testing"

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
		inv.Lines[0].Taxes[0].Rate = tax.RateExempt
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "ISE", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyTaxRate].String())
		assert.ErrorContains(t, inv.Validate(), "lines: (0: (taxes: (0: (ext: (pt-saft-exemption: required.).).).).)")

		inv.Lines[0].Taxes[0].Ext[saft.ExtKeyExemption] = "M04"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})
}

func validReceipt() *bill.Receipt {
	return &bill.Receipt{
		Type: bill.ReceiptTypePayment,
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "PT",
			},
		},
		Code:      "123",
		IssueDate: cal.MakeDate(2024, 3, 10),
		Lines: []*bill.ReceiptLine{
			{
				Debit: num.NewAmount(100, 2),
			},
		},
		Method: &pay.Instructions{
			Key: "credit-transfer",
		},
	}
}

func TestReceipt(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		r := validReceipt()
		r.SetAddons(saft.V1)
		require.NoError(t, r.Calculate())
		assert.Equal(t, "RG", r.Ext[saft.ExtKeyReceiptType].String())
		assert.NoError(t, r.Validate())
	})

	t.Run("VAT cash", func(t *testing.T) {
		r := validReceipt()
		r.SetAddons(saft.V1)
		require.NoError(t, r.Calculate())
		assert.Equal(t, "RC", r.Ext[saft.ExtKeyReceiptType].String())
		assert.NoError(t, r.Validate())
	})
}
