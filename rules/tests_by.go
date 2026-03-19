package rules

type byTest struct {
	desc string
	test func(any) bool
}

// By creates a validation rule that checks if a value satisfies a custom condition defined by the provided test function.
func By(desc string, test func(any) bool) byTest {
	return byTest{
		desc: desc,
		test: test,
	}
}

// ByError creates a validation rule that checks if a value satisfies a custom condition defined by the provided test function that returns an error.
func ByError(desc string, test func(any) error) byTest {
	return byTest{
		desc: desc,
		test: func(value any) bool {
			err := test(value)
			return err == nil
		},
	}
}

func (t byTest) Check(value any) bool {
	return t.test(value)
}

func (t byTest) String() string {
	return t.desc
}

// ByContext creates a Test that has access to the validation context values.
// fn receives the context values slice (nil when called without context) and
// the current value being tested. Use this to implement tests that need to
// inspect values injected via rules.WithContext options or embedded ContextAdder
// fields (e.g. tax.Regime).
func ByContext(desc string, fn func(ctx []any, val any) bool) byContextTest {
	return byContextTest{desc: desc, fn: fn}
}

type byContextTest struct {
	desc string
	fn   func(ctx []any, val any) bool
}

// checkWithContext implements contextualTest, giving this test access to
// rc.values during a rules.Validate call.
func (t byContextTest) checkWithContext(rc *RunCtx, val any) bool {
	return t.fn(rc.values, val)
}

// Check satisfies the Test interface when no context is available.
func (t byContextTest) Check(val any) bool {
	return t.fn(nil, val)
}

func (t byContextTest) String() string {
	return t.desc
}
