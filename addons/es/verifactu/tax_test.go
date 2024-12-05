package verifactu

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxCombo(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyTaxRegime).String())
	})

	t.Run("valid with country", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  l10n.ES.Tax(),
			Category: tax.CategoryVAT,
			Rate:     tax.RateSuperReduced,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyTaxRegime).String())
	})

	t.Run("undefined rate", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt,
		}
		normalizeTaxCombo(tc)
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
	})

	t.Run("foreign country", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  "FR",
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
		}
		normalizeTaxCombo(tc)
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
	})

	t.Run("with tax regime", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
			Ext: tax.Extensions{
				ExtKeyTaxRegime: "03",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "03", tc.Ext.Get(ExtKeyTaxRegime).String())
	})
}
