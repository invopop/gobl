package tax

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
)

// Error is a general wrapper around tax errors produced during run
// time, typically during calculations. Not to be confused with errors
// produced from definition validation.
type Error cbc.Key

// Standard list of tax errors
const (
	ErrMissingRegime        Error = "missing-regime"
	ErrInvalid              Error = "invalid"
	ErrInvalidCategory      Error = "invalid-category"
	ErrInvalidTag           Error = "invalid-tag"
	ErrInvalidDate          Error = "invalid-date"
	ErrInvalidPricesInclude Error = "invalid-prices-include"
)

// Error serializes the error's message.
func (e Error) Error() string {
	return string(e)
}

// WithMessage wraps around the original error so we can use if for matching
// and adds a human message.
func (e Error) WithMessage(msg string, s ...interface{}) error {
	msg = fmt.Sprintf(msg, s...)
	return fmt.Errorf("%w: %v", e, msg)
}
