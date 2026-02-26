package tr

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
	)
}

// validateSupplier ensures the supplier always has a tax ID (VKN or TCKN) on Turkish invoices.
func validateSupplier(value any) error {
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
