package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPaymentLineValidation(t *testing.T) {
	t.Run("missing line document", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Document = nil

		assert.ErrorContains(t, rules.Validate(pl, withAddonContext()), "cannot be blank")

		pl = nil
		assert.NoError(t, rules.Validate(pl, withAddonContext()))
	})

	t.Run("missing line document issue date", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Document.IssueDate = nil

		assert.ErrorContains(t, rules.Validate(pl, withAddonContext()), "cannot be blank")
	})

	t.Run("missing VAT category in line tax", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax = nil

		assert.ErrorContains(t, rules.Validate(pl, withAddonContext()), "cannot be blank")

		pl.Tax = new(tax.Total)
		assert.ErrorContains(t, rules.Validate(pl, withAddonContext()), "missing category VAT")
	})

	t.Run("missing line tax required extensions", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext = tax.Extensions{}

		err := rules.Validate(pl, withAddonContext())
		assert.ErrorContains(t, err, "region and tax rate are required")

		pl.Tax.Categories[0].Rates[0] = nil
		err = rules.Validate(pl, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("missing line tax exemption", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext = pl.Tax.Categories[0].Rates[0].Ext.Set(saft.ExtKeyTaxRate, saft.TaxRateExempt)

		// First check that the exemption extension is required
		err := rules.Validate(pl, withAddonContext())
		assert.ErrorContains(t, err, "exemption is required when tax rate is exempt")

		// Then add the exemption extension
		pl.Tax.Categories[0].Rates[0].Ext = pl.Tax.Categories[0].Rates[0].Ext.Set(saft.ExtKeyExemption, "M01")

		// Now it should fail because the exemption note is missing
		err = rules.Validate(pl, withAddonContext())
		assert.ErrorContains(t, err, "exemption notes invalid")

		// Add the required note
		pl.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  saft.ExtKeyExemption,
				Code: "M01",
				Text: "Artigo 13.º do CIVA",
			},
		}
		assert.NoError(t, rules.Validate(pl, withAddonContext()))
	})

	t.Run("nil tax category", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories = append(pl.Tax.Categories, nil)
		assert.NoError(t, rules.Validate(pl, withAddonContext()))
	})

	t.Run("too many VAT rates", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates = append(pl.Tax.Categories[0].Rates, &tax.RateTotal{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				pt.ExtKeyRegion:    "PT",
				saft.ExtKeyTaxRate: "INT",
			}),
		})

		err := rules.Validate(pl, withAddonContext())
		assert.ErrorContains(t, err, "only one rate allowed per line")
		// Note: the error format is now from the rules framework
	})

	t.Run("payment line with no notes", func(t *testing.T) {
		pl := validPaymentLine()
		assert.NoError(t, rules.Validate(pl, withAddonContext()))
	})

	t.Run("payment line with valid exemption note", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext = pl.Tax.Categories[0].Rates[0].Ext.Set(saft.ExtKeyTaxRate, saft.TaxRateExempt)
		pl.Tax.Categories[0].Rates[0].Ext = pl.Tax.Categories[0].Rates[0].Ext.Set(saft.ExtKeyExemption, "M04")
		pl.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  saft.ExtKeyExemption,
				Code: "M04",
				Text: "Artigo 13.º do CIVA",
			},
		}
		assert.NoError(t, rules.Validate(pl, withAddonContext()))
	})

	t.Run("payment line missing exemption note", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext = pl.Tax.Categories[0].Rates[0].Ext.Set(saft.ExtKeyTaxRate, saft.TaxRateExempt)
		pl.Tax.Categories[0].Rates[0].Ext = pl.Tax.Categories[0].Rates[0].Ext.Set(saft.ExtKeyExemption, "M05")
		// No notes added
		err := rules.Validate(pl, withAddonContext())
		assert.ErrorContains(t, err, "exemption notes invalid")
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
		err := rules.Validate(pl, withAddonContext())
		assert.ErrorContains(t, err, "exemption notes invalid")
	})

	t.Run("payment line with mismatched exemption note code", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext = pl.Tax.Categories[0].Rates[0].Ext.Set(saft.ExtKeyTaxRate, saft.TaxRateExempt)
		pl.Tax.Categories[0].Rates[0].Ext = pl.Tax.Categories[0].Rates[0].Ext.Set(saft.ExtKeyExemption, "M03")
		pl.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  saft.ExtKeyExemption,
				Code: "M01",
				Text: "Artigo 13.º do CIVA",
			},
		}
		err := rules.Validate(pl, withAddonContext())
		assert.ErrorContains(t, err, "exemption notes invalid")
	})

	t.Run("payment line with too many exemption notes", func(t *testing.T) {
		pl := validPaymentLine()
		pl.Tax.Categories[0].Rates[0].Ext = pl.Tax.Categories[0].Rates[0].Ext.Set(saft.ExtKeyTaxRate, saft.TaxRateExempt)
		pl.Tax.Categories[0].Rates[0].Ext = pl.Tax.Categories[0].Rates[0].Ext.Set(saft.ExtKeyExemption, "M02")
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
		err := rules.Validate(pl, withAddonContext())
		assert.ErrorContains(t, err, "exemption notes invalid")
	})
}

func validPaymentLine() *bill.PaymentLine {
	return &bill.PaymentLine{
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
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								pt.ExtKeyRegion:    "PT",
								saft.ExtKeyTaxRate: "NOR",
							}),
						},
					},
				},
			},
		},
	}
}
