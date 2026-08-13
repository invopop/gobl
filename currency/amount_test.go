package currency_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestAmountValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		a := currency.Amount{
			Currency: currency.USD,
			Value:    num.MakeAmount(100, 2),
		}
		assert.Nil(t, rules.Validate(&a))
	})
	t.Run("missing currency", func(t *testing.T) {
		a := currency.Amount{
			Value: num.MakeAmount(100, 2),
		}
		assert.NotNil(t, rules.Validate(&a))
	})
}
