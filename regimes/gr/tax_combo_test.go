package gr_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/gr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/require"
)

func TestValidateTaxCombo(t *testing.T) {

	greece := gr.New()
	ctx := greece.WithContext(context.Background())

	t.Run("vat category presence", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(4, 2),
		}
		err := tc.ValidateWithContext(ctx)
		require.ErrorContains(t, err, "ext: (gr-iapr-vat-cat: required.)")
	})

	t.Run("exemption presence", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateExempt,
			Ext:      tax.Extensions{"gr-iapr-vat-cat": "1"},
		}
		err := tc.ValidateWithContext(ctx)
		require.ErrorContains(t, err, "gr-iapr-exemption: required")
	})
}
