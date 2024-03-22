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
		{
			name:     "with template",
			currency: currency.AED,
			amt:      num.MakeAmount(123456, 2),
			exp:      "1,234.56 د.إ",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := test.currency.Def()
			assert.Equal(t, test.exp, d.Format(test.amt))
		})
	}
}

func TestDefFormat(t *testing.T) {
	d := currency.USD.Def()
	f := d.Formatter(currency.WithDisambiguateSymbol())
	a := num.MakeAmount(123456, 2)
	assert.Equal(t, "US$1,234.56", f.Format(a))
}
