package tbai_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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
		assert.Equal(t, "BI", inv.Tax.Ext[tbai.ExtKeyRegion].String())
	})

	t.Run("standard invoice in Gipuzkoa", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Gipuzkoa",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "SS", inv.Tax.Ext[tbai.ExtKeyRegion].String())
	})

	t.Run("standard invoice in Álava (accent)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Álava",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "VI", inv.Tax.Ext[tbai.ExtKeyRegion].String())
	})

	t.Run("standard invoice in Araba", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Araba",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "VI", inv.Tax.Ext[tbai.ExtKeyRegion].String())
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
			Ext: tax.Extensions{
				tbai.ExtKeyRegion: "BI", // not Alaba
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "BI", inv.Tax.Ext[tbai.ExtKeyRegion].String())
	})
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("with services", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0].Ext[tbai.ExtKeyProduct] = "services"
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing customer tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (tax_id: cannot be blank.)")
	})

	t.Run("with exemption reason", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0].Ext = nil
		assertValidationError(t, inv, "es-tbai-exemption: required")
	})

	t.Run("without series", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Series = ""
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("without notes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Notes = nil
		assertValidationError(t, inv, "notes: with key 'general' missing")
	})

	t.Run("correction", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Correct(
			bill.Credit,
			bill.WithExtension(tbai.ExtKeyCorrection, "R4"),
		))
		assert.Len(t, inv.Preceding, 1)
		assert.NoError(t, inv.Validate())
	})
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	t.Helper()
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.ErrorContains(t, err, expected)
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Addons: tax.WithAddons(tbai.V1),
		Series: "ABC",
		Code:   "123",
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				tbai.ExtKeyRegion: "BI",
			},
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
						Rate:     "exempt",
						Ext: tax.Extensions{
							tbai.ExtKeyExemption: "E1",
						},
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
