package adecf_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/adecf"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeOrgItem(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ad := tax.AddonForKey(adecf.V1)
		var item *org.Item
		assert.NotPanics(t, func() {
			ad.Normalizer(item)
		})
	})
	t.Run("with standard invoice", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "services", inv.Lines[0].Item.Ext[adecf.ExtKeyProduct].String())
	})
}
