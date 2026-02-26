package dk_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/dk"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("DK"),
		Series: "TEST",
		Code:   "0002",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "DK",
				Code:    "12345674",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "DK",
				Code:    "12345674",
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

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "supplier: (identities: missing type 'CVR'; tax_id: (code: cannot be blank.).).")
	})

	t.Run("supplier with CVR identity instead of tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		inv.Supplier.Identities = []*org.Identity{
			{
				Type: dk.IdentityTypeCVR,
				Code: "12345674",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("supplier with nil tax ID and CVR identity", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		inv.Supplier.Identities = []*org.Identity{
			{
				Type: dk.IdentityTypeCVR,
				Code: "12345674",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})
}
