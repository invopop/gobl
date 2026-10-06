package ro_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Unit tests to check the correction types the regime supports
func TestCorrections(t *testing.T) {
	tests := []struct {
		name string
		want cbc.Key
		err  string
	}{
		{name: "credit note", want: bill.InvoiceTypeCreditNote},
		{name: "corrective", want: bill.InvoiceTypeCorrective},
		{name: "debit note", want: bill.InvoiceTypeDebitNote, err: "invalid correction type: debit-note"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := testInvoice()
			require.NoError(t, inv.Calculate())

			err := inv.Correct(correctionOption(tt.want), bill.WithReason("Corectare"))
			if tt.err != "" {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, inv.Type)
			if assert.Len(t, inv.Preceding, 1) {
				assert.Equal(t, cbc.Code("0001"), inv.Preceding[0].Code)
				assert.Equal(t, cal.NewDate(2026, 2, 10), inv.Preceding[0].IssueDate)
			}
			assert.NoError(t, rules.Validate(inv))
		})
	}
}

// correctionOption returns the Correct option for each invoice type
func correctionOption(kind cbc.Key) schema.Option {
	switch kind {
	case bill.InvoiceTypeCorrective:
		return bill.Corrective
	case bill.InvoiceTypeDebitNote:
		return bill.Debit
	default:
		return bill.Credit
	}
}

func testInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime("RO"),
		Series:    "FACT",
		Code:      "0001",
		IssueDate: cal.MakeDate(2026, 2, 10),
		Currency:  "RON",
		Supplier: &org.Party{
			Name:  "Furnizor Exemplu S.R.L.",
			TaxID: &tax.Identity{Country: "RO", Code: "16281620"},
		},
		Customer: &org.Party{
			Name:  "Client Exemplu S.R.L.",
			TaxID: &tax.Identity{Country: "RO", Code: "23456783"},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Servicii de consultanță",
					Price: num.NewAmount(25000, 2),
				},
				Taxes: tax.Set{
					{Category: tax.CategoryVAT, Rate: tax.RateGeneral},
				},
			},
		},
	}
}
