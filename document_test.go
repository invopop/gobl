package gobl_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/yaml"
)

func TestDocument(t *testing.T) {
	msg := &note.Message{
		Content: "test message",
	}

	doc, err := gobl.NewDocument(msg)
	require.NoError(t, err)

	id := schema.Lookup(&note.Message{})
	assert.Contains(t, id.String(), "https://gobl.org/")
	assert.Contains(t, id.String(), "/note/message")

	dig := "82a5cddc56f069ff17705f310161dd17cd8b00d94728e6be3fafdad980522a27"
	assert.Equal(t, id, doc.Schema())
	sha, err := doc.Digest()
	require.NoError(t, err)
	assert.Equal(t, dig, sha.Value)
	assert.Equal(t, doc.Instance(), msg)

	data, err := json.Marshal(doc)
	require.NoError(t, err)
	assert.Equal(t, `{"$schema":"`+id.String()+`","content":"test message"}`, string(data))
	digest := dsig.NewSHA256Digest(data) // this works as the JSON is very simple!
	assert.Equal(t, dig, digest.Value)

	doc = new(gobl.Document)
	err = json.Unmarshal(data, doc)
	require.NoError(t, err)

	assert.Equal(t, doc.Schema(), id)
	sha, err = doc.Digest()
	require.NoError(t, err)
	assert.Equal(t, dig, sha.Value)

	obj, ok := doc.Instance().(*note.Message)
	assert.True(t, ok)
	assert.Equal(t, msg.Content, obj.Content)
}

func TestDocumentValidation(t *testing.T) {
	msg := &note.Message{}

	doc, err := gobl.NewDocument(msg)
	require.NoError(t, err)

	err = doc.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "content: cannot be blank")
	}

	doc = new(gobl.Document)
	data, err := os.ReadFile("./regimes/es/examples/invoice-es-es.yaml")
	require.NoError(t, err)
	err = yaml.Unmarshal(data, doc)
	require.NoError(t, err)

	inv := doc.Instance().(*bill.Invoice)
	inv.Code = "" // blank, which will not be accepted if not a draft
	require.NoError(t, doc.Normalize())
	assert.NoError(t, doc.Validate())
	inv.IssueDate = cal.Date{}
	err = doc.Validate()
	if assert.Error(t, err) {
		// Double check to make sure validation working
		assert.Contains(t, err.Error(), "issue_date: required")
	}

}
