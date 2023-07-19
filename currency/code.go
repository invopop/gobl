package currency

import (
	"github.com/invopop/gobl/num"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Code is the ISO currency code
type Code string

// CodeEmpty is used when there is no code.
const CodeEmpty Code = ""

// Def provides a structure for the currencies
type Def struct {
	Code  Code   `json:"code"`  // three-letter currency code
	Name  string `json:"name"`  // name of the currency
	Num   string `json:"num"`   // three-digit currency code
	Units uint32 `json:"units"` // how many cents are used for the currency
}

func validCodes() []interface{} {
	list := make([]interface{}, len(CodeDefinitions))
	for i, d := range CodeDefinitions {
		list[i] = string(d.Code)
	}
	return list
}

var isValidCode = validation.In(validCodes()...)

// Validate ensures the currency code is valid according
// to the ISO 4217 three-letter list.
func (c Code) Validate() error {
	return validation.Validate(string(c), isValidCode)
}

// Def provides the currency definition for the code.
func (c Code) Def() Def {
	d, _ := Get(c)
	return d
}

// Get provides the code's currency definition, or
// false if none is found.
func Get(c Code) (Def, bool) {
	for _, d := range CodeDefinitions {
		if d.Code == c {
			return d, true
		}
	}
	return Def{}, false
}

// Zero provides the currency's zero amount which is pre-set with the
// minimum precision for the currency.
func (d Def) Zero() num.Amount {
	return num.MakeAmount(0, d.Units)
}

// BaseAmount provides a definition's zero amount with the correct decimal
// places so that it can be used as a base for calculating totals.
//
// Deprecated: please now use the Zero method instead.
func (d Def) BaseAmount() num.Amount {
	return num.MakeAmount(0, d.Units)
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Code) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Currency Code",
		Type:        "string",
		OneOf:       make([]*jsonschema.Schema, len(CodeDefinitions)),
		Description: "ISO Currency Code",
	}
	for i, v := range CodeDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const: v.Code,
			Title: v.Name,
		}
	}
	return s
}
