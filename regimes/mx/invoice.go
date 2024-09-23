package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

const (
	// constants copied from CFDI
	extKeyIssuePlace cbc.Key = "mx-cfdi-issue-place"
	extKeyPostCode   cbc.Key = "mx-cfdi-post-code"
)

func normalizeInvoice(inv *bill.Invoice) {
	// 2024-04-26: copy suppliers post code to invoice, if not already
	// set.
	normalizeParty(inv.Supplier) // first do party
	ext := make(tax.Extensions)
	if inv.Tax != nil {
		ext = inv.Tax.Ext
	}
	if ext.Has(extKeyIssuePlace) {
		return
	}
	if inv.Supplier != nil && inv.Supplier.Ext.Has(extKeyPostCode) {
		ext[extKeyIssuePlace] = inv.Supplier.Ext[extKeyPostCode]
	} else if len(inv.Supplier.Addresses) > 0 {
		addr := inv.Supplier.Addresses[0]
		if addr.Code != "" {
			ext[extKeyIssuePlace] = tax.ExtValue(addr.Code)
		}
	}
	if len(ext) > 0 {
		if inv.Tax == nil {
			inv.Tax = new(bill.Tax)
		}
		inv.Tax.Ext = ext
	}
}
