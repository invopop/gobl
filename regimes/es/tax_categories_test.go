package es_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPSICategoryKeys(t *testing.T) {
	cd := es.New().CategoryDef(es.TaxCategoryIPSI)
	require.NotNil(t, cd)

	keys := make([]cbc.Key, 0, len(cd.Keys))
	for _, kd := range cd.Keys {
		keys = append(keys, kd.Key)
	}
	assert.ElementsMatch(t, []cbc.Key{
		tax.KeyStandard,
		tax.KeyZero,
		tax.KeyReverseCharge,
		tax.KeyExempt,
		tax.KeyExport,
		tax.KeyOutsideScope,
	}, keys)
	assert.Nil(t, cd.KeyDef(tax.KeyIntraCommunity), "intra-community has no meaning in Ceuta and Melilla")
}

func TestIPSICombos(t *testing.T) {
	newInvoice := func(tc *tax.Combo) *bill.Invoice {
		return &bill.Invoice{
			Code:     "IPSI-1",
			Currency: "EUR",
			Supplier: &org.Party{
				Name: "Ceuta Supplier",
				TaxID: &tax.Identity{
					Country: "ES",
					Code:    "B98602642",
				},
			},
			Customer: &org.Party{
				Name: "Customer",
				TaxID: &tax.Identity{
					Country: "ES",
					Code:    "54387763P",
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Service",
						Price: num.NewAmount(10000, 2),
					},
					Taxes: tax.Set{tc},
				},
			},
		}
	}

	tests := []struct {
		name    string
		combo   *tax.Combo
		err     string
		key     cbc.Key
		percent *string
	}{
		{
			name:  "percent only defaults to standard",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Percent: num.NewPercentage(4, 2)},
			key:   tax.KeyStandard,
		},
		{
			name:  "standard with percent",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyStandard, Percent: num.NewPercentage(4, 2)},
			key:   tax.KeyStandard,
		},
		{
			name:  "standard without percent",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyStandard},
			err:   "percent required or invalid",
		},
		{
			name:  "no key without percent",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI},
			err:   "percent required or invalid",
		},
		{
			name:  "zero with percent",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyZero, Percent: num.NewPercentage(0, 2)},
			key:   tax.KeyZero,
		},
		{
			name:  "zero without percent",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyZero},
			err:   "percent required or invalid",
		},
		{
			name:  "exempt without percent",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyExempt},
			key:   tax.KeyExempt,
		},
		{
			name:  "exempt clears percent",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyExempt, Percent: num.NewPercentage(4, 2)},
			key:   tax.KeyExempt,
		},
		{
			name:  "outside scope without percent",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyOutsideScope},
			key:   tax.KeyOutsideScope,
		},
		{
			name:  "reverse charge without percent",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyReverseCharge},
			key:   tax.KeyReverseCharge,
		},
		{
			name:  "export without percent",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyExport},
			key:   tax.KeyExport,
		},
		{
			name:  "intra-community rejected",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyIntraCommunity},
			err:   "tax combo key not valid for category in regime",
		},
		{
			name:  "unknown key rejected",
			combo: &tax.Combo{Category: es.TaxCategoryIPSI, Key: "bogus"},
			err:   "tax combo key not valid for category in regime",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := newInvoice(tt.combo)
			require.NoError(t, inv.Calculate())
			err := rules.Validate(inv)
			if tt.err != "" {
				assert.ErrorContains(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			tc := inv.Lines[0].Taxes[0]
			assert.Equal(t, tt.key, tc.Key)
			kd := es.New().CategoryDef(es.TaxCategoryIPSI).KeyDef(tc.Key)
			require.NotNil(t, kd)
			if kd.NoPercent {
				assert.Nil(t, tc.Percent)
			} else {
				assert.NotNil(t, tc.Percent)
			}
		})
	}
}
