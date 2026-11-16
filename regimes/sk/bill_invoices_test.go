package sk_test

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

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime("SK"),
		Series:    "FV",
		Code:      "2026-001",
		IssueDate: cal.MakeDate(2026, 11, 16),
		Supplier: &org.Party{
			Name: "Vzorová s.r.o.",
			TaxID: &tax.Identity{
				Country: "SK",
				Code:    "2020273893",
			},
		},
		Customer: &org.Party{
			Name: "Zákazník s.r.o.",
			TaxID: &tax.Identity{
				Country: "SK",
				Code:    "2099999902",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Softvérové služby",
					Price: num.NewAmount(100000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("IČ DPH with prefix", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = "SK2020273893"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
		assert.Equal(t, "2020273893", inv.Supplier.TaxID.Code.String())
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-01]")
	})

	t.Run("supplier tax ID without code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{Country: "SK"}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-01]")
	})

	t.Run("supplier tax ID from another country", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{Country: "CZ", Code: "12345679"}
		inv.SetRegime("SK")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-SK-BILL-INVOICE-02]")
		assert.NotContains(t, err.Error(), "[GOBL-SK-BILL-INVOICE-01]")
	})

	t.Run("supplier tax ID with invalid DIČ", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = "2020273894"
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-TAX-IDENTITY-02]")
	})

	// The rule covers the supplier only; § 74(1)(a) does not require the
	// buyer's tax ID on every invoice.
	t.Run("customer without tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("customer with IČO", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.Identities = []*org.Identity{{Type: "ICO", Code: "36 396 567"}}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
		assert.Equal(t, "36396567", inv.Customer.Identities[0].Code.String())
	})
}

func TestInvoiceKeepsSharedScenarios(t *testing.T) {
	inv := validInvoice()
	inv.Tags = tax.WithTags(tax.TagReverseCharge)
	require.NoError(t, inv.Calculate())
	require.NotNil(t, inv.Tax, "the shared reverse charge note must be added")
	require.Len(t, inv.Tax.Notes, 1)
	assert.Equal(t, tax.KeyReverseCharge, inv.Tax.Notes[0].Key)
}
