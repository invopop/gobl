package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyDefinitionsWithValues(t *testing.T) {
	kd := &cbc.KeyDefinition{
		Key: "key",
		Name: i18n.String{
			i18n.EN: "Name",
			i18n.ES: "Nombre",
		},
		Desc: i18n.String{
			i18n.EN: "Description",
			i18n.ES: "Descripción",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "value",
				Name: i18n.String{
					i18n.EN: "Value",
					i18n.ES: "Valor",
				},
			},
		},
	}
	assert.True(t, kd.HasValue("value"))
	assert.False(t, kd.HasValue("invalid"))
	vd := kd.ValueDef("value")
	require.NotNil(t, vd)
	assert.Equal(t, "Value", vd.Name[i18n.EN])
}

func TestKeyDefinitionWithPattern(t *testing.T) {
	kd := &cbc.KeyDefinition{
		Key: "key",
		Name: i18n.String{
			i18n.EN: "Name",
			i18n.ES: "Nombre",
		},
		Desc: i18n.String{
			i18n.EN: "Description",
			i18n.ES: "Descripción",
		},
		Pattern: "^[0-9]{3}$",
	}
	err := kd.Validate()
	assert.NoError(t, err)

	kd.Pattern = "[foo]["
	err = kd.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pattern: error parsing regexp: missing closing ]: `[`")
}
