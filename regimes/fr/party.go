package fr

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// normalizeParty normalizes the tax identity code for the party and adds any
// identities that are present.
func normalizeParty(party *org.Party) {
	if party == nil {
		return
	}

	if party.TaxID == nil || party.TaxID.Code == "" || party.TaxID.Country != "FR" {
		return
	}

	tax.NormalizeIdentity(party.TaxID)

	str := party.TaxID.Code.String()
	siret := false

	// First check if we have a SIRET and keep SIREN
	if len(str) == 14 {
		str = str[:9]
		siret = true
	}

	// Now check if we have a SIREN
	if len(str) == 9 {
		if err := validateSIRENTaxCode(party.TaxID.Code); err != nil {
			return
		}
		chk := calculateVATCheckDigit(str)
		// Add SIREN or SIRET identity
		identity := IdentityTypeSIREN
		if siret {
			identity = IdentityTypeSIRET
		}
		party.Identities = org.AddIdentity(party.Identities, &org.Identity{
			Type: identity,
			Code: party.TaxID.Code,
		})

		party.TaxID.Code = cbc.Code(fmt.Sprintf("%s%s", chk, str))
	}
}
