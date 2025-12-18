package ro_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ro"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateInvoice(t *testing.T) {
	t.Run("valid standard invoice", func(t *testing.T) {
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeStandard,
			Code:      "INV-001",
			IssueDate: cal.MakeDate(2024, 12, 15),
			Supplier: &org.Party{
				Name: "Test Supplier SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "18547290",
				},
			},
			Customer: &org.Party{
				Name: "Test Customer SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "27",
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(10, 0),
					Item: &org.Item{
						Name:  "Test Product",
						Price: num.NewAmount(10000, 2),
					},
					Taxes: []*tax.Combo{
						{
							Category: "VAT",
							Rate:     "standard",
						},
					},
				},
			},
		}

		err := ro.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("valid credit note", func(t *testing.T) {
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeCreditNote,
			Code:      "CN-001",
			IssueDate: cal.MakeDate(2024, 12, 15),
			Supplier: &org.Party{
				Name: "Test Supplier SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "18547290",
				},
			},
			Customer: &org.Party{
				Name: "Test Customer SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "27",
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Refund",
						Price: num.NewAmount(10000, 2),
					},
				},
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "INV",
					Code:   "001",
				},
			},
		}

		err := ro.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeStandard,
			Code:      "INV-002",
			IssueDate: cal.MakeDate(2024, 12, 15),
			Supplier: &org.Party{
				Name: "Test Supplier SRL",
				// Missing TaxID
			},
			Customer: &org.Party{
				Name: "Test Customer",
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Product",
						Price: num.NewAmount(10000, 2),
					},
				},
			},
		}

		err := ro.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tax ID")
	})

	t.Run("missing line item", func(t *testing.T) {
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeStandard,
			Code:      "INV-003",
			IssueDate: cal.MakeDate(2024, 12, 15),
			Supplier: &org.Party{
				Name: "Test Supplier SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "18547290",
				},
			},
			Customer: &org.Party{
				Name: "Test Customer",
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					// Missing Item
				},
			},
		}

		err := ro.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "item")
	})

	t.Run("B2C invoice without customer tax ID", func(t *testing.T) {
		// B2C invoices are mandatory in e-Factura since Jan 1, 2025 (OUG 69/2024).
		// CNP is optional; if missing, e-Factura uses a standard placeholder.
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeStandard,
			Code:      "INV-005",
			IssueDate: cal.MakeDate(2024, 12, 15),
			Supplier: &org.Party{
				Name: "Test Supplier SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "18547290",
				},
			},
			Customer: &org.Party{
				Name: "Individual Customer",
				// No TaxID - B2C scenario
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Product",
						Price: num.NewAmount(10000, 2),
					},
				},
			},
		}

		err := ro.Validate(inv)
		// B2C invoices are valid without customer tax ID
		assert.NoError(t, err)
	})

	t.Run("valid 2025 invoice with new VAT rates", func(t *testing.T) {
		// Test invoice after Aug 1, 2025 when VAT increased to 21% (Law 141/2025)
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeStandard,
			Code:      "INV-2025-001",
			IssueDate: cal.MakeDate(2025, 12, 15),
			Supplier: &org.Party{
				Name: "Test Supplier SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "18547290",
				},
			},
			Customer: &org.Party{
				Name: "Test Customer SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "27",
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(10, 0),
					Item: &org.Item{
						Name:  "Test Product",
						Price: num.NewAmount(10000, 2),
					},
					Taxes: []*tax.Combo{
						{
							Category: "VAT",
							Rate:     "standard", // Should apply 21% for 2025 dates
						},
					},
				},
			},
		}

		err := ro.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("valid 2025 B2C invoice with new mandate", func(t *testing.T) {
		// B2C mandatory reporting since Jan 1, 2025 (OUG 69/2024)
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeStandard,
			Code:      "INV-B2C-2025",
			IssueDate: cal.MakeDate(2025, 12, 15),
			Supplier: &org.Party{
				Name: "Test Supplier SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "18547290",
				},
			},
			Customer: &org.Party{
				Name: "Individual Customer",
				// No TaxID - B2C scenario, valid under OUG 69/2024
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Consumer Product",
						Price: num.NewAmount(5000, 2),
					},
					Taxes: []*tax.Combo{
						{
							Category: "VAT",
							Rate:     "reduced", // Should apply 11% for 2025 dates
						},
					},
				},
			},
		}

		err := ro.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("valid simplified invoice without customer", func(t *testing.T) {
		// Simplified invoices (<100 EUR) are mandatory in e-Factura since Jan 1, 2025 (OUG 138/2024).
		// Customer information is optional for simplified invoices.
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeStandard,
			Code:      "SIMP-001",
			IssueDate: cal.MakeDate(2025, 12, 15),
			Tags: tax.Tags{
				List: []cbc.Key{tax.TagSimplified},
			},
			Supplier: &org.Party{
				Name: "Test Supplier SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "18547290",
				},
			},
			// No customer - valid for simplified invoices
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Coffee",
						Price: num.NewAmount(1500, 2), // 15 EUR
					},
					Taxes: []*tax.Combo{
						{
							Category: "VAT",
							Rate:     "standard",
						},
					},
				},
			},
		}

		err := ro.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("valid simplified invoice with customer", func(t *testing.T) {
		// Simplified invoices can include customer information
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeStandard,
			Code:      "SIMP-002",
			IssueDate: cal.MakeDate(2025, 12, 15),
			Tags: tax.Tags{
				List: []cbc.Key{tax.TagSimplified},
			},
			Supplier: &org.Party{
				Name: "Test Supplier SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "18547290",
				},
			},
			Customer: &org.Party{
				Name: "Walk-in Customer",
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(2, 0),
					Item: &org.Item{
						Name:  "Pastry",
						Price: num.NewAmount(800, 2), // 8 EUR
					},
					Taxes: []*tax.Combo{
						{
							Category: "VAT",
							Rate:     "reduced",
						},
					},
				},
			},
		}

		err := ro.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("non-simplified invoice requires customer", func(t *testing.T) {
		// Regular invoices must have a customer
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeStandard,
			Code:      "INV-006",
			IssueDate: cal.MakeDate(2025, 12, 15),
			// No TagSimplified
			Supplier: &org.Party{
				Name: "Test Supplier SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "18547290",
				},
			},
			// No customer - should fail for non-simplified
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Product",
						Price: num.NewAmount(10000, 2),
					},
				},
			},
		}

		err := ro.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "customer")
	})
}

