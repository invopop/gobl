package se

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// invoiceValidator adds validation checks to invoices which are relevant
// for the region.
func validateInvoice(inv *bill.Invoice) error {
	if inv.Tags.HasTags(tax.TagSimplified) {
		// Simplified invoices only require a supplier tax ID.
		return validation.ValidateStruct(inv,
			validation.Field(&inv.Customer,
				validation.When(
					inv.Customer != nil,
					validation.Empty,
				),
				validation.Nil,
				validation.Skip,
			),
			validation.Field(&inv.Supplier,
				validation.Required,
				validation.By(validateSupplier),
				validation.By(validateSupplierSimplifiedInvoice),
				validation.Skip,
			),
		)
	}

	// Standard invoices require a supplier and customer.
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(validateOrgParty),
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.Required,
			validation.By(validateOrgParty),
			validation.Skip,
		),
	)
}

// validateOrgParty holds the common checks for both the supplier and customer.
func validateOrgParty(value any) error {
	party, ok := value.(*org.Party)
	if !ok || party == nil {
		return nil
	}

	// Name and addresses are always required.
	return validation.ValidateStruct(party,
		validation.Field(&party.Name,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&party.Addresses,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&party.Identities,
			// If the party is registered in Sweden for tax purposes,
			// then its identities must be one of the allowed types.
			validation.When(
				party.TaxID != nil && party.TaxID.Country == l10n.TaxCountryCode(l10n.SE),
				validation.Each(
					validation.By(func(value any) error {
						id, ok := value.(*org.Identity)
						if !ok || id == nil {
							return nil
						}
						if !id.Type.In(IdentityTypeOrgNr, IdentityTypePersonNr, IdentityTypeCoordinationNr) {
							return validation.NewError("type", "must be one of: SE-ON, SE-PN, SE-CN")
						}
						return nil
					}),
				),
			),
			validation.Skip,
		),
	)
}

// validateSupplier checks the supplier's tax ID requirements.
// The supplier's VAT number is always required.
func validateSupplier(value any) error {
	party, ok := value.(*org.Party)
	if !ok || party == nil {
		return nil
	}

	return validation.ValidateStruct(party,
		// Swedish Tax ID (VAT ID) is always required for the supplier,
		// and must have the correct format.
		validation.Field(&party.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.By(func(value any) error {
				tID, ok := value.(*tax.Identity)
				if !ok || tID == nil {
					return nil
				}
				return validateTaxIdentity(tID)
			}),
			validation.Skip,
		),
	)
}

func validateSupplierSimplifiedInvoice(value any) error {
	party, ok := value.(*org.Party)
	if !ok || party == nil {
		return nil
	}

	return validation.ValidateStruct(party,
		validation.Field(&party.Name,
			validation.NilOrNotEmpty,
			validation.Skip,
		),
		validation.Field(&party.Addresses,
			validation.NilOrNotEmpty,
			validation.Skip,
		),
	)
}
