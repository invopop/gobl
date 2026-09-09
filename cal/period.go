package cal

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// Period represents a span of time bounded by a start and/or an end date.
type Period struct {
	// Label is a short description of the period.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Start indicates when this period starts, required if there is no end date.
	Start *Date `json:"start,omitempty" jsonschema:"title=Start"`
	// End indicates when the period ends, required if there is no start date,
	// and must be on or after the start date.
	End *Date `json:"end,omitempty" jsonschema:"title=End"`
}

func periodRules() *rules.Set {
	return rules.For(new(Period),
		rules.Field("start",
			rules.AssertIfPresent("01", "start date cannot be zero",
				DateNotZero(),
			),
		),
		rules.Field("end",
			rules.AssertIfPresent("02", "end date cannot be zero",
				DateNotZero(),
			),
		),
		rules.Object(
			rules.Assert("10", "end date must be on or after start date",
				is.Func("end not before start", periodEndNotBeforeStart),
			),
			rules.Assert("11", "either a start or end date is required",
				is.Func("has start or end", periodHasBound),
			),
		),
	)
}

func periodHasBound(val any) bool {
	p, ok := val.(*Period)
	if !ok || p == nil {
		return true
	}
	return p.Start != nil || p.End != nil
}

func periodEndNotBeforeStart(val any) bool {
	p, ok := val.(*Period)
	if !ok || p == nil || p.Start == nil || p.End == nil {
		return true
	}
	if p.Start.IsZero() || p.End.IsZero() {
		return true
	}
	return p.End.DaysSince(p.Start.Date) >= 0
}
