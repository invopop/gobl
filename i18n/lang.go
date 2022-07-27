package i18n

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/jsonschema"
)

// To simplify language management GoBL does not support full localization
// and instead focusses on simple multi-language based on the ISO 639-1 set
// of two letter codes. For business documents, this is sufficient as they
// are generally issued in a given country context.

// Lang represents the two letter language code.
type Lang string

// LangDef serves to handle language definitions
type LangDef struct {
	// Language Code
	Code Lang `json:"code" jsonschema:"title=Code"`
	// English name of the language
	Name string `json:"name" jsonschema:"title=Name"`
}

var isLang = validation.In(validLang()...)

func validLang() []interface{} {
	list := make([]interface{}, len(LangDefinitions))
	for i, d := range LangDefinitions {
		list[i] = string(d.Code)
	}
	return list
}

// Validate ensures the language code is valid according
// to the ISO 639-1 two-letter list.
func (l Lang) Validate() error {
	return validation.Validate(string(l), isLang)
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Lang) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Language Code",
		OneOf:       make([]*jsonschema.Schema, len(LangDefinitions)),
		Description: "Identifies the ISO639-1 language code",
	}
	for i, v := range LangDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const:       v.Code,
			Description: v.Name,
		}
	}
	return s
}
