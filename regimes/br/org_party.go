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

// normalizeParty fills in the country of a party's addresses when missing,
// inferring it from the party's tax ID or first identity that declares a country.
func normalizeParty(p *org.Party) {
	if p == nil || len(p.Addresses) == 0 {
		return
	}
	country := partyCountry(p)
	if country.Empty() {
		return
	}
	for _, addr := range p.Addresses {
		if addr != nil && addr.Country.Empty() {
			addr.Country = country
		}
	}
}

// partyCountry derives a party's country from its tax ID or, failing that, from
// the first identity that declares a country.
func partyCountry(p *org.Party) l10n.ISOCountryCode {
	if p == nil {
		return ""
	}
	if p.TaxID != nil && !p.TaxID.Country.Empty() {
		return taxCountryToISO(p.TaxID.Country)
	}
	for _, id := range p.Identities {
		if id != nil && !id.Country.Empty() {
			return id.Country
		}
	}
	return ""
}

// taxCountryToISO converts a tax country code into its ISO country code.
func taxCountryToISO(c l10n.TaxCountryCode) l10n.ISOCountryCode {
	if iso := c.Code().ISO(); iso.Validate() == nil {
		return iso
	}
	// The tax code is not a valid ISO country code, try the alternative code in the
	// country definition (e.g. Greece uses EL for tax, GR for ISO).
	if def := l10n.Countries().Code(c.Code()); def != nil && !def.AltCode.Empty() {
		if iso := def.AltCode.ISO(); iso.Validate() == nil {
			return iso
		}
	}
	return ""
}

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
