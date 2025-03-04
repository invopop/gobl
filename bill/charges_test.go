package bill

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChargeTotals(t *testing.T) {
	t.Run("base line", func(t *testing.T) {
		ls := []*Charge{
			{
				Reason: "First Charge",
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
		zero := num.MakeAmount(0, 2)
		base := num.MakeAmount(30000, 2)
		calculateCharges(ls, base, zero, tax.RoundingRuleSumThenRound)
		sum := calculateChargeSum(ls, zero)
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

		ls = []*Charge{}
		calculateCharges(ls, base, zero, tax.RoundingRuleSumThenRound)
		sum = calculateChargeSum(ls, zero)
		assert.Nil(t, sum)
	})

	t.Run("with precision", func(t *testing.T) {
		ls := []*Charge{
			{
				Reason: "First Charge",
				Amount: num.MakeAmount(50, 0),
			},
			{
				Reason:  "Second discount",
				Percent: num.NewPercentage(20, 2),
			},
		}
		zero := num.MakeAmount(0, 2)
		base := num.MakeAmount(30844212, 6)
		calculateCharges(ls, base, zero, tax.RoundingRuleSumThenRound)
		sum := calculateChargeSum(ls, zero)
		require.NotNil(t, sum)
		assert.Equal(t, "6.17", ls[1].Amount.String())
		assert.Equal(t, "56.168842", sum.String())
	})

	t.Run("with fixed base", func(t *testing.T) {
		ls := []*Charge{
			{
				Reason:  "Charge",
				Base:    num.NewAmount(5012, 2),
				Percent: num.NewPercentage(20, 2),
			},
		}
		zero := num.MakeAmount(0, 2)
		base := num.MakeAmount(30844212, 6)
		calculateCharges(ls, base, zero, tax.RoundingRuleSumThenRound)
		sum := calculateChargeSum(ls, zero)
		require.NotNil(t, sum)
		assert.Equal(t, "10.02", ls[0].Amount.String())
		assert.Equal(t, "10.02", sum.String())
	})

	t.Run("with fixed base high precision", func(t *testing.T) {
		ls := []*Charge{
			{
				Reason:  "Charge",
				Base:    num.NewAmount(501234, 4),
				Percent: num.NewPercentage(20, 2),
			},
		}
		zero := num.MakeAmount(0, 2)
		base := num.MakeAmount(30844212, 6)
		calculateCharges(ls, base, zero, tax.RoundingRuleSumThenRound)
		sum := calculateChargeSum(ls, zero)
		require.NotNil(t, sum)
		assert.Equal(t, "10.0247", ls[0].Amount.String())
		assert.Equal(t, "10.0247", sum.String())
	})
}
