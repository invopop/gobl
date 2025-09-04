package org

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
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

// Validate ensures email address looks valid.
func (e *Email) Validate() error {
	return e.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures email address looks valid inside the provided context.
func (e *Email) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, e,
		validation.Field(&e.Address, validation.Required, is.EmailFormat),
	)
}
