package my_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetRegime("MY")
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (tax_id: cannot be blank.)")
	})
}

// testInvoiceStandard provides a valid minimal Malaysian invoice
func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()

	return &bill.Invoice{
		Regime:    tax.WithRegime("MY"),
		Currency:  "MYR",
		Code:      "INV-0001",
		IssueDate: cal.MakeDate(2025, 4, 24),
		Supplier: &org.Party{
			Name: "Tech Solutions Sdn Bhd",
			TaxID: &tax.Identity{
				Country: "MY",
				Code:    "201901234567",
			},
		},
		Customer: &org.Party{
			Name: "Innovatech Berhad",
			TaxID: &tax.Identity{
				Country: "MY",
				Code:    "202001234567",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Cloud Hosting Services",
					Price: num.NewAmount(10000, 2), // 100.00 MYR
					Unit:  org.UnitHour,
				},
				Taxes: tax.Set{
					{
						Category: "service-tax",
						Percent:  num.NewPercentage(6, 1),
					},
				},
			},
		},
	}
}
