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
		// Simplified invoices without customer tax id or identities should get F2
		i.Customer.TaxID = nil
		i.Customer.Identities = nil
		require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 0)
		assert.Equal(t, i.Tax.Ext[verifactu.ExtKeyDocType].String(), "F2")
	})

	t.Run("credit note", func(t *testing.T) {
		i := testInvoiceStandard(t)
		require.NoError(t, i.Correct(bill.Credit, bill.WithExtension(verifactu.ExtKeyDocType, "R1")))
		// require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 0)
		assert.Equal(t, i.Tax.Ext[verifactu.ExtKeyDocType].String(), "R1")
		assert.Equal(t, i.Tax.Ext.Get(verifactu.ExtKeyCorrectionType).String(), "I")
	})

	t.Run("corrective", func(t *testing.T) {
		i := testInvoiceStandard(t)
		require.NoError(t, i.Correct(bill.Corrective, bill.WithExtension(verifactu.ExtKeyDocType, "R2")))
		assert.Len(t, i.Notes, 0)
		assert.Equal(t, i.Tax.Ext.Get(verifactu.ExtKeyDocType).String(), "R2")
		assert.Equal(t, i.Tax.Ext.Get(verifactu.ExtKeyCorrectionType).String(), "S")
	})

	t.Run("corrective special", func(t *testing.T) {
		i := testInvoiceStandard(t)
		require.NoError(t, i.Calculate())
		require.NoError(t, i.Correct(bill.Corrective, bill.WithExtension(verifactu.ExtKeyDocType, "R2")))
		assert.Equal(t, "R2", i.Tax.Ext.Get(verifactu.ExtKeyDocType).String())
	})

	t.Run("simplified corrective invoice", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.SetTags(tax.TagSimplified)
		// Simplified invoices without customer tax id or identities should get R5 when corrected
		i.Customer.TaxID = nil
		i.Customer.Identities = nil
		require.NoError(t, i.Calculate())
		require.NoError(t, i.Correct(bill.Corrective))
		assert.Equal(t, i.Tax.Ext[verifactu.ExtKeyDocType].String(), "R5")
	})

	t.Run("replacement", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.SetTags(tax.TagReplacement)
		require.NoError(t, i.Calculate())
		assert.Equal(t, i.Tax.Ext[verifactu.ExtKeyDocType].String(), "F3")
	})
}
