package arca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

const (
	// TagMonotax is used for Invoice C - when the supplier is under the
	// Monotributo regime (simplified unified tax for small taxpayers in Argentina).
	TagMonotax cbc.Key = "monotax"
)

var invoiceCorrectionDefinitions = tax.CorrectionSet{
	{
		Schema: bill.ShortSchemaInvoice,
		Types: []cbc.Key{
			bill.InvoiceTypeCreditNote,
			bill.InvoiceTypeDebitNote,
		},
		Extensions: []cbc.Key{
			ExtKeyDocType,
		},
	},
}

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagMonotax,
			Name: i18n.String{
				i18n.EN: "Monotax",
				i18n.ES: "Monotributo",
			},
			Desc: i18n.String{
				i18n.EN: "Invoice C: Supplier is under the Monotributo regime (simplified unified tax for small taxpayers).",
				i18n.ES: "Factura C: El proveedor está bajo el régimen de Monotributo.",
			},
		},
	},
}

func normalizeBillInvoice(inv *bill.Invoice) {
	normalizeBillInvoiceCustomerVATStatus(inv.Customer)
	normalizeBillInvoiceTaxDocType(inv)
	normalizeBillInvoiceTaxConcept(inv)
}

func normalizeBillInvoiceCustomerVATStatus(p *org.Party) {
	if p == nil {
		return
	}
	switch {
	case p.TaxID == nil:
		// No tax ID: Final Consumer or Uncategorized
		p.Ext = p.Ext.SetOneOf(ExtKeyVATStatus,
			"5",  // Final Consumer
			"4",  // Exempt Subject
			"7",  // Uncategorized Subject
			"10", // VAT Exempt Law 19640
			"15", // VAT Not Applicable
		)
	case p.TaxID.Country != l10n.AR.Tax():
		// Foreign tax ID: Foreign Customer or Supplier
		p.Ext = p.Ext.SetOneOf(ExtKeyVATStatus,
			"9", // Foreign Customer
			"8", // Foreign Supplier
		)
	default:
		// AR tax ID: any valid AR customer status
		p.Ext = p.Ext.SetOneOf(ExtKeyVATStatus,
			"1",  // Registered VAT Company
			"6",  // Monotributo Responsible
			"13", // Social Monotributista
			"16", // Promoted Independent Worker Monotributista
			"4",  // Exempt Subject
			"7",  // Uncategorized Subject
			"10", // VAT Exempt Law 19640
			"15", // VAT Not Applicable
		)
	}
}

func normalizeBillInvoiceTaxDocType(inv *bill.Invoice) {
	if inv.Tax == nil {
		inv.Tax = new(bill.Tax)
	}

	// Skip if doc type is already set
	if inv.Tax.GetExt(ExtKeyDocType) != "" {
		return
	}

	// Determine the doc type category (A, B, or C)
	var docType cbc.Code

	// Check for monotax tag (Type C)
	if inv.Tags.HasTags(TagMonotax) {
		docType = getDocTypeForCategory("C", inv.Type)
	} else if inv.Customer != nil && inv.Customer.Ext != nil {
		// Check customer VAT status
		vatStatus := inv.Customer.Ext[ExtKeyVATStatus]
		if vatStatus.In(vatStatusesTypeA...) {
			// Type A for VAT status 1 (Registered VAT Company), 6 (Monotributo Responsible), 13 (Social Monotributista), 16 (Promoted Independent Worker Monotributista)
			docType = getDocTypeForCategory("A", inv.Type)
		} else {
			// Type B for other VAT statuses
			docType = getDocTypeForCategory("B", inv.Type)
		}
	} else {
		// Default to Type B if no customer
		docType = getDocTypeForCategory("B", inv.Type)
	}

	// Set the doc type extension
	if docType != "" {
		inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
			ExtKeyDocType: docType,
		})
	}
}

// getDocTypeForCategory returns the doc type code based on the category (A, B, C)
// and the invoice type (standard, credit-note, debit-note)
func getDocTypeForCategory(category string, invType cbc.Key) cbc.Code {
	switch category {
	case "A":
		switch invType {
		case bill.InvoiceTypeStandard:
			return "1" // Invoice A (standard)
		case bill.InvoiceTypeDebitNote:
			return "2" // Debit Note A
		case bill.InvoiceTypeCreditNote:
			return "3" // Credit Note A
		}
	case "B":
		switch invType {
		case bill.InvoiceTypeStandard:
			return "6" // Invoice B (standard)
		case bill.InvoiceTypeDebitNote:
			return "7" // Debit Note B
		case bill.InvoiceTypeCreditNote:
			return "8" // Credit Note B
		}
	case "C":
		switch invType {
		case bill.InvoiceTypeStandard:
			return "11" // Invoice C (standard)
		case bill.InvoiceTypeDebitNote:
			return "12" // Debit Note C
		case bill.InvoiceTypeCreditNote:
			return "13" // Credit Note C
		}
	}
	return ""
}

func normalizeBillInvoiceTaxConcept(inv *bill.Invoice) {
	var hasGoods, hasServices bool
	for _, line := range inv.Lines {
		if line.Item == nil || line.Item.Key != org.ItemKeyGoods {
			hasServices = true
		} else {
			hasGoods = true
		}
	}

	var code cbc.Code
	switch {
	case hasGoods && hasServices:
		code = "3" // Products and services
	case hasGoods:
		code = "1" // Products
	case hasServices:
		code = "2" // Services
	default:
		return
	}

	if inv.Tax == nil {
		inv.Tax = new(bill.Tax)
	}
	inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
		ExtKeyConcept: code,
	})
}
