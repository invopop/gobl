package gr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/gr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeTaxCombo(t *testing.T) {
	tests := []struct {
		name string
		rate cbc.Key
		vcat tax.ExtValue
	}{
		{
			name: "standard rate",
			rate: "standard",
			vcat: "1",
		},
		{
			name: "exempt rate",
			rate: "exempt",
			vcat: "7",
		},
		{
			name: "no rate",
			rate: "",
			vcat: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &tax.Combo{Category: tax.CategoryVAT, Rate: tt.rate}

			err := gr.Calculate(tc)
			require.NoError(t, err)

			vcat := tc.Ext[gr.ExtKeyIAPRVATCat]
			assert.Equal(t, tt.vcat, vcat)
		})
	}
}

func TestValidateTaxCombo(t *testing.T) {
	t.Run("vat category presence", func(t *testing.T) {
		tc := &tax.Combo{Category: tax.CategoryVAT, Percent: num.NewPercentage(4, 2)}

		err := gr.Calculate(tc)
		require.NoError(t, err)

		err = gr.Validate(tc)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "gr-iapr-vat-cat: required")
	})

	t.Run("exemption presence", func(t *testing.T) {
		tc := &tax.Combo{Category: tax.CategoryVAT, Rate: tax.RateExempt}

		err := gr.Calculate(tc)
		require.NoError(t, err)

		err = gr.Validate(tc)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "gr-iapr-exemption: required")
	})
}
