package sdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxCombo(t *testing.T) {
	ad := tax.AddonForKey(sdi.V1)
	t.Run("nil", func(t *testing.T) {
		var tc *tax.Combo
		assert.NotPanics(t, func() {
			ad.Normalizer(tc)
		})
	})

	t.Run("B2C EU export", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  "ES",
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Percent:  num.NewPercentage(21, 3),
		}
		ad.Normalizer(tc)
		assert.Equal(t, cbc.Code("N7"), tc.Ext.Get(sdi.ExtKeyExempt))
	})

	t.Run("B2C other export", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  "GB",
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Percent:  num.NewPercentage(20, 3),
		}
		ad.Normalizer(tc)
		assert.Equal(t, cbc.Code("N2.1"), tc.Ext.Get(sdi.ExtKeyExempt))
	})

	t.Run("check all extensions", func(t *testing.T) {
		ex := cbc.GetKeyDefinition(sdi.ExtKeyExempt, ad.Extensions)
		codes := cbc.DefinitionCodes(ex.Values)

		for _, code := range codes {
			t.Run("with code "+code.String(), func(t *testing.T) {
				tc := &tax.Combo{
					Category: tax.CategoryVAT,
					Ext:      tax.Extensions{sdi.ExtKeyExempt: code},
				}
				ad.Normalizer(tc)
				assert.Equal(t, code, tc.Ext.Get(sdi.ExtKeyExempt), "extension should be set correctly")
				assert.NotEmpty(t, tc.Key, "key should be set based on extension")
			})
		}
	})

	t.Run("check exempt handling", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Ext:      tax.Extensions{sdi.ExtKeyExempt: "N3.2"},
		}
		ad.Normalizer(tc)
		assert.Equal(t, "N3.2", tc.Ext.Get(sdi.ExtKeyExempt).String())
		assert.Equal(t, tax.KeyIntraCommunity, tc.Key)
	})

}
