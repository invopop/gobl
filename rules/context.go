package rules

import "reflect"

// RunCtx holds context values accumulated during a rules.Validate call.
// It is threaded through the validation engine to allow context-aware tests
// to check values injected from the root object or explicit options.
type RunCtx struct {
	values []any
}

// Add appends a value to the validation context.
func (rc *RunCtx) Add(v any) {
	rc.values = append(rc.values, v)
}

// WithContext is a functional option for rules.Validate that injects values
// into the validation context before validation begins.
type WithContext func(*RunCtx)

// ContextAdder is implemented by objects that want to automatically inject
// values into the validation context when encountered by the rules engine.
type ContextAdder interface {
	RulesContext() WithContext
}

// contextualTest is an internal interface for tests that need access to the
// validation context. The engine checks for this interface before falling
// back to the standard Test.Check method.
type contextualTest interface {
	checkWithContext(rc *RunCtx, val any) bool
}

// runTest evaluates test t against val. When rc is non-nil and t implements
// contextualTest, it delegates to checkWithContext; otherwise it calls Check.
func runTest(rc *RunCtx, t Test, val any) bool {
	if rc != nil {
		if ct, ok := t.(contextualTest); ok {
			return ct.checkWithContext(rc, val)
		}
	}
	return t.Check(val)
}

// collectContext builds the validation context from explicit options and by
// scanning the root object's exported fields for ContextAdder implementations.
// Since tax.Regime and tax.Addons are always embedded at the top of document
// structs, a single-level field scan is sufficient.
func collectContext(rc *RunCtx, obj any) {
	// Check the root object itself first.
	if ca, ok := obj.(ContextAdder); ok {
		ca.RulesContext()(rc)
	}

	// Scan exported struct fields for embedded ContextAdders.
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
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
	}
}
