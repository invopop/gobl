package arca

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	// TagSimplifiedScheme is used for Invoice C - when the supplier is under a
	// simplified tax regime (Monotributo in Argentina).
	TagSimplifiedScheme cbc.Key = "simplified-scheme"
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
			Key: TagSimplifiedScheme,
			Name: i18n.String{
				i18n.EN: "Simplified Tax Scheme",
				i18n.ES: "Monotributo",
			},
			Desc: i18n.String{
				i18n.EN: "Invoice C: Supplier is under a simplified tax scheme (Monotributo).",
				i18n.ES: "Factura C: El proveedor está bajo un esquema tributario simplificado (Monotributo).",
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
			VATStatusFinalConsumer,
			VATStatusExemptSubject,
			VATStatusUncategorizedSubject,
			VATStatusVATExemptLaw19640,
			VATStatusVATNotApplicable,
		)
	case p.TaxID.Country != l10n.AR.Tax():
		// Foreign tax ID: Foreign Customer or Supplier
		p.Ext = p.Ext.SetOneOf(ExtKeyVATStatus,
			VATStatusForeignCustomer,
			VATStatusForeignSupplier,
		)
	default:
		// AR tax ID: any valid AR customer status
		p.Ext = p.Ext.SetOneOf(ExtKeyVATStatus,
			VATStatusRegisteredCompany,
			VATStatusMonotributoResponsible,
			VATStatusSocialMonotributista,
			VATStatusPromotedIndependentWorkerMonotributista,
			VATStatusExemptSubject,
			VATStatusUncategorizedSubject,
			VATStatusVATExemptLaw19640,
			VATStatusVATNotApplicable,
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

	// Check for simplified-scheme tag (Type C)
	if inv.Tags.HasTags(TagSimplifiedScheme) {
		docType = getDocTypeForCategory("C", inv.Type)
	} else if inv.Customer != nil && inv.Customer.Ext != nil {
		// Check customer VAT status
		vatStatus := inv.Customer.Ext[ExtKeyVATStatus]
		if vatStatus.In(vatStatusesTypeA...) {
			// Type A for VAT status 1, 6, 13, 16
			docType = getDocTypeForCategory("A", inv.Type)
		} else {
			// Type B for other VAT statuses
			docType = getDocTypeForCategory("B", inv.Type)
		}
	} else {
		// Default to Type B if no customer or no VAT status
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

func validateBillInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Series,
			validation.Required,
			validation.By(validateBillInvoiceSeries),
		),
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateBillInvoiceTax),
			validation.Skip,
		),
		validation.Field(&inv.Type,
			validation.By(validateBillInvoiceType(inv.Tax.GetExt(ExtKeyDocType))),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				!inv.Tax.GetExt(ExtKeyDocType).In(append(DocTypesB, TypeUsedGoodsPurchaseInvoice)...),
				validation.Required,
			),
			validation.By(validateBillInvoiceCustomer(inv.Tax.GetExt(ExtKeyDocType))),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateBillInvoiceLineTaxes(inv.Tax.GetExt(ExtKeyDocType))),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Ordering,
			validation.When(
				inv.Tax.GetExt(ExtKeyConcept).In(ConceptServices, ConceptProductsAndServices),
				validation.Required,
				validation.By(validateBillOrderingPeriod),
			),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.When(
				inv.Tax.GetExt(ExtKeyConcept).In(ConceptServices, ConceptProductsAndServices),
				validation.Required,
				validation.By(validateBillPaymentDetailsServices),
			),
			validation.When(
				inv.Tax.GetExt(ExtKeyConcept).In(ConceptGoods),
				validation.By(validateBillPaymentDetailsGoods),
			),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
				validation.Required,
			),
			validation.Each(
				validation.By(validateBillInvoicePreceding),
			),
			validation.Skip,
		),
	)
}

func validateBillInvoiceTax(val any) error {
	tx, ok := val.(*bill.Tax)
	if !ok || tx == nil {
		return nil
	}
	return validation.ValidateStruct(tx,
		validation.Field(&tx.Ext,
			tax.ExtensionsRequire(ExtKeyDocType),
		),
	)
}

