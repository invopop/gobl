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
				Country: l10n.HU.Tax(),
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
	// Test 1: Valid Invoice
	t.Run("Valid Invoice", func(t *testing.T) {
		inv := baseInvoice()
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	// Test 2: Customer Invalid Group VAT Code
	t.Run("Customer Invalid Group VAT Code", func(t *testing.T) {
		inv := baseInvoice()
		inv.Customer.TaxID.Code = "98109858512" // Group VAT ID with 9th character 5
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "customer: (identities: cannot be blank.)")
	})

	// Test 3: Supplier Group VAT ID without Member ID
	t.Run("Supplier Group VAT ID without Member ID", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier.TaxID.Code = "88212131503" // Group VAT ID
		inv.Supplier.Identities = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "supplier: (identities: cannot be blank.)")
	})

	// Test 4: Supplier Group VAT ID with Invalid Member ID
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

	// Test 5: Supplier Group VAT ID with Valid Member ID
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

	// Test 6: Invoice Date Validation
	t.Run("Invoice Date Before 2010", func(t *testing.T) {
		inv := baseInvoice()
		inv.IssueDate = cal.MakeDate(2009, 12, 31)
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "issue_date: too early")
	})

	// Test 7: Invalid Invoice Type
	t.Run("Invalid Invoice Type", func(t *testing.T) {
		inv := baseInvoice()
		inv.Type = "InvalidType"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type: must be a valid value")
	})

	// Test 8: Credit Note without Preceding Invoice
	t.Run("Credit Note without Preceding Invoice", func(t *testing.T) {
		inv := baseInvoice()
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "preceding: cannot be blank")
	})

	// Test 9: Supplier without Address
	t.Run("Supplier without Address", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "supplier: (addresses: cannot be blank.)")
	})

	// Test 10: Valid Credit Note
	t.Run("Valid Credit Note", func(t *testing.T) {
		inv := baseInvoice()
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = []*org.DocumentRef{
			{
				Code: "TEST-001",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}
