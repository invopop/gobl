package pe_test

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
	t.Parallel()

	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = nil
		inv.SetRegime("PE")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-PE-BILL-INVOICE-01]")
	})

	t.Run("supplier tax ID without code", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = &tax.Identity{Country: "PE"}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-PE-BILL-INVOICE-01]")
	})

	t.Run("supplier tax ID from another country", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Country = "US"
		inv.SetRegime("PE")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-PE-BILL-INVOICE-02]")
	})

	t.Run("supplier tax ID with invalid RUC", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Code = "20547825443"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-PE-TAX-IDENTITY-01]")
	})
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Currency:  "PEN",
		IssueDate: cal.MakeDate(2026, 6, 15),
		Series:    "F001",
		Code:      "00000123",
		Supplier: &org.Party{
			Name: "Proveedor Ejemplo S.A.C.",
			TaxID: &tax.Identity{
				Country: "PE",
				Code:    "20131312955",
			},
		},
		Customer: &org.Party{
			Name: "Cliente Comercial S.A.",
			TaxID: &tax.Identity{
				Country: "PE",
				Code:    "20100047218",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Laptops Dell Latitude",
					Price: num.NewAmount(150000, 2),
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
