package mx_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/mx"
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
				Ext: tax.ExtMap{
					mx.ExtKeyCFDIProdServ: "12345678",
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
				Ext: tax.ExtMap{},
			},
			err: "ext: (mx-cfdi-prod-serv: required.)",
		},
		{
			name: "missing SAT identity",
			item: &org.Item{
				Ext: tax.ExtMap{
					"random": "12345678",
				},
			},
			err: "ext: (mx-cfdi-prod-serv: required.).",
		},
		{
			name: "invalid code format",
			item: &org.Item{
				Ext: tax.ExtMap{
					mx.ExtKeyCFDIProdServ: "AbC2",
				},
			},
			err: "ext: (mx-cfdi-prod-serv: must have 8 digits.)",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := mx.Validate(ts.item)
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
	r := mx.New()
	tests := []struct {
		Code     cbc.KeyOrCode
		Expected cbc.KeyOrCode
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
		item := &org.Item{Ext: tax.ExtMap{mx.ExtKeyCFDIProdServ: ts.Code}}
		err := r.CalculateObject(item)
		assert.NoError(t, err)
		assert.Equal(t, ts.Expected, item.Ext[mx.ExtKeyCFDIProdServ])
	}

	// In context of invoice
	inv := validInvoice()
	inv.Lines[0].Item.Ext[mx.ExtKeyCFDIProdServ] = "010101"
	err := inv.Calculate()
	require.NoError(t, err)
	assert.Equal(t, cbc.KeyOrCode("01010100"), inv.Lines[0].Item.Ext[mx.ExtKeyCFDIProdServ])
}

func TestItemIdentityMigration(t *testing.T) {
	inv := validInvoice()

	inv.Lines[0].Item.Ext = nil
	inv.Lines[0].Item.Identities = []*org.Identity{
		{
			Key:  mx.ExtKeyCFDIProdServ,
			Code: "01010101",
		},
	}

	err := inv.Calculate()
	require.NoError(t, err)
	assert.Equal(t, cbc.KeyOrCode("01010101"), inv.Lines[0].Item.Ext[mx.ExtKeyCFDIProdServ])
}
