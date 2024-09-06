package cbc

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// ValueDefinition describes a specific value and how it maps to a human name
// and description if appropriate.
type ValueDefinition struct {
	// Value for which the definition is for.
	Value string `json:"value" jsonschema:"title=Value"`
	// Short name for the value, if relevant.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`
	// Description offering more details about when the value should be used.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
	// Meta defines any additional details that may be useful or associated
	// with the value.
	Meta Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures the contents of the value definition are valid.
func (cd *ValueDefinition) Validate() error {
	return validation.ValidateStruct(cd,
		validation.Field(&cd.Value, validation.Required),
		validation.Field(&cd.Name, validation.Required),
		validation.Field(&cd.Desc),
	)
}
