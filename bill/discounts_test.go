package bill

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscountTotals(t *testing.T) {
	ls := []*Discount{
		{
			Reason: "First Discount",
			Amount: num.MakeAmount(100, 0),
		},
		{
			Reason:  "Second discount",
			Percent: num.NewPercentage(20, 2),
		},
	}
	zero := num.MakeAmount(0, 2)
	base := num.MakeAmount(30000, 2)
	err := calculateDiscounts(zero, base, ls)
	require.NoError(t, err)
	sum := calculateDiscountSum(zero, ls)
	require.NotNil(t, sum)
	assert.Equal(t, 1, ls[0].Index)
	assert.Equal(t, 2, ls[1].Index)
	assert.Equal(t, "160.00", sum.String())
	assert.Equal(t, "100.00", ls[0].Amount.String())
	assert.Equal(t, "300.00", ls[1].Base.String())
	assert.Equal(t, "20%", ls[1].Percent.String())
	assert.Equal(t, "60.00", ls[1].Amount.String())

	ls = []*Discount{}
	require.NoError(t, calculateDiscounts(zero, base, ls))
	sum = calculateDiscountSum(zero, ls)
	assert.Nil(t, sum)
}
