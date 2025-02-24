package adecf_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/adecf"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Addons:   tax.WithAddons(adecf.V1),
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
		},
		IssueDate: cal.MakeDate(2024, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateStandard,
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(13, 0),
				Item: &org.Item{
					Name:  "Test Item 2",
					Price: num.MakeAmount(1300, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Ext: tax.Extensions{
							adecf.ExtKeyExempt: "N2",
						},
					},
				},
			},
		},
	}
	return i
}

func TestInvoiceValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestSupplierValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Supplier.TaxID = &tax.Identity{
		Country: "IT",
		Code:    "RSSGNN60R30H501U",
	}
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "code: contains invalid characters")
}

func TestInvoiceLineTaxes(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Lines = append(inv.Lines, &bill.Line{
		Quantity: num.MakeAmount(10, 0),
		Item: &org.Item{
			Name:  "Test Item",
			Price: num.MakeAmount(10000, 2),
		},
		// No taxes!
	})
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.EqualError(t, err, "lines: (1: (taxes: missing category VAT.).).")
}

func TestTaxValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Tax = nil
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.EqualError(t, err, "tax: cannot be blank")
}
