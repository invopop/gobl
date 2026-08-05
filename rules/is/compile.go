package is

import (
	"github.com/invopop/gobl/rules"
)

// compilable is implemented by tests that need to be prepared before use,
// such as Expr or Matches. It mirrors the unexported interface used by the
// rules package so that tests wrapping other tests can pass compilation on.
type compilable interface {
	Compile(val any) error
}

// compileTest prepares the given test when it requires compilation.
func compileTest(val any, test rules.Test) error {
	if ct, ok := test.(compilable); ok {
		return ct.Compile(val)
	}
	return nil
}

// compileTests prepares each of the given tests that require compilation.
func compileTests(val any, tests []rules.Test) error {
	for _, test := range tests {
		if err := compileTest(val, test); err != nil {
			return err
		}
	}
	return nil
}
