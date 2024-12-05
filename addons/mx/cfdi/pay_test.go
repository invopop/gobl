package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPaymentMeansExtensions(t *testing.T) {
	ext := cfdi.PaymentMeansExtensions()
	assert.Equal(t, "01", ext.Get(pay.MeansKeyCash).String())
}

func TestNormalizePayInstructions(t *testing.T) {
	ad := tax.AddonForKey(cfdi.V4)

	t.Run("nil", func(t *testing.T) {
		var instr *pay.Instructions
		ad.Normalizer(instr)
	})

	t.Run("with match", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyOnline.With(cfdi.MeansKeyWallet),
		}
		ad.Normalizer(instr)
		assert.Equal(t, "05", instr.Ext.Get(cfdi.ExtKeyPaymentMeans).String())
	})
}

func TestNormalizePayAdvance(t *testing.T) {
	ad := tax.AddonForKey(cfdi.V4)

	t.Run("nil", func(t *testing.T) {
		var adv *pay.Advance
		ad.Normalizer(adv)
	})

	t.Run("with match", func(t *testing.T) {
		adv := &pay.Advance{
			Key: pay.MeansKeyOnline.With(cfdi.MeansKeyWallet),
		}
		ad.Normalizer(adv)
		assert.Equal(t, "05", adv.Ext.Get(cfdi.ExtKeyPaymentMeans).String())
	})
}

func TestValidatePayTerms(t *testing.T) {
	ad := tax.AddonForKey(cfdi.V4)

	t.Run("nil", func(t *testing.T) {
		var terms *pay.Terms
		ad.Validator(terms)
	})

	t.Run("valid", func(t *testing.T) {
		terms := &pay.Terms{
			Key:   pay.MeansKeyOnline.With(cfdi.MeansKeyWallet),
			Notes: "test",
		}
		err := ad.Validator(terms)
		assert.NoError(t, err)
	})
}
