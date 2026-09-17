package dian_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBillLineQuantityValidation(t *testing.T) {
	t.Run("rejects zero quantities", func(t *testing.T) {
		inv := baseInvoice()
		inv.Lines[0].Quantity = num.MakeAmount(0, 0)
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "GOBL-CO-DIAN-BILL-LINE-01")
	})

	t.Run("rejects negative quantities", func(t *testing.T) {
		inv := baseInvoice()
		inv.Lines[0].Quantity = num.MakeAmount(-1, 0)
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "GOBL-CO-DIAN-BILL-LINE-01")
	})
}
