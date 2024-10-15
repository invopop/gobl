package xrechnung_test

import (
	"testing"

	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxComboValidation(t *testing.T) {
	t.Run("standard VAT rate", func(t *testing.T) {
		p := num.MakePercentage(19, 2)
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
			Percent:  &p,
		}
		xrechnung.NormalizeTaxCombo(c)
		assert.NoError(t, xrechnung.ValidateTaxCombo(c))
		assert.Equal(t, "S", c.Ext[xrechnung.ExtKeyTaxRate].String())
		assert.Equal(t, "19%", c.Percent.String())
	})

	t.Run("missing rate", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
		}
		err := xrechnung.ValidateTaxCombo(c)
		assert.EqualError(t, err, "percent: cannot be blank.")
	})

}
