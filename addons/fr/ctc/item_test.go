package ctc_test

import (
	"testing"

	"github.com/invopop/gobl/addons/fr/ctc"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestItemMetaValidation(t *testing.T) {
	ad := tax.AddonForKey(ctc.Flow2)

	t.Run("valid item with meta values", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
			Meta: cbc.Meta{
				"order-id":   "12345",
				"batch-code": "ABC-123",
			},
		}
		err := ad.Validator(item)
		assert.NoError(t, err)
	})

	t.Run("item with blank meta value", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
			Meta: cbc.Meta{
				"order-id":   "12345",
				"batch-code": "",
			},
		}
		err := ad.Validator(item)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "batch-code")
		assert.ErrorContains(t, err, "cannot be blank")
	})

	t.Run("item with whitespace-only meta value", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
			Meta: cbc.Meta{
				"order-id":   "12345",
				"batch-code": "   ",
			},
		}
		err := ad.Validator(item)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "batch-code")
		assert.ErrorContains(t, err, "cannot be blank")
	})

	t.Run("item without meta", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
		}
		err := ad.Validator(item)
		assert.NoError(t, err)
	})

	t.Run("item with empty meta map", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
			Meta: cbc.Meta{},
		}
		err := ad.Validator(item)
		assert.NoError(t, err)
	})

	t.Run("multiple blank values", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
			Meta: cbc.Meta{
				"order-id":   "",
				"batch-code": "ABC-123",
			},
		}
		err := ad.Validator(item)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "order-id")
	})

	t.Run("nil item", func(t *testing.T) {
		err := ad.Validator((*org.Item)(nil))
		assert.NoError(t, err)
	})
}
