package ar_test

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
		inv.SetRegime("AR")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-AR-BILL-INVOICE-01]")
	})
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Currency:  "ARS",
		IssueDate: cal.MakeDate(2025, 1, 15),
		Series:    "1",
		Code:      "00000123",
		Supplier: &org.Party{
			Name: "Proveedor Ejemplo S.A.",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30712345671",
			},
		},
		Customer: &org.Party{
			Name: "Cliente Comercial S.R.L.",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30987654321",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Notebooks Dell Latitude",
					Price: num.NewAmount(450000, 0),
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
}
