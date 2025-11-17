package se

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateBillInvoice(inv *bill.Invoice) error {
	simplified := inv.Tags.HasTags(tax.TagSimplified)
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateBillInvoiceParty(!simplified)),
			validation.By(validateBillInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				simplified,
				validation.Empty,
			).Else(
				validation.Required,
				validation.By(validateBillInvoiceParty(true)),
			),
			validation.Skip,
		),
	)
}

func validateBillInvoiceParty(withAddress bool) validation.RuleFunc {
	return func(value any) error {
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
				validation.When(withAddress,
					validation.Required,
				),
				validation.Skip,
			),
			validation.Field(&party.Identities,
				// If the party is registered in Sweden for tax purposes,
				// then its identities must be one of the allowed types.
				validation.When(
					isSwedishParty(party) && party.TaxID.Code == "",
					org.RequireIdentityType(IdentityTypeOrgNr, IdentityTypePersonNr, IdentityTypeCoordinationNr),
				),
				validation.Skip,
			),
		)
	}
}

func isSwedishParty(party *org.Party) bool {
	if party == nil || party.TaxID == nil {
		return false
	}
	return party.TaxID.Country == l10n.TaxCountryCode(l10n.SE)
}

func validateBillInvoiceSupplier(value any) error {
	party, _ := value.(*org.Party)
	return validation.ValidateStruct(party,
		validation.Field(&party.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}
