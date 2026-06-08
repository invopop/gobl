package norm

import (
	"reflect"

	"github.com/invopop/gobl/rules"
)

// maxPasses caps the number of normalization passes. A second pass is needed
// only when a meta-addon normalizer appends further addons during the first
// pass; one extra pass is enough in practice, the rest is a safety net.
const maxPasses = 4

// Normalize applies every registered normalizer that matches a value reachable
// from doc. doc must be a pointer so that mutations persist. The object graph
// is walked depth-first and post-order: a node's children are normalized before
// the node's own normalizers run, so addon/regime normalizers always see
// fully-normalized children.
//
// Guards are evaluated against a validation context built from opts and from
// the root's embedded rules.ContextAdder fields (tax.Regime, tax.Addons), so
// context guards such as is.InContext(tax.AddonIn(V4)) apply to nested values
// without needing access to the root.
//
// Meta-addons may append further addons while normalizing (tax.Addons.AddAddons).
// Normalize detects this by re-collecting the context after each pass and walks
// again until the set of context keys stabilises (bounded by maxPasses). Because
// normalizers are idempotent, re-applying already-run normalizers on a later
// pass is harmless.
func Normalize(doc any, opts ...rules.WithContext) {
	if doc == nil {
		return
	}
	rv := reflect.ValueOf(doc)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return
	}

	var prevKeys map[rules.ContextKey]bool
	for pass := 0; pass < maxPasses; pass++ {
		rc := &rules.Context{}
		for _, opt := range opts {
			opt(rc)
		}
		collectContext(rc, doc)
		keys := keySet(rc)
		if pass > 0 && sameKeys(keys, prevKeys) {
			break // no addon was appended by the previous pass: fixpoint reached
		}
		walk(rc, rv)
		prevKeys = keys
	}
}

// walk visits v, normalizing children first (post-order) and then v itself.
func walk(rc *rules.Context, v reflect.Value) {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return
		}
		walk(rc, v.Elem())
		return
	case reflect.Interface:
		if v.IsNil() {
			return
		}
		walk(rc, v.Elem())
		return
	case reflect.Struct:
		t := v.Type()
		for i := range v.NumField() {
			if !t.Field(i).IsExported() {
				continue
			}
			walk(rc, v.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := range v.Len() {
			walk(rc, v.Index(i))
		}
	case reflect.Map:
		// Map values are not addressable, so only pointer (or interface)
		// values can be normalized in place by recursing into what they
		// point to. Scalar and struct map values are left untouched.
		iter := v.MapRange()
		for iter.Next() {
			mv := iter.Value()
			if mv.Kind() == reflect.Pointer || mv.Kind() == reflect.Interface {
				walk(rc, mv)
			}
		}
	}
	apply(rc, v)
}

// apply runs the normalizers registered for v's type whose guards all pass.
// It requires an addressable value so a pointer can be handed to the normalizer.
func apply(rc *rules.Context, v reflect.Value) {
	if !v.CanAddr() {
		return
	}
	regs := forType(v.Type())
	if len(regs) == 0 {
		return
	}
	ptr := v.Addr().Interface()
	for _, reg := range regs {
		if guardsPass(rc, reg.guards, ptr) {
			reg.fn(ptr)
		}
	}
}

// guardsPass reports whether every guard passes for the given value.
func guardsPass(rc *rules.Context, guards []rules.Test, val any) bool {
	for _, g := range guards {
		if !runTest(rc, g, val) {
			return false
		}
	}
	return true
}

// runTest evaluates a guard, using the validation context when the test
// supports it (e.g. is.InContext) and falling back to a plain value check.
func runTest(rc *rules.Context, t rules.Test, val any) bool {
	if ct, ok := t.(rules.ContextualTest); ok {
		return ct.CheckWithContext(rc, val)
	}
	return t.Check(val)
}

// collectContext builds the normalization context from the root object's
// exported fields that implement rules.ContextAdder. Since tax.Regime and
// tax.Addons are always embedded at the top of document structs, a single-level
// field scan over the root is sufficient (mirroring rules.collectContext).
func collectContext(rc *rules.Context, obj any) {
	if ca, ok := obj.(rules.ContextAdder); ok {
		ca.RulesContext()(rc)
	}
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return
	}
	rt := rv.Type()
	for i := range rt.NumField() {
		if !rt.Field(i).IsExported() {
			continue
		}
		fv := rv.Field(i)
		var fieldObj any
		if fv.CanAddr() {
			fieldObj = fv.Addr().Interface()
		} else {
			fieldObj = fv.Interface()
		}
		if ca, ok := fieldObj.(rules.ContextAdder); ok {
			ca.RulesContext()(rc)
		}
	}
}

// keySet returns the distinct context keys as a set.
func keySet(rc *rules.Context) map[rules.ContextKey]bool {
	keys := rc.Keys()
	set := make(map[rules.ContextKey]bool, len(keys))
	for _, k := range keys {
		set[k] = true
	}
	return set
}

// sameKeys reports whether two context key sets are identical.
func sameKeys(a, b map[rules.ContextKey]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}
