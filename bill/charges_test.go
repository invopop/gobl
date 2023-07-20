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
	sum := chargeTotal(zero, base, ls)
	require.NotNil(t, sum)
	assert.Equal(t, 1, ls[0].Index)
	assert.Equal(t, 2, ls[1].Index)
	assert.Equal(t, "160.00", sum.String())
	assert.Equal(t, "100.00", ls[0].Amount.String())
	assert.Equal(t, "300.00", ls[1].Base.String())
	assert.Equal(t, "20%", ls[1].Percent.String())
	assert.Equal(t, "60.00", ls[1].Amount.String())

	ls = []*Charge{}
	sum = chargeTotal(zero, base, ls)
	assert.Nil(t, sum)
}
