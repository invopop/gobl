package schema_test

import (
	"testing"

	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	id := schema.GOBL.Add("test/bar")
	base := schema.BaseURL + schema.Version

	assert.EqualValues(t, base+"/test/bar", id)

	id = id.Anchor("foo")
	assert.EqualValues(t, base+"/test/bar#foo", id)

	id = id.Anchor("new")
	assert.EqualValues(t, base+"/test/bar#new", id)

	id = schema.ID("bad-url")
	err := id.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "valid URL")

}

func TestExtract(t *testing.T) {
	base := schema.BaseURL + schema.Version + ``
	data := []byte(`{"$schema":"` + base + `/test/bar","random":"message"}`)

	id, err := schema.Extract(data)
	assert.NoError(t, err)
	assert.Equal(t, "https://gobl.org/draft-0/test/bar", id.String())

	data = []byte(`{"random":"message"}`)
	id, err = schema.Extract(data)
	assert.NoError(t, err)
	assert.Equal(t, schema.UnknownID, id)

	data = []byte(`bad-data`)
	_, err = schema.Extract(data)
	assert.Error(t, err)
}

func TestInsert(t *testing.T) {
	id := schema.ID(schema.BaseURL + schema.Version + `/test/bar`)
	data := []byte(`{"random":"message"}`)
	var err error
	data, err = schema.Insert(id, data)
	assert.NoError(t, err)
	assert.Equal(t, "{\"$schema\":\"https://gobl.org/draft-0/test/bar\",\"random\":\"message\"}", string(data))
}