func validateBillInvoiceSeries(value interface{}) error {
	s, ok := value.(cbc.Code)
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

func validateBillInvoiceType(docType cbc.Code) validation.RuleFunc {
	return func(val any) error {
		invType, ok := val.(cbc.Key)
		if !ok {
			return nil
		}

		// Check if invoice type is credit-note but doc type is not a credit note
		if invType == bill.InvoiceTypeCreditNote && !docType.In(DocTypesCreditNote...) {
			return errors.New("invoice type is credit-note but ar-arca-doc-type is not a credit note")
		}

		// Check if invoice type is debit-note but doc type is not a debit note
		if invType == bill.InvoiceTypeDebitNote && !docType.In(DocTypesDebitNote...) {
			return errors.New("invoice type is debit-note but ar-arca-doc-type is not a debit note")
		}

		// Check if doc type is a credit note but invoice type is not credit-note
		if docType.In(DocTypesCreditNote...) && invType != bill.InvoiceTypeCreditNote {
			return errors.New("ar-arca-doc-type is a credit note but invoice type is not credit-note")
		}

		// Check if doc type is a debit note but invoice type is not debit-note
		if docType.In(DocTypesDebitNote...) && invType != bill.InvoiceTypeDebitNote {
			return errors.New("doc type is a debit note but invoice type is not debit-note")
		}

		return nil
	}
}

func validateBillInvoiceCustomer(docType cbc.Code) validation.RuleFunc {
	return func(val any) error {
		p, ok := val.(*org.Party)
		if !ok || p == nil {
			return nil
		}
		if p.TaxID == nil && org.IdentityForExtKey(p.Identities, ExtKeyIdentityType) == nil {
			return fmt.Errorf("must have a tax_id, or an identity with ext '%s'", ExtKeyIdentityType)
		}

		return validation.ValidateStruct(p,
			validation.Field(&p.TaxID,
				tax.RequireIdentityCode,
				validation.Skip,
			),
			validation.Field(&p.Ext,
				tax.ExtensionsRequire(ExtKeyVATStatus),
				validation.By(validateVATStatusMatchesDocType(docType)),
				validation.Skip,
			),
		)
	}
}

func validateVATStatusMatchesDocType(docType cbc.Code) validation.RuleFunc {
	return func(val any) error {
		ext, ok := val.(tax.Extensions)
		if !ok {
			return nil
		}

		vatStatus := ext[ExtKeyVATStatus]

		if vatStatus == "" {
			return nil
		}

		// Doc type 49 (Used Goods Purchase Invoice) requires Final Consumer (5)
		if docType == TypeUsedGoodsPurchaseInvoice {
			if vatStatus != VATStatusFinalConsumer {
				return fmt.Errorf("document type 49 (Used Goods Purchase Invoice) requires customer VAT status to be 5 (Final Consumer)")
			}
			return nil
		}

		if docType.In(DocTypesA...) {
			// Type A invoices require VAT status 1, 6, 13, or 16
			return validation.Validate(vatStatus,
				validation.In(vatStatusesTypeA...).Error(
					fmt.Sprintf("document type A requires customer VAT status to be one of %v", vatStatusesTypeA),
				),
			)
		}

		if docType.In(DocTypesB...) {
			// Type B invoices require VAT status other than 1, 6, 13, or 16
			return validation.Validate(vatStatus,
				validation.NotIn(vatStatusesTypeA...).Error(
					fmt.Sprintf("document type B cannot have customer VAT status %v", vatStatusesTypeA),
				),
			)
		}

		return nil
	}
}

func validateBillOrderingPeriod(val any) error {
	ordering, ok := val.(*bill.Ordering)
	if !ok || ordering == nil {
		return nil
	}
	return validation.ValidateStruct(ordering,
		validation.Field(&ordering.Period,
			validation.Required,
		),
	)
}

func validateBillPaymentDetailsServices(val any) error {
	payment, ok := val.(*bill.PaymentDetails)
	if !ok || payment == nil {
		return nil
	}
	return validation.ValidateStruct(payment,
		validation.Field(&payment.Terms,
			validation.Required,
			validation.By(validatePaymentTermsDueDates),
		),
	)
}

func validatePaymentTermsDueDates(val any) error {
	terms, ok := val.(*pay.Terms)
	if !ok || terms == nil {
		return nil
	}
	return validation.ValidateStruct(terms,
		validation.Field(&terms.DueDates,
			validation.Required),
	)
}

func validateBillPaymentDetailsGoods(val any) error {
	payment, ok := val.(*bill.PaymentDetails)
	if !ok || payment == nil || payment.Terms == nil {
		return nil
	}
	return validation.ValidateStruct(payment.Terms,
		validation.Field(&payment.Terms.DueDates,
			validation.Empty),
	)
}

func validateBillInvoicePreceding(val any) error {
	preceding, ok := val.(*org.DocumentRef)
	if !ok || preceding == nil {
		return nil
	}
	return validation.ValidateStruct(preceding,
		validation.Field(&preceding.Ext,
			tax.ExtensionsRequire(ExtKeyDocType),
		),
	)
}

func validateBillInvoiceLineTaxes(docType cbc.Code) validation.RuleFunc {
	return func(val any) error {
		line, ok := val.(*bill.Line)
		if !ok || line == nil {
			return nil
		}
		if docType.In(DocTypesC...) {
			return validation.ValidateStruct(line,
				validation.Field(&line.Taxes,
					validation.Empty.Error("type C invoices (simplified tax scheme) must not have taxes on lines"),
					validation.Skip,
				),
			)
		}
		return nil
	}
}
