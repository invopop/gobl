package arca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// Document tag keys for Argentine invoice types
const (
	// TagVATRegistered is used for Invoice A - transactions where the customer
	// is VAT registered (Responsable Inscripto or Monotributista).
	TagVATRegistered cbc.Key = "vat-registered"

	// TagSimplifiedRegime is used for Invoice C - when the supplier is under a
	// simplified tax regime (Monotributo in Argentina).
	TagSimplifiedRegime cbc.Key = "simplified-regime"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagVATRegistered,
			Name: i18n.String{
				i18n.EN: "VAT Registered Customer",
				i18n.ES: "Cliente Registrado en IVA",
			},
			Desc: i18n.String{
				i18n.EN: "Invoice A: Customer is VAT registered (Responsable Inscripto or Monotributista).",
				i18n.ES: "Factura A: El cliente está registrado en IVA (Responsable Inscripto o Monotributista).",
			},
		},
		{
			Key: TagSimplifiedRegime,
			Name: i18n.String{
				i18n.EN: "Simplified Tax Regime",
				i18n.ES: "Régimen Simplificado",
			},
			Desc: i18n.String{
				i18n.EN: "Invoice C: Supplier is under a simplified tax regime (Monotributo).",
				i18n.ES: "Factura C: El proveedor está bajo un régimen tributario simplificado (Monotributo).",
			},
		},
	},
}

var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// ** Invoice A - Customer is VAT registered **
			{
				Name: i18n.String{
					i18n.EN: "Invoice A",
					i18n.ES: "Factura A",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used when the invoice is issued by a VAT registered company to another
						VAT registered company or a monotributista.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					TagVATRegistered,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "1",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Debit Note A",
					i18n.ES: "Nota de Débito A",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used when the debit note is issued by a VAT registered company to another
						VAT registered company or a monotributista.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeDebitNote,
				},
				Tags: []cbc.Key{
					TagVATRegistered,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "2",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Credit Note A",
					i18n.ES: "Nota de Crédito A",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used when the credit note is issued by a VAT registered company to another
						VAT registered company or a monotributista.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
				Tags: []cbc.Key{
					TagVATRegistered,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "3",
				},
			},
			// ** Invoice B - Final consumers and VAT exempt customers **
			{
				Name: i18n.String{
					i18n.EN: "Invoice B",
					i18n.ES: "Factura B",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used when the invoice is issued by a VAT registered company to final
						consumers, exempt subjects, non-categorized subjects, or foreign customers.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "6",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Debit Note B",
					i18n.ES: "Nota de Débito B",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used when the debit note is issued by a VAT registered company to final
						consumers, exempt subjects, non-categorized subjects, or foreign customers.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeDebitNote,
				},
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "7",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Credit Note B",
					i18n.ES: "Nota de Crédito B",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used when the credit note is issued by a VAT registered company to final
						consumers, exempt subjects, non-categorized subjects, or foreign customers.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "8",
				},
			},
			// ** Invoice C - Monotributista transactions **
			{
				Name: i18n.String{
					i18n.EN: "Invoice C",
					i18n.ES: "Factura C",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used when the invoice is issued by a monotributista (simplified tax regime)
						to any type of customer.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					TagSimplifiedRegime,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "11",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Debit Note C",
					i18n.ES: "Nota de Débito C",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used when the debit note is issued by a monotributista (simplified tax regime)
						to any type of customer.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeDebitNote,
				},
				Tags: []cbc.Key{
					TagSimplifiedRegime,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "12",
				},
			},
			{
				Name: i18n.String{
					i18n.EN: "Credit Note C",
					i18n.ES: "Nota de Crédito C",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Used when the credit note is issued by a monotributista (simplified tax regime)
						to any type of customer.
					`),
				},
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
				Tags: []cbc.Key{
					TagSimplifiedRegime,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "13",
				},
			},
		},
	},
}
