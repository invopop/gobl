package org_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
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
			Price: num.NewAmount(100, 2),
			Ext: tax.Extensions{
				"test": "",
			},
		}
		i.Normalize(nil)
		assert.Equal(t, "test item", i.Name)
		assert.Equal(t, num.NewAmount(100, 2), i.Price)
		assert.Nil(t, i.Ext)
	})
}

func TestItemValidation(t *testing.T) {
	// Check that the item is valid
	t.Run("basics", func(t *testing.T) {
		i := &org.Item{
			Name: "test item",
		}
		assert.NoError(t, i.Validate())
	})
	t.Run("missing name", func(t *testing.T) {
		i := new(org.Item)
		assert.ErrorContains(t, i.Validate(), "name: cannot be blank.")
	})
}

func TestItemPriceRequired(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var item *org.Item
		assert.NoError(t, validation.Validate(item, org.ItemPriceRequired()))
	})
	t.Run("success", func(t *testing.T) {
		item := &org.Item{
			Name:  "test item",
			Price: num.NewAmount(100, 2),
		}
		assert.NoError(t, validation.Validate(item, org.ItemPriceRequired()))
	})
	t.Run("missing", func(t *testing.T) {
		obj := struct {
			Item *org.Item `json:"item"`
		}{
			Item: &org.Item{
				Name: "test item",
			},
		}
		err := validation.ValidateStruct(&obj,
			validation.Field(&obj.Item, org.ItemPriceRequired()),
		)
		assert.ErrorContains(t, err, "item: (price: cannot be blank.)")
	})

}
