package cal

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// Period represents a span of time bounded by a start and/or an end date.
// At least one of the two bounds must be provided, but either may be omitted
// on its own, mirroring EN 16931 (BR-CO-19 and BR-CO-20) where the invoicing
// and line period start and end dates are each optional.
type Period struct {
	// Label is a short description of the period.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Start indicates when this period starts. Required if no end date is provided.
	Start Date `json:"start,omitzero" jsonschema:"title=Start"`
	// End indicates when the period ends, and must be on or after the start
	// date when both are present. Required if no start date is provided.
	End Date `json:"end,omitzero" jsonschema:"title=End"`
}

func periodRules() *rules.Set {
	return rules.For(new(Period),
		rules.When(is.Expr(`End.IsZero()`),
			rules.Field("start",
				rules.Assert("01", "start date cannot be zero",
					DateNotZero(),
				),
			),
		),
		rules.When(is.Expr(`Start.IsZero()`),
			rules.Field("end",
				rules.Assert("02", "end date cannot be zero",
					DateNotZero(),
				),
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
