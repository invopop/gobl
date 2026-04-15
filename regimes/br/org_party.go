package br

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var (
	// Official state acronyms as per IBGE in Brazil
	validStates = []cbc.Code{
		"AC", "AL", "AM", "AP", "BA", "CE", "DF", "ES", "GO",
		"MA", "MG", "MS", "MT", "PA", "PB", "PE", "PI", "PR",
		"RJ", "RN", "RO", "RR", "RS", "SC", "SE", "SP", "TO",
	}
	validPostCode = `^\d{5}-?\d{3}$`
)

func orgPartyRules() *rules.Set {
	return rules.For(new(org.Party),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				isBrazilianParty,
				rules.Field("ext",
					rules.AssertIfPresent("01", fmt.Sprintf("Brazilian party ext must define a valid '%s' code", ExtKeyMunicipality),
						tax.ExtensionHasValidCode(ExtKeyMunicipality),
					),
				),
				rules.Field("addresses",
					rules.Each(
						rules.When(
							is.Expr("string(Country) in ['','BR']"),
							rules.Field("state",
								rules.AssertIfPresent("02", "Brazilian state must be one of the valid states", cbc.InCodes(validStates...)),
							),
							rules.Field("code",
								rules.AssertIfPresent("03", "Brazilian postal code must match the valid format", is.Matches(validPostCode)),
							),
						),
					),
				),
			),
		),
	)
}

var isBrazilianParty = is.Func("is Brazilian party",
	func(value any) bool {
		party, _ := value.(*org.Party)
		if party == nil {
			return false
		}
		if party.TaxID != nil && party.TaxID.Country == l10n.BR.Tax() {
			return true
		}
		for _, addr := range party.Addresses {
			if addr != nil && addr.Country == l10n.BR.ISO() {
				return true
			}
		}
		return false
	})
