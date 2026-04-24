package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/cbc"
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
				Name: "Test Item",
				Unit: "kg",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					saft.ExtKeyProductType: "P",
				}),
			},
		},
		{
			name: "nil item",
			item: nil,
		},
		{
			name: "missing extensions",
			item: &org.Item{},
			err:  "product type is required",
		},
		{
			name: "empty extensions",
			item: &org.Item{
				Ext: tax.ExtensionsOf(tax.ExtMap{}),
			},
			err: "product type is required",
		},
		{
			name: "missing extension",
			item: &org.Item{
				Ext: tax.ExtensionsOf(tax.ExtMap{
					"random": "12345678",
				}),
			},
			err: "product type is required",
		},
		{
			name: "missing unit",
			item: &org.Item{},
			err:  "cannot be blank",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := rules.Validate(ts.item, withAddonContext())
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
				Ext: tax.ExtensionsOf(tax.ExtMap{
					saft.ExtKeyProductType: "P",
				}),
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
				Ext: tax.ExtensionsOf(tax.ExtMap{}),
			},
			out: "S",
		},
		{
			name: "missing extension",
			item: &org.Item{
				Ext: tax.ExtensionsOf(tax.ExtMap{
					"random": "12345678",
				}),
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
				assert.Equal(t, tt.out, tt.item.Ext.Get(saft.ExtKeyProductType))
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
			out:  "one",
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
