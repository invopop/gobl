package de_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/de"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("DE"),
		Series: "TEST",
		Code:   "0002",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "111111125",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "282741168",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("normal invoice", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-DE-BILL-INVOICE-01]")
	})

	t.Run("simplified invoice", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags("simplified")
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("regular invoice - only tax number", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  de.IdentityKeyTaxNumber,
				Code: "92/345/67894",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("regular invoice - only tax number nil tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  de.IdentityKeyTaxNumber,
				Code: "92/345/67894",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}
