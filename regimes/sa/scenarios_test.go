package sa_test

import (
	"testing"

	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceScenarios(t *testing.T) {
	t.Run("reverse charge note", func(t *testing.T) {
		i := testStandardInvoice(t)
		i.Tags = tax.WithTags(tax.TagReverseCharge)
		require.NoError(t, i.Calculate())
		require.NoError(t, i.Validate())
		require.Len(t, i.Notes, 1)
		assert.Equal(t, tax.TagReverseCharge, i.Notes[0].Src)
		assert.Contains(t, i.Notes[0].Text, "Reverse Charge")
		assert.Contains(t, i.Notes[0].Text, "آلية الاحتساب العكسي")
	})

	t.Run("simplified invoice note", func(t *testing.T) {
		i := testSimplifiedInvoice(t)
		require.NoError(t, i.Calculate())
		require.NoError(t, i.Validate())
		require.Len(t, i.Notes, 1)
		assert.Equal(t, tax.TagSimplified, i.Notes[0].Src)
		assert.Contains(t, i.Notes[0].Text, "Simplified Tax Invoice")
		assert.Contains(t, i.Notes[0].Text, "فاتورة ضريبية مبسطة")
	})

	t.Run("export note", func(t *testing.T) {
		i := testStandardInvoice(t)
		i.Tags = tax.WithTags(tax.TagExport)
		require.NoError(t, i.Calculate())
		require.NoError(t, i.Validate())
		require.Len(t, i.Notes, 1)
		assert.Equal(t, tax.TagExport, i.Notes[0].Src)
		assert.Contains(t, i.Notes[0].Text, "zero-rated")
		assert.Contains(t, i.Notes[0].Text, "تصدير سلع أو خدمات")
	})
}
