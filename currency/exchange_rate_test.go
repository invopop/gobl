package currency_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
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
				Into:   currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
			exp: "",
		},
		{
			name: "missing currency",
			rate: currency.ExchangeRate{
				Amount: num.MakeAmount(875967, 6),
			},
			exp: "from: cannot be blank, into: cannot be blank",
		},
		{
			name: "missing amount",
			rate: currency.ExchangeRate{
				From: currency.USD,
				Into: currency.EUR,
			},
			exp: "amount: must not be zero",
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
	rates := []*currency.ExchangeRate{
		{
			From:   currency.USD,
			Into:   currency.EUR,
			Amount: num.MakeAmount(875967, 6),
		},
		{
			From:   currency.EUR,
			Into:   currency.USD,
			Amount: num.MakeAmount(1141860, 6),
		},
	}
	a := currency.MatchExchangeRate(rates, currency.USD, currency.EUR)
	assert.Equal(t, "0.875967", a.String())

	a = currency.MatchExchangeRate(rates, currency.EUR, currency.USD)
	assert.Equal(t, "1.141860", a.String())

	a = currency.MatchExchangeRate(rates, currency.USD, currency.USD)
	assert.Equal(t, "1", a.String())

	a = currency.MatchExchangeRate(rates, currency.USD, currency.GBP)
	assert.Nil(t, a)
}

func TestExchange(t *testing.T) {
	rates := []*currency.ExchangeRate{
		{
			From:   currency.USD,
			Into:   currency.EUR,
			Amount: num.MakeAmount(875967, 6),
		},
		{
			From:   currency.EUR,
			Into:   currency.USD,
			Amount: num.MakeAmount(1141860, 6),
		},
		{
			From:   currency.EUR,
			Into:   currency.CLP,
			Amount: num.MakeAmount(100629, 2),
		},
	}
	a := num.MakeAmount(10000, 2)
	b := currency.Exchange(rates, currency.USD, currency.EUR, a)
	assert.Equal(t, "87.60", b.String())

	b = currency.Exchange(rates, currency.EUR, currency.USD, *b)
	assert.Equal(t, "100.03", b.String())

	b = currency.Exchange(rates, currency.USD, currency.USD, *b)
	assert.Equal(t, "100.03", b.String())

	b = currency.Exchange(rates, currency.USD, currency.GBP, *b)
	assert.Nil(t, b)

	b = currency.Exchange(rates, currency.EUR, currency.CLP, a)
	assert.Equal(t, "100629", b.String())
}
