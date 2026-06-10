package choruspro_test

import (
	"testing"

	"github.com/invopop/gobl/addons/fr/choruspro"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeParty(t *testing.T) {

	t.Run("normalizes SIRET identity with scheme 1", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
					// No scheme extension initially
				},
			},
		}

		norm.Normalize(party, tax.AddonContext(choruspro.V1))

		assert.NotNil(t, party.Ext)
		assert.Equal(t, cbc.Code("1"), party.Ext.Get(choruspro.ExtKeyScheme))
	})

	t.Run("preserves existing SIRET scheme", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "1",
			}),
		}

		norm.Normalize(party, tax.AddonContext(choruspro.V1))

		assert.Equal(t, cbc.Code("1"), party.Ext.Get(choruspro.ExtKeyScheme))
	})

	t.Run("Finds SIRET identity and adds scheme extension", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: "OTHER",
					Code: "123456789",
				},
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
					// No scheme extension initially
				},
			},
		}

		norm.Normalize(party, tax.AddonContext(choruspro.V1))

		// First SIRET should be normalized
		assert.NotNil(t, party.Ext)
		assert.Equal(t, cbc.Code("1"), party.Ext.Get(choruspro.ExtKeyScheme))

	})

	t.Run("does not normalize non-SIRET identities", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: "OTHER",
					Code: "123456789",
				},
			},
		}

		norm.Normalize(party, tax.AddonContext(choruspro.V1))
		assert.True(t, party.Ext.IsZero())
	})

	t.Run("Normalizes EU company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
			},
		}

		norm.Normalize(party, tax.AddonContext(choruspro.V1))
		assert.Equal(t, cbc.Code("2"), party.Ext.Get(choruspro.ExtKeyScheme))
	})

	t.Run("Normalizes non-EU company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
		}
		norm.Normalize(party, tax.AddonContext(choruspro.V1))
		assert.Equal(t, cbc.Code("3"), party.Ext.Get(choruspro.ExtKeyScheme))
	})

	t.Run("handles nil identities", func(t *testing.T) {
		party := &org.Party{
			Name:       "Test Party",
			Identities: nil,
		}

		norm.Normalize(party, tax.AddonContext(choruspro.V1))

		assert.Nil(t, party.Identities)
	})

	t.Run("handles empty identities", func(t *testing.T) {
		party := &org.Party{
			Name:       "Test Party",
			Identities: []*org.Identity{},
		}

		norm.Normalize(party, tax.AddonContext(choruspro.V1))

		assert.Empty(t, party.Identities)
	})

	t.Run("handles nil identity elements in identities array", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				nil,
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
				nil,
			},
		}

		norm.Normalize(party, tax.AddonContext(choruspro.V1))

		// Should find the SIRET identity and add scheme extension despite nil elements
		assert.False(t, party.Ext.IsZero())
		assert.Equal(t, cbc.Code("1"), party.Ext.Get(choruspro.ExtKeyScheme))
	})

	t.Run("handles all nil identity elements in identities array", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				nil,
				nil,
			},
		}

		norm.Normalize(party, tax.AddonContext(choruspro.V1))

		// Should not panic and should not add any extension
		assert.True(t, party.Ext.IsZero())
	})
}

func withAddonContext() rules.WithContext {
	return func(rc *rules.Context) {
		rc.Set(rules.ContextKey(choruspro.V1), tax.AddonForKey(choruspro.V1))
	}
}

func TestValidateParty(t *testing.T) {
	t.Run("validates party with SIRET identity", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "1",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.NoError(t, err)
		assert.Equal(t, cbc.Code("1"), party.Ext.Get(choruspro.ExtKeyScheme))
	})

	t.Run("validates party with tax ID", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "12345678901234",
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "1",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "identities must have a SIRET entry for scheme '1'")
	})

	t.Run("scheme 1 requires SIRET identity", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: "OTHER",
					Code: "123456789",
				},
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "1",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "identities must have a SIRET entry for scheme '1'")
	})

	t.Run("scheme 1 requires French tax ID", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
			},
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "1",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "tax ID must be 'FR' for scheme '1'")
	})

	t.Run("scheme 2 requires EU non-French company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "2",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("scheme 2 rejects French company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "123456789",
			},
			Identities: []*org.Identity{
				{
					Type: "OTHER",
					Code: "123456789",
				},
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "2",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "tax ID country must be a non-French, EU company with scheme '2'")
	})

	t.Run("scheme 2 rejects non-EU company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "2",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "tax ID country must be a member of the EU with scheme '2'")
	})

	t.Run("scheme 3 accepts non-EU company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "3",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("scheme 3 rejects EU company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "3",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "tax ID country must be a non-EU company with scheme '3'")
	})

	t.Run("scheme 1 rejects Foreign tax ID", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "1",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "tax ID must be 'FR' for scheme '1'")
	})

	t.Run("scheme 4 ignores tax ID", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "4",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("missing scheme extension", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "123456789",
			},
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
			},
		}

		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "scheme extension is required")
	})

	t.Run("wrong scheme extension", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "123456789",
			},
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				choruspro.ExtKeyScheme: "2",
			}),
		}

		err := rules.Validate(party, withAddonContext())
		assert.ErrorContains(t, err, "identities cannot have a SIRET entry when not '1' scheme")
	})
}
