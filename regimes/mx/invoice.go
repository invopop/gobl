package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
)

const (
	// constants copied from CFDI
	extKeyIssuePlace = "mx-cfdi-issue-place"
	extKeyPostCode   = "mx-cfdi-post-code"
)

func normalizeInvoice(inv *bill.Invoice) {
	// 2024-04-26: copy suppliers post code to invoice, if not already
	// set.
	var it bill.Tax
	if inv.Tax != nil {
		it = *inv.Tax
	} else {
		it = bill.Tax{}
	}
	if it.Ext == nil {
		it.Ext = make(tax.Extensions)
	}
	if it.Ext.Has(extKeyIssuePlace) {
		return
	}
	if inv.Supplier != nil && inv.Supplier.Ext.Has(extKeyPostCode) {
		it.Ext[extKeyIssuePlace] = inv.Supplier.Ext[extKeyPostCode]
		inv.Tax = &it
		return
	}
	if len(inv.Supplier.Addresses) > 0 {
		addr := inv.Supplier.Addresses[0]
		if addr.Code != "" {
			it.Ext[extKeyIssuePlace] = tax.ExtValue(addr.Code)
			inv.Tax = &it
		}
	}
}
