package rules

import (
	"fmt"
	"unicode/utf8"
)

type lengthTest struct {
	min  int
	max  int
	rune bool
}

// Length returns a validation rule that checks if a value's length is within the specified range.
// For strings, it counts bytes. If max is 0, there is no upper bound.
// Empty values are skipped; use Required to enforce presence.
func Length(min, max int) lengthTest {
	return lengthTest{
		min: min,
		max: max,
	}
}

// RuneLength returns a validation rule that checks if a string's rune (character) count is
// within the specified range. Unlike Length, it correctly handles multi-byte Unicode characters.
// For non-string types, it behaves identically to Length. If max is 0, there is no upper bound.
// Empty values are skipped; use Required to enforce presence.
func RuneLength(min, max int) lengthTest {
	return lengthTest{
		min:  min,
		max:  max,
		rune: true,
	}
}

func (t lengthTest) Check(value any) bool {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return true // ignore
	}

	var l int
	if t.rune {
		isString, str, _, _ := StringOrBytes(value)
		if isString {
			l = utf8.RuneCountInString(str)
		} else {
			var err error
			l, err = LengthOfValue(value)
			if err != nil {
				return false
			}
		}
	} else {
		var err error
		l, err = LengthOfValue(value)
		if err != nil {
			return false
		}
	}

	return !(t.min > 0 && l < t.min || t.max > 0 && l > t.max || t.min == 0 && t.max == 0 && l > 0)
}

func (t lengthTest) String() string {
	if t.rune {
		return fmt.Sprintf("rune length between %d and %d", t.min, t.max)
	}
	return fmt.Sprintf("length between %d and %d", t.min, t.max)
}
