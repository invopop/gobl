package schema_test

import (
	"testing"

	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	tt := schema.Type("test/bar")
	id := tt.ID()

	assert.EqualValues(t, "https://gobl.org/"+schema.VERSION+"/test/bar", id)
	assert.EqualValues(t, "test/bar", id.Type())

	id = schema.ID("bad-url")
	err := id.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "valid URL")
}
