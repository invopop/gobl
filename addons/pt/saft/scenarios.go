package saft

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
)

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Ext: tax.ExtensionsOf(tax.ExtMap{
				ExtKeyInvoiceType: InvoiceTypeStandard,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{tax.TagSimplified},
			Ext: tax.ExtensionsOf(tax.ExtMap{
				ExtKeyInvoiceType: InvoiceTypeSimplified,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Filter: func(doc any) bool {
				inv, ok := doc.(*bill.Invoice)
				if !ok {
					return false
				}
				return inv.HasTags(pt.TagInvoiceReceipt) || inv.Totals.Paid()
			},
			Ext: tax.ExtensionsOf(tax.ExtMap{
				ExtKeyInvoiceType: InvoiceTypeInvoiceReceipt,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Ext: tax.ExtensionsOf(tax.ExtMap{
				ExtKeyInvoiceType: InvoiceTypeDebitNote,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Ext: tax.ExtensionsOf(tax.ExtMap{
				ExtKeyInvoiceType: InvoiceTypeCreditNote,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeProforma},
			Ext: tax.ExtensionsOf(tax.ExtMap{
				ExtKeyWorkType: WorkTypeProforma,
			}),
		},
	},
}
