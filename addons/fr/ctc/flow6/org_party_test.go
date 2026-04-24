package flow6

import (
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPartyUnknownRoleRejected(t *testing.T) {
	p := &org.Party{
		Name: "Agent",
		Ext:  tax.ExtensionsOf(tax.ExtMap{ExtKeyRole: "XXX"}),
	}
	err := rules.Validate(p, addonContext())
	assert.ErrorContains(t, err, "UNCL 3035")
}

func TestPartyKnownRoleAccepted(t *testing.T) {
	p := &org.Party{
		Name: "Platform",
		Ext:  tax.ExtensionsOf(tax.ExtMap{ExtKeyRole: RoleWK}),
	}
	assert.NoError(t, rules.Validate(p, addonContext()))
}

func TestPartyUnknownIdentitySchemeRejected(t *testing.T) {
	p := &org.Party{
		Name: "Agent",
		Identities: []*org.Identity{{
			Code: "X",
			Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "9999"}),
		}},
	}
	err := rules.Validate(p, addonContext())
	assert.ErrorContains(t, err, "ICD 6523")
}

func TestPartyIdentityWithoutSchemeAccepted(t *testing.T) {
	p := &org.Party{
		Name:       "Agent",
		Identities: []*org.Identity{{Code: "X"}},
	}
	assert.NoError(t, rules.Validate(p, addonContext()))
}

// --- Internal helpers ---------------------------------------------------

func TestPartyIdentitySchemeAllowedEmptyScheme(t *testing.T) {
	// An Ext without the scheme ID key falls through the scheme check.
	e := tax.ExtensionsOf(tax.ExtMap{"some-other": "x"})
	assert.True(t, partyIdentitySchemeAllowed(e))
}
