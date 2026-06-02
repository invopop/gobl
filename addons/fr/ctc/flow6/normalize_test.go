package flow6

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// statusWithLine builds a status of the given type carrying a single
// line with the given key, then normalizes it.
func statusWithLine(typ cbc.Key, key cbc.Key) *bill.Status {
	st := &bill.Status{Type: typ, Lines: []*bill.StatusLine{{Key: key}}}
	normalizeStatus(st)
	return st
}

func TestNormalizeStatusLineForward(t *testing.T) {
	cases := []struct {
		typ  cbc.Key
		key  cbc.Key
		want cbc.Code
	}{
		{bill.StatusTypeUpdate, bill.StatusLineIssued, "200"},
		{bill.StatusTypeSystem, bill.StatusLineIssued, "200"},
		{bill.StatusTypeResponse, bill.StatusLineIssued, "201"},
		{bill.StatusTypeResponse, bill.StatusLineAcknowledged, "202"},
		{bill.StatusTypeResponse, bill.StatusLineProcessing, "204"},
		{bill.StatusTypeResponse, bill.StatusLineAccepted, "205"},
		{bill.StatusTypeResponse, bill.StatusLineQuerying, "208"},
		{bill.StatusTypeResponse, bill.StatusLineRejected, "210"},
	}
	for _, c := range cases {
		st := statusWithLine(c.typ, c.key)
		assert.Equal(t, c.want, st.Lines[0].Ext.Get(ExtKeyStatus), "%s/%s", c.typ, c.key)
	}

	t.Run("unhandled key under update clears the status ext", func(t *testing.T) {
		// No ext, so prepareStatusWithLine leaves the explicit type/key
		// intact; Accepted is not an update event so the switch deletes.
		st := &bill.Status{Type: bill.StatusTypeUpdate, Lines: []*bill.StatusLine{
			{Key: bill.StatusLineAccepted},
		}}
		normalizeStatus(st)
		assert.True(t, st.Lines[0].Ext.Get(ExtKeyStatus).IsEmpty())
	})

	t.Run("unhandled key under response clears the status ext", func(t *testing.T) {
		st := &bill.Status{Type: bill.StatusTypeResponse, Lines: []*bill.StatusLine{
			{Key: bill.StatusLineError},
		}}
		normalizeStatus(st)
		assert.True(t, st.Lines[0].Ext.Get(ExtKeyStatus).IsEmpty())
	})

	t.Run("nil line is a no-op", func(t *testing.T) {
		assert.NotPanics(t, func() { normalizeStatusLine(&bill.Status{}, nil) })
	})
}

func TestPrepareStatusWithLineFromExt(t *testing.T) {
	cases := []struct {
		code    cbc.Code
		wantTyp cbc.Key
		wantKey cbc.Key
	}{
		{"200", bill.StatusTypeUpdate, bill.StatusLineIssued},
		{"201", bill.StatusTypeResponse, bill.StatusLineIssued},
		{"202", bill.StatusTypeResponse, bill.StatusLineAcknowledged},
		{"204", bill.StatusTypeResponse, bill.StatusLineProcessing},
		{"205", bill.StatusTypeResponse, bill.StatusLineAccepted},
		{"206", bill.StatusTypeResponse, bill.StatusLineRejected},
		{"207", bill.StatusTypeResponse, bill.StatusLineRejected},
		{"210", bill.StatusTypeResponse, bill.StatusLineRejected},
		{"208", bill.StatusTypeResponse, bill.StatusLineQuerying},
		{"213", bill.StatusTypeResponse, bill.StatusLineError},
	}
	for _, c := range cases {
		s := &bill.Status{}
		line := &bill.StatusLine{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyStatus: c.code})}
		prepareStatusWithLine(s, line)
		assert.Equal(t, c.wantTyp, s.Type, "code %s type", c.code)
		assert.Equal(t, c.wantKey, line.Key, "code %s key", c.code)
	}
}

