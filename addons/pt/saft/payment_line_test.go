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

		// First check that the exemption extension is required
		err := addon.Validator(pl)
		assert.ErrorContains(t, err, "pt-saft-exemption: required")

		// Then add the exemption extension
		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyExemption] = "M01"

		// Now it should fail because the exemption note is missing
		err = addon.Validator(pl)
		assert.ErrorContains(t, err, "notes: missing exemption note for code M01")

		// Add the required note
		pl.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  saft.ExtKeyExemption,
				Code: "M01",
				Text: "Artigo 13.º do CIVA",
			},
		}
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

	t.Run("payment line with no notes", func(t *testing.T) {
		pl := validPaymentLine()
		assert.NoError(t, addon.Validator(pl))
	})

	t.Run("payment line with valid exemption note", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyTaxRate] = saft.TaxRateExempt
		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyExemption] = "M04"
		pl.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  saft.ExtKeyExemption,
				Code: "M04",
				Text: "Artigo 13.º do CIVA",
			},
		}
		assert.NoError(t, addon.Validator(pl))
	})

	t.Run("payment line missing exemption note", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyTaxRate] = saft.TaxRateExempt
		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyExemption] = "M05"
		// No notes added
		err := addon.Validator(pl)
		assert.ErrorContains(t, err, "notes: missing exemption note for code M05")
	})

	t.Run("payment line with unexpected exemption note", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  saft.ExtKeyExemption,
				Code: "M04",
				Text: "Artigo 13.º do CIVA",
			},
		}
		err := addon.Validator(pl)
		assert.ErrorContains(t, err, "notes: (0: unexpected exemption note")
	})

	t.Run("payment line with mismatched exemption note code", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyTaxRate] = saft.TaxRateExempt
		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyExemption] = "M03"
		pl.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  saft.ExtKeyExemption,
				Code: "M01", // Different code than extension
				Text: "Artigo 13.º do CIVA",
			},
		}
		err := addon.Validator(pl)
		assert.ErrorContains(t, err, "notes: (0: note code M01 must match extension M03)")
	})

	t.Run("payment line with too many exemption notes", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyTaxRate] = saft.TaxRateExempt
		pl.Tax.Categories[0].Rates[0].Ext[saft.ExtKeyExemption] = "M02"
		pl.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  saft.ExtKeyExemption,
				Code: "M02",
				Text: "Artigo 13.º do CIVA",
			},
			{
				Key:  org.NoteKeyLegal,
				Src:  saft.ExtKeyExemption,
				Code: "M02",
				Text: "Duplicate exemption note",
			},
		}
		err := addon.Validator(pl)
		assert.ErrorContains(t, err, "notes: (1: too many exemption notes)")
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
