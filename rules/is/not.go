package is

import (
	"github.com/invopop/gobl/rules"
)

type notTest struct {
	test rules.Test
}

// Not defines a test that will pass when the provided test does not.
func Not(test rules.Test) rules.Test {
	return notTest{test: test}
}

// Check will run the wrapped test on the object and invert the result.
func (t notTest) Check(obj any) bool {
	return !t.test.Check(obj)
}

// CheckWithContext implements rules.ContextualTest so that Not(InContext(...))
// correctly threads the context through to the wrapped test.
func (t notTest) CheckWithContext(rc *rules.Context, val any) bool {
	if ct, ok := t.test.(rules.ContextualTest); ok {
		return !ct.CheckWithContext(rc, val)
	}
	return !t.test.Check(val)
}

// Compile prepares the wrapped test when it requires compilation, such as
// Expr or Matches.
func (t notTest) Compile(val any) error {
	return compileTest(val, t.test)
}

// String provides the string representation of the test.
func (t notTest) String() string {
	return "not (" + t.test.String() + ")"
}
