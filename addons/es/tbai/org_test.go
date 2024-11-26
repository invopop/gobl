package tbai_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeOrgItem(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ad := tax.AddonForKey(tbai.V1)
		var item *org.Item
		assert.NotPanics(t, func() {
			ad.Normalizer(item)
		})
	})
	t.Run("with standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "services", inv.Lines[0].Item.Ext[tbai.ExtKeyProduct].String())
	})
}
