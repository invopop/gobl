package verifactu

import (
	"fmt"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/require"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
	t.Run("missing customer tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		assertValidationError(t, inv, "customer: (tax_id: cannot be blank.)")
	})

	t.Run("with exemption reason", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0].Ext = nil
		assertValidationError(t, inv, "es-verifactu-tax-classification: required")
	})

	t.Run("without notes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Notes = nil
		assertValidationError(t, inv, "notes: with key 'general' missing")
	})

	t.Run("missing doc type", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = nil
		err := inv.Validate()
		require.ErrorContains(t, err, "es-verifactu-doc-type: required")
	})

	t.Run("no customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("correction invoice requires preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		assertValidationError(t, inv, "preceding: cannot be blank")
	})

	t.Run("correction invoice preceding requires issue date", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = []*org.DocumentRef{
			{
				Code: "123",
			},
		}
		assertValidationError(t, inv, "preceding: (0: (issue_date: cannot be blank.).)")
	})

	t.Run("correction invoice with preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		d := cal.MakeDate(2024, 1, 1)
		inv.Preceding = []*org.DocumentRef{
			{
				Series:    "ABC",
				Code:      "122",
				IssueDate: &d,
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	if inv.Preceding != nil {
		fmt.Println(inv.Preceding[0].IssueDate)
	}
	require.ErrorContains(t, err, expected)
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Addons: tax.WithAddons(V1),
		Code:   "123",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "NL",
				Code:    "000099995B57",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.MakeAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "exempt",
						Ext: tax.Extensions{
							ExtKeyTaxClassification: "E1",
						},
					},
				},
			},
		},
		Notes: []*cbc.Note{
			{
				Key:  cbc.NoteKeyGeneral,
				Text: "This is a test invoice",
			},
		},
	}
}
