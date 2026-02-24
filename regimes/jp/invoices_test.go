package jp_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/jp"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceSupplierValidation(t *testing.T) {
	t.Run("valid with tax ID", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "supplier: (tax_id: cannot be blank.).")
	})

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "supplier: (tax_id: (code: cannot be blank.).).")
	})
}

func TestInvoiceSimplified(t *testing.T) {
	t.Run("simplified - no supplier tax ID code", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags("simplified")
		inv.Supplier.TaxID.Code = ""
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("simplified - no supplier tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags("simplified")
		inv.Supplier.TaxID = nil
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})
}

func TestInvoiceWithRegistrationNumber(t *testing.T) {
	t.Run("valid with registration number", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  jp.IdentityKeyRegistrationNumber,
				Code: "T5010401067252",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("normalized during calculation", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  jp.IdentityKeyRegistrationNumber,
				Code: "t5010401067252",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "T5010401067252", inv.Supplier.Identities[0].Code.String())
		assert.NoError(t, inv.Validate())
	})

	t.Run("invalid registration number format", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  jp.IdentityKeyRegistrationNumber,
				Code: "INVALID",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "code")
	})

	t.Run("missing T prefix", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  jp.IdentityKeyRegistrationNumber,
				Code: "5010401067252",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "code")
	})
}

func TestInvoiceWithWithholdingTax(t *testing.T) {
	t.Run("invoice with WHT", func(t *testing.T) {
		inv := &bill.Invoice{
			Regime: tax.WithRegime("JP"),
			Series: "TEST",
			Code:   "0002",
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: "JP",
					Code:    "5010401067252",
				},
			},
			Customer: &org.Party{
				Name: "Test Customer",
				TaxID: &tax.Identity{
					Country: "JP",
					Code:    "1130001011420",
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Consulting Service",
						Price: num.NewAmount(500000, 0),
						Unit:  org.UnitPackage,
					},
					Taxes: tax.Set{
						{
							Category: "VAT",
							Rate:     "general",
						},
						{
							Category: jp.TaxCategoryWHT,
							Rate:     jp.TaxRatePro,
						},
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())

		// Verify tax totals include both VAT and retained WHT
		require.NotNil(t, inv.Totals)
		require.NotNil(t, inv.Totals.Taxes)
		assert.Len(t, inv.Totals.Taxes.Categories, 2)

		// Check consumption tax category
		vatCat := inv.Totals.Taxes.Categories[0]
		assert.Equal(t, tax.CategoryVAT, vatCat.Code)
		assert.False(t, vatCat.Retained)
		assert.Equal(t, "50000", vatCat.Amount.String())

		// Check withholding tax category
		whtCat := inv.Totals.Taxes.Categories[1]
		assert.Equal(t, jp.TaxCategoryWHT, whtCat.Code)
		assert.True(t, whtCat.Retained)
		assert.Equal(t, "51050", whtCat.Amount.String())

		// Retained total should be set
		require.NotNil(t, inv.Totals.Taxes.Retained)
		assert.Equal(t, "51050", inv.Totals.Taxes.Retained.String())
	})
}
