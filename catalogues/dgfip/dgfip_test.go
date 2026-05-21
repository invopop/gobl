package dgfip_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/catalogues/dgfip"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	ext := tax.ExtensionForKey(dgfip.ExtKeyBillingMode)
	assert.NotNil(t, ext)

	ed := ext.CodeDef(dgfip.BillingModeM1)
	assert.NotNil(t, ed)
	assert.Equal(t, "M1", ed.Code.String())
	assert.Equal(t, "Mixed - Deposit invoice", ed.Name.String())

	assert.Len(t, ext.Values, 13)
}
