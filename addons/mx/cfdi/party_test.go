package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestMigratePartyIdentities(t *testing.T) {
	customer := &org.Party{
		Name: "Test Customer",
		Identities: []*org.Identity{
			{
				Key:  cfdi.ExtKeyFiscalRegime,
				Code: "608",
			},
			{
				Key:  cfdi.ExtKeyUse,
				Code: "G01",
			},
			{
				Key:  "random",
				Code: "12345678",
			},
		},
	}

	addon := tax.AddonForKey(cfdi.V4)
	addon.Normalizer(customer)

	assert.Len(t, customer.Identities, 1)
	assert.Len(t, customer.Ext, 2)
	assert.Equal(t, "608", customer.Ext[cfdi.ExtKeyFiscalRegime].String())
	assert.Equal(t, "G01", customer.Ext[cfdi.ExtKeyUse].String())
	assert.Equal(t, "12345678", customer.Identities[0].Code.String())
}
