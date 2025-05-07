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
	"github.com/invopop/gobl/pay"
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

	t.Run("missing doc type", func(t *testing.T) {
		inv := validInvoice()

		inv.Tax = nil
		assert.ErrorContains(t, addon.Validator(inv), "tax: (ext: either `pt-saft-work-type` or `pt-saft-invoice-type` must be set")

		inv.Tax = new(bill.Tax)
		assert.ErrorContains(t, addon.Validator(inv), "tax: (ext: either `pt-saft-work-type` or `pt-saft-invoice-type` must be set")
	})

	t.Run("both doc types set", func(t *testing.T) {
		inv := validInvoice()

		inv.Tax.Ext = tax.Extensions{
			saft.ExtKeyInvoiceType: saft.InvoiceTypeStandard,
			saft.ExtKeyWorkType:    saft.WorkTypeProforma,
		}
		assert.ErrorContains(t, addon.Validator(inv), "tax: (ext: either `pt-saft-work-type` or `pt-saft-invoice-type` must be set, but not both")
	})

	t.Run("work doc type only", func(t *testing.T) {
		inv := validInvoice()

		inv.Series = "PF SERIES-A"
		inv.Tax.Ext = tax.Extensions{
			saft.ExtKeyWorkType: saft.WorkTypeProforma,
		}
		require.NoError(t, addon.Validator(inv))
	})

	t.Run("invalid work type", func(t *testing.T) {
		inv := validInvoice()

		inv.Tax.Ext = tax.Extensions{
			saft.ExtKeyWorkType: saft.WorkTypeBudgets, // Budgets is not valid in invoices, only in orders
		}
		assert.ErrorContains(t, addon.Validator(inv), "value 'OR' invalid")
	})

	t.Run("missing VAT category in lines", func(t *testing.T) {
		inv := validInvoice()

		inv.Lines[0].Taxes = nil
		assert.ErrorContains(t, addon.Validator(inv), "lines: (0: (taxes: missing category VAT")
	})

	t.Run("unpaid invoice-receipt", func(t *testing.T) {
		inv := validInvoice()
		inv.Series = "FR SERIES-A"
		inv.Tax.Ext = tax.Extensions{
			saft.ExtKeyInvoiceType: saft.InvoiceTypeInvoiceReceipt,
		}
		inv.Totals = &bill.Totals{
			Due: num.NewAmount(10, 2), // Some payment due
		}

		assert.ErrorContains(t, addon.Validator(inv), "totals: (due: must be no greater than 0")

		inv.Totals.Due = num.NewAmount(0, 2) // Zero payment due
		require.NoError(t, addon.Validator(inv))

		inv.Totals.Due = nil // No payment due
		require.NoError(t, addon.Validator(inv))
	})
}

func TestSimplifiedInvoiceValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	tests := []struct {
		name  string
		typ   cbc.Code
		total num.Amount
		err   string
	}{
		{
			name:  "simplified invoice just below limit",
			typ:   saft.InvoiceTypeSimplified,
			total: num.MakeAmount(99999, 2), // 999.99
		},
		{
			name:  "simplified invoice exactly at limit",
			typ:   saft.InvoiceTypeSimplified,
			total: num.MakeAmount(100000, 2), // 1000.00
		},
		{
			name:  "simplified invoice just over limit",
			typ:   saft.InvoiceTypeSimplified,
			total: num.MakeAmount(100001, 2), // 1000.01
			err:   "must be no greater than 1000",
		},
		{
			name:  "standard invoice over simplified limit",
			typ:   saft.InvoiceTypeStandard,
			total: num.MakeAmount(100001, 2), // 1000.01
		},
		{
			name:  "invoice-receipt over simplified limit",
			typ:   saft.InvoiceTypeInvoiceReceipt,
			total: num.MakeAmount(100001, 2), // 1000.01
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inv := validInvoice()
			series := fmt.Sprintf("%s SERIES-A", test.typ)
			inv.Series = cbc.Code(series)
			inv.Tax.Ext = tax.Extensions{
				saft.ExtKeyInvoiceType: test.typ,
			}
			inv.Totals = &bill.Totals{
				TotalWithTax: test.total,
			}

			err := addon.Validator(inv)

			if test.err != "" {
				assert.ErrorContains(t, err, test.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
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

func TestInvoicePaymentValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("invoice with valid advance", func(t *testing.T) {
		inv := validInvoice()
		date := cal.MakeDate(2023, 1, 1) // Same as invoice date
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Date:   &date,
					Amount: num.MakeAmount(50, 0),
				},
			},
		}
		require.NoError(t, addon.Validator(inv))
	})

	t.Run("advance with different date than invoice", func(t *testing.T) {
		inv := validInvoice()
		date := cal.MakeDate(2023, 1, 2) // Different than invoice date
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Date:   &date,
					Amount: num.MakeAmount(50, 0),
				},
			},
		}
		assert.ErrorContains(t, addon.Validator(inv), "advances: (0: (date: must be the same as the invoice issue date")
	})

	t.Run("advance with nil date", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Date:   nil,
					Amount: num.MakeAmount(50, 0),
				},
			},
		}
		assert.ErrorContains(t, addon.Validator(inv), "advances: (0: (date: cannot be blank")
	})
}

func TestInvoicePrecedingValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("preceding document with no date", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				IssueDate: nil,
			},
		}
		require.NoError(t, addon.Validator(inv))
	})

	t.Run("valid preceding document date", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				IssueDate: cal.NewDate(2022, 12, 31), // Before invoice date
			},
		}
		require.NoError(t, addon.Validator(inv))
	})

	t.Run("preceding document with same date", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				IssueDate: cal.NewDate(2023, 1, 1), // Same as invoice date
			},
		}
		require.NoError(t, addon.Validator(inv))
	})

	t.Run("preceding document with future date", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				IssueDate: cal.NewDate(2023, 1, 2), // After invoice date
			},
		}
		assert.ErrorContains(t, addon.Validator(inv), "preceding: (0: (issue_date: too late")
	})
}

func TestInvoiceNormalization(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("set default advance date", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Date:   nil,
					Amount: num.MakeAmount(50, 0),
				},
			},
		}

		addon.Normalizer(inv)

		assert.Equal(t, &inv.IssueDate, inv.Payment.Advances[0].Date)
	})

	t.Run("nil payment details", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = nil

		addon.Normalizer(inv)

		assert.Nil(t, inv.Payment)
	})
}
