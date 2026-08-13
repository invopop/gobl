package bill

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

// supplierTaxCountry determines the tax country for the invoice based on the supplier tax
// identity.
func partyTaxCountry(party *org.Party) l10n.TaxCountryCode {
	if party == nil {
		return l10n.CodeEmpty.Tax()
	}
	if party.TaxID == nil {
		return l10n.CodeEmpty.Tax()
	}
	return party.TaxID.Country
}
