package flow2

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/dgfip"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeBillingMode(t *testing.T) {
	t.Run("keeps caller value", func(t *testing.T) {
		inv := &bill.Invoice{Tax: &bill.Tax{Ext: tax.ExtensionsOf(cbc.CodeMap{dgfip.ExtKeyBillingMode: dgfip.BillingModeB2})}}
		normalizeBillingMode(inv)
		assert.Equal(t, dgfip.BillingModeB2, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	})
	t.Run("M1 when unpaid", func(t *testing.T) {
		inv := &bill.Invoice{}
		normalizeBillingMode(inv)
		require.NotNil(t, inv.Tax)
		assert.Equal(t, dgfip.BillingModeM1, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	})
	t.Run("M2 when fully paid", func(t *testing.T) {
		due := num.MakeAmount(0, 2)
		inv := &bill.Invoice{Totals: &bill.Totals{Due: &due}}
		normalizeBillingMode(inv)
		assert.Equal(t, dgfip.BillingModeM2, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	})
}

func stcSupplier() *org.Party {
	return &org.Party{Identities: []*org.Identity{
		{Code: "1", Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: "0231"})},
	}}
}

func TestNormalizeSTCNote(t *testing.T) {
	t.Run("non-STC supplier is a no-op", func(t *testing.T) {
		inv := &bill.Invoice{Supplier: &org.Party{}}
		normalizeSTCNote(inv)
		assert.Empty(t, inv.Notes)
	})
	t.Run("adds TXD note for STC supplier", func(t *testing.T) {
		inv := &bill.Invoice{Supplier: stcSupplier()}
		normalizeSTCNote(inv)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, stcMembreAssujettiUnique, inv.Notes[0].Text)
		assert.Equal(t, noteSubjectTXD, inv.Notes[0].Ext.Get(untdid.ExtKeyTextSubject))
	})
	t.Run("does not duplicate an existing TXD note", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: stcSupplier(),
			Notes: []*org.Note{{
				Text: stcMembreAssujettiUnique,
				Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: noteSubjectTXD}),
			}},
		}
		normalizeSTCNote(inv)
		assert.Len(t, inv.Notes, 1)
	})
}

func TestNormalizeRequiredNotes(t *testing.T) {
	inv := &bill.Invoice{}
	normalizeRequiredNotes(inv)
	// PMT, PMD and AAB are appended when absent.
	assert.Len(t, inv.Notes, 3)
	// Idempotent: a second pass adds nothing.
	normalizeRequiredNotes(inv)
	assert.Len(t, inv.Notes, 3)
}

func TestNormalizeInvoiceDispatch(t *testing.T) {
	assert.NotPanics(t, func() { normalizeInvoice(nil) })
	inv := &bill.Invoice{}
	normalizeInvoice(inv)
	require.NotNil(t, inv.Tax)
	assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
	assert.NotEmpty(t, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	assert.GreaterOrEqual(t, len(inv.Notes), 3)
}

func TestNormalizeParty(t *testing.T) {
	assert.NotPanics(t, func() { normalizeParty(nil) })

	t.Run("derives SIREN from French tax ID", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "FR", Code: "732829320"}}
		normalizeParty(p)
		require.Len(t, p.Identities, 1)
		assert.Equal(t, fr.IdentityTypeSIREN, p.Identities[0].Type)
		assert.Equal(t, org.IdentityScopeLegal, p.Identities[0].Scope)
	})
	t.Run("non-French tax ID leaves identities untouched", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "ES", Code: "B98602642"}}
		normalizeParty(p)
		assert.Empty(t, p.Identities)
	})
	t.Run("empty tax-id code is a no-op", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "FR"}}
		normalizeParty(p)
		assert.Empty(t, p.Identities)
	})
	t.Run("derives SIREN prefix from SIRET", func(t *testing.T) {
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
	t.Run("tags private-id identity with scheme 0224", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{{Key: identityKeyPrivateID, Code: "ABC123"}}}
		normalizeParty(p)
		assert.Equal(t, cbc.Code(identitySchemeIDPrivate), p.Identities[0].Ext.Get(iso.ExtKeySchemeID))
	})
	t.Run("flags SIREN-scope inbox as peppol", func(t *testing.T) {
		p := &org.Party{Inboxes: []*org.Inbox{{Scheme: inboxSchemeSIREN, Code: "732829320_PEP"}}}
		normalizeParty(p)
		assert.Equal(t, org.InboxKeyPeppol, p.Inboxes[0].Key)
	})
	t.Run("does not override an existing peppol inbox", func(t *testing.T) {
		p := &org.Party{Inboxes: []*org.Inbox{
			{Key: org.InboxKeyPeppol, Scheme: "9999", Code: "X"},
			{Scheme: inboxSchemeSIREN, Code: "Y"},
		}}
		normalizeParty(p)
		assert.Equal(t, cbc.Key(""), p.Inboxes[1].Key)
	})
}

func TestSirenFromFrenchTaxID(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Type: fr.IdentityTypeSIRET, Code: "73282932000074"}}}
	assert.Equal(t, "732829320", sirenFromFrenchTaxID("x", p))
	assert.Equal(t, "732829320", sirenFromFrenchTaxID("FR44732829320", &org.Party{}))
	assert.Equal(t, "12", sirenFromFrenchTaxID("FR12", &org.Party{}))
}

func TestEnsureSIRENIdentity(t *testing.T) {
	p := &org.Party{}
	ensureSIRENIdentity(p, "")
	assert.Empty(t, p.Identities)

	ensureSIRENIdentity(p, "732829320")
	require.Len(t, p.Identities, 1)

	// Skips when a SIREN-scheme identity already exists.
	ensureSIRENIdentity(p, "999999999")
	assert.Len(t, p.Identities, 1)
}

func TestNormalizeIdentity(t *testing.T) {
	assert.NotPanics(t, func() { normalizeIdentity(nil) })
	siren := &org.Identity{Type: fr.IdentityTypeSIREN, Code: "1"}
	normalizeIdentity(siren)
	assert.Equal(t, cbc.Code(identitySchemeIDSIREN), siren.Ext.Get(iso.ExtKeySchemeID))
	siret := &org.Identity{Type: fr.IdentityTypeSIRET, Code: "1"}
	normalizeIdentity(siret)
	assert.Equal(t, cbc.Code(identitySchemeIDSIRET), siret.Ext.Get(iso.ExtKeySchemeID))
}

func TestIsPartyIdentitySTC(t *testing.T) {
	assert.False(t, isPartyIdentitySTC(nil))
	assert.False(t, isPartyIdentitySTC(&org.Party{}))
	assert.True(t, isPartyIdentitySTC(stcSupplier()))
}
