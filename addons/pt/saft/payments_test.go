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
				Code:    "123456789",
			},
		},
		Customer: &org.Party{
			Name: "Customer Name",
			TaxID: &tax.Identity{
				Country: "PT",
				Code:    "987654321",
			},
		},
		Ext: tax.Extensions{
			saft.ExtKeyPaymentType:   saft.PaymentTypeOther,
			saft.ExtKeySourceBilling: saft.SourceBillingProduced,
		},
		Series:    "RG SERIES-A",
		Code:      "123",
		IssueDate: cal.MakeDate(2024, 3, 10),
		Lines: []*bill.PaymentLine{
			{
				Document: &org.DocumentRef{
					IssueDate: cal.NewDate(2024, 3, 1),
				},
				Amount: num.MakeAmount(100, 2),
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: tax.CategoryVAT,
							Rates: []*tax.RateTotal{
								{
									Ext: tax.Extensions{
										pt.ExtKeyRegion:    "PT",
										saft.ExtKeyTaxRate: "NOR",
									},
								},
							},
						},
					},
				},
			},
		},
		Method: &pay.Instructions{
			Key: "credit-transfer",
		},
	}
}

func TestPaymentValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("valid payment", func(t *testing.T) {
		pmt := validPayment()
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("invalid series", func(t *testing.T) {
		pmt := validPayment()

		pmt.Series = "SERIES-A"
		assert.ErrorContains(t, addon.Validator(pmt), "series: must start with 'RG '")
	})

	t.Run("invalid code", func(t *testing.T) {
		pmt := validPayment()

		pmt.Code = "ABCD"
		assert.ErrorContains(t, addon.Validator(pmt), "code: must be in a valid format")
	})

	t.Run("valid full code", func(t *testing.T) {
		pmt := validPayment()

		pmt.Series = ""
		pmt.Code = "RG SERIES-A/123"
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("missing extension", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = nil

		assert.ErrorContains(t, addon.Validator(pmt), "ext: (pt-saft-payment-type: required")
	})

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		pmt := validPayment()
		pmt.Supplier.TaxID.Code = cbc.CodeEmpty

		assert.ErrorContains(t, addon.Validator(pmt), "supplier: (tax_id: (code: cannot be blank")

		pmt.Supplier.TaxID = nil
		assert.ErrorContains(t, addon.Validator(pmt), "supplier: (tax_id: cannot be blank.")

		pmt.Supplier = nil
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("missing customer name", func(t *testing.T) {
		pmt := validPayment()
		pmt.Customer.Name = ""

		assert.ErrorContains(t, addon.Validator(pmt), "customer: (name: cannot be blank")

		pmt.Customer.TaxID.Code = ""
		assert.NoError(t, addon.Validator(pmt))

		pmt.Customer.TaxID = nil
		assert.NoError(t, addon.Validator(pmt))

		pmt.Customer = nil
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("missing line document", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Document = nil

		assert.ErrorContains(t, addon.Validator(pmt), "lines: (0: (document: cannot be blank.).)")

		pmt.Lines[0] = nil
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("missing line document issue date", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Document.IssueDate = nil

		assert.ErrorContains(t, addon.Validator(pmt), "lines: (0: (document: (issue_date: cannot be blank")
	})

	t.Run("missing VAT category in line tax", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Tax = nil

		assert.ErrorContains(t, addon.Validator(pmt), "lines: (0: (tax: cannot be blank")

		pmt.Lines[0].Tax = new(tax.Total)
		assert.ErrorContains(t, addon.Validator(pmt), "lines: (0: (tax: missing category VAT")
	})

	t.Run("missing line tax required extensions", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Tax.Categories[0].Rates[0].Ext = nil

		err := addon.Validator(pmt)
		assert.ErrorContains(t, err, "pt-region: required")
		assert.ErrorContains(t, err, "pt-saft-tax-rate: required")

		pmt.Lines[0].Tax.Categories[0].Rates[0] = nil
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("missing source billing", func(t *testing.T) {
		pmt := validPayment()
		delete(pmt.Ext, saft.ExtKeySourceBilling)
		assert.ErrorContains(t, addon.Validator(pmt), "ext: (pt-saft-source-billing: required")
	})

	t.Run("source billing produced - no source doc ref required", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = tax.Extensions{
			saft.ExtKeyPaymentType:   saft.PaymentTypeOther,
			saft.ExtKeySourceBilling: saft.SourceBillingProduced,
		}
		require.NoError(t, addon.Validator(pmt))
	})

	t.Run("source billing integrated - source doc ref required", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = tax.Extensions{
			saft.ExtKeyPaymentType:   saft.PaymentTypeOther,
			saft.ExtKeySourceBilling: saft.SourceBillingIntegrated,
		}
		assert.ErrorContains(t, addon.Validator(pmt), "ext: (pt-saft-source-ref: required")

		// Add source doc ref - should pass
		pmt.Ext[saft.ExtKeySourceRef] = "RGM abc/00001"
		require.NoError(t, addon.Validator(pmt))
	})

	t.Run("source billing manual - source doc ref required", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = tax.Extensions{
			saft.ExtKeyPaymentType:   saft.PaymentTypeOther,
			saft.ExtKeySourceBilling: saft.SourceBillingManual,
		}
		assert.ErrorContains(t, addon.Validator(pmt), "ext: (pt-saft-source-ref: required")

		// Add source doc ref - should pass
		pmt.Ext[saft.ExtKeySourceRef] = "RGD RG SERIESA/123"
		require.NoError(t, addon.Validator(pmt))
	})

	t.Run("nil tax category", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Tax.Categories = append(pmt.Lines[0].Tax.Categories, nil)
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("too many VAT rates", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Tax.Categories[0].Rates = append(pmt.Lines[0].Tax.Categories[0].Rates, &tax.RateTotal{
			Ext: tax.Extensions{
				pt.ExtKeyRegion:    "PT",
				saft.ExtKeyTaxRate: "INT",
			},
		})

		err := addon.Validator(pmt)
		assert.ErrorContains(t, err, "lines: (0: (tax: (categories: (0: (rates: only one rate allowed per line")
	})

}

func TestPaymentSourceRefFormatValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("missing source ref", func(t *testing.T) {
		pmt := validPayment()
		delete(pmt.Ext, saft.ExtKeySourceRef)
		require.NoError(t, addon.Validator(pmt))
	})

	t.Run("missing payment type", func(t *testing.T) {
		pmt := validPayment()
		delete(pmt.Ext, saft.ExtKeyPaymentType)
		assert.ErrorContains(t, addon.Validator(pmt), "ext: (pt-saft-payment-type: required")
	})

	t.Run("integrated document", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext[saft.ExtKeySourceBilling] = saft.SourceBillingIntegrated
		pmt.Ext[saft.ExtKeySourceRef] = "RGR abc/00001"
		require.NoError(t, addon.Validator(pmt))
	})

	tests := []struct {
		ref string
		err string
	}{
		{"", ""},
		{"RGM abc/00001", ""},
		{"RGD RG SERIESA/123", ""},
		{"RGR abc/00001", "must be in valid format"},
		{"RGM a/bc/00001", "must be in valid format"},
		{"RGDA RG abc/00001", "must be in valid format"},
		{"ABC abc/00001", "must be in valid format"},
		{"RGM RG abc/00001", "must be in valid format"},
		{"FRM abc/00001", "must start with the document type 'RG' not 'FR'"},
		{"FRD RG SERIESA/123", "must start with the document type 'RG' not 'FR'"},
		{"RGD FR SERIESA/123", "must refer to an original document 'RG' not 'FR'"},
	}

	for _, test := range tests {
		t.Run(test.ref, func(t *testing.T) {
			pmt := validPayment()
			pmt.Ext[saft.ExtKeySourceBilling] = saft.SourceBillingManual
			pmt.Ext[saft.ExtKeySourceRef] = cbc.Code(test.ref)

			err := addon.Validator(pmt)
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
		pmt.Ext = nil
		addon.Normalizer(pmt)
		assert.Equal(t, "RG", pmt.Ext[saft.ExtKeyPaymentType].String())
	})

	t.Run("VAT cash", func(t *testing.T) {
		pmt := validPayment()
		pmt.SetTags("vat-cash")
		addon.Normalizer(pmt)
		assert.Equal(t, "RC", pmt.Ext[saft.ExtKeyPaymentType].String())
	})

	t.Run("normalize payment with nil extensions", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = nil

		addon.Normalizer(pmt)

		require.NotNil(t, pmt.Ext)
		assert.Equal(t, saft.SourceBillingProduced, pmt.Ext[saft.ExtKeySourceBilling])
	})

	t.Run("normalize payment with missing source billing", func(t *testing.T) {
		pmt := validPayment()
		delete(pmt.Ext, saft.ExtKeySourceBilling)

		addon.Normalizer(pmt)

		assert.Equal(t, saft.SourceBillingProduced, pmt.Ext[saft.ExtKeySourceBilling])
	})

	t.Run("normalize payment with existing source billing", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext[saft.ExtKeySourceBilling] = saft.SourceBillingIntegrated

		addon.Normalizer(pmt)

		assert.Equal(t, saft.SourceBillingIntegrated, pmt.Ext[saft.ExtKeySourceBilling])
	})
}

func TestPaymentTotalValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("valid total amount", func(t *testing.T) {
		pmt := validPayment()
		pmt.Total = num.MakeAmount(100, 2)
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("negative total amount", func(t *testing.T) {
		pmt := validPayment()
		pmt.Total = num.MakeAmount(-10, 2)
		assert.ErrorContains(t, addon.Validator(pmt), "total: must be no less than 0")
	})

	t.Run("nil total", func(t *testing.T) {
		pmt := validPayment()
		pmt.Total = num.Amount{}
		assert.NoError(t, addon.Validator(pmt))
	})
}
