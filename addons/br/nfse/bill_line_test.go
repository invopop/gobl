package nfse_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfse"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
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
			name: "valid line",
			line: &bill.Line{
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryISS,
						Ext: tax.Extensions{
							nfse.ExtKeyISSLiability: "1",
						},
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
			err:  "line taxes must include the ISS category",
		},
		{
			name: "empty taxes",
			line: &bill.Line{
				Taxes: tax.Set{},
			},
			err: "line taxes must include the ISS category",
		},
		{
			name: "missing ISS tax",
			line: &bill.Line{
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryPIS,
					},
				},
			},
			err: "line taxes must include the ISS category",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := rules.Validate(ts.line, withAddonContext())
			if ts.err != "" {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), ts.err)
				}
			}
		})
	}
}
