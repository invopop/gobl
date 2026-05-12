package zatca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// scenarios only maps the scenarios needed for
// ZATCA rules validation.
var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// ZATCA standard invoice overrides en16931 untdid document type.
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType:     "388",
					ExtKeyInvoiceTypeTransactions: "0100000",
				}),
			},
			// Default credit/debit notes transaction type
			{
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceTypeTransactions: "0100000",
				}),
			},
			// Simplified invoices and associated credit/debit notes.
			{
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceTypeTransactions: "0200000",
				}),
			},
			// Export invoices and associated credit/debit notes.
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
				Tags: []cbc.Key{
					tax.TagExport,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceTypeTransactions: "0100100",
				}),
			},
			// Summary and associated credit/debit notes.
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
				Tags: []cbc.Key{
					TagSummary,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceTypeTransactions: "0100010",
				}),
			},
			{
				Tags: []cbc.Key{
					TagSummary,
					tax.TagSimplified,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceTypeTransactions: "0200010",
				}),
			},
			// Simplified and summary and associated credit/debit notes.
			{
				Tags: []cbc.Key{
					tax.TagSimplified,
					TagSummary,
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyInvoiceTypeTransactions: "0200010",
				}),
			},
		},
	},
}
