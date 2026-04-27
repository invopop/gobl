package flow6

import (
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestExtValueNilPointer(t *testing.T) {
	assert.True(t, extValue((*tax.Extensions)(nil)).IsZero())
}

func TestExtValueUnknownType(t *testing.T) {
	assert.True(t, extValue(42).IsZero())
}

func TestExtValueFromValue(t *testing.T) {
	e := tax.ExtensionsOf(tax.ExtMap{"k": "v"})
	assert.False(t, extValue(e).IsZero())
}

func TestExtValueFromPointer(t *testing.T) {
	e := tax.ExtensionsOf(tax.ExtMap{"k": "v"})
	assert.False(t, extValue(&e).IsZero())
}
