package saft_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/pt/saft"
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
			Rate:     tax.RateStandard,
		}
		ad.Normalizer(combo)
		assert.Equal(t, "NOR", combo.Ext[saft.ExtKeyTaxRate].String())
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
}

func TestTaxRateKeyMap(t *testing.T) {
	m := saft.TaxRateExtensions()
	assert.Equal(t, tax.RateReduced, m.Lookup("RED"))
	assert.Equal(t, tax.RateIntermediate, m.Lookup("INT"))
	assert.Equal(t, tax.RateStandard, m.Lookup("NOR"))
	assert.Equal(t, tax.RateExempt, m.Lookup("ISE"))
	assert.Equal(t, tax.RateOther, m.Lookup("OUT"))
}
