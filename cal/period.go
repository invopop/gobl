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
		rules.Object(
			rules.Assert("03", "either a start or end date is required",
				is.Func("has start or end", periodHasBound),
			),
			rules.Assert("10", "end date must be on or after start date",
				is.Func("end not before start", periodEndNotBeforeStart),
			),
		),
	)
}

// periodHasBound checks that at least one of the two dates is set. Which one
// is irrelevant: EN 16931 BR-CO-19 and BR-CO-20 allow either bound alone, so
// this is reported once against the period rather than against each field.
func periodHasBound(val any) bool {
	p, ok := val.(*Period)
	if !ok || p == nil {
		return true
	}
	return !p.Start.IsZero() || !p.End.IsZero()
}

func periodEndNotBeforeStart(val any) bool {
	p, ok := val.(*Period)
	if !ok || p == nil || p.Start.IsZero() || p.End.IsZero() {
		return true
	}
	return p.End.DaysSince(p.Start.Date) >= 0
}
