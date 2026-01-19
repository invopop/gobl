package nfse_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfse"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestItemValidation(t *testing.T) {
	tests := []struct {
		name string
		item *org.Item
		err  string
	}{
		{
			name: "valid item",
			item: &org.Item{
				Ext: tax.Extensions{
					nfse.ExtKeyService: "12345678",
				},
			},
		},
		{
			name: "valid item with all IBS/CBS extensions",
			item: &org.Item{
				Ext: tax.Extensions{
					nfse.ExtKeyService:   "12345678",
					nfse.ExtKeyOperation: "030101",
					nfse.ExtKeyTaxStatus: "000",
					nfse.ExtKeyTaxClass:  "000001",
				},
			},
		},
		{
			name: "nil item",
			item: nil,
		},
		{
			name: "missing extensions",
			item: &org.Item{},
			err:  "ext: (br-nfse-service: required.)",
		},
		{
			name: "empty extensions",
			item: &org.Item{
				Ext: tax.Extensions{},
			},
			err: "ext: (br-nfse-service: required.)",
		},
		{
			name: "missing extension",
			item: &org.Item{
				Ext: tax.Extensions{
					"random": "12345678",
				},
			},
			err: "ext: (br-nfse-service: required.).",
		},
		{
			name: "only operation extension",
			item: &org.Item{
				Ext: tax.Extensions{
					nfse.ExtKeyService:   "12345678",
					nfse.ExtKeyOperation: "030101",
				},
			},
			err: "ext: (br-nfse-tax-class: required; br-nfse-tax-status: required.)",
		},
		{
			name: "only tax status extension",
			item: &org.Item{
				Ext: tax.Extensions{
					nfse.ExtKeyService:   "12345678",
					nfse.ExtKeyTaxStatus: "000",
				},
			},
			err: "ext: (br-nfse-operation: required; br-nfse-tax-class: required.)",
		},
		{
			name: "only tax class extension",
			item: &org.Item{
				Ext: tax.Extensions{
					nfse.ExtKeyService:  "12345678",
					nfse.ExtKeyTaxClass: "000001",
				},
			},
			err: "ext: (br-nfse-operation: required; br-nfse-tax-status: required.)",
		},
	}

	addon := tax.AddonForKey(nfse.V1)
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := addon.Validator(ts.item)
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
