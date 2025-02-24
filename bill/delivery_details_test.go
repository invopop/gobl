package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/require"
)

func TestDeliveryDetailsValidation(t *testing.T) {
	t.Run("check is used", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Delivery = &bill.DeliveryDetails{
			Date: cal.NewDate(2020, 1, 1),
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
}
