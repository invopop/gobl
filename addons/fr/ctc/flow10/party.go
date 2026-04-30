package flow10

import (
	"slices"
	"strings"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
)

// ICD 6523 scheme IDs accepted for Flow 10 party identification (G2.19).
const (
	schemeIDSIREN  = "0002" // French SIREN (9 digits)
	schemeIDEUVAT  = "0223" // EU (non-French) intra-community VAT ID
	schemeIDNonEU  = "0227" // Outside EU: country code + first 16 chars of name
	schemeIDRIDET  = "0228" // New Caledonia RIDET
	schemeIDTAHITI = "0229" // French Polynesia TAHITI
)

// allowedPartySchemeIDs lists the scheme IDs permitted for the legal
// identity of a Flow 10 B2B party (supplier or customer), per G2.19.
var allowedPartySchemeIDs = []string{
	schemeIDSIREN,
	schemeIDEUVAT,
	schemeIDNonEU,
	schemeIDRIDET,
	schemeIDTAHITI,
}

// schemeIDsRequiringVAT are the scheme IDs for which party.TaxID must also
// be present (G2.33): SIREN (French) and EU non-French VAT identifiers.
var schemeIDsRequiringVAT = []string{
	schemeIDSIREN,
	schemeIDEUVAT,
}

// normalizeParty attempts to derive the Flow 10 party identity from
// information already present on the party (TaxID, SIRET) so that the
// downstream rules can succeed without the caller having to hand-craft
// the ICD 6523 identity entry.
func normalizeParty(party *org.Party) {
	if party == nil || party.TaxID == nil {
		return
	}

	country := l10n.Code(party.TaxID.Country)
	code := string(party.TaxID.Code)
	if code == "" {
		return
	}

	switch {
	case country == l10n.FR:
		ensureIdentity(party, fr.IdentityTypeSIREN, cbc.Code(sirenFromFrenchTaxID(code, party)), schemeIDSIREN)
	case isEUNonFrance(country):
		ensureIdentity(party, "", cbc.Code(country.String()+code), schemeIDEUVAT)
	}
}

// sirenFromFrenchTaxID extracts the 9-digit SIREN from a French TaxID. The
// French VAT format is FR + 2 check digits + 9 digit SIREN; if the TaxID
// has already been stripped to just the 9 digits we return it as-is. If a
// SIRET identity is already present we prefer its first 9 digits.
func sirenFromFrenchTaxID(taxCode string, party *org.Party) string {
	for _, id := range party.Identities {
		if id != nil && id.Type == fr.IdentityTypeSIRET {
			s := string(id.Code)
			if len(s) == 14 {
				return s[:9]
			}
		}
	}
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, taxCode)
	if len(digits) >= 9 {
		return digits[len(digits)-9:]
	}
	return digits
}

// ensureIdentity adds an identity matching the given scheme ID if none is
// already present; identities that already declare the scheme (via
// iso.ExtKeySchemeID) are left untouched so user-supplied data wins.
func ensureIdentity(party *org.Party, typ cbc.Code, code cbc.Code, schemeID string) {
	if code == "" {
		return
	}
	for _, id := range party.Identities {
		if id != nil && !id.Ext.IsZero() && id.Ext.Get(iso.ExtKeySchemeID).String() == schemeID {
			return
		}
	}
	party.Identities = append(party.Identities, &org.Identity{
		Type: typ,
		Code: code,
		Ext: tax.ExtensionsOf(tax.ExtMap{
			iso.ExtKeySchemeID: cbc.Code(schemeID),
		}),
		Scope: org.IdentityScopeLegal,
	})
}

// partyLegalSchemeID returns the ICD 6523 scheme ID of the identity the
// party presents as its legal identifier for Flow 10. It prefers an
// identity scoped as "legal"; failing that, the first identity that
// declares a known Flow 10 scheme ID.
func partyLegalSchemeID(party *org.Party) string {
	if party == nil {
		return ""
	}
	var fallback string
	for _, id := range party.Identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		scheme := id.Ext.Get(iso.ExtKeySchemeID).String()
		if scheme == "" {
			continue
		}
		if id.Scope == org.IdentityScopeLegal {
			return scheme
		}
		if fallback == "" && slices.Contains(allowedPartySchemeIDs, scheme) {
			fallback = scheme
		}
	}
	return fallback
}

func isEUNonFrance(c l10n.Code) bool {
	if c == l10n.FR || c == "" {
		return false
	}
	eu := l10n.Union(l10n.EU)
	return eu != nil && eu.HasMember(c)
}
