package tax_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeExtMap(t *testing.T) {
	var em tax.ExtMap

	em2 := tax.NormalizeExtMap(em)
	assert.Nil(t, em2)

	em = tax.ExtMap{
		"key": "",
	}
	em2 = tax.NormalizeExtMap(em)
	assert.Nil(t, em2)

	em = tax.ExtMap{
		"key": "foo",
		"bar": "",
	}
	em2 = tax.NormalizeExtMap(em)
	assert.NotNil(t, em2)
	assert.Len(t, em2, 1)
	assert.Equal(t, cbc.KeyOrCode("foo"), em2["key"])
}
