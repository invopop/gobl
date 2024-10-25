package nfse_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfse"
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
			name: "valid line",
			line: &bill.Line{
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryISS,
					},
				},
			},
		},
		{
			name: "missing taxes",
			line: &bill.Line{},
			err:  "taxes: missing category ISS.",
		},
		{
			name: "empty taxes",
			line: &bill.Line{
				Taxes: tax.Set{},
			},
			err: "taxes: missing category ISS.",
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
			err: "taxes: missing category ISS.",
		},
	}

	addon := tax.AddonForKey(nfse.V1)
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
