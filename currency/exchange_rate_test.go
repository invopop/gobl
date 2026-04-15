package currency_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestExchangeRateValidation(t *testing.T) {
	tests := []struct {
		name string
		rate currency.ExchangeRate
		ok   bool
	}{
		{
			name: "valid",
			rate: currency.ExchangeRate{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
			ok: true,
		},
		{
			name: "missing currency",
			rate: currency.ExchangeRate{
				Amount: num.MakeAmount(875967, 6),
			},
		},
		{
			name: "missing amount",
			rate: currency.ExchangeRate{
				From: currency.USD,
				To:   currency.EUR,
			},
		},
		{
			name: "negative amount",
			rate: currency.ExchangeRate{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(-87596, 3),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := rules.Validate(&test.rate)
			if test.ok {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
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
