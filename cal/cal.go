// Package cal provides simple date handling.
package cal

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("cal"),
		Date{},
		DateTime{},
		Period{},
		Time{},
		Timestamp{},
	)
	rules.Register(
		"cal",
		rules.GOBL.Add("CAL"),
		dateRules(),
		dateTimeRules(),
		periodRules(),
		timestampRules(),
	)
}
