package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegistryType used for testing.
type TestRegistryType struct {
	Value string
}

func TestRegistry(t *testing.T) {
	r := newRegistry()
	id := ID("https://gobl.org/test01/schema")
	x := new(TestRegistryType)
	assert.NoError(t, r.addWithAnchor(id, x))
	assert.Len(t, r.entries, 1)

	assert.EqualValues(t, "https://gobl.org/test01/schema#TestRegistryType", r.entries[0].id)

	for _, id := range r.ids() {
		t.Logf("schema: %v", id)
	}
}
