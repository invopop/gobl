package norm

import (
	"reflect"

	"github.com/invopop/gobl/rules"
)

// registered is a flattened normalizer: a function bound to a type with the
// full list of guards (registration-level plus any When guards) that must all
// pass before it runs.
type registered struct {
	guards []rules.Test
	fn     Func
}

// typeIndex maps a target type to the normalizers registered for it, in
// registration order. Because the engine looks normalizers up by the concrete
// type of each node it visits, guards are only ever evaluated for types that
// actually have a normalizer, so no separate keyed/unkeyed guard partitioning
// (as in rules) is needed.
var typeIndex = make(map[reflect.Type][]*registered)

// Register adds one or more normalizers to the global registry. The pkg label
// is informational (it identifies the owning package, mirroring rules.Register)
// and plays no part in matching.
func Register(pkg string, sets ...*Set) {
	RegisterWithGuard(pkg, nil, sets...)
}

// RegisterWithGuard adds one or more normalizers behind a shared guard. The
// guard is typically a context test such as is.InContext(tax.AddonIn(V4)) or
// is.InContext(tax.RegimeIn("ES")) so that the normalizers only run for
// documents using that addon or regime.
func RegisterWithGuard(pkg string, guard rules.Test, sets ...*Set) {
	var base []rules.Test
	if guard != nil {
		base = []rules.Test{guard}
	}
	for _, s := range sets {
		flatten(s, base)
	}
}

// flatten walks a Set tree accumulating guards and emits a registered entry for
// every For leaf (one carrying a target type and function).
func flatten(s *Set, guards []rules.Test) {
	if s == nil {
		return
	}
	next := guards
	if s.guard != nil {
		next = append(append([]rules.Test(nil), guards...), s.guard)
	}
	if s.objType != nil && s.fn != nil {
		typeIndex[s.objType] = append(typeIndex[s.objType], &registered{guards: next, fn: s.fn})
	}
	for _, ss := range s.subsets {
		flatten(ss, next)
	}
}

// forType returns the normalizers registered for the given type ordered so that
// unguarded (intrinsic) normalizers run before guarded ones. This keeps a
// type's own normalization ahead of regime/addon adjustments, matching the
// previous "intrinsic body then normalizers.Each" ordering. Within each group
// registration order is preserved.
func forType(t reflect.Type) []*registered {
	regs := typeIndex[t]
	if len(regs) == 0 {
		return nil
	}
	ordered := make([]*registered, 0, len(regs))
	for _, r := range regs {
		if len(r.guards) == 0 {
			ordered = append(ordered, r)
		}
	}
	for _, r := range regs {
		if len(r.guards) > 0 {
			ordered = append(ordered, r)
		}
	}
	return ordered
}
