package au_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	t.Run("normalize tax IDs", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "AU",
			Code:    "51 824 753 556",
		}
		au.New().Normalizer(tID)
		assert.Equal(t, "51824753556", tID.Code.String())
	})

	t.Run("other object", func(_ *testing.T) {
		// Ensure Normalize ignores non-tax identity objects without panic.
		au.New().Normalizer(struct{}{})
	})
}

func TestNewRegimeDef(t *testing.T) {
	r := au.New()
	if assert.NotNil(t, r) {
		assert.Equal(t, "AU", r.Country.String())
		assert.Equal(t, currency.AUD, r.Currency)
		assert.Equal(t, tax.CategoryGST, r.TaxScheme)
		assert.Equal(t, "Australia/Sydney", r.TimeZone)
		assert.NotEmpty(t, r.Categories)
		assert.Equal(t, "Australia", r.Name.String())
	}
}