func TestPrepareStatusWithLineFromKey(t *testing.T) {
	t.Run("issued defaults to update", func(t *testing.T) {
		s := &bill.Status{}
		prepareStatusWithLine(s, &bill.StatusLine{Key: bill.StatusLineIssued})
		assert.Equal(t, bill.StatusTypeUpdate, s.Type)
	})
	t.Run("other keys imply response", func(t *testing.T) {
		for _, k := range []cbc.Key{
			bill.StatusLineAcknowledged, bill.StatusLineProcessing,
			bill.StatusLineAccepted, bill.StatusLineQuerying,
			bill.StatusLineRejected, bill.StatusLineError,
		} {
			s := &bill.Status{}
			prepareStatusWithLine(s, &bill.StatusLine{Key: k})
			assert.Equal(t, bill.StatusTypeResponse, s.Type, "key %s", k)
		}
	})
	t.Run("existing type is preserved", func(t *testing.T) {
		s := &bill.Status{Type: bill.StatusTypeResponse}
		prepareStatusWithLine(s, &bill.StatusLine{Key: bill.StatusLineIssued})
		assert.Equal(t, bill.StatusTypeResponse, s.Type)
	})
}

func TestNormalizeStatusRoles(t *testing.T) {
	st := &bill.Status{
		Type:     bill.StatusTypeResponse,
		Supplier: &org.Party{Name: "S"},
		Customer: &org.Party{Name: "C"},
		Lines:    []*bill.StatusLine{{Key: bill.StatusLineAccepted}},
	}
	normalizeStatus(st)
	assert.Equal(t, RoleSeller, st.Supplier.Ext.Get(ExtKeyRole))
	assert.Equal(t, RoleBuyer, st.Customer.Ext.Get(ExtKeyRole))
}

func TestPrepareActionKeyReverse(t *testing.T) {
	cases := map[cbc.Code]cbc.Key{
		"NOA": bill.ActionKeyNone,
		"PIN": bill.ActionKeyProvide,
		"NIN": bill.ActionKeyReissue,
		"CNF": bill.ActionKeyCreditFull,
		"CNP": bill.ActionKeyCreditPartial,
		"CNA": bill.ActionKeyCreditAmount,
		"OTH": bill.ActionKeyOther,
	}
	for code, want := range cases {
		a := &bill.Action{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyAction: code})}
		normalizeAction(a)
		assert.Equal(t, want, a.Key, "code %s", code)
	}
	assert.NotPanics(t, func() { normalizeAction(nil) })
}

func TestPrepareReasonKeyReverse(t *testing.T) {
	cases := map[cbc.Code]cbc.Key{
		"COORD_BANC_ERR": bill.ReasonKeyFinanceTerms,
		"AUTRE":          bill.ReasonKeyOther,
		"NON_CONFORME":   bill.ReasonKeyLegal,
		"DOUBLON":        bill.ReasonKeyNotRecognized,
		"DEST_INC":       bill.ReasonKeyUnknownReceiver,
		"CMD_ERR":        bill.ReasonKeyReferences,
		"PU_ERR":         bill.ReasonKeyPrices,
		"QTE_ERR":        bill.ReasonKeyQuantity,
		"ART_ERR":        bill.ReasonKeyItems,
		"MODPAI_ERR":     bill.ReasonKeyPaymentTerms,
		"QUALITE_ERR":    bill.ReasonKeyQuality,
		"LIVR_INCOMP":    bill.ReasonKeyDelivery,
	}
	for code, want := range cases {
		r := &bill.Reason{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReason: code})}
		prepareReasonKey(r)
		assert.Equal(t, want, r.Key, "code %s", code)
	}
	assert.NotPanics(t, func() { normalizeReason(nil) })
}

func TestNormalizePaymentAdvice(t *testing.T) {
	pmt := &bill.Payment{
		Type:     bill.PaymentTypeAdvice,
		Supplier: &org.Party{Name: "S"},
		Customer: &org.Party{Name: "C"},
	}
	normalizePayment(pmt)
	assert.Equal(t, cbc.Code("211"), pmt.Ext.Get(ExtKeyStatus))
	assert.Equal(t, ConditionAmountPaid, pmt.Ext.Get(ExtKeyCondition))
	assert.Equal(t, RoleBuyer, pmt.Supplier.Ext.Get(ExtKeyRole))
	assert.Equal(t, RoleSeller, pmt.Customer.Ext.Get(ExtKeyRole))
	assert.NotPanics(t, func() { normalizePayment(nil) })
}

