package schema_test

import (
	"testing"

	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// See also document tests performed in `gobl` package.

func TestObjectUUID(t *testing.T) {
	msg := &note.Message{
		UUID:    uuid.V1(),
		Title:   "just a test",
		Content: "this is a test message",
	}

	obj, err := schema.NewObject(msg)
	require.NoError(t, err)

	assert.Equal(t, msg.UUID, obj.UUID())
}
