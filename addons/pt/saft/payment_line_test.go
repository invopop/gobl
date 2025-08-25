package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentLineValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)
	require.NotNil(t, addon)

	t.Run("missing line document", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Document = nil

		assert.ErrorContains(t, addon.Validator(pl), "document: cannot be blank")

		pl = nil
		assert.NoError(t, addon.Validator(pl))
	})

	t.Run("missing line document issue date", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Document.IssueDate = nil

		assert.ErrorContains(t, addon.Validator(pl), "document: (issue_date: cannot be blank")
	})

	t.Run("missing VAT category in line tax", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax = nil

		assert.ErrorContains(t, addon.Validator(pl), "tax: cannot be blank")

		pl.Tax = new(tax.Total)
		assert.ErrorContains(t, addon.Validator(pl), "tax: missing category VAT")
	})

	t.Run("missing line tax required extensions", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext = nil

		err := addon.Validator(pl)
		assert.ErrorContains(t, err, "pt-region: required")
		assert.ErrorContains(t, err, "pt-saft-tax-rate: required")

		pl.Tax.Categories[0].Rates[0] = nil
		assert.NoError(t, addon.Validator(pl))
	})

	t.Run("missing line tax exemption", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyTaxRate] = saft.TaxRateExempt

		assert.ErrorContains(t, addon.Validator(pl), "pt-saft-exemption: required")

		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyExemption] = "M01"
		assert.NoError(t, addon.Validator(pl))
	})

	t.Run("nil tax category", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories = append(pl.Tax.Categories, nil)
		assert.NoError(t, addon.Validator(pl))
	})

	t.Run("too many VAT rates", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates = append(pl.Tax.Categories[0].Rates, &tax.RateTotal{
			Ext: tax.Extensions{
				pt.ExtKeyRegion:    "PT",
				saft.ExtKeyTaxRate: "INT",
			},
		})

		err := addon.Validator(pl)
		assert.ErrorContains(t, err, "tax: (categories: (0: (rates: only one rate allowed per line")
	})
}

func validPaymentLine() *bill.PaymentLine {
	return &bill.PaymentLine{
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
	}
}
