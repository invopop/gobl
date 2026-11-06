package cz_test

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
		Regime:    tax.WithRegime("CZ"),
		Series:    "TEST",
		Code:      "0001",
		IssueDate: cal.MakeDate(2026, 3, 1),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "CZ",
				Code:    "12345679",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Development services",
					Price: num.NewAmount(100000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Accommodation",
					Price: num.NewAmount(100000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "reduced",
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

	t.Run("defaults to CZK", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "CZK", inv.Currency.String())
	})

	t.Run("invalid supplier tax ID code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = "12345678"
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-CZ-TAX-IDENTITY-02]")
	})
}

func TestInvoiceRatesByDate(t *testing.T) {
	tests := []struct {
		date     cal.Date
		standard string
		reduced  string
	}{
		{cal.MakeDate(2026, 3, 1), "21.0%", "12.0%"},
		{cal.MakeDate(2024, 1, 1), "21.0%", "12.0%"},
		{cal.MakeDate(2023, 12, 31), "21.0%", "15.0%"},
		{cal.MakeDate(2013, 1, 1), "21.0%", "15.0%"},
		{cal.MakeDate(2012, 6, 1), "20.0%", "14.0%"},
		{cal.MakeDate(2010, 1, 1), "20.0%", "10.0%"},
	}

	for _, tt := range tests {
		t.Run(tt.date.String(), func(t *testing.T) {
			inv := validInvoice()
			inv.IssueDate = tt.date
			require.NoError(t, inv.Calculate())
			assert.Equal(t, tt.standard, inv.Lines[0].Taxes[0].Percent.String())
			assert.Equal(t, tt.reduced, inv.Lines[1].Taxes[0].Percent.String())
		})
	}
}
