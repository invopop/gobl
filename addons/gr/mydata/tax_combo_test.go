package mydata_test

import (
	"testing"

	"github.com/invopop/gobl/addons/gr/mydata"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxCombo(t *testing.T) {
	ad := tax.AddonForKey(mydata.V1)

	t.Run("vat category presence", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateStandard,
			Percent:  num.NewPercentage(4, 2),
		}
		err := ad.Validator(tc)
		assert.ErrorContains(t, err, "ext: (gr-mydata-vat-rate: required.)")
		ad.Normalizer(tc)
		assert.NoError(t, ad.Validator(tc))
	})

	t.Run("exemption presence", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt,
		}
		ad.Normalizer(tc)
		err := ad.Validator(tc)
		assert.ErrorContains(t, err, "gr-mydata-exemption: required")
	})

	t.Run("nil", func(t *testing.T) {
		var tc *tax.Combo
		assert.NoError(t, ad.Validator(tc))
	})

	t.Run("non-vat category", func(t *testing.T) {
		tc := &tax.Combo{
			Category: "FOO",
			Percent:  num.NewPercentage(4, 2),
		}
		assert.NoError(t, ad.Validator(tc))
	})
}
