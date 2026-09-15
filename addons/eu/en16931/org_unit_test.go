package en16931_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

// The mapping itself is covered by the UNTDID catalogue; these only check that
// the addon's converters still reach it.
func TestUnitToUNTDID(t *testing.T) {
	assert.Equal(t, cbc.Code("HUR"), en16931.UnitToUNTDID(org.UnitHour))
	assert.Equal(t, cbc.Code("13"), en16931.UnitToUNTDID(org.UnitPortion))
	assert.Equal(t, cbc.CodeEmpty, en16931.UnitToUNTDID(cbc.KeyEmpty))
}

func TestUnitFromUNTDID(t *testing.T) {
	assert.Equal(t, org.UnitHour, en16931.UnitFromUNTDID("HUR"))
	assert.Equal(t, org.UnitPortion, en16931.UnitFromUNTDID("13"))
	assert.Equal(t, cbc.KeyEmpty, en16931.UnitFromUNTDID("XZZ"))
}
