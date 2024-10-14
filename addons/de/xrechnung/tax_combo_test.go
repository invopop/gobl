package xrechnung_test

import (
	"testing"

	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxComboValidation(t *testing.T) {
	t.Run("standard VAT rate", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
		}
		assert.NoError(t, xrechnung.ValidateTaxCombo(combo))
		assert.Equal(t, "S", combo.Ext["de-xrechnung-tax-rate"])
	})

	t.Run("missing rate", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
		}
		err := xrechnung.ValidateTaxCombo(combo)
		assert.EqualError(t, err, "VAT category rate is required")
	})

}
