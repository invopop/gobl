package arca

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
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
	vatStatus := inv.Customer.Ext[ExtKeyVATStatus]
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
	vatStatus := inv.Customer.Ext[ExtKeyVATStatus]
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
	vatStatus := inv.Customer.Ext[ExtKeyVATStatus]
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
	concept := inv.Tax.GetExt(ExtKeyConcept)
	return concept.In("2", "3") // Services or Products and services
}

func invoiceConceptIsGoods(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
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
