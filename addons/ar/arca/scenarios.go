package arca

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
					i18n.EN: "Standard Invoice - A",
					i18n.ES: "Factura Estándar - A",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used for B2B transactions when the invoice is issued by a VAT registered company.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "001",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Standard Debit Note - A",
					i18n.ES: "Nota de Débito Estándar - A",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used for B2B transactions when the debit note is issued by a VAT registered company.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeDebitNote,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "002",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Standard Credit Note - A",
					i18n.ES: "Nota de Crédito Estándar - A",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used for B2B transactions when the credit note is issued by a VAT registered company.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "003",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Simplified Invoice - B",
					i18n.ES: "Factura Simplificada - B",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used for B2C transactions when the invoice is issued by a VAT registered company.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "006",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Simplified Debit Note - B",
					i18n.ES: "Nota de Débito Simplificada - B",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used for B2C transactions when the debit note is issued by a VAT registered company.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeDebitNote,
				},
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "007",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Simplified Credit Note - B",
					i18n.ES: "Nota de Crédito Simplificada - B",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used for B2C transactions when the credit note is issued by a VAT registered company.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "008",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Export Invoice",
					i18n.ES: "Factura de Exportación",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used for B2B transactions when the invoice is issued by a registered taxpayer or a monotributista that exports goods or services to clients outside the country.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagExport,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "019",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Debit Note for Foreign Operations",
					i18n.ES: "Nota de Débito por Operaciones con el Exterior",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used for B2B transactions when the debit note is issued by a company or business that exports goods or services to clients outside the country.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeDebitNote,
				},
				Tags: []cbc.Key{
					tax.TagExport,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "020",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Credit Note for Foreign Operations",
					i18n.ES: "Nota de Crédito por Operaciones con el Exterior",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used for B2B transactions when the credit note is issued by a company or business that exports goods or services to clients outside the country.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
				Tags: []cbc.Key{
					tax.TagExport,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "021",
				},
			},
		},
	},
}
