package bill

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
)

func TestPaymentCalculations(t *testing.T) {
	zero := num.MakeAmount(0, 2)
	total := num.MakeAmount(20000, 2)
	p := &Payment{
		Advances: []*pay.Advance{
			{
				Description: "Paid in advance",
				Percent:     num.NewPercentage(10, 2),
			},
		},
	}
	p.calculateAdvances(zero, total)
	assert.Equal(t, "20.00", p.Advances[0].Amount.String())

	p = &Payment{
		Advances: []*pay.Advance{
			{
				Description: "Paid in advance",
				Amount:      num.MakeAmount(10, 0),
			},
		},
	}
	assert.Equal(t, "10", p.Advances[0].Amount.String())
	p.calculateAdvances(zero, total)
	assert.Equal(t, "10.00", p.Advances[0].Amount.String())
	ta := p.totalAdvance(zero)
	assert.Equal(t, "10.00", ta.String())

	p = &Payment{
		Advances: []*pay.Advance{
			{
				Description: "Paid in advance",
				Amount:      num.MakeAmount(10, 0),
			},
			{
				Description: "Paid in advance %",
				Percent:     num.NewPercentage(10, 2),
			},
		},
	}
	p.calculateAdvances(zero, total)
	sum := p.totalAdvance(zero)
	assert.Equal(t, "30.00", sum.String())
}
