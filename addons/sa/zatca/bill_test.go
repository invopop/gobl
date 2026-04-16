package zatca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/sa/zatca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeInvoice(t *testing.T) {
	ad := tax.AddonForKey(zatca.V1)

	t.Run("nil invoice does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ad.Normalizer((*bill.Invoice)(nil))
		})
	})

	t.Run("sets rounding to currency", func(t *testing.T) {
		inv := validStandardInvoice()
		ad.Normalizer(inv)
		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
	})

	t.Run("creates tax object when nil", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Tax = nil
		ad.Normalizer(inv)
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
	})

	t.Run("creates issue time when nil", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.IssueTime = nil
		ad.Normalizer(inv)
		require.NotNil(t, inv.IssueTime)
	})

	t.Run("outside scope line gets zero percent", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Out of scope item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyOutsideScope,
					},
				},
			},
		}
		ad.Normalizer(inv)
		vat := inv.Lines[0].Taxes.Get(tax.CategoryVAT)
		require.NotNil(t, vat.Percent)
		assert.True(t, vat.Percent.IsZero())
	})

	t.Run("line without VAT combo is skipped", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "No tax item",
					Price: num.NewAmount(50, 0),
				},
				Taxes: tax.Set{},
			},
		}
		assert.NotPanics(t, func() {
			ad.Normalizer(inv)
		})
	})
}

func TestBillDiscountRules(t *testing.T) {
	t.Run("discount with taxes is valid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Discounts = []*bill.Discount{
			{
				Reason: "Loyalty discount",
				Amount: num.MakeAmount(50, 0),
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("discount without taxes fails", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Discounts = []*bill.Discount{
			{
				Reason: "Loyalty discount",
				Amount: num.MakeAmount(50, 0),
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "taxes are required (BR-32)")
	})
}
