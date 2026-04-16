package zatca_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Rule 01 (BR-KSA-35): if End is set, Start must also be set ---

func TestCalPeriodRule01_StartRequiredWhenEndPresent(t *testing.T) {
	t.Run("both start and end present is valid", func(t *testing.T) {
		inv := validSummaryInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("end set but start missing fires BR-KSA-35", func(t *testing.T) {
		inv := validSummaryInvoice()
		inv.Delivery.Period.Start = cal.Date{}
		_ = inv.Calculate()
		err := rules.Validate(inv)
		assert.ErrorContains(t, err,
			"if the invoice has a supply end date, it must also have a start date (BR-KSA-35)")
	})
}
