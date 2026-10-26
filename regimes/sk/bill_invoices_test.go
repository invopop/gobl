package sk_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
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

	t.Run("supplier without address", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-02]")
	})

	t.Run("supplier address missing the postal code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Addresses[0].Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-02]")
	})

	t.Run("supplier address missing the locality", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Addresses[0].Locality = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-02]")
	})

	t.Run("supplier address missing the street", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Addresses[0].Street = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-02]")
	})

	t.Run("second address complete", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Addresses = append([]*org.Address{{Locality: "Žilina"}}, inv.Supplier.Addresses...)
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("customer without address", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-03]")
	})

	t.Run("customer with a foreign tax ID still needs an address", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.TaxID = &tax.Identity{Country: "CZ", Code: "12345678"}
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-03]")
	})

	t.Run("no customer", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified invoice does not need a customer address", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(tax.TagSimplified)
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified invoice still needs the supplier address", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(tax.TagSimplified)
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-02]")
	})
}

func TestInvoiceCurrency(t *testing.T) {
	usd := func() *bill.Invoice {
		inv := validInvoice()
		inv.Currency = "USD"
		return inv
	}

	t.Run("EUR needs no rate", func(t *testing.T) {
		inv := validInvoice()
		inv.Currency = currency.EUR
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("not EUR without rate", func(t *testing.T) {
		inv := usd()
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-04]")
	})

	t.Run("not EUR with a rate to another currency", func(t *testing.T) {
		inv := usd()
		inv.ExchangeRates = []*currency.ExchangeRate{
			{From: "USD", To: "CZK", Amount: num.MakeAmount(2300, 2)},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-BILL-INVOICE-04]")
	})

	t.Run("not EUR with a rate to EUR", func(t *testing.T) {
		inv := usd()
		inv.ExchangeRates = []*currency.ExchangeRate{
			{From: "USD", To: currency.EUR, Amount: num.MakeAmount(92, 2)},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}
