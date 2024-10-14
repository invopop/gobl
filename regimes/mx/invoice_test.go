package mx_test

import (
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeInvoice(t *testing.T) {
	t.Run("no tax", func(t *testing.T) {
		inv := baseInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.ExtValue("21000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
	})

	t.Run("no ext", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax = &bill.Tax{}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.ExtValue("21000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
	})

	t.Run("with supplier address code", func(t *testing.T) {
		inv := baseInvoice()
		delete(inv.Supplier.Ext, cfdi.ExtKeyPostCode)
		inv.Supplier.Addresses = append(inv.Supplier.Addresses,
			&org.Address{
				Locality: "Mexico",
				Code:     "21000",
			},
		)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.ExtValue("21000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
	})
	t.Run("migrate supplier issue place", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax = nil
		inv.Supplier.Ext = tax.Extensions{
			cfdi.ExtKeyPostCode: "12345",
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.ExtValue("12345"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
	})
	t.Run("migrate supplier issue place", func(t *testing.T) {
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
		assert.Equal(t, tax.ExtValue("12345"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
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
				cfdi.ExtKeyPostCode:     "21000",
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
				cfdi.ExtKeyPostCode:     "65000",
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
