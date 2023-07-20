package bill

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutlayTotals(t *testing.T) {
	os := []*Outlay{
		{
			Description: "First outlay",
			Amount:      num.MakeAmount(10000, 2),
		},
		{
			Description: "Second outlay",
			Amount:      num.MakeAmount(200, 0),
		},
	}
	zero := num.MakeAmount(0, 2)
	sum := totalOutlays(zero, os)
	require.NotNil(t, sum)
	assert.Equal(t, 1, os[0].Index)
	assert.Equal(t, 2, os[1].Index)
	assert.Equal(t, "300.00", sum.String())
	assert.Equal(t, "200.00", os[1].Amount.String())

	os = []*Outlay{}
	sum = totalOutlays(zero, os)
	assert.Nil(t, sum)
}
