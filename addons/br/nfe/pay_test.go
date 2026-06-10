package nfe_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfe"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizePayInstructions(t *testing.T) {

	t.Run("nil", func(t *testing.T) {
		var instr *pay.Instructions
		assert.NotPanics(t, func() {
			norm.Normalize(instr, tax.AddonContext(nfe.V4))
		})
	})

	t.Run("with match", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyCash,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				nfe.ExtKeyPaymentMeans: "15", // must be overridden
			}),
		}
		norm.Normalize(instr, tax.AddonContext(nfe.V4))
		assert.Equal(t, "01", instr.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("without match", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: "unknown-payment-means",
		}
		norm.Normalize(instr, tax.AddonContext(nfe.V4))
		assert.Empty(t, instr.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("with other key and extension", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyOther,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				nfe.ExtKeyPaymentMeans: "13", // must be kept
			}),
		}
		norm.Normalize(instr, tax.AddonContext(nfe.V4))
		assert.Equal(t, "13", instr.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("with other key and no extension", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyOther,
		}
		norm.Normalize(instr, tax.AddonContext(nfe.V4))
		assert.Equal(t, "99", instr.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("preserves existing extensions", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyCard,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				"other-extension": "value",
			}),
		}
		norm.Normalize(instr, tax.AddonContext(nfe.V4))
		assert.Equal(t, "03", instr.Ext.Get(nfe.ExtKeyPaymentMeans).String())
		assert.Equal(t, "value", instr.Ext.Get("other-extension").String())
	})
}

func TestNormalizePayAdvance(t *testing.T) {

	t.Run("nil", func(t *testing.T) {
		var adv *pay.Record
		assert.NotPanics(t, func() {
			norm.Normalize(adv, tax.AddonContext(nfe.V4))
		})
	})

	t.Run("with match", func(t *testing.T) {
		adv := &pay.Record{
			Key: pay.MeansKeyCard,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				nfe.ExtKeyPaymentMeans: "14", // must be overridden
			}),
		}
		norm.Normalize(adv, tax.AddonContext(nfe.V4))
		assert.Equal(t, "03", adv.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("without match", func(t *testing.T) {
		adv := &pay.Record{
			Key: "unknown-payment-means",
		}
		norm.Normalize(adv, tax.AddonContext(nfe.V4))
		assert.Empty(t, adv.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("with other key and extension", func(t *testing.T) {
		adv := &pay.Record{
			Key: pay.MeansKeyOther,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				nfe.ExtKeyPaymentMeans: "13", // must be kept
			}),
		}
		norm.Normalize(adv, tax.AddonContext(nfe.V4))
		assert.Equal(t, "13", adv.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})

	t.Run("with other key and no extension", func(t *testing.T) {
		adv := &pay.Record{
			Key: pay.MeansKeyOther,
		}
		norm.Normalize(adv, tax.AddonContext(nfe.V4))
		assert.Equal(t, "99", adv.Ext.Get(nfe.ExtKeyPaymentMeans).String())
	})
}

func TestValidatePayInstructions(t *testing.T) {
	t.Run("with payment means", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyCash,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				nfe.ExtKeyPaymentMeans: "01",
			}),
		}
		err := rules.Validate(instr, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("without payment means", func(t *testing.T) {
		instr := &pay.Instructions{
			Key: pay.MeansKeyCash,
		}
		err := rules.Validate(instr, withAddonContext())
		assert.ErrorContains(t, err, "payment instructions require 'br-nfe-payment-means' extension")
	})

	t.Run("nil", func(t *testing.T) {
		var instr *pay.Instructions
		err := rules.Validate(instr, withAddonContext())
		assert.NoError(t, err)
	})
}

func TestValidatePayAdvance(t *testing.T) {
	t.Run("with payment means", func(t *testing.T) {
		adv := &pay.Record{
			Key:         pay.MeansKeyCard,
			Description: "Card payment",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				nfe.ExtKeyPaymentMeans: "03",
			}),
		}
		err := rules.Validate(adv, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("without payment means", func(t *testing.T) {
		adv := &pay.Record{
			Key: pay.MeansKeyCard,
		}
		err := rules.Validate(adv, withAddonContext())
		assert.ErrorContains(t, err, "payment advance requires 'br-nfe-payment-means' extension")
	})

	t.Run("nil", func(t *testing.T) {
		var adv *pay.Record
		err := rules.Validate(adv, withAddonContext())
		assert.NoError(t, err)
	})
}
