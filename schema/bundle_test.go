package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBundleSchema(t *testing.T) {
	t.Run("valid schema", func(t *testing.T) {
		data, err := BundleSchema("schemas/bill/invoice.json")
		require.NoError(t, err)
		require.NotEmpty(t, data)

		var result schemaFile
		require.NoError(t, json.Unmarshal(data, &result))

		assert.NotEmpty(t, result.ID)
		assert.NotEmpty(t, result.Ref)
		assert.Greater(t, len(result.Defs), 10, "should have many inlined definitions")

		// All $ref values should be internal
		assert.NotContains(t, string(data), `"$ref":"https://gobl.org/`)
		assert.NotContains(t, string(data), `"$ref": "https://gobl.org/`)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := BundleSchema("schemas/nonexistent/thing.json")
		assert.Error(t, err)
	})
}

func TestLoadSchemaFile(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		sf, err := loadSchemaFile("schemas/bill/invoice.json")
		require.NoError(t, err)
		assert.NotEmpty(t, sf.ID)
		assert.NotEmpty(t, sf.Defs)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := loadSchemaFile("schemas/nonexistent.json")
		assert.Error(t, err)
	})
}

func TestLoadSchemaByURL(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		sf, err := loadSchemaByURL(GOBL.String() + "/bill/invoice")
		require.NoError(t, err)
		assert.NotEmpty(t, sf.ID)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := loadSchemaByURL(GOBL.String() + "/nonexistent/type")
		assert.Error(t, err)
	})
}

func TestFindExternalRefs(t *testing.T) {
	prefix := GOBL.String()

	t.Run("with refs", func(t *testing.T) {
		defs := map[string]json.RawMessage{
			"Foo": json.RawMessage(`{"$ref":"` + prefix + `/org/party"}`),
			"Bar": json.RawMessage(`{"type":"string"}`),
		}
		refs := findExternalRefs(defs)
		assert.Len(t, refs, 1)
		assert.Equal(t, prefix+"/org/party", refs[0])
	})

	t.Run("no refs", func(t *testing.T) {
		defs := map[string]json.RawMessage{
			"Foo": json.RawMessage(`{"type":"string"}`),
		}
		refs := findExternalRefs(defs)
		assert.Empty(t, refs)
	})

	t.Run("nil defs", func(t *testing.T) {
		refs := findExternalRefs(nil)
		assert.Empty(t, refs)
	})
}

func TestRewriteRefs(t *testing.T) {
	urlToDef := map[string]string{
		"https://gobl.org/draft-0/org/party": "Party",
	}
	raw := json.RawMessage(`{"$ref":"https://gobl.org/draft-0/org/party"}`)
	result := rewriteRefs(raw, urlToDef)
	assert.Contains(t, string(result), `"#/$defs/Party"`)
	assert.NotContains(t, string(result), "https://gobl.org/")
}
