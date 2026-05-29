package flow6

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeIdentity(t *testing.T) {
	assert.NotPanics(t, func() { normalizeIdentity(nil) })

	t.Run("private-id key maps to scheme 0224", func(t *testing.T) {
		id := &org.Identity{Key: identityKeyPrivateID, Code: "ABC"}
		normalizeIdentity(id)
		assert.Equal(t, cbc.Code(identitySchemeIDPrivate), id.Ext.Get(iso.ExtKeySchemeID))
	})
	t.Run("SIREN type tagged with 0002", func(t *testing.T) {
		id := &org.Identity{Type: fr.IdentityTypeSIREN, Code: "1"}
		normalizeIdentity(id)
		assert.Equal(t, cbc.Code(identitySchemeIDSIREN), id.Ext.Get(iso.ExtKeySchemeID))
	})
	t.Run("SIRET type tagged with 0009", func(t *testing.T) {
		id := &org.Identity{Type: fr.IdentityTypeSIRET, Code: "1"}
		normalizeIdentity(id)
		assert.Equal(t, cbc.Code(identitySchemeIDSIRET), id.Ext.Get(iso.ExtKeySchemeID))
	})
	t.Run("existing scheme not overwritten", func(t *testing.T) {
		id := &org.Identity{Type: fr.IdentityTypeSIREN, Code: "1", Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: "9999"})}
		normalizeIdentity(id)
		assert.Equal(t, cbc.Code("9999"), id.Ext.Get(iso.ExtKeySchemeID))
	})
}

func TestNormalizeParty(t *testing.T) {
	assert.NotPanics(t, func() { normalizeParty(nil) })

	t.Run("tags identities and flags SIREN inbox as peppol", func(t *testing.T) {
		p := &org.Party{
			Identities: []*org.Identity{{Type: fr.IdentityTypeSIREN, Code: "1"}},
			Inboxes:    []*org.Inbox{{Scheme: inboxSchemeSIREN, Code: "1_PEP"}},
		}
		normalizeParty(p)
		assert.Equal(t, cbc.Code(identitySchemeIDSIREN), p.Identities[0].Ext.Get(iso.ExtKeySchemeID))
		assert.Equal(t, org.InboxKeyPeppol, p.Inboxes[0].Key)
	})
}

func TestNormalizeInboxes(t *testing.T) {
	assert.NotPanics(t, func() { normalizeInboxes(nil) })
	assert.NotPanics(t, func() { normalizeInboxes(&org.Party{}) })

	t.Run("does not override when a peppol inbox exists", func(t *testing.T) {
		p := &org.Party{Inboxes: []*org.Inbox{
			nil,
			{Key: org.InboxKeyPeppol, Scheme: "9999"},
			{Scheme: inboxSchemeSIREN, Code: "X"},
		}}
		normalizeInboxes(p)
		assert.Equal(t, cbc.Key(""), p.Inboxes[2].Key)
	})
}

func TestPartyRoleKnown(t *testing.T) {
	// partyRoleKnown is scoped to the party's ext field, so it receives
	// a tax.Extensions value rather than the party itself.
	// empty / missing role passes
	assert.True(t, partyRoleKnown(tax.Extensions{}))
	assert.True(t, partyRoleKnown(tax.ExtensionsOf(cbc.CodeMap{})))
	// known role
	assert.True(t, partyRoleKnown(tax.ExtensionsOf(cbc.CodeMap{ExtKeyRole: RoleSeller})))
	// unknown role
	assert.False(t, partyRoleKnown(tax.ExtensionsOf(cbc.CodeMap{ExtKeyRole: "ZZ"})))
}

func TestExtValue(t *testing.T) {
	ext := tax.ExtensionsOf(cbc.CodeMap{ExtKeyRole: RoleSeller})
	assert.Equal(t, RoleSeller, extValue(ext).Get(ExtKeyRole))
	assert.Equal(t, RoleSeller, extValue(&ext).Get(ExtKeyRole))
	assert.True(t, extValue((*tax.Extensions)(nil)).IsZero())
	assert.True(t, extValue("wrong-type").IsZero())
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
	// nil entries skipped
	assert.True(t, identitiesSIRETSIRENCoherent([]*org.Identity{nil}))
}

