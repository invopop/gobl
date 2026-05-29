package flow10

import (
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func legalIdentity(scheme cbc.Code, code string) *org.Identity {
	return &org.Identity{
		Code:  cbc.Code(code),
		Scope: org.IdentityScopeLegal,
		Ext:   tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: scheme}),
	}
}

func TestNormalizeParty(t *testing.T) {
	t.Run("nil safe", func(t *testing.T) {
		assert.NotPanics(t, func() { normalizeParty(nil) })
	})

	t.Run("French tax ID derives SIREN", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "FR", Code: "732829320"}}
		normalizeParty(p)
		require.Len(t, p.Identities, 1)
		assert.Equal(t, cbc.Code(identitySchemeIDSIREN), p.Identities[0].Ext.Get(iso.ExtKeySchemeID))
	})

	t.Run("EU non-France tax ID derives EU VAT identity", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "DE", Code: "111111125"}}
		normalizeParty(p)
		require.Len(t, p.Identities, 1)
		assert.Equal(t, cbc.Code(identitySchemeIDEUVAT), p.Identities[0].Ext.Get(iso.ExtKeySchemeID))
		assert.Equal(t, cbc.Code("DE111111125"), p.Identities[0].Code)
	})

	t.Run("empty code is a no-op", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "FR"}}
		normalizeParty(p)
		assert.Empty(t, p.Identities)
	})

	t.Run("nil tax ID is a no-op", func(t *testing.T) {
		p := &org.Party{}
		normalizeParty(p)
		assert.Empty(t, p.Identities)
	})

	t.Run("generates SIREN from SIRET", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{
			{Type: fr.IdentityTypeSIRET, Code: "73282932000074"},
		}}
		normalizeParty(p)
		var siren *org.Identity
		for _, id := range p.Identities {
			if id.Type == fr.IdentityTypeSIREN {
				siren = id
			}
		}
		require.NotNil(t, siren)
		assert.Equal(t, cbc.Code("732829320"), siren.Code)
		assert.Equal(t, org.IdentityScopeLegal, siren.Scope)
	})
}

func TestSirenFromFrenchTaxID(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Type: fr.IdentityTypeSIRET, Code: "73282932000074"}}}
	assert.Equal(t, "732829320", sirenFromFrenchTaxID("x", p))
	assert.Equal(t, "732829320", sirenFromFrenchTaxID("FR44732829320", &org.Party{}))
	assert.Equal(t, "12", sirenFromFrenchTaxID("FR12", &org.Party{}))
}

func TestEnsureIdentity(t *testing.T) {
	t.Run("empty code no-op", func(t *testing.T) {
		p := &org.Party{}
		ensureIdentity(p, fr.IdentityTypeSIREN, "", identitySchemeIDSIREN)
		assert.Empty(t, p.Identities)
	})
	t.Run("skips when scheme present", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{legalIdentity(identitySchemeIDSIREN, "1")}}
		ensureIdentity(p, fr.IdentityTypeSIREN, "2", identitySchemeIDSIREN)
		assert.Len(t, p.Identities, 1)
	})
}

func TestIsEUNonFrance(t *testing.T) {
	assert.False(t, isEUNonFrance("FR"))
	assert.False(t, isEUNonFrance(""))
	assert.True(t, isEUNonFrance("DE"))
	assert.False(t, isEUNonFrance("US"))
}

func TestPartyLegalSchemeID(t *testing.T) {
	assert.Equal(t, "", partyLegalSchemeID(nil))
	assert.Equal(t, "", partyLegalSchemeID(&org.Party{}))

	t.Run("legal-scope identity wins", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{legalIdentity(identitySchemeIDSIREN, "1")}}
		assert.Equal(t, identitySchemeIDSIREN, partyLegalSchemeID(p))
	})
	t.Run("falls back to allowed scheme without legal scope", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{
			{Code: "1", Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: identitySchemeIDNonEU})},
		}}
		assert.Equal(t, identitySchemeIDNonEU, partyLegalSchemeID(p))
	})
	t.Run("ignores identities with no scheme", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{{Code: "1"}}}
		assert.Equal(t, "", partyLegalSchemeID(p))
	})
}

