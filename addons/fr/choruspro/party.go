package choruspro

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// validateParty validates the party and its tax ID and identities.
func validateParty(party *org.Party) error {
	return validation.ValidateStruct(party,
		validation.Field(&party.TaxID,
			validation.When(
				!hasIdentityWithScheme(party),
				validation.Required.Error("tax ID scheme must be set when no identity has scheme extension"),
				validation.By(validateTaxID),
			),
		),
		validation.Field(&party.Identities,
			validation.Each(
				validation.By(validateIdentityScheme),
			),
		),
	)
}

func hasIdentityWithScheme(party *org.Party) bool {
	if party == nil {
		return false
	}
	for _, identity := range party.Identities {
		if identity != nil && identity.Ext != nil && identity.Ext.Has(ExtKeyScheme) {
			return true
		}
	}
	return false
}

func validateTaxID(value interface{}) error {
	taxID, ok := value.(*tax.Identity)
	if !ok || taxID == nil {
		return nil
	}

	return validation.ValidateStruct(taxID,
		validation.Field(&taxID.Country,
			validation.Required,
			validation.NotIn(l10n.TaxCountryCode(l10n.FR)).Error("French companies cannot have scheme set at tax identity level"),
			validation.Skip,
		),
		validation.Field(&taxID.Scheme,
			validation.Required,
			validation.In(
				cbc.Code("2"),
				cbc.Code("3"),
				cbc.Code("4"),
				cbc.Code("5"),
				cbc.Code("6"),
			),
			validation.Skip,
		),
	)
}

func validateIdentityScheme(value interface{}) error {
	identity, ok := value.(*org.Identity)
	if !ok || identity == nil {
		return nil
	}

	return validation.ValidateStruct(identity,
		validation.Field(&identity.Ext,
			validation.When(
				identity.Type == fr.IdentityTypeSIRET,
				tax.ExtensionsHasCodes(ExtKeyScheme, "1"),
			),
		),
	)
}
