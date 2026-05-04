package facturae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
)

var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// ** Invoice Document Types **
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
					bill.InvoiceTypeCorrective,
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyDocType: "FC", // default
				}),
			},
			{
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyDocType: "FA",
				}),
			},
			{
				Tags: []cbc.Key{
					tax.TagSelfBilled,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyDocType: "AF",
				}),
			},
			// ** Invoice Class **
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceClass: "OO", // Original Invoice
				}),
			},
			{
				Types: []cbc.Key{
					bill.InvoiceTypeCorrective,
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceClass: "OR", // Corrective
				}),
			},
			{
				Tags: []cbc.Key{es.TagSummary},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceClass: "OC", // Summary
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{es.TagCopy},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceClass: "CO", // Copy of the original
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeCorrective},
				Tags:  []cbc.Key{es.TagCopy},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceClass: "CR", // Copy of the corrective
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{es.TagCopy, es.TagSummary},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceClass: "CC", // Copy of the summary
				}),
			},
		},
	},
}
