package en16931_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestUnitToUNTDID(t *testing.T) {
	tests := []struct {
		unit cbc.Key
		code cbc.Code
	}{
		{org.UnitWeek, "WEE"},
		{org.UnitYear, "ANN"},
		{org.UnitDecilitre, "DLT"},
		{org.UnitKilolitre, "K6"},
		{org.UnitCentigram, "CGM"},
		{org.UnitLinearMetre, "LM"},
		{org.UnitLinearFoot, "LF"},
		{org.UnitBlock, "XOK"},
		{org.UnitPacket, "XPA"},
		{org.UnitBundle, "XBE"},
		{org.UnitPortion, ""},
		{cbc.KeyEmpty, ""},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.code, en16931.UnitToUNTDID(tt.unit))
	}
}

func TestUnitFromUNTDID(t *testing.T) {
	assert.Equal(t, org.UnitHour, en16931.UnitFromUNTDID("HUR"))
	assert.Equal(t, org.UnitUnit, en16931.UnitFromUNTDID("XUN"))
	assert.Equal(t, cbc.KeyEmpty, en16931.UnitFromUNTDID("XZZ"))
	assert.Equal(t, cbc.KeyEmpty, en16931.UnitFromUNTDID(""))
}

func TestUnitUNTDIDMapCoverage(t *testing.T) {
	unmapped := map[cbc.Key]bool{
		org.UnitPortion:   true,
		org.UnitSixPack:   true,
		org.UnitTetraBrik: true,
	}
	mapped := 0
	for _, def := range org.UnitDefinitions {
		unit := def.Key
		code := en16931.UnitToUNTDID(unit)
		if unmapped[unit] {
			assert.Empty(t, code, "non-standard unit %s should not be mapped", unit)
			continue
		}
		if assert.NotEmpty(t, code, "unit %s should be mapped", unit) {
			assert.Equal(t, unit, en16931.UnitFromUNTDID(code))
			mapped++
		}
	}
	assert.Equal(t, 86, mapped)
}
