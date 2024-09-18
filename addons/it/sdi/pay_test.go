package sdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/sdi"
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
				Key:         pay.MeansKeyDirectDebit.With(sdi.MeansKeyRID),
				Description: "Test advance",
				Amount:      num.MakeAmount(100, 0),
			},
		},
	}
	require.NoError(t, inv.Calculate())
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
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.ErrorContains(t, err, "payment: (advances: (0: (ext: (it-sdi-payment-means: required.).).).)")
}
