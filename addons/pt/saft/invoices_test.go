package saft_test

import (
	"fmt"
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("PT"),
		Addons: tax.WithAddons(saft.V1),
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				saft.ExtKeyInvoiceType: saft.InvoiceTypeStandard,
			},
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Code:    "123456789",
				Country: "PT",
			},
			Name: "Test Supplier",
		},
		Customer: &org.Party{
			Name: "Test Customer",
		},
		Series:    "FT SERIES-A",
		Code:      "123",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2023, 1, 1),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("valid invoice", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, addon.Validator(inv))
	})

	t.Run("missing invoice type", func(t *testing.T) {
		inv := validInvoice()

		inv.Tax = nil
		assert.ErrorContains(t, addon.Validator(inv), "tax: (ext: (pt-saft-invoice-type: required")

		inv.Tax = new(bill.Tax)
		assert.ErrorContains(t, addon.Validator(inv), "tax: (ext: (pt-saft-invoice-type: required")
	})

	t.Run("missing VAT category in lines", func(t *testing.T) {
		inv := validInvoice()

		inv.Lines[0].Taxes = nil
		assert.ErrorContains(t, addon.Validator(inv), "lines: (0: (taxes: missing category VAT")
	})
}

func TestInvoiceSeriesValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	tests := []struct {
		series cbc.Code
		code   cbc.Code
		err    string
	}{
		// Nil case
		{"", "", ""},

		// Valid code and series
		{"FT SERIES-A", "123", ""},
		{"", "FT SERIES-A/123", ""},

		// Invalid series
		{"SERIES-A", "123", "series: must start with 'FT '"},
		{"FT SERIES A", "123", "series: must be in a valid format"},
		{"FT SERIES-A/", "123", "series: must be in a valid format"},
		{"XX SERIES-A", "123", "series: must start with 'FT '"},

		// Invalid code (with series present)
		{"FT SERIES-A", "ABCD", "code: must be in a valid format"},
		{"FT SERIES-A", "FT SERIES-A/1234", "code: must be in a valid format"},

		// Invalid code (with series missing)
		{"", "ABCDEF", "code: must start with 'FT '"},
		{"", "SERIES-A/123", "code: must start with 'FT '"},
		{"", "FT SERIES-A 123", "code: must be in a valid format"},
		{"", "XX SERIES-A/1234", "code: must start with 'FT '"},
	}

	for _, test := range tests {
		name := fmt.Sprintf("Series %s Code %s", test.series, test.code)
		t.Run(name, func(t *testing.T) {
			inv := validInvoice()

			inv.Series = test.series
			inv.Code = test.code
			err := addon.Validator(inv)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.err)
			}
		})
	}
}
