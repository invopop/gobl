package nl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
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
		validation.Field(&inv.Currency,
			validation.In(currency.EUR),
		),
		validation.Field(&inv.Supplier,
			validation.By(v.supplier),
		),
		validation.Field(&inv.Customer,
			validation.By(v.customer),
		),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	obj, ok := value.(*org.Party)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID, validation.Required, tax.RequireIdentityCode),
	)
}

func (v *invoiceValidator) customer(_ interface{}) error {
	return nil // nothing to do yet
}
