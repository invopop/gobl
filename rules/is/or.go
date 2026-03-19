package is

import (
	"strings"

	"github.com/invopop/gobl/rules"
)

type orTest struct {
	desc  string
	tests []rules.Test
}

// Or defines a test that will pass if any of the provided tests pass.
func Or(tests ...rules.Test) rules.Test {
	var descs []string
	for _, t := range tests {
		descs = append(descs, t.String())
	}
	return orTest{
		desc:  strings.Join(descs, ", or "),
		tests: tests,
	}
}

// Check will run each of the tests on the object and return true
// if any of the tests pass.
func (t orTest) Check(obj any) bool {
	for _, test := range t.tests {
		if test.Check(obj) {
			return true
		}
	}
	return false
}

// CheckWithContext implements rules.ContextualTest so that Or(HasContext(...), ...)
// correctly threads the context through to each inner test.
func (t orTest) CheckWithContext(rc *rules.Context, val any) bool {
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
func (t orTest) String() string {
	return t.desc
}
