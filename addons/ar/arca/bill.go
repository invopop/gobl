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
	if inv.Customer != nil && !inv.Customer.Ext.Has(ExtKeyVATStatus) {
		switch {
		case inv.Customer.TaxID == nil:
			inv.Customer.Ext = inv.Customer.Ext.Set(ExtKeyVATStatus, "5") // Final Consumer
		case inv.Customer.TaxID.Country == l10n.AR.Tax():
			inv.Customer.Ext = inv.Customer.Ext.Set(ExtKeyVATStatus, "1") // Registered VAT Company
		default:
			inv.Customer.Ext = inv.Customer.Ext.Set(ExtKeyVATStatus, "9") // Foreign Customer
		}
	}
	normalizeConcept(inv)
}

func normalizeConcept(inv *bill.Invoice) {
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
				!inv.Tax.GetExt(ExtKeyDocType).In("6", "7", "8", "49"),
				validation.Required,
			),
			validation.By(validateInvoiceCustomer),
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
				inv.Tax.GetExt(ExtKeyConcept).In("2", "3"),
				validation.Required,
				validation.By(validateOrdering),
			),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.When(
				inv.Tax.GetExt(ExtKeyConcept).In("2", "3"),
				validation.Required,
				validation.By(validatePayment),
			),
			validation.When(
				inv.Tax.GetExt(ExtKeyConcept).In("1"),
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

func validateInvoiceCustomer(val any) error {
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
	)
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
		if docType.In(TypeCDocTypes...) {
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
