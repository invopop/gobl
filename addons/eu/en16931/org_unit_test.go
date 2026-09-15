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
		{org.UnitPortion, "13"},
		{org.UnitSixPack, "NMP"},
		{org.UnitTetraBrik, ""},
		{cbc.KeyEmpty, ""},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.code, en16931.UnitToUNTDID(tt.unit))
	}
}

func TestUnitFromUNTDID(t *testing.T) {
	assert.Equal(t, org.UnitHour, en16931.UnitFromUNTDID("HUR"))
	assert.Equal(t, org.UnitUnit, en16931.UnitFromUNTDID("XUN"))
	assert.Equal(t, org.UnitPortion, en16931.UnitFromUNTDID("13"))
	assert.Equal(t, cbc.KeyEmpty, en16931.UnitFromUNTDID("XZZ"))
	assert.Equal(t, cbc.KeyEmpty, en16931.UnitFromUNTDID(""))

	// A lossy code belongs to no unit, so it never converts back.
	assert.Equal(t, cbc.KeyEmpty, en16931.UnitFromUNTDID("NMP"))
}

func TestUnitUNTDIDMapCoverage(t *testing.T) {
	// Units mapped to a broader code, which cannot convert back.
	lossy := map[cbc.Key]bool{
		org.UnitSixPack: true,
	}
	// Units with no UNTDID code at all.
	unmapped := map[cbc.Key]bool{
		org.UnitTetraBrik: true,
	}
	mapped := 0
	for _, def := range org.UnitDefinitions {
		unit := def.Key
		code := en16931.UnitToUNTDID(unit)
		switch {
		case unmapped[unit]:
			assert.Empty(t, code, "non-standard unit %s should not be mapped", unit)
		case lossy[unit]:
			if assert.NotEmpty(t, code, "unit %s should be mapped", unit) {
				assert.NotEqual(t, unit, en16931.UnitFromUNTDID(code),
					"lossy unit %s should not convert back", unit)
			}
		default:
			if assert.NotEmpty(t, code, "unit %s should be mapped", unit) {
				assert.Equal(t, unit, en16931.UnitFromUNTDID(code))
				mapped++
			}
		}
	}
	assert.Equal(t, 87, mapped)
}
