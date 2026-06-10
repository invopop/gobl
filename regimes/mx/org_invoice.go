package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

const (
	// constants copied from CFDI
	extKeyIssuePlace cbc.Key = "mx-cfdi-issue-place"
	extKeyPostCode   cbc.Key = "mx-cfdi-post-code"
)

func normalizeInvoice(inv *bill.Invoice) {
	migrateSupplierExtPostCodeToInvoice(inv)
	migrateCustomerExtPostCodeToAddress(inv)

	deleteExtPostCode(inv.Supplier)
	deleteExtPostCode(inv.Customer)
}

// 2024-04-26: copy suppliers post code to invoice, if not alread set.
func migrateSupplierExtPostCodeToInvoice(inv *bill.Invoice) {
	ext := tax.MakeExtensions()
	if inv.Tax != nil && !inv.Tax.Ext.IsZero() {
		ext = inv.Tax.Ext
	}
	if ext.Has(extKeyIssuePlace) || inv.Supplier == nil {
		return
	}
	if inv.Supplier.Ext.Has(extKeyPostCode) {
		ext = ext.Set(extKeyIssuePlace, inv.Supplier.Ext.Get(extKeyPostCode))
	} else if len(inv.Supplier.Addresses) > 0 {
		addr := inv.Supplier.Addresses[0]
		if addr.Code != "" {
			ext = ext.Set(extKeyIssuePlace, addr.Code)
		}
	}
	if ext.Len() > 0 {
		if inv.Tax == nil {
			inv.Tax = new(bill.Tax)
		}
		inv.Tax.Ext = ext
	}
}

// 2024-10-29: move customers post code from extension to address
func migrateCustomerExtPostCodeToAddress(inv *bill.Invoice) {
	if inv.Customer != nil && inv.Customer.Ext.Has(extKeyPostCode) {
		if len(inv.Customer.Addresses) == 0 {
			inv.Customer.Addresses = []*org.Address{{}}
		}
		inv.Customer.Addresses[0].Code = inv.Customer.Ext.Get(extKeyPostCode)
	}
}

// 2024-10-29: remove post codes from extensions (no longer valid extension)
func deleteExtPostCode(p *org.Party) {
	if p == nil {
		return
	}

	p.Ext = p.Ext.Delete(extKeyPostCode)
}
