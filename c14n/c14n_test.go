package c14n_test

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"github.com/invopop/gobl/c14n"
	"github.com/stretchr/testify/assert"
)

func TestJSONToArray(t *testing.T) {
	// some random data found on the internets
	data := `{
		"actors": [
			{
				"name": "Tom Cruise",
				"age": 56,
				"born_at": "Syracuse, NY",
				"birthdate": "July 3, 1962",
				"photo": "https://jsonformatter.org/img/tom-cruise.jpg",
				"wife": null,
				"weight": 67.5,
				"has_children": true,
				"has_grey_hair": false,
				"children": [
					"Suri",
					"Isabella Jane",
					"Connor"
				],
				"icon": "🤩"
			},
			{
				"name": "Robert Downey Jr.",
				"age": 53,
				"born_at": "New York City, NY",
				"birthdate": "April 4, 1965",
				"photo": "https://jsonformatter.org/img/Robert-Downey-Jr.jpg",
				"wife": {
					"name": "Susan Downey",
					"age": 35
				},
				"weight": 77.1,
				"has_children": true,
				"has_grey_hair": false,
				"children": [
					"Indio Falconer",
					"Avri Roel",
					"Exton Elias"
				]
			}
		]
	}`
	r := strings.NewReader(data)
	obj, err := c14n.UnmarshalJSON(r)
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
		return
	}

	d, err := obj.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect JSON marshal error: %v", err.Error())
		return
	}
	s := fmt.Sprintf("%x", sha256.Sum256(d))
	if s != "f35a55c7bba2df8438802603db442976a8238ceb0a610d1eea38cae1b9fd9013" {
		t.Logf("marshaled data:\n%v\n", string(d))
		t.Errorf("unexpected sum, please check marshaled data, got: %v", s)
	}
}

func TestMarshalJSON(t *testing.T) {
	obj := struct {
		Title string `json:"title"`
		Idx   int64  `json:"idx"`
		Body  string `json:"body,omitempty"`
	}{
		Title: "test",
		Idx:   1,
		Body:  "Test body to play around with",
	}
	d, err := c14n.MarshalJSON(obj)
	assert.NoError(t, err)
	out := `{"body":"Test body to play around with","idx":1,"title":"test"}`
	assert.Contains(t, string(d), out)

	t.Run("with encoding error", func(t *testing.T) {
		invalidObj := make(chan int) // channels can't be marshaled to JSON
		_, err := c14n.MarshalJSON(invalidObj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "encoding")
	})
}

func TestCanonicalJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		data := `{"name":"test","age":25}`
		r := strings.NewReader(data)
		result, err := c14n.CanonicalJSON(r)
		assert.NoError(t, err)
		assert.Equal(t, `{"age":25,"name":"test"}`, string(result))
	})

	t.Run("invalid JSON", func(t *testing.T) {
		data := `{"name":"test","age":}`
		r := strings.NewReader(data)
		_, err := c14n.CanonicalJSON(r)
		assert.Error(t, err)
	})

	t.Run("empty input", func(t *testing.T) {
		r := strings.NewReader("")
		out, err := c14n.CanonicalJSON(r)
		assert.NoError(t, err)
		assert.Nil(t, out)
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("valid object", func(t *testing.T) {
		data := `{"name":"test","age":25}`
		r := strings.NewReader(data)
		result, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.IsType(t, &c14n.Object{}, result)
	})

	t.Run("valid array", func(t *testing.T) {
		data := `["test", 123, true]`
		r := strings.NewReader(data)
		result, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.IsType(t, &c14n.Array{}, result)
	})

	t.Run("simple string", func(t *testing.T) {
		data := `"test"`
		r := strings.NewReader(data)
		result, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.Equal(t, c14n.String("test"), result)
	})

	t.Run("simple number", func(t *testing.T) {
		data := `42`
		r := strings.NewReader(data)
		result, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.Equal(t, c14n.Integer(42), result)
	})

	t.Run("simple float", func(t *testing.T) {
		data := `42.5`
		r := strings.NewReader(data)
		result, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.Equal(t, c14n.Float(42.5), result)
	})

	t.Run("simple boolean", func(t *testing.T) {
		data := `true`
		r := strings.NewReader(data)
		result, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.Equal(t, c14n.Bool(true), result)
	})

	t.Run("null value", func(t *testing.T) {
		data := `null`
		r := strings.NewReader(data)
		result, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.Equal(t, c14n.Null{}, result)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		data := `{"name":}`
		r := strings.NewReader(data)
		_, err := c14n.UnmarshalJSON(r)
		assert.Error(t, err)
	})

	t.Run("empty input", func(t *testing.T) {
		r := strings.NewReader("")
		out, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.Nil(t, out)
	})

	t.Run("nested structures", func(t *testing.T) {
		data := `{"users":[{"name":"John","age":30},{"name":"Jane","age":25}],"count":2}`
		r := strings.NewReader(data)
		result, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.IsType(t, &c14n.Object{}, result)
	})

	t.Run("object with non-string key error", func(t *testing.T) {
		// This is a bit tricky to test directly since JSON objects always have string keys
		// But we can test the error path by creating malformed JSON
		data := `{123:"value"}`
		r := strings.NewReader(data)
		_, err := c14n.UnmarshalJSON(r)
		assert.Error(t, err)
	})

	t.Run("large numbers", func(t *testing.T) {
		data := `9223372036854775807` // max int64
		r := strings.NewReader(data)
		result, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.Equal(t, c14n.Integer(9223372036854775807), result)
	})

	t.Run("number too large for int64", func(t *testing.T) {
		data := `92233720368547758070` // larger than max int64
		r := strings.NewReader(data)
		result, err := c14n.UnmarshalJSON(r)
		assert.NoError(t, err)
		assert.IsType(t, c14n.Float(0), result)
	})

}
