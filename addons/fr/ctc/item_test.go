package ctc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestItemMetaValidation(t *testing.T) {
	t.Run("valid item with meta values", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
			Meta: cbc.Meta{
				"order-id":   "12345",
				"batch-code": "ABC-123",
			},
		}
		err := rules.Validate(item, withAddonContext())
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
		err := rules.Validate(item, withAddonContext())
		assert.Error(t, err)
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
		err := rules.Validate(item, withAddonContext())
		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot be blank")
	})

	t.Run("item without meta", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
		}
		err := rules.Validate(item, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("item with empty meta map", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
			Meta: cbc.Meta{},
		}
		err := rules.Validate(item, withAddonContext())
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
		err := rules.Validate(item, withAddonContext())
		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot be blank")
	})

	t.Run("nil item", func(t *testing.T) {
		err := rules.Validate((*org.Item)(nil), withAddonContext())
		assert.NoError(t, err)
	})
}
