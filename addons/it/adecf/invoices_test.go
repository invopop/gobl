package adecf_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/it/adecf"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func exampleStandardInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Regime:   tax.WithRegime("IT"),
		Addons:   tax.WithAddons(adecf.V1),
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
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item 0",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
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
					Name:  "Test Item 1",
					Price: num.MakeAmount(1000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Ext: tax.Extensions{
							adecf.ExtKeyExempt: "N4",
						},
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

func TestInvoiceValidation(t *testing.T) {
	inv := exampleStandardInvoice(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestSupplierValidation(t *testing.T) {
	inv := exampleStandardInvoice(t)
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
	t.Run("Item with no taxes", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item 2",
				Price: num.MakeAmount(10000, 2),
			},
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.EqualError(t, err, "lines: (2: (taxes: missing category VAT.).).")
	})

	t.Run("Item with no Rate and missing Ext", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item 2",
				Price: num.MakeAmount(10000, 2),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
				},
			},
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.EqualError(t, err, "lines: (2: (taxes: (0: (ext: (it-adecf-exempt: required.).).).).).")
	})

	t.Run("Item with Invalid Percentage", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item 2",
				Price: num.MakeAmount(10000, 2),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Percent:  num.NewPercentage(24, 2),
				},
			},
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.EqualError(t, err, "lines: (2: (taxes: (0: (percent: Invalid percentage.).).).).")
	})
}

func TestInvoiceTax(t *testing.T) {
	t.Run("Invalid PricesInclude", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Tax.PricesInclude = "invalid"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.EqualError(t, err, "tax: (prices_include: must be a valid value.).")
	})

	t.Run("Missing PricesInclude", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Tax.PricesInclude = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.EqualError(t, err, "tax: (prices_include: cannot be blank.).")
	})

	t.Run("Missing Tax", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.EqualError(t, err, "tax: cannot be blank.")
	})
}
