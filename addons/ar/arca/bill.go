package arca

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
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

func normalizeInvoice(inv *bill.Invoice) {
	normalizePartyVATStatus(inv.Customer)
	normalizeTransactionType(inv)
}

// VAT statuses valid for customers without a tax ID (final consumers)
// Default: 5 (Consumidor Final)
// Uncertain: 4, 7, 10, 15 - unclear if they need tax ID, allowing for flexibility
var vatStatusesNoTaxID = []cbc.Code{"5", "4", "7", "10", "15"}

// VAT statuses valid for customers with an Argentine tax ID
// Default: 1 (Responsable Inscripto)
// Certain (Type A): 1 (Responsable Inscripto), 6 (Monotributo), 13 (Monotributista Social),
// 16 (Monotributo Trabajador Independiente Promovido)
//
// Uncertain (Type B): 4 (VAT Exempt), 7 (Uncategorized), 10 (Tierra del Fuego), 15 (Not subject to VAT)
//   - these may or may not require an AR tax ID, allowing for flexibility
var vatStatusesARTaxID = []cbc.Code{"1", "6", "13", "16", "4", "7", "10", "15"}

// VAT statuses valid for customers with a foreign tax ID
// Default: 9 (Cliente del Exterior)
// Also valid: 8 (Proveedor del Exterior)
var vatStatusesForeignTaxID = []cbc.Code{"9", "8"}

func normalizePartyVATStatus(p *org.Party) {
	if p == nil {
		return
	}
	switch {
	case p.TaxID == nil:
		// No tax ID: Final Consumer or Uncategorized
		p.Ext = p.Ext.SetOneOf(ExtKeyVATStatus, "5", vatStatusesNoTaxID...)
	case p.TaxID.Country == l10n.AR.Tax():
		// AR tax ID: any valid AR customer status
		p.Ext = p.Ext.SetOneOf(ExtKeyVATStatus, "1", vatStatusesARTaxID...)
	default:
		// Foreign tax ID: Foreign Customer or Supplier
		p.Ext = p.Ext.SetOneOf(ExtKeyVATStatus, "9", vatStatusesForeignTaxID...)
	}
}

func normalizeTransactionType(inv *bill.Invoice) {
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
		ExtKeyTransactionType: code,
	})
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Series,
			validation.Required,
			validation.By(validateSeries),
		),
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateInvoiceTax),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				!inv.Tax.GetExt(ExtKeyDocType).In(append(DocTypesB, TypeUsedGoodsPurchaseInvoice)...),
				validation.Required,
			),
			validation.By(validateInvoiceCustomer(inv.Tax)),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateInvoiceLine(inv.Tax.GetExt(ExtKeyDocType))),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Ordering,
			validation.When(
				inv.Tax.GetExt(ExtKeyTransactionType).In("2", "3"),
				validation.Required,
				validation.By(validateOrdering),
			),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.When(
				inv.Tax.GetExt(ExtKeyTransactionType).In("2", "3"),
				validation.Required,
				validation.By(validatePayment),
			),
			validation.When(
				inv.Tax.GetExt(ExtKeyTransactionType).In("1"),
				validation.By(validatePaymentNoDueDates),
			),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
				validation.Required,
			),
			validation.Each(
				validation.By(validateInvoicePreceding),
			),
			validation.Skip,
		),
	)
}

func validateSeries(value interface{}) error {
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

func validateInvoiceTax(val any) error {
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

func validateInvoiceCustomer(invTax *bill.Tax) validation.RuleFunc {
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
				validation.By(validateVATStatusMatchesDocType(invTax)),
				validation.Skip,
			),
		)
	}
}

func validateVATStatusMatchesDocType(invTax *bill.Tax) validation.RuleFunc {
	return func(val any) error {
		if invTax == nil {
			return nil
		}

		ext, ok := val.(tax.Extensions)
		if !ok {
			return nil
		}

		docType := invTax.GetExt(ExtKeyDocType)
		vatStatus := ext[ExtKeyVATStatus]

		if vatStatus == "" {
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

func validateOrdering(val any) error {
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

func validatePayment(val any) error {
	payment, ok := val.(*bill.PaymentDetails)
	if !ok || payment == nil {
		return nil
	}
	return validation.ValidateStruct(payment,
		validation.Field(&payment.Terms,
			validation.Required,
			validation.By(validatePaymentTerms),
		),
	)
}

func validatePaymentTerms(val any) error {
	terms, ok := val.(*pay.Terms)
	if !ok || terms == nil {
		return nil
	}
	return validation.ValidateStruct(terms,
		validation.Field(&terms.DueDates,
			validation.Required),
	)
}

func validatePaymentNoDueDates(val any) error {
	payment, ok := val.(*bill.PaymentDetails)
	if !ok || payment == nil || payment.Terms == nil {
		return nil
	}
	return validation.ValidateStruct(payment.Terms,
		validation.Field(&payment.Terms.DueDates,
			validation.Empty),
	)
}

func validateInvoicePreceding(val any) error {
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

func validateInvoiceLine(docType cbc.Code) validation.RuleFunc {
	return func(val any) error {
		line, ok := val.(*bill.Line)
		if !ok || line == nil {
			return nil
		}
		if docType.In(DocTypesC...) {
			return validation.ValidateStruct(line,
				validation.Field(&line.Taxes,
					validation.Empty,
					validation.Skip,
				),
			)
		}
		return nil
	}
}
