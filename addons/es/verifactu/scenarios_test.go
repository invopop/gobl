package verifactu_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceDocumentScenarios(t *testing.T) {

	t.Run("with addon", func(t *testing.T) {
		i := testInvoiceStandard(t)
		require.NoError(t, i.Calculate())
		assert.Equal(t, i.Tax.Ext[verifactu.ExtKeyDocType].String(), "F1")
	})

	t.Run("simplified invoice", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.SetTags(tax.TagSimplified)
		require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 1)
		assert.Equal(t, i.Tax.Ext[verifactu.ExtKeyDocType].String(), "F2")
	})

	t.Run("substitution invoice", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.SetTags(verifactu.TagSubstitution)
		require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 1)
		assert.Equal(t, i.Tax.Ext[verifactu.ExtKeyDocType].String(), "F3")
	})

	t.Run("credit note", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.Type = bill.InvoiceTypeCreditNote
		require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 1)
		assert.Equal(t, i.Tax.Ext[verifactu.ExtKeyDocType].String(), "R1")
	})

	t.Run("corrective", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.Type = bill.InvoiceTypeCorrective
		require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 1)
		assert.Equal(t, i.Tax.Ext[verifactu.ExtKeyDocType].String(), "R1")
	})

	t.Run("simplified credit note", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.Type = bill.InvoiceTypeCreditNote
		i.SetTags(tax.TagSimplified)
		require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 1)
		assert.Equal(t, i.Tax.Ext[verifactu.ExtKeyDocType].String(), "R5")
	})
}
