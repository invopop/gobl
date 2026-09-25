package facturae_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/addons/es/facturae"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPayDueDates(t *testing.T) {
	t.Run("with amount", func(t *testing.T) {
		p := &pay.Terms{
			DueDates: []*pay.DueDate{
				{
					Date:   cal.NewDate(2025, time.January, 1),
					Amount: num.NewAmount(1000, 2),
				},
			},
		}
		err := rules.Validate(p, tax.AddonContext(facturae.V3))
		assert.NoError(t, err)
	})

	t.Run("missing amount", func(t *testing.T) {
		p := &pay.Terms{
			DueDates: []*pay.DueDate{
				{
					Date: cal.NewDate(2025, time.January, 1),
				},
			},
		}
		err := rules.Validate(p, tax.AddonContext(facturae.V3))
		assert.ErrorContains(t, err, "amount is required")
	})

	t.Run("zero amount", func(t *testing.T) {
		p := &pay.Terms{
			DueDates: []*pay.DueDate{
				{
					Date:   cal.NewDate(2025, time.January, 1),
					Amount: num.NewAmount(0, 2),
				},
			},
		}
		err := rules.Validate(p, tax.AddonContext(facturae.V3))
		assert.ErrorContains(t, err, "amount is required")
	})

	t.Run("missing amount without addon", func(t *testing.T) {
		p := &pay.Terms{
			DueDates: []*pay.DueDate{
				{
					Date: cal.NewDate(2025, time.January, 1),
				},
			},
		}
		err := rules.Validate(p)
		assert.NoError(t, err)
	})
}
