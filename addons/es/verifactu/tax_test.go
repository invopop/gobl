package verifactu

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/es"
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
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("valid with country", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  l10n.ES.Tax(),
			Category: tax.CategoryVAT,
			Rate:     tax.RateSuperReduced,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("exempt export", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt.With(tax.TagExport).With(tax.TagEEA),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "02", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "E5", tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("exempt export with non-EU customer", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt.With(tax.TagExport),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "02", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "E2", tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("surcharge", func(t *testing.T) {
		tc := &tax.Combo{
			Category:  tax.CategoryVAT,
			Rate:      tax.RateStandard.With(es.TaxRateEquivalence),
			Percent:   num.NewPercentage(210, 3),
			Surcharge: num.NewPercentage(50, 3),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "18", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
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
				ExtKeyRegime: "03",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "03", tc.Ext.Get(ExtKeyRegime).String())
	})
}

func TestValidateTaxCombo(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
			Ext: tax.Extensions{
				ExtKeyOpClass: "S1",
				ExtKeyRegime:  "01",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("not in category", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryGST,
			Rate:     tax.RateStandard,
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("exempt with valid reason", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyRegime: "01",
				ExtKeyExempt: "E1",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("excludes E2 exemption code with regime 01", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyRegime: "01",
				ExtKeyExempt: "E2",
			},
		}
		err := validateTaxCombo(tc)
		assert.ErrorContains(t, err, "E2")
	})

	t.Run("excludes E3 exemption code with regime 01", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyRegime: "01",
				ExtKeyExempt: "E3",
			},
		}
		err := validateTaxCombo(tc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "E3")
	})

	t.Run("allows E2 exemption code with non-01 regime", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyRegime: "02",
				ExtKeyExempt: "E2",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("excludes E2 exemption code with regime 01 and IGIC category", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIGIC,
			Ext: tax.Extensions{
				ExtKeyRegime: "01",
				ExtKeyExempt: "E2",
			},
		}
		err := validateTaxCombo(tc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "E2")
	})
}
