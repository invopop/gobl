package is

import (
	"strings"

	"github.com/invopop/gobl/rules"
)

type oneOfTest struct {
	desc  string
	tests []rules.Test
}

// OneOf defines a test that will pass only if exactly one of the provided
// tests passes. If zero or more than one pass, the test fails.
func OneOf(tests ...rules.Test) rules.Test {
	descs := make([]string, 0, len(tests))
	for _, t := range tests {
		descs = append(descs, t.String())
	}
	return oneOfTest{
		desc:  "exactly one of: " + strings.Join(descs, ", "),
		tests: tests,
	}
}

// Check will run each of the tests on the object and return true only
// if exactly one of the tests passes.
func (t oneOfTest) Check(obj any) bool {
	passed := 0
	for _, test := range t.tests {
		if test.Check(obj) {
			passed++
			if passed > 1 {
				return false
			}
		}
	}
	return passed == 1
}

// CheckWithContext implements rules.ContextualTest so that OneOf(InContext(...), ...)
// correctly threads the context through to each inner test.
func (t oneOfTest) CheckWithContext(rc *rules.Context, val any) bool {
	passed := 0
	for _, test := range t.tests {
		var ok bool
		if ct, isCtx := test.(rules.ContextualTest); isCtx {
			ok = ct.CheckWithContext(rc, val)
		} else {
			ok = test.Check(val)
		}
		if ok {
			passed++
			if passed > 1 {
				return false
			}
		}
	}
	return passed == 1
}

// String provides the string representation of the test.
func (t oneOfTest) String() string {
	return t.desc
}
