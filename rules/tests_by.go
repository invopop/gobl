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
