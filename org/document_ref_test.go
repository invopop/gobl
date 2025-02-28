package org_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestDocumentRefValidation(t *testing.T) {
	dr := new(org.DocumentRef)
	dr.Code = "FOO"
	dr.IssueDate = cal.NewDate(2022, 11, 6)

	err := dr.Validate()
	assert.NoError(t, err)
}

func TestDocumentRefNormalize(t *testing.T) {
	t.Run("basic", func(t *testing.T) {

		dr := &org.DocumentRef{
			Code: " Foo ",
			Ext: tax.Extensions{
				"fooo": "",
			},
		}
		dr.Normalize(nil)
		assert.Equal(t, "Foo", dr.Code.String())
		assert.Empty(t, dr.Ext)
	})
	t.Run("nil", func(t *testing.T) {
		var dr *org.DocumentRef
		assert.NotPanics(t, func() {
			dr.Normalize(nil)
		})
	})

}

func TestDocumentRefCalculate(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		dr := &org.DocumentRef{
			IssueDate: cal.NewDate(2022, 11, 6),
		}
		assert.NotPanics(t, func() {
			dr.Calculate(currency.EUR, tax.RoundingRuleSumThenRound)
		})
	})
	t.Run("with tax", func(t *testing.T) {
		dr := &org.DocumentRef{
			Code: "FOO",
			Tax: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code: tax.CategoryVAT,
						Rates: []*tax.RateTotal{
							{
								Base:    num.MakeAmount(1000, 2),
								Percent: num.NewPercentage(21, 2),
							},
						},
					},
				},
			},
		}
		assert.NotPanics(t, func() {
			dr.Calculate(currency.EUR, tax.RoundingRuleSumThenRound)
		})
		assert.Equal(t, "2.10", dr.Tax.Sum.String())
	})
}
