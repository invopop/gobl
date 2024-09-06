package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyDefinitionsWithCodes(t *testing.T) {
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
		Codes: []*cbc.CodeDefinition{
			{
				Code: cbc.Code("CODE"),
				Name: i18n.String{
					i18n.EN: "Code",
					i18n.ES: "Código",
				},
			},
		},
	}
	assert.True(t, kd.HasCode(cbc.Code("CODE")))
	assert.False(t, kd.HasCode(cbc.Code("INVALID")))
	cd := kd.CodeDef(cbc.Code("CODE"))
	require.NotNil(t, cd)
	assert.Equal(t, "Code", cd.Name[i18n.EN])
}

func TestKeyDefinitionsWithKeys(t *testing.T) {
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
		Keys: []*cbc.KeyDefinition{
			{
				Key: cbc.Key("code"),
				Name: i18n.String{
					i18n.EN: "Code",
					i18n.ES: "Código",
				},
			},
		},
	}
	assert.True(t, kd.HasKey(cbc.Key("code")))
	assert.False(t, kd.HasKey(cbc.Key("invalid")))
	cd := kd.KeyDef(cbc.Key("code"))
	require.NotNil(t, cd)
	assert.Equal(t, "Code", cd.Name[i18n.EN])
}

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
