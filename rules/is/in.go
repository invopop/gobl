package is

import (
	"fmt"
	"strings"
)

// InTest is a Test that checks if a value is one of a set.
type InTest struct {
	set []any // normalized by compile
}

// In provides a validation rule that checks if a value is one of the provided set.
// Named types (e.g. type Code string) are compared against their underlying primitive,
// so In("A", "B") will match both string("A") and Code("A").
// Nil values are skipped; use Present to enforce presence.
func In(set ...any) *InTest {
	return &InTest{set: set}
}

// Compile normalizes the set values for comparison.
func (t *InTest) Compile(_ any) error {
	for i, v := range t.set {
		t.set[i] = normalizeValue(v)
	}
	return nil
}

// Check returns true if the value is present in the set.
func (t InTest) Check(value any) bool {
	value, isNil := Indirect(value)
	if isNil {
		return true // skip nil
	}
	norm := normalizeValue(value)
	for _, v := range t.set {
		if norm == normalizeValue(v) {
			return true
		}
	}
	return false
}

func (t InTest) String() string {
	parts := make([]string, len(t.set))
	for i, v := range t.set {
		parts[i] = fmt.Sprintf("%v", v)
	}
	return "one of [" + strings.Join(parts, ", ") + "]"
}