func TestIdentitiesSchemesUnique(t *testing.T) {
	assert.True(t, identitiesSchemesUnique("wrong-type"))
	assert.True(t, identitiesSchemesUnique([]*org.Identity{}))
	mk := func(s cbc.Code) *org.Identity {
		return &org.Identity{Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: s})}
	}
	assert.True(t, identitiesSchemesUnique([]*org.Identity{nil, {Code: "x"}, mk("0002"), mk("0009")}))
	assert.False(t, identitiesSchemesUnique([]*org.Identity{mk("0002"), mk("0002")}))
}

func TestIdentitySchemeIsPrivate(t *testing.T) {
	assert.False(t, identitySchemeIsPrivate("wrong-type"))
	assert.False(t, identitySchemeIsPrivate(&org.Identity{}))
	assert.True(t, identitySchemeIsPrivate(&org.Identity{Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: identitySchemeIDPrivate})}))
}

// TestPrepareStatusWithLineFromCode covers the reverse mapping from a
// pre-set CDAR ProcessConditionCode to the (Status.Type, line.Key) pair,
// and the forward derivation of Type from a line.Key when no code is set.
func TestPrepareStatusWithLineFromCode(t *testing.T) {
	type want struct {
		typ cbc.Key
		key cbc.Key
	}
	codeCases := map[cbc.Code]want{
		"200": {bill.StatusTypeUpdate, bill.StatusLineIssued},
		"201": {bill.StatusTypeResponse, bill.StatusLineIssued},
		"202": {bill.StatusTypeResponse, bill.StatusLineAcknowledged},
		"204": {bill.StatusTypeResponse, bill.StatusLineProcessing},
		"205": {bill.StatusTypeResponse, bill.StatusLineAccepted},
		"206": {bill.StatusTypeResponse, bill.StatusLineRejected},
		"208": {bill.StatusTypeResponse, bill.StatusLineQuerying},
		"213": {bill.StatusTypeResponse, bill.StatusLineError},
	}
	for code, w := range codeCases {
		t.Run("code "+string(code), func(t *testing.T) {
			st := &bill.Status{Lines: []*bill.StatusLine{
				{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyStatus: code})},
			}}
			prepareStatusWithLine(st, st.Lines[0])
			assert.Equal(t, w.typ, st.Type)
			assert.Equal(t, w.key, st.Lines[0].Key)
		})
	}

	t.Run("derives type from key when no code", func(t *testing.T) {
		st := &bill.Status{Lines: []*bill.StatusLine{{Key: bill.StatusLineAccepted}}}
		prepareStatusWithLine(st, st.Lines[0])
		assert.Equal(t, bill.StatusTypeResponse, st.Type)
	})
	t.Run("issued key defaults to update", func(t *testing.T) {
		st := &bill.Status{Lines: []*bill.StatusLine{{Key: bill.StatusLineIssued}}}
		prepareStatusWithLine(st, st.Lines[0])
		assert.Equal(t, bill.StatusTypeUpdate, st.Type)
	})
	t.Run("explicit type preserved when no code", func(t *testing.T) {
		st := &bill.Status{Type: bill.StatusTypeResponse, Lines: []*bill.StatusLine{{Key: bill.StatusLineIssued}}}
		prepareStatusWithLine(st, st.Lines[0])
		assert.Equal(t, bill.StatusTypeResponse, st.Type)
	})
}

func TestInboxCodeValid(t *testing.T) {
	assert.True(t, inboxCodeValid("wrong-type"))
	assert.True(t, inboxCodeValid((*org.Inbox)(nil)))
	assert.True(t, inboxCodeValid(&org.Inbox{Scheme: "9999", Code: "anything goes"}))
	assert.True(t, inboxCodeValid(&org.Inbox{Scheme: inboxSchemeSIREN}))
	assert.True(t, inboxCodeValid(&org.Inbox{Scheme: inboxSchemeSIREN, Code: "1_PEP"}))
	assert.False(t, inboxCodeValid(&org.Inbox{Scheme: inboxSchemeSIREN, Code: cbc.Code(strings.Repeat("A", 126))}))
	assert.False(t, inboxCodeValid(&org.Inbox{Scheme: inboxSchemeSIREN, Code: "bad code"}))
}
