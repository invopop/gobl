// Package pay handles models related to payments.
package pay

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("pay"),
		Card{},
		CreditTransfer{},
		DirectDebit{},
		Instructions{},
		Online{},
		Record{},
		Terms{},
	)
	rules.Register(
		"pay",
		rules.GOBL.Add("PAY"),
		dueDateRules(),
		instructionsRules(),
		onlineRules(),
		recordRules(),
		termsRules(),
	)
}
