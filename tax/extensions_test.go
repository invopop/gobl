package tax_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeExtensions(t *testing.T) {
	var em tax.Extensions

	em2 := tax.NormalizeExtensions(em)
	assert.Nil(t, em2)

	em = tax.Extensions{
		"key": "",
	}
	em2 = tax.NormalizeExtensions(em)
	assert.Nil(t, em2)

	em = tax.Extensions{
		"key": "foo",
		"bar": "",
	}
	em2 = tax.NormalizeExtensions(em)
	assert.NotNil(t, em2)
	assert.Len(t, em2, 1)
	assert.Equal(t, "foo", em2["key"].String())
}

func TestExtValue(t *testing.T) {
	ev := tax.ExtValue("IT")
	assert.Equal(t, "IT", ev.String())
	assert.Equal(t, cbc.Code("IT"), ev.Code())
	assert.Equal(t, cbc.KeyEmpty, ev.Key())

	ev = tax.ExtValue("testing")
	assert.Equal(t, "testing", ev.String())
	assert.Equal(t, cbc.Key("testing"), ev.Key())
	assert.Equal(t, cbc.CodeEmpty, ev.Code())

	ev = tax.ExtValue("A string")
	assert.Equal(t, cbc.CodeEmpty, ev.Code())
	assert.Equal(t, cbc.KeyEmpty, ev.Key())
	assert.Equal(t, "A string", ev.String())
}
