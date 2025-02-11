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
		key  cbc.Key
		ext  tax.Extensions
		out  cbc.Code
	}{
		{
			name: "card, no ext",
			key:  pay.MeansKeyCard,
			out:  "CC",
		},
		{
			name: "card, ext",
			key:  pay.MeansKeyCard,
			ext: tax.Extensions{
				saft.ExtKeyPaymentMeans: "CB",
			},
			out: "CC",
		},
		{
			name: "other, no ext",
			key:  pay.MeansKeyOther,
			out:  "OU",
		},
		{
			name: "other, ext",
			key:  pay.MeansKeyOther,
			ext: tax.Extensions{
				saft.ExtKeyPaymentMeans: "CB",
			},
			out: "CB",
		},
	}

	addon := tax.AddonForKey(saft.V1)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := &pay.Advance{
				Key: test.key,
				Ext: test.ext,
			}
			addon.Normalizer(a)
			assert.Equal(t, test.out, a.Ext[saft.ExtKeyPaymentMeans])
		})
	}
}
