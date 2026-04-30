package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validPayment() *bill.Payment {
	return &bill.Payment{
		Type: bill.PaymentTypeReceipt,
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "PT",
				Code:    "545259045",
			},
		},
		Customer: &org.Party{
			Name: "Customer Name",
			TaxID: &tax.Identity{
				Country: "PT",
				Code:    "545259045",
			},
		},
		Currency: "EUR",
		Ext: tax.ExtensionsOf(tax.ExtMap{
			saft.ExtKeyPaymentType: saft.PaymentTypeOther,
			saft.ExtKeySource:      saft.SourceBillingProduced,
		}),
		Series:    "RG SERIES-A",
		Code:      "123",
		IssueDate: cal.MakeDate(2024, 3, 10),
		Lines: []*bill.PaymentLine{
			{
				Document: &org.DocumentRef{
					Code:      "INV/1",
					IssueDate: cal.NewDate(2024, 3, 1),
				},
				Amount: num.MakeAmount(100, 2),
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: tax.CategoryVAT,
							Rates: []*tax.RateTotal{
								{
									Ext: tax.ExtensionsOf(tax.ExtMap{
										pt.ExtKeyRegion:    "PT",
										saft.ExtKeyTaxRate: "NOR",
									}),
								},
							},
						},
					},
				},
			},
		},
		Methods: []*pay.Record{
			{Key: "credit-transfer"},
		},
	}
}

func TestPaymentValidation(t *testing.T) {
	t.Run("valid payment", func(t *testing.T) {
		pmt := validPayment()
		assert.NoError(t, rules.Validate(pmt, withAddonContext()))
	})

	t.Run("invalid series", func(t *testing.T) {
		pmt := validPayment()

		pmt.Series = "SERIES-A"
		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "series format must be valid")
	})

	t.Run("invalid code", func(t *testing.T) {
		pmt := validPayment()

		pmt.Code = "ABCD"
		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "code format must be valid")
	})

	t.Run("valid full code", func(t *testing.T) {
		pmt := validPayment()

		pmt.Series = ""
		pmt.Code = "RG SERIES-A/123"
		assert.NoError(t, rules.Validate(pmt, withAddonContext()))
	})

	t.Run("missing extension", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = tax.Extensions{}

		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "'pt-saft-payment-type' extension is required")
	})

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		pmt := validPayment()
		pmt.Supplier.TaxID.Code = cbc.CodeEmpty

		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "supplier tax ID code is required")

		pmt.Supplier.TaxID = nil
		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "supplier tax ID is required")

		// pmt.Supplier = nil is caught by core GOBL rules (payment supplier is required)
	})

	t.Run("missing customer name", func(t *testing.T) {
		pmt := validPayment()
		pmt.Customer.Name = ""

		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "customer name is required when customer has tax ID code")

		pmt.Customer.TaxID.Code = ""
		assert.NoError(t, rules.Validate(pmt, withAddonContext()))

		pmt.Customer.TaxID = nil
		assert.NoError(t, rules.Validate(pmt, withAddonContext()))

		pmt.Customer = nil
		assert.NoError(t, rules.Validate(pmt, withAddonContext()))
	})

	t.Run("missing source billing", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = pmt.Ext.Delete(saft.ExtKeySource)
		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "'pt-saft-source' extension is required")
	})

	t.Run("source billing produced - no source doc ref required", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = tax.ExtensionsOf(tax.ExtMap{
			saft.ExtKeyPaymentType: saft.PaymentTypeOther,
			saft.ExtKeySource:      saft.SourceBillingProduced,
		})
		require.NoError(t, rules.Validate(pmt, withAddonContext()))
	})

	t.Run("source billing integrated - source doc ref required", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = tax.ExtensionsOf(tax.ExtMap{
			saft.ExtKeyPaymentType: saft.PaymentTypeOther,
			saft.ExtKeySource:      saft.SourceBillingIntegrated,
		})
		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "'pt-saft-source-ref' extension is required when source is not produced")

		// Add source doc ref - should pass
		pmt.Ext = pmt.Ext.Set(saft.ExtKeySourceRef, "RGM abc/00001")
		require.NoError(t, rules.Validate(pmt, withAddonContext()))
	})

	t.Run("source billing manual - source doc ref required", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = tax.ExtensionsOf(tax.ExtMap{
			saft.ExtKeyPaymentType: saft.PaymentTypeOther,
			saft.ExtKeySource:      saft.SourceBillingManual,
		})
		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "'pt-saft-source-ref' extension is required when source is not produced")

		// Add source doc ref - should pass
		pmt.Ext = pmt.Ext.Set(saft.ExtKeySourceRef, "RGD RG SERIESA/123")
		require.NoError(t, rules.Validate(pmt, withAddonContext()))
	})
}

