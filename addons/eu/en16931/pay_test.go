package en16931_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayInstructions(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)

	t.Run("valid", func(t *testing.T) {
		m := &pay.Instructions{
			Key: pay.MeansKeyCreditTransfer,
		}
		ad.Normalizer(m)
		assert.Equal(t, "30", m.Ext[untdid.ExtKeyPaymentMeans].String())
	})

	t.Run("nil", func(t *testing.T) {
		var m *pay.Instructions
		assert.NotPanics(t, func() {
			ad.Normalizer(m)
		})
	})

	t.Run("validation", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "30", inv.Payment.Instructions.Ext[untdid.ExtKeyPaymentMeans].String())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}
