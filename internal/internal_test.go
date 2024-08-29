package internal_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/internal"
	"github.com/stretchr/testify/assert"
)

func TestSignatureFlag(t *testing.T) {
	ctx := context.Background()
	assert.False(t, internal.IsSigned(ctx))

	ctx = internal.SignedContext(ctx)
	assert.True(t, internal.IsSigned(ctx))
}
