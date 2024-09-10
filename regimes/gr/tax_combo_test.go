package gr_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/gr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxCombo(t *testing.T) {
	greece := gr.New()
	ctx := greece.WithContext(context.Background())

	t.Run("vat category presence", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(4, 2),
		}
		err := tc.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "ext: (gr-mydata-vat-cat: required.)")
	})

	t.Run("exemption presence", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext:      tax.Extensions{"gr-mydata-vat-cat": "1"},
		}
		err := tc.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "gr-mydata-exemption: required")
	})

	t.Run("income classification keys presence", func(t *testing.T) {
		tx := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(24, 2),
			Ext: tax.Extensions{
				gr.ExtKeyMyDATAVATCat:    "1",
				gr.ExtKeyMyDATAIncomeCat: "category1_1",
			},
		}
		err := tx.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "gr-mydata-income-type: required")

		tx.Ext[gr.ExtKeyMyDATAIncomeType] = "E3_106"
		delete(tx.Ext, gr.ExtKeyMyDATAIncomeCat)

		err = tx.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "gr-mydata-income-cat: required")

		tx.Ext[gr.ExtKeyMyDATAIncomeCat] = "category1_1"

		err = tx.ValidateWithContext(ctx)
		assert.NoError(t, err)
	})
}
