package saft_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// withAddonContext returns a rules.WithContext that injects the SAF-T V1
// addon into the validation context. Used for standalone objects that don't
// have Addons embedded (like bill.Payment, bill.Delivery, etc.).
func withAddonContext() rules.WithContext {
	return tax.AddonContext(saft.V1)
}

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("PT"),
		Addons: tax.WithAddons(saft.V1),
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				saft.ExtKeyInvoiceType: saft.InvoiceTypeStandard,
				saft.ExtKeySource:      saft.SourceBillingProduced,
			}),
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
		ValueDate: cal.NewDate(2022, 12, 31),
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
						Rate:     "general",
					},
				},
			},
		},
	}
}

func calculatedInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
	return inv
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice", func(t *testing.T) {
		inv := calculatedInvoice(t)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing doc type", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext = tax.ExtensionsOf(tax.ExtMap{})
		assert.ErrorContains(t, rules.Validate(inv), "either 'pt-saft-work-type' or 'pt-saft-invoice-type' must be set")

		inv.Tax = nil
		assert.ErrorContains(t, rules.Validate(inv), "either 'pt-saft-work-type' or 'pt-saft-invoice-type' must be set")

		inv.Tax = new(bill.Tax)
		assert.ErrorContains(t, rules.Validate(inv), "either 'pt-saft-work-type' or 'pt-saft-invoice-type' must be set")
	})

	t.Run("both doc types set", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeyWorkType, saft.WorkTypeProforma)
		assert.ErrorContains(t, rules.Validate(inv), "but not both")
	})

	t.Run("work doc type only", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Series = "PF SERIES-A"
		inv.Tax.Ext = tax.ExtensionsOf(tax.ExtMap{
			saft.ExtKeyWorkType: saft.WorkTypeProforma,
			saft.ExtKeySource:   saft.SourceBillingProduced,
		})
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid work type", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Delete(saft.ExtKeyInvoiceType)
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeyWorkType, saft.WorkTypeBudgets)
		assert.ErrorContains(t, rules.Validate(inv), "invoice work type is not valid")
	})

	t.Run("missing VAT category in lines", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Lines[0].Taxes = nil
		assert.ErrorContains(t, rules.Validate(inv), "line taxes must include VAT category")
	})

	t.Run("non-EUR currency without exchange rates", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Currency = "USD"
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-INVOICE-15] invoice must be in EUR or provide exchange rate for conversion")
	})

	t.Run("non-EUR currency with exchange rates", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Currency = "USD"
		inv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   "USD",
				To:     "EUR",
				Amount: num.MakeAmount(875967, 6),
			},
		}
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("missing source billing", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Delete(saft.ExtKeySource)
		assert.ErrorContains(t, rules.Validate(inv), "tax requires 'pt-saft-source' extension")
	})

	t.Run("source billing produced - no source doc ref required", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeySource, saft.SourceBillingProduced)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("source billing integrated - source doc ref required", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeySource, saft.SourceBillingIntegrated)
		assert.ErrorContains(t, rules.Validate(inv), "tax requires 'pt-saft-source-ref' extension when source is not produced")

		// Add source doc ref - should pass
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeySourceRef, "FTM abc/00001")
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("source billing manual - source doc ref required", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeySource, saft.SourceBillingManual)
		assert.ErrorContains(t, rules.Validate(inv), "tax requires 'pt-saft-source-ref' extension when source is not produced")

		// Add source doc ref
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeySourceRef, "FTD FT SERIESA/123")
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("unpaid invoice-receipt", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Series = "FR SERIES-A"
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeyInvoiceType, saft.InvoiceTypeInvoiceReceipt)
		inv.Totals = &bill.Totals{
			Due: num.NewAmount(10, 2),
		}
		assert.ErrorContains(t, rules.Validate(inv), "must be equal to 0")

		inv.Totals.Due = num.NewAmount(0, 2)
		require.NoError(t, rules.Validate(inv))

		inv.Totals.Due = nil
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing value date", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.ValueDate = nil
		assert.ErrorContains(t, rules.Validate(inv), "cannot be blank")
	})
}

