package ticket_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/ticket"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeOrgItem(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var item *org.Item
		assert.NotPanics(t, func() {
			norm.Normalize(item, tax.AddonContext(ticket.V1))
		})
	})
	t.Run("with standard invoice", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "services", inv.Lines[0].Item.Ext.Get(ticket.ExtKeyProduct).String())
	})
}
