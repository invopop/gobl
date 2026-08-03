package is

import (
	"strings"

	"github.com/invopop/gobl/rules"
)

type allOfTest struct {
	desc  string
	tests []rules.Test
}

// AllOf defines a test that will pass when every one of the provided tests
// passes. Assertions already require all their tests to pass, so this is for
// composing a single named test out of several conditions.
func AllOf(tests ...rules.Test) rules.Test {
	var descs []string
	for _, t := range tests {
		descs = append(descs, t.String())
	}
	return allOfTest{
		desc:  strings.Join(descs, ", and "),
		tests: tests,
	}
}

// Check will run each of the tests on the object and return true only if
// all of them pass.
func (t allOfTest) Check(obj any) bool {
	for _, test := range t.tests {
		if !test.Check(obj) {
			return false
		}
	}
	return true
}

// CheckWithContext implements rules.ContextualTest so that AllOf(InContext(...), ...)
// correctly threads the context through to each inner test.
func (t allOfTest) CheckWithContext(rc *rules.Context, val any) bool {
	for _, test := range t.tests {
		if ct, ok := test.(rules.ContextualTest); ok {
			if !ct.CheckWithContext(rc, val) {
				return false
			}
		} else if !test.Check(val) {
			return false
		}
	}
	return true
}

// String provides the string representation of the test.
func (t allOfTest) String() string {
	return t.desc
}
