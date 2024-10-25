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
		  name: "nil",
		  item: nil,
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
