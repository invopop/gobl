package rules

import (
	"fmt"
	"reflect"
	"time"
)

type thresholdOp int

const (
	thresholdGTE thresholdOp = iota // >=
	thresholdGT                     // >
	thresholdLTE                    // <=
	thresholdLT                     // <
)

type minMaxTest struct {
	threshold any
	op        thresholdOp
}

// Min returns a validation rule that checks if a value is greater than or equal to threshold.
// Supports int, uint, float, and time.Time types. Empty values are skipped.
func Min(threshold any) minMaxTest {
	return minMaxTest{threshold: threshold, op: thresholdGTE}
}

// Max returns a validation rule that checks if a value is less than or equal to threshold.
// Supports int, uint, float, and time.Time types. Empty values are skipped.
func Max(threshold any) minMaxTest {
	return minMaxTest{threshold: threshold, op: thresholdLTE}
}

// Exclusive returns a copy of the rule using strict (exclusive) comparison.
func (t minMaxTest) Exclusive() minMaxTest {
	switch t.op {
	case thresholdGTE:
		t.op = thresholdGT
	case thresholdLTE:
		t.op = thresholdLT
	}
	return t
}

func (t minMaxTest) Check(value any) bool {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return true
	}

	rv := reflect.ValueOf(t.threshold)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := reflectInt(value)
		if err != nil {
			return false
		}
		return t.compareInt(rv.Int(), v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v, err := reflectUint(value)
		if err != nil {
			return false
		}
		return t.compareUint(rv.Uint(), v)
	case reflect.Float32, reflect.Float64:
		v, err := reflectFloat(value)
		if err != nil {
			return false
		}
		return t.compareFloat(rv.Float(), v)
	case reflect.Struct:
		th, ok := t.threshold.(time.Time)
		if !ok {
			return false
		}
		v, ok := value.(time.Time)
		if !ok {
			return false
		}
		if v.IsZero() {
			return true
		}
		return t.compareTime(th, v)
	}
	return false
}

func (t minMaxTest) String() string {
	switch t.op {
	case thresholdGT:
		return fmt.Sprintf("greater than %v", t.threshold)
	case thresholdGTE:
		return fmt.Sprintf("at least %v", t.threshold)
	case thresholdLT:
		return fmt.Sprintf("less than %v", t.threshold)
	default:
		return fmt.Sprintf("at most %v", t.threshold)
	}
}

func (t minMaxTest) compareInt(threshold, value int64) bool {
	switch t.op {
	case thresholdGT:
		return value > threshold
	case thresholdGTE:
		return value >= threshold
	case thresholdLT:
		return value < threshold
	default:
		return value <= threshold
	}
}

func (t minMaxTest) compareUint(threshold, value uint64) bool {
	switch t.op {
	case thresholdGT:
		return value > threshold
	case thresholdGTE:
		return value >= threshold
	case thresholdLT:
		return value < threshold
	default:
		return value <= threshold
	}
}

func (t minMaxTest) compareFloat(threshold, value float64) bool {
	switch t.op {
	case thresholdGT:
		return value > threshold
	case thresholdGTE:
		return value >= threshold
	case thresholdLT:
		return value < threshold
	default:
		return value <= threshold
	}
}

func (t minMaxTest) compareTime(threshold, value time.Time) bool {
	switch t.op {
	case thresholdGT:
		return value.After(threshold)
	case thresholdGTE:
		return value.After(threshold) || value.Equal(threshold)
	case thresholdLT:
		return value.Before(threshold)
	default:
		return value.Before(threshold) || value.Equal(threshold)
	}
}

func reflectInt(value any) (int64, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	}
	return 0, fmt.Errorf("cannot convert %v to int64", v.Kind())
}

func reflectUint(value any) (uint64, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint(), nil
	}
	return 0, fmt.Errorf("cannot convert %v to uint64", v.Kind())
}

func reflectFloat(value any) (float64, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
	}
	return 0, fmt.Errorf("cannot convert %v to float64", v.Kind())
}
