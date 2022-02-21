package gobl

import (
	"testing"

	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocument(t *testing.T) {
	msg := &note.Message{
		Content: "test message",
	}

	doc := new(Document)

	err := doc.Insert(msg)
	require.NoError(t, err)

	id := schema.Lookup(&note.Message{})
	assert.Contains(t, id.String(), "https://gobl.org/")
	assert.Contains(t, id.String(), "/note#Message")
}
