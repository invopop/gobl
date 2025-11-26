package sii

import (
	"testing"

	"github.com/invopop/gobl/cbc"
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
			Key:      tax.KeyStandard,
			Rate:     tax.RateGeneral,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})
	t.Run("valid - no key", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
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
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})

	t.Run("exempt", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "E1", tc.Ext.Get(ExtKeyExempt).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass))
	})
	t.Run("export", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExport,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "02", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "E2", tc.Ext.Get(ExtKeyExempt).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass))
	})
	t.Run("surcharge", func(t *testing.T) {
		tc := &tax.Combo{
			Category:  tax.CategoryVAT,
			Rate:      tax.RateGeneral.With(es.TaxRateEquivalence),
			Percent:   num.NewPercentage(210, 3),
			Surcharge: num.NewPercentage(50, 3),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "18", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})
	t.Run("intra-community", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyIntraCommunity,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "E5", tc.Ext.Get(ExtKeyExempt).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass))
	})
	t.Run("outside scope", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "N2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})

	t.Run("outside scope N1", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
			Ext: tax.Extensions{
				ExtKeyOpClass: "N1",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "N1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})
	t.Run("reverse charge", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyReverseCharge,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "S2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})

	t.Run("foreign country", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  "FR",
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, cbc.Code("N2"), tc.Ext.Get(ExtKeyOpClass))
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})

	t.Run("with tax regime", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
			Ext: tax.Extensions{
				ExtKeyRegime: "03",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "03", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("with exempt code set", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyExempt: "E6",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E6", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, tax.KeyExempt, tc.Key)
	})
	t.Run("with export code set", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyExempt: "E2",
				ExtKeyRegime: "02",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E2", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, "02", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, tax.KeyExport, tc.Key)
	})
	t.Run("with reverse-charge", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyOpClass: "S2",
				ExtKeyRegime:  "01",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, tax.KeyReverseCharge, tc.Key)
	})
	t.Run("with outside-scope", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyOpClass: "N1",
				ExtKeyRegime:  "01",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "N1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, tax.KeyOutsideScope, tc.Key)
	})
	t.Run("with intra-community", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyExempt: "E5",
				ExtKeyRegime: "01",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E5", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, tax.KeyIntraCommunity, tc.Key)
	})
}

func TestValidateTaxCombo(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
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
			Rate:     tax.RateGeneral,
		}
		err := validateTaxCombo(tc)
		assert.NoError(t, err)
	})

	t.Run("exempt with valid reason", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
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

