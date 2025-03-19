package currency_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestExchangeRateValidation(t *testing.T) {
	tests := []struct {
		name string
		rate currency.ExchangeRate
		exp  string
	}{
		{
			name: "valid",
			rate: currency.ExchangeRate{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
			exp: "",
		},
		{
			name: "missing currency",
			rate: currency.ExchangeRate{
				Amount: num.MakeAmount(875967, 6),
			},
			exp: "from: cannot be blank; to: cannot be blank",
		},
		{
			name: "missing amount",
			rate: currency.ExchangeRate{
				From: currency.USD,
				To:   currency.EUR,
			},
			exp: "amount: must be greater than 0",
		},
		{
			name: "negative amount",
			rate: currency.ExchangeRate{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(-87596, 3),
			},
			exp: "amount: must be greater than 0",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.rate.Validate()
			if test.exp == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.exp)
			}
		})
	}
}

func TestMatchExchangeRate(t *testing.T) {
	rates := sampleRates()
	a := currency.MatchExchangeRate(rates, currency.USD, currency.EUR)
	assert.Equal(t, rates[0], a)

	a = currency.MatchExchangeRate(rates, currency.EUR, currency.USD)
	assert.Equal(t, rates[1], a)

	a = currency.MatchExchangeRate(rates, currency.USD, currency.USD)
	assert.Nil(t, a)

	a = currency.MatchExchangeRate(rates, currency.USD, currency.GBP)
	assert.Nil(t, a)
}

func TestExchangeRateConvert(t *testing.T) {
	er := &currency.ExchangeRate{
		From:   currency.USD,
		To:     currency.EUR,
		Amount: num.MakeAmount(875967, 6),
	}
	a := er.Convert(num.MakeAmount(10000, 2))
	assert.Equal(t, "87.60", a.String())

	er = &currency.ExchangeRate{
		From:   currency.EUR,
		To:     currency.CLP,
		Amount: num.MakeAmount(100629, 2),
	}
	a = er.Convert(num.MakeAmount(10000, 2))
	assert.Equal(t, "100629", a.String())
}

func TestConvert(t *testing.T) {
	rates := sampleRates()
	a := num.MakeAmount(10000, 2)
	b := currency.Convert(rates, currency.USD, currency.EUR, a)
	assert.Equal(t, "87.60", b.String())

	b = currency.Convert(rates, currency.EUR, currency.USD, *b)
	assert.Equal(t, "100.03", b.String())

	b = currency.Convert(rates, currency.USD, currency.USD, *b)
	assert.Equal(t, "100.03", b.String())

	b = currency.Convert(rates, currency.USD, currency.GBP, *b)
	assert.Nil(t, b)

	b = currency.Convert(rates, currency.EUR, currency.CLP, a)
	assert.Equal(t, "100629", b.String())
}

func TestExchangeRateValidationRule(t *testing.T) {
	rates := sampleRates()
	cur := currency.USD
	err := validation.Validate(cur, currency.CanConvertInto(rates, currency.EUR))
	assert.NoError(t, err)

	err = validation.Validate(cur, currency.CanConvertInto(rates, currency.MXN))
	assert.ErrorContains(t, err, "no exchange rate defined for 'USD' to 'MXN")

	err = validation.Validate(currency.CodeEmpty, currency.CanConvertInto(rates, currency.EUR))
	assert.NoError(t, err)

	err = validation.Validate(currency.CodeEmpty, currency.CanConvertInto(rates, currency.USD))
	assert.NoError(t, err)

	t.Run("same rate", func(t *testing.T) {
		err := validation.Validate(currency.USD, currency.CanConvertInto(rates, currency.USD))
		assert.NoError(t, err)
	})
}

func sampleRates() []*currency.ExchangeRate {
	return []*currency.ExchangeRate{
		{
			From:   currency.USD,
			To:     currency.EUR,
			Amount: num.MakeAmount(875967, 6),
		},
		{
			From:   currency.EUR,
			To:     currency.USD,
			Amount: num.MakeAmount(1141860, 6),
		},
		{
			From:   currency.EUR,
			To:     currency.CLP,
			Amount: num.MakeAmount(100629, 2),
		},
	}
}
