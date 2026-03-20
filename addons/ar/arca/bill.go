package arca

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
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

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		// Series
		rules.Field("series",
			rules.Assert("01", "cannot be blank", is.Present),
			rules.AssertIfPresent("02", "must be a number", is.StringFunc("numeric", invoiceSeriesNumeric)),
			rules.AssertIfPresent("03", "must be between 1 and 99998", is.StringFunc("range", invoiceSeriesInRange)),
		),
		// Tax
		rules.Field("tax",
			rules.Assert("04", "cannot be blank", is.Present),
			rules.Field("ext",
				rules.Assert("05", "ar-arca-doc-type: required", tax.ExtensionsRequire(ExtKeyDocType)),
			),
		),
		// Invoice type vs doc type alignment
		rules.Assert("10", "invoice type is credit-note but ar-arca-doc-type is not a credit note",
			is.Func("cn-doctype", invoiceTypeCreditNoteMatchesDocType),
		),
		rules.Assert("11", "invoice type is debit-note but ar-arca-doc-type is not a debit note",
			is.Func("dn-doctype", invoiceTypeDebitNoteMatchesDocType),
		),
		rules.Assert("12", "ar-arca-doc-type is a credit note but invoice type is not credit-note",
			is.Func("doctype-cn", invoiceDocTypeCreditNoteMatchesType),
		),
		rules.Assert("13", "doc type is a debit note but invoice type is not debit-note",
			is.Func("doctype-dn", invoiceDocTypeDebitNoteMatchesType),
		),
		// Customer required (not for type B or type 49)
		rules.When(is.Func("needs customer", invoiceRequiresCustomer),
			rules.Field("customer",
				rules.Assert("14", "cannot be blank", is.Present),
			),
		),
		// Customer validation
		rules.Field("customer",
			rules.Assert("15", fmt.Sprintf("must have a tax_id, or an identity with ext '%s'", ExtKeyIdentityType),
				is.Func("customer has id", invoiceCustomerHasID),
			),
			rules.Field("tax_id",
				rules.Field("code",
					rules.Assert("16", "cannot be blank", is.Present),
				),
			),
			rules.Field("ext",
				rules.Assert("17", "ar-arca-vat-status: required", tax.ExtensionsRequire(ExtKeyVATStatus)),
			),
		),
		// Doc type 49 requires Final Consumer VAT status
		rules.Assert("18", "document type 49 (Used Goods Purchase Invoice) requires customer VAT status to be 5 (Final Consumer)",
			is.Func("type49-vat", invoiceType49CustomerVATStatus),
		),
		// Type A requires specific VAT statuses
		rules.Assert("19", fmt.Sprintf("document type A requires customer VAT status to be one of %v", vatStatusesTypeA),
			is.Func("typeA-vat", invoiceTypeACustomerVATStatus),
		),
		// Type B cannot have type A VAT statuses
		rules.Assert("20", fmt.Sprintf("document type B cannot have customer VAT status %v", vatStatusesTypeA),
			is.Func("typeB-vat", invoiceTypeBCustomerVATStatus),
		),
		// Lines: type C invoices must not have taxes
		rules.When(is.Func("type C", invoiceIsTypeC),
			rules.Field("lines",
				rules.Each(
					rules.Field("taxes",
						rules.Assert("21", "type C invoices (simplified tax scheme) must not have taxes on lines", is.Empty),
					),
				),
			),
		),
		// Ordering required for service invoices
		rules.When(is.Func("services", invoiceConceptIncludesServices),
			rules.Field("ordering",
				rules.Assert("22", "cannot be blank", is.Present),
				rules.Field("period",
					rules.Assert("23", "cannot be blank", is.Present),
				),
			),
		),
		// Payment required for service invoices
		rules.When(is.Func("services payment", invoiceConceptIncludesServices),
			rules.Field("payment",
				rules.Assert("24", "cannot be blank", is.Present),
				rules.Field("terms",
					rules.Assert("25", "cannot be blank", is.Present),
					rules.Field("due_dates",
						rules.Assert("26", "cannot be blank", is.Present),
					),
				),
			),
		),
		// Products: payment due dates must be blank
		rules.When(is.Func("goods only", invoiceConceptIsGoods),
			rules.Field("payment",
				rules.Field("terms",
					rules.Field("due_dates",
						rules.Assert("27", "must be blank", is.Empty),
					),
				),
			),
		),
		// Preceding required for credit/debit notes
		rules.When(bill.InvoiceTypeIn(bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
			rules.Field("preceding",
				rules.Assert("28", "cannot be blank", is.Present),
			),
		),
		// Preceding doc type required
		rules.Field("preceding",
			rules.Each(
				rules.Field("ext",
					rules.Assert("29", "ar-arca-doc-type: required", tax.ExtensionsRequire(ExtKeyDocType)),
				),
			),
		),
	)
}

