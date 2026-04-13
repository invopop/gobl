package it_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPartyNormalization(t *testing.T) {
	r := tax.RegimeDefFor("IT")

	t.Run("normalize customer", func(t *testing.T) {
		cus := testCustomer(t)
		cus.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "RSSGNN60R30H501U",
			Type:    "individual",
		}
		r.Normalizer(cus)
		assert.Empty(t, cus.TaxID.Code)
		assert.Empty(t, cus.TaxID.Type) //nolint:staticcheck
		assert.Len(t, cus.Identities, 1)
		assert.Equal(t, it.IdentityKeyFiscalCode, cus.Identities[0].Key)
		assert.Equal(t, "RSSGNN60R30H501U", cus.Identities[0].Code.String())
	})

	t.Run("ignore foreign party", func(t *testing.T) {
		cus := testCustomer(t)
		cus.TaxID = &tax.Identity{
			Country: "XX",
			Code:    "RSSGNN60R30H501U",
			Type:    "individual",
		}
		r.Normalizer(cus)
		assert.Equal(t, "RSSGNN60R30H501U", cus.TaxID.Code.String())
		assert.Len(t, cus.Identities, 0)
	})

	t.Run("party without tax ID", func(t *testing.T) {
		cus := testCustomer(t)
		cus.TaxID = nil
		assert.NotPanics(t, func() {
			r.Normalizer(cus)
		})
	})
}

func testCustomer(t *testing.T) *org.Party {
	t.Helper()
	return &org.Party{
		Name: "Test Customer",
		TaxID: &tax.Identity{
			Country: "IT",
			Code:    "13029381004",
		},
		Addresses: []*org.Address{
			{
				Street:   "Piazza di Test",
				Code:     "38342",
				Locality: "Venezia",
				Country:  "IT",
				Number:   "1",
			},
		},
	}
}
