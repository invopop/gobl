package org_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestUnitValidation(t *testing.T) {
	u := org.Unit("h")
	assert.NoError(t, u.Validate())

	u = org.Unit("XUN")
	assert.NoError(t, u.Validate())

	u = org.Unit("X")
	err := u.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "must be a valid value or UN/ECE code")
	}

	u = org.Unit("XUNX")
	err = u.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "must be a valid value or UN/ECE code")
	}
}

func TestUnitUNECE(t *testing.T) {
	u := org.Unit("h")
	assert.Equal(t, u.UNECE(), cbc.Code("HUR"))

	u = org.UnitTetraBrik
	assert.Equal(t, u.UNECE(), cbc.CodeEmpty, "valid but no code")

	u = org.Unit("XUN")
	assert.Equal(t, u.UNECE(), cbc.Code("XUN"))
}
