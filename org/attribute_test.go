package org_test

import (
	"testing"

	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestAttributeValidation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		a := &org.Attribute{
			Name:  "Colour",
			Value: "Black",
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("missing name", func(t *testing.T) {
		a := &org.Attribute{
			Value: "Black",
		}
		assert.ErrorContains(t, rules.Validate(a), "attribute name is required")
	})
	t.Run("missing value", func(t *testing.T) {
		a := &org.Attribute{
			Name: "Colour",
		}
		assert.ErrorContains(t, rules.Validate(a), "attribute value is required")
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
			Name:  "  Colour  ",
			Value: " Black ",
		}
		norm.Normalize(a)
		assert.Equal(t, "Colour", a.Name)
		assert.Equal(t, "Black", a.Value)
	})
}

func TestCleanAttributes(t *testing.T) {
	t.Run("removes nil and empty entries", func(t *testing.T) {
		attrs := []*org.Attribute{
			nil,
			{},
			{Name: "Colour", Value: "Black"},
		}
		out := org.CleanAttributes(attrs)
		assert.Len(t, out, 1)
		assert.Equal(t, "Colour", out[0].Name)
	})
	t.Run("returns nil when none remain", func(t *testing.T) {
		assert.Nil(t, org.CleanAttributes([]*org.Attribute{nil, {}}))
		assert.Nil(t, org.CleanAttributes(nil))
	})
}
