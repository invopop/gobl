package cbc

import (
	"regexp"

	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// Definition defines properties of a key, code, or other value that has a specific meaning or
// utility.
type Definition struct {
	// Key being defined.
	Key Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Code this definition represents.
	Code Code `json:"code,omitempty" jsonschema:"title=Code"`

	// Short name for the key.
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Description offering more details about when the key should be used.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
	// Meta defines any additional details that may be useful or associated
	// with the key.
	Meta Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	// Where the information was sourced from.
	Sources []*Source `json:"sources,omitempty" jsonschema:"title=Sources"`

	// Values defines the possible values associated with the key, which themselves will
	// either be keys or codes depending on the context.
	Values []*Definition `json:"values,omitempty" jsonschema:"title=Values"`

	// Pattern is used to validate the key value instead of using a fixed value
	// from the code or key definitions.
	Pattern string `json:"pattern,omitempty" jsonschema:"title=Pattern"`

	// Map helps map local keys to specific codes, useful for converting the
	// described key into a local code.
	Map CodeMap `json:"map,omitempty" jsonschema:"title=Code Map"`
}

// Validate ensures the definition looks correct.
func (d *Definition) Validate() error {
	err := validation.ValidateStruct(d,
		validation.Field(&d.Key,
			validation.When(
				d.Code == "",
				validation.Required.Error("or code are required"),
			),
		),
		validation.Field(&d.Code,
			validation.When(
				d.Key != "",
				validation.Empty.Error("must be empty when key is set"),
			),
		),
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.Desc),
		validation.Field(&d.Sources),
		validation.Field(&d.Values),
		validation.Field(&d.Pattern, validation.By(validRegexpPattern)),
	)
	return err
}

// HasCode loops through values and determines if there
// is a match for the code.
func (d *Definition) HasCode(code Code) bool {
	cd := d.CodeDef(code)
	return cd != nil
}

// CodeDef searches the list of values and provides the matching definition.
func (d *Definition) CodeDef(code Code) *Definition {
	for _, c := range d.Values {
		if c.Code == code {
			return c
		}
	}
	return nil
}

// HasKey loops through values and determines if there
// is a match for the key.
func (d *Definition) HasKey(key Key) bool {
	cd := d.KeyDef(key)
	return cd != nil
}

// KeyDef searches the list of values and provides the matching definition.
func (d *Definition) KeyDef(key Key) *Definition {
	for _, c := range d.Values {
		if c.Key == key {
			return c
		}
	}
	return nil
}

// DefinitionKeys helps extract the keys from a list of key definitions.
func DefinitionKeys(list []*Definition) []Key {
	keys := make([]Key, 0, len(list))
	for _, item := range list {
		if item.Key != KeyEmpty {
			keys = append(keys, item.Key)
		}
	}
	return keys
}

// DefinitionCodes helps extract the codes from a list of key definitions.
func DefinitionCodes(list []*Definition) []Code {
	codes := make([]Code, 0, len(list))
	for _, item := range list {
		if item.Code != CodeEmpty {
			codes = append(codes, item.Code)
		}
	}
	return codes
}

// GetKeyDefinition helps fetch the key definition instance from a list.
func GetKeyDefinition(key Key, list []*Definition) *Definition {
	for _, item := range list {
		if item.Key == key {
			return item
		}
	}
	return nil
}

// GetCodeDefinition helps fetch the code definition instance from a list.
func GetCodeDefinition(code Code, list []*Definition) *Definition {
	for _, item := range list {
		if item.Code == code {
			return item
		}
	}
	return nil
}

// InKeyDefs prepares a validation to provide a rule that will determine
// if the keys are in the provided set.
func InKeyDefs(list []*Definition) validation.Rule {
	defs := make([]interface{}, len(list))
	for i, item := range list {
		defs[i] = item.Key
	}
	return validation.In(defs...)
}

// InCodeDefs prepares a validation to provide a rule that will determine
// if the codes are in the provided set.
func InCodeDefs(list []*Definition) validation.Rule {
	defs := make([]interface{}, len(list))
	for i, item := range list {
		defs[i] = item.Code
	}
	return validation.In(defs...)
}

func validRegexpPattern(value any) error {
	pattern, ok := value.(string)
	if !ok || pattern == "" {
		return nil
	}
	_, err := regexp.Compile(pattern)
	return err
}
