package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPaymentMeansExtensions(t *testing.T) {
	m := saft.PaymentMeansExtensions()
	assert.NotEmpty(t, m)
	assert.Len(t, m, 10)
	assert.Equal(t, pay.MeansKeyCash, m.Lookup("NU"))
}

func TestPayInstructionsNormalization(t *testing.T) {
	tests := []struct {
		name  string
		instr *pay.Instructions
		out   cbc.Code
	}{
		{
			name: "nil",
		},
		{
			name: "card, no ext",
			instr: &pay.Instructions{
				Key: pay.MeansKeyCard,
			},
			out: "CC",
		},
		{
			name: "card, ext",
			instr: &pay.Instructions{
				Key: pay.MeansKeyCard,
				Ext: tax.Extensions{
					saft.ExtKeyPaymentMeans: "CB",
				},
			},
			out: "CC",
		},
		{
			name: "other, no ext",
			instr: &pay.Instructions{
				Key: pay.MeansKeyOther,
			},
			out: "OU",
		},
		{
			name: "other, ext",
			instr: &pay.Instructions{
				Key: pay.MeansKeyOther,
				Ext: tax.Extensions{
					saft.ExtKeyPaymentMeans: "CB",
				},
			},
			out: "CB",
		},
	}

	addon := tax.AddonForKey(saft.V1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addon.Normalizer(tt.instr)
			if tt.instr == nil {
				// Nothing to check. Not panicking is enough.
				return
			}
			assert.Equal(t, tt.out, tt.instr.Ext[saft.ExtKeyPaymentMeans])
		})
	}
}

func TestPayAdvanceNormalization(t *testing.T) {
	tests := []struct {
		name string
		adv  *pay.Advance
		out  cbc.Code
	}{
		{
			name: "nil",
		},
		{
			name: "card, no ext",
			adv: &pay.Advance{
				Key: pay.MeansKeyCard,
			},
			out: "CC",
		},
		{
			name: "card, ext",
			adv: &pay.Advance{
				Key: pay.MeansKeyCard,
				Ext: tax.Extensions{
					saft.ExtKeyPaymentMeans: "CB",
				},
			},
			out: "CC",
		},
		{
			name: "other, no ext",
			adv: &pay.Advance{
				Key: pay.MeansKeyOther,
			},
			out: "OU",
		},
		{
			name: "other, ext",
			adv: &pay.Advance{
				Key: pay.MeansKeyOther,
				Ext: tax.Extensions{
					saft.ExtKeyPaymentMeans: "CB",
				},
			},
			out: "CB",
		},
	}

	addon := tax.AddonForKey(saft.V1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addon.Normalizer(tt.adv)
			if tt.adv == nil {
				// Nothing to check. Not panicking is enough.
				return
			}
			assert.Equal(t, tt.out, tt.adv.Ext[saft.ExtKeyPaymentMeans])
		})
	}
}
