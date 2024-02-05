package bill

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutlayUnmarshal(t *testing.T) {
	o := new(Outlay)
	err := json.Unmarshal([]byte(`{"desc":"foo"}`), o)
	require.NoError(t, err)
	assert.Equal(t, "foo", o.Description)
	err = json.Unmarshal([]byte(`{"description":"foo"}`), o)
	require.NoError(t, err)
	assert.Equal(t, "foo", o.Description)
}

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
	sum := calculateOutlays(zero, os)
	require.NotNil(t, sum)
	assert.Equal(t, 1, os[0].Index)
	assert.Equal(t, 2, os[1].Index)
	assert.Equal(t, "300.00", sum.String())
	assert.Equal(t, "200.00", os[1].Amount.String())

	os = []*Outlay{}
	sum = calculateOutlays(zero, os)
	assert.Nil(t, sum)
}
