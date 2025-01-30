package tax_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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

func TestRegimeGetRoundingRule(t *testing.T) {
	t.Run("with", func(t *testing.T) {
		r := new(tax.RegimeDef)
		r.CalculatorRoundingRule = tax.RoundingRuleRoundThenSum
		assert.Equal(t, tax.RoundingRuleRoundThenSum, r.GetRoundingRule())
	})
	t.Run("without", func(t *testing.T) {
		r := new(tax.RegimeDef)
		assert.Equal(t, tax.RoundingRuleSumThenRound, r.GetRoundingRule())
	})
	t.Run("nil", func(t *testing.T) {
		var r *tax.RegimeDef
		assert.Equal(t, tax.RoundingRuleSumThenRound, r.GetRoundingRule())
	})
}

func TestRegimeInCategoryRates(t *testing.T) {
	var r *tax.RegimeDef // nil regime
	rate := cbc.Key("standard")
	err := validation.Validate(rate, r.InCategoryRates(tax.CategoryVAT))
	assert.ErrorContains(t, err, "must be blank when regime is undefine")
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
	t.Run("nil regime", func(t *testing.T) {
		var r *tax.RegimeDef
		assert.Nil(t, r.CategoryDef(tax.CategoryVAT))
	})
}

func TestRateDefValue(t *testing.T) {
	t.Run("with tags", func(t *testing.T) {
		rd := &tax.RateDef{
			Key:    tax.RateStandard,
			Name:   i18n.NewString("Standard"),
			Exempt: false,
			Values: []*tax.RateValueDef{
				{
					Tags:    []cbc.Key{"special"},
					Percent: num.MakePercentage(100, 3),
					Since:   cal.NewDate(2025, 1, 1),
				},
				{
					Percent: num.MakePercentage(200, 3),
					Since:   cal.NewDate(2025, 1, 1),
				},
			},
		}
		rdv := rd.Value(cal.MakeDate(2025, 1, 10), nil, nil)
		assert.Equal(t, "20.0%", rdv.Percent.String())
		rdv = rd.Value(cal.MakeDate(2025, 1, 10), []cbc.Key{"special"}, nil)
		assert.Equal(t, "10.0%", rdv.Percent.String())
	})
}
