package cbc

import (
	"errors"
	"regexp"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Code represents a string used to uniquely identify the data we're looking
// at. We use "code" instead of "id", to reenforce the fact that codes should
// be more easily set and used by humans within definitions than IDs or UUIDs.
// Codes are standardised so that when validated they must contain between
// 1 and 32 inclusive upper-case letters or numbers with optional periods (`.`),
// dashes (`-`), or forward slashes (`/`) to separate blocks.
type Code string

// CodeMap is a map of keys to specific codes, useful to determine regime specific
// codes from their key counterparts.
type CodeMap map[Key]Code

// Basic code constants.
var (
	CodePattern              = `^[A-Z0-9]+([\.\-\/]?[A-Z0-9]+)*$`
	CodePatternRegexp        = regexp.MustCompile(CodePattern)
	CodeMinLength     uint64 = 1
	CodeMaxLength     uint64 = 32
)

var (
	codeUnderscoreOrSpaceRegexp = regexp.MustCompile(`[_ ]`)
	codeInvalidCharsRegexp      = regexp.MustCompile(`[^A-Z0-9\.\-\/]`)
)

// CodeEmpty is used when no code is defined.
const CodeEmpty Code = ""

// NormalizeCode attempts to clean and normalize the provided code so that
// it matches what we'd expect instead of raising validation errors.
func NormalizeCode(c Code) Code {
	code := strings.ToUpper(c.String())
	code = strings.TrimSpace(code)
	code = codeUnderscoreOrSpaceRegexp.ReplaceAllString(code, "-")
	code = codeInvalidCharsRegexp.ReplaceAllString(code, "")
	return Code(code)
}

// Validate ensures that the code complies with the expected rules.
func (c Code) Validate() error {
	return validation.Validate(string(c),
		validation.Length(1, int(CodeMaxLength)),
		validation.Match(CodePatternRegexp),
	)
}

// IsEmpty returns true if no code is specified.
func (c Code) IsEmpty() bool {
	return c == CodeEmpty
}

// String returns string representation of code.
func (c Code) String() string {
	return string(c)
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

// JSONSchema provides a representation of the struct for usage in Schema.
func (Code) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Pattern:     CodePattern,
		Title:       "Code",
		MinLength:   &CodeMinLength,
		MaxLength:   &CodeMaxLength,
		Description: "Alphanumerical text identifier with upper-case letters, no whitespace, nor symbols.",
	}
}

// Validate ensures the code maps data looks correct.
func (cs CodeMap) Validate() error {
	err := make(validation.Errors)
	// values are already tested
	for k := range cs {
		if e := k.Validate(); e != nil {
			err[k.String()] = e
		}
	}
	if len(err) == 0 {
		return nil
	}
	return err
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
func CodeMapHas(keys ...Key) validation.Rule {
	return validateCodeMap{keys: keys}
}

type validateCodeMap struct {
	keys []Key
}

func (v validateCodeMap) Validate(value interface{}) error {
	cs, ok := value.(CodeMap)
	if !ok {
		return nil
	}
	var err validation.Errors
	for _, k := range v.keys {
		if _, ok := cs[k]; !ok {
			if err == nil {
				err = make(validation.Errors)
			}
			err[k.String()] = errors.New("required")
		}
	}
	if len(err) > 0 {
		return err
	}
	return nil
}

// JSONSchemaExtend ensures the pattern property is set correctly.
func (CodeMap) JSONSchemaExtend(schema *jsonschema.Schema) {
	prop := schema.AdditionalProperties
	schema.AdditionalProperties = nil
	schema.PatternProperties = map[string]*jsonschema.Schema{
		KeyPattern: prop,
	}
}
