package choruspro_test

import (
	"testing"

	"github.com/invopop/gobl/addons/fr/choruspro"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeParty(t *testing.T) {
	addon := tax.AddonForKey(choruspro.V1)
	require.NotNil(t, addon)

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

		addon.Normalizer(party)

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
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "1",
			},
		}

		addon.Normalizer(party)

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

		addon.Normalizer(party)

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

		addon.Normalizer(party)
		assert.Nil(t, party.Ext)
	})

	t.Run("Normalizes EU company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
			},
		}

		addon.Normalizer(party)
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
		addon.Normalizer(party)
		assert.Equal(t, cbc.Code("3"), party.Ext.Get(choruspro.ExtKeyScheme))
	})

	t.Run("handles nil identities", func(t *testing.T) {
		party := &org.Party{
			Name:       "Test Party",
			Identities: nil,
		}

		addon.Normalizer(party)

		assert.Nil(t, party.Identities)
	})

	t.Run("handles empty identities", func(t *testing.T) {
		party := &org.Party{
			Name:       "Test Party",
			Identities: []*org.Identity{},
		}

		addon.Normalizer(party)

		assert.Empty(t, party.Identities)
	})
}

func TestValidateParty(t *testing.T) {
	addon := tax.AddonForKey(choruspro.V1)
	require.NotNil(t, addon)

	t.Run("validates party with SIRET identity", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
				},
			},
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "1",
			},
		}

		err := addon.Validator(party)
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
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "1",
			},
		}

		err := addon.Validator(party)
		assert.ErrorContains(t, err, "invalid format")
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
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "1",
			},
		}

		err := addon.Validator(party)
		assert.ErrorContains(t, err, "No SIRET identity found")
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
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "1",
			},
		}

		err := addon.Validator(party)
		assert.ErrorContains(t, err, "Customer must be a French company")
	})

	t.Run("scheme 2 requires EU non-French company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
			},
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "2",
			},
		}

		err := addon.Validator(party)
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
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "2",
			},
		}

		err := addon.Validator(party)
		assert.ErrorContains(t, err, "Customer must be a non-French, EU company")
	})

	t.Run("scheme 2 rejects non-EU company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "2",
			},
		}

		err := addon.Validator(party)
		assert.ErrorContains(t, err, "Customer must be a member of the EU")
	})

	t.Run("scheme 3 accepts non-EU company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "3",
			},
		}

		err := addon.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("scheme 3 rejects EU company", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
			},
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "3",
			},
		}

		err := addon.Validator(party)
		assert.ErrorContains(t, err, "Customer must be a non-EU company")
	})

	t.Run("scheme 1 rejects Foreign tax ID", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "1",
			},
		}

		err := addon.Validator(party)
		assert.ErrorContains(t, err, "Customer must be a French company")
	})

	t.Run("scheme 4 ignores tax ID", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "4",
			},
		}

		err := addon.Validator(party)
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

		err := addon.Validator(party)
		assert.ErrorContains(t, err, "required")
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
			Ext: tax.Extensions{
				choruspro.ExtKeyScheme: "2",
			},
		}

		err := addon.Validator(party)
		assert.ErrorContains(t, err, "SIRET identity not allowed for this extension")
	})
}
