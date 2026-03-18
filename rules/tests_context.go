package rules

// HasContext returns a Test that passes when any value in the current
// validation context satisfies the provided inner test t.
//
// When used as a guard on a registered rule set, the context is built
// automatically from the root object's embedded tax.Regime and tax.Addons
// (or any type implementing ContextAdder), and from explicit WithContext
// options passed to rules.Validate.
//
// When called via Check without a context (e.g. Set.Validate directly),
// it falls back to testing the validated value itself with t.
func HasContext(t Test) Test {
	return &hasContextTest{test: t}
}

type hasContextTest struct {
	test Test
}

// checkWithContext implements contextualTest: iterates context values and
// returns true when any satisfies the inner test.
func (h *hasContextTest) checkWithContext(rc *RunCtx, _ any) bool {
	for _, v := range rc.values {
		if h.test.Check(v) {
			return true
		}
	}
	return false
}

// Check falls back to testing val directly when no context is available.
func (h *hasContextTest) Check(val any) bool {
	return h.test.Check(val)
}

// String returns a human-readable description of the test.
func (h *hasContextTest) String() string {
	return "context: " + h.test.String()
}
