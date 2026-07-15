package org_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageValidation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		p := &org.Package{
			Key:   org.PackageKeyBox,
			Count: 2,
			Attributes: []*org.Attribute{
				{
					Key:    org.AttributeKeyWeight,
					Amount: num.NewAmount(35, 1),
					Unit:   org.UnitKilogram,
				},
			},
		}
		assert.NoError(t, rules.Validate(p))
	})
	t.Run("valid with extended key", func(t *testing.T) {
		p := &org.Package{
			Key: org.PackageKeyBox.With("gift"),
		}
		assert.NoError(t, rules.Validate(p))
	})
	t.Run("valid without key", func(t *testing.T) {
		p := &org.Package{
			Label: "Custom wrap",
		}
		assert.NoError(t, rules.Validate(p))
	})
	t.Run("unknown key", func(t *testing.T) {
		p := &org.Package{
			Key: "spaceship",
		}
		assert.ErrorContains(t, rules.Validate(p), "package key must be or extend one of the pre-defined keys")
	})
	t.Run("negative count", func(t *testing.T) {
		p := &org.Package{
			Key:   org.PackageKeyBox,
			Count: -1,
		}
		assert.ErrorContains(t, rules.Validate(p), "package count must be zero or positive")
	})
	t.Run("duplicate attribute keys", func(t *testing.T) {
		p := &org.Package{
			Key: org.PackageKeyBox,
			Attributes: []*org.Attribute{
				{Key: org.AttributeKeyWeight, Amount: num.NewAmount(1, 0), Unit: org.UnitKilogram},
				{Key: org.AttributeKeyWeight, Amount: num.NewAmount(2, 0), Unit: org.UnitKilogram},
			},
		}
		assert.ErrorContains(t, rules.Validate(p), "package attributes must not contain duplicate keys")
	})
}

func TestPackageNormalization(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var p *org.Package
		assert.NotPanics(t, func() {
			norm.Normalize(p)
		})
	})
	t.Run("trims label and cleans attributes", func(t *testing.T) {
		p := &org.Package{
			Label: " Box 1 of 2 ",
			Attributes: []*org.Attribute{
				nil,
				{},
			},
		}
		norm.Normalize(p)
		assert.Equal(t, "Box 1 of 2", p.Label)
		assert.Empty(t, p.Attributes)
	})
}

func TestPackageUNECERec21(t *testing.T) {
	tests := []struct {
		name string
		pkg  *org.Package
		want cbc.Code
	}{
		{"nil", nil, cbc.CodeEmpty},
		{"empty key", &org.Package{}, cbc.CodeEmpty},
		{"box", &org.Package{Key: org.PackageKeyBox}, "BX"},
		{"pallet", &org.Package{Key: org.PackageKeyPallet}, "PX"},
		{"unpacked", &org.Package{Key: org.PackageKeyUnpacked}, "NE"},
		{"extended key", &org.Package{Key: org.PackageKeyBox.With("gift")}, "BX"},
		{"unknown key", &org.Package{Key: "spaceship"}, org.UNECERec21MutuallyDefined},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.pkg.UNECERec21())
		})
	}
}

func TestCleanPackages(t *testing.T) {
	t.Run("removes nil and empty entries", func(t *testing.T) {
		pkgs := []*org.Package{
			nil,
			{},
			{Key: org.PackageKeyBox},
		}
		out := org.CleanPackages(pkgs)
		require.Len(t, out, 1)
		assert.Equal(t, org.PackageKeyBox, out[0].Key)
	})
	t.Run("returns nil when none remain", func(t *testing.T) {
		assert.Nil(t, org.CleanPackages([]*org.Package{nil, {}}))
		assert.Nil(t, org.CleanPackages(nil))
	})
}

func TestPackageJSONSchemaExtend(t *testing.T) {
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
	org.Package{}.JSONSchemaExtend(js)

	prop, ok := js.Properties.Get("key")
	require.True(t, ok)
	require.Len(t, prop.AnyOf, len(org.PackageKeyDefinitions)+1)
	assert.Equal(t, org.PackageKeyBag, prop.AnyOf[0].Const)
	assert.Equal(t, "Bag", prop.AnyOf[0].Title)
	last := prop.AnyOf[len(prop.AnyOf)-1]
	assert.Equal(t, "Other", last.Title)
	assert.NotEmpty(t, last.Pattern)
}
