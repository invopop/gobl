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

	u = org.Unit("random-something")
	assert.Equal(t, u.UNECE(), cbc.Code("ZZ"))

	u = org.Unit("")
	assert.Equal(t, u.UNECE(), cbc.CodeEmpty)
}

func TestUnitJSONSChema(t *testing.T) {
	u := org.Unit("h")
	schema := u.JSONSchema()
	assert.Equal(t, schema.Title, "Unit")
	assert.Equal(t, schema.OneOf[0].Const, org.UnitDefinitions[0].Unit)
	last := schema.OneOf[len(schema.OneOf)-1]
	assert.Equal(t, last.Pattern, org.UnitPatternUNECE)
}

func TestNewUnitsValidation(t *testing.T) {
	// Test all new units validate successfully
	newUnits := []org.Unit{
		org.UnitWeek,
		org.UnitYear,
		org.UnitDecilitre,
		org.UnitKilolitre,
		org.UnitCentigram,
		org.UnitLinearMetre,
		org.UnitLinearFoot,
		org.UnitFlatFee,
		org.UnitBlock,
		org.UnitPacket,
		org.UnitBundle,
	}

	for _, u := range newUnits {
		assert.NoError(t, u.Validate(), "unit %s should validate", u)
	}
}

func TestNewUnitsUNECE(t *testing.T) {
	// Test UNECE mappings for new units
	tests := []struct {
		unit     org.Unit
		expected cbc.Code
	}{
		{org.UnitWeek, "WEE"},
		{org.UnitYear, "ANN"},
		{org.UnitDecilitre, "DLT"},
		{org.UnitKilolitre, cbc.CodeEmpty},
		{org.UnitCentigram, "CGM"},
		{org.UnitLinearMetre, "LM"},
		{org.UnitLinearFoot, "LF"},
		{org.UnitFlatFee, cbc.CodeEmpty},
		{org.UnitBlock, cbc.CodeEmpty},
		{org.UnitPacket, "PA"},
		{org.UnitBundle, "BE"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.unit.UNECE(), "UNECE code for %s should be %s", tt.unit, tt.expected)
	}
}
