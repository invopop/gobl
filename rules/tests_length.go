package rules

import "fmt"

type lengthTest struct {
	min int
	max int
}

func Length(min, max int) lengthTest {
	return lengthTest{
		min: min,
		max: max,
	}
}

func (t lengthTest) Check(value any) bool {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return true // ignore
	}

	l, err := LengthOfValue(value)
	if err != nil {
		return false
	}
	return !(t.min > 0 && l < t.min || t.max > 0 && l > t.max || t.min == 0 && t.max == 0 && l > 0)
}

func (t lengthTest) String() string {
	return fmt.Sprintf("length between %d and %d", t.min, t.max)
}
