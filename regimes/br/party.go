package br

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

var (
	validStates = []cbc.Code{
		"AC", "AL", "AM", "AP", "BA", "CE", "DF", "ES", "GO",
		"MA", "MG", "MS", "MT", "PA", "PB", "PE", "PI", "PR",
		"RJ", "RN", "RO", "RR", "RS", "SC", "SE", "SP", "TO",
	}
	validPostCode = regexp.MustCompile(`^\d{5}-?\d{3}$`)
)

func validateParty(party *org.Party) error {
	if party == nil {
		return nil
	}

	return validation.ValidateStruct(party,
		validation.Field(&party.Addresses,
			validation.Each(
				validation.By(validateAddress(party)),
			),
			validation.Skip,
		),
	)
}

func validateAddress(party *org.Party) validation.RuleFunc {
	return func(value interface{}) error {
		addr, _ := value.(*org.Address)
		if addr == nil {
			return nil
		}

		if !isBrazilianAddress(party, addr) {
			return nil
		}

		return validation.ValidateStruct(addr,
			validation.Field(&addr.State, validation.In(validStates...)),
			validation.Field(&addr.Code, validation.Match(validPostCode)),
		)
	}
}

func normalizeParty(p *org.Party) {
	if p == nil || p.Ext == nil {
		return
	}

	// migrate old addon extension keys to the regime
	for oldKey, newKey := range map[cbc.Key]cbc.Key{
		"br-nfse-fiscal-incentive": ExtKeyFiscalIncentive,
		"br-nfse-municipality":     ExtKeyMunicipality,
		"br-nfse-simples":          ExtKeySimples,
		"br-nfse-special-regime":   ExtKeySpecialRegime,
	} {
		if value, exists := p.Ext[oldKey]; exists {
			p.Ext[newKey] = value
			delete(p.Ext, oldKey)
		}
	}
}

func isBrazilianAddress(party *org.Party, addr *org.Address) bool {
	if addr.Country != "" {
		return addr.Country == l10n.BR.ISO()
	}
	return party.TaxID != nil && party.TaxID.Country == l10n.BR.Tax()
}
