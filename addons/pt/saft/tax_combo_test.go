package saft_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
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
		assert.Equal(t, "ISE", combo.Ext[saft.ExtKeyTaxRate].String())
		assert.Equal(t, "M99", combo.Ext[saft.ExtKeyExemption].String())
	})
}

func TestTaxComboValidate(t *testing.T) {
	ad := tax.AddonForKey(saft.V1)
	t.Run("valid", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(230, 3),
			Ext: tax.Extensions{
				saft.ExtKeyTaxRate: "NOR",
			},
		}
		err := ad.Validator(combo)
		assert.NoError(t, err)
	})

	t.Run("missing rate", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
		}
		err := ad.Validator(combo)
		assert.ErrorContains(t, err, "ext: (pt-saft-tax-rate: required.)")
	})

	t.Run("missing rate foreign", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Country:  l10n.EL.Tax(),
		}
		err := ad.Validator(combo)
		assert.ErrorContains(t, err, "ext: (pt-saft-tax-rate: required.)")
	})

	t.Run("valid exempt", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
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
