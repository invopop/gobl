package nfe_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfe"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestLineValidation(t *testing.T) {
	tests := []struct {
		name string
		line *bill.Line
		err  string
	}{
		{
			name: "valid line with all required taxes",
			line: &bill.Line{
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryICMS,
					},
					{
						Category: br.TaxCategoryPIS,
					},
					{
						Category: br.TaxCategoryCOFINS,
					},
				},
			},
		},
		{
			name: "nil line",
			line: nil,
		},
		{
			name: "missing taxes",
			line: &bill.Line{},
			err:  "taxes: missing category ICMS.",
		},
		{
			name: "empty taxes",
			line: &bill.Line{
				Taxes: tax.Set{},
			},
			err: "taxes: missing category ICMS.",
		},
		{
			name: "missing ICMS tax",
			line: &bill.Line{
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryPIS,
					},
					{
						Category: br.TaxCategoryCOFINS,
					},
				},
			},
			err: "taxes: missing category ICMS.",
		},
		{
			name: "missing PIS tax",
			line: &bill.Line{
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryICMS,
					},
					{
						Category: br.TaxCategoryCOFINS,
					},
				},
			},
			err: "taxes: missing category PIS.",
		},
		{
			name: "missing COFINS tax",
			line: &bill.Line{
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryICMS,
					},
					{
						Category: br.TaxCategoryPIS,
					},
				},
			},
			err: "taxes: missing category COFINS.",
		},
	}

	addon := tax.AddonForKey(nfe.V4)
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := addon.Validator(ts.line)
			if ts.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), ts.err)
				}
			}
		})
	}
}
