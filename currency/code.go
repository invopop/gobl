package currency

import (
	"fmt"

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
