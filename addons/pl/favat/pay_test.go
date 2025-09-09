package favat_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPaymentMeansExtensions(t *testing.T) {
	ext := favat.PaymentMeansExtensions()
	assert.Equal(t, "1", ext.Get(pay.MeansKeyCash).String())
}

func TestNormalizePayInstructions(t *testing.T) {
	ad := tax.AddonForKey(favat.V2)

	t.Run("nil", func(t *testing.T) {
		var instr *pay.Instructions
		assert.NotPanics(t, func() {
			ad.Normalizer(instr)
		})
	})

	t.Run("with match", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyOnline.With(favat.MeansKeyLoan),
		}
		ad.Normalizer(instr)
		assert.Equal(t, "5", instr.Ext.Get(favat.ExtKeyPaymentMeans).String())
	})
}

func TestNormalizePayAdvance(t *testing.T) {
	ad := tax.AddonForKey(favat.V2)

	t.Run("nil", func(t *testing.T) {
		var adv *pay.Advance
		assert.NotPanics(t, func() {
			ad.Normalizer(adv)
		})
	})

	t.Run("with match", func(t *testing.T) {
		adv := &pay.Advance{
			Key: pay.MeansKeyOnline.With(favat.MeansKeyLoan),
		}
		ad.Normalizer(adv)
		assert.Equal(t, "5", adv.Ext.Get(favat.ExtKeyPaymentMeans).String())
	})
}

func TestValidatePay(t *testing.T) {
	ad := tax.AddonForKey(favat.V2)

	t.Run("advance nil", func(t *testing.T) {
		var adv *pay.Advance
		assert.NotPanics(t, func() {
			assert.NoError(t, ad.Validator(adv))
		})
	})

	t.Run("advance valid", func(t *testing.T) {
		adv := &pay.Advance{
			Key: pay.MeansKeyOnline.With(favat.MeansKeyLoan),
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
			Key: pay.MeansKeyOnline.With(favat.MeansKeyLoan),
		}
		ad.Normalizer(instr)
		err := ad.Validator(instr)
		assert.NoError(t, err)
	})

}
