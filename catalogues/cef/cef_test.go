package cef_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// Test that the catalogue is registered
	ed := tax.ExtensionForKey("cef-vatex")
	assert.NotNil(t, ed)
	assert.Equal(t, "cef-vatex", ed.Key.String())
}
