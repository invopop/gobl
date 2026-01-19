package saft_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxComboNormalize(t *testing.T) {
	ad := tax.AddonForKey(saft.V1)
	// most tests here done in TestInvoice
	t.Run("standard", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		}
		ad.Normalizer(combo)
		assert.Equal(t, "NOR", combo.Ext[saft.ExtKeyTaxRate].String())
	})

	t.Run("standard with exempt reason", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Rate:     tax.RateGeneral,
			Ext: tax.Extensions{
				saft.ExtKeyExemption: "M01",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "NOR", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Empty(t, combo.Ext[saft.ExtKeyExemption])
	})

	t.Run("reverse-charge", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyReverseCharge,
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M40", combo.Ext[saft.ExtKeyExemption].String())
	})

	t.Run("outside-scope", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M99", combo.Ext[saft.ExtKeyExemption].String())
	})
	t.Run("intra-community", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyIntraCommunity,
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M16", combo.Ext[saft.ExtKeyExemption].String())
	})

	t.Run("unsupported", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateSuperReduced,
		}
		ad.Normalizer(combo)
		assert.Empty(t, combo.Ext)
	})

	t.Run("nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var c *tax.Combo
			ad.Normalizer(c)
		})
	})

	t.Run("no rate", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
		}
		ad.Normalizer(combo)
		assert.Empty(t, combo.Ext)
	})

	t.Run("no rate, no override ext", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(21, 3),
			Ext: tax.Extensions{
				saft.ExtKeyTaxRate: "NOR",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "NOR", combo.Ext[saft.ExtKeyTaxRate].String())
	})

	t.Run("foreign rate", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Country:  l10n.EL.Tax(),
			Percent:  num.NewPercentage(24, 3),
			Ext: tax.Extensions{
				saft.ExtKeyTaxRate: "NOR",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "OUT", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.False(t, combo.Ext.Has(saft.ExtKeyExemption))
	})

	t.Run("reverse map exemption M30", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				saft.ExtKeyExemption: "M30",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M30", combo.Ext[saft.ExtKeyExemption].String())
		assert.Equal(t, tax.KeyReverseCharge, combo.Key)
	})

	t.Run("reverse map exemption M05", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				saft.ExtKeyExemption: "M05",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M05", combo.Ext[saft.ExtKeyExemption].String())
		assert.Equal(t, tax.KeyExport, combo.Key)
	})

	t.Run("reverse map exemption M16", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				saft.ExtKeyExemption: "M16",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M16", combo.Ext[saft.ExtKeyExemption].String())
		assert.Equal(t, tax.KeyIntraCommunity, combo.Key)
	})

	t.Run("reverse map exemption M99", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				saft.ExtKeyExemption: "M99",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M99", combo.Ext[saft.ExtKeyExemption].String())
		assert.Equal(t, tax.KeyOutsideScope, combo.Key)
	})

	t.Run("reverse map exemption M01", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				saft.ExtKeyExemption: "M01",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M01", combo.Ext[saft.ExtKeyExemption].String())
		assert.Equal(t, tax.KeyExempt, combo.Key)
	})

	t.Run("reverse map exemption M44", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				saft.ExtKeyExemption: "M44",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M44", combo.Ext[saft.ExtKeyExemption].String())
		assert.Equal(t, tax.KeyOutsideScope, combo.Key)
	})

	t.Run("reverse map exemption M45", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				saft.ExtKeyExemption: "M45",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M45", combo.Ext[saft.ExtKeyExemption].String())
		assert.Equal(t, tax.KeyExempt, combo.Key)
	})

	t.Run("reverse map exemption M46", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				saft.ExtKeyExemption: "M46",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M46", combo.Ext[saft.ExtKeyExemption].String())
		assert.Equal(t, tax.KeyExport, combo.Key)
	})

	t.Run("rate missing but extension present", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				saft.ExtKeyTaxRate: "INT",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "INT", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, tax.RateIntermediate, combo.Rate)
	})

	t.Run("rate and extension mismatching", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
			Ext: tax.Extensions{
				saft.ExtKeyTaxRate: "INT",
			},
		}
		ad.Normalizer(combo)
		assert.Equal(t, "NOR", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, tax.RateGeneral, combo.Rate)
	})
}

func TestTaxComboValidate(t *testing.T) {
	ad := tax.AddonForKey(saft.V1)
	t.Run("valid", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(230, 3),
			Ext: tax.Extensions{
				pt.ExtKeyRegion:    "PT",
				saft.ExtKeyTaxRate: "NOR",
			},
		}
		err := ad.Validator(combo)
		assert.NoError(t, err)
	})

	t.Run("missing rate", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				pt.ExtKeyRegion: "PT",
			},
		}
		err := ad.Validator(combo)
		assert.ErrorContains(t, err, "ext: (pt-saft-tax-rate: required.)")
	})

	t.Run("missing rate foreign", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Country:  l10n.EL.Tax(),
			Ext: tax.Extensions{
				pt.ExtKeyRegion: "PT",
			},
		}
		err := ad.Validator(combo)
		assert.ErrorContains(t, err, "ext: (pt-saft-tax-rate: required.)")
	})

	t.Run("valid exempt", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				pt.ExtKeyRegion:      "PT",
				saft.ExtKeyTaxRate:   "ISE",
				saft.ExtKeyExemption: "M01",
			},
		}
		err := ad.Validator(combo)
		assert.NoError(t, err)
	})

	t.Run("missing exempt", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				pt.ExtKeyRegion:    "PT",
				saft.ExtKeyTaxRate: "ISE",
			},
		}
		err := ad.Validator(combo)
		assert.ErrorContains(t, err, "ext: (pt-saft-exemption: required.)")
	})

	t.Run("other category", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryGST,
		}
		err := ad.Validator(combo)
		assert.NoError(t, err)
	})

}
