package cbc_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	var checks = []struct {
		key cbc.Key
		err bool
		val string
	}{
		{key: "h"},
		{key: "test"},
		{key: "1a"},
		{key: "ack1"},
		{key: cbc.Key("1a").With("foo"), val: "1a+foo"},
		{key: "-a", err: true},
		{key: "a-", err: true},
		{key: "1", err: true},
		{key: "-", err: true},
		{key: "+", err: true},
		{key: "a+", err: true},
	}
	for _, check := range checks {
		t.Run(check.key.String(), func(t *testing.T) {
			err := check.key.Validate()
			if check.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if check.val != "" {
				assert.Equal(t, check.key.String(), check.val)
			}
		})
	}
}

func TestStringsToKeys(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		list := []string{
			"key1",
			"key2",
			"key3"
		}
		out := cbc.StringsToKeys(list)
		assert.Equal(t, cbc.Key("key1"), out[0])
		assert.Equal(t, cbc.Key("key2"), out[1])
		assert.Equal(t, cbc.Key("key3"), out[2])
	})
}

func TestKeyHas(t *testing.T) {
	k := cbc.Key("standard")
	assert.True(t, k.Has("standard"))
	assert.False(t, k.Has("pro"))
	k = k.With("pro")
	assert.True(t, k.Has("standard"))
	assert.True(t, k.Has("pro"))
}

func TestKeyStrings(t *testing.T) {
	keys := []cbc.Key{"a", "b", "c"}
	assert.Equal(t, []string{"a", "b", "c"}, cbc.KeyStrings(keys))
}

func TestKeyHasPrefix(t *testing.T) {
	k := cbc.Key("standard")
	assert.True(t, k.HasPrefix("standard"))
	assert.False(t, k.HasPrefix("pro"))
	k = k.With("pro")
	assert.True(t, k.HasPrefix("standard"))
	assert.False(t, k.HasPrefix("pro"))
	k = cbc.KeyEmpty
	assert.False(t, k.HasPrefix("foo"))
}

func TestKeyIsEmpty(t *testing.T) {
	assert.True(t, cbc.KeyEmpty.IsEmpty())
	assert.False(t, cbc.Key("foo").IsEmpty())
}

func TestKeyIn(t *testing.T) {
	c := cbc.Key("standard")

	assert.True(t, c.In("pro", "reduced+eqs", "standard"))
	assert.False(t, c.In("pro", "reduced"))
}

func TestKeyPop(t *testing.T) {
	k := cbc.Key("a+b+c")
	assert.Equal(t, cbc.Key("a+b"), k.Pop())
	assert.Equal(t, cbc.Key("a"), k.Pop().Pop())
	assert.Equal(t, cbc.KeyEmpty, k.Pop().Pop().Pop())
}

func TestAppendUniqueKeys(t *testing.T) {
	keys := []cbc.Key{"a", "b", "c"}
	keys = cbc.AppendUniqueKeys(keys, "b", "d")
	assert.Equal(t, []cbc.Key{"a", "b", "c", "d"}, keys)
}

func TestHasValidKeyIn(t *testing.T) {
	k := cbc.Key("standard")
	err := validation.Validate(k, cbc.HasValidKeyIn("pro", "reduced+eqs"))
	assert.ErrorContains(t, err, "must be or start with a valid ke")

	err = validation.Validate(k, cbc.HasValidKeyIn("pro", "reduced+eqs", "standard"))
	assert.NoError(t, err)

	k = cbc.KeyEmpty
	err = validation.Validate(k, cbc.HasValidKeyIn("pro", "reduced+eqs", "standard"))
	assert.NoError(t, err)
}

func TestKeyJSONSchema(t *testing.T) {
	data := []byte(`{"description":"Text identifier to be used instead of a code for a more verbose but readable identifier.", "maxLength":64, "minLength":1, "pattern":"^(?:[a-z]|[a-z0-9][a-z0-9-+]*[a-z0-9])$", "title":"Key", "type":"string"}`)
	k := cbc.Key("standard")
	schema := k.JSONSchema()
	out, err := json.Marshal(schema)
	require.NoError(t, err)
	assert.JSONEq(t, string(data), string(out))
}
