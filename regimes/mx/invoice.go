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
	migrateTaxIDZoneToExtPostCode(inv.Supplier)
	migrateTaxIDZoneToExtPostCode(inv.Customer)

	migrateSupplierExtPostCodeToInvoice(inv)
	migrateCustomerExtPostCodeToAddress(inv)

	deleteExtPostCode(inv.Supplier)
	deleteExtPostCode(inv.Customer)
}

// 2024-03-14: Migrate Tax ID Zone to extensions "mx-cfdi-post-code"
func migrateTaxIDZoneToExtPostCode(p *org.Party) {
	if p == nil {
		return
	}
	if p.TaxID != nil && p.TaxID.Zone != "" { //nolint:staticcheck
		if p.Ext == nil {
			p.Ext = make(tax.Extensions)
		}
		p.Ext[extKeyPostCode] = tax.ExtValue(p.TaxID.Zone) //nolint:staticcheck
		p.TaxID.Zone = ""                                  //nolint:staticcheck
	}
}

// 2024-04-26: copy suppliers post code to invoice, if not alread set.
func migrateSupplierExtPostCodeToInvoice(inv *bill.Invoice) {
	ext := make(tax.Extensions)
	if inv.Tax != nil && inv.Tax.Ext != nil {
		ext = inv.Tax.Ext
	}
	if ext.Has(extKeyIssuePlace) || inv.Supplier == nil {
		return
	}
	if inv.Supplier.Ext.Has(extKeyPostCode) {
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

// 2024-10-29: move customers post code from extension to address
func migrateCustomerExtPostCodeToAddress(inv *bill.Invoice) {
	if inv.Customer != nil && inv.Customer.Ext.Has(extKeyPostCode) {
		if len(inv.Customer.Addresses) == 0 {
			inv.Customer.Addresses = []*org.Address{{}}
		}
		inv.Customer.Addresses[0].Code = inv.Customer.Ext[extKeyPostCode].Code()
	}
}

// 2024-10-29: remove post codes from extensions (no longer valid extension)
func deleteExtPostCode(p *org.Party) {
	if p == nil {
		return
	}

	delete(p.Ext, extKeyPostCode)
}
