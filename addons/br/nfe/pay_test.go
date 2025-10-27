package nfe_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfe"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizePayInstructions(t *testing.T) {
	ad := tax.AddonForKey(nfe.V4)

	t.Run("nil", func(t *testing.T) {
		var instr *pay.Instructions
		assert.NotPanics(t, func() {
			ad.Normalizer(instr)
		})
	})

	t.Run("with match", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyCash,
			Ext: tax.Extensions{
				nfe.ExtKeyPaymentMeans: "15", // must be overridden
			},
		}
		ad.Normalizer(instr)
		assert.Equal(t, "01", instr.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("without match", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: "unknown-payment-means",
		}
		ad.Normalizer(instr)
		assert.Empty(t, instr.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("with other key and extension", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyOther,
			Ext: tax.Extensions{
				nfe.ExtKeyPaymentMeans: "13", // must be kept
			},
		}
		ad.Normalizer(instr)
		assert.Equal(t, "13", instr.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("with other key and no extension", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyOther,
		}
		ad.Normalizer(instr)
		assert.Equal(t, "99", instr.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("preserves existing extensions", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyCard,
			Ext: tax.Extensions{
				"other-extension": "value",
			},
		}
		ad.Normalizer(instr)
		assert.Equal(t, "03", instr.Ext.Get(nfe.ExtKeyPaymentMeans).String())
		assert.Equal(t, "value", instr.Ext.Get("other-extension").String())
	})
}

func TestNormalizePayAdvance(t *testing.T) {
	ad := tax.AddonForKey(nfe.V4)

	t.Run("nil", func(t *testing.T) {
		var adv *pay.Advance
		assert.NotPanics(t, func() {
			ad.Normalizer(adv)
		})
	})

	t.Run("with match", func(t *testing.T) {
		adv := &pay.Advance{
			Key: pay.MeansKeyCard,
			Ext: tax.Extensions{
				nfe.ExtKeyPaymentMeans: "14", // must be overridden
			},
		}
		ad.Normalizer(adv)
		assert.Equal(t, "03", adv.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("without match", func(t *testing.T) {
		adv := &pay.Advance{
			Key: "unknown-payment-means",
		}
		ad.Normalizer(adv)
		assert.Empty(t, adv.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("with other key and extension", func(t *testing.T) {
		adv := &pay.Advance{
			Key: pay.MeansKeyOther,
			Ext: tax.Extensions{
				nfe.ExtKeyPaymentMeans: "13", // must be kept
			},
		}
		ad.Normalizer(adv)
		assert.Equal(t, "13", adv.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("with other key and no extension", func(t *testing.T) {
		adv := &pay.Advance{
			Key: pay.MeansKeyOther,
		}
		ad.Normalizer(adv)
		assert.Equal(t, "99", adv.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})
}

func TestValidatePayInstructions(t *testing.T) {
	ad := tax.AddonForKey(nfe.V4)

	t.Run("with payment means", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyCash,
			Ext: tax.Extensions{
				nfe.ExtKeyPaymentMeans: "01",
			},
		}
		err := ad.Validator(instr)
		assert.NoError(t, err)
	})

	t.Run("without payment means", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyCash,
		}
		err := ad.Validator(instr)
		assert.ErrorContains(t, err, "br-nfe-payment-means: required")
	})

	t.Run("nil", func(t *testing.T) {
		var instr *pay.Instructions
		err := ad.Validator(instr)
		assert.NoError(t, err)
	})
}

func TestValidatePayAdvance(t *testing.T) {
	ad := tax.AddonForKey(nfe.V4)

	t.Run("with payment means", func(t *testing.T) {
		adv := &pay.Advance{
			Key: pay.MeansKeyCard,
			Ext: tax.Extensions{
				nfe.ExtKeyPaymentMeans: "03",
			},
		}
		err := ad.Validator(adv)
		assert.NoError(t, err)
	})

	t.Run("without payment means", func(t *testing.T) {
		adv := &pay.Advance{
			Key: pay.MeansKeyCard,
		}
		err := ad.Validator(adv)
		assert.ErrorContains(t, err, "br-nfe-payment-means: required")
	})

	t.Run("nil", func(t *testing.T) {
		var adv *pay.Advance
		err := ad.Validator(adv)
		assert.NoError(t, err)
	})
}
