package cy_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CY-BILL-INVOICE-01]")
	})

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CY-BILL-INVOICE-02]")
	})

	t.Run("missing supplier address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CY-BILL-INVOICE-03]")
	})
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 8, 19),
		Code:      "INV-001",
		Supplier: &org.Party{
			Name: "Cyprus Supplier Ltd",
			TaxID: &tax.Identity{
				Country: "CY",
				Code:    "60000000A",
			},
			Addresses: []*org.Address{
				{
					Street:   "Example Street",
					Locality: "Nicosia",
					Country:  "CY",
				},
			},
		},
		Customer: &org.Party{
			Name: "Customer Ltd",
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Professional services",
					Price: num.NewAmount(10000, 2),
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
	inv.SetRegime("CY")
	return inv
}
