package sa_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validStandardInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "INV-001",
		Currency: "SAR",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "300075588700003",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "310122393500003",
			},
		},
		IssueDate: cal.MakeDate(2024, 1, 15),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(5, 0),
				Item: &org.Item{
					Name:  "Service Item",
					Price: num.NewAmount(200, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}

func validSimplifiedInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "INV-002",
		Currency: "SAR",
		Tags:     tax.WithTags(tax.TagSimplified),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "300075588700003",
			},
		},
		IssueDate: cal.MakeDate(2024, 1, 15),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Product Item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid standard invoice", func(t *testing.T) {
		i := validStandardInvoice(t)
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())
	})

	t.Run("valid simplified invoice without customer", func(t *testing.T) {
		i := validSimplifiedInvoice(t)
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())
	})

	t.Run("standard invoice missing supplier tax ID code", func(t *testing.T) {
		i := validStandardInvoice(t)
		i.Supplier.TaxID.Code = ""
		require.NoError(t, i.Calculate())
		err := i.Validate()
		assert.ErrorContains(t, err, "supplier")
		assert.ErrorContains(t, err, "tax_id")
	})

	t.Run("simplified invoice missing supplier tax ID code", func(t *testing.T) {
		i := validSimplifiedInvoice(t)
		i.Supplier.TaxID.Code = ""
		require.NoError(t, i.Calculate())
		err := i.Validate()
		assert.ErrorContains(t, err, "supplier")
		assert.ErrorContains(t, err, "tax_id")
	})

	t.Run("standard invoice missing customer", func(t *testing.T) {
		i := validStandardInvoice(t)
		i.Customer = nil
		require.NoError(t, i.Calculate())
		err := i.Validate()
		assert.ErrorContains(t, err, "customer")
	})

	t.Run("simplified invoice allows no customer", func(t *testing.T) {
		i := validSimplifiedInvoice(t)
		i.Customer = nil
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())
	})

	t.Run("reverse charge requires customer tax ID", func(t *testing.T) {
		i := validStandardInvoice(t)
		i.Tags = tax.WithTags(tax.TagReverseCharge)
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())

		// Remove customer tax ID code
		i = validStandardInvoice(t)
		i.Tags = tax.WithTags(tax.TagReverseCharge)
		i.Customer.TaxID.Code = ""
		require.NoError(t, i.Calculate())
		err := i.Validate()
		assert.ErrorContains(t, err, "customer")
		assert.ErrorContains(t, err, "tax_id")
	})

	t.Run("reverse charge without customer tax ID object", func(t *testing.T) {
		i := validStandardInvoice(t)
		i.Tags = tax.WithTags(tax.TagReverseCharge)
		i.Customer.TaxID = nil
		require.NoError(t, i.Calculate())
		err := i.Validate()
		assert.ErrorContains(t, err, "customer")
		assert.ErrorContains(t, err, "tax_id")
	})

	t.Run("standard invoice without reverse charge does not require customer tax ID", func(t *testing.T) {
		i := validStandardInvoice(t)
		i.Customer.TaxID = nil
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())
	})
}
