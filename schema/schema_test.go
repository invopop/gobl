package schema_test

import (
	"testing"

	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	id := schema.GOBL.Add("test/bar")
	base := "https://gobl.org/" + schema.VERSION

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
