package saft_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxComboNormalize(t *testing.T) {
	// most tests here done in TestInvoice
	t.Run("standard", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "NOR", combo.Ext.Get(saft.ExtKeyTaxRate).String())
	})

	t.Run("standard with exempt reason", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Rate:     tax.RateGeneral,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyExemption: "M01",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "NOR", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Empty(t, combo.Ext.Get(saft.ExtKeyExemption))
	})

	t.Run("reverse-charge", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyReverseCharge,
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M40", combo.Ext.Get(saft.ExtKeyExemption).String())
	})

	t.Run("outside-scope", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M99", combo.Ext.Get(saft.ExtKeyExemption).String())
	})
	t.Run("intra-community", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyIntraCommunity,
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M16", combo.Ext.Get(saft.ExtKeyExemption).String())
	})

	t.Run("unsupported", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateSuperReduced,
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.True(t, combo.Ext.IsZero())
	})

	t.Run("nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var c *tax.Combo
			norm.Normalize(c, tax.AddonContext(saft.V1))
		})
	})

	t.Run("no rate", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.True(t, combo.Ext.IsZero())
	})

	t.Run("no rate, no override ext", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(21, 3),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyTaxRate: "NOR",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "NOR", combo.Ext.Get(saft.ExtKeyTaxRate).String())
	})

	t.Run("foreign rate", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Country:  l10n.EL.Tax(),
			Percent:  num.NewPercentage(24, 3),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyTaxRate: "NOR",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "OUT", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.False(t, combo.Ext.Has(saft.ExtKeyExemption))
	})

	t.Run("reverse map exemption M30", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyExemption: "M30",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M30", combo.Ext.Get(saft.ExtKeyExemption).String())
		assert.Equal(t, tax.KeyReverseCharge, combo.Key)
	})

	t.Run("reverse map exemption M05", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyExemption: "M05",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M05", combo.Ext.Get(saft.ExtKeyExemption).String())
		assert.Equal(t, tax.KeyExport, combo.Key)
	})

	t.Run("reverse map exemption M16", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyExemption: "M16",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M16", combo.Ext.Get(saft.ExtKeyExemption).String())
		assert.Equal(t, tax.KeyIntraCommunity, combo.Key)
	})

	t.Run("reverse map exemption M99", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyExemption: "M99",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M99", combo.Ext.Get(saft.ExtKeyExemption).String())
		assert.Equal(t, tax.KeyOutsideScope, combo.Key)
	})

	t.Run("reverse map exemption M01", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyExemption: "M01",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M01", combo.Ext.Get(saft.ExtKeyExemption).String())
		assert.Equal(t, tax.KeyExempt, combo.Key)
	})

	t.Run("reverse map exemption M44", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyExemption: "M44",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M44", combo.Ext.Get(saft.ExtKeyExemption).String())
		assert.Equal(t, tax.KeyOutsideScope, combo.Key)
	})

	t.Run("reverse map exemption M45", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyExemption: "M45",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M45", combo.Ext.Get(saft.ExtKeyExemption).String())
		assert.Equal(t, tax.KeyOutsideScope, combo.Key)
	})

	t.Run("reverse map exemption M46", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyExemption: "M46",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "ISE", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, "M46", combo.Ext.Get(saft.ExtKeyExemption).String())
		assert.Equal(t, tax.KeyOutsideScope, combo.Key)
	})

	t.Run("rate missing but extension present", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyTaxRate: "INT",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "INT", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, tax.RateIntermediate, combo.Rate)
	})

	t.Run("rate and extension mismatching", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				saft.ExtKeyTaxRate: "INT",
			}),
		}
		norm.Normalize(combo, tax.AddonContext(saft.V1))
		assert.Equal(t, "NOR", combo.Ext.Get(saft.ExtKeyTaxRate).String())
		assert.Equal(t, tax.RateGeneral, combo.Rate)
	})
}

func TestTaxComboValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(230, 3),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				pt.ExtKeyRegion:    "PT",
				saft.ExtKeyTaxRate: "NOR",
			}),
		}
		err := rules.Validate(combo, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("missing rate", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				pt.ExtKeyRegion: "PT",
			}),
		}
		err := rules.Validate(combo, withAddonContext())
		assert.ErrorContains(t, err, "region and tax rate are required")
	})

	t.Run("missing rate foreign", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Country:  l10n.EL.Tax(),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				pt.ExtKeyRegion: "PT",
			}),
		}
		err := rules.Validate(combo, withAddonContext())
		assert.ErrorContains(t, err, "region and tax rate are required")
	})

	t.Run("valid exempt", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				pt.ExtKeyRegion:      "PT",
				saft.ExtKeyTaxRate:   "ISE",
				saft.ExtKeyExemption: "M01",
			}),
		}
		err := rules.Validate(combo, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("missing exempt", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				pt.ExtKeyRegion:    "PT",
				saft.ExtKeyTaxRate: "ISE",
			}),
		}
		err := rules.Validate(combo, withAddonContext())
		assert.ErrorContains(t, err, "exemption is required when tax rate is exempt")
	})

	t.Run("other category", func(t *testing.T) {
		combo := &tax.Combo{
			Category: tax.CategoryGST,
		}
		err := rules.Validate(combo, withAddonContext())
		assert.NoError(t, err)
	})

}
