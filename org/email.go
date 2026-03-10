package org

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
)

// Email describes the electronic mailing details.
type Email struct {
	uuid.Identify
	// Identifier for the email.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Electronic mailing address.
	Address string `json:"addr" jsonschema:"title=Address"`
	// Additional fields.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Normalize will try to clean the email object.
func (e *Email) Normalize() {
	if e == nil {
		return
	}
	uuid.Normalize(&e.UUID)
	e.Label = cbc.NormalizeString(e.Label)
	e.Address = cbc.NormalizeString(e.Address)
}

func emailRules() *rules.Set {
	e := new(Email)
	return rules.ForStruct(e,
		rules.Field(&e.Address,
			rules.Assert("010", "expected a valid email address",
				rules.Required,
				is.EmailFormat,
			),
		),
	)
}
