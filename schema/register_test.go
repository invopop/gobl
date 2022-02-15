package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegistryType used for testing.
type TestRegistryType struct {
	Value string
}
type TestRegistryType2 struct {
	Value string
}

func TestRegistry(t *testing.T) {
	r := newRegistry()
	base := ID("https://gobl.org/test01/schema")

	assert.NoError(t, r.addWithAnchor(base, TestRegistryType{}))
	assert.NoError(t, r.add(base, TestRegistryType2{}))
	assert.Len(t, r.entries, 2)

	assert.EqualValues(t, "https://gobl.org/test01/schema#TestRegistryType", r.entries[0].id)
	assert.EqualValues(t, "https://gobl.org/test01/schema/test-registry-type2", r.entries[1].id)

	for _, id := range r.ids() {
		t.Logf("schema: %v", id)
	}
}
