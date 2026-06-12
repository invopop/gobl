package cef_test

import (
	"regexp"
	"testing"

	_ "github.com/invopop/gobl"
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

func TestVATEXPattern(t *testing.T) {
	ed := tax.ExtensionForKey("cef-vatex")
	require.NotNil(t, ed)
	require.NotEmpty(t, ed.Pattern)
	re := regexp.MustCompile(ed.Pattern)

	valid := []string{
		"VATEX-EU-79-C",
		"VATEX-EU-132-1A",
		"VATEX-EU-143-1FA",
		"VATEX-EU-AE",
		"VATEX-FR-FRANCHISE",
		"VATEX-SA-34-1",
		"VATEX-FR-CGI261-1", // outside the official list
	}
	for _, code := range valid {
		assert.True(t, re.MatchString(code), code)
	}

	invalid := []string{
		"EXEMPT-132",
		"VATEX",
		"VATEX-EU",
		"VATEX-eu-132",
		"vatex-EU-132",
		"VATEX-E-132",
		"VATEX-EU-132-",
	}
	for _, code := range invalid {
		assert.False(t, re.MatchString(code), code)
	}
}
