package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefinitionValidation(t *testing.T) {
	t.Run("check pattern", func(t *testing.T) {
		kd := &cbc.Definition{
			Key: "key",
			Name: i18n.String{
				i18n.EN: "Name",
				i18n.ES: "Nombre",
			},
			Pattern: "^[0-9]{3}$",
		}
		err := rules.Validate(kd)
		assert.NoError(t, err)

		kd.Pattern = ""
		err = rules.Validate(kd)
		assert.NoError(t, err)

		kd.Pattern = "[foo]["
		err = rules.Validate(kd)
		assert.ErrorContains(t, err, "[GOBL-CBC-DEFINITION-03] (pattern) pattern must be a valid regular expression")
	})
	t.Run("with source", func(t *testing.T) {
		kd := &cbc.Definition{
			Key: "key",
			Name: i18n.String{
				i18n.EN: "Name",
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.NewString("Test"),
					URL:   "http://example.com",
				},
			},
		}
		err := rules.Validate(kd)
		assert.NoError(t, err)
	})
	t.Run("with bad source", func(t *testing.T) {
		kd := &cbc.Definition{
			Key: "key",
			Name: i18n.String{
				i18n.EN: "Name",
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.NewString("Test"),
					URL:   "http:\\\\example.com",
				},
			},
		}
		err := rules.Validate(kd)
		assert.ErrorContains(t, err, "[GOBL-CBC-SOURCE-01] (sources[0].url) url is required and must be a URL")
	})
}

func TestDefinitionsWithValues(t *testing.T) {
	kd := &cbc.Definition{
		Key: "key",
		Name: i18n.String{
			i18n.EN: "Name",
			i18n.ES: "Nombre",
		},
		Desc: i18n.String{
			i18n.EN: "Description",
			i18n.ES: "Descripción",
		},
		Values: []*cbc.Definition{
			{
				Code: "CODE1",
				Name: i18n.String{
					i18n.EN: "Code 1",
					i18n.ES: "Código 1",
				},
			},
			{
				Code: "CODE2",
				Name: i18n.String{
					i18n.EN: "Code 2",
					i18n.ES: "Código 2",
				},
			},
			{
				Key: "key1",
				Name: i18n.String{
					i18n.EN: "Key 1",
					i18n.ES: "Clave 1",
				},
			},
			{
				Key: "key2",
				Name: i18n.String{
					i18n.EN: "Key 2",
					i18n.ES: "Clave 2",
				},
			},
		},
	}
	t.Run("for codes", func(t *testing.T) {
		assert.True(t, kd.HasCode("CODE1"))
		assert.False(t, kd.HasCode("INV"))
		vd := kd.CodeDef("CODE1")
		require.NotNil(t, vd)
		assert.Equal(t, "Code 1", vd.Name[i18n.EN])

		codes := cbc.DefinitionCodes(kd.Values)
		assert.Len(t, codes, 2)
		assert.Contains(t, codes, cbc.Code("CODE1"))
		assert.Contains(t, codes, cbc.Code("CODE2"))

		cdn := cbc.GetCodeDefinition("FOO", kd.Values)
		assert.Nil(t, cdn)

		cd := cbc.GetCodeDefinition("CODE2", kd.Values)
		require.NotNil(t, cd)
		assert.Equal(t, "Code 2", cd.Name[i18n.EN])
	})
	t.Run("for keys", func(t *testing.T) {
		assert.True(t, kd.HasKey("key1"))
		assert.False(t, kd.HasKey("bad"))
		vd := kd.KeyDef("key1")
		require.NotNil(t, vd)
		assert.Equal(t, "Key 1", vd.Name[i18n.EN])

		keys := cbc.DefinitionKeys(kd.Values)
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, cbc.Key("key1"))
		assert.Contains(t, keys, cbc.Key("key2"))

		kdn := cbc.GetKeyDefinition("bad", kd.Values)
		assert.Nil(t, kdn)

		kdn = cbc.GetKeyDefinition("key2", kd.Values)
		require.NotNil(t, kdn)
		assert.Equal(t, "Key 2", kdn.Name[i18n.EN])
	})
}

func TestDefinitionWithPattern(t *testing.T) {
	kd := &cbc.Definition{
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
	err := rules.Validate(kd)
	assert.NoError(t, err)

	kd.Pattern = "[foo]["
	err = rules.Validate(kd)
	assert.ErrorContains(t, err, "[GOBL-CBC-DEFINITION-03] (pattern) pattern must be a valid regular expression")
}
