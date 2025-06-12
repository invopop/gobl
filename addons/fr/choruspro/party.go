package choruspro

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeParty(party *org.Party) {
	if party == nil {
		return
	}

	if party.Ext != nil {
		if party.Ext.Get(ExtKeyScheme) != "" {
			return
		}
	}

	if party.TaxID != nil && party.TaxID.Country != "FR" {
		if party.Ext == nil {
			party.Ext = make(tax.Extensions)
		}
		if l10n.Unions().Code(l10n.EU).HasMember(l10n.Code(party.TaxID.Country)) {
			party.Ext = party.Ext.Merge(
				tax.Extensions{
					ExtKeyScheme: "2",
				},
			)
		} else {
			party.Ext = party.Ext.Merge(
				tax.Extensions{
					ExtKeyScheme: "3",
				},
			)
		}
		return
	}

	// If FR or no tax ID we search for a SIRET identity and set the scheme to 1
	for _, identity := range party.Identities {
		if identity.Type == fr.IdentityTypeSIRET {
			if party.Ext == nil {
				party.Ext = make(tax.Extensions)
			}
			party.Ext = party.Ext.Merge(
				tax.Extensions{
					ExtKeyScheme: "1",
				},
			)
			return
		}
	}
}

func validateParty(value interface{}) error {
	party, ok := value.(*org.Party)
	if !ok || party == nil {
		return nil
	}

	return validation.ValidateStruct(party,
		validation.Field(&party.Ext,
			tax.ExtensionsRequire(ExtKeyScheme),
			validation.Skip,
		),
		validation.Field(&party.Identities,
			validation.When(
				party.Ext != nil,
				validation.By(validateIdentities(party.Ext.Get(ExtKeyScheme))),
			),
			validation.Skip,
		),
		validation.Field(&party.TaxID,
			validation.When(
				party.Ext != nil,
				validation.By(validateTaxID(party.Ext.Get(ExtKeyScheme))),
			),
			validation.Skip,
		),
	)
}

func validateIdentities(scheme cbc.Code) validation.RuleFunc {
	return func(value interface{}) error {
		identities, ok := value.([]*org.Identity)
		if !ok || identities == nil {
			return nil
		}

		foundSIRET := false

		for _, identity := range identities {
			if identity.Type == fr.IdentityTypeSIRET {
				foundSIRET = true
				break
			}
		}

		if scheme == "1" && !foundSIRET {
			return validation.NewError("identities", "No SIRET identity found")
		}
		if scheme != "1" && foundSIRET {
			return validation.NewError("identities", "SIRET identity not allowed for this extension")
		}
		return nil
	}
}

func validateTaxID(scheme cbc.Code) validation.RuleFunc {
	return func(value interface{}) error {
		taxID, ok := value.(*tax.Identity)
		if !ok || taxID == nil {
			return nil
		}

		switch scheme {
		case "1":
			return validation.ValidateStruct(taxID,
				validation.Field(&taxID.Country,
					validation.Required,
					validation.In(l10n.TaxCountryCode(l10n.FR)).Error("Customer must be a French company"),
					validation.Skip,
				),
			)
		case "2":
			return validation.ValidateStruct(taxID,
				validation.Field(&taxID.Country,
					validation.Required,
					validation.NotIn(l10n.TaxCountryCode(l10n.FR)).Error("Customer must be a non-French, EU company"),
					validation.By(validateEUCompany),
					validation.Skip,
				),
			)
		case "3":
			return validation.ValidateStruct(taxID,
				validation.Field(&taxID.Country,
					validation.Required,
					validation.NotIn(l10n.TaxCountryCode(l10n.FR)).Error("Customer must be a non-French, EU company"),
					validation.By(validateNonEUCompany),
					validation.Skip,
				),
			)
		}
		return nil
	}
}

func validateEUCompany(value interface{}) error {
	country, ok := value.(l10n.TaxCountryCode)
	if !ok || country == "" {
		return nil
	}

	if !l10n.Unions().Code(l10n.EU).HasMember(l10n.Code(country)) {
		return validation.NewError("taxID", "Customer must be a member of the EU")
	}

	return nil
}

func validateNonEUCompany(value interface{}) error {
	country, ok := value.(l10n.TaxCountryCode)
	if !ok || country == "" {
		return nil
	}

	if l10n.Unions().Code(l10n.EU).HasMember(l10n.Code(country)) {
		return validation.NewError("taxID", "Customer must be a non-EU company")
	}

	return nil
}
