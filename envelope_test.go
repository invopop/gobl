package gobl_test

import (
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
)

var testKey = dsig.NewES256Key()

func TestEnvelopePayload(t *testing.T) {
	m := note.NewMessage()
	m.Content = "This is test content."

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

	if assert.NotNil(t, e.Head.Digest) {
		assert.Equal(t, e.Head.Digest.Algorithm, dsig.DigestSHA256, "unexpected digest algorithm")
		assert.Equal(t, e.Head.Digest.Value, "94ba910f5f4baf7caf245c1ba15442272a02162b62321d50967243e51b4e73d9", "digest should be the same")
	}

	assert.Empty(t, e.Signatures)
	assert.NoError(t, e.Sign(testKey), "signing envelope")
	assert.NotEmpty(t, e.Signatures, "expected a signature")

	assert.NoError(t, e.Verify(), "did not expect verify error")

	nm := new(note.Message)
	assert.NoError(t, e.Extract(nm))
	assert.Equal(t, m.Content, nm.Content, "content mismatch")
}

func TestEnvelopeExtract(t *testing.T) {
	e := &gobl.Envelope{}
	inv := new(bill.Invoice)
	err := e.Extract(inv)
	assert.ErrorIs(t, err, gobl.ErrNoDocument)
}

func TestEnvelopeValidate(t *testing.T) {
	tests := []struct {
		name string
		env  *gobl.Envelope
		want string
	}{
		{
			name: "no head nor version",
			env:  &gobl.Envelope{},
			want: "$schema: cannot be blank; doc: cannot be blank; head: cannot be blank.",
		},
		{
			name: "missing sig, draft",
			env: &gobl.Envelope{
				Def: gobl.EnvelopeType.Def(),
				Head: &gobl.Header{
					Digest: &dsig.Digest{},
					Region: "ES",
					Draft:  true,
					UUID:   uuid.NewV1(),
				},
				Document: new(gobl.Document),
			},
		},
		{
			name: "missing sig, draft",
			env: &gobl.Envelope{
				Def: gobl.EnvelopeType.Def(),
				Head: &gobl.Header{
					Digest: &dsig.Digest{},
					Region: "ES",
					UUID:   uuid.NewV1(),
				},
				Document: new(gobl.Document),
			},
			want: "sigs: cannot be blank.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.env.Validate()
			if tt.want == "" && err == nil {
				return
			}
			assert.EqualError(t, err, tt.want)
		})
	}
}
