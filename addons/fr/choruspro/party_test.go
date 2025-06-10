package choruspro_test

import (
	"testing"

	"github.com/invopop/gobl/addons/fr/choruspro"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validParty() *org.Party {
	return &org.Party{
		TaxID: &tax.Identity{
			Country: "FR",
			Code:    "39356000000",
		},
	}
}

func TestValidateParty(t *testing.T) {
	addon := tax.AddonForKey(choruspro.V1)
	require.NotNil(t, addon)

	t.Run("nil party", func(t *testing.T) {
		err := addon.Validator(nil)
		assert.NoError(t, err)
	})

	t.Run("tax ID without scheme", func(t *testing.T) {
		party := validParty()
		party.TaxID = nil
		err := addon.Validator(party)
		assert.ErrorContains(t, err, "tax ID scheme must be set when no identity has scheme extension")
	})

	t.Run("Frenchtax ID with scheme", func(t *testing.T) {
		party := validParty()
		party.TaxID.Scheme = "2"
		err := addon.Validator(party)
		assert.ErrorContains(t, err, "French companies cannot have scheme set at tax identity level")
	})

	t.Run("valid EU tax ID with scheme 2", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Company",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
				Scheme:  "2",
			},
		}
		err := addon.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("valid non-EU tax ID with scheme 3", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Company",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
				Scheme:  "3",
			},
		}
		err := addon.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("valid scheme 4", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Company",
			TaxID: &tax.Identity{
				Country: "NC",
				Code:    "123456789",
				Scheme:  "4",
			},
		}
		err := addon.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("valid scheme 5", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Company",
			TaxID: &tax.Identity{
				Country: "PF",
				Code:    "123456789",
				Scheme:  "5",
			},
		}
		err := addon.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("valid scheme 6", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Company",
			TaxID: &tax.Identity{
				Country: "NC",
				Code:    "123456789",
				Scheme:  "6",
			},
		}
		err := addon.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("invalid scheme 1 for tax ID", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Company",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
				Scheme:  "1",
			},
		}
		err := addon.Validator(party)
		assert.ErrorContains(t, err, "must be a valid value")
	})

	t.Run("invalid scheme 7", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Company",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
				Scheme:  "7",
			},
		}
		err := addon.Validator(party)
		assert.ErrorContains(t, err, "must be a valid value")
	})

	t.Run("identity with scheme", func(t *testing.T) {
		t.Run("SIRET identity with correct scheme 1", func(t *testing.T) {
			party := &org.Party{
				Name: "Test Company",
				Identities: []*org.Identity{
					{
						Type: fr.IdentityTypeSIRET,
						Code: "12345678901234",
						Ext: tax.Extensions{
							choruspro.ExtKeyScheme: "1",
						},
					},
				},
			}
			err := addon.Validator(party)
			assert.NoError(t, err)
		})

		t.Run("SIRET identity with incorrect scheme", func(t *testing.T) {
			party := &org.Party{
				Name: "Test Company",
				Identities: []*org.Identity{
					{
						Type: fr.IdentityTypeSIRET,
						Code: "12345678901234",
						Ext: tax.Extensions{
							choruspro.ExtKeyScheme: "2",
						},
					},
				},
			}
			err := addon.Validator(party)
			assert.ErrorContains(t, err, "invalid value")
		})

		t.Run("non-SIRET identity with any scheme", func(t *testing.T) {
			party := &org.Party{
				Name: "Test Company",
				Identities: []*org.Identity{
					{
						Type: "OTHER",
						Code: "123456789",
						Ext: tax.Extensions{
							choruspro.ExtKeyScheme: "3",
						},
					},
				},
			}
			err := addon.Validator(party)
			assert.NoError(t, err)
		})

		t.Run("identity with scheme - tax ID ignored", func(t *testing.T) {
			party := &org.Party{
				Name: "Test Company",
				TaxID: &tax.Identity{
					Country: "DE",
					Code:    "123456789",
					// No scheme - should be ignored because identity has scheme
				},
				Identities: []*org.Identity{
					{
						Type: "OTHER",
						Code: "123456789",
						Ext: tax.Extensions{
							choruspro.ExtKeyScheme: "2",
						},
					},
				},
			}
			err := addon.Validator(party)
			assert.NoError(t, err)
		})

		t.Run("multiple identities - one with scheme", func(t *testing.T) {
			party := &org.Party{
				Name: "Test Company",
				Identities: []*org.Identity{
					{
						Type: "OTHER",
						Code: "123456789",
						// No scheme extension
					},
					{
						Type: fr.IdentityTypeSIRET,
						Code: "12345678901234",
						Ext: tax.Extensions{
							choruspro.ExtKeyScheme: "1",
						},
					},
				},
			}
			err := addon.Validator(party)
			assert.NoError(t, err)
		})

	})
}
