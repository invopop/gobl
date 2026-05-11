package ctc

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

// withAddonContext is the historical name for addonContext kept so the
// ported flow2 tests read unchanged.
func withAddonContext() rules.WithContext { return addonContext() }

// --- Inbox / SIREN inbox validation (BR-FR-CO-10 + BR-FR-13) ------------

func TestElectronicAddressValidation(t *testing.T) {
	t.Run("valid SIREN inbox matching VAT", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{Country: "FR", Code: "44732829320"},
			Inboxes: []*org.Inbox{
				{Scheme: cbc.Code("0225"), Code: "732829320"},
			},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("valid SIREN inbox matching SIREN identity", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIREN,
					Code: "123456789",
					Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0002"}),
				},
			},
			Inboxes: []*org.Inbox{{Scheme: cbc.Code("0225"), Code: "123456789"}},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("SIREN inbox with + character rejected by cbc.Code base rule", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{Country: "FR", Code: "44732829320"},
			Inboxes: []*org.Inbox{
				{Scheme: cbc.Code("0225"), Code: "732829320+routing"},
			},
		}
		err := rules.Validate(party, withAddonContext())
		assert.Error(t, err)
		assert.ErrorContains(t, err, "code")
	})

	t.Run("SIREN inbox with any valid format is accepted", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{Country: "FR", Code: "44732829320"},
			Inboxes: []*org.Inbox{
				{Scheme: cbc.Code("0225"), Code: "999999999"},
			},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("SIREN inbox invalid characters", func(t *testing.T) {
		party := &org.Party{
			Inboxes: []*org.Inbox{
				{Scheme: cbc.Code("0225"), Code: "123456789@invalid"},
			},
		}
		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "must be in a valid format")
	})

	t.Run("SIREN inbox without party context is valid", func(t *testing.T) {
		party := &org.Party{
			Inboxes: []*org.Inbox{
				{Scheme: cbc.Code("0225"), Code: "123456789"},
			},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("SIREN inbox with allowed separators", func(t *testing.T) {
		party := &org.Party{
			Inboxes: []*org.Inbox{
				{Scheme: cbc.Code("0225"), Code: "123456789-test"},
			},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("SIREN inbox at cbc.Code max length (64 characters)", func(t *testing.T) {
		longCode := strings.Repeat("1", 64)
		party := &org.Party{
			Inboxes: []*org.Inbox{
				{Scheme: cbc.Code("0225"), Code: cbc.Code(longCode)},
			},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("SIREN inbox exceeds cbc.Code max length (65 characters)", func(t *testing.T) {
		tooLong := strings.Repeat("1", 65)
		party := &org.Party{
			Inboxes: []*org.Inbox{
				{Scheme: cbc.Code("0225"), Code: cbc.Code(tooLong)},
			},
		}
		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "no longer than 64")
	})
}

// --- Peppol key normalisation (normalizeInboxes) ------------------------

func TestPeppolKeyNormalization(t *testing.T) {
	ad := tax.AddonForKey(V1)

	t.Run("peppol key set on SIREN inbox when none exist", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{Country: "FR", Code: "44732829320"},
			Identities: []*org.Identity{
				{Type: fr.IdentityTypeSIREN, Code: "732829320"},
			},
			Inboxes: []*org.Inbox{{Scheme: cbc.Code("0225"), Code: "732829320"}},
		}
		ad.Normalizer(party)
		assert.Equal(t, org.InboxKeyPeppol, party.Inboxes[0].Key)
	})

	t.Run("peppol key not duplicated when another inbox already has it", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{Country: "FR", Code: "44732829320"},
			Identities: []*org.Identity{
				{Type: fr.IdentityTypeSIREN, Code: "732829320"},
			},
			Inboxes: []*org.Inbox{
				{Key: org.InboxKeyPeppol, Scheme: "0088", Code: "1234567890123"},
				{Scheme: cbc.Code("0225"), Code: "732829320"},
			},
		}
		ad.Normalizer(party)
		assert.NotEqual(t, org.InboxKeyPeppol, party.Inboxes[1].Key)
	})

	t.Run("peppol key set even for non-French party", func(t *testing.T) {
		party := &org.Party{
			TaxID:   &tax.Identity{Country: "DE", Code: "123456789"},
			Inboxes: []*org.Inbox{{Scheme: cbc.Code("0225"), Code: "123456789"}},
		}
		ad.Normalizer(party)
		assert.Equal(t, org.InboxKeyPeppol, party.Inboxes[0].Key)
	})

	t.Run("peppol key not set if no SIREN inbox", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{Country: "FR", Code: "44732829320"},
			Identities: []*org.Identity{
				{Type: fr.IdentityTypeSIREN, Code: "732829320"},
			},
			Inboxes: []*org.Inbox{{Scheme: "0088", Code: "1234567890123"}},
		}
		ad.Normalizer(party)
		assert.Equal(t, cbc.Key(""), party.Inboxes[0].Key)
	})
}

// --- Identity scheme format (BR-FR-CO-10) -------------------------------

func TestIdentitySchemeFormatValidation(t *testing.T) {
	t.Run("valid 0224 alphanumeric", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{{
				Code: "ABC123XYZ",
				Ext:  tax.ExtensionsOf(tax.ExtMap{"iso-scheme-id": "0224"}),
			}},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("valid 0224 with allowed special characters", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{{
				Code: "ABC123-info_data/route",
				Ext:  tax.ExtensionsOf(tax.ExtMap{"iso-scheme-id": "0224"}),
			}},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("invalid 0224 special chars rejected", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{{
				Code: "ABC123@invalid",
				Ext:  tax.ExtensionsOf(tax.ExtMap{"iso-scheme-id": "0224"}),
			}},
		}
		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "must be in a valid format")
	})

	t.Run("scheme 0002 not subject to 0224 alphanumeric rules", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{{
				Code: "ABC123",
				Ext:  tax.ExtensionsOf(tax.ExtMap{"iso-scheme-id": "0002"}),
			}},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("identity without scheme ID rejected (BR-FR-CO-10)", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{{
				Type: fr.IdentityTypeSIREN,
				Code: "123456789",
			}},
		}
		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "BR-FR-CO-10")
	})

	t.Run("0224 at cbc.Code max length (64)", func(t *testing.T) {
		longCode := strings.Repeat("1", 64)
		party := &org.Party{
			Identities: []*org.Identity{{
				Code: cbc.Code(longCode),
				Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0224"}),
			}},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("0224 exceeds cbc.Code max length (65)", func(t *testing.T) {
		tooLong := strings.Repeat("1", 65)
		party := &org.Party{
			Identities: []*org.Identity{{
				Code: cbc.Code(tooLong),
				Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0224"}),
			}},
		}
		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "no longer than 64")
	})
}

// --- Private ID normalization (private-id key → 0224) -------------------

func TestPrivateIDNormalization(t *testing.T) {
	ad := tax.AddonForKey(V1)

	t.Run("private-id key sets scheme 0224", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{Key: cbc.Key("private-id"), Code: "ABC123XYZ"},
			},
		}
		ad.Normalizer(party)
		assert.Equal(t, "0224", party.Identities[0].Ext.Get(iso.ExtKeySchemeID).String())
	})

	t.Run("private-id keeps pre-existing extensions", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{{
				Key:  cbc.Key("private-id"),
				Code: "ABC123XYZ",
				Ext:  tax.ExtensionsOf(tax.ExtMap{"other-key": "other-value"}),
			}},
		}
		ad.Normalizer(party)
		assert.Equal(t, "0224", party.Identities[0].Ext.Get(iso.ExtKeySchemeID).String())
		assert.Equal(t, "other-value", party.Identities[0].Ext.Get("other-key").String())
	})

	t.Run("non-private-id identity not modified", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{Type: fr.IdentityTypeSIREN, Code: "123456789"},
			},
		}
		ad.Normalizer(party)
		assert.True(t, party.Identities[0].Ext.IsZero())
	})

	t.Run("private-id overrides pre-existing scheme ID", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{{
				Key:  cbc.Key("private-id"),
				Code: "ABC123XYZ",
				Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "9999"}),
			}},
		}
		ad.Normalizer(party)
		assert.Equal(t, "0224", party.Identities[0].Ext.Get(iso.ExtKeySchemeID).String())
	})
}

