package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// See also document tests performed in `gobl` package.

func TestObjectUUID(t *testing.T) {
	tr := &tax.Regime{} // doesn't have a UUID field!
	obj, err := schema.NewObject(tr)
	require.NoError(t, err)
	assert.Equal(t, uuid.Empty, obj.UUID())

	msg := &note.Message{
		UUID:    uuid.V1(),
		Title:   "just a test",
		Content: "this is a test message",
	}

	obj, err = schema.NewObject(msg)
	require.NoError(t, err)

	assert.Equal(t, msg.UUID, obj.UUID())
}

func TestObjectInjectUUID(t *testing.T) {
	tr := &tax.Regime{} // doesn't have a UUID field!
	id := uuid.V1()
	obj, err := schema.NewObject(tr)
	require.NoError(t, err)
	assert.NotPanics(t, func() {
		obj.InjectUUID(id)
	})

	msg := &note.Message{
		Title:   "just a test",
		Content: "this is a test message",
	}
	obj, err = schema.NewObject(msg)
	require.NoError(t, err)

	obj.InjectUUID(id)

	assert.Equal(t, id, obj.UUID())
	data, err := json.Marshal(obj)
	require.NoError(t, err)
	assert.Equal(t, `{"$schema":"https://gobl.org/draft-0/note/message","uuid":"`+id.String()+`","title":"just a test","content":"this is a test message"}`, string(data))
}
