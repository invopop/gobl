package ticket

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func normalizeOrgItem(item *org.Item) {
	if item == nil {
		return
	}
	if item.Ext.IsZero() {
		item.Ext = tax.MakeExtensions()
	}
	if !item.Ext.Has(ExtKeyProduct) {
		// Assume all items are services by default.
		item.Ext = item.Ext.Set(ExtKeyProduct, "services")
	}
}