func TestPaymentSourceRefFormatValidation(t *testing.T) {
	t.Run("missing source ref", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = pmt.Ext.Delete(saft.ExtKeySourceRef)
		require.NoError(t, rules.Validate(pmt, withAddonContext()))
	})

	t.Run("missing payment type", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = pmt.Ext.Delete(saft.ExtKeyPaymentType)
		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "'pt-saft-payment-type' extension is required")
	})

	t.Run("integrated document", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = pmt.Ext.Set(saft.ExtKeySource, saft.SourceBillingIntegrated)
		pmt.Ext = pmt.Ext.Set(saft.ExtKeySourceRef, "RGR abc/00001")
		require.NoError(t, rules.Validate(pmt, withAddonContext()))
	})

	tests := []struct {
		ref string
		err string
	}{
		{"RGM abc/00001", ""},
		{"RGD RG SERIESA/123", ""},
		{"RGR abc/00001", "source ref format is invalid"},
		{"RGM a/bc/00001", "source ref format is invalid"},
		{"RGDA RG abc/00001", "source ref format is invalid"},
		{"ABC abc/00001", "source ref format is invalid"},
		{"RGM RG abc/00001", "source ref format is invalid"},
		{"FRM abc/00001", "source ref format is invalid"},
		{"FRD RG SERIESA/123", "source ref format is invalid"},
		{"RGD FR SERIESA/123", "source ref format is invalid"},
	}

	for _, test := range tests {
		t.Run(test.ref, func(t *testing.T) {
			pmt := validPayment()
			pmt.Ext = pmt.Ext.Set(saft.ExtKeySource, saft.SourceBillingManual)
			pmt.Ext = pmt.Ext.Set(saft.ExtKeySourceRef, cbc.Code(test.ref))

			err := rules.Validate(pmt, withAddonContext())
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.err)
			}
		})
	}
}

func TestPaymentNormalization(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("general", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = tax.Extensions{}
		addon.Normalizer(pmt)
		assert.Equal(t, "RG", pmt.Ext.Get(saft.ExtKeyPaymentType).String())
	})

	t.Run("VAT cash", func(t *testing.T) {
		pmt := validPayment()
		pmt.SetTags("vat-cash")
		addon.Normalizer(pmt)
		assert.Equal(t, "RC", pmt.Ext.Get(saft.ExtKeyPaymentType).String())
	})

	t.Run("normalize payment with nil extensions", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = tax.Extensions{}

		addon.Normalizer(pmt)

		require.NotNil(t, pmt.Ext)
		assert.Equal(t, saft.SourceBillingProduced, pmt.Ext.Get(saft.ExtKeySource))
	})

	t.Run("normalize payment with missing source billing", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = pmt.Ext.Delete(saft.ExtKeySource)

		addon.Normalizer(pmt)

		assert.Equal(t, saft.SourceBillingProduced, pmt.Ext.Get(saft.ExtKeySource))
	})

	t.Run("normalize payment with existing source billing", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = pmt.Ext.Set(saft.ExtKeySource, saft.SourceBillingIntegrated)

		addon.Normalizer(pmt)

		assert.Equal(t, saft.SourceBillingIntegrated, pmt.Ext.Get(saft.ExtKeySource))
	})
}

func TestPaymentTotalValidation(t *testing.T) {
	t.Run("valid total amount", func(t *testing.T) {
		pmt := validPayment()
		pmt.Total = num.MakeAmount(100, 2)
		pmt.Methods[0].Amount = pmt.Total
		assert.NoError(t, rules.Validate(pmt, withAddonContext()))
	})

	t.Run("negative total amount", func(t *testing.T) {
		pmt := validPayment()
		pmt.Total = num.MakeAmount(-10, 2)
		assert.ErrorContains(t, rules.Validate(pmt, withAddonContext()), "must be no less than 0")
	})

	t.Run("nil total", func(t *testing.T) {
		pmt := validPayment()
		pmt.Total = num.Amount{}
		assert.NoError(t, rules.Validate(pmt, withAddonContext()))
	})
}
