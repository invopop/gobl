package choruspro

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func orgPartyRules() *rules.Set {
	return rules.For(new(org.Party),
		// Scheme extension required
		rules.Field("ext",
			rules.Assert("01", "scheme extension is required",
				tax.ExtensionsRequire(ExtKeyScheme),
			),
		),
		// When scheme is "1", require SIRET identity
		rules.When(
			partySchemeIs("1"),
			rules.Field("identities",
				rules.Assert("02", "identities must have a SIRET entry for scheme '1'",
					is.Func("has SIRET", partyHasSIRET),
				),
			),
			rules.Field("tax_id",
				rules.Field("country",
					rules.Assert("03", "tax ID must be 'FR' for scheme '1'",
						is.In(l10n.TaxCountryCode(l10n.FR)),
					),
				),
			),
		),
		// When scheme is not "1" and ext present, reject SIRET
		rules.When(
			partySchemeIsNot("1"),
			rules.Field("identities",
				rules.Assert("04", "identities cannot have a SIRET entry when not '1' scheme",
					is.Func("no SIRET", partyHasNoSIRET),
				),
			),
		),
		// When scheme is "2", require EU non-French
		rules.When(
			partySchemeIs("2"),
			rules.Field("tax_id",
				rules.Field("country",
					rules.Assert("05", "tax ID country must be a non-French, EU company with scheme '2'",
						is.NotIn(l10n.TaxCountryCode(l10n.FR)),
					),
					rules.Assert("06", "tax ID country must be a member of the EU with scheme '2'",
						is.Func("EU member", countryIsEU),
					),
				),
			),
		),
		// When scheme is "3", require non-EU
		rules.When(
			partySchemeIs("3"),
			rules.Field("tax_id",
				rules.Field("country",
					rules.Assert("07", "tax ID country must be a non-EU company with scheme '3'",
						is.Func("non-EU", countryIsNotEU),
					),
				),
			),
		),
	)
}

func partySchemeIs(code cbc.Code) rules.Test {
	return is.Func(
		"party scheme is "+string(code),
		func(val any) bool {
			p, ok := val.(*org.Party)
			return ok && p != nil && p.Ext != nil && p.Ext.Get(ExtKeyScheme) == code
		},
	)
}

func partySchemeIsNot(code cbc.Code) rules.Test {
	return is.Func(
		"party scheme is not "+string(code),
		func(val any) bool {
			p, ok := val.(*org.Party)
			if !ok || p == nil || p.Ext == nil {
				return false
			}
			scheme := p.Ext.Get(ExtKeyScheme)
			return scheme != "" && scheme != code
		},
	)
}

func partyHasSIRET(val any) bool {
	ids, ok := val.([]*org.Identity)
	if !ok {
		return false
	}
	for _, id := range ids {
		if id != nil && id.Type == fr.IdentityTypeSIRET {
			return true
		}
	}
	return false
}

func partyHasNoSIRET(val any) bool {
	ids, ok := val.([]*org.Identity)
	if !ok {
		return true
	}
	for _, id := range ids {
		if id != nil && id.Type == fr.IdentityTypeSIRET {
			return false
		}
	}
	return true
}

func countryIsEU(val any) bool {
	country, ok := val.(l10n.TaxCountryCode)
	if !ok || country == "" {
		return true // skip check for empty
	}
	return l10n.Unions().Code(l10n.EU).HasMember(l10n.Code(country))
}

func countryIsNotEU(val any) bool {
	country, ok := val.(l10n.TaxCountryCode)
	if !ok || country == "" {
		return true // skip check for empty
	}
	return !l10n.Unions().Code(l10n.EU).HasMember(l10n.Code(country))
}
