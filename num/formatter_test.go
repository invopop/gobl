package num_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
)

func TestFormatterAmount(t *testing.T) {
	// make a test table
	tests := []struct {
		name string
		f    num.Formatter
		amt  num.Amount
		exp  string
	}{
		{
			name: "no unit: zero",
			f:    num.MakeFormatter(".", ","),
			amt:  num.MakeAmount(0, 2),
			exp:  "0.00",
		},
		{
			name: "no unit: zero eu",
			f:    num.MakeFormatter(",", "."),
			amt:  num.MakeAmount(0, 2),
			exp:  "0,00",
		},
		{
			name: "no unit: thousands",
			f:    num.MakeFormatter(".", ","),
			amt:  num.MakeAmount(123456, 2),
			exp:  "1,234.56",
		},
		{
			name: "no unit: thousands eu",
			f:    num.MakeFormatter(",", "."),
			amt:  num.MakeAmount(123456, 2),
			exp:  "1.234,56",
		},
		{
			name: "no unit: millions",
			f:    num.MakeFormatter(".", ","),
			amt:  num.MakeAmount(123456789, 2),
			exp:  "1,234,567.89",
		},
		{
			name: "with unit default format zero",
			f:    num.MakeFormatter(".", ",").WithUnit("%"),
			amt:  num.MakeAmount(0, 2),
			exp:  "0.00%",
		},
		{
			name: "with unit default format thousands",
			f:    num.MakeFormatter(".", ",").WithUnit("%"),
			amt:  num.MakeAmount(123456, 2),
			exp:  "1,234.56%",
		},
		{
			name: "with unit default format millions",
			f:    num.MakeFormatter(".", ",").WithUnit("%"),
			amt:  num.MakeAmount(123456789, 2),
			exp:  "1,234,567.89%",
		},
		{
			name: "with custom template format zero",
			f: num.MakeFormatter(".", ",").
				WithUnit("$").
				WithTemplate("%u%n"),
			amt: num.MakeAmount(0, 2),
			exp: "$0.00",
		},
		{
			name: "with custom template format thousands",
			f:    num.MakeFormatter(".", ",").WithUnit("$").WithTemplate("%u%n"),
			amt:  num.MakeAmount(123456, 2),
			exp:  "$1,234.56",
		},
		{
			name: "with custom template format millions",
			f:    num.MakeFormatter(",", ".").WithUnit("€").WithTemplate("%n %u"),
			amt:  num.MakeAmount(123456789, 2),
			exp:  "1.234.567,89 €",
		},
		{
			name: "with custom template format millions",
			f:    num.MakeFormatter(",", ".").WithUnit("€").WithTemplate("%n %u").WithoutUnit(),
			amt:  num.MakeAmount(123456789, 2),
			exp:  "1.234.567,89",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.exp, test.f.Amount(test.amt))
		})
	}
}

func TestFormatterPercentage(t *testing.T) {
	tests := []struct {
		name string
		f    num.Formatter
		p    num.Percentage
		exp  string
	}{
		{
			name: "default format zero",
			f:    num.MakeFormatter(".", ","),
			p:    num.MakePercentage(0, 2),
			exp:  "0%",
		},
		{
			name: "default format 16%",
			f:    num.MakeFormatter(".", ","),
			p:    num.MakePercentage(160, 3),
			exp:  "16.0%",
		},
		{
			name: "default format with unit to ignore",
			f:    num.MakeFormatter(".", ",").WithUnit("XX"),
			p:    num.MakePercentage(160, 3),
			exp:  "16.0%",
		},
		{
			name: "default with format to ignore",
			f:    num.MakeFormatter(".", ",").WithTemplate("%u%n"),
			p:    num.MakePercentage(160, 3),
			exp:  "16.0%",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.exp, test.f.Percentage(test.p))
		})
	}

}
