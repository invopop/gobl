package flow6

import (
	"slices"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// allowedRoleCodes is the UNCL 3035 subset that the fr-ctc-role
// extension accepts, kept in sync with the extension definition.
var allowedRoleCodes = []cbc.Code{
	RoleSE, RoleBY, RoleWK, RoleDFH, RoleAB, RoleSR,
	RoleDL, RolePE, RolePR, RoleII, RoleIV,
}

// allowedIdentitySchemes is the ICD 6523 subset CDAR accepts on the
// Flow 6 party identities — SIREN plus the commonly used foreign
// identifier schemes. Parties with identities outside this set should
// not be reported in a Flow 6 CDV.
var allowedIdentitySchemes = []string{
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

func orgPartyRules() *rules.Set {
	return rules.For(new(org.Party),
		rules.Field("ext",
			rules.Assert("01", "fr-ctc-role must be one of the UNCL 3035 subset: SE, BY, WK, DFH, AB, SR, DL, PE, PR, II, IV",
				is.Func("known fr-ctc-role", partyRoleKnown),
			),
		),
		rules.Field("identities",
			rules.Each(
				rules.Field("ext",
					rules.Assert("02", "identity scheme (iso-scheme-id) must be one of the ICD 6523 codes accepted by Flow 6",
						is.Func("scheme in allowed set", partyIdentitySchemeAllowed),
					),
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

func partyIdentitySchemeAllowed(v any) bool {
	ext := extValue(v)
	if ext.IsZero() {
		return true
	}
	scheme := ext.Get(iso.ExtKeySchemeID).String()
	if scheme == "" {
		return true
	}
	return slices.Contains(allowedIdentitySchemes, scheme)
}
