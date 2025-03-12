package currency_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
)

func TestAmountValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		a := currency.Amount{
			Currency: currency.USD,
			Value:    num.MakeAmount(100, 2),
		}
		assert.NoError(t, a.Validate())
	})
	t.Run("missing currency", func(t *testing.T) {
		a := currency.Amount{
			Value: num.MakeAmount(100, 2),
		}
		assert.ErrorContains(t, a.Validate(), "currency: cannot be blank")
	})
}
