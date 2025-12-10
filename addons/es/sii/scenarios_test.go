package sii_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/sii"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceDocumentScenarios(t *testing.T) {
	t.Run("with addon", func(t *testing.T) {
		i := testInvoiceStandard(t)
		require.NoError(t, i.Calculate())
		assert.Equal(t, i.Tax.Ext[sii.ExtKeyDocType].String(), "F1")
	})

	t.Run("simplified invoice", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.SetTags(tax.TagSimplified)
		require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 0)
		assert.Equal(t, i.Tax.Ext[sii.ExtKeyDocType].String(), "F2")
	})

	t.Run("credit note", func(t *testing.T) {
		i := testInvoiceStandard(t)
		require.NoError(t, i.Correct(bill.Credit, bill.WithExtension(sii.ExtKeyDocType, "R1")))
		// require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 0)
		assert.Equal(t, i.Tax.Ext[sii.ExtKeyDocType].String(), "R1")
		assert.Equal(t, i.Tax.Ext.Get(sii.ExtKeyCorrectionType).String(), "I")
	})

	t.Run("corrective", func(t *testing.T) {
		i := testInvoiceStandard(t)
		require.NoError(t, i.Correct(bill.Corrective, bill.WithExtension(sii.ExtKeyDocType, "R2")))
		assert.Len(t, i.Notes, 0)
		assert.Equal(t, i.Tax.Ext.Get(sii.ExtKeyDocType).String(), "R2")
		assert.Equal(t, i.Tax.Ext.Get(sii.ExtKeyCorrectionType).String(), "S")
	})

	t.Run("corrective special", func(t *testing.T) {
		i := testInvoiceStandard(t)
		require.NoError(t, i.Calculate())
		require.NoError(t, i.Correct(bill.Corrective, bill.WithExtension(sii.ExtKeyDocType, "R2")))
		assert.Equal(t, "R2", i.Tax.Ext.Get(sii.ExtKeyDocType).String())
	})

	t.Run("simplified corrective invoice", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.SetTags(tax.TagSimplified)
		require.NoError(t, i.Calculate())
		require.NoError(t, i.Correct(bill.Corrective))
		assert.Equal(t, i.Tax.Ext[sii.ExtKeyDocType].String(), "R5")
	})

	t.Run("replacement", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.SetTags(tax.TagReplacement)
		require.NoError(t, i.Calculate())
		assert.Equal(t, i.Tax.Ext[sii.ExtKeyDocType].String(), "F3")
	})
}
