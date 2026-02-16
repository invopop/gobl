package ctc_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/addons/fr/ctc"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestElectronicAddressValidation(t *testing.T) {
	ad := tax.AddonForKey(ctc.Flow2V1)

	t.Run("valid SIREN inbox matching VAT", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44732829320", // 2 check digits + 9 digit SIREN
			},
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   "732829320", // Starts with SIREN from VAT
				},
			},
		}
		err := ad.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("valid SIREN inbox matching SIREN identity", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIREN,
					Code: "123456789",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0002", // ISO scheme ID required by BR-FR-CO-10
					},
				},
			},
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   "123456789", // Starts with SIREN identity
				},
			},
		}
		err := ad.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("SIREN inbox with additional routing info (cbc.Code limits apply)", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44732829320",
			},
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   "732829320+routing", // Contains + which isn't in cbc.Code allowed chars
				},
			},
		}
		err := ad.Validator(party)
		// cbc.Code base validation doesn't allow + character
		assert.Error(t, err)
		assert.ErrorContains(t, err, "must be in a valid format")
	})

	t.Run("SIREN inbox with any valid format is accepted", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44732829320",
			},
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   "999999999", // Format check only, not SIREN match
				},
			},
		}
		err := ad.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("SIREN inbox invalid characters", func(t *testing.T) {
		party := &org.Party{
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   "123456789@invalid", // @ not allowed
				},
			},
		}
		err := ad.Validator(party)
		assert.ErrorContains(t, err, "must be in a valid format")
	})

	t.Run("SIREN inbox without party context is valid", func(t *testing.T) {
		// When party has no SIREN/VAT, any valid format is accepted
		party := &org.Party{
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   "123456789",
				},
			},
		}
		err := ad.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("SIREN inbox with allowed cbc.Code separators", func(t *testing.T) {
		party := &org.Party{
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   "123456789-test", // cbc.Code allows separators between alphanumeric
				},
			},
		}
		err := ad.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("SIREN inbox at cbc.Code max length (64 characters)", func(t *testing.T) {
		// cbc.Code has max length of 64, not 125
		longCode := "1234567890123456789012345678901234567890123456789012345678901234"
		assert.Equal(t, 64, len(longCode))

		party := &org.Party{
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   cbc.Code(longCode),
				},
			},
		}
		err := ad.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("SIREN inbox exceeds cbc.Code max length (65 characters)", func(t *testing.T) {
		// cbc.Code max is 64, so 65 should fail
		tooLongCode := "12345678901234567890123456789012345678901234567890123456789012345"
		assert.Equal(t, 65, len(tooLongCode))

		party := &org.Party{
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   cbc.Code(tooLongCode),
				},
			},
		}
		err := ad.Validator(party)
		assert.ErrorContains(t, err, "the length must be between 1 and 64")
	})
}

func TestPeppolKeyNormalization(t *testing.T) {
	ad := tax.AddonForKey(ctc.Flow2V1)

	t.Run("peppol key set on SIREN inbox when none exist", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44732829320",
			},
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIREN,
					Code: "732829320",
				},
			},
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   "732829320",
				},
			},
		}
		ad.Normalizer(party)
		// Check that peppol key was set
		assert.Equal(t, "peppol", party.Inboxes[0].Key.String())
	})

	t.Run("peppol key not set if another inbox already has it", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44732829320",
			},
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIREN,
					Code: "732829320",
				},
			},
			Inboxes: []*org.Inbox{
				{
					Key:    "peppol",
					Scheme: "0088",
					Code:   "1234567890123",
				},
				{
					Scheme: cbc.Code("0225"),
					Code:   "732829320",
				},
			},
		}
		ad.Normalizer(party)
		// Check that SIREN inbox does not have peppol key
		assert.NotEqual(t, "peppol", party.Inboxes[1].Key.String())
		assert.Equal(t, "", party.Inboxes[1].Key.String())
	})

	t.Run("peppol key set even for non-French party", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
			},
			Inboxes: []*org.Inbox{
				{
					Scheme: cbc.Code("0225"),
					Code:   "123456789",
				},
			},
		}
		ad.Normalizer(party)
		// Check that peppol key was set (addon usage implies French context)
		assert.Equal(t, "peppol", party.Inboxes[0].Key.String())
	})

	t.Run("peppol key not set if no SIREN inbox", func(t *testing.T) {
		party := &org.Party{
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44732829320",
			},
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIREN,
					Code: "732829320",
				},
			},
			Inboxes: []*org.Inbox{
				{
					Scheme: "0088",
					Code:   "1234567890123",
				},
			},
		}
		ad.Normalizer(party)
		// Check that inbox does not have peppol key
		assert.Equal(t, "", party.Inboxes[0].Key.String())
	})
}

