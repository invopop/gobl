package no_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/no"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Create a standard valid invoice to use as a base
func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Code:     "123TEST",
		Currency: "NOK",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "NO",
				Code:    "290883970",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "NO",
				Code:    "974760673",
			},
		},
		IssueDate: cal.MakeDate(2023, 1, 15),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(5, 0),
				Item: &org.Item{
					Name:  "Standard Item",
					Price: num.NewAmount(5000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     tax.RateStandard,
					},
				},
			},
		},
	}
}

func reverseChargeInvoice() *bill.Invoice {
	inv := validInvoice()
	inv.Tags = tax.WithTags(tax.TagReverseCharge)
	return inv
}

func simplifiedInvoice() *bill.Invoice {
	inv := validInvoice()
	inv.Tags = tax.WithTags(tax.TagSimplified)
	inv.Customer = nil // Simplified invoices typically don't require customer details
	return inv
}

func marginSchemeInvoice() *bill.Invoice {
	inv := validInvoice()
	inv.Tags = tax.WithTags(no.TagSecondHand)
	inv.Lines[0].Item.Name = "Vintage Furniture"
	inv.Lines[0].Taxes = tax.Set{
		{
			Category: "VAT",
			Rate:     tax.RateExempt,
			Ext: tax.Extensions{
				no.ExtKeyExemptionCode: "M1",
			},
		},
	}
	inv.Tax = &bill.Tax{
		Ext: tax.Extensions{
			no.ExtKeyMarginScheme:  "second-hand",
			no.ExtKeyExemptionCode: "M1",
		},
	}
	return inv
}

func booksInvoice() *bill.Invoice {
	inv := validInvoice()
	inv.Tags = tax.WithTags(no.TagBooks)
	inv.Lines[0].Item.Name = "Educational Book"
	inv.Lines[0].Taxes = tax.Set{
		{
			Category: "VAT",
			Rate:     tax.RateZero.With(no.TagBooks),
			Ext: tax.Extensions{
				no.ExtKeyExemptionCode: "E3",
			},
		},
	}
	inv.Tax = &bill.Tax{
		Ext: tax.Extensions{
			no.ExtKeyExemptionCode: "E3",
		},
	}
	return inv
}

func voecInvoice() *bill.Invoice {
	inv := validInvoice()
	inv.Tags = tax.WithTags(no.TagECommerce)
	inv.Lines[0].Item.Name = "Imported Electronics"
	inv.Lines[0].Taxes = tax.Set{
		{
			Category: "VAT",
			Rate:     tax.RateStandard,
		},
	}
	inv.Tax = &bill.Tax{
		Ext: tax.Extensions{
			no.ExtKeyVOEC:          "registered",
			no.ExtKeyExemptionCode: "V1",
		},
	}
	return inv
}

func TestInvoiceValidation(t *testing.T) {
	// Test standard invoice validation
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestScenarios(t *testing.T) {
	t.Run("reverse charge", func(t *testing.T) {
		inv := reverseChargeInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Len(t, inv.Notes, 1)
		assert.Equal(t, tax.TagReverseCharge, inv.Notes[0].Src)
		assert.Equal(t, "Reverse Charge", inv.Notes[0].Text)
	})

	t.Run("simplified", func(t *testing.T) {
		inv := simplifiedInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Len(t, inv.Notes, 1)
		assert.Equal(t, tax.TagSimplified, inv.Notes[0].Src)
		assert.Equal(t, "Simplified Invoice (for transactions below NOK 1,000)", inv.Notes[0].Text)
	})

	t.Run("margin scheme", func(t *testing.T) {
		inv := marginSchemeInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Len(t, inv.Notes, 1)
		assert.Equal(t, no.TagSecondHand, inv.Notes[0].Src)
		assert.Equal(t, "Margin Scheme - Second-hand Goods (Chapter Va MVAL)", inv.Notes[0].Text)
		assert.Equal(t, cbc.Code("second-hand"), inv.Tax.Ext[no.ExtKeyMarginScheme])
		assert.Equal(t, cbc.Code("M1"), inv.Tax.Ext[no.ExtKeyExemptionCode])
	})

	t.Run("books", func(t *testing.T) {
		inv := booksInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Len(t, inv.Notes, 1)
		assert.Equal(t, no.TagBooks, inv.Notes[0].Src)
		assert.Equal(t, "Zero-Rated Books and Periodicals (§ 6-4 MVAL)", inv.Notes[0].Text)
		assert.Equal(t, cbc.Code("E3"), inv.Tax.Ext[no.ExtKeyExemptionCode])
	})

	t.Run("voec", func(t *testing.T) {
		inv := voecInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Len(t, inv.Notes, 1)
		assert.Equal(t, no.TagECommerce, inv.Notes[0].Src)
		assert.Equal(t, "VOEC Scheme - B2C E-commerce (§ 3-30 MVAL)", inv.Notes[0].Text)
		assert.Equal(t, cbc.Code("registered"), inv.Tax.Ext[no.ExtKeyVOEC])
		assert.Equal(t, cbc.Code("V1"), inv.Tax.Ext[no.ExtKeyExemptionCode])
	})
}

func TestInvalidTaxID(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.TaxID.Code = "123456789" // Invalid TRN (doesn't pass checksum)
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid checksum for TRN")
}
