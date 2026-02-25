package sa

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
		validation.Field(&inv.Customer,
			validation.When(
				!isSimplified(inv),
				validation.Required,
				validation.By(validateCustomer),
			),
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
		// ZATCA requires supplier tax ID and name on all invoices.
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

func validateCustomer(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		// Standard invoices require a customer name.
		validation.Field(&p.Name,
			validation.Required,
		),
	)
}

func isSimplified(inv *bill.Invoice) bool {
	return inv.HasTags(tax.TagSimplified)
}
