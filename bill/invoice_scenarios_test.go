package bill_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarios(t *testing.T) {
	t.Run("basic flow", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		require.NoError(t, inv.Calculate())
		ss := inv.ScenarioSummary() //nolint:staticcheck
		require.NotNil(t, ss)
	})

	t.Run("invalid tags", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.SetTags("random")
		assert.ErrorContains(t, inv.Calculate(), "'random' undefined")
	})

	t.Run("scenario for new note", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.SetTags(tax.TagReverseCharge)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "12345678903",
		}
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Tax)
		assert.Len(t, inv.Tax.Notes, 1)
		assert.Equal(t, "Reverse Charge / Inversione del soggetto passivo", inv.Tax.Notes[0].Text)
	})

	t.Run("scenario for existing note", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.SetTags(tax.TagReverseCharge)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "12345678903",
		}
		inv.Tax = &bill.Tax{
			Notes: []*tax.Note{
				{
					Category: tax.CategoryVAT,
					Key:      tax.KeyReverseCharge,
					Text:     "User provided note",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Tax)
		assert.Len(t, inv.Tax.Notes, 1)
		assert.Equal(t, "User provided note", inv.Tax.Notes[0].Text, "should keep user's existing note")
	})

	t.Run("without tax defined", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Addons = tax.WithAddons(sdi.V1)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "12345678903",
		}
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		assert.Equal(t, 2, inv.Tax.Ext.Len())
		assert.Equal(t, "TD01", inv.Tax.Ext.Get(sdi.ExtKeyDocumentType).String())
	})

	t.Run("overwrite previous values with tag", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Addons = tax.WithAddons(sdi.V1)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "12345678903",
		}
		inv.SetTags(tax.TagB2G)
		inv.Tax = &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				sdi.ExtKeyFormat: "XXXX",
			}),
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, 2, inv.Tax.Ext.Len())
		assert.Equal(t, "FPA12", inv.Tax.Ext.Get(sdi.ExtKeyFormat).String())
	})

	t.Run("overwrite previous values without tags", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Addons = tax.WithAddons(sdi.V1)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "12345678903",
		}
		inv.Tax = &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				sdi.ExtKeyFormat: "XXXX",
			}),
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, 2, inv.Tax.Ext.Len())
		assert.Equal(t, "FPR12", inv.Tax.Ext.Get(sdi.ExtKeyFormat).String())
	})
}

func TestInvoiceGetExtensions(t *testing.T) {
	t.Run("with lines", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Addons = tax.WithAddons(sdi.V1)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "12345678903",
		}
		require.NoError(t, inv.Calculate())
		ext := inv.GetExtensions()
		assert.Len(t, ext, 2)
		assert.Equal(t, "FPR12", ext[0].Get(sdi.ExtKeyFormat).String())
	})
	t.Run("missing lines", func(t *testing.T) {
		inv := baseInvoice(t)
		inv.Addons = tax.WithAddons(saft.V1)
		require.NoError(t, inv.Calculate())
		ext := inv.GetExtensions()
		assert.Len(t, ext, 1)
	})
}

func TestInvoiceScenarios(t *testing.T) {
	t.Run("provides list", func(t *testing.T) {
		list := bill.InvoiceScenarios()
		assert.Len(t, list.List, 1)
	})
}
