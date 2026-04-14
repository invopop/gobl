package hu_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("HU"),
		Series: "TEST",
		Code:   "0001",
		Supplier: &org.Party{
			Name: "Példa Kft.",
			TaxID: &tax.Identity{
				Country: "HU",
				Code:    "13895459",
			},
		},
		Customer: &org.Party{
			Name: "Vevő Zrt.",
			TaxID: &tax.Identity{
				Country: "HU",
				Code:    "10537914",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Szoftverfejlesztés",
					Price: num.NewAmount(50000, 2),
					Unit:  org.UnitHour,
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

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-HU-BILL-INVOICE-01]")
	})

	t.Run("invalid supplier tax ID code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = "13895450"
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-HU-TAX-IDENTITY-01]")
	})
}
