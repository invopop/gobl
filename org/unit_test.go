package org_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitValidation(t *testing.T) {
	assert.True(t, org.HasValidUnitKey.Check(cbc.Key("h")))
	assert.False(t, org.HasValidUnitKey.Check(cbc.Key("XUN")))
	assert.False(t, org.HasValidUnitKey.Check(cbc.Key("X")))
	assert.False(t, org.HasValidUnitKey.Check(cbc.Key("XUNX")))
}

func TestUnitDefinitions(t *testing.T) {
	for _, def := range org.UnitDefinitions {
		assert.Empty(t, def.Map)
	}

	def := cbc.GetKeyDefinition(org.UnitHour, org.UnitDefinitions)
	assert.NotNil(t, def)
	assert.Equal(t, "Hours", def.Name.String())
	assert.Empty(t, def.Map)

	def = cbc.GetKeyDefinition(org.UnitKilogram, org.UnitDefinitions)
	assert.NotNil(t, def)
	assert.Equal(t, "kg", def.Meta[org.UnitMetaKeySymbol])
}

func TestUnitJSONSchema(t *testing.T) {
	t.Run("with property", func(t *testing.T) {
		schema := &jsonschema.Schema{Properties: jsonschema.NewProperties()}
		schema.Properties.Set("unit", &jsonschema.Schema{})
		org.ExtendUnitKeySchema(schema, "unit")
		unit, ok := schema.Properties.Get("unit")
		assert.True(t, ok)
		assert.Equal(t, unit.OneOf[0].Const, org.UnitDefinitions[0].Key)
		assert.Len(t, unit.OneOf, len(org.UnitDefinitions))
	})
	t.Run("missing property", func(t *testing.T) {
		schema := &jsonschema.Schema{Properties: jsonschema.NewProperties()}
		schema.Properties.Set("other", &jsonschema.Schema{})
		org.ExtendUnitKeySchema(schema, "unit")
		other, ok := schema.Properties.Get("other")
		require.True(t, ok)
		assert.Empty(t, other.OneOf)
		_, ok = schema.Properties.Get("unit")
		assert.False(t, ok)
	})
}

func TestNewUnitsValidation(t *testing.T) {
	// Test all new units validate successfully
	newUnits := []cbc.Key{
		org.UnitWeek,
		org.UnitYear,
		org.UnitDecilitre,
		org.UnitKilolitre,
		org.UnitCentigram,
		org.UnitLinearMetre,
		org.UnitLinearFoot,
		org.UnitBlock,
		org.UnitPacket,
		org.UnitBundle,
	}

	for _, u := range newUnits {
		assert.True(t, org.HasValidUnitKey.Check(u), "unit %s should validate", u)
	}
}
