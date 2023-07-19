package it_test

import (
	"context"
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

	inv.Payment = &bill.Payment{
		Advances: []*pay.Advance{
			{
				Key:         pay.MeansKeyDirectDebit.With(it.MeansKeyRID),
				Description: "Test advance",
				Amount:      num.MakeAmount(100, 0),
			},
		},
	}
	ctx := context.Background()
	require.NoError(t, inv.Calculate(ctx))
	err := inv.Validate()
	require.NoError(t, err)

	inv.Payment = &bill.Payment{
		Advances: []*pay.Advance{
			{
				Key:         pay.MeansKeyDirectDebit.With("fooo"),
				Description: "Test advance",
				Amount:      num.MakeAmount(100, 0),
			},
		},
	}
	require.NoError(t, inv.Calculate(ctx))
	err = inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "key: must be a valid value")
}
