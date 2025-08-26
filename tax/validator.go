package tax

import (
	"context"

	"github.com/invopop/validation"
)

type contextKey string

const (
	validtorsKey contextKey = "validators"
)

// Validator is used for functions that will validate the provided object
// and provide an error if the object is not valid.
type Validator func(doc any) error

// ContextWithValidator will add the provided validator to the current validators.
// This should be prepared from the Addon or Regime itself.
func ContextWithValidator(ctx context.Context, v Validator) context.Context {
	if v == nil {
		return ctx
	}
	list := append(Validators(ctx), v)
	return context.WithValue(ctx, validtorsKey, list)
}

// Validators provides the list of validators that have been added to the current
// context.
func Validators(ctx context.Context) []Validator {
	list, ok := ctx.Value(validtorsKey).([]Validator)
	if !ok {
		return make([]Validator, 0)
	}
	return list
}

// ValidateStructWithContext wraps around the standard validation.ValidateStructWithContext
// method to add an additional tax specific validation checks from the context. See
// also the ContextWithValidator method.
func ValidateStructWithContext(ctx context.Context, obj any, fields ...*validation.FieldRules) error {
	// First run regular validation
	if err := validation.ValidateStructWithContext(ctx, obj, fields...); err != nil {
		return err
	}
	for _, validator := range Validators(ctx) {
		if err := validator(obj); err != nil {
			return err
		}
	}
	return nil
}
