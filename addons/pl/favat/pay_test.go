package favat_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizePayInstructions(t *testing.T) {
	ad := tax.AddonForKey(favat.V3)

	t.Run("nil", func(t *testing.T) {
		var instr *pay.Instructions
		assert.NotPanics(t, func() {
			ad.Normalizer(instr)
		})
	})

	t.Run("with match", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyOther.With(favat.MeansKeyCredit),
		}
		ad.Normalizer(instr)
		assert.Equal(t, "5", instr.Ext.Get(favat.ExtKeyPaymentMeans).String())
	})
}

func TestNormalizePayAdvance(t *testing.T) {
	ad := tax.AddonForKey(favat.V3)

	t.Run("nil", func(t *testing.T) {
		var adv *pay.Advance
		assert.NotPanics(t, func() {
			ad.Normalizer(adv)
		})
	})

	t.Run("with match", func(t *testing.T) {
		adv := &pay.Advance{
			Key: pay.MeansKeyOther.With(favat.MeansKeyCredit),
		}
		ad.Normalizer(adv)
		assert.Equal(t, "5", adv.Ext.Get(favat.ExtKeyPaymentMeans).String())
	})
}

func TestValidatePay(t *testing.T) {
	ad := tax.AddonForKey(favat.V3)

	t.Run("advance nil", func(t *testing.T) {
		var adv *pay.Advance
		assert.NotPanics(t, func() {
			assert.NoError(t, ad.Validator(adv))
		})
	})

	t.Run("advance valid", func(t *testing.T) {
		adv := &pay.Advance{
			Key: pay.MeansKeyOther.With(favat.MeansKeyCredit),
		}
		ad.Normalizer(adv)
		err := ad.Validator(adv)
		assert.NoError(t, err)
	})

	t.Run("instructions nil", func(t *testing.T) {
		var instr *pay.Instructions
		assert.NotPanics(t, func() {
			assert.NoError(t, ad.Validator(instr))
		})
	})

	t.Run("instructions valid", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyOther.With(favat.MeansKeyCredit),
		}
		ad.Normalizer(instr)
		err := ad.Validator(instr)
		assert.NoError(t, err)
	})

}

func TestPaymentMeansMapping(t *testing.T) {
	ad := tax.AddonForKey(favat.V3)

	tests := []struct {
		name     string
		key      cbc.Key
		expected string
	}{
		{
			name:     "cash",
			key:      pay.MeansKeyCash,
			expected: "1",
		},
		{
			name:     "card",
			key:      pay.MeansKeyCard,
			expected: "2",
		},
		{
			name:     "voucher",
			key:      pay.MeansKeyOther.With(favat.MeansKeyVoucher),
			expected: "3",
		},
		{
			name:     "cheque",
			key:      pay.MeansKeyCheque,
			expected: "4",
		},
		{
			name:     "credit",
			key:      pay.MeansKeyOther.With(favat.MeansKeyCredit),
			expected: "5",
		},
		{
			name:     "credit transfer",
			key:      pay.MeansKeyCreditTransfer,
			expected: "6",
		},
		{
			name:     "online",
			key:      pay.MeansKeyOnline,
			expected: "7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" instructions", func(t *testing.T) {
			instr := &pay.Instructions{
				Key: tt.key,
			}
			ad.Normalizer(instr)
			assert.Equal(t, tt.expected, instr.Ext.Get(favat.ExtKeyPaymentMeans).String())
		})

		t.Run(tt.name+" advance", func(t *testing.T) {
			adv := &pay.Advance{
				Key: tt.key,
			}
			ad.Normalizer(adv)
			assert.Equal(t, tt.expected, adv.Ext.Get(favat.ExtKeyPaymentMeans).String())
		})
	}
}

func TestPaymentMeansNoMatch(t *testing.T) {
	ad := tax.AddonForKey(favat.V3)

	t.Run("instructions with unknown key", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: "unknown-payment-means",
		}
		ad.Normalizer(instr)
		assert.Equal(t, "", instr.Ext.Get(favat.ExtKeyPaymentMeans).String())
	})

	t.Run("advance with unknown key", func(t *testing.T) {
		adv := &pay.Advance{
			Key: "unknown-payment-means",
		}
		ad.Normalizer(adv)
		assert.Equal(t, "", adv.Ext.Get(favat.ExtKeyPaymentMeans).String())
	})
}
