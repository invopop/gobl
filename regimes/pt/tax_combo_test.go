package pt_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/pt"
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
			},
			err: "ext: (pt-region: required.)",
		},
		{
			name: "empty extensions",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Ext:      tax.Extensions{},
			},
			err: "ext: (pt-region: required.)",
		},
		{
			name: "missing extension",
			tc: &tax.Combo{
				Category: tax.CategoryVAT,
				Ext: tax.Extensions{
					"random": "12345678",
				},
			},
			err: "ext: (pt-region: required.)",
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
			err := pt.Validate(tt.tc)
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
