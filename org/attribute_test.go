package org_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttributeValidation(t *testing.T) {
	t.Run("valid with text", func(t *testing.T) {
		a := &org.Attribute{
			Key:   org.AttributeKeyColor,
			Label: "Color",
			Text:  "Black",
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("valid with code", func(t *testing.T) {
		a := &org.Attribute{
			Key:  org.AttributeKeyColor,
			Code: "RAL 5015",
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("valid with amount and unit", func(t *testing.T) {
		a := &org.Attribute{
			Key:    org.AttributeKeyWeight,
			Amount: num.NewAmount(200, 0),
			Unit:   org.UnitGram,
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("valid with amount and no unit", func(t *testing.T) {
		a := &org.Attribute{
			Key:    "thread-count",
			Amount: num.NewAmount(300, 0),
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("valid with date", func(t *testing.T) {
		a := &org.Attribute{
			Key:  org.AttributeKeyExpiry,
			Date: cal.NewDate(2026, time.December, 31),
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("valid with type instead of key", func(t *testing.T) {
		a := &org.Attribute{
			Type: "X01",
			Text: "Black",
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("missing key and type", func(t *testing.T) {
		a := &org.Attribute{
			Text: "Black",
		}
		assert.ErrorContains(t, rules.Validate(a), "attribute must have either a key or a type, but not both")
	})
	t.Run("both key and type", func(t *testing.T) {
		a := &org.Attribute{
			Key:  org.AttributeKeyColor,
			Type: "X01",
			Text: "Black",
		}
		assert.ErrorContains(t, rules.Validate(a), "attribute must have either a key or a type, but not both")
	})
	t.Run("missing value", func(t *testing.T) {
		a := &org.Attribute{
			Key: org.AttributeKeyColor,
		}
		assert.ErrorContains(t, rules.Validate(a), "attribute must have exactly one of the text, code, amount, or date values")
	})
	t.Run("multiple values", func(t *testing.T) {
		a := &org.Attribute{
			Key:    org.AttributeKeyWeight,
			Text:   "200g",
			Amount: num.NewAmount(200, 0),
		}
		assert.ErrorContains(t, rules.Validate(a), "attribute must have exactly one of the text, code, amount, or date values")
	})
	t.Run("text and code", func(t *testing.T) {
		a := &org.Attribute{
			Key:  org.AttributeKeyColor,
			Text: "Black",
			Code: "RAL 9005",
		}
		assert.ErrorContains(t, rules.Validate(a), "attribute must have exactly one of the text, code, amount, or date values")
	})
	t.Run("unit without amount", func(t *testing.T) {
		a := &org.Attribute{
			Key:  org.AttributeKeyWeight,
			Text: "200",
			Unit: org.UnitGram,
		}
		assert.ErrorContains(t, rules.Validate(a), "attribute unit may only be used alongside an amount")
	})
}

func TestAttributeNormalization(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var a *org.Attribute
		assert.NotPanics(t, func() {
			norm.Normalize(a)
		})
	})
	t.Run("trims strings", func(t *testing.T) {
		a := &org.Attribute{
			Label: "  Color  ",
			Text:  " Black ",
		}
		norm.Normalize(a)
		assert.Equal(t, "Color", a.Label)
		assert.Equal(t, "Black", a.Text)
	})
}

func TestCleanAttributes(t *testing.T) {
	t.Run("removes nil and empty entries", func(t *testing.T) {
		attrs := []*org.Attribute{
			nil,
			{},
			{Key: "color", Text: "Black"},
		}
		out := org.CleanAttributes(attrs)
		require.Len(t, out, 1)
		assert.Equal(t, "Black", out[0].Text)
	})
	t.Run("keeps partially filled entries", func(t *testing.T) {
		attrs := []*org.Attribute{
			{Unit: org.UnitGram},
		}
		assert.Len(t, org.CleanAttributes(attrs), 1)
	})
	t.Run("returns nil when none remain", func(t *testing.T) {
		assert.Nil(t, org.CleanAttributes([]*org.Attribute{nil, {}}))
		assert.Nil(t, org.CleanAttributes(nil))
	})
}

func TestAttributesHaveUniqueKeys(t *testing.T) {
	test := org.AttributesHaveUniqueKeys()
	t.Run("unique keys", func(t *testing.T) {
		assert.True(t, test.Check([]*org.Attribute{
			{Key: org.AttributeKeyColor, Text: "Black"},
			{Key: org.AttributeKeySize, Text: "L"},
		}))
	})
	t.Run("duplicate keys", func(t *testing.T) {
		assert.False(t, test.Check([]*org.Attribute{
			{Key: org.AttributeKeyColor, Text: "Black"},
			{Key: org.AttributeKeyColor, Text: "White"},
		}))
	})
	t.Run("ignores nil entries and missing keys", func(t *testing.T) {
		assert.True(t, test.Check([]*org.Attribute{
			nil,
			{Type: "X01", Text: "Black"},
			{Type: "X01", Text: "White"},
			{Key: org.AttributeKeyColor, Text: "Black"},
		}))
	})
	t.Run("empty list", func(t *testing.T) {
		assert.True(t, test.Check([]*org.Attribute{}))
		assert.True(t, test.Check([]*org.Attribute(nil)))
	})
	t.Run("unexpected type", func(t *testing.T) {
		assert.False(t, test.Check("not a list"))
	})
}

func TestAttributeJSONSchemaExtend(t *testing.T) {
	base := `
		{
			"properties": {
				"key": {
					"$ref": "https://gobl.org/draft-0/cbc/key",
					"title": "Key"
				}
			}
		}
	`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(base), js))
	org.Attribute{}.JSONSchemaExtend(js)

	prop, ok := js.Properties.Get("key")
	require.True(t, ok)
	require.Len(t, prop.AnyOf, len(org.AttributeKeyDefinitions)+1)
	assert.Equal(t, org.AttributeKeyColor, prop.AnyOf[0].Const)
	assert.Equal(t, "Color", prop.AnyOf[0].Title)
	last := prop.AnyOf[len(prop.AnyOf)-1]
	assert.Equal(t, "Other", last.Title)
	assert.NotEmpty(t, last.Pattern)
}
