package tax

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// KeyDefinition defines properties of a key that is specific for a regime.
type KeyDefinition struct {
	// Actual key value.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Short name for the key.
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Description offering more details about when the key should be used.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
	// Codes describes the list of codes that can be used alongside the Key,
	// for example with identities.
	Codes []*CodeDefinition `json:"codes,omitempty" jsonschema:"title=Codes"`
	// Keys is used instead of codes to define a further sub-set of keys that
	// can be used alongside this one.
	Keys []*KeyDefinition `json:"keys,omitempty" jsonschema:"title=Keys"`
	// Pattern is used to validate the key value instead of using a fixed value
	// from the code or key definitions.
	Pattern string `json:"pattern,omitempty" jsonschema:"title=Pattern"`
	// Map helps map local keys to specific codes, useful for converting the
	// described key into a local code.
	Map cbc.CodeMap `json:"map,omitempty" jsonschema:"title=Code Map"`
}

// CodeDefinition describes a specific code and how it maps to a human name
// and description if appropriate. Regimes shouldn't typically do any additional
// conversion of codes, for that, regular keys should be used.
type CodeDefinition struct {
	// Code for which the definition is for.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`
	// Short name for the code, if relevant.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`
	// Description offering more details about when the code should be used.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
	// Meta defines any additional details that may be useful or associated
	// with the code.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// HasCode loops through the key definitions codes and determines if there
// is a match.
func (kd *KeyDefinition) HasCode(code cbc.Code) bool {
	cd := kd.CodeDef(code)
	return cd != nil
}

// CodeDef returns the code definition for the provided code, or nil.
func (kd *KeyDefinition) CodeDef(code cbc.Code) *CodeDefinition {
	for _, c := range kd.Codes {
		if c.Code == code {
			return c
		}
	}
	return nil
}

// HasKey loops through the key definitions keys and determines if there
// is a match.
func (kd *KeyDefinition) HasKey(key cbc.Key) bool {
	skd := kd.KeyDef(key)
	return skd != nil
}

// KeyDef returns the key definition for the provided key, or nil.
func (kd *KeyDefinition) KeyDef(key cbc.Key) *KeyDefinition {
	for _, skd := range kd.Keys {
		if skd.Key == key {
			return skd
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
		validation.Field(&kd.Codes),
		validation.Field(&kd.Keys,
			validation.When(len(kd.Codes) > 0,
				validation.Empty,
			),
		),
		validation.Field(&kd.Pattern, validation.By(patternMustCompile)),
	)
	return err
}

// Validate ensures the contents of the code definition are valid.
func (cd *CodeDefinition) Validate() error {
	return validation.ValidateStruct(cd,
		validation.Field(&cd.Code, validation.Required),
		validation.Field(&cd.Name, validation.Required),
		validation.Field(&cd.Desc),
	)
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

func patternMustCompile(value any) error {
	pattern, ok := value.(string)
	if !ok || pattern == "" {
		return nil
	}
	_, err := regexp.Compile(pattern)
	return err
}
