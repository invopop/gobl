package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/cbc"
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
				Unit: "kg",
				Ext: tax.Extensions{
					saft.ExtKeyProductType: "P",
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
			err:  "ext: (pt-saft-product-type: required.)",
		},
		{
			name: "empty extensions",
			item: &org.Item{
				Ext: tax.Extensions{},
			},
			err: "ext: (pt-saft-product-type: required.)",
		},
		{
			name: "missing extension",
			item: &org.Item{
				Ext: tax.Extensions{
					"random": "12345678",
				},
			},
			err: "ext: (pt-saft-product-type: required.)",
		},
		{
			name: "missing unit",
			item: &org.Item{},
			err:  "unit: cannot be blank.",
		},
	}

	addon := tax.AddonForKey(saft.V1)
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := addon.Validator(ts.item)
			if ts.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, ts.err)
			}
		})
	}
}

func TestItemExtProductTypeNormalization(t *testing.T) {
	tests := []struct {
		name string
		item *org.Item
		out  cbc.Code
	}{
		{
			name: "extension present",
			item: &org.Item{
				Ext: tax.Extensions{
					saft.ExtKeyProductType: "P",
				},
			},
			out: "P",
		},
		{
			name: "nil item",
			item: nil,
		},
		{
			name: "empty extensions",
			item: &org.Item{
				Ext: tax.Extensions{},
			},
			out: "S",
		},
		{
			name: "missing extension",
			item: &org.Item{
				Ext: tax.Extensions{
					"random": "12345678",
				},
			},
			out: "S",
		},
		{
			name: "goods unit set",
			item: &org.Item{
				Unit: "kg",
			},
			out: "P",
		},
		{
			name: "service unit set",
			item: &org.Item{
				Unit: "service",
			},
			out: "S",
		},
	}

	addon := tax.AddonForKey(saft.V1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addon.Normalizer(tt.item)
			if tt.item != nil {
				assert.Equal(t, tt.out, tt.item.Ext[saft.ExtKeyProductType])
			}
		})
	}
}

func TestItemUnitNormalization(t *testing.T) {
	tests := []struct {
		name string
		item *org.Item
		out  org.Unit
	}{
		{
			name: "unit present",
			item: &org.Item{
				Unit: "kg",
			},
			out: "kg",
		},
		{
			name: "nil item",
			item: nil,
		},
		{
			name: "unit not present",
			item: &org.Item{},
			out:  "item",
		},
	}

	addon := tax.AddonForKey(saft.V1)
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			addon.Normalizer(ts.item)
			if ts.out != "" {
				assert.Equal(t, ts.out, ts.item.Unit)
			}
		})
	}
}
