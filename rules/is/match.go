package is

import "regexp"

type matchTest struct {
	pattern string
	re      *regexp.Regexp
}

// Matches provides a validation rule that checks if a string value matches the specified regular expression
// pattern. Patterns will be compiled when used in rules.For() or rules.ForValue() and cached for future use.
func Matches(pattern string) *matchTest {
	return &matchTest{
		pattern: pattern,
	}
}

// MatchesRegexp provides a validation rule that checks if the provided Regular Expression matches the value.
func MatchesRegexp(re *regexp.Regexp) *matchTest {
	return &matchTest{
		pattern: re.String(),
		re:      re,
	}
}

func (t matchTest) Check(value any) bool {
	if t.re == nil {
		panic("match test was not compiled; use MatchesRegexp or wrap it in rules.For()")
	}
	value, isNil := Indirect(value)
	if isNil {
		return true // ignore
	}

	isString, str, isBytes, bs := StringOrBytes(value)
	if isString && (str == "" || t.re.MatchString(str)) {
		return true
	} else if isBytes && (len(bs) == 0 || t.re.Match(bs)) {
		return true
	}

	return false
}

func (t *matchTest) Compile(_ any) error {
	if t.re != nil {
		return nil
	}
	var err error
	t.re, err = regexp.Compile(t.pattern)
	if err != nil {
		return err
	}
	return nil
}

func (t matchTest) String() string {
	return "matches " + t.pattern
}
