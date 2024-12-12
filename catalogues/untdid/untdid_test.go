package untdid_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	ext := tax.ExtensionForKey("untdid-tax-category")
	assert.NotNil(t, ext)

	ext = tax.ExtensionForKey("untdid-charge")
	ed := ext.CodeDef("AAS")
	assert.NotNil(t, ed)
	assert.Equal(t, "AAS", ed.Code.String())
	assert.Equal(t, "Acceptance", ed.Name.String())
}
