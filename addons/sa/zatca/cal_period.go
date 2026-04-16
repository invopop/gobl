package zatca

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func calPeriodRules() *rules.Set {
	return rules.For(new(cal.Period),
		rules.Assert("01", "if the invoice has a supply end date, it must also have a start date (BR-KSA-35)",
			is.Func("start and end date must be valid", validStartAndEndDate),
		),
	)
}

func validStartAndEndDate(val any) bool {
	period, ok := val.(*cal.Period)
	if !ok || period == nil {
		return true
	}
	if period.End.IsZero() {
		return true
	}
	return !period.Start.IsZero()
}
