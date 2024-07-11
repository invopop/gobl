package gr

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Series, validation.Required),
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceParty),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(validateInvoiceParty),
			validation.By(validateInvoiceCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateInvoiceLine),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.Required,
			validation.By(validateInvoicePayment),
			validation.Skip,
		),
	)
}

func validateInvoiceParty(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

func validateInvoiceCustomer(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Addresses,
			validation.Required,
			validation.Length(1, 0),
			validation.Skip,
		),
	)
}

func validateInvoiceLine(value any) error {
	l, ok := value.(*bill.Line)
	if !ok || l == nil {
		return nil
	}
	return validation.ValidateStruct(l,
		validation.Field(&l.Total,
			num.Positive,
			num.NotZero,
			validation.Skip,
		),
	)
}

func validateInvoicePayment(value any) error {
	p, ok := value.(*bill.Payment)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Instructions,
			validation.Required,
			validation.By(validateInvoicePaymentInstructions),
			validation.Skip,
		),
	)
}

func validateInvoicePaymentInstructions(value any) error {
	i, ok := value.(*pay.Instructions)
	if !ok || i == nil {
		return nil
	}

	return validation.ValidateStruct(i,
		validation.Field(&i.Key,
			validation.Required,
			isValidPaymentMeanKey,
			validation.Skip,
		),
	)
}
