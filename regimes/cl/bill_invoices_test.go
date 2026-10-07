package cl_test

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
		inv.SetRegime("CL")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CL-BILL-INVOICE-01]")
	})

	t.Run("supplier tax ID without code", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = &tax.Identity{Country: "CL"}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CL-BILL-INVOICE-01]")
	})

	t.Run("supplier tax ID from another country", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Country = "US"
		inv.SetRegime("CL")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CL-BILL-INVOICE-02]")
	})

	t.Run("supplier tax ID with invalid RUT", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Code = "123456781"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CL-TAX-IDENTITY-01]")
	})

	t.Run("missing customer tax ID on standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CL-BILL-INVOICE-03]")
	})

	t.Run("customer tax ID without code on standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = &tax.Identity{Country: "CL"}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CL-BILL-INVOICE-03]")
	})

	t.Run("simplified invoice without customer tax ID is valid", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.SetTags(tax.TagSimplified)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified invoice still requires supplier tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = nil
		inv.SetTags(tax.TagSimplified)
		inv.SetRegime("CL")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CL-BILL-INVOICE-01]")
	})
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Currency:  "CLP",
		IssueDate: cal.MakeDate(2026, 6, 15),
		Series:    "F001",
		Code:      "00000123",
		Supplier: &org.Party{
			Name: "Proveedor Ejemplo SpA",
			TaxID: &tax.Identity{
				Country: "CL",
				Code:    "123456785",
			},
		},
		Customer: &org.Party{
			Name: "Cliente Comercial S.A.",
			TaxID: &tax.Identity{
				Country: "CL",
				Code:    "111111111",
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
