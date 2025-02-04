// Package tax encapsulates models related to taxation.
package tax

import (
	"context"
	"reflect"

	"github.com/invopop/gobl/schema"
	"github.com/invopop/validation"
)

func init() {
	schema.Register(schema.GOBL.Add("tax"),
		Identity{},
		Set{},
		Extensions{},
		Total{},
		RegimeDef{},
		AddonDef{},
		CatalogueDef{},
	)
}

type contextKey string

const (
	validtorsKey contextKey = "validators"
)

// Normalizer is used for functions that will normalize the provided object
// ensuring that all the data is aligned with expected values, and adding
// any additional data tha may be required.
//
// Normalizer cannot fail by design, they should always be designed to fail
// silently in case of issues and depend on the Validator to pick up on
// any issues.
type Normalizer func(doc any)

// Normalizers defines a list of normalizer methods with some helpers for
// execution.
type Normalizers []Normalizer

// Validator is used for functions that will validate the provided object
// and provide an error if the object is not valid.
type Validator func(doc any) error

// contextWithValidator will add the provided validator to the current validators.
// This should be prepared from the Addon or Regime itself.
func contextWithValidator(ctx context.Context, v Validator) context.Context {
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

// ExtractNormalizers will extract the normalizers from the provided object
// that is using either the regime or addons.
func ExtractNormalizers(obj any) Normalizers {
	if obj == nil {
		return nil
	}
	normalizers := make(Normalizers, 0)
	if n, ok := obj.(regimeImpl); ok {
		if r := n.RegimeDef(); r != nil {
			normalizers = normalizers.Append(r.Normalizer)
		}
	}
	if n, ok := obj.(addonsImpl); ok {
		n.normalizeAddons()
		for _, a := range n.AddonDefs() {
			normalizers = normalizers.Append(a.Normalizer)
		}
	}
	return normalizers
}

type regimeImpl interface {
	RegimeDef() *RegimeDef
}

type addonsImpl interface {
	normalizeAddons()
	AddonDefs() []*AddonDef
}

type normalizeImpl interface {
	Normalize(Normalizers)
}

// Each will run a simple loop over the normalizers on the provided object.
func (ns Normalizers) Each(doc any) {
	if doc == nil {
		return
	}
	if ns == nil {
		return
	}
	for _, n := range ns {
		n(doc)
	}
}

// Append adds the normalizer, but only if it is not nil.
func (ns Normalizers) Append(n Normalizer) Normalizers {
	if n == nil {
		return ns
	}
	return append(ns, n)
}

// Normalize will either run the "Normalize" method on the provided object,
// or directly go through the list of normalizers on the object.
// This supports arrays and slices, and will automatically normalize each
// element in the list.
func Normalize(list Normalizers, doc any) {
	if doc == nil {
		return
	}
	if n, ok := doc.(normalizeImpl); ok {
		n.Normalize(list)
	} else {
		switch reflect.TypeOf(doc).Kind() {
		case reflect.Slice, reflect.Array:
			s := reflect.ValueOf(doc)
			for i := 0; i < s.Len(); i++ {
				d := s.Index(i).Interface()
				Normalize(list, d)
			}
		default:
			list.Each(doc)
		}
	}
}

// ValidateStructWithContext wraps around the standard validation.ValidateStructWithContext
// method to add an additional check for the tax regime.
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
