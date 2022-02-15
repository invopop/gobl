package i18n

import "github.com/invopop/jsonschema"

const (
	defaultLanguage = EN
)

// String provides a simple map of locales to texts.
type String map[Lang]string

// String provides a single string from the map using the
// language requested or resorting to the default.
func (s String) String(lang Lang) string {
	if v, ok := s[lang]; ok {
		return v
	}
	if v, ok := s[defaultLanguage]; ok {
		return v
	}
	for _, v := range s {
		return v // provide first entry
	}
	return ""
}

// JSONSchema returns the json schema defintion
func (String) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		PatternProperties: map[string]*jsonschema.Schema{
			`^[a-z]{2}$`: {
				Type:  "string",
				Title: "Text in given language.",
			},
		},
		Title:       "Multi-language String",
		Description: "Map of 2-Letter language codes to their translations.",
	}
}
