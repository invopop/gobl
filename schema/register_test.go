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

func TestRegistryTypeFor(t *testing.T) {
	r := newRegistry()
	base := ID("https://gobl.org/test02/schema")

	assert.NoError(t, r.add(base, TestRegistryType{}))

	id := ID("https://gobl.org/test02/schema/test-registry-type")
	typ := r.typeFor(id)
	assert.NotNil(t, typ)
	assert.Equal(t, "TestRegistryType", typ.Name())

	// Unknown ID should return nil
	typ = r.typeFor(ID("https://gobl.org/unknown"))
	assert.Nil(t, typ)
}

func TestRegistryLookupNotFound(t *testing.T) {
	r := newRegistry()
	type UnknownType struct{}
	id := r.lookup(UnknownType{})
	assert.Equal(t, UnknownID, id)
}

func TestGlobalType(t *testing.T) {
	// schema.Object is registered globally
	id := ID(GOBL.String() + "/schema/object")
	typ := schemas.typeFor(id)
	assert.NotNil(t, typ)
	assert.Equal(t, "Object", typ.Name())
}

func TestGlobalList(t *testing.T) {
	ids := List()
	assert.Greater(t, len(ids), 0)
}

func TestGlobalTypes(t *testing.T) {
	types := Types()
	assert.Greater(t, len(types), 0)
}
