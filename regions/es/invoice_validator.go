package es

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
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
		// Only commercial and simplified supported at this time for spain.
		// Rectification state determined by Preceding value.
		validation.Field(&inv.TypeKey, validation.In(
			bill.TypeKeyCommercial,
			bill.TypeKeySimplified,
		)),
		validation.Field(&inv.Preceding, validation.By(v.preceding)),
		validation.Field(&inv.Supplier, validation.Required, validation.By(v.supplier)),
		validation.Field(&inv.Customer, validation.When(
			inv.TypeKey != bill.TypeKeySimplified,
			validation.Required,
			validation.By(v.customer),
		)),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID, validation.Required, ValidTaxID.RequireCode()),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return tax.ValidateTaxIdentity(obj.TaxID)
}

func (v *invoiceValidator) preceding(value interface{}) error {
	obj, _ := value.(*bill.Preceding)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Period, validation.Required),
		validation.Field(&obj.Corrections, validation.Required, validation.In(correctionReasonKeys()...)),
		validation.Field(&obj.CorrectionMethod, validation.Required, validation.In(correctionMethodKeys()...)),
	)
}
