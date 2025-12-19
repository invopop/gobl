package ro_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ro"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Valid Real Identities for testing
const (
	// DEDEMAN S.R.L.
	validSupplierCUI = "18547290"
	// EMZ SAH CONF S.R.L.
	validCustomerCUI = "36026362"
)

func TestCorrectionDefinitions(t *testing.T) {
	rd := ro.New()

	// Verify corrections are defined
	require.NotNil(t, rd.Corrections)
	require.Len(t, rd.Corrections, 1)

	// Verify correction types
	corrections := rd.Corrections[0]
	assert.Equal(t, bill.ShortSchemaInvoice, corrections.Schema)
	assert.Contains(t, corrections.Types, bill.InvoiceTypeCreditNote)
	assert.Contains(t, corrections.Types, bill.InvoiceTypeDebitNote)
	assert.Contains(t, corrections.Types, bill.InvoiceTypeCorrective)
}

func TestCreditNoteValidation(t *testing.T) {
	t.Run("credit note without preceding document should fail", func(t *testing.T) {
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeCreditNote,
			Code:      "TEST-CN-001",
			IssueDate: cal.MakeDate(2025, 1, 15),
			Supplier: &org.Party{
				Name: "DEDEMAN S.R.L.",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    validSupplierCUI,
				},
			},
			Customer: &org.Party{
				Name: "EMZ SAH CONF S.R.L.",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    validCustomerCUI,
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(10000, 2),
					},
					Taxes: tax.Set{
						{
							Category: "VAT",
							Percent:  num.NewPercentage(21, 2),
						},
					},
				},
			},
		}

		err := ro.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preceding")
	})

	t.Run("credit note with preceding document should succeed", func(t *testing.T) {
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeCreditNote,
			Code:      "TEST-CN-002",
			IssueDate: cal.MakeDate(2024, 1, 15),
			Supplier: &org.Party{
				Name: "DEDEMAN S.R.L.",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    validSupplierCUI,
				},
			},
			Customer: &org.Party{
				Name: "EMZ SAH CONF S.R.L.",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    validCustomerCUI,
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(10000, 2),
					},
					Taxes: tax.Set{
						{
							Category: "VAT",
							Percent:  num.NewPercentage(19, 2),
						},
					},
				},
			},
			Preceding: []*org.DocumentRef{
				{
					Series:    "INV",
					Code:      "001",
					IssueDate: cal.NewDate(2024, 1, 10),
				},
			},
		}

		err := ro.Validate(inv)
		assert.NoError(t, err)
	})
}

func TestDebitNoteValidation(t *testing.T) {
	t.Run("debit note without preceding document should fail", func(t *testing.T) {
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeDebitNote,
			Code:      "TEST-DN-001",
			IssueDate: cal.MakeDate(2024, 1, 15),
			Supplier: &org.Party{
				Name: "DEDEMAN S.R.L.",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    validSupplierCUI,
				},
			},
			Customer: &org.Party{
				Name: "EMZ SAH CONF S.R.L.",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    validCustomerCUI,
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(10000, 2),
					},
					Taxes: tax.Set{
						{
							Category: "VAT",
							Percent:  num.NewPercentage(19, 2),
						},
					},
				},
			},
		}

		err := ro.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preceding")
	})

	t.Run("debit note with preceding document should succeed", func(t *testing.T) {
		inv := &bill.Invoice{
			Type:      bill.InvoiceTypeDebitNote,
			Code:      "TEST-DN-002",
			IssueDate: cal.MakeDate(2024, 1, 15),
			Supplier: &org.Party{
				Name: "DEDEMAN S.R.L.",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    validSupplierCUI,
				},
			},
			Customer: &org.Party{
				Name: "EMZ SAH CONF S.R.L.",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    validCustomerCUI,
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(10000, 2),
					},
					Taxes: tax.Set{
						{
							Category: "VAT",
							Percent:  num.NewPercentage(19, 2),
						},
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
}
