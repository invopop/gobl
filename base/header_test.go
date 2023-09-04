package base_test

import (
	"testing"

	"github.com/invopop/gobl/base"
	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestNewHeader(t *testing.T) {
	h := base.NewHeader()
	assert.False(t, h.UUID.IsZero())
	assert.True(t, h.UUID.Version() == 1)
	assert.NotPanics(t, func() {
		h.Meta["foo"] = "bar"
		h.Tags = append(h.Tags, "foo")
		h.Stamps = append(h.Stamps, &cbc.Stamp{})
	}, "header and meta hash should have been initialized")
}

func TestHeaderAddStamp(t *testing.T) {
	h := base.NewHeader()
	h.AddStamp(&cbc.Stamp{Provider: "foo", Value: "bar"})
	assert.Len(t, h.Stamps, 1)
	h.AddStamp(&cbc.Stamp{Provider: "foo", Value: "baz"})
	assert.Len(t, h.Stamps, 1)
	assert.Equal(t, "baz", h.Stamps[0].Value)
	h.AddStamp(&cbc.Stamp{Provider: "bar", Value: "bax"})
	assert.Len(t, h.Stamps, 2)
	assert.Equal(t, "bax", h.Stamps[1].Value)
}
