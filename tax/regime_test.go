package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegimeJSONSchemaExtend(t *testing.T) {
	eg := `{
		"properties": {
			"$regime": {
				"$ref": "https://gobl.org/draft-0/cbc/key",
				"title": "Regime"
			}
		}
	}`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), js))

	r := tax.Regime{}
	r.JSONSchemaExtend(js)

	assert.Equal(t, js.Properties.Len(), 1)
	prop, ok := js.Properties.Get("$regime")
	require.True(t, ok)
	assert.Greater(t, len(prop.OneOf), 1)
	rd := tax.AllRegimeDefs()[0]
	assert.Equal(t, rd.Code().String(), prop.OneOf[0].Const)
}