// --- SIREN-from-SIRET normalization -------------------------------------

func TestSIRENGenerationFromSIRET(t *testing.T) {
	ad := tax.AddonForKey(V1)

	t.Run("SIREN generated from SIRET", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{Type: fr.IdentityTypeSIRET, Code: "12345678901234"},
			},
		}
		ad.Normalizer(party)
		assert.Len(t, party.Identities, 2)
		var siren *org.Identity
		for _, id := range party.Identities {
			if id.Type == fr.IdentityTypeSIREN {
				siren = id
			}
		}
		assert.NotNil(t, siren)
		assert.Equal(t, "123456789", siren.Code.String())
		assert.Equal(t, "0002", siren.Ext.Get(iso.ExtKeySchemeID).String())
	})

	t.Run("generated SIREN gets legal scope", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{Type: fr.IdentityTypeSIRET, Code: "12345678901234"},
			},
		}
		ad.Normalizer(party)
		var siren *org.Identity
		for _, id := range party.Identities {
			if id.Type == fr.IdentityTypeSIREN {
				siren = id
			}
		}
		assert.NotNil(t, siren)
		assert.Equal(t, org.IdentityScopeLegal, siren.Scope)
	})

	t.Run("SIREN not regenerated when already present", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{Type: fr.IdentityTypeSIRET, Code: "12345678901234"},
				{Type: fr.IdentityTypeSIREN, Code: "123456789"},
			},
		}
		ad.Normalizer(party)
		assert.Len(t, party.Identities, 2)
	})
}

