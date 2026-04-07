package rules

import "reflect"

// ContextKey is the key type for Context entries.
type ContextKey string

// contextEntry holds a single key-value pair in a validation Context.
type contextEntry struct {
	key   ContextKey
	value any
}

// Context holds key-value pairs accumulated during a rules.Validate call and
// is passed to ByContext test functions. Use Set to store values and Value to
// retrieve them.
type Context struct {
	entries []contextEntry
}

// Set appends a key-value pair to the validation context, preserving insertion order.
func (c *Context) Set(key ContextKey, value any) {
	c.entries = append(c.entries, contextEntry{key, value})
}

// Value returns the stored value for key, or nil if absent.
// Callers do a type assertion: v, ok := ctx.Value(key).(MyType)
func (c Context) Value(key ContextKey) any {
	for _, e := range c.entries {
		if e.key == key {
			return e.value
		}
	}
	return nil
}

// WithContext is a functional option for rules.Validate that injects values
// into the validation context before validation begins.
type WithContext func(*Context)

// ContextAdder is implemented by objects that want to automatically inject
// values into the validation context when encountered by the rules engine.
type ContextAdder interface {
	RulesContext() WithContext
}

// ContextualTest is implemented by tests that need access to the validation
// context. The engine checks for this interface before falling back to the
// standard Test.Check method.
type ContextualTest interface {
	CheckWithContext(rc *Context, val any) bool
}

// ContextKeyable is optionally implemented by guard tests that can report
// which context keys they depend on. This enables the engine to skip
// guard evaluation entirely when none of the required keys are present
// in the validation context.
type ContextKeyable interface {
	ContextKeys() []ContextKey
}

// runTest evaluates test t against val. When rc is non-nil and t implements
// ContextualTest, it delegates to CheckWithContext; otherwise it calls Check.
func runTest(rc *Context, t Test, val any) bool {
	if rc != nil {
		if ct, ok := t.(ContextualTest); ok {
			return ct.CheckWithContext(rc, val)
		}
	}
	return t.Check(val)
}

// Keys returns the set of distinct keys present in the context.
func (c Context) Keys() []ContextKey {
	seen := make(map[ContextKey]struct{}, len(c.entries))
	keys := make([]ContextKey, 0, len(c.entries))
	for _, e := range c.entries {
		if _, ok := seen[e.key]; !ok {
			seen[e.key] = struct{}{}
			keys = append(keys, e.key)
		}
	}
	return keys
}

// Each iterates over all values in the context, calling fn for each. Returns
// true as soon as fn returns true (short-circuit), false otherwise.
func (c Context) Each(fn func(value any) bool) bool {
	for _, e := range c.entries {
		if fn(e.value) {
			return true
		}
	}
	return false
}

// collectContext builds the validation context from explicit options and by
// scanning the root object's exported fields for ContextAdder implementations.
// Since tax.Regime and tax.Addons are always embedded at the top of document
// structs, a single-level field scan is sufficient.
func collectContext(rc *Context, obj any) {
	// Scan exported struct fields for embedded ContextAdders.
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return
		}
	}

	// Check the root object itself first.
	if ca, ok := obj.(ContextAdder); ok {
		ca.RulesContext()(rc)
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return
	}
	rt := rv.Type()
	for i := range rt.NumField() {
		sf := rt.Field(i)
		if !sf.IsExported() {
			continue
		}
		fv := rv.Field(i)
		var fieldObj any
		if fv.CanAddr() {
			fieldObj = fv.Addr().Interface()
		} else {
			fieldObj = fv.Interface()
		}
		if ca, ok := fieldObj.(ContextAdder); ok {
			ca.RulesContext()(rc)
		}
		// If the field wraps an embedded payload (e.g. schema.Object),
		// collect context from the inner value as well. Use fv.Interface()
		// directly since fieldObj may be double-pointer for pointer fields.
		if fv.Kind() != reflect.Ptr || !fv.IsNil() {
			if emb, ok := fv.Interface().(Embeddable); ok {
				if inner := emb.Embedded(); inner != nil {
					collectContext(rc, inner)
				}
			}
		}
	}

	// Claude's fix:
	// If the object exposes an embedded payload (e.g. schema.Object wrapping
	// a bill.Invoice), collect context from the inner object so that guards
	// depending on embedded ContextAdders (like tax.Addons) can match.
	if emb, ok := obj.(Embeddable); ok {
		if inner := emb.Embedded(); inner != nil {
			collectContext(rc, inner)
		}
	}
}
