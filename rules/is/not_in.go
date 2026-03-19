package is

import (
	"fmt"
	"strings"
)

type notInTest struct {
	set []any // normalized by compile
}

// NotIn provides a validation rule that checks if a value is absent from the provided set.
// Named types (e.g. type Code string) are compared against their underlying primitive,
// so NotIn("A", "B") will reject both string("A") and Code("A").
// Nil values are skipped; use Present to enforce presence.
func NotIn(set ...any) *notInTest {
	return &notInTest{set: set}
}

func (t *notInTest) Compile(_ any) error {
	for i, v := range t.set {
		t.set[i] = normalizeValue(v)
	}
	return nil
}

func (t notInTest) Check(value any) bool {
	value, isNil := Indirect(value)
	if isNil {
		return true // skip nil
	}
	norm := normalizeValue(value)
	for _, v := range t.set {
		if norm == normalizeValue(v) {
			return false
		}
	}
	return true
}

func (t notInTest) String() string {
	parts := make([]string, len(t.set))
	for i, v := range t.set {
		parts[i] = fmt.Sprintf("%v", v)
	}
	return "not one of [" + strings.Join(parts, ", ") + "]"
}
