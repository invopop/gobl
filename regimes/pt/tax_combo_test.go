package pt_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxComboValidation(t *testing.T) {
	tests := []struct {
		name string
		tc   *tax.Combo
		err  string
	}{
		{
			name: "valid combo",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Percent:  num.NewPercentage(210, 3),
				Ext: tax.Extensions{
					pt.ExtKeyRegion: "PT-AC",
				},
			},
		},
		{
			name: "nil combo",
			tc:   nil,
		},
		{
			name: "missing extensions",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Percent:  num.NewPercentage(210, 3),
			},
			err: "[GOBL-PT-TAX-COMBO-01]",
		},
		{
			name: "empty extensions",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Percent:  num.NewPercentage(210, 3),
				Ext:      tax.Extensions{},
			},
			err: "[GOBL-PT-TAX-COMBO-01]",
		},
		{
			name: "missing extension",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Percent:  num.NewPercentage(210, 3),
				Ext: tax.Extensions{
					"random": "12345678",
				},
			},
			err: "[GOBL-PT-TAX-COMBO-01]",
		},
		{
			name: "other category",
			tc: &tax.Combo{
				Category: tax.CategoryGST,
				Percent:  num.NewPercentage(210, 3),
			},
			err: "[GOBL-TAX-COMBO-01]", // general combo error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.Validate(tt.tc, tax.RegimeContext(l10n.PT.Tax()))
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}

func TestTaxComboNormalization(t *testing.T) {
	tests := []struct {
		name string
		tc   *tax.Combo
		out  cbc.Code
	}{
		{
			name: "extension present",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Ext: tax.Extensions{
					pt.ExtKeyRegion: "PT-AC",
				},
			},
			out: "PT-AC",
		},
		{
			name: "nil combo",
			tc:   nil,
		},
		{
			name: "empty extensions",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Ext:      tax.Extensions{},
			},
			out: "PT",
		},
		{
			name: "missing extension",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Ext: tax.Extensions{
					"random": "12345678",
				},
			},
			out: "PT",
		},
		{
			name: "foreign tax",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Country:  l10n.EL.Tax(),
			},
			out: "GR",
		},
		{
			name: "foreign tax override",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Country:  l10n.ES.Tax(),
				Ext: tax.Extensions{
					pt.ExtKeyRegion: "PT",
				},
			},
			out: "ES",
		},
		{
			name: "foreign tax EU",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Country:  l10n.EU.Tax(),
			},
			out: "",
		},
		{
			name: "other category",
			tc: &tax.Combo{
				Category: tax.CategoryGST,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt.Normalize(tt.tc)
			if tt.tc != nil {
				assert.Equal(t, tt.out, tt.tc.Ext[pt.ExtKeyRegion])
			}
		})
	}
}
