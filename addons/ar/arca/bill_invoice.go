package arca

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
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
	} else if inv.Customer != nil && !inv.Customer.Ext.IsZero() {
		// Check customer VAT status
		vatStatus := inv.Customer.Ext.Get(ExtKeyVATStatus)
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
		inv.Tax = inv.Tax.MergeExtensions(tax.ExtensionsOf(tax.ExtMap{
			ExtKeyDocType: docType,
		}))
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
	inv.Tax = inv.Tax.MergeExtensions(tax.ExtensionsOf(tax.ExtMap{
		ExtKeyConcept: code,
	}))
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("24", "invoice must be in ARS or provide exchange rate for conversion", currency.CanConvertTo(currency.ARS)),
		rules.Field("series",
			rules.Assert("01", "series is required", is.Present),
			rules.Assert("02", "series must be a valid number between 1 and 99998",
				is.FuncError("valid series", invoiceSeriesValid),
			),
		),
		rules.Field("tax",
			rules.Assert("03", "tax is required", is.Present),
			rules.Field("ext",
				rules.Assert("04",
					fmt.Sprintf("tax requires '%s' extension", ExtKeyDocType),
					tax.ExtensionsRequire(ExtKeyDocType),
				),
			),
		),
		// Invoice type vs doc type cross-checks (object-level)
		rules.Assert("05", "invoice type is credit-note but ar-arca-doc-type is not a credit note",
			is.Func("type matches doc type credit note", invoiceTypeMatchesDocTypeCreditNote),
		),
		rules.Assert("06", "invoice type is debit-note but ar-arca-doc-type is not a debit note",
			is.Func("type matches doc type debit note", invoiceTypeMatchesDocTypeDebitNote),
		),
		rules.Assert("07", "ar-arca-doc-type is a credit note but invoice type is not credit-note",
			is.Func("doc type credit note matches type", invoiceDocTypeCreditNoteMatchesType),
		),
		rules.Assert("08", "doc type is a debit note but invoice type is not debit-note",
			is.Func("doc type debit note matches type", invoiceDocTypeDebitNoteMatchesType),
		),
		// Customer requirement depends on doc type
		rules.When(is.Func("customer required", invoiceCustomerRequired),
			rules.Field("customer",
				rules.Assert("09", "customer is required", is.Present),
			),
		),
		// Customer field validations
		rules.Field("customer",
			rules.Assert("10",
				fmt.Sprintf("must have a tax_id, or an identity with ext '%s'", ExtKeyIdentityType),
				is.Func("has tax ID or identity", invoiceCustomerHasTaxIDOrIdentity),
			),
			rules.Field("ext",
				rules.Assert("11",
					fmt.Sprintf("customer requires '%s' extension", ExtKeyVATStatus),
					tax.ExtensionsRequire(ExtKeyVATStatus),
				),
			),
		),
		// VAT status vs doc type cross-checks (object-level)
		rules.Assert("12",
			"document type 49 (Used Goods Purchase Invoice) requires customer VAT status to be 5 (Final Consumer)",
			is.Func("doc type 49 vat status", invoiceDocType49VATStatus),
		),
		rules.Assert("13",
			fmt.Sprintf("document type A requires customer VAT status to be one of %v", vatStatusesTypeA),
			is.Func("doc type A vat status", invoiceDocTypeAVATStatus),
		),
		rules.Assert("14",
			fmt.Sprintf("document type B cannot have customer VAT status %v", vatStatusesTypeA),
			is.Func("doc type B vat status", invoiceDocTypeBVATStatus),
		),
		// Services require ordering and payment
		rules.When(is.Func("concept is services", invoiceConceptIsServices),
			rules.Field("ordering",
				rules.Assert("15", "ordering is required for services", is.Present),
				rules.Field("period",
					rules.Assert("16", "ordering period is required for services", is.Present),
				),
			),
			rules.Field("payment",
				rules.Assert("17", "payment is required for services", is.Present),
				rules.Field("terms",
					rules.Assert("18", "payment terms are required for services", is.Present),
					rules.Field("due_dates",
						rules.Assert("19", "payment due dates are required for services", is.Present),
					),
				),
			),
		),
		// Products must not have payment due dates
		rules.When(is.Func("concept is goods", invoiceConceptIsGoods),
			rules.Field("payment",
				rules.Field("terms",
					rules.Field("due_dates",
						rules.Assert("20", "payment due dates must not be set for goods", is.Empty),
					),
				),
			),
		),
		// Credit/debit notes require preceding documents
		rules.When(bill.InvoiceTypeIn(bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
			rules.Field("preceding",
				rules.Assert("21", "preceding documents are required for credit/debit notes", is.Present),
			),
		),
		// All preceding documents require doc type extension
		rules.Field("preceding",
			rules.Each(
				rules.Field("ext",
					rules.Assert("22",
						fmt.Sprintf("preceding document requires '%s' extension", ExtKeyDocType),
						tax.ExtensionsRequire(ExtKeyDocType),
					),
				),
			),
		),
		// Type C invoices must not have taxes on lines
		rules.When(is.Func("doc type is C", invoiceDocTypeIsC),
			rules.Field("lines",
				rules.Each(
					rules.Field("taxes",
						rules.Assert("23",
							"type C invoices (simplified tax scheme) must not have taxes on lines",
							is.Empty,
						),
					),
				),
			),
		),
		// Type T invoices require the tourism relation extension
		rules.When(is.Func("doc type is T", invoiceDocTypeIsT),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("25",
						fmt.Sprintf("tourism invoice requires '%s' extension", ExtKeyTourismRelation),
						tax.ExtensionsRequire(ExtKeyTourismRelation),
					),
				),
			),
			rules.Field("customer",
				rules.Field("addresses",
					rules.Assert("27", "tourism invoice customer requires an address", is.Present),
				),
			),
			rules.Field("lines",
				rules.Each(
					rules.Field("taxes",
						rules.Each(
							rules.Field("ext",
								rules.Assert("26",
									fmt.Sprintf("tourism invoice line requires '%s' extension", ExtKeyTourismCode),
									tax.ExtensionsRequire(ExtKeyTourismCode),
								),
								rules.Assert("28",
									fmt.Sprintf("tourism invoice line VAT rate must be '5' (21%%) via '%s'", ExtKeyVATRate),
									tax.ExtensionsHasCodes(ExtKeyVATRate, "5"),
								),
							),
						),
					),
				),
			),
		),
	)
}

