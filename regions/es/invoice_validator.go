package es

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
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
		validation.Field(&inv.TypeCode, validation.In(
			bill.CommercialTypeCode,
			bill.SimplifiedTypeCode,
		)),
		validation.Field(&inv.Preceding, validation.By(v.preceding)),
		validation.Field(&inv.Supplier, validation.Required, validation.By(v.supplier)),
		validation.Field(&inv.Customer, validation.When(
			inv.TypeCode != bill.SimplifiedTypeCode,
			validation.Required,
			validation.By(v.customer),
		)),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	obj, ok := value.(*org.Party)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID, validation.Required, ValidTaxID),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	obj, ok := value.(*org.Party)
	if !ok {
		return nil
	}
	if obj.TaxID == nil || obj.TaxID.Country != l10n.ES {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID, ValidTaxID),
	)
}

func (v *invoiceValidator) preceding(value interface{}) error {
	obj, ok := value.(*bill.Preceding)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Period, validation.Required),
		validation.Field(&obj.Corrections, validation.Required, validation.In(correctionReasonKeys()...)),
		validation.Field(&obj.CorrectionMethod, validation.Required, validation.In(correctionMethodKeys()...)),
	)
}
