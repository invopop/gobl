package tax_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestFilterValidate(t *testing.T) {
	assert.NoError(t, rules.Validate(&tax.Filter{Category: tax.CategoryVAT}))
	assert.ErrorContains(t, rules.Validate(&tax.Filter{Key: tax.KeyStandard}), "tax filter category is required")
}

func TestTotalFilteredAmount(t *testing.T) {
	zero := num.MakeAmount(0, 2)
	tt := &tax.Total{
		Categories: []*tax.CategoryTotal{
			{
				Code: tax.CategoryVAT,
				Rates: []*tax.RateTotal{
					{
						Key:    tax.KeyStandard,
						Ext:    tax.ExtensionsOf(cbc.CodeMap{"foo": "1"}),
						Amount: num.MakeAmount(2100, 2),
					},
					{
						Key:    tax.KeyStandard,
						Ext:    tax.ExtensionsOf(cbc.CodeMap{"foo": "2", "bar": "x"}),
						Amount: num.MakeAmount(1050, 2),
					},
					{
						Key:    tax.KeyExempt,
						Amount: num.MakeAmount(0, 2),
					},
				},
			},
			{
				Code: "IRPF",
				Rates: []*tax.RateTotal{
					{Amount: num.MakeAmount(1500, 2)},
				},
			},
		},
	}
	tests := []struct {
		name    string
		filters []*tax.Filter
		want    string
	}{
		{"category", []*tax.Filter{{Category: tax.CategoryVAT}}, "31.50"},
		{"other category", []*tax.Filter{{Category: "IRPF"}}, "15.00"},
		{"key", []*tax.Filter{{Category: tax.CategoryVAT, Key: tax.KeyExempt}}, "0.00"},
		{"ext", []*tax.Filter{{Category: tax.CategoryVAT, Ext: tax.ExtensionsOf(cbc.CodeMap{"foo": "2"})}}, "10.50"},
		{"ext mismatch", []*tax.Filter{{Category: tax.CategoryVAT, Ext: tax.ExtensionsOf(cbc.CodeMap{"foo": "3"})}}, "0.00"},
		{"any filter", []*tax.Filter{
			{Category: tax.CategoryVAT, Ext: tax.ExtensionsOf(cbc.CodeMap{"foo": "1"})},
			{Category: tax.CategoryVAT, Ext: tax.ExtensionsOf(cbc.CodeMap{"foo": "2"})},
		}, "31.50"},
		{"overlapping filters count once", []*tax.Filter{
			{Category: tax.CategoryVAT},
			{Category: tax.CategoryVAT, Key: tax.KeyStandard},
		}, "31.50"},
		{"none", nil, "0.00"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tt.FilteredAmount(zero, tc.filters).String())
		})
	}

	t.Run("nil total", func(t *testing.T) {
		var nt *tax.Total
		assert.Equal(t, "0.00", nt.FilteredAmount(zero, []*tax.Filter{{Category: tax.CategoryVAT}}).String())
	})
}
