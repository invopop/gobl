package head_test

import (
	"testing"

	"github.com/invopop/gobl/head"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	type embedOpts struct {
		head.CorrectionOptions
	}

	h := head.NewHeader()
	opt := head.WithHead(h)

	// Test that the option is applied to the options struct.
	opts := new(head.CorrectionOptions)
	opt(opts)

	assert.Equal(t, h, opts.Head)

	no := new(embedOpts)
	opt(no)
	assert.Equal(t, h, no.Head, "should be applied to embedded struct")
}
