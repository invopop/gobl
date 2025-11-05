package sdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
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
	err := inv.Validate()
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
	err = inv.Validate()
	assert.ErrorContains(t, err, "payment: (advances: (0: (ext: (it-sdi-payment-means: required.).).).)")

	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{
			Key: pay.MeansKeyDirectDebit.With("fooo"),
		},
	}
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.ErrorContains(t, err, "payment: (instructions: (ext: (it-sdi-payment-means: required.).).)")

	// Valid IBAN - 27 characters (minimum)
	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{
			Key: pay.MeansKeyCreditTransfer,
			CreditTransfer: []*pay.CreditTransfer{
				{
					IBAN: "IT60X0542811101000000123456",
				},
			},
		},
	}
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	require.NoError(t, err)

	// Valid IBAN - 34 characters (maximum)
	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{
			Key: pay.MeansKeyCreditTransfer,
			CreditTransfer: []*pay.CreditTransfer{
				{
					IBAN: "IT60X0542811101000000123456ABCDEFG",
				},
			},
		},
	}
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	require.NoError(t, err)

	// Invalid IBAN - too short (26 characters)
	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{
			Key: pay.MeansKeyCreditTransfer,
			CreditTransfer: []*pay.CreditTransfer{
				{
					IBAN: "IT60X054281110100000012345",
				},
			},
		},
	}
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.ErrorContains(t, err, "iban: the length must be between 27 and 34")

	// Invalid IBAN - too long (35 characters)
	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{
			Key: pay.MeansKeyCreditTransfer,
			CreditTransfer: []*pay.CreditTransfer{
				{
					IBAN: "IT60X0542811101000000123456ABCDEFGH",
				},
			},
		},
	}
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.ErrorContains(t, err, "iban: the length must be between 27 and 34")

	// Empty IBAN with other fields set - should pass
	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{
			Key: pay.MeansKeyCreditTransfer,
			CreditTransfer: []*pay.CreditTransfer{
				{
					IBAN: "",
					BIC:  "ABCDITMM",
					Name: "Test Bank",
				},
			},
		},
	}
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	require.NoError(t, err)
}
