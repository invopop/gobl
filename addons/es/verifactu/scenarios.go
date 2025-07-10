package verifactu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// ** Invoice Document Types **
			{
				Name: i18n.String{
					i18n.EN: "Standard Invoice",
					i18n.ES: "Factura Estándar",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Standard invoice used for B2B transactions, where the complete fiscal details of the customer
						are available.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "F1",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Simplified Invoice",
					i18n.ES: "Factura Simplificada",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used for B2C transactions when the client details are not available.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "F2",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Replacement Invoice",
					i18n.ES: "Factura Emitida en Sustitución",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used under special circumstances to indicate that this invoice replaces a previously
						issued simplified invoice. The previous document was correct, but the replacement is
						necessary to provide tax details of the customer.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagReplacement,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "F3",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Simplified Corrective Invoice",
					i18n.ES: "Factura Simplificada Correctiva",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						This scenario covers when a simplified invoice is being corrected either
						with a credit or debit note, or a corrective replacement invoice.

						In VERI*FACTU, only the document type ~R5~ is supported for corrective
						invoices.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeCorrective,
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "R5",
				},
			},
		},
	},
}
