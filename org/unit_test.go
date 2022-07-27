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

func TestUnitUNECE(t *testing.T) {
	u := org.Unit("h")
	assert.Equal(t, u.UNECE(), org.Code("HUR"))

	u = org.UnitTetraBrik
	assert.Equal(t, u.UNECE(), org.CodeEmpty, "valid but no code")

	u = org.Unit("FOO")
	assert.Equal(t, u.UNECE(), org.CodeEmpty)
}
