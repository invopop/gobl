package flow10

import (
	"slices"
	"strings"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Identity scheme constants used by Flow 10 reporting.
const (
	// identitySchemeIDSIREN is the ISO scheme ID for SIREN identities.
	identitySchemeIDSIREN = "0002"
	// identitySchemeIDSIRET is the ISO scheme ID for SIRET identities.
	identitySchemeIDSIRET = "0009"
	// identitySchemeIDEUVAT is the ISO scheme ID for EU (non-French)
	// intra-community VAT.
	identitySchemeIDEUVAT = "0223"
	// identitySchemeIDNonEU is the ISO scheme ID for non-EU party
	// identifiers.
	identitySchemeIDNonEU = "0227"
	// identitySchemeIDRIDET is the ISO scheme ID for New Caledonia
	// RIDET.
	identitySchemeIDRIDET = "0228"
	// identitySchemeIDTAHITI is the ISO scheme ID for French Polynesia
	// TAHITI.
	identitySchemeIDTAHITI = "0229"
)

// allowedPartySchemeIDs lists the scheme IDs permitted for the legal
// identity of a Flow 10 B2B party (supplier or customer), per G2.19.
var allowedPartySchemeIDs = []string{
	identitySchemeIDSIREN,
	identitySchemeIDEUVAT,
	identitySchemeIDNonEU,
	identitySchemeIDRIDET,
	identitySchemeIDTAHITI,
}

// schemeIDsRequiringVAT are the scheme IDs for which party.TaxID must
// also be present (G2.33): SIREN (French) and EU non-French VAT.
var schemeIDsRequiringVAT = []string{
	identitySchemeIDSIREN,
	identitySchemeIDEUVAT,
}

func normalizeParty(party *org.Party) {
	if party == nil {
		return
	}
	normalizePartyFromTaxID(party)
	normalizeIdentities(party)
}

// normalizePartyFromTaxID derives a legal identity from the party's
// TaxID when no matching identity is present. French TaxID → SIREN
// identity; other-EU TaxID → EU-VAT identity.
func normalizePartyFromTaxID(party *org.Party) {
	if party.TaxID == nil {
		return
	}
	country := l10n.Code(party.TaxID.Country)
	code := string(party.TaxID.Code)
	if code == "" {
		return
	}
	switch {
	case country == l10n.FR:
		ensureIdentity(party, fr.IdentityTypeSIREN, cbc.Code(sirenFromFrenchTaxID(code, party)), identitySchemeIDSIREN)
	case isEUNonFrance(country):
		ensureIdentity(party, "", cbc.Code(country.String()+code), identitySchemeIDEUVAT)
	}
}

// sirenFromFrenchTaxID extracts the 9-digit SIREN from a French TaxID.
// Prefers the first 9 digits of a present SIRET identity, otherwise
// falls back to the last 9 digits of the TaxID code.
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

// ensureIdentity adds an identity matching the given scheme ID if none
// is already present.
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
		Ext: tax.ExtensionsOf(cbc.CodeMap{
			iso.ExtKeySchemeID: cbc.Code(schemeID),
		}),
		Scope: org.IdentityScopeLegal,
	})
}

func normalizeIdentities(party *org.Party) {
	if party == nil || len(party.Identities) == 0 {
		return
	}
	var siret, siren *org.Identity
	hasLegalScope := false
	for _, id := range party.Identities {
		if id == nil {
			continue
		}
		normalizeIdentity(id)
		if id.Type == fr.IdentityTypeSIRET {
			siret = id
		}
		if id.Type == fr.IdentityTypeSIREN {
			siren = id
		}
		if id.Scope == org.IdentityScopeLegal {
			hasLegalScope = true
		}
	}
	// Generate SIREN from SIRET if needed.
	if siret != nil && siren == nil {
		siretCode := string(siret.Code)
		if len(siretCode) == 14 {
			sirenCode := siretCode[:9]
			siren = &org.Identity{
				Type: fr.IdentityTypeSIREN,
				Code: cbc.Code(sirenCode),
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					iso.ExtKeySchemeID: identitySchemeIDSIREN,
				}),
			}
			party.Identities = append(party.Identities, siren)
		}
	}
	if siren != nil && !hasLegalScope {
		siren.Scope = org.IdentityScopeLegal
	}
}

func normalizeIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	if id.Type == fr.IdentityTypeSIREN && id.Ext.Get(iso.ExtKeySchemeID) == "" {
		id.Ext = id.Ext.Set(iso.ExtKeySchemeID, identitySchemeIDSIREN)
	}
	if id.Type == fr.IdentityTypeSIRET && id.Ext.Get(iso.ExtKeySchemeID) == "" {
		id.Ext = id.Ext.Set(iso.ExtKeySchemeID, identitySchemeIDSIRET)
	}
}

func isEUNonFrance(c l10n.Code) bool {
	if c == l10n.FR || c == "" {
		return false
	}
	eu := l10n.Union(l10n.EU)
	return eu != nil && eu.HasMember(c)
}

// partyLegalSchemeID returns the ICD 6523 scheme ID of the identity
// the party presents as its legal identifier.
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

// partyHasSIREN reports whether the party carries a SIREN-scheme
// (0002) identity.
func partyHasSIREN(v any) bool {
	party, ok := v.(*org.Party)
	if !ok || party == nil {
		return false
	}
	for _, id := range party.Identities {
		if id == nil {
			continue
		}
		if id.Type == fr.IdentityTypeSIREN {
			return true
		}
		if !id.Ext.IsZero() && id.Ext.Get(iso.ExtKeySchemeID).String() == identitySchemeIDSIREN {
			return true
		}
	}
	return false
}

func partyHasAllowedLegalScheme(v any) bool {
	party, ok := v.(*org.Party)
	if !ok || party == nil {
		return false
	}
	return slices.Contains(allowedPartySchemeIDs, partyLegalSchemeID(party))
}

func partyHasTaxIDWhenRequired(v any) bool {
	party, ok := v.(*org.Party)
	if !ok || party == nil {
		return true
	}
	scheme := partyLegalSchemeID(party)
	if !slices.Contains(schemeIDsRequiringVAT, scheme) {
		return true
	}
	return party.TaxID != nil && party.TaxID.Code != ""
}

func partyHasVATCode(p *org.Party) bool {
	return p != nil && p.TaxID != nil && p.TaxID.Code != ""
}

func orgPartyRules() *rules.Set {
	return rules.For(new(org.Party),
		rules.Field("identities",
			rules.Assert("01", "party identities must not duplicate iso-scheme-id values (BR-FR-CO-10)",
				is.Func("unique iso-scheme-id", identitiesSchemesUnique),
			),
			rules.Each(
				rules.Field("ext",
					rules.Assert("02", "party identity ext iso-scheme-id is required (BR-FR-CO-10)",
						tax.ExtensionsRequire(iso.ExtKeySchemeID),
					),
				),
			),
		),
	)
}

// identitiesSchemesUnique reports whether the slice contains at most
// one identity per iso-scheme-id value. Empty extensions are ignored
// (the per-identity rule covers them).
func identitiesSchemesUnique(val any) bool {
	identities, ok := val.([]*org.Identity)
	if !ok || len(identities) == 0 {
		return true
	}
	seen := make(map[cbc.Code]bool, len(identities))
	for _, id := range identities {
		if id == nil {
			continue
		}
		schemeID := id.Ext.Get(iso.ExtKeySchemeID)
		if schemeID == cbc.CodeEmpty {
			continue
		}
		if seen[schemeID] {
			return false
		}
		seen[schemeID] = true
	}
	return true
}
