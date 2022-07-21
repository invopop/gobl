package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestKeyIn(t *testing.T) {
	c := org.Key("standard")

	assert.True(t, c.In("pro", "reduced+eqs", "standard"))
	assert.False(t, c.In("pro", "reduced"))
}
