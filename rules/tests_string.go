package rules

type StringTest struct {
	desc string
	test func(string) bool
}

func String(test func(string) bool) StringTest {
	return StringTest{
		desc: "custom string test",
		test: test,
	}
}

func ByString(desc string, test func(string) bool) StringTest {
	return StringTest{
		desc: desc,
		test: test,
	}
}

func (t StringTest) Check(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	return t.test(str)
}

func (t StringTest) String() string {
	return t.desc
}
