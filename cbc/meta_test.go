package cbc_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeta(t *testing.T) {
	err := rules.Validate(cbc.Key("test"))
	assert.NoError(t, err)

	err = rules.Validate(cbc.Key("bad_key"))
	assert.Error(t, err)
	assert.ErrorContains(t, err, "key must match the required pattern")
}

func TestMetaEquals(t *testing.T) {
	m1 := cbc.Meta{
		cbc.Key("test"): "bar",
	}
	m2 := cbc.Meta{
		cbc.Key("test"): "bar",
	}
	assert.True(t, m1.Equals(m2))

	m2 = cbc.Meta{
		cbc.Key("test"): "foo",
	}
	assert.False(t, m1.Equals(m2))

	m2 = cbc.Meta{
		cbc.Key("test"): "bar",
		cbc.Key("foo"):  "bar",
	}
	assert.False(t, m1.Equals(m2))
}

func TestMetaKeys(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var m cbc.Meta
		assert.Nil(t, m.Keys())
	})
	t.Run("sorted alphabetically", func(t *testing.T) {
		m := cbc.Meta{"zeta": "z", "alpha": "a", "mu": "m"}
		assert.Equal(t, []cbc.Key{"alpha", "mu", "zeta"}, m.Keys())
	})
}

func TestMetaValues(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var m cbc.Meta
		assert.Nil(t, m.Values())
	})
	t.Run("paired with sorted keys", func(t *testing.T) {
		m := cbc.Meta{"zeta": "z", "alpha": "a", "mu": "m"}
		assert.Equal(t, []string{"a", "m", "z"}, m.Values())
	})
}

func TestMetaAll(t *testing.T) {
	t.Run("iterates in alphabetical key order", func(t *testing.T) {
		m := cbc.Meta{"c": "cherry", "a": "apple", "b": "banana"}
		var keys []cbc.Key
		var vals []string
		for k, v := range m.All() {
			keys = append(keys, k)
			vals = append(vals, v)
		}
		assert.Equal(t, []cbc.Key{"a", "b", "c"}, keys)
		assert.Equal(t, []string{"apple", "banana", "cherry"}, vals)
	})
	t.Run("early break stops iteration", func(t *testing.T) {
		m := cbc.Meta{"a": "1", "b": "2", "c": "3"}
		var count int
		for range m.All() {
			count++
			if count == 2 {
				break
			}
		}
		assert.Equal(t, 2, count)
	})
	t.Run("empty meta yields nothing", func(t *testing.T) {
		var m cbc.Meta
		var count int
		for range m.All() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestMetaMarshalJSON(t *testing.T) {
	t.Run("empty marshals to null", func(t *testing.T) {
		var m cbc.Meta
		data, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})
	t.Run("keys are sorted alphabetically", func(t *testing.T) {
		// Two Meta maps with the same entries inserted in different orders
		// must produce byte-identical JSON output.
		m1 := cbc.Meta{"zeta": "z", "alpha": "a", "mu": "m"}
		m2 := cbc.Meta{"mu": "m", "zeta": "z", "alpha": "a"}
		b1, err := json.Marshal(m1)
		require.NoError(t, err)
		b2, err := json.Marshal(m2)
		require.NoError(t, err)
		assert.Equal(t, string(b1), string(b2))
		assert.Equal(t, `{"alpha":"a","mu":"m","zeta":"z"}`, string(b1))
	})
	t.Run("round-trip preserves entries", func(t *testing.T) {
		m := cbc.Meta{"a": "1", "b": "2"}
		data, err := json.Marshal(m)
		require.NoError(t, err)
		var out cbc.Meta
		require.NoError(t, json.Unmarshal(data, &out))
		assert.True(t, m.Equals(out))
	})
	t.Run("omitempty skips empty meta in struct", func(t *testing.T) {
		type wrapper struct {
			Meta cbc.Meta `json:"meta,omitempty"`
		}
		data, err := json.Marshal(wrapper{})
		require.NoError(t, err)
		assert.Equal(t, `{}`, string(data))
	})
}

func TestMetaJSONSchemaExtend(t *testing.T) {
	in := `{
		"additionalProperties": {
			"type": "string"
		},
		"type": "object"
	}`
	schema := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(in), schema))
	assert.NotNil(t, schema.AdditionalProperties)
	var m cbc.Meta
	m.JSONSchemaExtend(schema)
	assert.Nil(t, schema.AdditionalProperties)
	assert.NotEmpty(t, schema.PatternProperties)
}
