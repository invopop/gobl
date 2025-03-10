package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// RoundingRulePrecise is the default method of performing calculations
	// in GOBL, and provides the best results for most cases. Before calculations
	// are made, item and line prices are calculated with at least the currency's
	// precision plus 2 significant digits. In the case of Euros, all calculations
	// would be made with at least 4 decimal places. The tradeoff however is
	// that sometimes the totals may not sum exactly based on what is visible in the
	// document, which can cause problems with some electronic invoice formats or
	// regional requirements.
	RoundingRulePrecise cbc.Key = "precise"

	// RoundingRuleCurrency is the alternative method of performing calculations
	// whereby the currency's precision or subunits are used to round the amounts
	// **before** performing the sums. This can lead to rounding errors when converting
	// to and from prices that include taxes, but is a common approach in other digital
	// invoicing formats.
	RoundingRuleCurrency cbc.Key = "currency"
)

// RoundingRules defines the list of supported rounding rules.
var RoundingRules = []*cbc.Definition{
	{
		Key: RoundingRulePrecise,
		Name: i18n.String{
			i18n.EN: "Precise",
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
		Key: RoundingRuleCurrency,
		Name: i18n.String{
			i18n.EN: "Currency",
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
	case RoundingRuleCurrency:
		return amount.Rescale(exp)
	default:
		return amount.RescaleUp(exp)
	}
}
