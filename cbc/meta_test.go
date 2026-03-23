package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestMeta(t *testing.T) {
	err := rules.Validate(cbc.Key("test"))
	assert.NoError(t, err)

	err = rules.Validate(cbc.Key("bad_key"))
	assert.Error(t, err)
	assert.ErrorContains(t, err, "key must match the required pattern")
}

func TestMetaEquals(t *testing.T) {
	m1 := cbc.Meta{
		cbc.Key("test"): "bar",
	}
	m2 := cbc.Meta{
		cbc.Key("test"): "bar",
	}
	assert.True(t, m1.Equals(m2))

	m2 = cbc.Meta{
		cbc.Key("test"): "foo",
	}
	assert.False(t, m1.Equals(m2))

	m2 = cbc.Meta{
		cbc.Key("test"): "bar",
		cbc.Key("foo"):  "bar",
	}
	assert.False(t, m1.Equals(m2))
}