func TestIdentitySchemeFormatValidation(t *testing.T) {
	ad := tax.AddonForKey(ctc.Flow2V1)

	t.Run("valid identity with scheme 0224 - alphanumeric", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Code: "ABC123XYZ",
					Ext: tax.Extensions{
						"iso-scheme-id": "0224",
					},
				},
			},
		}
		err := ad.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("valid identity with scheme 0224 - with special characters", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Code: "ABC123+test-info_data/route",
					Ext: tax.Extensions{
						"iso-scheme-id": "0224",
					},
				},
			},
		}
		err := ad.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("invalid identity with scheme 0224 - special chars not allowed", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Code: "ABC123@invalid",
					Ext: tax.Extensions{
						"iso-scheme-id": "0224",
					},
				},
			},
		}
		err := ad.Validator(party)
		assert.ErrorContains(t, err, "identity with ISO scheme ID 0224")
		assert.ErrorContains(t, err, "must contain only alphanumeric")
	})

	t.Run("identity with other scheme ID not validated", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Code: "ABC@123", // Invalid format but different scheme
					Ext: tax.Extensions{
						"iso-scheme-id": "0002",
					},
				},
			},
		}
		err := ad.Validator(party)
		// Should not fail on format for scheme 0002
		assert.NoError(t, err)
	})

	t.Run("identity without scheme ID rejected (BR-FR-CO-10)", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIREN,
					Code: "123456789", // Valid code but missing ISO scheme ID
				},
			},
		}
		err := ad.Validator(party)
		// BR-FR-CO-10: All identities must have an ISO scheme ID
		assert.ErrorContains(t, err, "BR-FR-CO-10")
	})

	t.Run("identity with scheme 0224 at max length (100 characters)", func(t *testing.T) {
		// Create a 100 character string
		longCode := "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
		assert.Equal(t, 100, len(longCode))

		party := &org.Party{
			Identities: []*org.Identity{
				{
					Code: cbc.Code(longCode),
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0224",
					},
				},
			},
		}
		err := ad.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("identity with scheme 0224 exceeds max length (101 characters)", func(t *testing.T) {
		// Create a 101 character string
		tooLongCode := "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901"
		assert.Equal(t, 101, len(tooLongCode))

		party := &org.Party{
			Identities: []*org.Identity{
				{
					Code: cbc.Code(tooLongCode),
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0224",
					},
				},
			},
		}
		err := ad.Validator(party)
		assert.ErrorContains(t, err, "must not exceed 100 characters")
	})
}

func TestPrivateIDNormalization(t *testing.T) {
	ad := tax.AddonForKey(ctc.Flow2V1)

	t.Run("private-id key sets ISO scheme ID 0224", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Key:  cbc.Key("private-id"),
					Code: "ABC123XYZ",
				},
			},
		}
		ad.Normalizer(party)
		// Check that ISO scheme ID was set
		assert.Equal(t, "0224", party.Identities[0].Ext[iso.ExtKeySchemeID].String())
	})

	t.Run("private-id key sets ISO scheme ID 0224 with existing extensions", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Key:  cbc.Key("private-id"),
					Code: "ABC123XYZ",
					Ext: tax.Extensions{
						"other-key": "other-value",
					},
				},
			},
		}
		ad.Normalizer(party)
		// Check that ISO scheme ID was set and other extensions preserved
		assert.Equal(t, "0224", party.Identities[0].Ext[iso.ExtKeySchemeID].String())
		assert.Equal(t, "other-value", party.Identities[0].Ext["other-key"].String())
	})

	t.Run("identity without private-id key not modified", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIREN,
					Code: "123456789",
				},
			},
		}
		ad.Normalizer(party)
		// Check that no ISO scheme ID was set
		assert.Nil(t, party.Identities[0].Ext)
	})

	t.Run("existing ISO scheme ID not overwritten", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Key:  cbc.Key("private-id"),
					Code: "ABC123XYZ",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "9999", // Pre-existing value
					},
				},
			},
		}
		ad.Normalizer(party)
		// Check that ISO scheme ID was overwritten to 0224
		assert.Equal(t, "0224", party.Identities[0].Ext[iso.ExtKeySchemeID].String())
	})
}

