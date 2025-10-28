package tbai

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
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("valid - no key", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		}
		normalizeTaxCombo(tc)
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("valid with country", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  l10n.ES.Tax(),
			Category: tax.CategoryVAT,
			Rate:     tax.RateSuperReduced,
		}
		normalizeTaxCombo(tc)
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("exempt", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E1", tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("export", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExport,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E2", tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("surcharge", func(t *testing.T) {
		tc := &tax.Combo{
			Category:  tax.CategoryVAT,
			Rate:      tax.RateGeneral.With(es.TaxRateEquivalence),
			Percent:   num.NewPercentage(210, 3),
			Surcharge: num.NewPercentage(50, 3),
		}
		normalizeTaxCombo(tc)
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("intra-community", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyIntraCommunity,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E5", tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("outside scope", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "OT", tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("reverse charge", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyReverseCharge,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S2", tc.Ext.Get(ExtKeyExempt).String())
	})

	t.Run("foreign country", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  "FR",
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, cbc.Code("RL"), tc.Ext.Get(ExtKeyExempt))
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
		assert.Equal(t, tax.KeyExempt, tc.Key)
	})

	t.Run("with export code set", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyExempt: "E2",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E2", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, tax.KeyExport, tc.Key)
	})

	t.Run("with reverse-charge", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyExempt: "S2",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S2", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, tax.KeyReverseCharge, tc.Key)
	})

	t.Run("with outside-scope", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyExempt: "OT",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "OT", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, tax.KeyOutsideScope, tc.Key)
	})

	t.Run("with intra-community", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyExempt: "E5",
			},
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E5", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, tax.KeyIntraCommunity, tc.Key)
	})

	t.Run("with standard", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				ExtKeyExempt: "S1",
			},
		}
		normalizeTaxCombo(tc)
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, tax.KeyStandard, tc.Key)
	})
}
