package arca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/ar/arca"
	"github.com/invopop/gobl/regimes/ar"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxCombo(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("nil tax combo does not panic", func(t *testing.T) {
		var tc *tax.Combo
		assert.NotPanics(t, func() {
			ad.Normalizer(tc)
		})
	})

	t.Run("zero rate key sets VAT rate 3", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyZero,
		}
		ad.Normalizer(tc)
		assert.Equal(t, "3", tc.Ext[arca.ExtKeyVATRate].String())
	})

	t.Run("reduced rate sets VAT rate 4", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateReduced,
		}
		ad.Normalizer(tc)
		assert.Equal(t, "4", tc.Ext[arca.ExtKeyVATRate].String())
	})

	t.Run("general rate sets VAT rate 5", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		}
		ad.Normalizer(tc)
		assert.Equal(t, "5", tc.Ext[arca.ExtKeyVATRate].String())
	})

	t.Run("increased rate sets VAT rate 6", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     ar.RateIncreased,
		}
		ad.Normalizer(tc)
		assert.Equal(t, "6", tc.Ext[arca.ExtKeyVATRate].String())
	})

	t.Run("other rate does not set VAT rate", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateSpecial,
		}
		ad.Normalizer(tc)
		assert.Empty(t, tc.Ext[arca.ExtKeyVATRate])
	})

	t.Run("existing extensions are preserved", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
			Ext: tax.Extensions{
				"custom-key": "custom-value",
			},
		}
		ad.Normalizer(tc)
		assert.Equal(t, "5", tc.Ext[arca.ExtKeyVATRate].String())
		assert.Equal(t, "custom-value", tc.Ext["custom-key"].String())
	})
}

func TestValidateTaxCombo(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("valid VAT combo with rate extension", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.Extensions{
				arca.ExtKeyVATRate: "5",
			},
		}
		err := ad.Validator(tc)
		assert.NoError(t, err)
	})

	t.Run("VAT combo missing rate extension", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
		}
		err := ad.Validator(tc)
		assert.ErrorContains(t, err, "ar-arca-vat-rate: required")
	})

	t.Run("non-VAT combo does not require rate extension", func(t *testing.T) {
		tc := &tax.Combo{
			Category: "OTHER",
		}
		err := ad.Validator(tc)
		assert.NoError(t, err)
	})
}
