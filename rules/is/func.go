package is

import "github.com/invopop/gobl/rules"

type funcTest struct {
	desc string
	test func(any) bool
}

// Func creates a validation rule that checks if a value satisfies a custom
// condition defined by the provided boolean test function.
func Func(desc string, test func(any) bool) funcTest {
	return funcTest{
		desc: desc,
		test: test,
	}
}

// FuncError creates a validation rule that checks if a value satisfies a custom
// condition defined by the provided test function that returns an error.
func FuncError(desc string, test func(any) error) funcTest {
	return funcTest{
		desc: desc,
		test: func(value any) bool {
			err := test(value)
			return err == nil
		},
	}
}

func (t funcTest) Check(value any) bool {
	return t.test(value)
}

func (t funcTest) String() string {
	return t.desc
}

// FuncContext creates a Test that has access to the validation context.
// fn receives a rules.Context (use ctx.Value(key) to extract stored values) and the
// current value being tested. Use this to implement tests that need to
// inspect values injected via rules.WithContext options or embedded ContextAdder
// fields (e.g. tax.Regime).
func FuncContext(desc string, fn func(ctx rules.Context, val any) bool) funcContextTest {
	return funcContextTest{desc: desc, fn: fn}
}

type funcContextTest struct {
	desc string
	fn   func(ctx rules.Context, val any) bool
}

// CheckWithContext implements rules.ContextualTest, giving this test access to
// the validation context during a rules.Validate call.
func (t funcContextTest) CheckWithContext(rc *rules.Context, val any) bool {
	return t.fn(*rc, val)
}

// Check satisfies the rules.Test interface when no context is available.
func (t funcContextTest) Check(val any) bool {
	return t.fn(rules.Context{}, val)
}

func (t funcContextTest) String() string {
	return t.desc
}
