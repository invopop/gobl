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
		TaxID: &tax.Identity{
			Country: "MX",
			Code:    "ZZZ010101ZZZ",
			Zone:    "65000",
		},
	}

	addon := tax.AddonForKey(cfdi.KeyV4)
	addon.Normalizer(customer)

	assert.Empty(t, customer.Identities)
	assert.Len(t, customer.Ext, 3)
	assert.Equal(t, "608", customer.Ext[cfdi.ExtKeyFiscalRegime].String())
	assert.Equal(t, "G01", customer.Ext[cfdi.ExtKeyUse].String())
	assert.Equal(t, "65000", customer.Ext[cfdi.ExtKeyPostCode].String())
	assert.Empty(t, customer.TaxID.Zone) //nolint:staticcheck
}
