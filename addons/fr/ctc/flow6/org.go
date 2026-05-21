package flow6

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// Identity scheme constants used by Flow 6.
const (
	// identitySchemeIDSIREN is the ISO scheme ID for SIREN identities.
	identitySchemeIDSIREN = "0002"
	// identitySchemeIDSIRET is the ISO scheme ID for SIRET identities.
	identitySchemeIDSIRET = "0009"
	// identitySchemeIDPrivate is the ISO scheme ID for identities
	// requiring alphanumeric format (CTC-specific 0224 private ID).
	identitySchemeIDPrivate = "0224"

	// identityKeyPrivateID is the key for private ID identities.
	identityKeyPrivateID cbc.Key = "private-id"

	// inboxSchemeSIREN is the scheme code for SIREN-based addresses
	// (ISO/IEC 6523).
	inboxSchemeSIREN cbc.Code = "0225"
)

// sirenInboxFormatRegex enforces the alphanumeric + `-+_/` format
// shared by SIREN-scope inboxes and private-id identity codes.
var sirenInboxFormatRegex = regexp.MustCompile(`^[A-Za-z0-9+\-_/]+$`)

// allowedFlow6IdentitySchemes is the ICD 6523 subset CDAR accepts on
// Flow 6 (CDV lifecycle) party identities. STC (0231 — assujetti
// unique) is intentionally absent: it is a Flow 2 invoice concept and
// must not appear on a CDV.
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

// allowedRoleCodes is the UNCL 3035 subset that the fr-ctc-flow6-role
// extension accepts.
var allowedRoleCodes = []cbc.Code{
	RoleSE, RoleBY, RoleWK, RoleDFH, RoleAB, RoleSR,
	RoleDL, RolePE, RolePR, RoleII, RoleIV,
}

// normalizeParty handles the per-party normalisation Flow 6 requires:
// identity scheme tagging (SIREN/SIRET → 0002/0009) and Peppol-key
// flag on SIREN-scoped inbox.
func normalizeParty(party *org.Party) {
	if party == nil {
		return
	}
	for _, id := range party.Identities {
		normalizeIdentity(id)
	}
	normalizeInboxes(party)
}

// normalizeIdentity tags SIREN/SIRET identities with the ISO scheme
// extension when missing and maps the private-id key to scheme 0224.
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

func orgPartyRules() *rules.Set {
	return rules.For(new(org.Party),
		rules.Field("ext",
			rules.Assert("01", "fr-ctc-flow6-role must be one of the UNCL 3035 subset: SE, BY, WK, DFH, AB, SR, DL, PE, PR, II, IV",
				is.Func("known fr-ctc-flow6-role", partyRoleKnown),
			),
		),
		rules.Field("identities",
			rules.Assert("02", "SIRET and SIREN must be coherent (BR-FR-09/10)",
				is.Func("SIRET/SIREN coherent", identitiesSIRETSIRENCoherent),
			),
			rules.Assert("03", "identity scheme format invalid (BR-FR-CO-10)",
				is.FuncError("valid scheme format", identitiesSchemeFormatValid),
			),
		),
		rules.Field("inboxes",
			rules.Each(
				rules.Assert("05", "inbox code format invalid",
					is.Func("valid inbox", inboxCodeValid),
				),
			),
		),
	)
}

func partyRoleKnown(v any) bool {
	ext := extValue(v)
	role := ext.Get(ExtKeyRole)
	if role == "" {
		return true
	}
	return slices.Contains(allowedRoleCodes, role)
}

func identitiesSIRETSIRENCoherent(val any) bool {
	identities, ok := val.([]*org.Identity)
	if !ok || len(identities) == 0 {
		return true
	}
	var siret, siren *org.Identity
	for _, id := range identities {
		if id == nil {
			continue
		}
		if id.Type == fr.IdentityTypeSIRET {
			siret = id
		}
		if id.Type == fr.IdentityTypeSIREN {
			siren = id
		}
	}
	if siret != nil && siren != nil {
		siretCode := string(siret.Code)
		sirenCode := string(siren.Code)
		if len(siretCode) == 14 && len(sirenCode) == 9 {
			if !strings.HasPrefix(siretCode, sirenCode) {
				return false
			}
		}
	}
	return true
}

func identitiesSchemeFormatValid(val any) error {
	identities, ok := val.([]*org.Identity)
	if !ok || len(identities) == 0 {
		return nil
	}
	schemes := make(map[cbc.Code]bool)
	for _, id := range identities {
		if id == nil {
			continue
		}
		schemeID := id.Ext.Get(iso.ExtKeySchemeID)
		if schemeID == cbc.CodeEmpty {
			return errors.New("all identities must have an ISO scheme ID defined in extensions BR-FR-CO-10")
		}
		if schemes[schemeID] {
			return fmt.Errorf("duplicate identities with ISO scheme ID '%s' are not allowed (BR-FR-CO-10)", schemeID)
		}
		if schemeID == identitySchemeIDPrivate {
			code := string(id.Code)
			if code == "" {
				schemes[schemeID] = true
				continue
			}
			if len(code) > 100 {
				return errors.New("identity with ISO scheme ID 0224 (private-id) must not exceed 100 characters (BR-FR-26)")
			}
			if !sirenInboxFormatRegex.MatchString(code) {
				return errors.New("identity with ISO scheme ID 0224 (private-id) must contain only alphanumeric characters and +, -, _, / (BR-FR-24)")
			}
		}
		schemes[schemeID] = true
	}
	return nil
}

func inboxCodeValid(val any) bool {
	inbox, ok := val.(*org.Inbox)
	if !ok || inbox == nil {
		return true
	}
	if inbox.Scheme != inboxSchemeSIREN {
		return true
	}
	code := string(inbox.Code)
	if code == "" {
		return true
	}
	if len(code) > 125 {
		return false
	}
	return sirenInboxFormatRegex.MatchString(code)
}
