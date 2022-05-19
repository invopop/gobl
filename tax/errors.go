package tax

import "fmt"

// Error is a general wrapper around tax errors produced during run
// time, typically during calculations. Not to be confused with errors
// produced from definition validation.
type Error struct {
	Reason Key
}

// Standard list of tax errors
var (
	ErrInvalidCategory      = &Error{Reason: "invalid-category"}
	ErrInvalidRate          = &Error{Reason: "invalid-rate"}
	ErrInvalidDate          = &Error{Reason: "invalid-date"}
	ErrInvalidPricesInclude = &Error{Reason: "invalid-prices-include"}
)

// Error serializes the error's message.
func (e Error) Error() string {
	return e.Reason.String()
}

// WithMessage wraps around the original error so we can use if for matching
// and adds a human message.
func (e *Error) WithMessage(msg string, s ...interface{}) error {
	msg = fmt.Sprintf(msg, s...)
	return fmt.Errorf("%w: %v", e, msg)
}