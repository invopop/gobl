package ctc

import (
	"regexp"
	"slices"
	"strings"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
)

// Inbox / identity scheme constants used across the addon.
const (
	// inboxSchemeSIREN is the scheme code for SIREN-based addresses (ISO/IEC 6523).
	inboxSchemeSIREN cbc.Code = "0225"

	// identitySchemeIDSIREN is the ISO scheme ID for SIREN identities.
	identitySchemeIDSIREN = "0002"
	// identitySchemeIDSIRET is the ISO scheme ID for SIRET identities.
	identitySchemeIDSIRET = "0009"
	// identitySchemeIDEUVAT is the ISO scheme ID for EU (non-French) intra-community VAT.
	identitySchemeIDEUVAT = "0223"
	// identitySchemeIDNonEU is the ISO scheme ID for non-EU party identifiers.
	identitySchemeIDNonEU = "0227"
	// identitySchemeIDRIDET is the ISO scheme ID for New Caledonia RIDET.
	identitySchemeIDRIDET = "0228"
	// identitySchemeIDTAHITI is the ISO scheme ID for French Polynesia TAHITI.
	identitySchemeIDTAHITI = "0229"

	// identityKeyPrivateID is the key for private ID identities.
	identityKeyPrivateID cbc.Key = "private-id"
	// identitySchemeIDPrivate is the ISO scheme ID for identities requiring
	// alphanumeric format (CTC-specific 0224 private ID).
	identitySchemeIDPrivate = "0224"
)

// sirenInboxFormatRegex enforces the alphanumeric + `-+_/` format
// shared by SIREN-scope inboxes and private-id identity codes.
var sirenInboxFormatRegex = regexp.MustCompile(`^[A-Za-z0-9+\-_/]+$`)

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

// allowedFlow6IdentitySchemes is the ICD 6523 subset CDAR accepts on
// Flow 6 (CDV lifecycle) party identities. STC (0231 — assujetti
// unique) is intentionally absent: it is a Flow 2 invoice concept
// and must not appear on a CDV. Used by bill_status rule 22; not
// applied addon-wide because Flow 10 cross-border B2B may legitimately
// carry foreign identifier schemes outside this list.
var allowedFlow6IdentitySchemes = []string{
	"0002", // SIREN
	"0009", // SIRET
	"0223", // EU VAT
	"0224", // Private ID
	"0226", // European VAT
	"0227", // Non-EU
	"0228", // RIDET (New Caledonia)
	"0229", // TAHITI (French Polynesia)
	"0238", // Peppol participant ID
}

// allowedRoleCodes is the UNCL 3035 subset that the fr-ctc-role
// extension accepts.
var allowedRoleCodes = []cbc.Code{
	RoleSE, RoleBY, RoleWK, RoleDFH, RoleAB, RoleSR,
	RoleDL, RolePE, RolePR, RoleII, RoleIV,
}

func normalizeParty(party *org.Party) {
	if party == nil {
		return
	}

	// Derive identities from TaxID for Flow 10 reporting (SIREN for
	// French, EU-VAT identity for other EU countries). The Flow 2
	// SIREN-from-SIRET path is handled below in normalizeIdentities.
	normalizePartyFromTaxID(party)

	// Normalize identities (SIREN-from-SIRET, legal scope).
	normalizeIdentities(party)

	// Normalize inboxes (peppol key on SIREN inbox).
	normalizeInboxes(party)
}

// normalizePartyFromTaxID attempts to derive a legal identity from the
// party's TaxID when no matching identity is present. Mirrors the
// pre-merge flow10 behaviour: French TaxID → SIREN identity; other-EU
// TaxID → EU-VAT identity.
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

// normalizeIdentities handles SIRET → SIREN derivation and legal-scope
// assignment.
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

	// BR-FR-09/10: Generate SIREN from SIRET if needed.
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

	// Set SIREN scope to legal if no other identity has legal scope.
	if siren != nil && !hasLegalScope {
		siren.Scope = org.IdentityScopeLegal
	}
}

// normalizeIdentity handles per-identity normalization: maps the
// private-id key to scheme 0224 and the SIREN/SIRET identity types
// to schemes 0002/0009 respectively. The fr-ctc addon owns this
// mapping so it works even when eu-en16931 is not declared (Flow 6
// or standalone Flow 10 callers).
func normalizeIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	if id.Key == identityKeyPrivateID {
		id.Ext = id.Ext.Set(iso.ExtKeySchemeID, identitySchemeIDPrivate)
	}
	if id.Type == fr.IdentityTypeSIREN && id.Ext.Get(iso.ExtKeySchemeID) == "" {
		id.Ext = id.Ext.Set(iso.ExtKeySchemeID, identitySchemeIDSIREN)
	}
	if id.Type == fr.IdentityTypeSIRET && id.Ext.Get(iso.ExtKeySchemeID) == "" {
		id.Ext = id.Ext.Set(iso.ExtKeySchemeID, identitySchemeIDSIRET)
	}
}

// normalizeInboxes flags the SIREN-scope inbox with the peppol key
// when no other inbox already carries it.
func normalizeInboxes(party *org.Party) {
	if party == nil || len(party.Inboxes) == 0 {
		return
	}

	hasPeppol := false
	var sirenInbox *org.Inbox
	for _, inbox := range party.Inboxes {
		if inbox == nil {
			continue
		}
		if inbox.Key == org.InboxKeyPeppol {
			hasPeppol = true
		}
		if inbox.Scheme == inboxSchemeSIREN {
			sirenInbox = inbox
		}
	}
	if !hasPeppol && sirenInbox != nil {
		sirenInbox.Key = org.InboxKeyPeppol
	}
}

// -- Predicates used across addon files ---------------------------------

// partyLegalSchemeID returns the ICD 6523 scheme ID of the identity the
// party presents as its legal identifier. Prefers an identity scoped as
// "legal"; failing that, the first identity that declares a known
// scheme ID.
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

// partyHasSIREN reports whether the party carries a SIREN-scheme
// (0002) identity.
func partyHasSIREN(v any) bool {
	party, ok := v.(*org.Party)
	if !ok || party == nil {
		return false
	}
	return partyCarriesSIREN(party)
}

func partyCarriesSIREN(party *org.Party) bool {
	if party == nil {
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

// partyIsFrench returns true when the party is identifiable as French —
// either it carries a SIREN identity or its TaxID is registered under
// the French regime. Used by the invoice-rule dispatcher to pick the
// Flow 2 ruleset.
func partyIsFrench(party *org.Party) bool {
	if party == nil {
		return false
	}
	if partyCarriesSIREN(party) {
		return true
	}
	if party.TaxID != nil && l10n.Code(party.TaxID.Country) == l10n.FR {
		return true
	}
	return false
}
