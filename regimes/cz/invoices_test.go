package cz_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/cz"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("CZ"),
		Series: "TEST",
		Code:   "0002",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "CZ",
				Code:    "00177041",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "CZ",
				Code:    "45274649",
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
						Rate:     "standard",
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice with tax ID", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("valid invoice with ICO instead of tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  cz.IdentityKeyICO,
				Code: "00177041",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("missing both tax ID and ICO", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "supplier")
	})

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "identities")
	})

	t.Run("simplified invoice without supplier tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(tax.TagSimplified)
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})
}
