package untdid_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	ext := tax.ExtensionForKey("untdid-tax-category")
	assert.NotNil(t, ext)

	ext = tax.ExtensionForKey("untdid-charge")
	ed := ext.CodeDef("AAS")
	assert.NotNil(t, ed)
	assert.Equal(t, "AAS", ed.Code.String())
	assert.Equal(t, "Acceptance", ed.Name.String())

	ext = tax.ExtensionForKey("untdid-unit")
	assert.NotNil(t, ext)
	assert.Equal(t, "^[A-Z0-9]{2,3}$", ext.Pattern)
}

func TestDeliveryTerms(t *testing.T) {
	ext := tax.ExtensionForKey(untdid.ExtKeyDeliveryTerms)
	require.NotNil(t, ext)
	assert.Equal(t, []cbc.Code{"1", "2"}, cbc.DefinitionCodes(ext.Values))
	assert.Equal(t, []string{"bill/delivery-details"}, ext.Usage.Schemas)
}
