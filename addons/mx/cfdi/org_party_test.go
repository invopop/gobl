package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/norm"
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

	norm.Normalize(customer, tax.AddonContext(cfdi.V4))

	assert.Len(t, customer.Identities, 1)
	assert.Equal(t, 2, customer.Ext.Len())
	assert.Equal(t, "608", customer.Ext.Get(cfdi.ExtKeyFiscalRegime).String())
	assert.Equal(t, "G01", customer.Ext.Get(cfdi.ExtKeyUse).String())
	assert.Equal(t, "12345678", customer.Identities[0].Code.String())
}

func TestNormalizePartyWithNilIdentities(t *testing.T) {
	t.Run("party with nil identity in array", func(t *testing.T) {
		customer := &org.Party{
			Name: "Test Customer",
			Identities: []*org.Identity{
				nil,
				{
					Key:  cfdi.ExtKeyFiscalRegime,
					Code: "608",
				},
				nil,
			},
		}

		norm.Normalize(customer, tax.AddonContext(cfdi.V4))

		// Should not panic with nil identities
		assert.Len(t, customer.Identities, 0)
		assert.Equal(t, 1, customer.Ext.Len())
		assert.Equal(t, "608", customer.Ext.Get(cfdi.ExtKeyFiscalRegime).String())
	})

	t.Run("party with only nil identities", func(t *testing.T) {
		customer := &org.Party{
			Name:       "Test Customer",
			Identities: []*org.Identity{nil, nil},
		}

		norm.Normalize(customer, tax.AddonContext(cfdi.V4))

		// Should not panic with only nil identities
		assert.Len(t, customer.Identities, 0)
		assert.True(t, customer.Ext.IsZero())
	})
}
