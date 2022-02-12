package gobl_test

import (
	"testing"

	"github.com/invopop/gobl"
	"github.com/stretchr/testify/assert"
)

func TestNewHeader(t *testing.T) {
	h := gobl.NewHeader()
	assert.False(t, h.UUID.IsZero())
	assert.True(t, h.UUID.Version() == 1)
	assert.NotPanics(t, func() {
		h.Meta["foo"] = "bar"
		h.Tags = append(h.Tags, "foo")
		h.Stamps = append(h.Stamps, &gobl.Stamp{})
	})
}
