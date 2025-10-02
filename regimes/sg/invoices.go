package sg

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

func validateSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		// In Singapore, all invoices must have a tax ID and name of the issuer,
		// there are no exceptions to this rule.
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&p.Name,
			validation.Required,
		),
	)
}
