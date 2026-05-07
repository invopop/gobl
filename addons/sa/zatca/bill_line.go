package zatca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func billLineRules() *rules.Set {
	return rules.For(new(bill.Line),
		rules.Field("charges",
			rules.Assert("01", "line charges are not allowed (BR-KSA-EN16931-06)", is.Length(0, 0)),
		),
	)
}
