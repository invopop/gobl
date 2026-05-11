package currency

import (
	"fmt"

	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/jsonschema"
)

// Code is the ISO currency code
type Code string

// CodeEmpty is used when there is no code.
const CodeEmpty Code = ""

// Validate ensures the currency code is valid according
// to the ISO 4217 three-letter list.
func (c Code) Validate() error {
	return inDefinitions(c)
}

// IsCodeDefined provides a rule validation that ensures the currency code is
// defined in the GOBL source.
var IsCodeDefined = is.Func("valid currency code", func(val any) bool {
	if c, ok := val.(Code); ok {
		return c == CodeEmpty || Get(c) != nil
	}
	return false
})

func inDefinitions(code Code) error {
	if code == CodeEmpty {
		return nil
	}
	if d := Get(code); d == nil {
		return fmt.Errorf("currency code %s not defined", code)
	}
	return nil
}

// Def provides the currency definition for the code.
func (c Code) Def() *Def {
	return Get(c)
}

// String provides the currency code as a string.
func (c Code) String() string {
	return string(c)
}

// In checks if the code is in the provided list of codes.
func (c Code) In(codes ...Code) bool {
	for _, v := range codes {
		if c == v {
			return true
		}
	}
	return false
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Code) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Currency Code",
		Type:        "string",
		OneOf:       make([]*jsonschema.Schema, len(Definitions())),
		Description: "Currency Code as defined in the GOBL source which is expected to be ISO or commonly used alternative.",
	}
	for i, v := range Definitions() {
		s.OneOf[i] = &jsonschema.Schema{
			Const: v.ISOCode,
			Title: v.Name,
		}
	}
	return s
}
