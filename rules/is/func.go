package is

import "github.com/invopop/gobl/rules"

// FuncTest is a Test that validates a value against a custom boolean function.
type FuncTest struct {
	desc string
	test func(any) bool
}

// Func creates a validation rule that checks if a value satisfies a custom
// condition defined by the provided boolean test function.
func Func(desc string, test func(any) bool) FuncTest {
	return FuncTest{
		desc: desc,
		test: test,
	}
}

// FuncError creates a validation rule that checks if a value satisfies a custom
// condition defined by the provided test function that returns an error.
func FuncError(desc string, test func(any) error) FuncTest {
	return FuncTest{
		desc: desc,
		test: func(value any) bool {
			err := test(value)
			return err == nil
		},
	}
}

// Check returns true if the value satisfies the test function.
func (t FuncTest) Check(value any) bool {
	return t.test(value)
}

func (t FuncTest) String() string {
	return t.desc
}

// FuncContext creates a Test that has access to the validation context.
// fn receives a rules.Context (use ctx.Value(key) to extract stored values) and the
// current value being tested. Use this to implement tests that need to
// inspect values injected via rules.WithContext options or embedded ContextAdder
// fields (e.g. tax.Regime).
func FuncContext(desc string, fn func(ctx rules.Context, val any) bool) FuncContextTest {
	return FuncContextTest{desc: desc, fn: fn}
}

// FuncContextTest is a Test with access to the validation context.
type FuncContextTest struct {
	desc string
	fn   func(ctx rules.Context, val any) bool
}

// CheckWithContext implements rules.ContextualTest, giving this test access to
// the validation context during a rules.Validate call.
func (t FuncContextTest) CheckWithContext(rc *rules.Context, val any) bool {
	return t.fn(*rc, val)
}

// Check satisfies the rules.Test interface when no context is available.
func (t FuncContextTest) Check(val any) bool {
	return t.fn(rules.Context{}, val)
}

func (t FuncContextTest) String() string {
	return t.desc
}
