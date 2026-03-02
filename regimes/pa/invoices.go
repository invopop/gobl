package pa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	if inv == nil {
		return nil
	}

	simplified := inv.Tags.HasTags(tax.TagSimplified)

	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(!simplified,
				validation.Required,
			),
			validation.Skip,
		),
	)
}

// validateInvoiceSupplier requires a TaxID with code. In Panama, all businesses
// must have a RUC to operate — it is a general taxpayer ID, not a VAT registration.
func validateInvoiceSupplier(value any) error {
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
