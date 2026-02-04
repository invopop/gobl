package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemIdentityNormalization(t *testing.T) {
	addon := tax.AddonForKey(cfdi.V4)
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "123456",
			Expected: "12345600",
		},
		{
			Code:     "12345678",
			Expected: "12345678",
		},
		{
			Code:     "1234567",
			Expected: "1234567",
		},
	}
	for _, ts := range tests {
		item := &org.Item{Ext: tax.Extensions{cfdi.ExtKeyProdServ: ts.Code}}
		addon.Normalizer(item)
		assert.Equal(t, ts.Expected, item.Ext[cfdi.ExtKeyProdServ])
	}

	// In context of invoice
	inv := validInvoice()
	inv.Lines[0].Item.Ext[cfdi.ExtKeyProdServ] = "010101"
	err := inv.Calculate()
	require.NoError(t, err)
	assert.Equal(t, cbc.Code("01010100"), inv.Lines[0].Item.Ext[cfdi.ExtKeyProdServ])
}

func TestItemIdentityMigration(t *testing.T) {
	inv := validInvoice()

	inv.Lines[0].Item.Ext = nil
	inv.Lines[0].Item.Identities = []*org.Identity{
		{
			Key:  cfdi.ExtKeyProdServ,
			Code: "01010101",
		},
		{
			Key:  "other",
			Code: "1234",
		},
	}

	err := inv.Calculate()
	require.NoError(t, err)
	assert.Equal(t, cbc.Code("01010101"), inv.Lines[0].Item.Ext[cfdi.ExtKeyProdServ])
	assert.Equal(t, "1234", inv.Lines[0].Item.Identities[0].Code.String())
}

func TestItemNilIdentityHandling(t *testing.T) {
	t.Run("item with nil identity in array", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Identities = []*org.Identity{nil}
		require.NoError(t, inv.Calculate())
		// Should not panic with nil identity
	})

	t.Run("item with mixed nil and valid identities", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Ext = nil
		inv.Lines[0].Item.Identities = []*org.Identity{
			nil,
			{
				Key:  cfdi.ExtKeyProdServ,
				Code: "01010101",
			},
			nil,
			{
				Key:  "other",
				Code: "5678",
			},
		}
		require.NoError(t, inv.Calculate())
		// Should not panic and should migrate valid identities
		assert.Equal(t, cbc.Code("01010101"), inv.Lines[0].Item.Ext[cfdi.ExtKeyProdServ])
		assert.Len(t, inv.Lines[0].Item.Identities, 1)
		assert.Equal(t, "5678", inv.Lines[0].Item.Identities[0].Code.String())
	})
}
