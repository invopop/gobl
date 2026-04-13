package nfse_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfse"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
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
				Name: "Test Service",
				Ext: tax.Extensions{
					nfse.ExtKeyService: "12345678",
				},
			},
		},
		{
			name: "valid item with all IBS/CBS extensions",
			item: &org.Item{
				Name: "Test Service",
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
			item: &org.Item{
				Name: "Test",
			},
			err: "item requires 'br-nfse-service' extension",
		},
		{
			name: "empty extensions",
			item: &org.Item{
				Name: "Test",
				Ext:  tax.Extensions{},
			},
			err: "item requires 'br-nfse-service' extension",
		},
		{
			name: "missing extension",
			item: &org.Item{
				Name: "Test",
				Ext: tax.Extensions{
					"random": "12345678",
				},
			},
			err: "item requires 'br-nfse-service' extension",
		},
		{
			name: "only operation extension",
			item: &org.Item{
				Name: "Test",
				Ext: tax.Extensions{
					nfse.ExtKeyService:   "12345678",
					nfse.ExtKeyOperation: "030101",
				},
			},
			err: "item extensions 'br-nfse-operation', 'br-nfse-tax-status', and 'br-nfse-tax-class' must all be present or all absent",
		},
		{
			name: "only tax status extension",
			item: &org.Item{
				Name: "Test",
				Ext: tax.Extensions{
					nfse.ExtKeyService:   "12345678",
					nfse.ExtKeyTaxStatus: "000",
				},
			},
			err: "item extensions 'br-nfse-operation', 'br-nfse-tax-status', and 'br-nfse-tax-class' must all be present or all absent",
		},
		{
			name: "only tax class extension",
			item: &org.Item{
				Name: "Test",
				Ext: tax.Extensions{
					nfse.ExtKeyService:  "12345678",
					nfse.ExtKeyTaxClass: "000001",
				},
			},
			err: "item extensions 'br-nfse-operation', 'br-nfse-tax-status', and 'br-nfse-tax-class' must all be present or all absent",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := rules.Validate(ts.item, withAddonContext())
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
