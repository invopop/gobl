package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbeddingAddons(t *testing.T) {
	type testStruct struct {
		tax.Addons
		Name string `json:"test"`
	}

	ts := &testStruct{
		Addons: tax.WithAddons("mx-cfdi-v4"),
		Name:   "Test",
	}
	assert.NotNil(t, ts.Addons)
	assert.Equal(t, "Test", ts.Name)

	assert.Equal(t, []cbc.Key{"mx-cfdi-v4"}, ts.GetAddons())

	defs := ts.AddonDefs()
	assert.Len(t, defs, 1)
	assert.Equal(t, "mx-cfdi-v4", defs[0].Key.String())

	ts.Addons = tax.WithAddons("mx-cfdi-v4", "invalid-addon")

	err := rules.Validate(ts)
	assert.ErrorContains(t, err, "[GOBL-TAX-ADDONS-01] ($.$addons[1]) add-on must be registered")

	t.Run("test addon normalization", func(t *testing.T) {
		ts.Addons.List = tax.AddonList{"mx-cfdi-v4", "mx-cfdi-v4", "de-xrechnung-v3"}
		_ = tax.ExtractNormalizers(ts)
		assert.Equal(t, tax.AddonList{"mx-cfdi-v4", "eu-en16931-v2017", "de-xrechnung-v3"}, ts.Addons.List)
	})
}

func TestAddonForKey(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		a := tax.AddonForKey("unknown")
		assert.Nil(t, a)
	})

	t.Run("found", func(t *testing.T) {
		a := tax.AddonForKey("mx-cfdi-v4")
		require.NotNil(t, a)
		assert.NoError(t, rules.Validate(a))
	})
}

func TestSetAddons(t *testing.T) {
	var as tax.Addons
	as.SetAddons("mx-cfdi-v4", "es-verifactu-v1")
	assert.Equal(t, []cbc.Key{"mx-cfdi-v4", "es-verifactu-v1"}, as.GetAddons())

	// SetAddons replaces the existing list wholesale.
	as.SetAddons("pt-saft-v1")
	assert.Equal(t, []cbc.Key{"pt-saft-v1"}, as.GetAddons())
}

func TestAddAddons(t *testing.T) {
	t.Run("appends to an empty list", func(t *testing.T) {
		var as tax.Addons
		as.AddAddons("mx-cfdi-v4")
		assert.Equal(t, []cbc.Key{"mx-cfdi-v4"}, as.GetAddons())
	})

	t.Run("appends to an existing list", func(t *testing.T) {
		as := tax.WithAddons("mx-cfdi-v4")
		as.AddAddons("es-verifactu-v1")
		assert.Equal(t, []cbc.Key{"mx-cfdi-v4", "es-verifactu-v1"}, as.GetAddons())
	})

	t.Run("skips empty keys", func(t *testing.T) {
		var as tax.Addons
		as.AddAddons("", "mx-cfdi-v4", "")
		assert.Equal(t, []cbc.Key{"mx-cfdi-v4"}, as.GetAddons())
	})

	t.Run("skips keys already present", func(t *testing.T) {
		as := tax.WithAddons("mx-cfdi-v4")
		as.AddAddons("mx-cfdi-v4")
		assert.Equal(t, []cbc.Key{"mx-cfdi-v4"}, as.GetAddons())
	})

	t.Run("de-duplicates within a single call", func(t *testing.T) {
		var as tax.Addons
		as.AddAddons("mx-cfdi-v4", "es-verifactu-v1", "mx-cfdi-v4")
		assert.Equal(t, []cbc.Key{"mx-cfdi-v4", "es-verifactu-v1"}, as.GetAddons())
	})

	t.Run("no keys is a no-op", func(t *testing.T) {
		as := tax.WithAddons("mx-cfdi-v4")
		as.AddAddons()
		assert.Equal(t, []cbc.Key{"mx-cfdi-v4"}, as.GetAddons())
	})

	t.Run("nil receiver is safe", func(t *testing.T) {
		var as *tax.Addons
		assert.NotPanics(t, func() { as.AddAddons("mx-cfdi-v4") })
	})
}

