package cbc

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// CodeDefinition describes a specific code and how it maps to a human name
// and description if appropriate. Regimes shouldn't typically do any additional
// conversion of codes, for that, regular keys should be used.
type CodeDefinition struct {
	// Code for which the definition is for.
	Code Code `json:"code" jsonschema:"title=Code"`
	// Short name for the code, if relevant.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`
	// Description offering more details about when the code should be used.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
	// Meta defines any additional details that may be useful or associated
	// with the code.
	Meta Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures the contents of the code definition are valid.
func (cd *CodeDefinition) Validate() error {
	return validation.ValidateStruct(cd,
		validation.Field(&cd.Code, validation.Required),
		validation.Field(&cd.Name, validation.Required),
		validation.Field(&cd.Desc),
	)
}
