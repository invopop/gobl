package org

import validation "github.com/go-ozzo/ozzo-validation/v4"

// Period represents two dates with a start and finish.
type Period struct {
	Start Date `json:"start"`
	End   Date `json:"end"`
}

// Validate checks to ensure the period looks correct.
func (p *Period) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Start, DateNotZero(), DateBefore(p.End)),
		validation.Field(&p.End, DateNotZero(), DateAfter(p.Start)),
	)
}