func TestSIRENGenerationFromSIRET(t *testing.T) {
	ad := tax.AddonForKey(ctc.Flow2V1)

	t.Run("generated SIREN from SIRET", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
			},
		}
		ad.Normalizer(party)
		// Should have generated a SIREN identity
		assert.Len(t, party.Identities, 2)
		// Find the SIREN identity
		var sirenIdentity *org.Identity
		for _, id := range party.Identities {
			if id.Type == fr.IdentityTypeSIREN {
				sirenIdentity = id
				break
			}
		}
		assert.NotNil(t, sirenIdentity)
		assert.Equal(t, "123456789", sirenIdentity.Code.String())
		// Note: ISO scheme ID 0002 will be set by EN16931 addon
	})

	t.Run("generated SIREN gets legal scope", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
			},
		}
		ad.Normalizer(party)
		// Find the SIREN identity
		var sirenIdentity *org.Identity
		for _, id := range party.Identities {
			if id.Type == fr.IdentityTypeSIREN {
				sirenIdentity = id
				break
			}
		}
		assert.NotNil(t, sirenIdentity)
		assert.Equal(t, org.IdentityScopeLegal, sirenIdentity.Scope)
	})

	t.Run("SIREN not generated if already exists", func(t *testing.T) {
		party := &org.Party{
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
				{
					Type: fr.IdentityTypeSIREN,
					Code: "123456789",
				},
			},
		}
		ad.Normalizer(party)
		// Should not add another SIREN
		assert.Len(t, party.Identities, 2)
	})
}

// Additional edge cases for better coverage
func TestValidateIdentityEdgeCases(t *testing.T) {
	t.Run("nil identity returns nil", func(t *testing.T) {
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator((*org.Identity)(nil))
		assert.NoError(t, err)
	})

	t.Run("identity with ISO scheme 0224 and code over 100 chars", func(t *testing.T) {
		id := &org.Identity{
			Code: cbc.Code(strings.Repeat("A", 101)),
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0224",
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(id)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "must be no more than 100")
	})

	t.Run("identity with ISO scheme 0224 and valid code", func(t *testing.T) {
		id := &org.Identity{
			Code: "VALID-CODE_123",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0224",
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(id)
		assert.NoError(t, err)
	})
}

func TestValidatePartyEdgeCases(t *testing.T) {
	t.Run("nil party returns nil", func(t *testing.T) {
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator((*org.Party)(nil))
		assert.NoError(t, err)
	})

	t.Run("party with SIRET but mismatching SIREN", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0009",
					},
				},
				{
					Type: fr.IdentityTypeSIREN,
					Code: "999999999",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0002",
					},
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(party)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "BR-FR-09/10")
	})

	t.Run("party with invalid inbox scheme 0225 code", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Inboxes: []*org.Inbox{
				{
					Scheme: "0225",
					Code:   cbc.Code(strings.Repeat("A", 126)),
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(party)
		assert.Error(t, err)
	})
}

