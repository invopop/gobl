package flow10

import (
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

// --- isEUNonFrance -------------------------------------------------------

func TestIsEUNonFranceEmpty(t *testing.T) {
	assert.False(t, isEUNonFrance(""))
}

func TestIsEUNonFranceFrance(t *testing.T) {
	assert.False(t, isEUNonFrance(l10n.FR))
}

func TestIsEUNonFranceSpain(t *testing.T) {
	assert.True(t, isEUNonFrance(l10n.ES))
}

func TestIsEUNonFranceUSA(t *testing.T) {
	assert.False(t, isEUNonFrance(l10n.US))
}

// --- normalizeParty ------------------------------------------------------

func TestNormalizePartyNilSafe(t *testing.T) {
	assert.NotPanics(t, func() { normalizeParty(nil) })
}

func TestNormalizePartyNoTaxID(t *testing.T) {
	p := &org.Party{Name: "Solo"}
	normalizeParty(p)
	assert.Empty(t, p.Identities)
}

func TestNormalizePartyEmptyTaxIDCode(t *testing.T) {
	p := &org.Party{TaxID: &tax.Identity{Country: "FR"}}
	normalizeParty(p)
	assert.Empty(t, p.Identities)
}

func TestNormalizePartyNonEUNonFR(t *testing.T) {
	// Non-EU / non-FR countries are left alone.
	p := &org.Party{TaxID: &tax.Identity{Country: "US", Code: "12-3456789"}}
	normalizeParty(p)
	assert.Empty(t, p.Identities)
}

// --- sirenFromFrenchTaxID ------------------------------------------------

func TestSirenFromFrenchTaxIDSIRETFallback(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Code: "35600000000011"}}}
	got := sirenFromFrenchTaxID("FR39356000000", p)
	assert.Len(t, got, 9)
}

func TestSirenFromFrenchTaxIDSIRETWrongLength(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Code: "1234"}}}
	got := sirenFromFrenchTaxID("FR39356000000", p)
	assert.Equal(t, "356000000", got)
}

func TestSirenFromFrenchTaxIDShortInput(t *testing.T) {
	got := sirenFromFrenchTaxID("FR12", &org.Party{})
	assert.Equal(t, "12", got)
}

// --- ensureIdentity ------------------------------------------------------

func TestEnsureIdentityEmptyCode(t *testing.T) {
	p := &org.Party{}
	ensureIdentity(p, "", "", "0002")
	assert.Empty(t, p.Identities)
}

func TestEnsureIdentityExistingSchemeLeftUntouched(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{
		{
			Code: "existing",
			Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0002"}),
		},
	}}
	ensureIdentity(p, "", "new", "0002")
	assert.Len(t, p.Identities, 1)
	assert.Equal(t, cbc.Code("existing"), p.Identities[0].Code)
}

// --- partyLegalSchemeID --------------------------------------------------

func TestPartyLegalSchemeIDNil(t *testing.T) {
	assert.Equal(t, "", partyLegalSchemeID(nil))
}

func TestPartyLegalSchemeIDNoSchemeExt(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Code: "X"}}}
	assert.Equal(t, "", partyLegalSchemeID(p))
}

func TestPartyLegalSchemeIDLegalScopeWins(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{
		{Code: "A", Ext: tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0227"})},
		{Code: "B", Scope: org.IdentityScopeLegal, Ext: tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0002"})},
	}}
	assert.Equal(t, "0002", partyLegalSchemeID(p))
}

func TestPartyLegalSchemeIDFallbackUsed(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{
		{Code: "A", Ext: tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "9999"})},
		{Code: "B", Ext: tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0002"})},
	}}
	assert.Equal(t, "0002", partyLegalSchemeID(p))
}
