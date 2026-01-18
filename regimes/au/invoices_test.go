package au_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
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

// TestInvoiceBuyerIdentityThreshold tests the ATO requirement that invoices
// with taxable amount ≥A$1,000 must include buyer name OR ABN.
func TestInvoiceBuyerIdentityThreshold(t *testing.T) {
	tests := []struct {
		name         string
		price        int64 // in cents (e.g., 150000 = A$1,500.00)
		customerName string
		customerABN  cbc.Code // empty = no ABN
		wantErr      bool
	}{
		{
			name:         "under threshold with name",
			price:        50000, // A$500
			customerName: "Test Customer",
			wantErr:      false,
		},
		{
			name:        "under threshold with ABN",
			price:       50000, // A$500
			customerABN: "53004085616",
			wantErr:     false,
		},
		{
			name:         "under threshold without identity",
			price:        90000, // A$900
			customerName: "",
			customerABN:  "",
			wantErr:      false,
		},
		{
			name:         "exactly at threshold without identity",
			price:        100000, // A$1,000
			customerName: "",
			customerABN:  "",
			wantErr:      true,
		},
		{
			name:         "over threshold without identity",
			price:        150000, // A$1,500
			customerName: "",
			customerABN:  "",
			wantErr:      true,
		},
		{
			name:         "over threshold with name only",
			price:        150000, // A$1,500
			customerName: "Test Customer",
			customerABN:  "",
			wantErr:      false,
		},
		{
			name:         "over threshold with ABN only",
			price:        150000, // A$1,500
			customerName: "",
			customerABN:  "53004085616",
			wantErr:      false,
		},
		{
			name:         "over threshold with both name and ABN",
			price:        150000, // A$1,500
			customerName: "Test Customer",
			customerABN:  "53004085616",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := testInvoice()
			inv.Lines[0].Item.Price = num.NewAmount(tt.price, 2)
			inv.Customer.Name = tt.customerName

			if tt.customerABN != "" {
				inv.Customer.TaxID = &tax.Identity{
					Code:    tt.customerABN,
					Country: "AU",
				}
			} else {
				inv.Customer.TaxID = nil
			}

			require.NoError(t, inv.Calculate())
			err := au.Validate(inv)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "buyer name or ABN required")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestInvoiceZeroRatedGST verifies that GST-free (0%) supplies still count
// as taxable supplies for the A$1,000 threshold rule.
func TestInvoiceZeroRatedGST(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500
	inv.Lines[0].Taxes = tax.Set{
		{
			Category: tax.CategoryGST,
			Percent:  num.NewPercentage(0, 2), // 0% GST-free
		},
	}
	inv.Customer.Name = ""
	inv.Customer.TaxID = nil

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.Error(t, err, "Zero-rated GST should still count as taxable supply")
	assert.Contains(t, err.Error(), "buyer name or ABN required")
}

// TestInvoiceNonGSTLines verifies that lines without GST category
// don't count toward the threshold.
func TestInvoiceNonGSTLines(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Taxes = nil                           // No GST category
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500
	inv.Customer.Name = ""
	inv.Customer.TaxID = nil

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.NoError(t, err, "Non-GST lines should not count toward threshold")
}

// TestInvoiceMixedGSTAndNonGST verifies that only GST lines contribute
// to the taxable amount calculation.
func TestInvoiceMixedGSTAndNonGST(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(60000, 2) // A$600 with GST

	// Add a non-GST line
	inv.Lines = append(inv.Lines, &bill.Line{
		Quantity: num.MakeAmount(1, 0),
		Item: &org.Item{
			Name:  "Non-GST Item",
			Price: num.NewAmount(50000, 2), // A$500 without GST
		},
		Taxes: nil, // No GST
	})
	inv.Customer.Name = ""
	inv.Customer.TaxID = nil

	// Only GST lines count: A$600 < A$1,000
	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.NoError(t, err, "Only GST lines should count (A$600 total)")
}

// TestInvoiceMultipleGSTLines verifies that multiple GST lines
// are summed correctly for the threshold check.
func TestInvoiceMultipleGSTLines(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(50000, 2) // A$500

	// Add another GST line: A$500 + A$600 = A$1,100
	inv.Lines = append(inv.Lines, &bill.Line{
		Quantity: num.MakeAmount(1, 0),
		Item: &org.Item{
			Name:  "Another Item",
			Price: num.NewAmount(60000, 2), // A$600
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

// TestInvoiceNilCustomer verifies that nil customer doesn't cause errors
// (simplified invoices may not have customer details).
func TestInvoiceNilCustomer(t *testing.T) {
	inv := testInvoice()
	inv.Lines[0].Item.Price = num.NewAmount(150000, 2) // A$1,500
	inv.Customer = nil

	require.NoError(t, inv.Calculate())
	err := au.Validate(inv)
	require.NoError(t, err, "Nil customer should be allowed")
}
