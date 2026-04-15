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
	assert.NotEmpty(t, m)
	assert.Len(t, m, 24)
	assert.Equal(t, pay.MeansKeyCash, m.Lookup("MP01"))
}

func TestPayInstructionsNormalize(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{
			Key: "online",
			Ext: tax.Extensions{
				"random": "",
			},
		},
		Advances: []*pay.Advance{
			{
				Key:         pay.MeansKeyDirectDebit.With(sdi.MeansKeyRID),
				Description: "Test advance",
				Amount:      num.MakeAmount(100, 0),
				Ext: tax.Extensions{
					"random": "",
				},
			},
		},
	}
	_, ok := inv.Payment.Instructions.Ext["random"]
	assert.True(t, ok)
	_, ok = inv.Payment.Advances[0].Ext["random"]
	assert.True(t, ok)
	assert.NoError(t, inv.Calculate())
	_, ok = inv.Payment.Instructions.Ext["random"]
	assert.False(t, ok)
	_, ok = inv.Payment.Advances[0].Ext["random"]
	assert.False(t, ok)
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
				Key:         pay.MeansKeyDirectDebit.With("fooo"),
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
			Key: pay.MeansKeyDirectDebit.With("fooo"),
		},
	}
	require.NoError(t, inv.Calculate())
	err = rules.Validate(inv)
	assert.ErrorContains(t, err, fmt.Sprintf("payment instructions require '%s' extension", sdi.ExtKeyPaymentMeans))
}
