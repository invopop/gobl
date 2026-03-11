package rules

type requiredTest struct {
	desc    string
	skipNil bool
}

var (
	// Required provides a validation rule that checks if a value is present and not empty.
	// A value is considered not empty if
	// - integer, float: not zero
	// - bool: true
	// - string, array, slice, map: len() > 0
	// - interface, pointer: not nil and the referenced value is not empty
	// - any other types
	Required Test = requiredTest{"required", false}

	// NilOrNotEmpty checks if a value is a nil pointer or a value that is not empty.
	// NilOrNotEmpty differs from Required in that it treats a nil pointer as valid.
	NilOrNotEmpty Test = requiredTest{"nil or not empty", true}
)

func (r requiredTest) Check(value any) bool {
	value, isNil := Indirect(value)
	if r.skipNil {
		return isNil || !IsEmpty(value)
	}
	return !isNil && !IsEmpty(value)
}

func (r requiredTest) String() string {
	return r.desc
}
