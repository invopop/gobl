package arca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// Document tag keys for Argentine invoice types
const (
	// TagSimplifiedRegime is used for Invoice C - when the supplier is under a
	// simplified tax regime (Monotributo in Argentina).
	TagSimplifiedRegime cbc.Key = "simplified-regime"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
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

// VAT statuses that require Invoice Type A (B2B with VAT registered customers)
// Based on ARCA documentation: 1 (Responsable Inscripto), 6 (Monotributo),
// 13 (Monotributista Social), 16 (Monotributo Trabajador Independiente Promovido)
var vatStatusesTypeA = []cbc.Code{"1", "6", "13", "16"}

// invoiceCustomerIsB2B checks if the invoice should be Type A.
// Type A is used when the customer has an Argentine tax ID AND is VAT registered.
//
// This function checks VAT status if explicitly set (user provided), otherwise
// uses tax ID as a fallback. Validation ensures consistency after normalization.
func invoiceCustomerIsB2B(doc any) bool {
	inv, ok := doc.(*bill.Invoice)
	if !ok {
		return false
	}
	// Exclude simplified-regime invoices (type C)
	if inv.HasTags(TagSimplifiedRegime) {
		return false
	}
	if inv.Customer == nil {
		return false
	}

	// If VAT status is explicitly set, use it to determine type
	// This allows users to override the default behavior (e.g., AR tax ID with exempt status → Type B)
	vatStatus := inv.Customer.Ext[ExtKeyVATStatus]
	if vatStatus != "" {
		return vatStatus.In(vatStatusesTypeA...)
	}

	// VAT status not set - use tax ID as fallback
	// Customer with AR tax ID will be normalized to status "1" (Responsable Inscripto) → Type A
	if inv.Customer.TaxID != nil && inv.Customer.TaxID.Country == l10n.AR.Tax() {
		return true
	}

	return false
}

// invoiceCustomerIsB2C checks if the invoice should be Type B.
// Type B is used for final consumers, foreign customers, or when there's no customer.
//
// This is the logical opposite of invoiceCustomerIsB2B - all invoices without
// the simplified-regime tag are either B2B or B2C, making these mutually exclusive.
func invoiceCustomerIsB2C(doc any) bool {
	inv, ok := doc.(*bill.Invoice)
	if !ok {
		return false
	}
	// Exclude simplified-regime invoices (type C)
	if inv.HasTags(TagSimplifiedRegime) {
		return false
	}
	// B2C is the opposite of B2B
	return !invoiceCustomerIsB2B(doc)
}

var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// ** Invoice C - Monotributista transactions (simplified-regime tag required) **
			// These must be first as they require explicit tags and should take precedence
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
			// ** Invoice B - Final consumers and foreign customers (B2C - automatic) **
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
				Filter: invoiceCustomerIsB2C,
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
				Filter: invoiceCustomerIsB2C,
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
				Filter: invoiceCustomerIsB2C,
				Ext: tax.Extensions{
					ExtKeyDocType: "8",
				},
			},
			// ** Invoice A - B2B with VAT registered customer (automatic) **
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
				Filter: invoiceCustomerIsB2B,
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
				Filter: invoiceCustomerIsB2B,
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
				Filter: invoiceCustomerIsB2B,
				Ext: tax.Extensions{
					ExtKeyDocType: "3",
				},
			},
		},
	},
}
