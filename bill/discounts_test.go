package bill

import (
	"context"
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscountValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		d := &Discount{
			Amount: num.MakeAmount(100, 2),
		}
		err := d.ValidateWithContext(context.Background())
		require.Nil(t, err)
	})
	t.Run("base with percent", func(t *testing.T) {
		d := &Discount{
			Base:   num.NewAmount(1000, 2),
			Amount: num.MakeAmount(100, 2),
		}
		err := d.ValidateWithContext(context.Background())
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "percent: cannot be blank")
	})
}

func TestDiscountTotals(t *testing.T) {
	t.Run("base line", func(t *testing.T) {
		ls := []*Discount{
			{
				Reason: "First Discount",
				Amount: num.MakeAmount(100, 0),
			},
			{
				Reason:  "Second discount",
				Percent: num.NewPercentage(20, 2),
			},
			{
				Reason:  "Third discount",
				Base:    num.NewAmount(20000, 2),
				Percent: num.NewPercentage(20, 2),
			},
		}
		base := num.MakeAmount(30000, 2)
		calculateDiscounts(ls, currency.EUR, base, tax.RoundingRuleSumThenRound)
		sum := calculateDiscountSum(ls, currency.EUR)
		require.NotNil(t, sum)
		assert.Equal(t, 1, ls[0].Index)
		assert.Nil(t, ls[0].Base)
		assert.Equal(t, 2, ls[1].Index)
		assert.Equal(t, "200.00", sum.String())
		assert.Equal(t, "100.00", ls[0].Amount.String())
		assert.Nil(t, ls[1].Base)
		assert.Equal(t, "20%", ls[1].Percent.String())
		assert.Equal(t, "60.00", ls[1].Amount.String())
		assert.Equal(t, "200.00", ls[2].Base.String())
		assert.Equal(t, "40.00", ls[2].Amount.String())

		ls = []*Discount{}
		calculateDiscounts(ls, currency.EUR, base, tax.RoundingRuleSumThenRound)
		sum = calculateDiscountSum(ls, currency.EUR)
		assert.Nil(t, sum)
	})

	t.Run("with precision", func(t *testing.T) {
		ls := []*Discount{
			{
				Reason: "First Discount",
				Amount: num.MakeAmount(50, 0),
			},
			{
				Reason:  "Second discount",
				Percent: num.NewPercentage(20, 2),
			},
		}
		base := num.MakeAmount(30844212, 6)
		calculateDiscounts(ls, currency.EUR, base, tax.RoundingRuleSumThenRound)
		sum := calculateDiscountSum(ls, currency.EUR)
		require.NotNil(t, sum)
		assert.Equal(t, "50.00", ls[0].Amount.String())
		assert.Equal(t, "6.17", ls[1].Amount.String())
		assert.Equal(t, "6.168842", ls[1].amount.String())
		assert.Equal(t, "56.168842", sum.String())
	})

	t.Run("with precision, round-then-sum", func(t *testing.T) {
		ls := []*Discount{
			{
				Reason: "First Discount",
				Amount: num.MakeAmount(50, 0),
			},
			{
				Reason:  "Second discount",
				Percent: num.NewPercentage(20, 2),
			},
		}
		base := num.MakeAmount(30844212, 6)
		calculateDiscounts(ls, currency.EUR, base, tax.RoundingRuleRoundThenSum)
		sum := calculateDiscountSum(ls, currency.EUR)
		require.NotNil(t, sum)
		assert.Equal(t, "50.00", ls[0].amount.String())
		assert.Equal(t, "50.00", ls[0].Amount.String())
		assert.Equal(t, "6.17", ls[1].Amount.String())
		assert.Equal(t, "6.17", ls[1].amount.String())
		assert.Equal(t, "56.17", sum.String())
	})

	t.Run("with fixed base", func(t *testing.T) {
		ls := []*Discount{
			{
				Reason:  "Discount",
				Base:    num.NewAmount(5012, 2),
				Percent: num.NewPercentage(20, 2),
			},
		}
		base := num.MakeAmount(30844212, 6)
		calculateDiscounts(ls, currency.EUR, base, tax.RoundingRuleSumThenRound)
		sum := calculateDiscountSum(ls, currency.EUR)
		require.NotNil(t, sum)
		assert.Equal(t, "10.02", ls[0].Amount.String())
		assert.Equal(t, "10.02", sum.String())
	})

	t.Run("with fixed amount", func(t *testing.T) {
		ls := []*Discount{
			{
				Reason: "Discount",
				Amount: num.MakeAmount(501762, 4),
			},
		}
		base := num.MakeAmount(30844212, 6)
		calculateDiscounts(ls, currency.EUR, base, tax.RoundingRuleSumThenRound)
		sum := calculateDiscountSum(ls, currency.EUR)
		require.NotNil(t, sum)
		assert.Equal(t, "50.18", ls[0].amount.String())
		assert.Equal(t, "50.18", ls[0].Amount.String())
		assert.Equal(t, "50.18", sum.String())
	})

	t.Run("with fixed base high precision", func(t *testing.T) {
		ls := []*Discount{
			{
				Reason:  "Discount",
				Base:    num.NewAmount(501234, 4),
				Percent: num.NewPercentage(20, 2),
			},
		}
		base := num.MakeAmount(30844212, 6)
		calculateDiscounts(ls, currency.EUR, base, tax.RoundingRuleSumThenRound)
		sum := calculateDiscountSum(ls, currency.EUR)
		require.NotNil(t, sum)
		assert.Equal(t, "10.0247", ls[0].amount.String())
		assert.Equal(t, "10.02", ls[0].Amount.String())
		assert.Equal(t, "10.0247", sum.String())
	})

	t.Run("with fixed base high precision, round-then-sum", func(t *testing.T) {
		ls := []*Discount{
			{
				Reason:  "Discount",
				Base:    num.NewAmount(501234, 4),
				Percent: num.NewPercentage(20, 2),
			},
		}
		base := num.MakeAmount(30844212, 6)
		calculateDiscounts(ls, currency.EUR, base, tax.RoundingRuleRoundThenSum)
		sum := calculateDiscountSum(ls, currency.EUR)
		require.NotNil(t, sum)
		assert.Equal(t, "50.1234", ls[0].Base.String())
		assert.Equal(t, "10.02", ls[0].amount.String())
		assert.Equal(t, "10.02", ls[0].Amount.String())
		assert.Equal(t, "10.02", sum.String())
	})
}
