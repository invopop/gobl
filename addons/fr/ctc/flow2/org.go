package flow2

import (
	"errors"
	"fmt"
	"regexp"
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

// Inbox / identity scheme constants used across Flow 2.
const (
	// inboxSchemeSIREN is the scheme code for SIREN-based addresses
	// (ISO/IEC 6523).
	inboxSchemeSIREN cbc.Code = "0225"
	// identitySchemeIDSIREN is the ISO scheme ID for SIREN identities.
	identitySchemeIDSIREN = "0002"
	// identitySchemeIDSIRET is the ISO scheme ID for SIRET identities.
	identitySchemeIDSIRET = "0009"
	// identitySchemeIDPrivate is the ISO scheme ID for identities
	// requiring alphanumeric format (CTC-specific 0224 private ID).
	identitySchemeIDPrivate = "0224"
	// identityKeyPrivateID is the key for private ID identities.
	identityKeyPrivateID cbc.Key = "private-id"
)

// sirenInboxFormatRegex enforces the alphanumeric + `-+_/` format
// shared by SIREN-scope inboxes and private-id identity codes.
var sirenInboxFormatRegex = regexp.MustCompile(`^[A-Za-z0-9+\-_/]+$`)

func normalizeParty(party *org.Party) {
	if party == nil {
		return
	}
	normalizePartyFromTaxID(party)
	normalizeIdentities(party)
	normalizeInboxes(party)
}

// normalizePartyFromTaxID derives a legal identity from the party's
// TaxID when no matching identity is present.
func normalizePartyFromTaxID(party *org.Party) {
	if party.TaxID == nil {
		return
	}
	country := l10n.Code(party.TaxID.Country)
	code := string(party.TaxID.Code)
	if code == "" || country != l10n.FR {
		return
	}
	ensureSIRENIdentity(party, cbc.Code(sirenFromFrenchTaxID(code, party)))
}

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

// ensureSIRENIdentity appends a SIREN legal identity (ISO scheme 0002)
// when the party does not already carry one.
func ensureSIRENIdentity(party *org.Party, code cbc.Code) {
	if code == "" {
		return
	}
	for _, id := range party.Identities {
		if id != nil && !id.Ext.IsZero() && id.Ext.Get(iso.ExtKeySchemeID).String() == identitySchemeIDSIREN {
			return
		}
	}
	party.Identities = append(party.Identities, &org.Identity{
		Type: fr.IdentityTypeSIREN,
		Code: code,
		Ext: tax.ExtensionsOf(cbc.CodeMap{
			iso.ExtKeySchemeID: cbc.Code(identitySchemeIDSIREN),
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
	if siren != nil && !hasLegalScope {
		siren.Scope = org.IdentityScopeLegal
	}
}

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

// -- Helpers --------------------------------------------------------------

func getPartySIREN(party *org.Party) string {
	if party == nil {
		return ""
	}
	for _, id := range party.Identities {
		if id != nil && (id.Type == fr.IdentityTypeSIREN ||
			(!id.Ext.IsZero() && id.Ext.Get(iso.ExtKeySchemeID) == identitySchemeIDSIREN)) {
			return string(id.Code)
		}
	}
	return ""
}

func isPartyIdentitySTC(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	for _, id := range party.Identities {
		if id != nil && !id.Ext.IsZero() {
			if code := id.Ext.Get(iso.ExtKeySchemeID); code == "0231" {
				return true
			}
		}
	}
	return false
}

func identitiesHasLegalSIREN(val any) bool {
	identities, ok := val.([]*org.Identity)
	if !ok {
		return true
	}
	for _, id := range identities {
		if id != nil && !id.Ext.IsZero() {
			if code := id.Ext.Get(iso.ExtKeySchemeID); code == identitySchemeIDSIREN && id.Scope.Has(org.IdentityScopeLegal) {
				return true
			}
		}
	}
	return false
}

func partyHasSIRENInbox(val any) bool {
	party, ok := val.(*org.Party)
	if !ok || party == nil {
		return true
	}
	siren := getPartySIREN(party)
	if siren == "" {
		return true
	}
	for _, inbox := range party.Inboxes {
		if inbox != nil && inbox.Scheme == inboxSchemeSIREN {
			return strings.HasPrefix(string(inbox.Code), siren)
		}
	}
	return false
}

// -- Rules ----------------------------------------------------------------

func orgPartyRules() *rules.Set {
	return rules.For(new(org.Party),
		rules.Field("identities",
			rules.Assert("01", "SIRET and SIREN must be coherent (BR-FR-09/10)",
				is.Func("SIRET/SIREN coherent", identitiesSIRETSIRENCoherent),
			),
			rules.Assert("02", "identity scheme format invalid (BR-FR-CO-10)",
				is.FuncError("valid scheme format", identitiesSchemeFormatValid),
			),
		),
		rules.Field("inboxes",
			rules.Each(
				rules.Assert("03", "inbox code format invalid",
					is.Func("valid inbox", inboxCodeValid),
				),
			),
		),
	)
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.Func("scheme 0224", identitySchemeIs0224),
			rules.Field("code",
				rules.Assert("01", "must be no more than 100 characters long",
					is.Length(0, 100),
				),
				rules.Assert("02", "must be in a valid format",
					is.Matches(`^[A-Za-z0-9\-\+_/]+$`),
				),
			),
		),
	)
}

func orgInboxRules() *rules.Set {
	return rules.For(new(org.Inbox),
		rules.When(
			is.Func("scheme 0225", inboxSchemeIs0225),
			rules.Field("code",
				rules.Assert("01", "the length must be between 0 and 125",
					is.Length(0, 125),
				),
				rules.Assert("02", "must be in a valid format",
					is.Matches(`^[A-Za-z0-9\-\+_/]+$`),
				),
			),
		),
	)
}

func orgItemRules() *rules.Set {
	return rules.For(new(org.Item),
		rules.Field("meta",
			rules.Assert("01", "meta values cannot be blank (BR-FR-28)",
				is.FuncError("no blank meta", metaNoBlankValues),
			),
		),
	)
}

// -- Validation helpers ---------------------------------------------------

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

func identitySchemeIs0224(val any) bool {
	id, ok := val.(*org.Identity)
	return ok && id != nil && !id.Ext.IsZero() && id.Ext.Get(iso.ExtKeySchemeID) == identitySchemeIDPrivate
}

func inboxSchemeIs0225(val any) bool {
	inbox, ok := val.(*org.Inbox)
	return ok && inbox != nil && inbox.Scheme == inboxSchemeSIREN
}

func metaNoBlankValues(val any) error {
	meta, ok := val.(cbc.Meta)
	if !ok || len(meta) == 0 {
		return nil
	}
	for key, v := range meta {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("%s: value cannot be blank (BR-FR-28)", key)
		}
	}
	return nil
}
