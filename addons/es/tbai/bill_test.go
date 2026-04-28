package tbai_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceNormalization(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ad := tax.AddonForKey(tbai.V1)
		var inv *bill.Invoice
		assert.NotPanics(t, func() {
			ad.Normalizer(inv)
		})
	})

	t.Run("standard invoice, no address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		assert.Nil(t, inv.Tax)
	})

	t.Run("standard invoice in Vizcaya", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Vizcaya",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "BI", inv.Tax.Ext.Get(tbai.ExtKeyRegion).String())
	})

	t.Run("standard invoice in Gipuzkoa", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Gipuzkoa",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "SS", inv.Tax.Ext.Get(tbai.ExtKeyRegion).String())
	})

	t.Run("standard invoice in Álava (accent)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Álava",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "VI", inv.Tax.Ext.Get(tbai.ExtKeyRegion).String())
	})

	t.Run("standard invoice in Araba", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Araba",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "VI", inv.Tax.Ext.Get(tbai.ExtKeyRegion).String())
	})

	t.Run("standard invoice in Araba", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Madrid",
		})
		require.NoError(t, inv.Calculate())
		assert.Nil(t, inv.Tax)
	})

	t.Run("with existing region", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Araba",
		})
		inv.Tax = &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				tbai.ExtKeyRegion: "BI", // not Alaba
			}),
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "BI", inv.Tax.Ext.Get(tbai.ExtKeyRegion).String())
	})
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("with services", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0].Ext = inv.Lines[0].Taxes[0].Ext.Set(tbai.ExtKeyProduct, "services")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer is required for non-simplified invoices")
	})

	t.Run("missing customer tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer tax ID is required")
	})

	t.Run("simplified invoice without customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified invoice with customer without tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified invoice with customer tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer tax ID must not be set for simplified invoices")
	})

	t.Run("with exemption reason", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0].Ext = tax.Extensions{}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("without series", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Series = ""
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("without notes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Notes = nil
		assertValidationError(t, inv, "with key 'general' missing")
	})

	t.Run("correction", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Correct(
			bill.Credit,
			bill.WithExtension(tbai.ExtKeyCorrection, "R4"),
		))
		assert.Len(t, inv.Preceding, 1)
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestBillLineNormalization(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ad := tax.AddonForKey(tbai.V1)
		var line *bill.Line
		assert.NotPanics(t, func() {
			ad.Normalizer(line)
		})
	})
	t.Run("with standard invoice, set default", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "services", inv.Lines[0].Taxes[0].Ext.Get(tbai.ExtKeyProduct).String())
	})
	t.Run("with standard invoice, set override for goods", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Item.Key = org.ItemKeyGoods
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "goods", inv.Lines[0].Taxes[0].Ext.Get(tbai.ExtKeyProduct).String())
	})
	t.Run("with standard invoice, set override for resale", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Item.Key = org.ItemKeyGoods
		inv.Lines[0].Taxes[0].Ext = inv.Lines[0].Taxes[0].Ext.Set(tbai.ExtKeyProduct, "resale")
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "resale", inv.Lines[0].Taxes[0].Ext.Get(tbai.ExtKeyProduct).String())
	})
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	t.Helper()
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	require.ErrorContains(t, err, expected)
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Addons: tax.WithAddons(tbai.V1),
		Series: "ABC",
		Code:   "123",
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				tbai.ExtKeyRegion: "BI",
			}),
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "NL",
				Code:    "000099995B57",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Key:      "exempt",
						Ext: tax.ExtensionsOf(tax.ExtMap{
							tbai.ExtKeyExempt: "E1",
						}),
					},
				},
			},
		},
		Notes: []*org.Note{
			{
				Key:  org.NoteKeyGeneral,
				Text: "This is a test invoice",
			},
		},
	}
}
