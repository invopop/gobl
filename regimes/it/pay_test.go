package it_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayInstructionsValidation(t *testing.T) {
	inv := testInvoiceStandard(t)

	inv.Payment = &bill.Payment{
		Advances: []*pay.Advance{
			{
				Key:         pay.MeansKeyOther,
				Description: "Test advance",
				Amount:      num.MakeAmount(100, 0),
			},
		},
	}
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "payment: (advances: (0: (code: cannot be blank.).).).")

	inv.Payment = &bill.Payment{
		Instructions: &pay.Instructions{
			Key:  pay.MeansKeyOther,
			Code: "",
		},
	}
	err = inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "payment: (instructions: (code: cannot be blank.).).")
}
