package mydata_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/gr/mydata"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPaymentMeans(t *testing.T) {
	m := mydata.PaymentMeansExtensions()
	assert.False(t, m.IsZero())
	assert.Equal(t, 7, m.Len())
}

func TestPayInstructions(t *testing.T) {
	t.Run("valid cash", func(t *testing.T) {
		i := &pay.Instructions{
			Key: pay.MeansKeyCash,
		}
		norm.Normalize(i, tax.AddonContext(mydata.V1))
		assert.False(t, i.Ext.IsZero())
		err := rules.Validate(i, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("invalid key", func(t *testing.T) {
		i := &pay.Instructions{
			Key: cbc.Key("foo"),
		}
		norm.Normalize(i, tax.AddonContext(mydata.V1))
		assert.True(t, i.Ext.IsZero())
		err := rules.Validate(i, withAddonContext())
		assert.ErrorContains(t, err, "payment instructions require 'gr-mydata-payment-means' extension")
	})

	t.Run("nil", func(t *testing.T) {
		var i *pay.Instructions
		norm.Normalize(i, tax.AddonContext(mydata.V1))
		assert.NoError(t, rules.Validate(i, withAddonContext()))
	})
}

func TestPayAdvance(t *testing.T) {
	t.Run("valid cash", func(t *testing.T) {
		i := &pay.Record{
			Key:         pay.MeansKeyCash,
			Description: "Cash advance",
		}
		norm.Normalize(i, tax.AddonContext(mydata.V1))
		assert.False(t, i.Ext.IsZero())
		err := rules.Validate(i, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("invalid key", func(t *testing.T) {
		i := &pay.Record{
			Key:         cbc.Key("foo"),
			Description: "Bad advance",
		}
		norm.Normalize(i, tax.AddonContext(mydata.V1))
		assert.True(t, i.Ext.IsZero())
		err := rules.Validate(i, withAddonContext())
		assert.ErrorContains(t, err, "payment advance requires 'gr-mydata-payment-means' extension")
	})

	t.Run("nil", func(t *testing.T) {
		var i *pay.Record
		norm.Normalize(i, tax.AddonContext(mydata.V1))
		assert.NoError(t, rules.Validate(i, withAddonContext()))
	})
}
