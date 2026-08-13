// Package currency provides models for dealing with currencies.
package currency

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/schema"
)

func init() {
	definitions = new(defs)
	if err := definitions.load(data.Content, "currency"); err != nil {
		panic(err)
	}
	schema.Register(schema.GOBL.Add("currency"),
		Code(""),
		ExchangeRate{},
		Amount{},
	)
	rules.Register(
		"currency",
		rules.GOBL.Add("CURRENCY"),
		codeRules(),
		amountRules(),
		exchangeRateRules(),
	)
}

func codeRules() *rules.Set {
	return rules.For(Code(""),
		rules.AssertIfPresent("01", "currency code must be defined in GOBL", IsCodeDefined),
	)
}

func amountRules() *rules.Set {
	return rules.For(new(Amount),
		rules.Field("currency",
			rules.Assert("01", "currency is required", is.Present),
		),
	)
}

func exchangeRateRules() *rules.Set {
	return rules.For(new(ExchangeRate),
		rules.Field("from",
			rules.Assert("01", "from currency is required", is.Present),
		),
		rules.Field("to",
			rules.Assert("02", "to currency is required", is.Present),
		),
		rules.Field("at",
			rules.AssertIfPresent("03", "date/time must not be zero", cal.DateTimeNotZero()),
		),
		rules.Field("amount",
			rules.Assert("04", "amount must be positive", num.Positive),
		),
	)
}
