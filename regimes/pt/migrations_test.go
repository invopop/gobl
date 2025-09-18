package pt_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxRateMigration(t *testing.T) {
	// Valid old rate
	inv := validInvoice()
	inv.SetAddons(saft.V1)
	inv.Lines[0].Taxes[0].Rate = "exempt+outlay"

	err := inv.Calculate()
	require.NoError(t, err)

	t0 := inv.Lines[0].Taxes[0]
	assert.Equal(t, tax.KeyExempt, t0.Key)
	assert.Equal(t, cbc.Code("M01"), t0.Ext[saft.ExtKeyExemption])

	// Valid new rate
	inv = validInvoice()
	inv.SetAddons(saft.V1)
	inv.Lines[0].Taxes[0].Rate = "exempt"
	inv.Lines[0].Taxes[0].Ext = tax.Extensions{saft.ExtKeyExemption: "M02"}

	err = inv.Calculate()
	require.NoError(t, err)

	t0 = inv.Lines[0].Taxes[0]
	assert.Equal(t, tax.KeyExempt, t0.Key)
	assert.Equal(t, cbc.Code("M02"), t0.Ext[saft.ExtKeyExemption])
}

func TestTaxZoneMigration(t *testing.T) {
	tests := []struct {
		name     string
		supplier *org.Party
		region   cbc.Code
	}{
		{
			name: "Azores zone set",
			supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: "PT",
					Zone:    "20", //nolint:staticcheck
				},
			},
			region: "PT-AC",
		},
		{
			name: "Madeira zone set",
			supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: "PT",
					Zone:    "30", //nolint:staticcheck
				},
			},
			region: "PT-MA",
		},
		{
			name: "Other zone set",
			supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: "PT",
					Zone:    "40", //nolint:staticcheck
				},
			},
			region: "PT",
		},
		{
			name: "No zone set",
			supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: "PT",
				},
			},
			region: "PT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := validInvoice()
			inv.Discounts = []*bill.Discount{{Taxes: tax.Set{{Category: tax.CategoryVAT}}}}
			inv.Charges = []*bill.Charge{{Taxes: tax.Set{{Category: tax.CategoryVAT}}}}
			inv.SetAddons(saft.V1)

			inv.Supplier = tt.supplier
			err := inv.Calculate()
			require.NoError(t, err)

			t0 := inv.Lines[0].Taxes[0]
			assert.Equal(t, tt.region, t0.Ext[pt.ExtKeyRegion])
			t0 = inv.Discounts[0].Taxes[0]
			assert.Equal(t, tt.region, t0.Ext[pt.ExtKeyRegion])
			t0 = inv.Charges[0].Taxes[0]
			assert.Equal(t, tt.region, t0.Ext[pt.ExtKeyRegion])
		})
	}
}
