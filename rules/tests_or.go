package rules

import "strings"

type orTest struct {
	desc  string
	tests []Test
}

// Or defines a test that will pass if any of the provided tests pass.
func Or(tests ...Test) Test {
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

// checkWithContext implements contextualTest so that Or(HasContext(...), ...)
// correctly threads the context through to each inner test.
func (t orTest) checkWithContext(rc *RunCtx, val any) bool {
	for _, test := range t.tests {
		if runTest(rc, test, val) {
			return true
		}
	}
	return false
}

// String provides the string representation of the test
func (t orTest) String() string {
	return t.desc
}
