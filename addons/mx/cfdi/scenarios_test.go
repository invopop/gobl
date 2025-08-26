package cfdi_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceScenarios(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "PPD", inv.Tax.Ext[cfdi.ExtKeyPaymentMethod].String())
		assert.NoError(t, inv.Validate())
	})
	t.Run("prepaid", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Key:         "card",
					Description: "Pago anticipado",
					Percent:     num.NewPercentage(100, 2),
				},
			},
		}
		require.NoError(t, inv.Calculate())
		data, _ := json.MarshalIndent(inv, "", "  ")
		t.Logf("DOC: %s", string(data))
		assert.NoError(t, inv.Validate())
		assert.Equal(t, "PUE", inv.Tax.Ext[cfdi.ExtKeyPaymentMethod].String())
	})
}
