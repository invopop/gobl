package gobl_test

import (
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/note"
	"github.com/stretchr/testify/assert"
)

var parseExampleDoc = `{
	"$schema": "https://gobl.org/draft-0/note/message",
	"title": "Test Message",
	"content": "We hope you like this test message!"
}`
var parseExampleEnvelope = `{
	"$schema": "https://gobl.org/draft-0/envelope",
	"head": {
		"uuid": "1aec7332-ee67-11ec-8f98-3e7e00ce5635",
		"dig": {
			"alg": "sha256",
			"val": "45ac3115c8569a1789e58af8d0dc91ef3baa1fb71daaf38f5aef94f82b4d0033"
		}
	},"doc":` + parseExampleDoc + `,"sigs": [
		"eyJhbGciOiJFUzI1NiIsImtpZCI6IjY0ZjBhMDFhLWUzOTktNDZmMy04ZmI0LTgxNDNhMzc3MGY4OSJ9.eyJ1dWlkIjoiMWFlYzczMzItZWU2Ny0xMWVjLThmOTgtM2U3ZTAwY2U1NjM1IiwiZGlnIjp7ImFsZyI6InNoYTI1NiIsInZhbCI6IjQ1YWMzMTE1Yzg1NjlhMTc4OWU1OGFmOGQwZGM5MWVmM2JhYTFmYjcxZGFhZjM4ZjVhZWY5NGY4MmI0ZDAwMzMifX0.QJQS83I7Q7Hl0pySv49pNo1fr0tL7IIgI73HTsqiwO-cOjZak1CUnpE6-us6797Q856spghseX0cD4yogAhmqg"
	]
}`

func TestParse(t *testing.T) {
	t.Run("basic note message", func(t *testing.T) {
		doc, err := gobl.Parse([]byte(parseExampleDoc))
		assert.NoError(t, err)
		assert.IsType(t, &note.Message{}, doc)
		n, ok := doc.(*note.Message)
		assert.True(t, ok)
		assert.Equal(t, "Test Message", n.Title)
	})

	t.Run("check digest", func(t *testing.T) {
		doc, err := gobl.Parse([]byte(parseExampleEnvelope))
		assert.NoError(t, err)
		assert.IsType(t, &gobl.Envelope{}, doc)
		env, ok := doc.(*gobl.Envelope)
		assert.True(t, ok)
		assert.Contains(t, env.Head.Digest.Value, "45ac3115c8569a1")
		n := env.Extract().(*note.Message)
		assert.NotNil(t, n)
		assert.Equal(t, "Test Message", n.Title)

	})

	t.Run("unknown schema", func(t *testing.T) {
		_, err := gobl.Parse([]byte(`{"$schema": "https://gobl.org/draft-0/unknown"}`))
		assert.Error(t, err)
		assert.ErrorIs(t, err, gobl.ErrUnmarshal)
		assert.ErrorContains(t, err, "unmarshal: json: Unmarshal(nil)")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := gobl.Parse([]byte(`{"$schema": "https://gobl.org/draft-0/note/message", "title": "Test Message"`))
		assert.Error(t, err)
		assert.ErrorIs(t, err, gobl.ErrUnmarshal)
		assert.ErrorContains(t, err, "unmarshal: unexpected end of JSON input")
	})

	t.Run("empty schema", func(t *testing.T) {
		_, err := gobl.Parse([]byte(`{"$schema": ""}`))
		assert.Error(t, err)
		assert.ErrorIs(t, err, gobl.ErrUnknownSchema)
		assert.ErrorContains(t, err, "unknown-schema")
	})
}
