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
				Currency: currency.EUR,
				Amount:   num.MakeAmount(875967, 6),
			},
			exp: "",
		},
		{
			name: "missing currency",
			rate: currency.ExchangeRate{
				Amount: num.MakeAmount(875967, 6),
			},
			exp: "currency: cannot be blank",
		},
		{
			name: "missing amount",
			rate: currency.ExchangeRate{
				Currency: currency.EUR,
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
