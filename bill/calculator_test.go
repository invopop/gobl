package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: many calculation tests are distributed throughout this package.

func TestRemoveIncludedTaxes(t *testing.T) {
	t.Run("no included tax", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
	})

	t.Run("from discounts", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Discounts = []*bill.Discount{
			{
				Amount: num.MakeAmount(1000, 2),
				Reason: "testing",
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     "standard",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
		assert.Equal(t, "8.26", inv.Totals.Discount.String())
	})

	t.Run("from charges", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Charges = []*bill.Charge{
			{
				Amount: num.MakeAmount(1000, 2),
				Reason: "testing",
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     "standard",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
		assert.Equal(t, "8.26", inv.Totals.Charge.String())
	})
}
