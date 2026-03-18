package cbc

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/jsonschema"
)

const (
	// DefaultCodeSeparator is the default separator used to join codes.
	DefaultCodeSeparator Code = "-"
)

// Code represents a string used to uniquely identify the data we're looking
// at. We use "code" instead of "id", to re-enforce the fact that codes should
// be more easily set and used by humans within definitions than IDs or UUIDs.
// Codes are standardized so that when validated they must contain between
// 1 and 64 inclusive english alphabet letters or numbers with optional
// periods (`.`), dashes (`-`), underscores (`_`), forward slashes (`/`),
// colons (`:`), commas (`,`), ampersand (`&`), or spaces (` `) to separate blocks.
// Each block must only be separated by a single symbol.
//
// The objective is to have a code that is easy to read and understand, while
// still being unique and easy to validate.
type Code string

// CodeMap is a map of keys to specific codes, useful to determine regime specific
// codes from their key counterparts.
type CodeMap map[Key]Code

// Basic code constants.
var (
	CodeSeparators           = `\.\-\:/,_\& `
	CodeDigits               = `A-ZÑa-z0-9`
	CodePattern              = `^[` + CodeDigits + `]+([` + CodeSeparators + `]?[` + CodeDigits + `]+)*$`
	CodePatternRegexp        = regexp.MustCompile(CodePattern)
	CodeMinLength     uint64 = 1
	CodeMaxLength     uint64 = 64
)

var (
	codeSeparatorRegexp         = regexp.MustCompile(`([` + CodeSeparators + `])[^` + CodeDigits + `]+`)
	codeInvalidCharsRegexp      = regexp.MustCompile(`[^` + CodeDigits + CodeSeparators + `]+`)
	codeNonAlphanumericalRegexp = regexp.MustCompile(`[^A-Z\d]`)
	codeNonNumericalRegexp      = regexp.MustCompile(`[^\d]`)
)

// CodeEmpty is used when no code is defined.
const CodeEmpty Code = ""

// NormalizeCode attempts to clean and normalize the provided code so that
// it matches what we'd expect instead of raising validation errors.
func NormalizeCode(c Code) Code {
	code := c.String()
	code = strings.TrimSpace(code)
	code = codeSeparatorRegexp.ReplaceAllString(code, "$1")
	code = codeInvalidCharsRegexp.ReplaceAllString(code, "")
	return Code(code)
}

// NormalizeAlphanumericalCode cleans and normalizes the code,
// ensuring all letters are uppercase while also removing
// non-alphanumerical characters.
func NormalizeAlphanumericalCode(c Code) Code {
	code := NormalizeCode(c).String()
	code = strings.ToUpper(code)
	code = codeNonAlphanumericalRegexp.ReplaceAllString(code, "")
	return Code(code)
}

// NormalizeNumericalCode cleans and normalizes the code, while also
// removing non-numerical characters.
func NormalizeNumericalCode(c Code) Code {
	code := NormalizeCode(c).String()
	code = codeNonNumericalRegexp.ReplaceAllString(code, "")
	return Code(code)
}

func codeRules() *rules.Set {
	return rules.For(Code(""),
		rules.Assert("01", fmt.Sprintf("codes must be no longer than %d characters", CodeMaxLength),
			rules.Length(0, int(CodeMaxLength)),
		),
		rules.Assert("02", "codes must only contain letters, numbers, and optionally separated by .-:/,_& or space",
			rules.Matches(CodePattern),
		),
	)
}

// Validate ensures that the code complies with the expected rules.
func (c Code) Validate() error {
	return rules.Validate(c)
}

// IsEmpty returns true if no code is specified.
func (c Code) IsEmpty() bool {
	return c == CodeEmpty
}

// String returns string representation of code.
func (c Code) String() string {
	return string(c)
}

// CodeStrings is a convenience method to convert a list of codes
// into a list of strings.
func CodeStrings(codes []Code) []string {
	l := make([]string, len(codes))
	for i, v := range codes {
		l[i] = v.String()
	}
	return l
}

// In returns true if the code's value matches one of those
// in the provided list.
func (c Code) In(ary ...Code) bool {
	for _, v := range ary {
		if v == c {
			return true
		}
	}
	return false
}

// Join returns a new code that is the result of joining the provided
// code with the current one using a default separator.
func (c Code) Join(c2 Code) Code {
	return c.JoinWith(DefaultCodeSeparator, c2)
}

// JoinWith returns a new code that is the result of joining the provided
// code with the current one using the provided separator. If any of the codes
// are empty, no separator will be added.
func (c Code) JoinWith(separator Code, c2 Code) Code {
	if c == CodeEmpty {
		return c2
	}
	if c2 == CodeEmpty {
		return c
	}
	return c + separator + c2
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Code) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:      "string",
		Pattern:   CodePattern,
		Title:     "Code",
		MinLength: &CodeMinLength,
		MaxLength: &CodeMaxLength,
		Description: here.Doc(`
			Alphanumerical text identifier with upper-case letters and limits on using
			special characters or whitespace to separate blocks.
		`),
	}
}

// InCodes provides a rules test that checks if a code's value is one of the provided codes.
func InCodes(codes ...Code) rules.Test {
	return rules.By("code in ["+strings.Join(CodeStrings(codes), ", ")+"]",
		func(value any) bool {
			c, ok := value.(Code)
			if !ok {
				return false
			}
			return c.In(codes...)
		},
	)
}

func codeMapRules() *rules.Set {
	return rules.For(CodeMap{},
		rules.Assert("01", "all code map keys must be valid",
			rules.By("valid keys", codeMapKeysValid),
		),
	)
}

func codeMapKeysValid(v any) bool {
	m, ok := v.(CodeMap)
	if !ok {
		return true
	}
	for k := range m {
		if rules.Validate(k) != nil {
			return false
		}
	}
	return true
}

// Has returns true if the code map has values for all the provided keys.
func (cs CodeMap) Has(keys ...Key) bool {
	for _, k := range keys {
		if _, ok := cs[k]; !ok {
			return false
		}
	}
	return true
}

// Equals returns true if the code map has the same keys and values as the provided
// map.
func (cs CodeMap) Equals(other CodeMap) bool {
	if len(cs) != len(other) {
		return false
	}
	for k, v := range cs {
		v2, ok := other[k]
		if !ok {
			return false
		}
		if v2 != v {
			return false
		}
	}
	return true
}

// CodeMapHas returns a validation rule that ensures the code set contains
// the provided keys.
func CodeMapHas(keys ...Key) rules.Test {
	return codeMapTest{keys: keys}
}

type codeMapTest struct {
	keys []Key
}

// String returns a string representation of the rule.
func (r codeMapTest) String() string {
	return fmt.Sprintf("have keys [%s]", strings.Join(KeyStrings(r.keys), ", "))
}

// Check returns true if the code map has all the required keys.
func (r codeMapTest) Check(value any) bool {
	cs, ok := value.(CodeMap)
	if !ok {
		return false
	}
	for _, k := range r.keys {
		if _, ok := cs[k]; !ok {
			return false
		}
	}
	return true
}

// JSONSchemaExtend ensures the pattern property is set correctly.
func (CodeMap) JSONSchemaExtend(schema *jsonschema.Schema) {
	prop := schema.AdditionalProperties
	schema.AdditionalProperties = nil
	schema.PatternProperties = map[string]*jsonschema.Schema{
		KeyPattern: prop,
	}
}