func invoiceSeriesValid(val any) error {
	s, ok := val.(cbc.Code)
	if !ok || s == "" {
		return nil
	}
	if !regexp.MustCompile(`^\d+$`).MatchString(s.String()) {
		return errors.New("must be a number")
	}
	i, err := strconv.Atoi(s.String())
	if err != nil {
		return errors.New("must be a number")
	}
	if i < 1 || i > 99998 {
		return errors.New("must be between 1 and 99998")
	}
	return nil
}

func invoiceTypeMatchesDocTypeCreditNote(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	if inv.Type == bill.InvoiceTypeCreditNote && !docType.In(DocTypesCreditNote...) {
		return false
	}
	return true
}

func invoiceTypeMatchesDocTypeDebitNote(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	if inv.Type == bill.InvoiceTypeDebitNote && !docType.In(DocTypesDebitNote...) {
		return false
	}
	return true
}

func invoiceDocTypeCreditNoteMatchesType(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	if docType.In(DocTypesCreditNote...) && inv.Type != bill.InvoiceTypeCreditNote {
		return false
	}
	return true
}

func invoiceDocTypeDebitNoteMatchesType(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	if docType.In(DocTypesDebitNote...) && inv.Type != bill.InvoiceTypeDebitNote {
		return false
	}
	return true
}

func invoiceCustomerRequired(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	return !docType.In(append(DocTypesB, TypeUsedGoodsPurchaseInvoice)...)
}

func invoiceCustomerHasTaxIDOrIdentity(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true // nil customer is handled by the "customer required" rule
	}
	if p.TaxID != nil {
		return true
	}
	if org.IdentityForExtKey(p.Identities, ExtKeyIdentityType) != nil {
		return true
	}
	return false
}

func invoiceDocType49VATStatus(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Customer == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	if docType != TypeUsedGoodsPurchaseInvoice {
		return true
	}
	vatStatus := inv.Customer.Ext.Get(ExtKeyVATStatus)
	if vatStatus == "" {
		return true
	}
	return vatStatus == "5" // Final Consumer
}

func invoiceDocTypeAVATStatus(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Customer == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	if !docType.In(DocTypesA...) {
		return true
	}
	vatStatus := inv.Customer.Ext.Get(ExtKeyVATStatus)
	if vatStatus == "" {
		return true
	}
	return vatStatus.In(vatStatusesTypeA...)
}

func invoiceDocTypeBVATStatus(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Customer == nil {
		return true
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	if !docType.In(DocTypesB...) {
		return true
	}
	vatStatus := inv.Customer.Ext.Get(ExtKeyVATStatus)
	if vatStatus == "" {
		return true
	}
	return !vatStatus.In(vatStatusesTypeA...)
}

func invoiceConceptIsServices(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	if invoiceDocTypeIsT(inv) {
		return false // type T (tourism) invoices don't require ordering/payment
	}
	concept := inv.Tax.GetExt(ExtKeyConcept)
	return concept.In("2", "3") // Services or Products and services
}

func invoiceConceptIsGoods(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	if invoiceDocTypeIsT(inv) {
		return false // type T (tourism) invoices don't use concept-based payment rules
	}
	concept := inv.Tax.GetExt(ExtKeyConcept)
	return concept == "1" // Products only
}

func invoiceDocTypeIsC(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	return docType.In(DocTypesC...)
}

func invoiceDocTypeIsT(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	docType := inv.Tax.GetExt(ExtKeyDocType)
	return docType.In(DocTypesT...)
}
