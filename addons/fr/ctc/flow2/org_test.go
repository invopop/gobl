package flow2

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sirenIdentity(code string) *org.Identity {
	return &org.Identity{
		Type:  fr.IdentityTypeSIREN,
		Code:  cbc.Code(code),
		Scope: org.IdentityScopeLegal,
		Ext:   tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: identitySchemeIDSIREN}),
	}
}

func TestNormalizeParty(t *testing.T) {
	t.Run("nil safe", func(t *testing.T) {
		assert.NotPanics(t, func() { normalizeParty(nil) })
	})

	t.Run("derives SIREN from French tax ID", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "FR", Code: "732829320"}}
		normalizeParty(p)
		require.Len(t, p.Identities, 1)
		assert.Equal(t, fr.IdentityTypeSIREN, p.Identities[0].Type)
		assert.Equal(t, cbc.Code("732829320"), p.Identities[0].Code)
		assert.Equal(t, org.IdentityScopeLegal, p.Identities[0].Scope)
	})

	t.Run("non-French tax ID leaves identities untouched", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "ES", Code: "B98602642"}}
		normalizeParty(p)
		assert.Empty(t, p.Identities)
	})

	t.Run("empty tax ID code is a no-op", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "FR"}}
		normalizeParty(p)
		assert.Empty(t, p.Identities)
	})

	t.Run("derives SIREN prefix from SIRET", func(t *testing.T) {
		p := &org.Party{
			TaxID: &tax.Identity{Country: "FR", Code: "FR12"},
			Identities: []*org.Identity{
				{Type: fr.IdentityTypeSIRET, Code: "73282932000074"},
			},
		}
		normalizeParty(p)
		// SIRET gets the 0009 scheme; a SIREN (first 9 digits) is generated.
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

	t.Run("tags private-id identity with scheme 0224", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{
			{Key: identityKeyPrivateID, Code: "ABC123"},
		}}
		normalizeParty(p)
		assert.Equal(t, cbc.Code(identitySchemeIDPrivate), p.Identities[0].Ext.Get(iso.ExtKeySchemeID))
	})

	t.Run("flags SIREN-scope inbox as peppol", func(t *testing.T) {
		p := &org.Party{Inboxes: []*org.Inbox{
			{Scheme: inboxSchemeSIREN, Code: "732829320_PEP"},
		}}
		normalizeParty(p)
		assert.Equal(t, org.InboxKeyPeppol, p.Inboxes[0].Key)
	})

	t.Run("does not override existing peppol inbox", func(t *testing.T) {
		p := &org.Party{Inboxes: []*org.Inbox{
			{Key: org.InboxKeyPeppol, Scheme: "9999", Code: "X"},
			{Scheme: inboxSchemeSIREN, Code: "Y"},
		}}
		normalizeParty(p)
		assert.Equal(t, cbc.Key(""), p.Inboxes[1].Key)
	})
}

func TestSirenFromFrenchTaxID(t *testing.T) {
	t.Run("from SIRET identity", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{
			{Type: fr.IdentityTypeSIRET, Code: "73282932000074"},
		}}
		assert.Equal(t, "732829320", sirenFromFrenchTaxID("anything", p))
	})
	t.Run("strips non-digits and takes last 9", func(t *testing.T) {
		assert.Equal(t, "732829320", sirenFromFrenchTaxID("FR44732829320", &org.Party{}))
	})
	t.Run("short code returned as-is", func(t *testing.T) {
		assert.Equal(t, "1234", sirenFromFrenchTaxID("FR1234", &org.Party{}))
	})
}

func TestEnsureSIRENIdentity(t *testing.T) {
	t.Run("empty code is a no-op", func(t *testing.T) {
		p := &org.Party{}
		ensureSIRENIdentity(p, "")
		assert.Empty(t, p.Identities)
	})
	t.Run("skips when scheme already present", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{sirenIdentity("111111111")}}
		ensureSIRENIdentity(p, "222222222")
		assert.Len(t, p.Identities, 1)
		assert.Equal(t, cbc.Code("111111111"), p.Identities[0].Code)
	})
	t.Run("appends when missing", func(t *testing.T) {
		p := &org.Party{}
		ensureSIRENIdentity(p, "333333333")
		require.Len(t, p.Identities, 1)
		assert.Equal(t, cbc.Code("333333333"), p.Identities[0].Code)
	})
}

func TestGetPartySIREN(t *testing.T) {
	assert.Equal(t, "", getPartySIREN(nil))
	assert.Equal(t, "", getPartySIREN(&org.Party{}))
	assert.Equal(t, "732829320", getPartySIREN(&org.Party{Identities: []*org.Identity{sirenIdentity("732829320")}}))
	// matches via ext scheme even without SIREN type
	p := &org.Party{Identities: []*org.Identity{
		{Code: "999", Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: identitySchemeIDSIREN})},
	}}
	assert.Equal(t, "999", getPartySIREN(p))
}

func TestIsPartyIdentitySTC(t *testing.T) {
	assert.False(t, isPartyIdentitySTC(nil))
	assert.False(t, isPartyIdentitySTC(&org.Party{}))
	assert.False(t, isPartyIdentitySTC(&org.Party{Identities: []*org.Identity{sirenIdentity("1")}}))
	stc := &org.Party{Identities: []*org.Identity{
		{Code: "1", Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: "0231"})},
	}}
	assert.True(t, isPartyIdentitySTC(stc))
}

func TestIdentitiesHasLegalSIREN(t *testing.T) {
	assert.True(t, identitiesHasLegalSIREN("wrong-type"))
	assert.False(t, identitiesHasLegalSIREN([]*org.Identity{}))
	assert.True(t, identitiesHasLegalSIREN([]*org.Identity{sirenIdentity("1")}))
	// SIREN scheme but not legal scope
	nonLegal := []*org.Identity{
		{Code: "1", Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: identitySchemeIDSIREN})},
	}
	assert.False(t, identitiesHasLegalSIREN(nonLegal))
}

