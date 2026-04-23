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
		item := &org.Item{Ext: tax.ExtensionsOf(tax.ExtMap{cfdi.ExtKeyProdServ: ts.Code})}
		addon.Normalizer(item)
		assert.Equal(t, ts.Expected, item.Ext.Get(cfdi.ExtKeyProdServ))
	}

	// In context of invoice
	inv := validInvoice()
	inv.Lines[0].Item.Ext = inv.Lines[0].Item.Ext.Set(cfdi.ExtKeyProdServ, "010101")
	err := inv.Calculate()
	require.NoError(t, err)
	assert.Equal(t, cbc.Code("01010100"), inv.Lines[0].Item.Ext.Get(cfdi.ExtKeyProdServ))
}

func TestItemIdentityMigration(t *testing.T) {
	inv := validInvoice()

	inv.Lines[0].Item.Ext = tax.Extensions{}
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
	assert.Equal(t, cbc.Code("01010101"), inv.Lines[0].Item.Ext.Get(cfdi.ExtKeyProdServ))
	assert.Equal(t, "1234", inv.Lines[0].Item.Identities[0].Code.String())
}
