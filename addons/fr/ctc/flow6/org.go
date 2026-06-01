package flow6

import (
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// Identity scheme constants used by the Flow 6 normalizers.
const (
	// identitySchemeIDSIREN is the ISO scheme ID for SIREN identities.
	identitySchemeIDSIREN = "0002"
	// identitySchemeIDSIRET is the ISO scheme ID for SIRET identities.
	identitySchemeIDSIRET = "0009"
	// identitySchemeIDPrivate is the ISO scheme ID for the CTC-specific
	// 0224 private ID.
	identitySchemeIDPrivate = "0224"

	// identityKeyPrivateID is the key for private ID identities.
	identityKeyPrivateID cbc.Key = "private-id"

	// inboxSchemeSIREN is the scheme code for SIREN-based addresses.
	inboxSchemeSIREN cbc.Code = "0225"
)

// allowedFlow6IdentitySchemes is the ICD 6523 subset CDAR accepts on
// Flow 6 party identities (STC 0231 is intentionally absent).
var allowedFlow6IdentitySchemes = []cbc.Code{
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

// orgPartyRules validates the integrity of the addon's own party
// extensions: the role code and the identity iso-scheme-id must be
// recognised Flow 6 values. French CTC format/business rules are the
// converter's responsibility — see the package doc.
func orgPartyRules() *rules.Set {
	return rules.For(new(org.Party),
		rules.Field("ext",
			rules.Assert("01", "party ext fr-ctc-flow6-role must be a recognised CDAR RoleCode",
				tax.ExtensionHasValidCode(ExtKeyRole),
			),
		),
		rules.Field("identities",
			rules.Each(
				rules.Field("ext",
					rules.Assert("02", "party identity ext iso-scheme-id must be one of the Flow 6 allowed schemes",
						tax.ExtensionsHasCodes(iso.ExtKeySchemeID, allowedFlow6IdentitySchemes...),
					),
				),
			),
		),
	)
}
