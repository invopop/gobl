package org_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestItemNormalization(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var i *org.Item
		assert.NotPanics(t, func() {
			i.Normalize(nil)
		})
	})
	t.Run("extensions", func(t *testing.T) {
		i := &org.Item{
			Name:  "test item",
			Price: num.MakeAmount(100, 2),
			Ext: tax.Extensions{
				"test": "",
			},
		}
		i.Normalize(nil)
		assert.Equal(t, "test item", i.Name)
		assert.Equal(t, num.MakeAmount(100, 2), i.Price)
		assert.Nil(t, i.Ext)
	})
}

func TestItemValidation(t *testing.T) {
	// Check that the item is valid
	t.Run("basics", func(t *testing.T) {
		i := &org.Item{
			Name:  "test item",
			Price: num.MakeAmount(100, 2),
		}
		assert.NoError(t, i.Validate())
	})
}
