package currency_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
)

func TestDefAmount(t *testing.T) {
	// make a test table
	tests := []struct {
		name     string
		currency currency.Code
		amt      num.Amount
		exp      string
	}{
		{
			name:     "zero EUR",
			currency: currency.EUR,
			amt:      num.MakeAmount(0, 2),
			exp:      "€0,00",
		},
		{
			name:     "thousand EUR",
			currency: currency.EUR,
			amt:      num.MakeAmount(123456, 2),
			exp:      "€1.234,56",
		},
		{
			name:     "million EUR",
			currency: currency.EUR,
			amt:      num.MakeAmount(123456789, 2),
			exp:      "€1.234.567,89",
		},
		{
			name:     "zero USD",
			currency: currency.USD,
			amt:      num.MakeAmount(0, 2),
			exp:      "$0.00",
		},
		{
			name:     "thousand USD",
			currency: currency.USD,
			amt:      num.MakeAmount(123456, 2),
			exp:      "$1,234.56",
		},
		{
			name:     "million USD",
			currency: currency.USD,
			amt:      num.MakeAmount(123456789, 2),
			exp:      "$1,234,567.89",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := test.currency.Def()
			assert.Equal(t, test.exp, d.Amount(test.amt))
		})
	}
}

func TestDefPercentage(t *testing.T) {
	// make a test table
	tests := []struct {
		name     string
		currency currency.Code
		amt      num.Percentage
		exp      string
	}{
		{
			name:     "zero EUR",
			currency: currency.EUR,
			amt:      num.MakePercentage(0, 2),
			exp:      "0%",
		},
		{
			name:     "low EUR",
			currency: currency.EUR,
			amt:      num.MakePercentage(210, 3),
			exp:      "21,0%",
		},
		{
			name:     "thousand EUR",
			currency: currency.EUR,
			amt:      num.MakePercentage(12345, 3),
			exp:      "1.234,5%",
		},
		{
			name:     "zero USD",
			currency: currency.USD,
			amt:      num.MakePercentage(0, 0),
			exp:      "0%",
		},
		{
			name:     "low USD",
			currency: currency.USD,
			amt:      num.MakePercentage(1600, 4),
			exp:      "16.00%",
		},
		{
			name:     "thousand USD",
			currency: currency.USD,
			amt:      num.MakePercentage(12345, 3),
			exp:      "1,234.5%",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := test.currency.Def()
			assert.Equal(t, test.exp, d.Percentage(test.amt))
		})
	}
}
