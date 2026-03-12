package rules

type emptyTest struct{}

var (
	// Empty checks that a value is not present (nil or empty).
	// A value is considered empty if
	// - integer, float: zero
	// - bool: false
	// - string, array, slice, map: len() == 0
	// - interface, pointer: nil or the referenced value is empty
	// It is the inverse of Required.
	Empty Test = emptyTest{}
)

func (e emptyTest) Check(value any) bool {
	value, isNil := Indirect(value)
	return isNil || IsEmpty(value)
}

func (e emptyTest) String() string {
	return "empty"
}
