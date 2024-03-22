package num_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
)

func TestFormatter(t *testing.T) {
	// make a test table
	f := &num.Formatter{
		DecimalMark:        ".",
		ThousandsSeparator: ",",
	}
	tests := []struct {
		name string
		f    *num.Formatter
		amt  num.Amount
		exp  string
	}{
		{
			name: "no unit: zero",
			f:    &num.Formatter{".", ",", "", ""},
			amt:  num.MakeAmount(0, 2),
			exp:  "0.00",
		},
		{
			name: "no unit: zero eu",
			f:    &num.Formatter{",", ".", "", ""},
			amt:  num.MakeAmount(0, 2),
			exp:  "0,00",
		},
		{
			name: "no unit: thousands",
			f:    &num.Formatter{".", ",", "", ""},
			amt:  num.MakeAmount(123456, 2),
			exp:  "1,234.56",
		},
		{
			name: "no unit: thousands eu",
			f:    &num.Formatter{",", ".", "", ""},
			amt:  num.MakeAmount(123456, 2),
			exp:  "1.234,56",
		},
		{
			name: "no unit: millions",
			f:    &num.Formatter{".", ",", "", ""},
			amt:  num.MakeAmount(123456789, 2),
			exp:  "1,234,567.89",
		},
		{
			name: "with unit default format zero",
			f:    &num.Formatter{".", ",", "%", ""},
			amt:  num.MakeAmount(0, 2),
			exp:  "0.00%",
		},
		{
			name: "with unit default format thousands",
			f:    &num.Formatter{".", ",", "%", ""},
			amt:  num.MakeAmount(123456, 2),
			exp:  "1,234.56%",
		},
		{
			name: "with unit default format millions",
			f:    &num.Formatter{".", ",", "%", ""},
			amt:  num.MakeAmount(123456789, 2),
			exp:  "1,234,567.89%",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.exp, f.Format(test.amt))
		})
	}
}
