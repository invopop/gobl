package rules

type nilTest struct{}

var (
	// Nil checks that a value is a nil pointer.
	// Unlike Empty, it does not pass for empty non-nil values such as "".
	Nil Test = nilTest{}
)

func (n nilTest) Check(value any) bool {
	_, isNil := Indirect(value)
	return isNil
}

func (n nilTest) String() string {
	return "nil"
}