// invoiceSeriesNumeric returns true when the series consists only of digits
// and can be parsed as an integer (no overflow).
func invoiceSeriesNumeric(s string) bool {
	if !regexp.MustCompile(`^\d+$`).MatchString(s) {
		return false
	}
	_, err := strconv.Atoi(s)
	return err == nil
}

// invoiceSeriesInRange returns true when the series, parsed as an integer, is in the range [1, 99998].
func invoiceSeriesInRange(s string) bool {
	i, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return i >= 1 && i <= 99998
}

func invoiceTypeCreditNoteMatchesDocType(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	return inv.Type != bill.InvoiceTypeCreditNote || docType.In(DocTypesCreditNote...)
}

func invoiceTypeDebitNoteMatchesDocType(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	return inv.Type != bill.InvoiceTypeDebitNote || docType.In(DocTypesDebitNote...)
}

func invoiceDocTypeCreditNoteMatchesType(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	return !docType.In(DocTypesCreditNote...) || inv.Type == bill.InvoiceTypeCreditNote
}

func invoiceDocTypeDebitNoteMatchesType(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	return !docType.In(DocTypesDebitNote...) || inv.Type == bill.InvoiceTypeDebitNote
}

// invoiceRequiresCustomer returns true when the customer field is required
// (i.e. doc type is not type B and not type 49).
func invoiceRequiresCustomer(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	return !docType.In(append(DocTypesB, TypeUsedGoodsPurchaseInvoice)...)
}

// invoiceCustomerHasID returns true when the customer has a tax ID or
// an identity with the ARCA identity type extension.
func invoiceCustomerHasID(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true // nil customer is valid (handled by Required check)
	}
	return p.TaxID != nil || org.IdentityForExtKey(p.Identities, ExtKeyIdentityType) != nil
}

func invoiceIsTypeC(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.Tax.GetExt(ExtKeyDocType).In(DocTypesC...)
}

func invoiceConceptIncludesServices(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.Tax.GetExt(ExtKeyConcept).In("2", "3")
}

func invoiceConceptIsGoods(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.Tax.GetExt(ExtKeyConcept) == "1"
}

// invoiceType49CustomerVATStatus returns true unless the doc type is 49
// and the customer VAT status is not "5" (Final Consumer).
func invoiceType49CustomerVATStatus(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	if inv.Tax.GetExt(ExtKeyDocType) != TypeUsedGoodsPurchaseInvoice {
		return true // not type 49, skip
	}
	if inv.Customer == nil {
		return true // no customer, skip
	}
	if !invoiceCustomerHasID(inv.Customer) {
		return true // customer has no valid ID, skip VAT status check
	}
	return inv.Customer.Ext[ExtKeyVATStatus] == "5"
}

// invoiceTypeACustomerVATStatus returns true unless the doc type is type A
// and the customer VAT status is not in vatStatusesTypeA.
func invoiceTypeACustomerVATStatus(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	if !inv.Tax.GetExt(ExtKeyDocType).In(DocTypesA...) {
		return true // not type A, skip
	}
	if inv.Customer == nil {
		return true // no customer, skip
	}
	if !invoiceCustomerHasID(inv.Customer) {
		return true // customer has no valid ID, skip VAT status check
	}
	vatStatus := inv.Customer.Ext[ExtKeyVATStatus]
	if vatStatus == "" {
		return true // no VAT status, handled by Required check
	}
	return vatStatus.In(vatStatusesTypeA...)
}

// invoiceTypeBCustomerVATStatus returns true unless the doc type is type B
// and the customer VAT status is in vatStatusesTypeA (which is not allowed for type B).
func invoiceTypeBCustomerVATStatus(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	if !inv.Tax.GetExt(ExtKeyDocType).In(DocTypesB...) {
		return true // not type B, skip
	}
	if inv.Customer == nil {
		return true // no customer, skip
	}
	if !invoiceCustomerHasID(inv.Customer) {
		return true // customer has no valid ID, skip VAT status check
	}
	vatStatus := inv.Customer.Ext[ExtKeyVATStatus]
	if vatStatus == "" {
		return true // no VAT status, handled by Required check
	}
	return !vatStatus.In(vatStatusesTypeA...)
}
