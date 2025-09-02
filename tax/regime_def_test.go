package tax_test

import (
	"context"
	"testing"
	"time"

	"github.com/invopop/gobl/i18n"
	_ "github.com/invopop/gobl/regimes"
	"github.com/invopop/gobl/regimes/pt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegimeTimeLocation(t *testing.T) {
	r := new(tax.RegimeDef)
	r.TimeZone = "Europe/Amsterdam"
	loc, err := time.LoadLocation("Europe/Amsterdam")
	require.NoError(t, err)

	assert.Equal(t, loc, r.TimeLocation())

	r.TimeZone = "INVALID"
	loc = r.TimeLocation()
	assert.Equal(t, loc, time.UTC)
}

func TestRegimeGetCurrency(t *testing.T) {
	t.Run("with", func(t *testing.T) {
		r := new(tax.RegimeDef)
		r.Currency = currency.EUR
		assert.Equal(t, currency.EUR, r.GetCurrency())
	})
	t.Run("without", func(t *testing.T) {
		r := new(tax.RegimeDef)
		assert.Empty(t, r.GetCurrency())
	})
	t.Run("nil", func(t *testing.T) {
		var r *tax.RegimeDef
		assert.Empty(t, r.GetCurrency())
	})
	t.Run("currency def", func(t *testing.T) {
		r := new(tax.RegimeDef)
		r.Currency = currency.EUR
		assert.Equal(t, r.CurrencyDef().Name, "Euro")
	})
}

func TestRegimeDefScenarioSet(t *testing.T) {
	t.Run("with scenario", func(t *testing.T) {
		r := es.New()
		ss := r.ScenarioSet("bill/invoice")
		assert.NotNil(t, ss)
	})
	t.Run("without scenario", func(t *testing.T) {
		r := es.New()
		ss := r.ScenarioSet("unknown")
		assert.Nil(t, ss)
	})
}

func TestRegimeDefGetCountry(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var r *tax.RegimeDef
		assert.Empty(t, r.GetCountry())
	})
	t.Run("with", func(t *testing.T) {
		r := new(tax.RegimeDef)
		r.Country = "DE"
		assert.Equal(t, l10n.DE.Tax(), r.GetCountry())
	})
}

func TestRegimeGetRoundingRule(t *testing.T) {
	t.Run("with", func(t *testing.T) {
		r := new(tax.RegimeDef)
		r.CalculatorRoundingRule = tax.RoundingRuleCurrency
		assert.Equal(t, tax.RoundingRuleCurrency, r.GetRoundingRule())
	})
	t.Run("without", func(t *testing.T) {
		r := new(tax.RegimeDef)
		assert.Equal(t, tax.RoundingRulePrecise, r.GetRoundingRule())
	})
	t.Run("nil", func(t *testing.T) {
		var r *tax.RegimeDef
		assert.Equal(t, tax.RoundingRulePrecise, r.GetRoundingRule())
	})
}

func TestRegimeInCategoryRates(t *testing.T) {
	var r *tax.RegimeDef // nil regime
	rate := cbc.Key("general")
	err := validation.Validate(rate, r.InCategoryRates(tax.CategoryVAT, tax.KeyStandard))
	assert.ErrorContains(t, err, "must be blank when regime is undefine")
}

func TestRegimeInCategoryRule(t *testing.T) {
	t.Run("no rates", func(t *testing.T) {
		r := es.New()
		err := validation.Validate(tax.RateGeneral, r.InCategoryRates(es.TaxCategoryIPSI, cbc.KeyEmpty))
		assert.ErrorContains(t, err, "must be blank for category 'IPSI' with no key")
	})
	t.Run("invalid rate", func(t *testing.T) {
		r := es.New()
		err := validation.Validate(cbc.Key("foo"), r.InCategoryRates(tax.CategoryVAT, tax.KeyStandard))
		assert.ErrorContains(t, err, "'foo' not defined in 'VAT' category for key 'standard'")
	})
}

func TestRegimeDefValidateObject(t *testing.T) {
	t.Run("nil regime", func(t *testing.T) {
		var r *tax.RegimeDef
		err := r.ValidateObject(&org.Note{})
		assert.Nil(t, err)
	})
	t.Run("without validator", func(t *testing.T) {
		r := new(tax.RegimeDef)
		err := r.ValidateObject(&org.Note{})
		assert.Nil(t, err)
	})
}

func TestRegimeDefNormalizeObject(t *testing.T) {
	t.Run("nil regime", func(t *testing.T) {
		var r *tax.RegimeDef
		assert.NotPanics(t, func() {
			r.NormalizeObject(&org.Note{})
		})
	})
}

func TestRegimeDefCategoryDef(t *testing.T) {
	t.Run("nil regime for known category", func(t *testing.T) {
		var r *tax.RegimeDef
		cd := r.CategoryDef(tax.CategoryVAT)
		assert.NotNil(t, cd)
		assert.Equal(t, tax.CategoryVAT, cd.Code)
	})
	t.Run("nil regime for unknown category", func(t *testing.T) {
		var r *tax.RegimeDef
		cd := r.CategoryDef(cbc.Code("UNKNOWN"))
		assert.Nil(t, cd)
	})
}

func TestRegimeDefNormalizers(t *testing.T) {
	t.Run("nil regime", func(t *testing.T) {
		var r *tax.RegimeDef
		assert.Nil(t, r.Normalizers())
	})

	t.Run("with normalizer", func(t *testing.T) {
		r := &tax.RegimeDef{
			Normalizer: func(_ any) {
				// nothing here
			},
		}
		assert.NotNil(t, r.Normalizers())
		assert.Len(t, r.Normalizers(), 1)
	})

	t.Run("without normalizer", func(t *testing.T) {
		r := &tax.RegimeDef{}
		assert.Nil(t, r.Normalizers())
	})
}

func TestCategoryDefValidations(t *testing.T) {
	r := tax.RegimeDefFor("PT")
	ctx := r.WithContext(context.Background())

	t.Run("valid", func(t *testing.T) {
		c := baseCategoryDef()
		err := c.ValidateWithContext(ctx)
		require.NoError(t, err)
	})

	t.Run("informative", func(t *testing.T) {
		c := baseCategoryDef()
		c.Informative = true

		err := c.ValidateWithContext(ctx)
		require.NoError(t, err)
	})

	t.Run("retained", func(t *testing.T) {
		c := baseCategoryDef()
		c.Retained = true

		err := c.ValidateWithContext(ctx)
		require.NoError(t, err)
	})

	t.Run("informative and retained", func(t *testing.T) {
		c := baseCategoryDef()
		c.Informative = true
		c.Retained = true

		err := c.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "cannot be true when informative is true")
	})

	t.Run("with valid extensions", func(t *testing.T) {
		c := baseCategoryDef()
		c.Extensions = []cbc.Key{pt.ExtKeyRegion}

		err := c.ValidateWithContext(ctx)
		require.NoError(t, err)
	})

	t.Run("with invalid extensions", func(t *testing.T) {
		c := baseCategoryDef()
		c.Extensions = []cbc.Key{"INVALID"}

		err := c.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "must be a valid value")
	})
}

func baseCategoryDef() *tax.CategoryDef {
	return &tax.CategoryDef{
		Code:  "TEST",
		Name:  i18n.NewString("TEST"),
		Title: i18n.NewString("Test tax"),
	}
}
