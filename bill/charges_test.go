package bill

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChargeTotals(t *testing.T) {
	ls := []*Charge{
		{
			Reason: "First Charge",
			Amount: num.MakeAmount(100, 0),
		},
		{
			Reason:  "Second Charge",
			Percent: num.NewPercentage(20, 2),
		},
	}
	zero := num.MakeAmount(0, 2)
	base := num.MakeAmount(30000, 2)
	err := calculateCharges(ls, base, zero)
	require.NoError(t, err)
	sum := calculateChargeSum(ls, zero)
	require.NotNil(t, sum)
	assert.Equal(t, 1, ls[0].Index)
	assert.Equal(t, 2, ls[1].Index)
	assert.Equal(t, "160.00", sum.String())
	assert.Equal(t, "100.00", ls[0].Amount.String())
	assert.Nil(t, ls[1].Base)
	assert.Equal(t, "20%", ls[1].Percent.String())
	assert.Equal(t, "60.00", ls[1].Amount.String())

	ls = []*Charge{}
	require.NoError(t, calculateCharges(ls, base, zero))
	sum = calculateChargeSum(ls, zero)
	assert.Nil(t, sum)
}
