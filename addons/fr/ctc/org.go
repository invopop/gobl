package ctc

import (
	"regexp"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
)

// Inbox validation patterns
var sirenInboxFormatRegex = regexp.MustCompile(`^[A-Za-z0-9+\-_/]+$`)

const (
	// inboxSchemeSIREN is the scheme code for SIREN-based addresses (ISO/IEC 6523)
	inboxSchemeSIREN cbc.Code = "0225"
	// identitySchemeIDSIREN is the ISO scheme ID for SIREN identities
	identitySchemeIDSIREN = "0002"

	// identityKeyPrivateID is the key for private ID identities
	identityKeyPrivateID cbc.Key = "private-id"
	// identitySchemeIDPrivate is the ISO scheme ID for identities requiring alphanumeric format
	identitySchemeIDPrivate = "0224"
)

func normalizeParty(party *org.Party) {
	if party == nil {
		return
	}

	// Normalize identities
	normalizeIdentities(party)

	// Normalize inboxes
	normalizeInboxes(party)
}

// normalizeIdentities handles all identity-related normalizations
func normalizeIdentities(party *org.Party) {
	if party == nil || len(party.Identities) == 0 {
		return
	}

	var siret, siren *org.Identity
	hasLegalScope := false

	// First pass: normalize each identity and collect information
	for _, id := range party.Identities {
		if id == nil {
			continue
		}

		// Normalize individual identity (sets type from scheme ID, private-id scheme)
		normalizeIdentity(id)

		// Check for SIRET and SIREN (after normalization may have set the type)
		if id.Type == fr.IdentityTypeSIRET {
			siret = id
		}
		if id.Type == fr.IdentityTypeSIREN {
			siren = id
		}

		// Check for legal scope
		if id.Scope == org.IdentityScopeLegal {
			hasLegalScope = true
		}
	}

	// BR-FR-09/10: Generate SIREN from SIRET if needed
	if siret != nil && siren == nil {
		siretCode := string(siret.Code)
		if len(siretCode) == 14 {
			// Create SIREN identity from first 9 digits of SIRET.
			// We must set the ISO scheme ID here because EN16931 has already
			// processed the identities slice before FR CTC runs normalizeParty.
			sirenCode := siretCode[:9]
			siren = &org.Identity{
				Type: fr.IdentityTypeSIREN,
				Code: cbc.Code(sirenCode),
				Ext: tax.ExtensionsOf(tax.ExtMap{
					iso.ExtKeySchemeID: identitySchemeIDSIREN,
				}),
			}
			party.Identities = append(party.Identities, siren)
		}
	}

	// Set SIREN scope to legal if no other identity has legal scope
	if siren != nil && !hasLegalScope {
		siren.Scope = org.IdentityScopeLegal
	}
}

// normalizeIdentity handles normalization for a single identity
func normalizeIdentity(id *org.Identity) {
	if id == nil {
		return
	}

	// Set ISO scheme ID 0224 for private-id key (CTC-specific)
	if id.Key == identityKeyPrivateID {
		id.Ext = id.Ext.Set(iso.ExtKeySchemeID, identitySchemeIDPrivate)
	}
	// Note: Type ↔ ISO scheme ID mapping for SIREN/SIRET is handled by EN16931 addon
}

// normalizeInboxes handles all inbox-related normalizations
func normalizeInboxes(party *org.Party) {
	if party == nil || len(party.Inboxes) == 0 {
		return
	}

	// Check if any inbox already has the peppol key
	hasPeppol := false
	var sirenInbox *org.Inbox
	for _, inbox := range party.Inboxes {
		if inbox == nil {
			continue
		}
		if inbox.Key == "peppol" {
			hasPeppol = true
		}
		if inbox.Scheme == inboxSchemeSIREN {
			sirenInbox = inbox
		}
	}

	// If no inbox has peppol key and we have a SIREN inbox, set it
	if !hasPeppol && sirenInbox != nil {
		sirenInbox.Key = "peppol"
	}
}
