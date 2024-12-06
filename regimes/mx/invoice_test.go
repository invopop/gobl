package mx_test

import (
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeInvoice(t *testing.T) {
	t.Run("does not migrate issue place when already present", func(t *testing.T) {
		inv := baseInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, cbc.Code("21000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
		assert.False(t, inv.Supplier.Ext.Has("mx-cfdi-post-code"))
	})
	t.Run("migrate issue place from supplier ext when no tax ext", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax = &bill.Tax{}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, cbc.Code("22000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
		assert.False(t, inv.Supplier.Ext.Has("mx-cfdi-post-code"))
	})
	t.Run("migrate issue place from supplier tax ID zone", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier.Ext = nil
		inv.Supplier.TaxID.Zone = "21000" //nolint:staticcheck
		inv.Tax = &bill.Tax{}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, cbc.Code("21000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
		assert.False(t, inv.Supplier.Ext.Has("mx-cfdi-post-code"))
	})
	t.Run("does not migrate issue place from address when already present", func(t *testing.T) {
		inv := baseInvoice()
		delete(inv.Supplier.Ext, "mx-cfdi-post-code")
		inv.Supplier.Addresses = append(inv.Supplier.Addresses,
			&org.Address{
				Locality: "Mexico",
				Code:     "22000",
			},
		)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, cbc.Code("21000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
	})
	t.Run("migrate issue place from supplier address code", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax = nil
		inv.Supplier.Ext = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses,
			&org.Address{
				Locality: "Mexico",
				Code:     "12345",
			},
		)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, cbc.Code("12345"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
	})
	t.Run("migrate customer post code from ext", func(t *testing.T) {
		inv := baseInvoice()
		inv.Customer.Ext = tax.Extensions{
			"mx-cfdi-post-code": "12345",
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Customer.Addresses)
		assert.Equal(t, "12345", inv.Customer.Addresses[0].Code.String())
		assert.False(t, inv.Customer.Ext.Has("mx-cfdi-post-code"))
	})
	t.Run("migrate customer post code from zone", func(t *testing.T) {
		inv := baseInvoice()
		inv.Customer.TaxID.Zone = "12345" //nolint:staticcheck
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Customer.Addresses)
		assert.Equal(t, "12345", inv.Customer.Addresses[0].Code.String())
		assert.False(t, inv.Customer.Ext.Has("mx-cfdi-post-code"))
	})
	t.Run("does not migrate anything when the customer is missing", func(t *testing.T) {
		inv := baseInvoice()
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Nil(t, inv.Customer)
	})
}

func baseInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime("MX"),
		Code:      "123",
		Currency:  "MXN",
		IssueDate: cal.MakeDate(2023, 1, 1),
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				cfdi.ExtKeyIssuePlace: "21000",
			},
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			Ext: tax.Extensions{
				"mx-cfdi-post-code":     "22000",
				cfdi.ExtKeyFiscalRegime: "601",
			},
			TaxID: &tax.Identity{
				Country: "MX",
				Code:    "AAA010101AAA",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			Ext: tax.Extensions{
				"mx-cfdi-post-code":     "65000",
				cfdi.ExtKeyFiscalRegime: "608",
				cfdi.ExtKeyUse:          "G01",
			},
			TaxID: &tax.Identity{
				Country: "MX",
				Code:    "ZZZ010101ZZZ",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.MakeAmount(10000, 2),
					Unit:  org.UnitPackage,
					Ext: tax.Extensions{
						cfdi.ExtKeyProdServ: "01010101",
					},
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
	}
}
