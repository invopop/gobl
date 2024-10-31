package de

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// validateInvoice checks to ensure the German invoice is not simplified
// and the supplier contains either a Tax ID (VAT) *or* a Tax Number.
func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.When(
				!isSimplified(inv),
				validation.By(validateInvoiceSupplier),
			),
			validation.Skip,
		),
	)
}

func validateInvoiceSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.When(
				!hasIdentityTaxNumber(p),
				validation.Required,
				tax.RequireIdentityCode,
			),
			validation.Skip,
		),
		validation.Field(&p.Identities,
			validation.When(
				!hasTaxIDCode(p),
				org.RequireIdentityKey(IdentityKeyTaxNumber),
			),
			validation.Skip,
		),
	)
}

func isSimplified(inv *bill.Invoice) bool {
	return inv.HasTags(tax.TagSimplified)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasIdentityTaxNumber(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentityForKey(party.Identities, IdentityKeyTaxNumber) != nil
}
