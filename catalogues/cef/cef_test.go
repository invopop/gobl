package cef_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	// Test that the catalogue is registered
	ed := tax.ExtensionForKey("cef-vatex")
	assert.NotNil(t, ed)
	assert.Equal(t, "cef-vatex", ed.Key.String())
}

func TestVATEXDefinition(t *testing.T) {
	ed := tax.ExtensionForKey("cef-vatex")
	require.NotNil(t, ed)

	t.Run("code list version", func(t *testing.T) {
		assert.Equal(t, "8.0", ed.Meta["version"])
	})

	t.Run("sources", func(t *testing.T) {
		require.Len(t, ed.Sources, 3)
		for _, src := range ed.Sources {
			assert.Contains(t, src.URL, "digital-building-blocks")
		}
	})

	t.Run("codes are enumerated, not matched by pattern", func(t *testing.T) {
		assert.Empty(t, ed.Pattern)
		assert.Len(t, ed.Values, 88)
	})

	t.Run("known codes", func(t *testing.T) {
		valid := []cbc.Code{
			"VATEX-EU-79-C",
			"VATEX-EU-132-1A",
			"VATEX-EU-135-1", // added in version 8
			"VATEX-EU-143-1FA",
			"VATEX-EU-AE",
			"VATEX-FR-FRANCHISE",
			"VATEX-FR-CGI261-1",
			"VATEX-FR-CGI275",
			"VATEX-FR-CGI261D-1BIS",
		}
		for _, code := range valid {
			assert.NotNil(t, ed.CodeDef(code), code.String())
		}
	})

	t.Run("unknown codes", func(t *testing.T) {
		invalid := []cbc.Code{
			"EXEMPT-132",
			"VATEX",
			"VATEX-EU",
			"VATEX-EU-132-1Z",
			"VATEX-SA-34-1", // Saudi codes belong to the ZATCA catalogue
			"VATEX-ES-EXEMPT",
		}
		for _, code := range invalid {
			assert.Nil(t, ed.CodeDef(code), code.String())
		}
	})

	t.Run("every code carries its source metadata", func(t *testing.T) {
		for _, cd := range ed.Values {
			assert.NotEmpty(t, cd.Name.In("en"), cd.Code.String())
			assert.NotEmpty(t, cd.Desc.In("en"), cd.Code.String())
			assert.Contains(t, []string{"true", "false"}, cd.Meta["deprecated"], cd.Code.String())
			assert.NotEmpty(t, cd.Meta["first-version"], cd.Code.String())
			assert.NotEmpty(t, cd.Meta["nationality"], cd.Code.String())
		}
	})

	t.Run("code metadata", func(t *testing.T) {
		cd := ed.CodeDef("VATEX-EU-G")
		require.NotNil(t, cd)
		assert.Equal(t, "Export outside the EU", cd.Name.In("en"))
		assert.Equal(t, "EU", cd.Meta["nationality"])
		assert.Equal(t, "false", cd.Meta["deprecated"])
		assert.Equal(t, "1", cd.Meta["first-version"])
		assert.Empty(t, cd.Meta["last-version"])
		assert.Equal(t, "Only use with VAT category code G", cd.Meta["remark"])
	})
}