func TestInvoiceSeriesValidation(t *testing.T) {
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
		{"SERIES-A", "123", "series format must be valid"},
		{"FT SERIES A", "123", "series format must be valid"},
		{"FT SERIES-A/", "123", "series format must be valid"},
		{"XX SERIES-A", "123", "series format must be valid"},

		// Invalid code (with series present)
		{"FT SERIES-A", "ABCD", "code format must be valid"},
		{"FT SERIES-A", "FT SERIES-A/1234", "code format must be valid"},

		// Invalid code (with series missing)
		{"", "ABCDEF", "code format must be valid"},
		{"", "SERIES-A/123", "code format must be valid"},
		{"", "FT SERIES-A 123", "code format must be valid"},
		{"", "XX SERIES-A/1234", "code format must be valid"},
	}

	for _, test := range tests {
		name := "Series " + string(test.series) + " Code " + string(test.code)
		t.Run(name, func(t *testing.T) {
			inv := calculatedInvoice(t)
			inv.Series = test.series
			inv.Code = test.code
			err := rules.Validate(inv)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.err)
			}
		})
	}
}

func TestSourceRefFormatValidation(t *testing.T) {
	t.Run("missing source ref", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Delete(saft.ExtKeySourceRef)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing invoice type", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Delete(saft.ExtKeyInvoiceType)
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeyWorkType, saft.WorkTypeProforma)
		inv.Series = "PF SERIES-A"
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("integrated document", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeySource, saft.SourceBillingIntegrated)
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeySourceRef, "FTR abc/00001")
		require.NoError(t, rules.Validate(inv))
	})

	tests := []struct {
		ref string
		err string
	}{
		{"FTM abc/00001", ""},
		{"FTD FT SERIESA/123", ""},
		{"FTR abc/00001", "source ref format is invalid"},
		{"FTM a/bc/00001", "source ref format is invalid"},
		{"FTDA FT abc/00001", "source ref format is invalid"},
		{"ABC abc/00001", "source ref format is invalid"},
		{"FTM FT abc/00001", "source ref format is invalid"},
		{"FRM abc/00001", "source ref format is invalid"},
		{"FRD FT SERIESA/123", "source ref format is invalid"},
		{"FTD FR SERIESA/123", "source ref format is invalid"},
	}

	for _, test := range tests {
		t.Run(test.ref, func(t *testing.T) {
			inv := calculatedInvoice(t)
			inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeySource, saft.SourceBillingManual)
			inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeySourceRef, cbc.Code(test.ref))

			err := rules.Validate(inv)
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
		assert.Equal(t, saft.SourceBillingProduced, inv.Tax.Ext.Get(saft.ExtKeySource))
	})

	t.Run("normalize invoice with nil tax extensions", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{}

		addon.Normalizer(inv)

		require.NotNil(t, inv.Tax.Ext)
		assert.Equal(t, saft.SourceBillingProduced, inv.Tax.Ext.Get(saft.ExtKeySource))
	})

	t.Run("normalize invoice with missing source billing", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax.Ext = inv.Tax.Ext.Delete(saft.ExtKeySource)

		addon.Normalizer(inv)

		assert.Equal(t, saft.SourceBillingProduced, inv.Tax.Ext.Get(saft.ExtKeySource))
	})

	t.Run("normalize invoice with existing source billing", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax.Ext = inv.Tax.Ext.Set(saft.ExtKeySource, saft.SourceBillingIntegrated)

		addon.Normalizer(inv)

		assert.Equal(t, saft.SourceBillingIntegrated, inv.Tax.Ext.Get(saft.ExtKeySource))
	})

	t.Run("nil invoice", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var inv *bill.Invoice
			addon.Normalizer(inv)
		})
	})

	t.Run("sets default value date from issue date", func(t *testing.T) {
		inv := validInvoice()
		inv.ValueDate = nil
		addon.Normalizer(inv)
		assert.Equal(t, &inv.IssueDate, inv.ValueDate)
	})

	t.Run("sets default value date from operation date", func(t *testing.T) {
		inv := validInvoice()
		inv.OperationDate = cal.NewDate(2022, 12, 30)
		inv.ValueDate = nil
		addon.Normalizer(inv)
		assert.Equal(t, inv.OperationDate, inv.ValueDate)
	})

	t.Run("keeps existing value date", func(t *testing.T) {
		inv := validInvoice()
		inv.ValueDate = cal.NewDate(2022, 12, 30)
		addon.Normalizer(inv)
		assert.Equal(t, cal.NewDate(2022, 12, 30), inv.ValueDate)
	})

	t.Run("sets today as value date when no issue date is set", func(t *testing.T) {
		inv := validInvoice()
		inv.IssueDate = cal.Date{}
		inv.ValueDate = nil

		addon.Normalizer(inv)

		loc, err := time.LoadLocation("Europe/Lisbon")
		require.NoError(t, err)
		today := cal.TodayIn(loc)
		assert.Equal(t, &today, inv.ValueDate)
	})
	t.Run("normalize invoice with reverse charge", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Taxes[0] = &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyReverseCharge,
		}
		inv.Normalize(tax.ExtractNormalizers(inv))
		assert.Equal(t, tax.KeyReverseCharge, inv.Lines[0].Taxes[0].Key)
		assert.Equal(t, "ISE", inv.Lines[0].Taxes[0].Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M40", inv.Lines[0].Taxes[0].Ext.Get(saft.ExtKeyExemption).String())
		require.Len(t, inv.Lines[0].Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Lines[0].Notes[0].Key)
		assert.Equal(t, "M40", inv.Lines[0].Notes[0].Code.String())
		assert.Equal(t, saft.ExtKeyExemption, inv.Lines[0].Notes[0].Src)
		assert.Equal(t, "Artigo 6.º n.º 6 alínea a) do CIVA, a contrário", inv.Lines[0].Notes[0].Text)
	})
}

