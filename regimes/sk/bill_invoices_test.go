package sk_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	_ "github.com/invopop/gobl/regimes"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func skAddress() []*org.Address {
	return []*org.Address{
		{
			Number:   "12",
			Street:   "Hlavná",
			Locality: "Bratislava",
			Code:     "81101",
			Country:  "SK",
		},
	}
}

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime("SK"),
		Series:    "TEST",
		Code:      "0001",
		IssueDate: cal.MakeDate(2027, 1, 15),
		Supplier: &org.Party{
			Name:      "Vzorová s.r.o.",
			TaxID:     &tax.Identity{Country: "SK", Code: "2020273893"},
			Addresses: skAddress(),
		},
		Customer: &org.Party{
			Name:      "Zákazník s.r.o.",
			TaxID:     &tax.Identity{Country: "SK", Code: "7120001713"},
			Addresses: skAddress(),
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Služby",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{Category: "VAT", Rate: "general"},
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

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-01]")
	})

	t.Run("supplier tax ID without code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-01]")
	})

	// The address and the currency are not checked: the regime keeps to what
	// other regimes validate, and e-invoicing rules belong to the formats.
	t.Run("no addresses and no rate to EUR", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Addresses = nil
		inv.Customer.Addresses = nil
		inv.Currency = "USD"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}
