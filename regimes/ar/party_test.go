package ar_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ar"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPartyNormalization(t *testing.T) {
	t.Run("should normalize party tax ID", func(t *testing.T) {
		p := &org.Party{
			Name: "Empresa Argentina S.A.",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30-71458984-0",
			},
		}

		ar.Normalize(p)

		assert.Equal(t, cbc.Code("30714589840"), p.TaxID.Code)
	})

	t.Run("should normalize postal code", func(t *testing.T) {
		p := &org.Party{
			Name: "Empresa Argentina S.A.",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30714589840",
			},
			Addresses: []*org.Address{
				{
					Street:   "Av. Corrientes 1234",
					Locality: "Buenos Aires",
					Code:     "C 1043 AAZ",
				},
			},
		}

		ar.Normalize(p)

		// Postal code should be normalized to remove spaces and letters (keeping only numbers)
		assert.Equal(t, cbc.Code("1043"), p.Addresses[0].Code)
	})

	t.Run("should set country for Argentine party addresses", func(t *testing.T) {
		p := &org.Party{
			Name: "Empresa Argentina S.A.",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30714589840",
			},
			Addresses: []*org.Address{
				{
					Street:   "Av. Corrientes 1234",
					Locality: "Buenos Aires",
				},
			},
		}

		ar.Normalize(p)

		assert.Equal(t, "AR", string(p.Addresses[0].Country))
	})

	t.Run("should not override existing country", func(t *testing.T) {
		p := &org.Party{
			Name: "Empresa Argentina S.A.",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30714589840",
			},
			Addresses: []*org.Address{
				{
					Street:   "Av. Corrientes 1234",
					Locality: "Buenos Aires",
					Country:  "UY", // Different country specified
				},
			},
		}

		ar.Normalize(p)

		assert.Equal(t, "UY", string(p.Addresses[0].Country))
	})

	t.Run("should handle party without tax ID", func(t *testing.T) {
		p := &org.Party{
			Name: "Consumer",
			Addresses: []*org.Address{
				{
					Street:   "Calle Principal 123",
					Locality: "Córdoba",
					Code:     "X5000",
				},
			},
		}

		ar.Normalize(p)

		// Should not crash and should normalize postal code
		assert.Equal(t, cbc.Code("5000"), p.Addresses[0].Code)
	})

	t.Run("should handle party without addresses", func(t *testing.T) {
		p := &org.Party{
			Name: "Empresa Sin Dirección",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30-71458984-0",
			},
		}

		ar.Normalize(p)

		// Should normalize tax ID without errors
		assert.Equal(t, cbc.Code("30714589840"), p.TaxID.Code)
	})

	t.Run("should handle nil party", func(t *testing.T) {
		var p *org.Party

		// Should not crash
		assert.NotPanics(t, func() {
			ar.Normalize(p)
		})
	})

	t.Run("should normalize multiple addresses", func(t *testing.T) {
		p := &org.Party{
			Name: "Empresa con Múltiples Direcciones",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30714589840",
			},
			Addresses: []*org.Address{
				{
					Street:   "Av. Corrientes 1234",
					Locality: "Buenos Aires",
					Code:     "C1043AAZ",
				},
				{
					Street:   "Av. Santa Fe 5678",
					Locality: "Buenos Aires",
					Code:     "C1425",
				},
			},
		}

		ar.Normalize(p)

		assert.Equal(t, cbc.Code("1043"), p.Addresses[0].Code)
		assert.Equal(t, cbc.Code("1425"), p.Addresses[1].Code)
		assert.Equal(t, "AR", string(p.Addresses[0].Country))
		assert.Equal(t, "AR", string(p.Addresses[1].Country))
	})
}

func TestPartyValidation(t *testing.T) {
	t.Run("valid party should pass", func(t *testing.T) {
		p := &org.Party{
			Name: "Empresa Argentina S.A.",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30714589840",
			},
		}

		err := ar.Validate(p)
		assert.NoError(t, err)
	})

	t.Run("party with invalid tax ID should fail", func(t *testing.T) {
		p := &org.Party{
			Name: "Empresa Argentina S.A.",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "12345678901", // Invalid check digit
			},
		}

		err := ar.Validate(p)
		assert.Error(t, err)
	})

	t.Run("party without tax ID should pass validation", func(t *testing.T) {
		p := &org.Party{
			Name: "Consumer Final",
		}

		err := ar.Validate(p)
		assert.NoError(t, err)
	})
}
