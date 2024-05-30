package it_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/it"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayInstructionsValidation(t *testing.T) {
	inv := testInvoiceStandard(t)

	inv.Payment = &bill.InvoicePayment{
		Advances: []*pay.Advance{
			{
				Key:         pay.MeansKeyDirectDebit.With(it.MeansKeyRID),
				Description: "Test advance",
				Amount:      num.MakeAmount(100, 0),
			},
		},
	}
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.NoError(t, err)

	inv.Payment = &bill.InvoicePayment{
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
	require.Error(t, err)
	assert.Contains(t, err.Error(), "key: must be a valid value")
}
