package it

import "github.com/invopop/gobl/org"

func normalizeParty(party *org.Party) {
	if party == nil || party.TaxID == nil {
		return
	}
	if !party.TaxID.Country.In("IT") {
		return
	}
	// If the party is an individual, move the fiscal code to the identities.
	if party.TaxID.Type == "individual" { //nolint:staticcheck
		id := &org.Identity{
			Key:  IdentityKeyFiscalCode,
			Code: party.TaxID.Code,
		}
		party.TaxID.Code = ""
		party.TaxID.Type = "" //nolint:staticcheck
		party.Identities = org.AddIdentity(party.Identities, id)
	}
}
