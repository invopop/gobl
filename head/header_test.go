package head_test

import (
	"testing"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/stretchr/testify/assert"
)

func TestNewHeader(t *testing.T) {
	h := head.NewHeader()
	assert.False(t, h.UUID.IsZero())
	assert.True(t, h.UUID.Version() == 1)
	assert.NotPanics(t, func() {
		h.Meta["foo"] = "bar"
		h.Tags = append(h.Tags, "foo")
		h.Stamps = append(h.Stamps, &head.Stamp{})
	}, "header and meta hash should have been initialized")
}

func TestHeaderAddStamp(t *testing.T) {
	h := head.NewHeader()
	h.AddStamp(&head.Stamp{Provider: "foo", Value: "bar"})
	assert.Len(t, h.Stamps, 1)
	h.AddStamp(&head.Stamp{Provider: "foo", Value: "baz"})
	assert.Len(t, h.Stamps, 1)
	assert.Equal(t, "baz", h.Stamps[0].Value)
	h.AddStamp(&head.Stamp{Provider: "bar", Value: "bax"})
	assert.Len(t, h.Stamps, 2)
	assert.Equal(t, "bax", h.Stamps[1].Value)
}

func TestHeaderContains(t *testing.T) {
	h1 := head.NewHeader()
	h2 := head.NewHeader()

	// UUID
	assert.False(t, h1.Contains(h2))
	h2.UUID = h1.UUID
	assert.True(t, h1.Contains(h2))

	// Digest
	h1.Digest = dsig.NewSHA256Digest([]byte("testing"))
	h2.Digest = dsig.NewSHA256Digest([]byte("testing 2"))
	assert.False(t, h1.Contains(h2))
	h2.Digest = h1.Digest
	assert.True(t, h1.Contains(h2))

	// Stamps
	h1.AddStamp(&head.Stamp{Provider: "foo", Value: "bar"})
	h1.AddStamp(&head.Stamp{Provider: "foo2", Value: "bar"})
	assert.True(t, h1.Contains(h2))
	h2.AddStamp(&head.Stamp{Provider: "foo", Value: "bar"})
	assert.True(t, h1.Contains(h2))
	h2.AddStamp(&head.Stamp{Provider: "foo2", Value: "bar"})
	assert.True(t, h1.Contains(h2))
	h2.AddStamp(&head.Stamp{Provider: "foo3", Value: "bar"})
	assert.False(t, h1.Contains(h2))
	h1.AddStamp(&head.Stamp{Provider: "foo3", Value: "bar"})
	assert.True(t, h1.Contains(h2))

	// Tags
	h1.Tags = append(h1.Tags, "foo")
	assert.True(t, h1.Contains(h2))
	h2.Tags = append(h2.Tags, "foo")
	assert.True(t, h1.Contains(h2))
	h2.Tags = append(h2.Tags, "foo2")
	assert.False(t, h1.Contains(h2))
	h1.Tags = append(h1.Tags, "foo2")
	assert.True(t, h1.Contains(h2))

	// Meta
	h1.Meta["foo"] = "bar"
	assert.True(t, h1.Contains(h2))
	h2.Meta["foo"] = "bar"
	assert.True(t, h1.Contains(h2))
	h2.Meta["foo2"] = "bar"
	assert.False(t, h1.Contains(h2))
	h1.Meta["foo2"] = "bar"
	assert.True(t, h1.Contains(h2))

	// Notes
	h1.Notes = "test notes"
	assert.True(t, h1.Contains(h2))
	h2.Notes = "bad notes"
	assert.False(t, h1.Contains(h2))
	h2.Notes = h1.Notes
	assert.True(t, h1.Contains(h2))
}
