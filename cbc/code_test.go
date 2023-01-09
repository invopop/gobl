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
	c := cbc.Code("ABC")
	if err := c.Validate(); err != nil {
		t.Errorf("did not expect error: %v", err)
	}
	c = cbc.Code("abc")
	if err := c.Validate(); err == nil {
		t.Errorf("expected a validation error")
	}
	c = cbc.Code("ab")
	if err := c.Validate(); err == nil {
		t.Errorf("expected a validation error")
	}
}