// --- Identity edge cases ------------------------------------------------

func TestValidateIdentityEdgeCases(t *testing.T) {
	t.Run("nil identity returns nil", func(t *testing.T) {
		err := rules.Validate((*org.Identity)(nil), withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("0224 code over 100 chars rejected", func(t *testing.T) {
		id := &org.Identity{
			Code: cbc.Code(strings.Repeat("A", 101)),
			Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0224"}),
		}
		err := rules.Validate(id, withAddonContext())
		assert.ErrorContains(t, err, "must be no more than 100")
	})

	t.Run("0224 valid code", func(t *testing.T) {
		id := &org.Identity{
			Code: "VALID-CODE_123",
			Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0224"}),
		}
		assert.NoError(t, rules.Validate(id, withAddonContext()))
	})
}

func TestValidatePartyEdgeCases(t *testing.T) {
	t.Run("nil party returns nil", func(t *testing.T) {
		err := rules.Validate((*org.Party)(nil), withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("SIRET with mismatching SIREN rejected (BR-FR-09/10)", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
					Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0009"}),
				},
				{
					Type: fr.IdentityTypeSIREN,
					Code: "999999999",
					Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0002"}),
				},
			},
		}
		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "BR-FR-09/10")
	})

	t.Run("inbox scheme 0225 code over 125 chars rejected", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Inboxes: []*org.Inbox{
				{Scheme: "0225", Code: cbc.Code(strings.Repeat("A", 126))},
			},
		}
		err := rules.Validate(party, withAddonContext())
		assert.Error(t, err)
	})
}

func TestNormalizePartyEdgeCases(t *testing.T) {
	ad := tax.AddonForKey(V1)

	t.Run("nil party is safe", func(_ *testing.T) {
		ad.Normalizer((*org.Party)(nil))
	})

	t.Run("party without identities", func(t *testing.T) {
		party := &org.Party{Name: "Test Party"}
		ad.Normalizer(party)
		assert.Len(t, party.Identities, 0)
	})

	t.Run("party with nil identity in array", func(t *testing.T) {
		var nilID *org.Identity
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				nilID,
				{Type: fr.IdentityTypeSIRET, Code: "12345678901234"},
			},
		}
		ad.Normalizer(party)
		assert.Len(t, party.Identities, 3)
		nonNilCount := 0
		var hasSIREN, hasSIRET bool
		for _, id := range party.Identities {
			if id != nil {
				nonNilCount++
				if id.Type == fr.IdentityTypeSIREN {
					hasSIREN = true
				}
				if id.Type == fr.IdentityTypeSIRET {
					hasSIRET = true
				}
			}
		}
		assert.Equal(t, 2, nonNilCount)
		assert.True(t, hasSIREN)
		assert.True(t, hasSIRET)
	})

	t.Run("normalize inbox with nil element in array", func(t *testing.T) {
		var nilInbox *org.Inbox
		party := &org.Party{
			Name: "Test Party",
			Inboxes: []*org.Inbox{
				nilInbox,
				{Scheme: "0225", Code: "123456789-test"},
				nilInbox,
			},
		}
		ad.Normalizer(party)
		assert.Len(t, party.Inboxes, 3)
		var hasPeppol bool
		for _, inbox := range party.Inboxes {
			if inbox != nil && inbox.Key == org.InboxKeyPeppol {
				hasPeppol = true
			}
		}
		assert.True(t, hasPeppol)
	})
}

