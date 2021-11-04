package gobl_test

import (
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/region"
	"github.com/stretchr/testify/assert"
)

var testKey = dsig.NewES256Key()

func TestEnvelopePayload(t *testing.T) {
	m := &note.Message{
		Content: "This is test content.",
	}
	e := gobl.NewEnvelope(region.ES)
	if assert.NotNil(t, e.Head) {
		assert.NotEmpty(t, e.Head.UUID, "empty header uuid")
	}
	assert.NotNil(t, e.Document)
	if assert.NotNil(t, e.Region()) {
		assert.Equal(t, region.ES, e.Region().Code())
	}

	if err := e.Insert(m); err != nil {
		t.Errorf("failed to insert payload: %v", err)
		return
	}

	assert.Equal(t, e.Head.Type, "note.Message", "type should match")
	if assert.NotNil(t, e.Head.Digest) {
		assert.Equal(t, e.Head.Digest.Algorithm, dsig.DigestSHA256, "unexpected digest algorithm")
		assert.Equal(t, e.Head.Digest.Value, "2c24a95a0141a3e74c7a910fecda9ed69d67396f4e3000999a9e3acc722208ef", "digest should be the same")
	}

	assert.Empty(t, e.Signatures)
	assert.NoError(t, e.Sign(testKey), "signing envelope")
	assert.NotEmpty(t, e.Signatures, "expected a signature")

	assert.NoError(t, e.Verify(), "did not expect verify error")

	nm := new(note.Message)
	assert.NoError(t, e.Extract(nm))
	assert.Equal(t, m.Content, nm.Content, "content mismatch")
}
