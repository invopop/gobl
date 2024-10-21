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
		},
	}

	addon := tax.AddonForKey(cfdi.V4)
	addon.Normalizer(customer)

	assert.Empty(t, customer.Identities)
	assert.Len(t, customer.Ext, 2)
	assert.Equal(t, "608", customer.Ext[cfdi.ExtKeyFiscalRegime].String())
	assert.Equal(t, "G01", customer.Ext[cfdi.ExtKeyUse].String())
}
