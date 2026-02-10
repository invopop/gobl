package ctc

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
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
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

	// Attachment for format category for BR-FR-18
	attachmentFormat = "LISIBLE"
)

// Allowed attachment description values for French CTC (BR-FR-17)
var allowedAttachmentDescriptions = []cbc.Code{
	"RIB",                        // Bank account details (Relevé d'Identité Bancaire)
	"LISIBLE",                    // Human-readable representation of the invoice
	"FEUILLE_DE_STYLE",           // Style sheet for document presentation
	"PJA",                        // Additional supporting document (Pièce Jointe Additionnelle)
	"BON_LIVRAISON",              // Delivery note
	"BON_COMMANDE",               // Purchase order
	"DOCUMENT_ANNEXE",            // Annex document
	"BORDEREAU_SUIVI",            // Follow-up form
	"BORDEREAU_SUIVI_VALIDATION", // Validated follow-up form
	"ETAT_ACOMPTE",               // Payment status statement
	"FACTURE_PAIEMENT_DIRECT",    // Direct payment invoice
	"RECAPITULATIF_COTRAITANCE",  // Co-contracting summary
}

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
			// Create SIREN identity from first 9 digits of SIRET
			sirenCode := siretCode[:9]
			siren = &org.Identity{
				Type: fr.IdentityTypeSIREN,
				Code: cbc.Code(sirenCode),
				// Note: ISO scheme ID will be set by EN16931 addon
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
		if id.Ext == nil {
			id.Ext = make(tax.Extensions)
		}
		id.Ext.Set(iso.ExtKeySchemeID, identitySchemeIDPrivate)
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

// validateIdentity validates a single identity
func validateIdentity(id *org.Identity) error {
	if id == nil {
		return nil
	}

	return validation.ValidateStruct(&id,
		validation.Field(&id.Code,
			validation.When(
				id.Ext != nil && id.Ext.Get(iso.ExtKeySchemeID) == identitySchemeIDPrivate,
				validation.Length(0, 100),
				validation.Match(sirenInboxFormatRegex),
			),
		),
	)
}

func validateParty(party *org.Party) error {
	if party == nil {
		return nil
	}

	return validation.ValidateStruct(party,
		validation.Field(&party.Identities,
			validation.By(validateSIRETSIRENCoherence),
			validation.By(validateIdentitySchemeFormat),
			validation.Skip,
		),
		validation.Field(&party.Inboxes,
			validation.Each(validation.By(validateInbox)),
			validation.Skip,
		),
	)
}

// validateIdentitySchemeFormat validates format for identities with specific ISO scheme IDs
func validateIdentitySchemeFormat(value any) error {
	identities, ok := value.([]*org.Identity)
	if !ok || len(identities) == 0 {
		return nil
	}

	schemes := make(map[cbc.Code]bool)

	for _, id := range identities {
		if id == nil {
			continue
		}

		schemeID := id.Ext[iso.ExtKeySchemeID]

		if schemeID == cbc.CodeEmpty {
			// All identites must have an ISO scheme ID
			return errors.New("all identities must have an ISO scheme ID defined in extensions BR-FR-CO-10")
		}

		if schemes[schemeID] {
			return fmt.Errorf("duplicate identities with ISO scheme ID '%s' are not allowed (BR-FR-CO-10)", schemeID)
		}

		// Check if identity has ISO scheme ID 0224 (private-id)
		if schemeID == identitySchemeIDPrivate {
			code := string(id.Code)
			if code == "" {
				continue
			}

			// Validate length: max 100 characters
			if len(code) > 100 {
				return errors.New("identity with ISO scheme ID 0224 (private-id) must not exceed 100 characters (BR-FR-26)")
			}

			// Validate format: alphanumeric plus +, -, _, /
			if !sirenInboxFormatRegex.MatchString(code) {
				return errors.New("identity with ISO scheme ID 0224 (private-id) must contain only alphanumeric characters and +, -, _, / (BR-FR-24)")
			}
		}

		schemes[schemeID] = true
	}

	return nil
}

// validateSIRETSIRENCoherence validates that SIRET and SIREN identities are coherent (BR-FR-09/10)
// If both are present, the first 9 digits of SIRET must match the SIREN
func validateSIRETSIRENCoherence(value any) error {
	identities, ok := value.([]*org.Identity)
	if !ok || len(identities) == 0 {
		return nil
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

	// If both SIRET and SIREN are present, check coherence
	if siret != nil && siren != nil {
		siretCode := string(siret.Code)
		sirenCode := string(siren.Code)

		if len(siretCode) == 14 && len(sirenCode) == 9 {
			if !strings.HasPrefix(siretCode, sirenCode) {
				return errors.New("SIRET must start with the same 9 digits as SIREN (BR-FR-09/10)")
			}
		}
	}

	return nil
}

// validateInbox validates a single inbox with the expected SIREN
func validateInbox(value any) error {
	inbox, ok := value.(*org.Inbox)
	if !ok || inbox == nil {
		return nil
	}

	return validation.ValidateStruct(&inbox,
		validation.Field(&inbox.Code,
			validation.When(
				inbox.Scheme == inboxSchemeSIREN,
				validation.Length(0, 125),
				validation.Match(sirenInboxFormatRegex),
			),
			validation.Skip,
		),
	)
}

// validateOrgAttachments validates attachment descriptions and uniqueness constraints (BR-FR-17/18)
// BR-FR-17: Description must be in list of valid descriptions
// BR-FR-18: Only one attachment can have the "LISIBLE" description per invoice
//
// When creating attachments via API, use a value from the allowed list in the description field.
func validateOrgAttachments(attachments []*org.Attachment) error {
	if len(attachments) == 0 {
		return nil
	}

	lisibleCount := 0

	for _, att := range attachments {
		if att == nil {
			continue
		}

		// BR-FR-17: Validate that description is one of the allowed values
		if att.Description != "" {
			desc := cbc.Code(strings.TrimSpace(att.Description))
			if !slices.Contains(allowedAttachmentDescriptions, desc) {
				allowed := make([]string, len(allowedAttachmentDescriptions))
				for i, code := range allowedAttachmentDescriptions {
					allowed[i] = string(code)
				}
				return fmt.Errorf("attachment description '%s' is not allowed, must be one of: %v (BR-FR-17)", desc, allowed)
			}

			// BR-FR-18: Count LISIBLE attachments
			if desc == attachmentFormat {
				lisibleCount++
			}
		}
	}

	// BR-FR-18: Only one LISIBLE attachment allowed
	if lisibleCount > 1 {
		return errors.New("only one attachment with description 'LISIBLE' is allowed per invoice (BR-FR-18)")
	}

	return nil
}