func TestValidateIdentitySchemeFormatEdgeCases(t *testing.T) {
	t.Run("empty identities returns nil", func(t *testing.T) {
		party := &org.Party{Name: "Test Party", Identities: []*org.Identity{}}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("identity without ext returns error", func(t *testing.T) {
		party := &org.Party{
			Name:       "Test Party",
			Identities: []*org.Identity{{Code: "123"}},
		}
		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "BR-FR-CO-10")
	})

	t.Run("duplicate ISO scheme IDs return error", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{Code: "123", Ext: tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0002"})},
				{Code: "456", Ext: tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0002"})},
			},
		}
		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "BR-FR-CO-10")
	})

	t.Run("nil identity in array is skipped", func(t *testing.T) {
		var nilID *org.Identity
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				nilID,
				{Code: "123", Ext: tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0002"})},
			},
		}
		assert.NoError(t, rules.Validate(party, withAddonContext()))
	})

	t.Run("0224 with empty code rejected by base identity rules", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{Code: "", Ext: tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0224"})},
				{Code: "valid-id", Ext: tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0002"})},
			},
		}
		err := rules.Validate(party, withAddonContext())
		assert.Error(t, err)
	})
}

func TestValidateInboxEdgeCases(t *testing.T) {
	t.Run("nil inbox returns nil", func(t *testing.T) {
		err := rules.Validate((*org.Inbox)(nil), withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("0225 scheme with valid code", func(t *testing.T) {
		inbox := &org.Inbox{Scheme: "0225", Code: "123456789-valid-code"}
		assert.NoError(t, rules.Validate(inbox, withAddonContext()))
	})

	t.Run("non-0225 scheme not validated", func(t *testing.T) {
		inbox := &org.Inbox{Scheme: "9999", Code: "ANY-CODE-FORMAT"}
		assert.NoError(t, rules.Validate(inbox, withAddonContext()))
	})
}

func TestItemMetaValidation(t *testing.T) {
	t.Run("valid item with meta values", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
			Meta: cbc.Meta{"order-id": "12345", "batch-code": "ABC-123"},
		}
		assert.NoError(t, rules.Validate(item, withAddonContext()))
	})

	t.Run("blank meta value rejected", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
			Meta: cbc.Meta{"order-id": "12345", "batch-code": ""},
		}
		err := rules.Validate(item, withAddonContext())
		assert.ErrorContains(t, err, "cannot be blank")
	})

	t.Run("whitespace-only meta value rejected", func(t *testing.T) {
		item := &org.Item{
			Name: "Test Item",
			Meta: cbc.Meta{"order-id": "12345", "batch-code": "   "},
		}
		err := rules.Validate(item, withAddonContext())
		assert.ErrorContains(t, err, "cannot be blank")
	})

	t.Run("item without meta", func(t *testing.T) {
		assert.NoError(t, rules.Validate(&org.Item{Name: "Test"}, withAddonContext()))
	})

	t.Run("empty meta map", func(t *testing.T) {
		item := &org.Item{Name: "Test", Meta: cbc.Meta{}}
		assert.NoError(t, rules.Validate(item, withAddonContext()))
	})

	t.Run("nil item", func(t *testing.T) {
		assert.NoError(t, rules.Validate((*org.Item)(nil), withAddonContext()))
	})
}

// --- Flow 10 normalizeParty edge cases ----------------------------------

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

func TestNormalizePartyNilSafe(t *testing.T) {
	assert.NotPanics(t, func() { normalizeParty(nil) })
}

func TestNormalizePartyNoTaxIDNoChange(t *testing.T) {
	p := &org.Party{Name: "Solo"}
	normalizeParty(p)
	assert.Empty(t, p.Identities)
}

func TestNormalizePartyEmptyTaxIDCodeNoChange(t *testing.T) {
	p := &org.Party{TaxID: &tax.Identity{Country: "FR"}}
	normalizeParty(p)
	assert.Empty(t, p.Identities)
}

func TestNormalizePartyNonEUNonFRNoChange(t *testing.T) {
	p := &org.Party{TaxID: &tax.Identity{Country: "US", Code: "12-3456789"}}
	normalizeParty(p)
	assert.Empty(t, p.Identities)
}

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

func TestEnsureIdentityEmptyCode(t *testing.T) {
	p := &org.Party{}
	ensureIdentity(p, "", "", "0002")
	assert.Empty(t, p.Identities)
}

func TestEnsureIdentityExistingSchemeLeftUntouched(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{
		Code: "existing",
		Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0002"}),
	}}}
	ensureIdentity(p, "", "new", "0002")
	assert.Len(t, p.Identities, 1)
	assert.Equal(t, cbc.Code("existing"), p.Identities[0].Code)
}

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

// --- Flow 6 role / scheme validation ------------------------------------

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

func TestPartyIdentitySchemeAllowedEmptyScheme(t *testing.T) {
	e := tax.ExtensionsOf(tax.ExtMap{"some-other": "x"})
	assert.True(t, partyIdentitySchemeAllowed(e))
}

func TestPartyRoleKnownEmptyExtPasses(t *testing.T) {
	assert.True(t, partyRoleKnown(tax.Extensions{}))
}
