package tax_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
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

	err := ts.Addons.Validate()
	assert.ErrorContains(t, err, "1: addon 'invalid-addon' not registered")

	t.Run("test addon normalization", func(t *testing.T) {
		ts.Addons.List = []cbc.Key{"mx-cfdi-v4", "mx-cfdi-v4", "de-xrechnung-v3"}
		_ = tax.ExtractNormalizers(ts)
		assert.Equal(t, []cbc.Key{"mx-cfdi-v4", "eu-en16931-v2017", "de-xrechnung-v3"}, ts.Addons.List)
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
		assert.NoError(t, a.Validate())
	})
}

func TestAllAddonDefs(t *testing.T) {
	as := tax.AllAddonDefs()
	assert.NotEmpty(t, as)
}

func TestAddonWithContext(t *testing.T) {
	ad := tax.AddonForKey("mx-cfdi-v4")
	ctx := ad.WithContext(context.Background())

	vs := tax.Validators(ctx)
	assert.Len(t, vs, 1)
	// no reliable way to check the function is actually the same :-(
}

func TestAddonsJSONSchemaEmbed(t *testing.T) {
	eg := `{
		"properties": {
			"$addons": {
				"items": {
            		"$ref": "https://gobl.org/draft-0/cbc/key",
					"type": "array",
					"title": "Addons",
					"description": "Addons defines a list of keys used to identify tax addons that apply special\nnormalization, scenarios, and validation rules to a document."
				}
			}
		}
	}`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), js))

	as := tax.Addons{}
	as.JSONSchemaExtend(js)

	assert.Equal(t, js.Properties.Len(), 1)
	prop, ok := js.Properties.Get("$addons")
	require.True(t, ok)
	assert.Greater(t, len(prop.Items.OneOf), 1)
	ao := tax.AllAddonDefs()[0]
	assert.Equal(t, ao.Key.String(), prop.Items.OneOf[0].Const)
}
