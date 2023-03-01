package org_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/invopop/gobl/l10n"
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

func TestPartyCalculate(t *testing.T) {
	party := org.Party{
		Name: "Invopop",
		TaxID: &tax.Identity{
			Country: "ES",
			Code:    "423 429 12.G",
		},
	}
	assert.NoError(t, party.Calculate())
	assert.Equal(t, l10n.ES, party.TaxID.Country)
	assert.Equal(t, "ES42342912G", party.TaxID.String())

	party = org.Party{
		Name: "Invopop",
		TaxID: &tax.Identity{
			Country: "ZZ", // no country has ZZ!
			Code:    "423 429 12.G",
		},
	}
	assert.NoError(t, party.Calculate(), "unknown entry should not cause problem")
}
