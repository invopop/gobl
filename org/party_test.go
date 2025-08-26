package org_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invopop/gobl/org"
	_ "github.com/invopop/gobl/regimes"
	"github.com/invopop/gobl/tax"
)

func TestEmailValidation(t *testing.T) {
	valid := org.Email{
		Address: "foobar@invopop.example.com",
	}
	assert.NoError(t, valid.Validate())

	invalid := org.Email{
		Address: "foobar",
	}
	assert.EqualError(t, invalid.Validate(), "addr: must be a valid email address.")
}

func TestPartyNormalize(t *testing.T) {
	t.Run("for known regime", func(t *testing.T) {
		party := org.Party{
			Name: "Invopop",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "423 429 12.G",
			},
		}
		party.Normalize(nil)
		assert.Empty(t, party.GetRegime())
		assert.Equal(t, "ES", party.TaxID.Country.String())
		assert.Equal(t, "ES42342912G", party.TaxID.String())
	})

	t.Run("for known regime with Calculate", func(t *testing.T) {
		party := org.Party{
			Name: "Invopop",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "423 429 12.G",
			},
		}
		assert.NoError(t, party.Calculate())
		assert.Empty(t, party.GetRegime())
		assert.Equal(t, "ES", party.TaxID.Country.String())
		assert.Equal(t, "ES42342912G", party.TaxID.String())
	})

	t.Run("for unknown regime", func(t *testing.T) {
		party := org.Party{
			Name: "Invopop",
			TaxID: &tax.Identity{
				Country: "ZZ", // no country has ZZ!
				Code:    "423 429 12.G",
			},
		}
		party.Normalize(nil) // unknown entry should not cause problem
		assert.Equal(t, "42342912G", party.TaxID.Code.String())
	})

	t.Run("for specific regime", func(t *testing.T) {
		party := org.Party{
			Regime: tax.WithRegime("DE"),
			Name:   "Invopop",
			Identities: []*org.Identity{
				{
					Key:  "de-tax-number",
					Code: "123 456 78901",
				},
			},
		}
		require.NoError(t, party.Calculate())
		assert.Equal(t, "DE", party.GetRegime().String())
		assert.Equal(t, "123/456/78901", party.Identities[0].Code.String())
	})

	t.Run("with telephone", func(t *testing.T) {
		party := org.Party{
			Name: "Invopop",
			Telephones: []*org.Telephone{
				{
					Number: " +49 123 4567890 ",
				},
			},
		}
		party.Normalize(nil)
		assert.Equal(t, "+49 123 4567890", party.Telephones[0].Number)
	})
}

func TestPartyAddressNill(t *testing.T) {
	party := org.Party{
		Addresses: []*org.Address{nil},
	}
	party.Normalize(nil)
	assert.NoError(t, party.Validate())
}

func TestPartyValidation(t *testing.T) {
	t.Run("with regime", func(t *testing.T) {
		party := org.Party{
			Regime: tax.WithRegime("DE"),
			Name:   "Invopop",
			Identities: []*org.Identity{
				{
					Key:  "de-tax-number",
					Code: "123 456 78901",
				},
			},
		}
		require.NoError(t, party.Calculate())
		assert.NoError(t, party.Validate())
		assert.Equal(t, "DE", party.GetRegime().String())
	})
	t.Run("with regime and bad code", func(t *testing.T) {
		party := org.Party{
			Regime: tax.WithRegime("DE"),
			Name:   "Invopop",
			Identities: []*org.Identity{
				{
					Key:  "de-tax-number",
					Code: "1231312423432422",
				},
			},
		}
		require.NoError(t, party.Calculate())
		err := party.Validate()
		assert.ErrorContains(t, err, "identities: (0: (code: must be in a valid format.).).")
	})
}
