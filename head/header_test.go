package head_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewHeader(t *testing.T) {
	h := head.NewHeader()
	assert.False(t, h.UUID.IsZero())
	assert.Equal(t, uuid.Version(7), h.UUID.Version())
	assert.NotPanics(t, func() {
		h.Meta["foo"] = "bar"
		h.Tags = append(h.Tags, "foo")
		h.Stamps = append(h.Stamps, &head.Stamp{})
	}, "header and meta hash should have been initialized")
}

func TestHeaderValidation(t *testing.T) {
	h := head.NewHeader()
	h.Stamps = []*head.Stamp{
		{Provider: "foo", Value: "bar"},
	}
	err := h.Validate()
	assert.ErrorContains(t, err, "dig: cannot be blank; stamps: must be blank")

	h.Digest = dsig.NewSHA256Digest([]byte("testing"))

	ctx := internal.SignedContext(context.Background())
	err = h.ValidateWithContext(ctx)
	assert.NoError(t, err)

	h.Stamps = append(h.Stamps, &head.Stamp{Provider: "foo", Value: "bar"})
	err = h.ValidateWithContext(ctx)
	assert.ErrorContains(t, err, "stamps: duplicate stamp 'foo'")
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

func TestHeaderAddLink(t *testing.T) {
	h := head.NewHeader()
	assert.Len(t, h.Links, 0)
	h.AddLink(&head.Link{Key: "foo", URL: "bar.com"})
	assert.Len(t, h.Links, 1)
}

func TestHeaderLink(t *testing.T) {
	t.Run("without category", func(t *testing.T) {
		h := head.NewHeader()
		h.AddLink(&head.Link{Key: "foo", URL: "bar.com"})
		l := h.Link("", "foo")
		assert.NotNil(t, l)
		assert.Equal(t, "bar.com", l.URL)
		l = h.Link("", "baa")
		assert.Nil(t, l)
	})
	t.Run("with category", func(t *testing.T) {
		h := head.NewHeader()
		h.AddLink(&head.Link{Category: head.LinkCategoryKeyPortal, Key: "foo", URL: "bar.com"})
		l := h.Link(head.LinkCategoryKeyPortal, "foo")
		assert.NotNil(t, l)
		assert.Equal(t, "bar.com", l.URL)
		l = h.Link(head.LinkCategoryKeyPortal, "baa")
		assert.Nil(t, l)
	})
}

func TestHeaderStamp(t *testing.T) {
	h := head.NewHeader()
	h.AddStamp(&head.Stamp{Provider: "foo", Value: "boo"})
	h.AddStamp(&head.Stamp{Provider: "foo2", Value: "bling"})
	st := h.GetStamp("foo")
	assert.NotNil(t, st)
	assert.Equal(t, "boo", st.Value)
	st = h.GetStamp("bad")
	assert.Nil(t, st)
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
	h1.AddStamp(&head.Stamp{Provider: "foo", Value: "boo"})
	h1.AddStamp(&head.Stamp{Provider: "foo2", Value: "bling"})
	assert.True(t, h1.Contains(h2))
	h2.AddStamp(&head.Stamp{Provider: "foo", Value: "boo"})
	assert.True(t, h1.Contains(h2))
	h2.AddStamp(&head.Stamp{Provider: "foo2", Value: "bling"})
	assert.True(t, h1.Contains(h2))
	h2.AddStamp(&head.Stamp{Provider: "foo3", Value: "bow"})
	assert.False(t, h1.Contains(h2))
	h1.AddStamp(&head.Stamp{Provider: "foo3", Value: "bow"})
	assert.True(t, h1.Contains(h2))

	// Links
	h1.AddLink(&head.Link{Key: "foo", URL: "bar.com"})
	h1.AddLink(&head.Link{Key: "foo2", URL: "bar2.com"})
	assert.True(t, h1.Contains(h2))
	h2.AddLink(&head.Link{Key: "foo", URL: "bar.com"})
	assert.True(t, h1.Contains(h2))
	h2.AddLink(&head.Link{Key: "foo2", URL: "bar2.com"})
	assert.True(t, h1.Contains(h2))
	h2.AddLink(&head.Link{Key: "foo3", URL: "bar3.com"})
	assert.False(t, h1.Contains(h2))
	h1.AddLink(&head.Link{Key: "foo3", URL: "bar3.com"})

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
	h1.Meta["foo"] = "bang"
	assert.True(t, h1.Contains(h2))
	h2.Meta["foo"] = "bang"
	assert.True(t, h1.Contains(h2))
	h2.Meta["foo2"] = "bo2"
	assert.False(t, h1.Contains(h2))
	h1.Meta["foo2"] = "bo2"
	assert.True(t, h1.Contains(h2))

	// Notes
	h1.Notes = "test notes"
	assert.True(t, h1.Contains(h2))
	h2.Notes = "bad notes"
	assert.False(t, h1.Contains(h2))
	h2.Notes = h1.Notes
	assert.True(t, h1.Contains(h2))
}
