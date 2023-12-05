package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestKeyOrCode(t *testing.T) {
	kc := cbc.KeyOrCode("IT")
	assert.Equal(t, "IT", kc.String())
	assert.NoError(t, kc.Validate())
	assert.Equal(t, cbc.Code("IT"), kc.Code())
	assert.Equal(t, cbc.KeyEmpty, kc.Key())

	kc = cbc.KeyOrCode("testing")
	assert.Equal(t, "testing", kc.String())
	assert.NoError(t, kc.Validate())
	assert.Equal(t, cbc.Key("testing"), kc.Key())
	assert.Equal(t, cbc.CodeEmpty, kc.Code())

	kc = cbc.KeyOrCode("INvalid")
	if assert.Error(t, kc.Validate()) {
		assert.Contains(t, kc.Validate().Error(), "value is not a key or code")
	}
	assert.Equal(t, cbc.CodeEmpty, kc.Code())
	assert.Equal(t, cbc.KeyEmpty, kc.Key())
}
