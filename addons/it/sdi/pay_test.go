package sdi_test

import (
	"fmt"
	"testing"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentMeansExtensions(t *testing.T) {
	m := sdi.PaymentMeansExtensions()
	assert.False(t, m.IsZero())
	assert.Equal(t, 24, m.Len())
	assert.Equal(t, pay.MeansKeyCash, m.Lookup("MP01"))
}

func TestPayInstructionsNormalize(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{
			Key: "online",
			Ext: tax.ExtensionsOf(tax.ExtMap{
				"random": "",
			}),
		},
		Advances: []*pay.Advance{
			{
				Key:         pay.MeansKeyDirectDebit.With(sdi.MeansKeyRID),
				Description: "Test advance",
				Amount:      num.MakeAmount(100, 0),
				Ext: tax.ExtensionsOf(tax.ExtMap{
					"random": "",
				}),
			},
		},
	}
	assert.True(t, inv.Payment.Instructions.Ext.Has("random"))
	assert.True(t, inv.Payment.Advances[0].Ext.Has("random"))
	assert.NoError(t, inv.Calculate())
	assert.False(t, inv.Payment.Instructions.Ext.Has("random"))
	assert.False(t, inv.Payment.Advances[0].Ext.Has("random"))
}

func TestPayInstructionsCardSubKeysFallback(t *testing.T) {
	addon := tax.AddonForKey(sdi.V1)

	t.Run("bare card", func(t *testing.T) {
		instr := &pay.Instructions{Key: pay.MeansKeyCard}
		addon.Normalizer(instr)
		assert.Equal(t, "MP08", instr.Ext.Get(sdi.ExtKeyPaymentMeans).String())
	})

	t.Run("card+credit falls back to MP08", func(t *testing.T) {
		instr := &pay.Instructions{Key: pay.MeansKeyCard.With(pay.MeansKeyCredit)}
		addon.Normalizer(instr)
		assert.Equal(t, "MP08", instr.Ext.Get(sdi.ExtKeyPaymentMeans).String())
	})

	t.Run("card+debit falls back to MP08", func(t *testing.T) {
		instr := &pay.Instructions{Key: pay.MeansKeyCard.With(pay.MeansKeyDebit)}
		addon.Normalizer(instr)
		assert.Equal(t, "MP08", instr.Ext.Get(sdi.ExtKeyPaymentMeans).String())
	})
}

func TestPayInstructionsValidation(t *testing.T) {
	inv := testInvoiceStandard(t)

	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{
			Key: "cash",
		},
		Advances: []*pay.Advance{
			{
				Key:         pay.MeansKeyDirectDebit.With(sdi.MeansKeyRID),
				Description: "Test advance",
				Amount:      num.MakeAmount(100, 0),
			},
		},
	}
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	require.NoError(t, err)

	inv.Payment = &bill.PaymentDetails{
		Advances: []*pay.Advance{
			{
				Key:         pay.MeansKeyAny.With("fooo"),
				Description: "Test advance",
				Amount:      num.MakeAmount(100, 0),
			},
		},
	}
	require.NoError(t, inv.Calculate())
	err = rules.Validate(inv)
	assert.ErrorContains(t, err, fmt.Sprintf("payment advance requires '%s' extension", sdi.ExtKeyPaymentMeans))

	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{
			Key: pay.MeansKeyAny.With("fooo"),
		},
	}
	require.NoError(t, inv.Calculate())
	err = rules.Validate(inv)
	assert.ErrorContains(t, err, fmt.Sprintf("payment instructions require '%s' extension", sdi.ExtKeyPaymentMeans))
}