func TestAllAddonDefs(t *testing.T) {
	as := tax.AllAddonDefs()
	assert.NotEmpty(t, as)
}

func TestHasAddon(t *testing.T) {
	type withAddons struct {
		tax.Addons
	}

	t.Run("key present", func(t *testing.T) {
		obj := &withAddons{Addons: tax.WithAddons("mx-cfdi-v4")}
		assert.True(t, tax.HasAddon("mx-cfdi-v4").Check(obj))
	})

	t.Run("key absent", func(t *testing.T) {
		obj := &withAddons{Addons: tax.WithAddons("mx-cfdi-v4")}
		assert.False(t, tax.HasAddon("sa-zatca-v1").Check(obj))
	})

	t.Run("value does not implement GetAddons", func(t *testing.T) {
		assert.False(t, tax.HasAddon("mx-cfdi-v4").Check("not an addon carrier"))
		assert.False(t, tax.HasAddon("mx-cfdi-v4").Check(nil))
	})
}

// TestAddonIn covers tax.AddonIn and the negative-assertion pattern that
// addons such as EN 16931 use to skip a rule when a specific addon is in
// scope. The pattern mirrors `addonIsNotZATCA` in
// addons/eu/en16931/tax_combo.go
func TestAddonIn(t *testing.T) {
	t.Run("Check matches addon def", func(t *testing.T) {
		def := tax.AddonForKey("mx-cfdi-v4")
		require.NotNil(t, def)
		assert.True(t, tax.AddonIn("mx-cfdi-v4").Check(def))
		assert.True(t, tax.AddonIn("eu-en16931-v2017", "mx-cfdi-v4").Check(def))
	})

	t.Run("Check rejects non-matching addon def", func(t *testing.T) {
		def := tax.AddonForKey("mx-cfdi-v4")
		require.NotNil(t, def)
		assert.False(t, tax.AddonIn("sa-zatca-v1").Check(def))
	})

	t.Run("Check rejects non-AddonDef value", func(t *testing.T) {
		assert.False(t, tax.AddonIn("mx-cfdi-v4").Check("not an addon def"))
		assert.False(t, tax.AddonIn("mx-cfdi-v4").Check(nil))
	})

	// addonAbsent mirrors `addonIsNotZATCA` in
	// addons/eu/en16931/tax_combo.go: it returns true when the given addon
	// key is NOT present anywhere in the validation context.
	addonAbsent := func(ctx rules.Context, key cbc.Key) bool {
		return !ctx.Each(func(v any) bool {
			return tax.AddonIn(key).Check(v)
		})
	}

	t.Run("negative pattern: addon absent from context", func(t *testing.T) {
		var ctx rules.Context
		tax.AddonContext("eu-en16931-v2017")(&ctx)
		assert.True(t, addonAbsent(ctx, "sa-zatca-v1"))
	})

	t.Run("negative pattern: addon present in context", func(t *testing.T) {
		var ctx rules.Context
		tax.AddonContext("eu-en16931-v2017", "sa-zatca-v1")(&ctx)
		assert.False(t, addonAbsent(ctx, "sa-zatca-v1"))
	})

	t.Run("negative pattern: empty context", func(t *testing.T) {
		var ctx rules.Context
		assert.True(t, addonAbsent(ctx, "sa-zatca-v1"))
	})
}

func TestAddonsJSONSchemaEmbed(t *testing.T) {
	eg := `{
		"type": "array",
		"items": {
			"$ref": "https://gobl.org/draft-0/cbc/key"
		},
		"description": "AddonList defines the slice of keys to use for addons."
	}`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), js))

	al := tax.AddonList{}
	al.JSONSchemaExtend(js)

	assert.Greater(t, len(js.Items.OneOf), 1)
	ao := tax.AllAddonDefs()[0]
	assert.Equal(t, ao.Key.String(), js.Items.OneOf[0].Const)
}
