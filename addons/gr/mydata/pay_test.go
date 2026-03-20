package mydata_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/gr/mydata"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPaymentMeans(t *testing.T) {
	m := mydata.PaymentMeansExtensions()
	assert.NotEmpty(t, m)
	assert.Len(t, m, 7)
}

func TestPayInstructions(t *testing.T) {
	ad := tax.AddonForKey(mydata.V1)
	t.Run("valid cash", func(t *testing.T) {
		i := &pay.Instructions{
			Key: pay.MeansKeyCash,
		}
		ad.Normalizer(i)
		assert.NotEmpty(t, i.Ext)
		err := rules.Validate(i, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("invalid key", func(t *testing.T) {
		i := &pay.Instructions{
			Key: cbc.Key("foo"),
		}
		ad.Normalizer(i)
		assert.Empty(t, i.Ext)
		err := rules.Validate(i, withAddonContext())
		assert.ErrorContains(t, err, "payment instructions require 'gr-mydata-payment-means' extension")
	})

	t.Run("nil", func(t *testing.T) {
		var i *pay.Instructions
		ad.Normalizer(i)
		assert.NoError(t, rules.Validate(i, withAddonContext()))
	})
}

func TestPayAdvance(t *testing.T) {
	ad := tax.AddonForKey(mydata.V1)
	t.Run("valid cash", func(t *testing.T) {
		i := &pay.Advance{
			Key:         pay.MeansKeyCash,
			Description: "Cash advance",
		}
		ad.Normalizer(i)
		assert.NotEmpty(t, i.Ext)
		err := rules.Validate(i, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("invalid key", func(t *testing.T) {
		i := &pay.Advance{
			Key:         cbc.Key("foo"),
			Description: "Bad advance",
		}
		ad.Normalizer(i)
		assert.Empty(t, i.Ext)
		err := rules.Validate(i, withAddonContext())
		assert.ErrorContains(t, err, "payment advance requires 'gr-mydata-payment-means' extension")
	})

	t.Run("nil", func(t *testing.T) {
		var i *pay.Advance
		ad.Normalizer(i)
		assert.NoError(t, rules.Validate(i, withAddonContext()))
	})
}
