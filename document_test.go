package gobl_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocument(t *testing.T) {
	data := []byte(`{"$schema":"https://gobl.org/draft0/note/message","content":"test message"}`)

	doc := new(gobl.Document)
	err := json.Unmarshal(data, doc)
	require.NoError(t, err)

	typ, err := doc.Type()
	assert.NoError(t, err)
	assert.EqualValues(t, "note/message", typ)
}
