package org

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/uuid"
)

// Identity is used to define a code for a specific context.
type Identity struct {
	// Unique identity for this identity object.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Optional label useful for non-standard identities to give a bit more context.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// The type of Code being represented and usually specific for
	// a particular context, country, or tax regime.
	Type cbc.Code `json:"type,omitempty" jsonschema:"title=Type"`
	// The actual value of the identity code.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`
}

// Validate ensures the identity looks valid.
func (i *Identity) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Label),
		validation.Field(&i.Type),
		validation.Field(&i.Code, validation.Required),
	)
}
