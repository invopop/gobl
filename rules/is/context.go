package is

import "github.com/invopop/gobl/rules"

// InContext returns a Test that passes when any value in the current
// validation context satisfies the provided inner test t.
//
// When used as a guard on a registered rule set, the context is built
// automatically from the root object's embedded tax.Regime and tax.Addons
// (or any type implementing rules.ContextAdder), and from explicit
// rules.WithContext options passed to rules.Validate.
//
// When called via Check without a context (e.g. Set.Validate directly),
// it falls back to testing the validated value itself with t.
func InContext(t rules.Test) rules.Test {
	return &inContextTest{test: t}
}

type inContextTest struct {
	test rules.Test
}

// CheckWithContext implements rules.ContextualTest: iterates context values and
// returns true when any satisfies the inner test.
func (h *inContextTest) CheckWithContext(rc *rules.Context, _ any) bool {
	return rc.Each(func(v any) bool { return h.test.Check(v) })
}

// Check falls back to testing val directly when no context is available.
func (h *inContextTest) Check(val any) bool {
	return h.test.Check(val)
}

// String returns a human-readable description of the test.
func (h *inContextTest) String() string {
	return "context: " + h.test.String()
}

// ContextKeys implements rules.ContextKeyable by delegating to the inner test.
// If the inner test does not implement ContextKeyable, nil is returned.
func (h *inContextTest) ContextKeys() []rules.ContextKey {
	if ck, ok := h.test.(rules.ContextKeyable); ok {
		return ck.ContextKeys()
	}
	return nil
}
