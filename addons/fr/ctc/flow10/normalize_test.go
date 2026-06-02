package flow10

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/dgfip"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeInvoiceDispatch(t *testing.T) {
	assert.NotPanics(t, func() { normalizeInvoice(nil) })

	t.Run("B2C gets a category, no billing mode", func(t *testing.T) {
		inv := &bill.Invoice{} // no customer => B2C
		normalizeInvoice(inv)
		assert.Equal(t, B2CCategoryNotTaxable, inv.Tax.Ext.Get(ExtKeyB2CCategory))
		assert.Empty(t, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	})

	t.Run("B2B gets a billing mode, no category", func(t *testing.T) {
		inv := &bill.Invoice{Customer: &org.Party{Name: "X"}}
		normalizeInvoice(inv)
		assert.NotEmpty(t, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
		assert.Empty(t, inv.Tax.Ext.Get(ExtKeyB2CCategory))
	})
}

func TestNormalizeBillingMode(t *testing.T) {
	t.Run("keeps caller value", func(t *testing.T) {
		inv := &bill.Invoice{Tax: &bill.Tax{Ext: tax.ExtensionsOf(cbc.CodeMap{dgfip.ExtKeyBillingMode: dgfip.BillingModeB2})}}
		normalizeBillingMode(inv)
		assert.Equal(t, dgfip.BillingModeB2, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	})
	t.Run("M1 when unpaid", func(t *testing.T) {
		inv := &bill.Invoice{}
		normalizeBillingMode(inv)
		assert.Equal(t, dgfip.BillingModeM1, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	})
	t.Run("M2 when fully paid", func(t *testing.T) {
		due := num.MakeAmount(0, 2)
		inv := &bill.Invoice{Totals: &bill.Totals{Due: &due}}
		normalizeBillingMode(inv)
		assert.Equal(t, dgfip.BillingModeM2, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	})
}

func TestNormalizeB2CCategoryOnInvoice(t *testing.T) {
	t.Run("defaults to TNT1", func(t *testing.T) {
		inv := &bill.Invoice{}
		normalizeB2CCategoryOnInvoice(inv)
		require.NotNil(t, inv.Tax)
		assert.Equal(t, B2CCategoryNotTaxable, inv.Tax.Ext.Get(ExtKeyB2CCategory))
	})
	t.Run("keeps caller value", func(t *testing.T) {
		inv := &bill.Invoice{Tax: &bill.Tax{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyB2CCategory: B2CCategoryGoods})}}
		normalizeB2CCategoryOnInvoice(inv)
		assert.Equal(t, B2CCategoryGoods, inv.Tax.Ext.Get(ExtKeyB2CCategory))
	})
}

func TestNormalizeParty(t *testing.T) {
	assert.NotPanics(t, func() { normalizeParty(nil) })

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
	t.Run("empty code / nil tax ID is a no-op", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "FR"}}
		normalizeParty(p)
		assert.Empty(t, p.Identities)
		normalizeParty(&org.Party{})
	})
	t.Run("generates SIREN from SIRET", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{{Type: fr.IdentityTypeSIRET, Code: "73282932000074"}}}
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
	p := &org.Party{}
	ensureIdentity(p, fr.IdentityTypeSIREN, "", identitySchemeIDSIREN)
	assert.Empty(t, p.Identities)

	ensureIdentity(p, fr.IdentityTypeSIREN, "732829320", identitySchemeIDSIREN)
	require.Len(t, p.Identities, 1)

	// Skips when the scheme is already present.
	ensureIdentity(p, fr.IdentityTypeSIREN, "999", identitySchemeIDSIREN)
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

func TestIsEUNonFrance(t *testing.T) {
	assert.False(t, isEUNonFrance("FR"))
	assert.False(t, isEUNonFrance(""))
	assert.True(t, isEUNonFrance("DE"))
	assert.False(t, isEUNonFrance("US"))
}
