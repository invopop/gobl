package untdid_test

import (
	"testing"

	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestUnitCode(t *testing.T) {
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
		{cbc.KeyEmpty, ""},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.code, untdid.UnitCode(tt.unit))
	}
}

func TestUnitKey(t *testing.T) {
	assert.Equal(t, org.UnitHour, untdid.UnitKey("HUR"))
	assert.Equal(t, org.UnitUnit, untdid.UnitKey("XUN"))
	assert.Equal(t, org.UnitPortion, untdid.UnitKey("13"))
	assert.Equal(t, cbc.KeyEmpty, untdid.UnitKey("XZZ"))
	assert.Equal(t, cbc.KeyEmpty, untdid.UnitKey(""))
}

func TestUnitCodeCoverage(t *testing.T) {
	mapped := 0
	for _, def := range org.UnitDefinitions {
		unit := def.Key
		code := untdid.UnitCode(unit)
		if assert.NotEmpty(t, code, "unit %s should be mapped", unit) {
			assert.Equal(t, unit, untdid.UnitKey(code))
			mapped++
		}
	}
	assert.Equal(t, len(org.UnitDefinitions), mapped)
}
