package rules

import "regexp"

type matchTest struct {
	pattern string
	re      *regexp.Regexp
}

// Matches provides a validation rule that checks if a string value matches the specified regular expression pattern.
func Matches(pattern string) Test {
	return &matchTest{
		pattern: pattern,
	}
}

func (t matchTest) Check(value any) bool {
	if t.re == nil {
		panic("match test was not compiled; wrap it in ForStruct() or ForValue()")
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

func (t *matchTest) compile(_ any) error {
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
