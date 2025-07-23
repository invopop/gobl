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

	// Test all tax rate mappings to operation classes
	t.Run("reduced rate maps to S1", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateReduced,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("zero rate maps to S1", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateZero,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("zero rate with location tag maps to N2", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateZero.With(tax.TagForeignVAT),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "N2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "08", tc.Ext.Get(ExtKeyRegime).String())
	})

	// Test reverse charge scenarios
	t.Run("reverse charge standard rate maps to S2", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard.With(tax.TagReverseCharge),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("reverse charge reduced rate maps to S2", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateReduced.With(tax.TagReverseCharge),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("reverse charge super reduced rate maps to S2", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateSuperReduced.With(tax.TagReverseCharge),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})

	// Test simplified regime scenarios
	t.Run("simplified tag sets regime to 20", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard.With(tax.TagSimplified),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "20", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("simplified with reduced rate", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateReduced.With(tax.TagSimplified),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "20", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("simplified with zero rate", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateZero.With(tax.TagSimplified),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "20", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("simplified with exempt", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt.With(tax.TagSimplified),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E1", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, "20", tc.Ext.Get(ExtKeyRegime).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
	})

	t.Run("simplified export with EEA", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt.With(tax.TagExport).With(tax.TagEEA).With(tax.TagSimplified),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E5", tc.Ext.Get(ExtKeyExempt).String())
		// Simplified takes precedence over export for regime
		assert.Equal(t, "20", tc.Ext.Get(ExtKeyRegime).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
	})

	t.Run("location tag with IGIC", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIGIC,
			Rate:     tax.RateZero.With(tax.TagForeignVAT),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "N2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "08", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("simplified tag precedence - simplified wins over export", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt.With(tax.TagExport).With(tax.TagSimplified),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E2", tc.Ext.Get(ExtKeyExempt).String())
		// Simplified should set regime to 20, overriding export's 02
		assert.Equal(t, "20", tc.Ext.Get(ExtKeyRegime).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
	})

	t.Run("surcharge takes precedence over simplified", func(t *testing.T) {
		tc := &tax.Combo{
			Category:  tax.CategoryVAT,
			Rate:      tax.RateStandard.With(tax.TagSimplified),
			Surcharge: num.NewPercentage(50, 3),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		// Surcharge should set regime to 18, not simplified's 20
		assert.Equal(t, "18", tc.Ext.Get(ExtKeyRegime).String())
	})

	// Test exempt scenarios
	t.Run("basic exempt maps to E1", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E1", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
	})

	// Test IGIC category
	t.Run("IGIC category with standard rate", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIGIC,
			Rate:     tax.RateStandard,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("IGIC category with exempt", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIGIC,
			Rate:     tax.RateExempt,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E1", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
	})

	t.Run("IGIC with export and EEA", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIGIC,
			Rate:     tax.RateExempt.With(tax.TagExport).With(tax.TagEEA),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "02", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "E5", tc.Ext.Get(ExtKeyExempt).String())
	})

	// Test that operation class overrides exempt
	t.Run("operation class removes exempt", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
			Ext: tax.Extensions{
				ExtKeyExempt: "E1", // This should be removed
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt).String())
	})

	// Test non-VAT/IGIC categories are ignored
	t.Run("non-VAT category ignored", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryGST,
			Rate:     tax.RateStandard,
		}
		normalizeTaxCombo(tc)
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("empty rate with VAT category", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			// Rate field not set (empty)
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt).String())
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

	// Additional validation tests
	t.Run("missing regime for VAT", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
			Ext: tax.Extensions{
				ExtKeyOpClass: "S1",
			},
		}
		err := validateTaxCombo(tc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "regime")
	})

	t.Run("missing regime for IGIC", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIGIC,
			Rate:     tax.RateStandard,
			Ext: tax.Extensions{
				ExtKeyOpClass: "S1",
			},
		}
		err := validateTaxCombo(tc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "regime")
	})

	t.Run("missing operation class when taxed", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
			Percent:  num.NewPercentage(210, 3),
			Ext: tax.Extensions{
				ExtKeyRegime: "01",
			},
		}
		err := validateTaxCombo(tc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "es-verifactu-op-class")
	})

	t.Run("cannot have both operation class and exempt", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
			Ext: tax.Extensions{
				ExtKeyRegime:  "01",
				ExtKeyOpClass: "S1",
				ExtKeyExempt:  "E1",
			},
		}
		err := validateTaxCombo(tc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exempt")
	})

	t.Run("allows E3 exemption code with regime 02", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyRegime: "02",
				ExtKeyExempt: "E3",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("allows E5 exemption code with any regime", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyRegime: "01",
				ExtKeyExempt: "E5",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("valid reverse charge with S2", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard.With(tax.TagReverseCharge),
			Percent:  num.NewPercentage(210, 3),
			Ext: tax.Extensions{
				ExtKeyRegime:  "01",
				ExtKeyOpClass: "S2",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("valid surcharge scenario", func(t *testing.T) {
		tc := &tax.Combo{
			Category:  tax.CategoryVAT,
			Rate:      tax.RateStandard,
			Percent:   num.NewPercentage(210, 3),
			Surcharge: num.NewPercentage(50, 3),
			Ext: tax.Extensions{
				ExtKeyRegime:  "18",
				ExtKeyOpClass: "S1",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("exempt without percent is valid", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt,
			Ext: tax.Extensions{
				ExtKeyRegime: "01",
				ExtKeyExempt: "E1",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("valid N2 operation class", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateZero.With(tax.TagForeignVAT),
			Percent:  num.NewPercentage(0, 3),
			Ext: tax.Extensions{
				ExtKeyRegime:  "01",
				ExtKeyOpClass: "N2",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("valid simplified regime 20", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard.With(tax.TagSimplified),
			Percent:  num.NewPercentage(210, 3),
			Ext: tax.Extensions{
				ExtKeyRegime:  "20",
				ExtKeyOpClass: "S1",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("valid simplified exempt scenario", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt.With(tax.TagSimplified),
			Ext: tax.Extensions{
				ExtKeyRegime: "20",
				ExtKeyExempt: "E1",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("valid N2 with IGIC", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIGIC,
			Rate:     tax.RateZero.With(tax.TagForeignVAT),
			Percent:  num.NewPercentage(0, 3),
			Ext: tax.Extensions{
				ExtKeyRegime:  "01",
				ExtKeyOpClass: "N2",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("allows E2 and E3 exemption codes with regime 20", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyRegime: "20",
				ExtKeyExempt: "E2",
			},
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)

		tc.Ext[ExtKeyExempt] = "E3"
		err = validateTaxCombo(tc)
		assert.NoError(t, err)
	})
}
