package flow2

import (
	"errors"
	"fmt"
	"strings"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

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

// --- Helper functions ---

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
