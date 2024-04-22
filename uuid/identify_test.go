package uuid_test

import (
	"testing"

	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIdentify(t *testing.T) {
	type Document struct {
		uuid.Identify
		Name string
	}
	doc := new(Document)
	id := uuid.V1()
	doc.SetUUID(id)
	assert.Equal(t, id, doc.GetUUID())
}