func TestPartyHasSIREN(t *testing.T) {
	assert.False(t, partyHasSIREN("wrong-type"))
	assert.False(t, partyHasSIREN((*org.Party)(nil)))
	assert.False(t, partyHasSIREN(&org.Party{}))
	assert.True(t, partyHasSIREN(&org.Party{Identities: []*org.Identity{{Type: fr.IdentityTypeSIREN, Code: "1"}}}))
	assert.True(t, partyHasSIREN(&org.Party{Identities: []*org.Identity{legalIdentity(identitySchemeIDSIREN, "1")}}))
	assert.False(t, partyHasSIREN(&org.Party{Identities: []*org.Identity{nil, legalIdentity(identitySchemeIDNonEU, "1")}}))
}

func TestPartyHasAllowedLegalScheme(t *testing.T) {
	assert.False(t, partyHasAllowedLegalScheme("wrong-type"))
	assert.False(t, partyHasAllowedLegalScheme((*org.Party)(nil)))
	assert.True(t, partyHasAllowedLegalScheme(&org.Party{Identities: []*org.Identity{legalIdentity(identitySchemeIDSIREN, "1")}}))
	assert.False(t, partyHasAllowedLegalScheme(&org.Party{Identities: []*org.Identity{legalIdentity("9999", "1")}}))
}

func TestPartyHasTaxIDWhenRequired(t *testing.T) {
	assert.True(t, partyHasTaxIDWhenRequired("wrong-type"))
	assert.True(t, partyHasTaxIDWhenRequired((*org.Party)(nil)))

	t.Run("non-VAT-requiring scheme passes without tax ID", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{legalIdentity(identitySchemeIDNonEU, "1")}}
		assert.True(t, partyHasTaxIDWhenRequired(p))
	})
	t.Run("SIREN scheme requires tax ID", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{legalIdentity(identitySchemeIDSIREN, "1")}}
		assert.False(t, partyHasTaxIDWhenRequired(p))
		p.TaxID = &tax.Identity{Country: "FR", Code: "44732829320"}
		assert.True(t, partyHasTaxIDWhenRequired(p))
	})
}

func TestIdentitiesSchemesUnique(t *testing.T) {
	assert.True(t, identitiesSchemesUnique("wrong-type"))
	assert.True(t, identitiesSchemesUnique([]*org.Identity{}))
	// nil and empty-ext entries skipped
	assert.True(t, identitiesSchemesUnique([]*org.Identity{nil, {Code: "x"}}))
	unique := []*org.Identity{legalIdentity(identitySchemeIDSIREN, "1"), legalIdentity(identitySchemeIDNonEU, "2")}
	assert.True(t, identitiesSchemesUnique(unique))
	dup := []*org.Identity{legalIdentity(identitySchemeIDSIREN, "1"), legalIdentity(identitySchemeIDSIREN, "2")}
	assert.False(t, identitiesSchemesUnique(dup))
}

func TestNormalizeIdentityFlow10(t *testing.T) {
	assert.NotPanics(t, func() { normalizeIdentity(nil) })
	siren := &org.Identity{Type: fr.IdentityTypeSIREN, Code: "1"}
	normalizeIdentity(siren)
	assert.Equal(t, cbc.Code(identitySchemeIDSIREN), siren.Ext.Get(iso.ExtKeySchemeID))
	siret := &org.Identity{Type: fr.IdentityTypeSIRET, Code: "1"}
	normalizeIdentity(siret)
	assert.Equal(t, cbc.Code(identitySchemeIDSIRET), siret.Ext.Get(iso.ExtKeySchemeID))
}
