package favat_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxCombo(t *testing.T) {
	ad := tax.AddonForKey(favat.V3)

	tests := []struct {
		name     string
		key      cbc.Key
		rate     cbc.Key
		expected string
	}{
		{
			name:     "standard rate general",
			key:      tax.KeyStandard,
			rate:     tax.RateGeneral,
			expected: "1",
		},
		{
			name:     "standard rate reduced",
			key:      tax.KeyStandard,
			rate:     tax.RateReduced,
			expected: "2",
		},
		{
			name:     "standard rate super reduced",
			key:      tax.KeyStandard,
			rate:     tax.RateSuperReduced,
			expected: "3",
		},
		{
			name:     "zero rate",
			key:      tax.KeyZero,
			rate:     "",
			expected: "6.1",
		},
		{
			name:     "intra community",
			key:      tax.KeyIntraCommunity,
			rate:     "",
			expected: "6.2",
		},
		{
			name:     "export",
			key:      tax.KeyExport,
			rate:     "",
			expected: "6.3",
		},
		{
			name:     "exempt",
			key:      tax.KeyExempt,
			rate:     "",
			expected: "7",
		},
		{
			name:     "outside scope",
			key:      tax.KeyOutsideScope,
			rate:     "",
			expected: "8",
		},
		{
			name:     "reverse charge",
			key:      tax.KeyReverseCharge,
			rate:     "",
			expected: "9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &tax.Combo{
				Key:     tt.key,
				Rate:    tt.rate,
				Percent: num.NewPercentage(23, 0),
			}
			ad.Normalizer(tc)
			assert.Equal(t, tt.expected, tc.Ext.Get(favat.ExtKeyTaxCategory).String())
		})
	}
}

func TestValidateTaxCombo(t *testing.T) {
	ad := tax.AddonForKey(favat.V3)

	t.Run("valid tax combo with category", func(t *testing.T) {
		tc := &tax.Combo{
			Key:     tax.KeyStandard,
			Rate:    tax.RateGeneral,
			Percent: num.NewPercentage(23, 0),
		}
		ad.Normalizer(tc)
		err := ad.Validator(tc)
		assert.NoError(t, err)
	})

	t.Run("invalid tax combo without category", func(t *testing.T) {
		tc := &tax.Combo{
			Key:     "unknown-key",
			Rate:    tax.RateGeneral,
			Percent: num.NewPercentage(23, 0),
		}
		err := ad.Validator(tc)
		assert.ErrorContains(t, err, "ext: (pl-favat-tax-category: required.)")
	})
}
