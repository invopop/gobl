package is

import (
	"strings"

	"github.com/invopop/gobl/rules"
)

type anyOfTest struct {
	desc  string
	tests []rules.Test
}

// Or defines a test that will pass if one or more of the provided
// sub-tests pass.
// Deprecated: use AnyOf instead.
func Or(tests ...rules.Test) rules.Test {
	return AnyOf(tests...)
}

// AnyOf defines a test that will pass if one or more of the provided tests pass.
func AnyOf(tests ...rules.Test) rules.Test {
	var descs []string
	for _, t := range tests {
		descs = append(descs, t.String())
	}
	return anyOfTest{
		desc:  strings.Join(descs, ", or "),
		tests: tests,
	}
}

// Check will run each of the tests on the object and return true
// if any of the tests pass.
func (t anyOfTest) Check(obj any) bool {
	for _, test := range t.tests {
		if test.Check(obj) {
			return true
		}
	}
	return false
}

// CheckWithContext implements rules.ContextualTest so that AnyOf(InContext(...), ...)
// correctly threads the context through to each inner test.
func (t anyOfTest) CheckWithContext(rc *rules.Context, val any) bool {
	for _, test := range t.tests {
		if ct, ok := test.(rules.ContextualTest); ok {
			if ct.CheckWithContext(rc, val) {
				return true
			}
		} else if test.Check(val) {
			return true
		}
	}
	return false
}

// String provides the string representation of the test.
func (t anyOfTest) String() string {
	return t.desc
}
