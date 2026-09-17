package dian

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
)

func billLineRules() *rules.Set {
	return rules.For(new(bill.Line),
		// Code 01: DIAN rule FAW01 rejects zero-amount lines, so quantities must be positive
		rules.Field("quantity",
			rules.Assert("01", "must be greater than 0", num.Positive),
		),
	)
}
