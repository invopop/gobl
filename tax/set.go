package tax

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// Set defines a list of tax categories and their rates to be used alongside taxable items.
type Set []*Combo

// CleanSet removes any nil values from the set.
func CleanSet(s Set) Set {
	if s == nil {
		return nil
	}
	ns := make(Set, 0)
	for _, c := range s {
		if c == nil {
			continue
		}
		ns = append(ns, c)
	}
	if len(ns) == 0 {
		return nil
	}
	return ns
}

// Equals returns true if the sets match, regardless of order.
func (s Set) Equals(s2 Set) bool {
	for _, a := range s {
		match := false
		for _, b := range s2 {
			if a.Category == b.Category && a.Key == b.Key && a.Country == b.Country {
				match = true
			}
		}
		if !match {
			// implies the code defined in the base, was not present in the second
			// array.
			return false
		}
	}
	return true
}

// Get the Rate key for the given category
func (s Set) Get(cat cbc.Code) *Combo {
	for _, c := range s {
		if c.Category == cat {
			return c
		}
	}
	return nil
}

// Key returns the key from the matching category, if set.
func (s Set) Key(cat cbc.Code) cbc.Key {
	for _, c := range s {
		if c.Category == cat {
			return c.Key
		}
	}
	return ""
}

func setRules() *rules.Set {
	return rules.For(new(Set),
		rules.Assert("01", "all tax categories in a set must be unique",
			is.Func("no duplicate categories", setNoDuplicateCategories),
		),
	)
}

func setNoDuplicateCategories(val any) bool {
	s, ok := val.(Set)
	if !ok {
		return true
	}
	seen := make(map[cbc.Code]bool)
	for _, c := range s {
		if c == nil {
			continue
		}
		if seen[c.Category] {
			return false
		}
		seen[c.Category] = true
	}
	return true
}

// SetTest defines a validation rule for tax sets, checking for the presence of certain categories.
type SetTest struct {
	desc       string
	categories []cbc.Code
	oneOf      bool
}

// SetHasCategory validates that the set contains the given category.
func SetHasCategory(categories ...cbc.Code) *SetTest {
	return &SetTest{
		desc:       fmt.Sprintf("all of [%s]", strings.Join(cbc.CodeStrings(categories), ", ")),
		categories: categories,
	}
}

// SetHasOneOf checks that the tax set has at least one of the provided
// categories.
func SetHasOneOf(categories ...cbc.Code) *SetTest {
	return &SetTest{
		desc:       fmt.Sprintf("one of [%s]", strings.Join(cbc.CodeStrings(categories), ", ")),
		categories: categories,
		oneOf:      true,
	}
}

// Check returns true if the value passes the validation.
func (sv *SetTest) Check(value any) bool {
	return sv.Validate(value) == nil
}

// String returns a description of the rule.
func (sv *SetTest) String() string {
	return sv.desc
}

// Validate checks that the tax set contains the required categories, and if oneOf is true,
// that at least one of them is present.
func (sv *SetTest) Validate(value any) error {
	s, ok := value.(Set)
	if !ok {
		return nil
	}
	found := false
	for _, c := range sv.categories {
		if s.Get(c) == nil {
			if !sv.oneOf {
				return fmt.Errorf("missing category %s", c.String())
			}
		} else {
			found = true
		}
	}
	if sv.oneOf && !found {
		return fmt.Errorf("missing category in %s", strings.Join(cbc.CodeStrings(sv.categories), ", "))
	}
	return nil
}
