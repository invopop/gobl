package facturae

import (
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// payDueDateRules ensures every due date carries the non-zero amount.
func payDueDateRules() *rules.Set {
	return rules.For(new(pay.DueDate),
		rules.Field("amount",
			rules.Assert("01", "amount is required", is.Present, num.NotZero),
		),
	)
}