func TestNormalizeIdentityFlow6(t *testing.T) {
	assert.NotPanics(t, func() { normalizeIdentity(nil) })

	priv := &org.Identity{Key: identityKeyPrivateID, Code: "X"}
	normalizeIdentity(priv)
	assert.Equal(t, cbc.Code(identitySchemeIDPrivate), priv.Ext.Get(iso.ExtKeySchemeID))

	siren := &org.Identity{Type: fr.IdentityTypeSIREN, Code: "1"}
	normalizeIdentity(siren)
	assert.Equal(t, cbc.Code(identitySchemeIDSIREN), siren.Ext.Get(iso.ExtKeySchemeID))

	siret := &org.Identity{Type: fr.IdentityTypeSIRET, Code: "1"}
	normalizeIdentity(siret)
	assert.Equal(t, cbc.Code(identitySchemeIDSIRET), siret.Ext.Get(iso.ExtKeySchemeID))

	// An identity that already carries a scheme is left untouched.
	pre := &org.Identity{Type: fr.IdentityTypeSIREN, Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: "0238"})}
	normalizeIdentity(pre)
	assert.Equal(t, cbc.Code("0238"), pre.Ext.Get(iso.ExtKeySchemeID))
}

func TestNormalizePartyFlow6(t *testing.T) {
	assert.NotPanics(t, func() { normalizeParty(nil) })

	p := &org.Party{Identities: []*org.Identity{{Type: fr.IdentityTypeSIREN, Code: "1"}}}
	normalizeParty(p)
	assert.Equal(t, cbc.Code(identitySchemeIDSIREN), p.Identities[0].Ext.Get(iso.ExtKeySchemeID))
}

func TestNormalizeInboxesFlow6(t *testing.T) {
	assert.NotPanics(t, func() { normalizeInboxes(nil) })
	assert.NotPanics(t, func() { normalizeInboxes(&org.Party{}) })

	t.Run("flags the SIREN inbox as peppol", func(t *testing.T) {
		p := &org.Party{Inboxes: []*org.Inbox{nil, {Scheme: inboxSchemeSIREN, Code: "Y"}}}
		normalizeInboxes(p)
		assert.Equal(t, org.InboxKeyPeppol, p.Inboxes[1].Key)
	})
	t.Run("leaves SIREN inbox alone when peppol already present", func(t *testing.T) {
		p := &org.Party{Inboxes: []*org.Inbox{
			{Key: org.InboxKeyPeppol, Scheme: "9999", Code: "X"},
			{Scheme: inboxSchemeSIREN, Code: "Y"},
		}}
		normalizeInboxes(p)
		assert.Equal(t, cbc.Key(""), p.Inboxes[1].Key)
	})
}

func TestSetPartyRoleDefaultFlow6(t *testing.T) {
	assert.NotPanics(t, func() { setPartyRoleDefault(nil, RoleSeller) })

	p := &org.Party{}
	setPartyRoleDefault(p, RoleSeller)
	assert.Equal(t, RoleSeller, p.Ext.Get(ExtKeyRole))

	// Existing role is preserved.
	setPartyRoleDefault(p, RoleBuyer)
	assert.Equal(t, RoleSeller, p.Ext.Get(ExtKeyRole))
}

func TestNormalizeReasonForwardBuckets(t *testing.T) {
	// Hit each forward bucket arm so the SetOneOf chains are exercised.
	for _, key := range []cbc.Key{
		bill.ReasonKeyFinanceTerms, bill.ReasonKeyOther, bill.ReasonKeyLegal,
		bill.ReasonKeyNotRecognized, bill.ReasonKeyUnknownReceiver,
		bill.ReasonKeyReferences, bill.ReasonKeyPrices, bill.ReasonKeyQuantity,
		bill.ReasonKeyItems, bill.ReasonKeyPaymentTerms, bill.ReasonKeyQuality,
		bill.ReasonKeyDelivery,
	} {
		r := &bill.Reason{Key: key}
		normalizeReason(r)
		assert.False(t, r.Ext.Get(ExtKeyReason).IsEmpty(), "key %s reason", key)
		require.NotNil(t, r.Ext)
	}
}
