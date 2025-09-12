package mydata_test

import (
	"testing"

	"github.com/invopop/gobl/addons/gr/mydata"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxCombo(t *testing.T) {
	ad := tax.AddonForKey(mydata.V1)

	t.Run("vat category presence", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Rate:     tax.RateGeneral,
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
			Key:      tax.KeyExempt,
		}
		ad.Normalizer(tc)
		err := ad.Validator(tc)
		assert.NoError(t, err)
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

func TestNormalizeTaxCombo(t *testing.T) {
	ad := tax.AddonForKey(mydata.V1)

	t.Run("nil", func(t *testing.T) {
		var tc *tax.Combo
		ad.Normalizer(tc)
		assert.Nil(t, tc)
	})

	t.Run("standard - general", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Rate:     tax.RateGeneral,
			Percent:  num.NewPercentage(4, 2),
		}
		ad.Normalizer(tc)
		assert.Equal(t, "4%", tc.Percent.String())
		assert.Equal(t, "standard", tc.Key.String())
		assert.Equal(t, "general", tc.Rate.String())
		assert.Equal(t, "1", tc.Ext.Get(mydata.ExtKeyVATRate).String())
	})

	t.Run("standard - reduced", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Rate:     tax.RateReduced,
			Percent:  num.NewPercentage(4, 2),
		}
		ad.Normalizer(tc)
		assert.Equal(t, "4%", tc.Percent.String())
		assert.Equal(t, "standard", tc.Key.String())
		assert.Equal(t, "reduced", tc.Rate.String())
		assert.Equal(t, "2", tc.Ext.Get(mydata.ExtKeyVATRate).String())
	})

	t.Run("standard - with exempt", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Rate:     tax.RateReduced,
			Percent:  num.NewPercentage(4, 2),
			Ext:      tax.Extensions{mydata.ExtKeyExemption: "3"},
		}
		ad.Normalizer(tc)
		assert.Equal(t, "4%", tc.Percent.String())
		assert.Equal(t, "standard", tc.Key.String())
		assert.Equal(t, "reduced", tc.Rate.String())
		assert.Equal(t, "2", tc.Ext.Get(mydata.ExtKeyVATRate).String())
		assert.Empty(t, tc.Ext.Get(mydata.ExtKeyExemption).String())
	})

	t.Run("exempt", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
		}
		ad.Normalizer(tc)
		assert.Equal(t, "7", tc.Ext.Get(mydata.ExtKeyExemption).String())
		assert.Equal(t, "exempt", tc.Key.String())
	})

	t.Run("all exemption codes", func(t *testing.T) {
		ex := cbc.GetKeyDefinition(mydata.ExtKeyExemption, ad.Extensions)
		codes := cbc.DefinitionCodes(ex.Values)
		for _, code := range codes {
			tc := &tax.Combo{
				Category: tax.CategoryVAT,
				Ext:      tax.Extensions{mydata.ExtKeyExemption: code},
			}
			ad.Normalizer(tc)
			assert.Equal(t, code, tc.Ext.Get(mydata.ExtKeyExemption))
			assert.Equal(t, cbc.Code("7"), tc.Ext.Get(mydata.ExtKeyVATRate))
			switch code {
			case "1", "2", "24", "29", "30", "31":
				assert.Equal(t, tax.KeyOutsideScope, tc.Key)
			case "3", "4", "5", "6", "7", "9", "10", "11", "12", "13", "15", "17",
				"18", "19", "20", "21", "22", "23", "25", "26", "27":
				assert.Equal(t, tax.KeyExempt, tc.Key)
			case "8", "28":
				assert.Equal(t, tax.KeyExport, tc.Key)
			case "14":
				assert.Equal(t, tax.KeyIntraCommunity, tc.Key)
			case "16":
				assert.Equal(t, tax.KeyReverseCharge, tc.Key)
			default:
				assert.Equal(t, tax.KeyStandard, tc.Key)
			}
		}
	})

	t.Run("sale with destination EU country tax rates", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  "ES",
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		}
		ad.Normalizer(tc)
		assert.Equal(t, "30", tc.Ext.Get(mydata.ExtKeyExemption).String())
		assert.Equal(t, "7", tc.Ext.Get(mydata.ExtKeyVATRate).String())
		assert.Equal(t, "outside-scope", tc.Key.String())
	})

	t.Run("sale with destination outside-EU country tax rates", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  "GB",
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		}
		ad.Normalizer(tc)
		assert.Equal(t, "29", tc.Ext.Get(mydata.ExtKeyExemption).String())
		assert.Equal(t, "7", tc.Ext.Get(mydata.ExtKeyVATRate).String())
		assert.Equal(t, "outside-scope", tc.Key.String())
	})

}
