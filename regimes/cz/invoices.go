package cz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// validateInvoice checks the Czech invoice requirements.
// Simplified invoices skip supplier validation since businesses
// below the VAT threshold may not have a DIČ.
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

// validateInvoiceSupplier requires either a TaxID (DIČ) or an
// IČO identity. All Czech businesses have an IČO, but only
// VAT-registered businesses are required to show a DIČ.
func validateInvoiceSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.When(
				!hasIdentityICO(p),
				validation.Required,
				tax.RequireIdentityCode,
			),
			validation.Skip,
		),
		validation.Field(&p.Identities,
			validation.When(
				!hasTaxIDCode(p),
				org.RequireIdentityKey(IdentityKeyICO),
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

func hasIdentityICO(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentityForKey(party.Identities, IdentityKeyICO) != nil
}
