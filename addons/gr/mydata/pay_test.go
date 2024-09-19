package mydata_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/gr/mydata"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPayInstructions(t *testing.T) {
	ad := tax.AddonForKey(mydata.V1)
	t.Run("valid cash", func(t *testing.T) {
		i := &pay.Instructions{
			Key: pay.MeansKeyCash,
		}
		ad.Normalizer(i)
		assert.NotEmpty(t, i.Ext)
		err := ad.Validator(i)
		assert.NoError(t, err)
	})

	t.Run("invalid key", func(t *testing.T) {
		i := &pay.Instructions{
			Key: cbc.Key("foo"),
		}
		ad.Normalizer(i)
		assert.Empty(t, i.Ext)
		err := ad.Validator(i)
		assert.ErrorContains(t, err, "ext: (gr-mydata-payment-means: required.)")
	})

	t.Run("nil", func(t *testing.T) {
		var i *pay.Instructions
		ad.Normalizer(i)
		assert.NoError(t, ad.Validator(i))
	})
}

func TestPayAdvance(t *testing.T) {
	ad := tax.AddonForKey(mydata.V1)
	t.Run("valid cash", func(t *testing.T) {
		i := &pay.Advance{
			Key: pay.MeansKeyCash,
		}
		ad.Normalizer(i)
		assert.NotEmpty(t, i.Ext)
		err := ad.Validator(i)
		assert.NoError(t, err)
	})

	t.Run("invalid key", func(t *testing.T) {
		i := &pay.Advance{
			Key: cbc.Key("foo"),
		}
		ad.Normalizer(i)
		assert.Empty(t, i.Ext)
		err := ad.Validator(i)
		assert.ErrorContains(t, err, "ext: (gr-mydata-payment-means: required.)")
	})

	t.Run("nil", func(t *testing.T) {
		var i *pay.Advance
		ad.Normalizer(i)
		assert.NoError(t, ad.Validator(i))
	})
}
