package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// invoiceValidator adds validation checks to invoices which are relevant
// for the region.
type invoiceValidator struct {
	inv *bill.Invoice
}

func validateInvoice(inv *bill.Invoice) error {
	v := &invoiceValidator{inv: inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Tax.ContainsTag(TagTicketBAI) && inv.Type.In(correctionTypes...),
				validation.Required,
			),
			validation.Each(validation.By(v.preceding)),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(v.supplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(v.customer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(v.validateLine),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	// Customers must have a tax ID to at least set the country,
	// and Spanish ones should also have an ID. There are more complex
	// rules for exports.
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			validation.When(
				obj.TaxID != nil && obj.TaxID.Country.In("ES"),
				tax.RequireIdentityCode,
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) preceding(value interface{}) error {
	obj, _ := value.(*bill.Preceding)
	if obj == nil {
		return nil
	}

	if v.inv.Tax.ContainsTag(TagTicketBAI) {
		return validation.ValidateStruct(obj,
			validation.Field(&obj.IssueDate, validation.Required),
			validation.Field(&obj.Ext, tax.ExtensionsRequires(ExtKeyTBAICorrection)),
		)
	}
	if v.inv.Tax.ContainsTag(TagFacturaE) {
		return validation.ValidateStruct(obj,
			validation.Field(&obj.IssueDate, validation.Required),
			validation.Field(&obj.Ext, tax.ExtensionsRequires(ExtKeyFacturaECorrection)),
		)
	}

	return nil
}

func (v *invoiceValidator) validateLine(value interface{}) error {
	obj, _ := value.(*bill.Line)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Taxes,
			validation.Each(
				validation.By(v.validateLineTax),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) validateLineTax(value interface{}) error {
	obj, ok := value.(*tax.Combo)
	if obj == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			validation.When(
				v.inv.Tax.ContainsTag(TagTicketBAI) && obj.Rate == tax.RateExempt,
				tax.ExtensionsRequires(ExtKeyTBAIExemption),
			),
			validation.Skip,
		),
	)
}
