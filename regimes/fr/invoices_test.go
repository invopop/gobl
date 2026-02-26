package fr_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series: "TEST",
		Code:   "0002",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44732829320",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44391838042",
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
	t.Run("valid with tax ID", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("missing tax ID code requires SIREN or SIRET", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "identities: missing type")
	})

	t.Run("valid with SIREN identity instead of tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Regime = tax.WithRegime("FR")
		inv.Supplier.TaxID = nil
		inv.Supplier.Identities = []*org.Identity{
			{
				Type: fr.IdentityTypeSIREN,
				Code: "732829320",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("valid with SIRET identity instead of tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Regime = tax.WithRegime("FR")
		inv.Supplier.TaxID = nil
		inv.Supplier.Identities = []*org.Identity{
			{
				Type: fr.IdentityTypeSIRET,
				Code: "73282932000015",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("missing both tax ID and identity", func(t *testing.T) {
		inv := validInvoice()
		inv.Regime = tax.WithRegime("FR")
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "tax_id: cannot be blank")
	})
}
