package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestCodeIn(t *testing.T) {
	c := cbc.Code("FOO")

	assert.True(t, c.In("BAR", "FOO", "DOM"))
	assert.False(t, c.In("BAR", "DOM"))
}

func TestCodeValidation(t *testing.T) {
	c := cbc.Code("ABC123")
	err := c.Validate()
	assert.NoError(t, err)

	c = cbc.Code("12345678901234567890ABCD")
	err = c.Validate()
	assert.NoError(t, err)

	c = cbc.Code("")
	err = c.Validate()
	assert.NoError(t, err)

	c = cbc.Code("B-1234567")
	err = c.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "valid format")

	c = cbc.Code("12345678901234567890ABCDE")
	err = c.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "length must be between")

	c = cbc.Code("ab")
	err = c.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "valid format")
}
