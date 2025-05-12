package bill

import (
	"context"
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPaymentDetailsValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		p := &PaymentDetails{}
		assert.NoError(t, p.ValidateWithContext(context.Background()))
	})
}

func TestPaymentDetailsResetAdvances(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		p := &PaymentDetails{
			Advances: []*pay.Advance{
				{
					Description: "Paid in advance",
					Amount:      num.MakeAmount(10, 0),
				},
			},
		}
		p.ResetAdvances()
		assert.Empty(t, p.Advances)
	})
	t.Run("nil", func(t *testing.T) {
		var p *PaymentDetails
		assert.NotPanics(t, func() {
			p.ResetAdvances()
		})
	})
}

func TestPaymentDetailsNormalize(t *testing.T) {
	p := &PaymentDetails{
		Instructions: &pay.Instructions{
			Key:    "online",
			Detail: "Some random payment",
			Ext: tax.Extensions{
				"random": "",
			},
		},
	}
	p.Normalize(nil)
	assert.Empty(t, p.Instructions.Ext)
	assert.NotPanics(t, func() {
		p.Normalize(nil)
	})
}

func TestPaymentDetailsCalculations(t *testing.T) {
	zero := num.MakeAmount(0, 2)
	total := num.MakeAmount(20000, 2)
	p := &PaymentDetails{
		Advances: []*pay.Advance{
			{
				Description: "Paid in advance",
				Percent:     num.NewPercentage(10, 2),
			},
		},
	}
	p.calculateAdvances(zero, total)
	assert.Equal(t, "20.00", p.Advances[0].Amount.String())

	p = &PaymentDetails{
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

	p = &PaymentDetails{
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

	t.Run("nil", func(t *testing.T) {
		var p *PaymentDetails
		assert.Nil(t, p.totalAdvance(zero))
	})

	t.Run("uses currency precision", func(t *testing.T) {
		zero := num.MakeAmount(0, 2)
		total := num.MakeAmount(20845, 3)
		p := &PaymentDetails{
			Advances: []*pay.Advance{
				{
					Description: "Paid in advance",
					Percent:     num.NewPercentage(100, 2),
				},
			},
		}
		p.calculateAdvances(zero, total)
		a := p.totalAdvance(zero)
		assert.Equal(t, "20.85", p.Advances[0].Amount.String())
		assert.Equal(t, "20.85", a.String())
	})
}
