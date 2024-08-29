package hu_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
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
				Country: l10n.HU.Tax(),
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

	// Test 4: Supplier Validation (Group VAT ID without Member ID)
	t.Run("Supplier Group VAT ID without Member ID", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier.TaxID.Code = "88212131503" // Group VAT ID
		inv.Supplier.Identities = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "supplier: (identities: cannot be blank.).")
	})

	// Test 5: Supplier Validation (Group VAT ID with Invalid Member ID)
	t.Run("Supplier Group VAT ID with Invalid Member ID", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier.TaxID.Code = "88212131503" // Group VAT ID
		inv.Supplier.Identities = []*org.Identity{
			{Code: "12345678302"}, // Invalid member ID (9th digit is not 4)
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "supplier: (identities: (code: must be a group member ID.).).")
	})

	// Test 6: Supplier Validation (Group VAT ID with Valid Member ID)
	t.Run("Supplier Group VAT ID with Valid Member ID", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier.TaxID.Code = "88212131503" // Group VAT ID
		inv.Supplier.Identities = []*org.Identity{
			{Code: "12345678402"}, // Valid member ID (9th digit is 4)
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	// Test 7: Invoice Date Validation
	t.Run("Invoice Date Before 2010", func(t *testing.T) {
		inv := baseInvoice()
		inv.IssueDate = cal.MakeDate(2009, 12, 31)
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "issue_date: too early")
	})
}
