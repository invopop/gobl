package au_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoice() *bill.Invoice {
	return &bill.Invoice{
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Code:    "51824753556",
				Country: "AU",
			},
			Name: "Test Supplier Pty Ltd",
			Addresses: []*org.Address{
				{
					Street:  "123 Test Street",
					Code:    "2000",
					Country: l10n.AU.ISO(),
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
		},
		Code:     "INV-001",
		Currency: "AUD",
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(50000, 2), // A$500.00
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryGST,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}

func TestInvoiceUnderThreshold(t *testing.T) {
	inv := testInvoice()
	// Total is A$500, under A$1,000 threshold
	// Customer has name but no ABN - should pass
	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.NoError(t, err)
}

func TestInvoiceOverThresholdWithName(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500.00
	inv.Customer.Name = "Test Customer"
	inv.Customer.TaxID = nil

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.NoError(t, err, "Invoice over threshold should pass with buyer name")
}

func TestInvoiceOverThresholdWithABN(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500.00
	inv.Customer.Name = ""
	inv.Customer.TaxID = &tax.Identity{
		Code:    "53004085616",
		Country: "AU",
	}

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.NoError(t, err, "Invoice over threshold should pass with buyer ABN")
}

func TestInvoiceOverThresholdWithBothNameAndABN(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500.00
	inv.Customer.Name = "Test Customer"
	inv.Customer.TaxID = &tax.Identity{
		Code:    "53004085616",
		Country: "AU",
	}

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.NoError(t, err, "Invoice over threshold should pass with both name and ABN")
}

func TestInvoiceOverThresholdNoIdentity(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500.00
	inv.Customer.Name = ""
	inv.Customer.TaxID = nil

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.Error(t, err, "Invoice over threshold should fail without buyer identity")
	assert.Contains(t, err.Error(), "buyer name or ABN required")
	assert.Contains(t, err.Error(), "A$1,000")
}

func TestInvoiceExactlyAtThreshold(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(100000, 2) // Exactly A$1,000.00
	inv.Customer.Name = ""
	inv.Customer.TaxID = nil

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.Error(t, err, "Invoice at exactly A$1,000 should trigger the rule")
	assert.Contains(t, err.Error(), "buyer name or ABN required")
}

func TestInvoiceJustUnderThreshold(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(90000, 2) // A$900.00
	inv.Customer.Name = ""
	inv.Customer.TaxID = nil

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.NoError(t, err, "Invoice under A$1,000 should not require buyer identity")
}

func TestInvoiceMultipleLinesOverThreshold(t *testing.T) {
	inv := testInvoice()
	// Add another line, total = A$500 + A$600 = A$1,100
	inv.Lines = append(inv.Lines, &bill.Line{
		Quantity: num.MakeAmount(1, 0),
		Item: &org.Item{
			Name:  "Another Item",
			Price: num.NewAmount(60000, 2), // A$600.00
		},
		Taxes: tax.Set{
			{
				Category: tax.CategoryGST,
				Rate:     tax.RateGeneral,
			},
		},
	})
	inv.Customer.Name = ""
	inv.Customer.TaxID = nil

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.Error(t, err, "Multiple lines totaling over A$1,000 should trigger rule")
	assert.Contains(t, err.Error(), "buyer name or ABN required")
}

func TestInvoiceWithZeroRatedGST(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500.00
	inv.Lines[0].Taxes = tax.Set{
		{
			Category: tax.CategoryGST,
			Percent:  num.NewPercentage(0, 2), // 0% GST-free
		},
	}
	inv.Customer.Name = ""
	inv.Customer.TaxID = nil

	// Even with 0% GST, it's still a taxable supply and should trigger the rule
	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.Error(t, err, "Zero-rated GST should still count as taxable supply")
	assert.Contains(t, err.Error(), "buyer name or ABN required")
}

func TestInvoiceNonGSTLinesOnly(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Taxes = nil // No GST category
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500.00
	inv.Customer.Name = ""
	inv.Customer.TaxID = nil

	// No GST lines, so threshold not reached
	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.NoError(t, err, "Non-GST lines should not count toward threshold")
}

func TestInvoiceMixedGSTAndNonGST(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(60000, 2) // A$600.00 with GST

	// Add a non-GST line
	inv.Lines = append(inv.Lines, &bill.Line{
		Quantity: num.MakeAmount(1, 0),
		Item: &org.Item{
			Name:  "Non-GST Item",
			Price: num.NewAmount(50000, 2), // A$500.00 without GST
		},
		Taxes: nil, // No GST
	})
	inv.Customer.Name = ""
	inv.Customer.TaxID = nil

	// Only GST lines count: A$600 < A$1,000
	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.NoError(t, err, "Only GST lines should count, total A$600 under threshold")
}

func TestInvoiceNilCustomer(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500.00
	inv.Customer = nil

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	// Nil customer is allowed (simplified invoices)
	require.NoError(t, err, "Nil customer should not cause validation error")
}

func TestInvoiceWithEmptyABN(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500.00
	inv.Customer.Name = ""
	inv.Customer.TaxID = &tax.Identity{
		Code:    "", // Empty ABN
		Country: "AU",
	}

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.Error(t, err, "Empty ABN should be treated as no ABN")
	assert.Contains(t, err.Error(), "buyer name or ABN required")
}
