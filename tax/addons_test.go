package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
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
		Addons: tax.WithAddons("gr-mydata-v1"),
		Name:   "Test",
	}
	assert.NotNil(t, ts.Addons)
	assert.Equal(t, "Test", ts.Name)

	assert.Equal(t, []cbc.Key{"gr-mydata-v1"}, ts.GetAddons())

	defs := ts.AddonDefs()
	assert.Len(t, defs, 1)
	assert.Equal(t, "gr-mydata-v1", defs[0].Key.String())

	ts.Addons = tax.WithAddons("gr-mydata-v1", "invalid-addon")

	err := rules.Validate(ts)
	assert.ErrorContains(t, err, "[GOBL-TAX-ADDONS-01] ($.$addons[1]) add-on must be registered")

	t.Run("test addon normalization", func(t *testing.T) {
		ts.Addons.List = tax.AddonList{"gr-mydata-v1", "gr-mydata-v1", "de-xrechnung-v3"}
		norm.Normalize(ts)
		assert.Equal(t, tax.AddonList{"gr-mydata-v1", "eu-en16931-v2017", "de-xrechnung-v3"}, ts.Addons.List)
	})
}

func TestAddonForKey(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		a := tax.AddonForKey("unknown")
		assert.Nil(t, a)
	})

	t.Run("found", func(t *testing.T) {
		a := tax.AddonForKey("gr-mydata-v1")
		require.NotNil(t, a)
		assert.NoError(t, rules.Validate(a))
	})
}

func TestSetAddons(t *testing.T) {
	var as tax.Addons
	as.SetAddons("gr-mydata-v1", "es-verifactu-v1")
	assert.Equal(t, []cbc.Key{"gr-mydata-v1", "es-verifactu-v1"}, as.GetAddons())

	// SetAddons replaces the existing list wholesale.
	as.SetAddons("pt-saft-v1")
	assert.Equal(t, []cbc.Key{"pt-saft-v1"}, as.GetAddons())
}

func TestAddAddons(t *testing.T) {
	t.Run("appends to an empty list", func(t *testing.T) {
		var as tax.Addons
		as.AddAddons("gr-mydata-v1")
		assert.Equal(t, []cbc.Key{"gr-mydata-v1"}, as.GetAddons())
	})

	t.Run("appends to an existing list", func(t *testing.T) {
		as := tax.WithAddons("gr-mydata-v1")
		as.AddAddons("es-verifactu-v1")
		assert.Equal(t, []cbc.Key{"gr-mydata-v1", "es-verifactu-v1"}, as.GetAddons())
	})

	t.Run("skips empty keys", func(t *testing.T) {
		var as tax.Addons
		as.AddAddons("", "gr-mydata-v1", "")
		assert.Equal(t, []cbc.Key{"gr-mydata-v1"}, as.GetAddons())
	})

	t.Run("skips keys already present", func(t *testing.T) {
		as := tax.WithAddons("gr-mydata-v1")
		as.AddAddons("gr-mydata-v1")
		assert.Equal(t, []cbc.Key{"gr-mydata-v1"}, as.GetAddons())
	})

	t.Run("de-duplicates within a single call", func(t *testing.T) {
		var as tax.Addons
		as.AddAddons("gr-mydata-v1", "es-verifactu-v1", "gr-mydata-v1")
		assert.Equal(t, []cbc.Key{"gr-mydata-v1", "es-verifactu-v1"}, as.GetAddons())
	})

	t.Run("no keys is a no-op", func(t *testing.T) {
		as := tax.WithAddons("gr-mydata-v1")
		as.AddAddons()
		assert.Equal(t, []cbc.Key{"gr-mydata-v1"}, as.GetAddons())
	})

	t.Run("nil receiver is safe", func(t *testing.T) {
		var as *tax.Addons
		assert.NotPanics(t, func() { as.AddAddons("gr-mydata-v1") })
	})
}

func TestAllAddonDefs(t *testing.T) {
	as := tax.AllAddonDefs()
	assert.NotEmpty(t, as)
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
