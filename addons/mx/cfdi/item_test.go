package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemValidation(t *testing.T) {
	tests := []struct {
		name string
		item *org.Item
		err  string
	}{
		{
			name: "valid item",
			item: &org.Item{
				Ext: tax.Extensions{
					cfdi.ExtKeyProdServ: "12345678",
				},
			},
		},
		{
			name: "missing extension",
			item: &org.Item{},
			err:  "ext: (mx-cfdi-prod-serv: required.)",
		},
		{
			name: "empty extension",
			item: &org.Item{
				Ext: tax.Extensions{},
			},
			err: "ext: (mx-cfdi-prod-serv: required.)",
		},
		{
			name: "missing SAT identity",
			item: &org.Item{
				Ext: tax.Extensions{
					"random": "12345678",
				},
			},
			err: "ext: (mx-cfdi-prod-serv: required.).",
		},
		{
			name: "invalid code format",
			item: &org.Item{
				Ext: tax.Extensions{
					cfdi.ExtKeyProdServ: "AbC2",
				},
			},
			err: "ext: (mx-cfdi-prod-serv: must have 8 digits.)",
		},
	}

	addon := tax.AddonForKey(cfdi.KeyV4)
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := addon.Validate(ts.item)
			if ts.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), ts.err)
				}
			}
		})
	}
}

func TestItemIdentityNormalization(t *testing.T) {
	addon := tax.AddonForKey(cfdi.KeyV4)
	tests := []struct {
		Code     tax.ExtValue
		Expected tax.ExtValue
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
		err := addon.Normalize(item)
		assert.NoError(t, err)
		assert.Equal(t, ts.Expected, item.Ext[cfdi.ExtKeyProdServ])
	}

	// In context of invoice
	inv := validInvoice()
	inv.Lines[0].Item.Ext[cfdi.ExtKeyProdServ] = "010101"
	err := inv.Calculate()
	require.NoError(t, err)
	assert.Equal(t, tax.ExtValue("01010100"), inv.Lines[0].Item.Ext[cfdi.ExtKeyProdServ])
}

func TestInvoiceLineExtensions(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())

	l := inv.Lines[0]
	assert.Equal(t, "002", l.Taxes[0].Ext[cfdi.ExtKeyTaxType].String())
}

func TestItemIdentityMigration(t *testing.T) {
	inv := validInvoice()

	inv.Lines[0].Item.Ext = nil
	inv.Lines[0].Item.Identities = []*org.Identity{
		{
			Key:  cfdi.ExtKeyProdServ,
			Code: "01010101",
		},
	}

	err := inv.Calculate()
	require.NoError(t, err)
	assert.Equal(t, tax.ExtValue("01010101"), inv.Lines[0].Item.Ext[cfdi.ExtKeyProdServ])
}
