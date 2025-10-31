package ar_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ar"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarios(t *testing.T) {
	t.Run("should have invoice tags", func(t *testing.T) {
		r := ar.New()
		require.NotNil(t, r)
		require.NotNil(t, r.Tags)
		require.Len(t, r.Tags, 1)

		invoiceTags := r.Tags[0]
		assert.Equal(t, bill.ShortSchemaInvoice, invoiceTags.Schema)
		assert.Greater(t, len(invoiceTags.List), 0)
	})

	t.Run("should have invoice scenarios", func(t *testing.T) {
		r := ar.New()
		require.NotNil(t, r)
		require.NotNil(t, r.Scenarios)
		require.Len(t, r.Scenarios, 1)

		invoiceScenarios := r.Scenarios[0]
		assert.Equal(t, bill.ShortSchemaInvoice, invoiceScenarios.Schema)
		assert.Greater(t, len(invoiceScenarios.List), 0)
	})

	t.Run("should add legal note for invoice type A", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(ar.TagInvoiceTypeA),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
		assert.Contains(t, inv.Notes[0].Text, "Factura Tipo A")
	})

	t.Run("should add legal note for invoice type B", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(ar.TagInvoiceTypeB),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
		assert.Contains(t, inv.Notes[0].Text, "Factura Tipo B")
	})

	t.Run("should add legal note for invoice type C", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(ar.TagInvoiceTypeC),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
		assert.Contains(t, inv.Notes[0].Text, "Factura Tipo C")
	})

	t.Run("should add legal note for invoice type E", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(ar.TagInvoiceTypeE),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
		assert.Contains(t, inv.Notes[0].Text, "Factura Tipo E")
	})

	t.Run("should add legal note for monotributo", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(ar.TagMonotributo),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
		assert.Contains(t, inv.Notes[0].Text, "Monotributo")
	})

	t.Run("should add legal note for export services", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(ar.TagExportServices),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
		assert.Contains(t, inv.Notes[0].Text, "Exportación de Servicios")
	})

	t.Run("should add legal note for export goods", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(ar.TagExportGoods),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
		assert.Contains(t, inv.Notes[0].Text, "Exportación de Bienes")
	})

	t.Run("should add legal note for VAT exempt", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(ar.TagExento),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
		assert.Contains(t, inv.Notes[0].Text, "IVA Exento")
	})

	t.Run("should add legal note for reverse charge", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(tax.TagReverseCharge),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
		assert.Contains(t, inv.Notes[0].Text, "Inversión del Sujeto Pasivo")
	})

	t.Run("should add legal note for simplified invoice", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(tax.TagSimplified),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
		assert.Contains(t, inv.Notes[0].Text, "Factura Simplificada")
	})

	t.Run("should handle multiple tags with multiple notes", func(t *testing.T) {
		inv := &bill.Invoice{
			Tags: tax.WithTags(ar.TagInvoiceTypeA, ar.TagResponsableInscripto),
		}

		err := inv.Calculate()
		require.NoError(t, err)

		// Only invoice type A should add a note (responsable-inscripto doesn't have a scenario)
		require.NotNil(t, inv.Notes)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
	})
}
