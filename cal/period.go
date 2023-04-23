package cal

import "github.com/invopop/validation"

// Period represents two dates with a start and finish.
type Period struct {
	// Label is a short description of the period.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Start indicates when this period starts.
	Start Date `json:"start" jsonschema:"title=Start"`
	// End indicates when the period ends, and must be after the start date.
	End Date `json:"end" jsonschema:"title=End"`
}

// Validate checks to ensure the period looks correct.
func (p *Period) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Start, DateNotZero(), DateBefore(p.End)),
		validation.Field(&p.End, DateNotZero(), DateAfter(p.Start)),
	)
}