func TestInvoiceTypes(t *testing.T) {
	validTypes := []cbc.Key{
		bill.InvoiceTypeStandard,
		bill.InvoiceTypeCreditNote,
		bill.InvoiceTypeDebitNote,
		bill.InvoiceTypeCorrective, // Added to match implementation coverage
	}

	for _, invType := range validTypes {
		t.Run(string(invType), func(t *testing.T) {
			inv := &bill.Invoice{
				Type:      invType,
				Code:      "TEST-001",
				IssueDate: cal.MakeDate(2024, 12, 15),
				Supplier: &org.Party{
					Name: "Test Supplier SRL",
					TaxID: &tax.Identity{
						Country: "RO",
						Code:    "18547290",
					},
				},
				Customer: &org.Party{
					Name: "Test Customer",
				},
				Lines: []*bill.Line{
					{
						Quantity: num.MakeAmount(1, 0),
						Item: &org.Item{
							Name:  "Item",
							Price: num.NewAmount(10000, 2),
						},
					},
				},
			}

			// Add preceding document for all correction types
			if invType.In(bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote, bill.InvoiceTypeCorrective) {
				inv.Preceding = []*org.DocumentRef{
					{
						Series: "INV",
						Code:   "001",
					},
				}
			}

			err := ro.Validate(inv)
			assert.NoError(t, err, "Invoice type %s should be valid", invType)
		})
	}
}
