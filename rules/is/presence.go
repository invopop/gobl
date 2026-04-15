package is

type presenceTest struct {
	desc    string
	skipNil bool
}

var (
	// Present checks if a value is non-nil and non-empty.
	// A value is considered not empty if
	// - integer, float: not zero
	// - bool: true
	// - string, array, slice, map: len() > 0
	// - interface, pointer: not nil and the referenced value is not empty
	// - any other types
	Present = presenceTest{"present", false}

	// NilOrNotEmpty checks if a value is a nil pointer or a value that is not empty.
	// NilOrNotEmpty differs from Present in that it treats a nil pointer as valid.
	NilOrNotEmpty = presenceTest{"nil or not empty", true}
)

func (r presenceTest) Check(value any) bool {
	value, isNil := Indirect(value)
	if r.skipNil {
		return isNil || !emptyValue(value)
	}
	return !isNil && !emptyValue(value)
}

func (r presenceTest) String() string {
	return r.desc
}
