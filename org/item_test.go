package org_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				"test": "",
			}),
		}
		i.Normalize(nil)
		assert.Equal(t, "test item", i.Name)
		assert.Equal(t, num.NewAmount(100, 2), i.Price)
		assert.True(t, i.Ext.IsZero())
	})
	t.Run("clean ref", func(t *testing.T) {
		i := &org.Item{
			Name: "test item",
			Ref:  "  test-ref  ",
		}
		i.Normalize(nil)
		assert.Equal(t, "test-ref", i.Ref.String())
	})
}

func TestItemValidation(t *testing.T) {
	// Check that the item is valid
	t.Run("basics", func(t *testing.T) {
		i := &org.Item{
			Name: "test item",
		}
		assert.NoError(t, rules.Validate(i))
	})
	t.Run("missing name", func(t *testing.T) {
		i := new(org.Item)
		assert.ErrorContains(t, rules.Validate(i), "item name is required")
	})
	t.Run("without key", func(t *testing.T) {
		i := &org.Item{
			Name: "test item",
		}
		assert.NoError(t, rules.Validate(i))
	})
	t.Run("with key", func(t *testing.T) {
		i := &org.Item{
			Name: "test item",
			Key:  org.ItemKeyServices,
		}
		assert.NoError(t, rules.Validate(i))
	})
	t.Run("invalid key", func(t *testing.T) {
		i := &org.Item{
			Name: "test item",
			Key:  "invalid_key",
		}
		assert.ErrorContains(t, rules.Validate(i), "key must match the required pattern")
	})
}

func TestItemJSONSchema(t *testing.T) {
	base := here.Doc(`
		{
			"properties": {
				"key": {
					"$ref": "https://gobl.org/draft-0/cbc/key",
					"title": "Key",
					"description": "Special key used to classify the item sometimes required by some regimes."
				}
			}
		}
	`)
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(base), js))
	org.Item{}.JSONSchemaExtend(js)

	prop, ok := js.Properties.Get("key")
	assert.True(t, ok)
	assert.Len(t, prop.AnyOf, 3)
	assert.Equal(t, org.ItemKeyGoods, prop.AnyOf[0].Const)
	assert.Equal(t, "Goods", prop.AnyOf[0].Title)
	assert.Equal(t, org.ItemKeyServices, prop.AnyOf[1].Const)
	assert.Equal(t, "Services", prop.AnyOf[1].Title)
}