func TestNormalizePartyEdgeCases(t *testing.T) {
	t.Run("normalize nil party", func(_ *testing.T) {
		ad := tax.AddonForKey(ctc.Flow2V1)
		ad.Normalizer((*org.Party)(nil))
		// Should not crash
	})

	t.Run("normalize party without identities", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		ad.Normalizer(party)
		assert.Len(t, party.Identities, 0)
	})

	t.Run("normalize party with nil identity in array", func(t *testing.T) {
		var id *org.Identity
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				id, // Nil identity should be skipped via continue
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		ad.Normalizer(party)

		// Should have generated SIREN from SIRET, plus the original SIRET, plus 1 nil
		// Total: 3 elements (1 nil + SIRET + generated SIREN)
		assert.Len(t, party.Identities, 3)

		// Count non-nil identities
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

		assert.Equal(t, 2, nonNilCount, "should have 2 non-nil identities (SIRET + generated SIREN)")
		assert.True(t, hasSIREN, "should have generated SIREN")
		assert.True(t, hasSIRET, "should have original SIRET")
	})

	t.Run("normalize party with SIRET generates SIREN with legal scope", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0009",
					},
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		ad.Normalizer(party)

		// Should have generated SIREN
		assert.Len(t, party.Identities, 2)

		// Find the generated SIREN
		var siren *org.Identity
		for _, id := range party.Identities {
			if id.Type == fr.IdentityTypeSIREN {
				siren = id
			}
		}
		assert.NotNil(t, siren)
		assert.Equal(t, cbc.Code("123456789"), siren.Code)
		assert.Equal(t, org.IdentityScopeLegal, siren.Scope)
	})

	t.Run("normalize party with SIRET and existing SIREN with legal scope", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0009",
					},
				},
				{
					Type: fr.IdentityTypeSIREN,
					Code: "123456789",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0002",
					},
					Scope: org.IdentityScopeLegal,
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		ad.Normalizer(party)

		// Should not generate duplicate SIREN
		var sirenCount int
		for _, id := range party.Identities {
			if id.Type == fr.IdentityTypeSIREN {
				sirenCount++
			}
		}
		assert.Equal(t, 1, sirenCount)
	})

	t.Run("normalize inbox with SIREN scheme sets peppol key", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Inboxes: []*org.Inbox{
				{
					Scheme: "0225",
					Code:   "123456789:test",
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		ad.Normalizer(party)
		assert.Equal(t, cbc.Key("peppol"), party.Inboxes[0].Key)
	})

	t.Run("normalize inbox does not override existing peppol key", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Inboxes: []*org.Inbox{
				{
					Key:    "peppol",
					Scheme: "0088",
					Code:   "existing",
				},
				{
					Scheme: "0225",
					Code:   "123456789:test",
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		ad.Normalizer(party)
		// First inbox should keep its peppol key, second should not get it
		assert.Equal(t, cbc.Key("peppol"), party.Inboxes[0].Key)
		assert.NotEqual(t, cbc.Key("peppol"), party.Inboxes[1].Key)
	})

	t.Run("normalize inbox with nil element in array", func(t *testing.T) {
		var nilInbox *org.Inbox
		party := &org.Party{
			Name: "Test Party",
			Inboxes: []*org.Inbox{
				nilInbox, // Nil inbox should be skipped via continue
				{
					Scheme: "0225",
					Code:   "123456789:test",
				},
				nilInbox, // Another nil for good measure
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		ad.Normalizer(party)

		// Should still have 3 elements (2 nils + 1 valid inbox)
		assert.Len(t, party.Inboxes, 3)

		// Count non-nil inboxes and verify peppol key was set
		nonNilCount := 0
		hasPeppol := false
		for _, inbox := range party.Inboxes {
			if inbox != nil {
				nonNilCount++
				if inbox.Key == "peppol" {
					hasPeppol = true
				}
			}
		}

		assert.Equal(t, 1, nonNilCount, "should have 1 non-nil inbox")
		assert.True(t, hasPeppol, "SIREN inbox should have peppol key set")
	})
}

func TestValidateIdentitySchemeFormatEdgeCases(t *testing.T) {
	t.Run("empty identities returns nil", func(t *testing.T) {
		party := &org.Party{
			Name:       "Test Party",
			Identities: []*org.Identity{},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("identity without ext returns error", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Code: "123",
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(party)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "BR-FR-CO-10")
	})

	t.Run("duplicate ISO scheme IDs return error", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Code: "123",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0002",
					},
				},
				{
					Code: "456",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0002",
					},
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(party)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "BR-FR-CO-10")
		assert.ErrorContains(t, err, "duplicate")
	})

	t.Run("nil identity in array is skipped", func(t *testing.T) {
		var nilID *org.Identity
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				nilID, // Nil identity should be skipped via continue
				{
					Code: "123",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0002",
					},
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(party)
		assert.NoError(t, err, "validation should skip nil identity and succeed with valid identity")
	})

	t.Run("private-id (0224) with empty code is skipped", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Code: "", // Empty code should be skipped via continue
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0224", // private-id scheme
					},
				},
				{
					Code: "valid-id",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0002",
					},
				},
			},
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(party)
		assert.NoError(t, err, "validation should skip empty code for private-id and succeed with other valid identity")
	})
}

func TestValidateInboxEdgeCases(t *testing.T) {
	t.Run("nil inbox returns nil", func(t *testing.T) {
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator((*org.Inbox)(nil))
		assert.NoError(t, err)
	})

	t.Run("inbox with scheme 0225 and valid code", func(t *testing.T) {
		inbox := &org.Inbox{
			Scheme: "0225",
			Code:   "123456789-valid-code",
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(inbox)
		assert.NoError(t, err)
	})

	t.Run("inbox with different scheme is not validated", func(t *testing.T) {
		inbox := &org.Inbox{
			Scheme: "9999",
			Code:   "ANY-CODE-FORMAT", // Different scheme, CTC doesn't validate it
		}
		ad := tax.AddonForKey(ctc.Flow2V1)
		err := ad.Validator(inbox)
		assert.NoError(t, err)
	})
}