func TestInvoicePaymentValidation(t *testing.T) {
	t.Run("advance with nil date", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Date:   nil,
					Amount: num.MakeAmount(50, 0),
				},
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "cannot be blank")
	})

	t.Run("nil advance", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{nil},
		}
		require.NoError(t, rules.Validate(inv))
	})
}

func TestInvoicePaymentNormalization(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("set default advance date", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Date: nil,
				},
			},
		}

		addon.Normalizer(inv)

		assert.Equal(t, &inv.IssueDate, inv.Payment.Advances[0].Date)
	})

	t.Run("no issue date", func(t *testing.T) {
		inv := validInvoice()
		inv.IssueDate = cal.Date{}
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Date: nil,
				},
			},
		}

		addon.Normalizer(inv)

		loc, err := time.LoadLocation("Europe/Lisbon")
		require.NoError(t, err)

		today := cal.TodayIn(loc)
		assert.Equal(t, &today, inv.Payment.Advances[0].Date)
	})

	t.Run("nil payment details", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = nil

		addon.Normalizer(inv)

		assert.Nil(t, inv.Payment)
	})
}

func TestInvoicePrecedingValidation(t *testing.T) {
	t.Run("nil preceding", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Preceding = nil
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid preceding", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV/1",
				IssueDate: cal.NewDate(2023, 1, 1),
			},
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("several preceding documents", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV/1",
				IssueDate: cal.NewDate(2023, 1, 1),
			},
			{
				Code:      "INV/2",
				IssueDate: cal.NewDate(2023, 1, 1),
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "the length must be no more than 1")
	})
}

func TestInvoiceTotalsValidation(t *testing.T) {
	t.Run("valid payable amount", func(t *testing.T) {
		inv := calculatedInvoice(t)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("negative payable amount", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Totals = &bill.Totals{
			Payable: num.MakeAmount(-10, 2),
		}
		assert.ErrorContains(t, rules.Validate(inv), "must be no less than 0")
	})

	// nil totals is now caught by core GOBL rules (invoice totals are required),
	// so no addon-specific test is needed.
}

func TestCorrectionDefinitions(t *testing.T) {
	t.Run("correction definitions exist for credit and debit notes", func(t *testing.T) {
		addon := tax.AddonForKey(saft.V1)
		require.NotNil(t, addon.Corrections)
		def := addon.Corrections.Def(bill.ShortSchemaInvoice)
		require.NotNil(t, def)
		assert.True(t, def.HasType(bill.InvoiceTypeCreditNote))
		assert.True(t, def.HasType(bill.InvoiceTypeDebitNote))
		assert.True(t, def.ReasonRequired)
	})
}
