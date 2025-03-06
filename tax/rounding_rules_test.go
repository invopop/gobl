package tax_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestApplyRoundingRule(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		a := num.MakeAmount(512462, 4)
		a = tax.ApplyRoundingRule(tax.RoundingRuleSumThenRound, "USD", a)
		assert.Equal(t, "51.2462", a.String())
	})
	t.Run("default, less precision", func(t *testing.T) {
		a := num.MakeAmount(51, 0)
		a = tax.ApplyRoundingRule(tax.RoundingRuleSumThenRound, "USD", a)
		assert.Equal(t, "51.00", a.String())
	})
	t.Run("round-then-sum", func(t *testing.T) {
		a := num.MakeAmount(512462, 4)
		a = tax.ApplyRoundingRule(tax.RoundingRuleRoundThenSum, "USD", a)
		assert.Equal(t, "51.25", a.String())
	})
	t.Run("round-then-sum, less precision", func(t *testing.T) {
		a := num.MakeAmount(51, 0)
		a = tax.ApplyRoundingRule(tax.RoundingRuleRoundThenSum, "USD", a)
		assert.Equal(t, "51.00", a.String())
	})
}
