package facturae_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/facturae"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceDocumentScenarios(t *testing.T) {
	t.Run("no addons", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.Addons.List = []cbc.Key{} // reset
		require.NoError(t, i.Calculate())
		assert.Empty(t, i.Tax.Ext[facturae.ExtKeyDocType])
	})

	t.Run("with addon", func(t *testing.T) {
		i := testInvoiceStandard(t)
		require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 0)
		assert.Equal(t, i.Tax.Ext[facturae.ExtKeyDocType].String(), "FC")
	})

	t.Run("self-billed", func(t *testing.T) {
		i := testInvoiceStandard(t)
		i.SetTags(tax.TagSelfBilled)
		require.NoError(t, i.Calculate())
		assert.Len(t, i.Notes, 1)
		assert.Equal(t, i.Tax.Ext[facturae.ExtKeyDocType].String(), "AF")
	})
}
