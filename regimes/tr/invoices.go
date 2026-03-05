package tr

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// validateInvoice checks the Turkish invoice for required fields.
func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
	)
}

// validateSupplier ensures the supplier has either a tax ID (VKN) or a
// TCKN identity on Turkish invoices.
func validateSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.When(
				!hasIdentityTCKN(p),
				validation.Required,
				tax.RequireIdentityCode,
			),
			validation.Skip,
		),
		validation.Field(&p.Identities,
			validation.When(
				!hasTaxIDCode(p),
				org.RequireIdentityType(IdentityTypeTCKN),
			),
			validation.Skip,
		),
	)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasIdentityTCKN(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentityForType(party.Identities, IdentityTypeTCKN) != nil
}
