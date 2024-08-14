package hu_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func baseInvoice() *bill.Invoice {
	inv := &bill.Invoice{
		Currency:  currency.HUF,
		Code:      "TEST",
		IssueDate: cal.MakeDate(2023, 8, 14),
		Type:      bill.InvoiceTypeStandard,
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "HU",
				Code:    "88212131503", // Group VAT ID
			},
			Identities: []*org.Identity{
				{Code: "12345678402"}, // Invalid VAT code
			},
			Addresses: []*org.Address{
				{
					Locality: "Budapest",
					Region:   "Pest",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "HU",
				Code:    "98109858",
			},
			Addresses: []*org.Address{
				{
					Locality: "Debrecen",
					Region:   "Hajdú-Bihar",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(1000, 0),
				},
			},
		},
	}
	return inv
}

func TestInvoiceValidation(t *testing.T) {
	// Test 1: Basic Validation
	t.Run("Valid Invoice", func(t *testing.T) {
		inv := baseInvoice()
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	// Test 2: Customer Validation (Missing TaxID)
	t.Run("Customer Missing TaxID", func(t *testing.T) {
		inv := baseInvoice()
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		println(err.Error())
		assert.Contains(t, err.Error(), "customer: (tax_id: cannot be blank.).")
	})

	// Test 3: Customer Validation (Invalid Group VAT Code)
	t.Run("Customer Invalid Group VAT Code", func(t *testing.T) {
		inv := baseInvoice()
		inv.Customer.TaxID.Code = "98109858512" // Group VAT ID with 9th character 5
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "customer: (identities: cannot be blank.).")
	})
}
