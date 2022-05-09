package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestUnitValidation(t *testing.T) {
	u := org.Unit("h")
	assert.NoError(t, u.Validate())

	u = org.Unit("FOO")
	err := u.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "valid format")
	}
}
