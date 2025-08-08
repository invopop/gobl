package it_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxComboNormalization(t *testing.T) {
	t.Run("replace natura with exempt extenions", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0] = &tax.Combo{
			Category: "VAT",
			Percent:  nil, // exempt
			Ext: tax.Extensions{
				"it-sdi-nature": "N1",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "N1", inv.Lines[0].Taxes[0].Ext[sdi.ExtKeyExempt].String())
		assert.NotContains(t, inv.Lines[0].Taxes[0].Ext, "it-sdi-nature")
	})

	t.Run("replace retained tax extenion", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0] = &tax.Combo{
			Category: "IRPEF",
			Percent:  num.NewPercentage(8, 3),
			Ext: tax.Extensions{
				"it-sdi-retained-tax": "A",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "A", inv.Lines[0].Taxes[0].Ext[sdi.ExtKeyRetained].String())
		assert.NotContains(t, inv.Lines[0].Taxes[0].Ext, "it-sdi-retained-tax")
	})

}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Regime:   tax.WithRegime("IT"),
		Code:     "123TEST",
		Currency: "EUR",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Type: bill.InvoiceTypeStandard,
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "IT",
				Code:    "12345678903",
			},
			Addresses: []*org.Address{
				{
					Street:   "Via di Test",
					Code:     "12345",
					Locality: "Rome",
					Country:  "IT",
					Number:   "3",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "IT",
				Code:    "13029381004",
			},
			Addresses: []*org.Address{
				{
					Street:   "Piazza di Test",
					Code:     "38342",
					Locality: "Venezia",
					Country:  "IT",
					Number:   "1",
				},
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Key:      "standard",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}
	return i
}
