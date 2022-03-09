package schema

import (
	"reflect"
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

func TestFind(t *testing.T) {
	const (
		idInvoice = "https://gobl.org/draft-0/bill/invoice"
	)
	type Invoice struct{}
	r := &registry{
		entries: []*entry{
			{
				id:  idInvoice,
				typ: reflect.TypeOf(Invoice{}),
			},
		},
	}

	t.Run("exact schema match", func(t *testing.T) {
		conf, got := r.find(idInvoice)
		if conf != 1 {
			t.Errorf("Unexpected confidence: %v", conf)
		}
		if got != idInvoice {
			t.Errorf("Unexpected result: %v", got)
		}
	})
	t.Run("exact type match", func(t *testing.T) {
		conf, got := r.find("Invoice")
		if conf != 1 {
			t.Errorf("Unexpected confidence: %v", conf)
		}
		if got != idInvoice {
			t.Errorf("Unexpected result: %v", got)
		}
	})
	t.Run("exact type match with package", func(t *testing.T) {
		conf, got := r.find("schema.Invoice")
		if conf != 1 {
			t.Errorf("Unexpected confidence: %v", conf)
		}
		if got != idInvoice {
			t.Errorf("Unexpected result: %v", got)
		}
	})
}
