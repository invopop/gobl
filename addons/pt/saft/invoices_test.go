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
				saft.ExtKeyInvoiceType:   saft.InvoiceTypeStandard,
				saft.ExtKeySourceBilling: saft.SourceBillingProduced,
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
			saft.ExtKeyWorkType:      saft.WorkTypeProforma,
			saft.ExtKeySourceBilling: saft.SourceBillingProduced,
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

	t.Run("missing source billing", func(t *testing.T) {
		inv := validInvoice()
		delete(inv.Tax.Ext, saft.ExtKeySourceBilling)
		assert.ErrorContains(t, addon.Validator(inv), "tax: (ext: (pt-saft-source-billing: required.).).")
	})

	t.Run("source billing produced - no source doc ref required", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax.Ext = tax.Extensions{
			saft.ExtKeyInvoiceType:   saft.InvoiceTypeStandard,
			saft.ExtKeySourceBilling: saft.SourceBillingProduced,
		}
		require.NoError(t, addon.Validator(inv))
	})

	t.Run("source billing integrated - source doc ref required", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax.Ext = tax.Extensions{
			saft.ExtKeyInvoiceType:   saft.InvoiceTypeStandard,
			saft.ExtKeySourceBilling: saft.SourceBillingIntegrated,
		}
		assert.ErrorContains(t, addon.Validator(inv), "tax: (ext: (pt-saft-source-ref: required.).).")

		// Add source doc ref - should pass
		inv.Tax.Ext[saft.ExtKeySourceRef] = "FTM abc/00001"
		require.NoError(t, addon.Validator(inv))
	})

	t.Run("source billing manual - source doc ref required", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax.Ext = tax.Extensions{
			saft.ExtKeyInvoiceType:   saft.InvoiceTypeStandard,
			saft.ExtKeySourceBilling: saft.SourceBillingManual,
		}
		assert.ErrorContains(t, addon.Validator(inv), "tax: (ext: (pt-saft-source-ref: required.).).")

		// Add source doc ref - should pass
		inv.Tax.Ext[saft.ExtKeySourceRef] = "FTD FT SERIESA/123"
		require.NoError(t, addon.Validator(inv))
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

func TestSourceRefFormatValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	tests := []struct {
		ref string
		err string
	}{
		{"FTM abc/00001", ""},
		{"FTD FT SERIESA/123", ""},
		{"FTR abc/00001", "must be in valid format"},
		{"FTM a/bc/00001", "must be in valid format"},
		{"FTDA FT abc/00001", "must be in valid format"},
		{"ABC abc/00001", "must be in valid format"},
		{"FTM FT abc/00001", "must be in valid format"},
		{"FRM abc/00001", "must start with the document type 'FT' not 'FR'"},
		{"FRD FT SERIESA/123", "must start with the document type 'FT' not 'FR'"},
		{"FTD FR SERIESA/123", "must refer to an original document 'FT' not 'FR'"},
	}

	for _, test := range tests {
		t.Run(test.ref, func(t *testing.T) {
			inv := validInvoice()
			inv.Tax.Ext[saft.ExtKeySourceBilling] = saft.SourceBillingManual
			inv.Tax.Ext[saft.ExtKeySourceRef] = cbc.Code(test.ref)

			err := addon.Validator(inv)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.err)
			}
		})
	}
}

func TestInvoiceNormalization(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("normalize invoice with nil tax", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = nil

		addon.Normalizer(inv)

		require.NotNil(t, inv.Tax)
		require.NotNil(t, inv.Tax.Ext)
		assert.Equal(t, saft.SourceBillingProduced, inv.Tax.Ext[saft.ExtKeySourceBilling])
	})

	t.Run("normalize invoice with nil tax extensions", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{}

		addon.Normalizer(inv)

		require.NotNil(t, inv.Tax.Ext)
		assert.Equal(t, saft.SourceBillingProduced, inv.Tax.Ext[saft.ExtKeySourceBilling])
	})

	t.Run("normalize invoice with missing source billing", func(t *testing.T) {
		inv := validInvoice()
		delete(inv.Tax.Ext, saft.ExtKeySourceBilling)

		addon.Normalizer(inv)

		assert.Equal(t, saft.SourceBillingProduced, inv.Tax.Ext[saft.ExtKeySourceBilling])
	})

	t.Run("normalize invoice with existing source billing", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax.Ext[saft.ExtKeySourceBilling] = saft.SourceBillingIntegrated

		addon.Normalizer(inv)

		assert.Equal(t, saft.SourceBillingIntegrated, inv.Tax.Ext[saft.ExtKeySourceBilling])
	})
}
