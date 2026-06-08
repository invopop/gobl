// Package norm provides a framework for defining and applying normalization
// to data structures. It is the normalization counterpart to the rules
// package: instead of declaring and calling normalizer methods and manually
// recursing into every property, normalizers are registered against a type and
// a reflective engine walks the object graph applying the matching ones.
//
// A normalizer is a plain Go function bound to a single type via For. It may
// be gated by a guard Test (typically built with the is package, e.g.
// is.InContext(tax.AddonIn(V4))) exactly as guards work in the rules package.
//
// norm is intentionally simpler than rules: normalizers are not identified by
// code, cannot fail, and there are no field/each sub-trees — the engine
// discovers nested values by reflection. Normalizers must be idempotent: the
// engine may apply them more than once when meta-addons append further addons
// during normalization (see Normalize).
package norm

import (
	"reflect"

	"github.com/invopop/gobl/rules"
)

// Func is a normalizer bound to a concrete type by For. It receives a pointer
// to the value being normalized so that mutations persist.
type Func func(obj any)

// Set binds a normalizer to a type (via For), or groups guarded subsets (via
// When). It is the building block passed to Register and RegisterWithGuard.
// Unlike rules.Set there are no IDs, assertions, or field/each subsets.
type Set struct {
	objType reflect.Type // type this set normalizes; nil for a When grouping
	guard   rules.Test   // optional guard; from When or RegisterWithGuard
	fn      Func         // the normalizer to run on a matching value; nil for a When grouping
	subsets []*Set       // grouping only (When); each carries its own objType
}

// For binds a typed normalizer to its type. The target type T is inferred from
// the function signature, so each normalizer targets exactly one type and no
// type switch is required.
//
// The function receives a non-nil pointer to the value being normalized:
//
//	norm.For(func(inv *bill.Invoice) {
//	    if inv.Type == cbc.KeyEmpty {
//	        inv.Type = bill.InvoiceTypeStandard
//	    }
//	})
func For[T any](fn func(*T)) *Set {
	return &Set{
		objType: reflect.TypeOf((*T)(nil)).Elem(),
		fn: func(obj any) {
			if v, ok := obj.(*T); ok {
				fn(v)
			}
		},
	}
}

// When wraps sets so that they only run when guard passes. It composes with
// any registration-level guard (RegisterWithGuard): all guards along the path
// from registration to a normalizer must pass for it to run. Guards are
// evaluated against the same value (a pointer to the node) and may inspect the
// validation context when they implement rules.ContextualTest, so context
// guards such as is.InContext(tax.AddonIn(V4)) work unchanged.
func When(guard rules.Test, sets ...*Set) *Set {
	return &Set{guard: guard, subsets: sets}
}
