package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// RoundingRuleSumThenRound is the default method of calculating the totals
	// in GOBL, and provides the best results for most cases as the precision
	// is maintained to the maximum amount possible. The tradeoff however is
	// that sometimes the totals may not sum exactly based on what is visible.
	RoundingRuleSumThenRound cbc.Key = "sum-then-round"

	// RoundingRuleRoundThenSum is the alternative method of calculating the totals
	// that will first round all the amounts to the currency's precision before
	// making the sums. Totals using this approach can always be recalculated using
	// the amounts presented, but can lead to rounding errors in the case of
	// pre-payments and when line item prices include tax.
	RoundingRuleRoundThenSum cbc.Key = "round-then-sum"
)

// RoundingRules defines the list of supported rounding rules.
var RoundingRules = []*cbc.Definition{
	{
		Key: RoundingRuleSumThenRound,
		Name: i18n.String{
			i18n.EN: "Sum Then Round",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The default method of calculating the totals in GOBL, and provides the best results
				for most cases as the precision is maintained to the maximum amount possible. The
				tradeoff however is that sometimes the totals may not sum exactly based on what is visible.
			`),
		},
	},
	{
		Key: RoundingRuleRoundThenSum,
		Name: i18n.String{
			i18n.EN: "Round Then Sum",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The alternative method of calculating the totals that will first round all the amounts
				to the currency's precision before making the sums. Totals using this approach can always
				be recalculated using the amounts presented, but can lead to rounding errors in the case
				of pre-payments and when line item prices include tax.
			`),
		},
	},
}

// ApplyRoundingRule applies the given rounding rule to the amount
// using the currency's base precision as a reference.
func ApplyRoundingRule(rr cbc.Key, cur currency.Code, amount num.Amount) num.Amount {
	exp := cur.Def().Subunits
	switch rr {
	case RoundingRuleRoundThenSum:
		return amount.Rescale(exp)
	default:
		return amount.RescaleUp(exp)
	}
}
