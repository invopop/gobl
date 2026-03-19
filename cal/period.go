package cal

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// Period represents two dates with a start and finish.
type Period struct {
	// Label is a short description of the period.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Start indicates when this period starts.
	Start Date `json:"start" jsonschema:"title=Start"`
	// End indicates when the period ends, and must be after the start date.
	End Date `json:"end" jsonschema:"title=End"`
}

func periodRules() *rules.Set {
	return rules.For(new(Period),
		rules.Field("start",
			rules.Assert("01", "start date cannot be zero",
				DateNotZero(),
			),
		),
		rules.Field("end",
			rules.Assert("02", "end date cannot be zero",
				DateNotZero(),
			),
		),
		rules.Object(
			rules.Assert("10", "end date must be on or after start date",
				is.Func("end not before start", periodEndNotBeforeStart),
			),
		),
	)
}

func periodEndNotBeforeStart(val any) bool {
	p, ok := val.(*Period)
	if !ok || p == nil || p.Start.IsZero() || p.End.IsZero() {
		return true
	}
	return p.End.DaysSince(p.Start.Date) >= 0
}
