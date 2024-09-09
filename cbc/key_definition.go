package cbc

import (
	"regexp"

	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// KeyDefinition defines properties of a key that is specific for a regime.
type KeyDefinition struct {
	// Actual key value.
	Key Key `json:"key" jsonschema:"title=Key"`
	// Short name for the key.
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Description offering more details about when the key should be used.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
	// Meta defines any additional details that may be useful or associated
	// with the key.
	Meta Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	// Values defines the possible values associated with the key.
	Values []*ValueDefinition `json:"values,omitempty" jsonschema:"title=Values"`

	// Pattern is used to validate the key value instead of using a fixed value
	// from the code or key definitions.
	Pattern string `json:"pattern,omitempty" jsonschema:"title=Pattern"`

	// Map helps map local keys to specific codes, useful for converting the
	// described key into a local code.
	Map CodeMap `json:"map,omitempty" jsonschema:"title=Code Map"`
}

// HasValue loops through values and determines if there
// is a match.
func (kd *KeyDefinition) HasValue(val string) bool {
	cd := kd.ValueDef(val)
	return cd != nil
}

// CodeDef returns the code definition for the provided code, or nil.
func (kd *KeyDefinition) ValueDef(val string) *ValueDefinition {
	for _, c := range kd.Values {
		if c.Value == val {
			return c
		}
	}
	return nil
}

// Validate ensures the key definition looks correct in the context of the regime.
func (kd *KeyDefinition) Validate() error {
	err := validation.ValidateStruct(kd,
		validation.Field(&kd.Key, validation.Required),
		validation.Field(&kd.Name, validation.Required),
		validation.Field(&kd.Desc),
		validation.Field(&kd.Values),
		validation.Field(&kd.Pattern, validation.By(validRegexpPattern)),
	)
	return err
}

// DefinitionKeys helps extract the keys from a list of key definitions.
func DefinitionKeys(list []*KeyDefinition) []Key {
	keys := make([]Key, len(list))
	for i, item := range list {
		keys[i] = item.Key
	}
	return keys
}

// GetKeyDefinition helps fetch the key definition instance from a list.
func GetKeyDefinition(key Key, list []*KeyDefinition) *KeyDefinition {
	for _, item := range list {
		if item.Key == key {
			return item
		}
	}
	return nil
}

// InKeyDefs prepares a validation to provide a rule that will determine
// if the keys are in the provided set.
func InKeyDefs(list []*KeyDefinition) validation.Rule {
	defs := make([]interface{}, len(list))
	for i, item := range list {
		defs[i] = item.Key
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
