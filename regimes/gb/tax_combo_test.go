package gb_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxSetRules(t *testing.T) {
	p20 := num.MakePercentage(200, 3)
	gbCtx := tax.RegimeContext(l10n.GB.Tax())
	tests := []struct {
		name string
		set  tax.Set
		err  string
	}{
		{
			name: "valid single",
			set:  tax.Set{{Category: tax.CategoryVAT, Key: tax.KeyStandard, Percent: &p20}},
		},
		{
			name: "valid exempt",
			set:  tax.Set{{Category: tax.CategoryVAT, Key: tax.KeyExempt}},
		},
		{
			name: "duplicate category",
			set: tax.Set{
				{Category: tax.CategoryVAT, Key: tax.KeyStandard},
				{Category: tax.CategoryVAT, Key: tax.KeyExempt},
			},
			err: "TAX-SET-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.Validate(tt.set, gbCtx)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}

func TestTaxComboRules(t *testing.T) {
	p20 := num.MakePercentage(200, 3)
	pZero := num.MakePercentage(0, 2)
	gbCtx := tax.RegimeContext(l10n.GB.Tax())
	tests := []struct {
		name  string
		combo *tax.Combo
		err   string
	}{
		{
			name:  "valid standard+general",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: tax.KeyStandard, Rate: tax.RateGeneral, Percent: &p20},
		},
		{
			name:  "valid standard+reduced",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: tax.KeyStandard, Rate: tax.RateReduced, Percent: &p20},
		},
		{
			name:  "valid zero key",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: tax.KeyZero, Percent: &pZero},
		},
		{
			name:  "valid reverse-charge",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: tax.KeyReverseCharge},
		},
		{
			name:  "valid exempt",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: tax.KeyExempt},
		},
		{
			name:  "valid export",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: tax.KeyExport},
		},
		{
			name:  "invalid category",
			combo: &tax.Combo{Category: tax.CategoryGST},
			err:   "GOBL-TAX-COMBO-01",
		},
		{
			name:  "invalid key",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: "unknown-key"},
			err:   "GOBL-TAX-COMBO-02",
		},
		{
			name:  "rate on no-rate key",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: tax.KeyExempt, Rate: tax.RateGeneral},
			err:   "GOBL-TAX-COMBO-03",
		},
		{
			name:  "unknown rate on standard key",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: tax.KeyStandard, Rate: "special"},
			err:   "GOBL-TAX-COMBO-03",
		},
		{
			name:  "percent on reverse-charge key",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: tax.KeyReverseCharge, Percent: &p20},
			err:   "GOBL-TAX-COMBO-04",
		},
		{
			name:  "percent on exempt key",
			combo: &tax.Combo{Category: tax.CategoryVAT, Key: tax.KeyExempt, Percent: &p20},
			err:   "GOBL-TAX-COMBO-04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.Validate(tt.combo, gbCtx)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}
