// Package pay handles models related to payments.
package pay

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("pay"),
		Advance{},
		Instructions{},
		Terms{},
	)
	rules.Register(
		"pay",
		rules.GOBL.Add("PAY"),
		advanceRules(),
		instructionsRules(),
		onlineRules(),
		termsRules(),
		dueDateRules(),
	)
}
