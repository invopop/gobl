package au

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	t.Run("normalize tax IDs", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "AU",
			Code:    "51 824 753 556",
		}
		New().Normalizer(tID)
		assert.Equal(t, "51824753556", tID.Code.String())
	})

	t.Run("other object", func(_ *testing.T) {
		// Ensure Normalize ignores non-tax identity objects without panic.
		New().Normalizer(struct{}{})
	})
}

func TestNewRegimeDef(t *testing.T) {
	r := New()
	if assert.NotNil(t, r) {
		assert.Equal(t, "AU", r.Country.String())
		assert.Equal(t, currency.AUD, r.Currency)
		assert.Equal(t, tax.CategoryGST, r.TaxScheme)
		assert.Equal(t, "Australia/Sydney", r.TimeZone)
		assert.NotEmpty(t, r.Categories)
		assert.Equal(t, "Australia", r.Name.String())
	}
}

func TestTaxCategories(t *testing.T) {
	if assert.Len(t, taxCategories, 1) {
		cat := taxCategories[0]
		assert.Equal(t, tax.CategoryGST, cat.Code)
		assert.False(t, cat.Retained)
		assert.NotEmpty(t, cat.Keys)

		std := cat.RateDef(tax.KeyStandard, tax.RateGeneral)
		assert.NotNil(t, std)

		zero := cat.RateDef(tax.KeyZero, tax.RateZero)
		assert.NotNil(t, zero)

		exempt := cat.RateDef(tax.KeyExempt, tax.RateZero)
		assert.NotNil(t, exempt)
	}
}
