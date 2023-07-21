package gb

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/regimes/common"
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
		validation.Field(&inv.Currency, validation.In(currency.GBP)),
		validation.Field(&inv.Supplier, validation.Required),
		validation.Field(&inv.Customer, validation.When(
			!inv.Tax.ContainsTag(common.TagSimplified),
			validation.Required,
		)),
	)
}
