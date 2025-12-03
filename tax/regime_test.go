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
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$id": "https://gobl.org/draft-0/tax/regime-code",
		"type": "string"
	}`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), js))

	rc := tax.RegimeCode("")
	rc.JSONSchemaExtend(js)

	assert.Greater(t, len(js.OneOf), 1)
	rd := tax.AllRegimeDefs()[0]
	assert.Equal(t, rd.Code().String(), js.OneOf[0].Const)
}

func TestRegimeCode(t *testing.T) {
	rc := tax.RegimeCode("US")
	assert.Equal(t, "US", rc.String())
	assert.Equal(t, "US", rc.Code().String())
}