func TestPartyHasSIRENInbox(t *testing.T) {
	assert.True(t, partyHasSIRENInbox("wrong-type"))
	assert.True(t, partyHasSIRENInbox((*org.Party)(nil)))
	// no SIREN at all → passes
	assert.True(t, partyHasSIRENInbox(&org.Party{}))
	// SIREN present, matching inbox
	ok := &org.Party{
		Identities: []*org.Identity{sirenIdentity("732829320")},
		Inboxes:    []*org.Inbox{{Scheme: inboxSchemeSIREN, Code: "732829320_PEP"}},
	}
	assert.True(t, partyHasSIRENInbox(ok))
	// SIREN present, no matching inbox
	bad := &org.Party{
		Identities: []*org.Identity{sirenIdentity("732829320")},
		Inboxes:    []*org.Inbox{{Scheme: "9999", Code: "X"}},
	}
	assert.False(t, partyHasSIRENInbox(bad))
}

func TestIdentitiesSIRETSIRENCoherent(t *testing.T) {
	assert.True(t, identitiesSIRETSIRENCoherent("wrong-type"))
	assert.True(t, identitiesSIRETSIRENCoherent([]*org.Identity{}))
	coherent := []*org.Identity{
		{Type: fr.IdentityTypeSIRET, Code: "73282932000074"},
		{Type: fr.IdentityTypeSIREN, Code: "732829320"},
	}
	assert.True(t, identitiesSIRETSIRENCoherent(coherent))
	incoherent := []*org.Identity{
		{Type: fr.IdentityTypeSIRET, Code: "73282932000074"},
		{Type: fr.IdentityTypeSIREN, Code: "999999999"},
	}
	assert.False(t, identitiesSIRETSIRENCoherent(incoherent))
}

func TestIdentitiesSchemeFormatValid(t *testing.T) {
	assert.NoError(t, identitiesSchemeFormatValid("wrong-type"))
	assert.NoError(t, identitiesSchemeFormatValid([]*org.Identity{}))

	t.Run("missing scheme errors", func(t *testing.T) {
		err := identitiesSchemeFormatValid([]*org.Identity{{Code: "1"}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ISO scheme ID")
	})
	t.Run("duplicate scheme errors", func(t *testing.T) {
		ids := []*org.Identity{sirenIdentity("1"), sirenIdentity("2")}
		err := identitiesSchemeFormatValid(ids)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")
	})
	t.Run("valid private-id", func(t *testing.T) {
		ids := []*org.Identity{
			{Code: "ABC-123", Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: identitySchemeIDPrivate})},
		}
		assert.NoError(t, identitiesSchemeFormatValid(ids))
	})
	t.Run("empty private-id code allowed", func(t *testing.T) {
		ids := []*org.Identity{
			{Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: identitySchemeIDPrivate})},
		}
		assert.NoError(t, identitiesSchemeFormatValid(ids))
	})
	t.Run("private-id too long errors", func(t *testing.T) {
		ids := []*org.Identity{
			{Code: cbc.Code(strings.Repeat("A", 101)), Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: identitySchemeIDPrivate})},
		}
		err := identitiesSchemeFormatValid(ids)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "100 characters")
	})
	t.Run("private-id bad format errors", func(t *testing.T) {
		ids := []*org.Identity{
			{Code: "bad code!", Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: identitySchemeIDPrivate})},
		}
		err := identitiesSchemeFormatValid(ids)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "alphanumeric")
	})
}

func TestInboxCodeValid(t *testing.T) {
	assert.True(t, inboxCodeValid("wrong-type"))
	assert.True(t, inboxCodeValid((*org.Inbox)(nil)))
	// non-SIREN scheme passes regardless
	assert.True(t, inboxCodeValid(&org.Inbox{Scheme: "9999", Code: "anything goes"}))
	// SIREN scheme empty code passes
	assert.True(t, inboxCodeValid(&org.Inbox{Scheme: inboxSchemeSIREN}))
	// SIREN scheme valid code
	assert.True(t, inboxCodeValid(&org.Inbox{Scheme: inboxSchemeSIREN, Code: "732829320_PEP"}))
	// SIREN scheme too long
	assert.False(t, inboxCodeValid(&org.Inbox{Scheme: inboxSchemeSIREN, Code: cbc.Code(strings.Repeat("A", 126))}))
	// SIREN scheme bad format
	assert.False(t, inboxCodeValid(&org.Inbox{Scheme: inboxSchemeSIREN, Code: "bad code"}))
}

func TestSchemeGuards(t *testing.T) {
	assert.False(t, identitySchemeIs0224("wrong-type"))
	assert.False(t, identitySchemeIs0224(&org.Identity{}))
	assert.True(t, identitySchemeIs0224(&org.Identity{Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: identitySchemeIDPrivate})}))

	assert.False(t, inboxSchemeIs0225("wrong-type"))
	assert.False(t, inboxSchemeIs0225(&org.Inbox{Scheme: "9999"}))
	assert.True(t, inboxSchemeIs0225(&org.Inbox{Scheme: inboxSchemeSIREN}))
}

func TestMetaNoBlankValues(t *testing.T) {
	assert.NoError(t, metaNoBlankValues("wrong-type"))
	assert.NoError(t, metaNoBlankValues(cbc.Meta{}))
	assert.NoError(t, metaNoBlankValues(cbc.Meta{"k": "v"}))
	err := metaNoBlankValues(cbc.Meta{"k": "   "})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be blank")
}
